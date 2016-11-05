package jwk

import (
	"sync"

	"encoding/json"

	"time"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RethinkManager struct {
	Session *r.Session
	Table   r.Term
	sync.RWMutex

	Cipher *AEAD

	Keys map[string]jose.JsonWebKeySet
}

func (m *RethinkManager) SetUpIndex() error {
	if _, err := m.Table.IndexWait("kid").Run(m.Session); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *RethinkManager) AddKey(set string, key *jose.JsonWebKey) error {
	if err := m.publishAdd(set, []jose.JsonWebKey{*key}); err != nil {
		return err
	}
	return nil
}

func (m *RethinkManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	if err := m.publishAdd(set, keys.Keys); err != nil {
		return err
	}
	return nil
}

func (m *RethinkManager) GetKey(set, kid string) (*jose.JsonWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	result := keys.Key(kid)
	if len(result) == 0 {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	return &jose.JsonWebKeySet{
		Keys: result,
	}, nil
}

func (m *RethinkManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	if len(keys.Keys) == 0 {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	return &keys, nil
}

func (m *RethinkManager) DeleteKey(set, kid string) error {
	keys, err := m.GetKey(set, kid)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if err := m.publishDelete(set, keys.Keys); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *RethinkManager) DeleteKeySet(set string) error {
	if err := m.publishDeleteAll(set); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *RethinkManager) alloc() {
	if m.Keys == nil {
		m.Keys = make(map[string]jose.JsonWebKeySet)
	}
}

type rethinkSchema struct {
	KID string `gorethink:"kid"`
	Set string `gorethink:"set"`
	Key string `gorethink:"key"`
}

func (m *RethinkManager) publishAdd(set string, keys []jose.JsonWebKey) error {
	raws := make([]string, len(keys))
	for k, key := range keys {
		out, err := json.Marshal(key)
		if err != nil {
			return errors.Wrap(err, "")
		}
		encrypted, err := m.Cipher.Encrypt(out)
		if err != nil {
			return errors.Wrap(err, "")
		}
		raws[k] = encrypted
	}

	for k, raw := range raws {
		if _, err := m.Table.Insert(&rethinkSchema{
			KID: keys[k].KeyID,
			Set: set,
			Key: raw,
		}).RunWrite(m.Session); err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}
func (m *RethinkManager) publishDeleteAll(set string) error {
	if err := m.Table.Filter(map[string]interface{}{
		"set": set,
	}).Delete().Exec(m.Session); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *RethinkManager) publishDelete(set string, keys []jose.JsonWebKey) error {
	for _, key := range keys {
		if _, err := m.Table.Filter(map[string]interface{}{
			"kid": key.KeyID,
			"set": set,
		}).Delete().RunWrite(m.Session); err != nil {
			return errors.Wrap(err, "")
		}
	}
	return nil
}

func (m *RethinkManager) Watch(ctx context.Context) {
	go pkg.Retry(time.Second*15, time.Minute, func() error {
		connections, err := m.Table.Changes().Run(m.Session)
		if err != nil {
			return errors.Wrap(err, "")
		}
		defer connections.Close()

		var update map[string]*rethinkSchema
		for connections.Next(&update) {
			logrus.Debug("Received update from RethinkDB Cluster in JSON Web Key manager.")
			newVal := update["new_val"]
			oldVal := update["old_val"]
			if newVal == nil && oldVal != nil {
				m.watcherRemove(oldVal)
			} else if newVal != nil && oldVal != nil {
				m.watcherRemove(oldVal)
				m.watcherInsert(newVal)
			} else {
				m.watcherInsert(newVal)
			}
		}

		if connections.Err() != nil {
			err = errors.Wrap(connections.Err(), "")
			pkg.LogError(err)
			return err
		}
		return nil
	})
}

func (m *RethinkManager) watcherInsert(val *rethinkSchema) {
	var c jose.JsonWebKey
	key, err := m.Cipher.Decrypt(val.Key)
	if err != nil {
		pkg.LogError(errors.Wrap(err, ""))
		return
	}

	if err := json.Unmarshal(key, &c); err != nil {
		pkg.LogError(errors.Wrap(err, ""))
		return
	}

	m.Lock()
	defer m.Unlock()
	keys := m.Keys[val.Set]
	keys.Keys = append(keys.Keys, c)
	m.Keys[val.Set] = keys
}

func (m *RethinkManager) watcherRemove(val *rethinkSchema) {
	keys, ok := m.Keys[val.Set]
	if !ok {
		return
	}

	m.Lock()
	defer m.Unlock()
	keys.Keys = filter(keys.Keys, func(k jose.JsonWebKey) bool {
		return k.KeyID != val.KID
	})
	m.Keys[val.Set] = keys
}

func (m *RethinkManager) ColdStart() error {
	m.Keys = map[string]jose.JsonWebKeySet{}
	clients, err := m.Table.Run(m.Session)
	if err != nil {
		return errors.Wrap(err, "")
	}

	var raw *rethinkSchema
	var key jose.JsonWebKey
	m.Lock()
	defer m.Unlock()
	for clients.Next(&raw) {
		pt, err := m.Cipher.Decrypt(raw.Key)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not decrypt JSON Web Keys because: %s. This usually happens when a wrong system secret is being used", err.Error()))
		}

		if err := json.Unmarshal(pt, &key); err != nil {
			return errors.Wrap(err, "")
		}

		keys, ok := m.Keys[raw.Set]
		if !ok {
			keys = jose.JsonWebKeySet{}
		}
		keys.Keys = append(keys.Keys, key)
		m.Keys[raw.Set] = keys
	}

	return nil
}

func filter(vs []jose.JsonWebKey, f func(jose.JsonWebKey) bool) []jose.JsonWebKey {
	vsf := make([]jose.JsonWebKey, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
