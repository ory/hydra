package jwk

import (
	"sync"

	"github.com/go-errors/errors"
	r "github.com/dancannon/gorethink"
	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
)

type RethinkManager struct {
	Session     *r.Session
	Table       r.Term

	Keys map[string]map[string]jose.JsonWebKey

	sync.RWMutex
}

func (m *RethinkManager) AddKey(set string, key *jose.JsonWebKey) error {
	m.Lock()
	defer m.Unlock()

	m.alloc(set)
	m.Keys[set][key.KeyID] = *key
	return nil
}

func (m *RethinkManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	m.Lock()
	defer m.Unlock()

	m.alloc(set)
	for _, key := range keys.Keys {
		m.Keys[set][key.KeyID] = key
	}
	return nil
}

func (m *RethinkManager) GetKey(set, kid string) (*jose.JsonWebKey, error) {
	m.Lock()
	defer m.Unlock()

	m.alloc("")
	if _, found := m.Keys[set]; !found {
		return nil, errors.New(pkg.ErrNotFound)
	}

	k, found := m.Keys[set][kid]
	if !found || &k == nil {
		return nil, errors.New(pkg.ErrNotFound)
	}

	return &k, nil
}

func (m *RethinkManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	m.Lock()
	defer m.Unlock()

	m.alloc("")
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.New(pkg.ErrNotFound)
	}

	ks := []jose.JsonWebKey{}
	for _, key := range keys {
		ks = append(ks, key)
	}

	return &jose.JsonWebKeySet{
		Keys: ks,
	}, nil
}

func (m *RethinkManager) DeleteKey(set, kid string) error {
	m.Lock()
	defer m.Unlock()

	m.alloc(set)
	delete(m.Keys[set], kid)
	return nil
}

func (m *RethinkManager) DeleteKeySet(set string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Keys, set)
	return nil
}

func (m *RethinkManager) alloc(set string) {
	if m.Keys == nil {
		m.Keys = make(map[string]map[string]jose.JsonWebKey)
	}
	if set != "" && m.Keys[set] == nil {
		m.Keys[set] = make(map[string]jose.JsonWebKey)
	}
}

type rdbSchema struct {
	ID string `gorethink:"id"`
	Key *jose.JsonWebKey `gorethink:"key"`
}

func (m *RethinkManager) publishCreate(set string, key *jose.JsonWebKey) error {
	keys := map[string]*jose.JsonWebKey{}
	keys[set] = key
	if err := m.Table.Exec(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) publishDelete(id string) error {
	if err := m.Table.Get(id).Delete().Exec(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) Watch(ctx context.Context) error {
	connections, err := m.Table.Changes().Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	go func() {
		var update map[string]*Connection
		defer connections.Close()
		for connections.Next(&update) {
			newVal := update["new_val"]
			oldVal := update["old_val"]
			m.Lock()
			if newVal == nil && oldVal != nil {
				delete(m.Connections, oldVal.GetID())
			} else if newVal != nil && oldVal != nil {
				delete(m.Connections, oldVal.GetID())
				m.Connections[newVal.GetID()] = newVal
			} else {
				m.Connections[newVal.GetID()] = newVal
			}
			m.Unlock()
		}
	}()

	return nil
}