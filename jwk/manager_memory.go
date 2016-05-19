package jwk

import (
	"sync"

	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
)

type MemoryManager struct {
	Keys map[string]*jose.JsonWebKeySet
	sync.RWMutex
}

func (m *MemoryManager) AddKey(set string, key *jose.JsonWebKey) error {
	m.Lock()
	defer m.Unlock()

	m.alloc()
	if m.Keys[set] == nil {
		m.Keys[set] = &jose.JsonWebKeySet{Keys: []jose.JsonWebKey{}}
	}
	m.Keys[set].Keys = append(m.Keys[set].Keys, *key)
	return nil
}

func (m *MemoryManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	for _, key := range keys.Keys {
		m.AddKey(set, &key)
	}
	return nil
}

func (m *MemoryManager) GetKey(set, kid string) ([]jose.JsonWebKey, error) {
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

func (m *MemoryManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	m.Lock()
	defer m.Unlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.New(pkg.ErrNotFound)
	}

	return keys, nil
}

func (m *MemoryManager) DeleteKey(set, kid string) error {
	return nil
}

func (m *MemoryManager) DeleteKeySet(set string) error {
	return nil
}

func (m *MemoryManager) alloc() {
	if m.Keys == nil {
		m.Keys = make(map[string]*jose.JsonWebKeySet)
	}
}
