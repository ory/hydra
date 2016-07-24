package client

import (
	"sync"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
)

type MemoryManager struct {
	Clients map[string]Client
	Hasher  hash.Hasher
	sync.RWMutex
}

func (m *MemoryManager) GetClient(id string) (fosite.Client, error) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}
	return &c, nil
}

func (m *MemoryManager) Authenticate(id string, secret []byte) (*Client, error) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.New(err)
	}

	return &c, nil
}

func (m *MemoryManager) CreateClient(c *Client) error {
	m.Lock()
	defer m.Unlock()

	if c.ID == "" {
		c.ID = uuid.New()
	}

	hash, err := m.Hasher.Hash([]byte(c.Secret))
	if err != nil {
		return errors.New(err)
	}
	c.Secret = string(hash)

	m.Clients[c.GetID()] = *c
	return nil
}

func (m *MemoryManager) DeleteClient(id string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Clients, id)
	return nil
}

func (m *MemoryManager) GetClients() (clients map[string]Client, err error) {
	m.RLock()
	defer m.RUnlock()
	clients = make(map[string]Client)
	for _, c := range m.Clients {
		clients[c.ID] = c
	}

	return clients, nil
}
