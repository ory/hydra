package jwk

import (
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
	"sync"
)

type MemoryManager struct {
	Keys map[string]map[string]jose.JsonWebKey
	sync.RWMutex
}

func (m *MemoryManager) AddKey(set string, key *jose.JsonWebKey) error {
	m.Lock()
	defer m.Unlock()

	m.Keys[set][key.KeyID] = *key
	return nil
}

func (m *MemoryManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	m.Lock()
	defer m.Unlock()

	m.alloc(set)
	for _, key := range keys.Keys {
		m.Keys[set][key.KeyID] = key
	}
	return nil
}

func (m *MemoryManager) GetKey(set, kid string) (*jose.JsonWebKey, error) {
	m.Lock()
	defer m.Unlock()

	m.alloc(set)

	if _, found := m.Keys[set]; !found {
		return nil, errors.New(pkg.ErrNotFound)
	}

	k, found := m.Keys[set][kid]
	if !found || &k == nil {
		return nil, errors.New(pkg.ErrNotFound)
	}

	return &k, nil
}

func (m *MemoryManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	m.Lock()
	defer m.Unlock()

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

func (m *MemoryManager) DeleteKey(set, kid string) error {
	m.Lock()
	defer m.Unlock()

	m.alloc(set)
	delete(m.Keys[set], kid)
	return nil
}

func (m *MemoryManager) DeleteKeySet(set string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Keys, set)
	return nil
}

func (m *MemoryManager) alloc(set string) {
	if m.Keys == nil {
		m.Keys = make(map[string]map[string]jose.JsonWebKey)
	}
	if m.Keys[set] == nil {
		m.Keys[set] = make(map[string]jose.JsonWebKey)
	}
}
