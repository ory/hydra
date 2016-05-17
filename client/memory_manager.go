package client

import (
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"sync"
)

type MemoryManager struct {
	Clients map[string]*fosite.DefaultClient
	Hasher  hash.Hasher
	sync.RWMutex
}

func (m *MemoryManager) GetClient(id string) (fosite.Client, error) {
	m.Lock()
	defer m.Unlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}
	return c, nil
}

func (m *MemoryManager) Authenticate(id string, secret []byte) (*fosite.DefaultClient, error) {
	m.Lock()
	defer m.Unlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.New(err)
	}

	return c, nil
}

func (m *MemoryManager) CreateClient(c *fosite.DefaultClient) error {
	m.Lock()
	defer m.Unlock()

	if c.ID == "" {
		c.ID = uuid.New()
	}

	hash, err := m.Hasher.Hash(c.Secret)
	if err != nil {
		return errors.New(err)
	}
	c.Secret = hash

	m.Clients[c.GetID()] = c
	return nil
}

func (m *MemoryManager) DeleteClient(id string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Clients, id)
	return nil
}

func (m *MemoryManager) GetClients() (map[string]*fosite.DefaultClient, error) {
	m.Lock()
	defer m.Unlock()

	return m.Clients, nil
}
