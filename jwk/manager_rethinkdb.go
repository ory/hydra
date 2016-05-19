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

type rethinkSchema struct {
	ID string `gorethink:"id"`
	jose.JsonWebKeySet
}

func (m *RethinkManager) publishAdd(set string, keys []jose.JsonWebKey) error {
	if err := m.Table.Get(set).Exec(m.Session); err == r.ErrEmptyResult {
		if m.Table.Get(set).Insert(&rethinkSchema{
			ID: set,
			JsonWebKeySet: jose.JsonWebKeySet{Keys:keys},
		}).Exec(m.Session); err != nil {
			return errors.New(err)
		}
		return nil
	}

	if err := m.Table.Get(set).Update(map[string]interface{}{
		"keys": r.Row.Field([]interface{}{"keys"}).Default([]interface{}{}).Append([]interface{}{keys}...),
	}).Exec(m.Session); err != nil {
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
		"keys": r.Row.Field([]interface{}{"keys"}).Merge(r.Row.Field("keys").Difference([]interface{}{keys}...)),
	}).Exec(m.Session); err != nil {
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
		var update map[string]*rethinkSchema
		defer connections.Close()
		for connections.Next(&update) {
			newVal := update["new_val"]
			oldVal := update["old_val"]
			m.Lock()
			if newVal == nil && oldVal != nil {
				delete(m.Keys, oldVal.ID)
			} else if newVal != nil && oldVal != nil {
				delete(m.Keys, oldVal.ID)
				m.Keys[newVal.ID] = newVal.JsonWebKeySet
			} else {
				m.Keys[newVal.ID] = newVal.JsonWebKeySet
			}
			m.Unlock()
		}
	}()

	return nil
}