package client

import (
	"sync"

	"github.com/imdario/mergo"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type MemoryManager struct {
	Clients map[string]Client
	Hasher  fosite.Hasher
	sync.RWMutex
}

func (m *MemoryManager) GetConcreteClient(id string) (*Client, error) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}
	return &c, nil
}

func (m *MemoryManager) GetClient(id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *MemoryManager) UpdateClient(c *Client) error {
	o, err := m.GetClient(c.ID)
	if err != nil {
		return err
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(c.Secret))
		if err != nil {
			return errors.Wrap(err, "")
		}
		c.Secret = string(h)
	}
	if err := mergo.Merge(c, o); err != nil {
		return errors.Wrap(err, "")
	}

	m.Clients[c.GetID()] = *c
	return nil
}

func (m *MemoryManager) Authenticate(id string, secret []byte) (*Client, error) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.Wrap(err, "")
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
		return errors.Wrap(err, "")
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
