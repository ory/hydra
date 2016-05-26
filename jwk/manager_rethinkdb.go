package jwk

import (
	"sync"

	"encoding/json"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
	"golang.org/x/net/context"
)

type RethinkManager struct {
	Session *r.Session
	Table   r.Term

	Keys map[string]jose.JsonWebKeySet

	sync.RWMutex
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

func (m *RethinkManager) GetKey(set, kid string) ([]jose.JsonWebKey, error) {
	m.Lock()
	defer m.Unlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.New(pkg.ErrNotFound)
	}

	result := keys.Key(kid)
	if len(result) == 0 {
		return nil, errors.New(pkg.ErrNotFound)
	}
	return result, nil
}

func (m *RethinkManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	m.Lock()
	defer m.Unlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.New(pkg.ErrNotFound)
	}

	return &keys, nil
}

func (m *RethinkManager) DeleteKey(set, kid string) error {
	keys, err := m.GetKey(set, kid)
	if err != nil {
		return errors.New(err)
	}

	if err := m.publishDelete(set, keys); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) DeleteKeySet(set string) error {
	if err := m.publishDeleteAll(set); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) alloc() {
	if m.Keys == nil {
		m.Keys = make(map[string]jose.JsonWebKeySet)
	}
}

type rethinkSchema struct {
	ID   string            `gorethink:"id"`
	Keys []json.RawMessage `gorethink:"keys"`
}

func (m *RethinkManager) publishAdd(set string, keys []jose.JsonWebKey) error {
	raws := make([]json.RawMessage, len(keys))
	for k, key := range keys {
		out, err := json.Marshal(key)
		if err != nil {
			return errors.New(err)
		}
		raws[k] = out
	}

	if _, err := m.GetKeySet(set); errors.Is(err, pkg.ErrNotFound) {
		if _, err := m.Table.Insert(&rethinkSchema{
			ID:   set,
			Keys: raws,
		}).RunWrite(m.Session); err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	if _, err := m.Table.Get(set).Update(map[string]interface{}{
		"keys": r.Row.Field("keys").Append([]interface{}{keys}...),
	}).RunWrite(m.Session); err != nil {
		return errors.New(err)
	}

	return nil
}
func (m *RethinkManager) publishDeleteAll(set string) error {
	if err := m.Table.Get(set).Delete().Exec(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) publishDelete(set string, keys []jose.JsonWebKey) error {
	if err := m.Table.Get(set).Update(map[string]interface{}{
		"keys": r.Row.Field("keys").SetDifference([]interface{}{keys}...),
	}).Exec(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func rawToKeys(raws []json.RawMessage) []jose.JsonWebKey {
	fmt.Printf("%s", raws)

	var keys = make([]jose.JsonWebKey, len(raws))
	var key = new(jose.JsonWebKey)
	for k, raw := range raws {
		err := key.UnmarshalJSON(raw)
		if err != nil {
			panic(err.Error())
		}
		keys[k] = *key
	}
	return keys

}

func (m *RethinkManager) Watch(ctx context.Context) error {
	connections, err := m.Table.Changes().Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	go func() {
		var update map[string]*rethinkSchema
		defer connections.Close()
		for connections.Next(&update) {
			newVal := update["new_val"]
			oldVal := update["old_val"]
			m.Lock()
			fmt.Printf("\n\nGot new data: %v\n\n", update)
			if newVal == nil && oldVal != nil {
				delete(m.Keys, oldVal.ID)
			} else if newVal != nil && oldVal != nil {
				delete(m.Keys, oldVal.ID)
				m.Keys[newVal.ID] = jose.JsonWebKeySet{
					Keys: rawToKeys(newVal.Keys),
				}
			} else {
				m.Keys[newVal.ID] = jose.JsonWebKeySet{
					Keys: rawToKeys(newVal.Keys),
				}
			}
			m.Unlock()
		}
	}()

	return nil
}
