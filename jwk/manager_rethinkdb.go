package jwk

import (
	"sync"

	"github.com/go-errors/errors"
	r "github.com/dancannon/gorethink"
	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
	"golang.org/x/net/context"
)

type RethinkManager struct {
	Session *r.Session
	Table   r.Term

	Keys    map[string]jose.JsonWebKeySet

	sync.RWMutex
}

func (m *RethinkManager) AddKey(set string, key *jose.JsonWebKey) error {
	return nil
}

func (m *RethinkManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
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

	return keys.Key(kid), nil
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

	return nil
}

func (m *RethinkManager) DeleteKeySet(set string) error {
	return nil
}

func (m *RethinkManager) alloc() {
	if m.Keys == nil {
		m.Keys = make(map[string]jose.JsonWebKeySet)
	}
}

type rdbSchema struct {
	ID  string `gorethink:"id"`
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
	_, err := m.Table.Changes().Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	return nil
}