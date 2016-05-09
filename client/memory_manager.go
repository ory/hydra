package client

import (
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/ory-am/fosite/hash"
)

type MemoryManager struct {
	Clients map[string]*fosite.DefaultClient
	Hasher  hash.Hasher
}

func (m *MemoryManager) GetClient(id string) (fosite.Client, error) {
	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}
	return c, nil
}

func (m *MemoryManager) Authenticate(id string, secret []byte) (*fosite.DefaultClient, error) {
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
	delete(m.Clients, id)
	return nil
}

func (m *MemoryManager) GetClients() (map[string]*fosite.DefaultClient, error) {
	return m.Clients, nil
}
