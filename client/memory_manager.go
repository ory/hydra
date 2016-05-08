package client

import (
	"bytes"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
)

type MemoryManager struct {
	Clients map[string]*fosite.DefaultClient
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

	if bytes.Compare(c.GetHashedSecret(), secret) != 0 {
		return nil, errors.New(pkg.ErrUnauthorized)
	}

	return c, nil
}

func (m *MemoryManager) CreateClient(c *fosite.DefaultClient) error {
	if c.ID == "" {
		c.ID = uuid.New()
	}

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
