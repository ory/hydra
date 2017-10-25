package jwk

import (
	"sync"

	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type MemoryManager struct {
	Keys map[string]*jose.JSONWebKeySet
	sync.RWMutex
}

func (m *MemoryManager) AddKey(set string, key *jose.JSONWebKey) error {
	m.Lock()
	defer m.Unlock()

	m.alloc()
	if m.Keys[set] == nil {
		m.Keys[set] = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	}
	m.Keys[set].Keys = append(m.Keys[set].Keys, *key)
	return nil
}

func (m *MemoryManager) AddKeySet(set string, keys *jose.JSONWebKeySet) error {
	for _, key := range keys.Keys {
		m.AddKey(set, &key)
	}
	return nil
}

func (m *MemoryManager) GetKey(set, kid string) (*jose.JSONWebKeySet, error) {
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

	return &jose.JSONWebKeySet{
		Keys: result,
	}, nil
}

func (m *MemoryManager) GetKeySet(set string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	return keys, nil
}

func (m *MemoryManager) DeleteKey(set, kid string) error {
	keys, err := m.GetKeySet(set)
	if err != nil {
		return err
	}

	m.Lock()
	var results []jose.JSONWebKey
	for _, key := range keys.Keys {
		if key.KeyID != kid {
			results = append(results)
		}
	}
	m.Keys[set].Keys = results
	defer m.Unlock()

	return nil
}

func (m *MemoryManager) DeleteKeySet(set string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Keys, set)
	return nil
}

func (m *MemoryManager) alloc() {
	if m.Keys == nil {
		m.Keys = make(map[string]*jose.JSONWebKeySet)
	}
}
