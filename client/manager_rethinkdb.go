package client

import (
	r "github.com/dancannon/gorethink"
	"sync"
	"github.com/ory-am/fosite"
	"github.com/go-errors/errors"
	"golang.org/x/net/context"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/fosite/hash"
	"github.com/pborman/uuid"
)

type RethinkManager struct {
	Session *r.Session
	Table   r.Term

	Clients map[string]*fosite.DefaultClient
	Hasher  hash.Hasher

	sync.RWMutex
}

func (m *RethinkManager) GetClient(id string) (fosite.Client, error) {
	m.Lock()
	defer m.Unlock()

	c, ok := m.Clients[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}
	return c, nil
}

func (m *RethinkManager) Authenticate(id string, secret []byte) (*fosite.DefaultClient, error) {
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

func (m *RethinkManager) CreateClient(c *fosite.DefaultClient) error {
	if c.ID == "" {
		c.ID = uuid.New()
	}

	hash, err := m.Hasher.Hash(c.Secret)
	if err != nil {
		return errors.New(err)
	}
	c.Secret = hash

	if err := m.publishCreate(c); err != nil {
		return err
	}

	return nil
}

func (m *RethinkManager) DeleteClient(id string) error {
	if err := m.publishDelete(id); err != nil {
		return err
	}

	return nil
}

func (m *RethinkManager) GetClients() (map[string]*fosite.DefaultClient, error) {
	m.Lock()
	defer m.Unlock()

	return m.Clients, nil
}

func (m *RethinkManager) ColdStart() error {
	m.Clients = map[string]*fosite.DefaultClient{}
	clients, err := m.Table.Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	var client fosite.DefaultClient
	m.Lock()
	defer m.Unlock()
	for clients.Next(&client) {
		m.Clients[client.ID] = &client
	}

	return nil
}

func (m *RethinkManager) publishCreate(client *fosite.DefaultClient) error {
	if _, err := m.Table.Insert(client).RunWrite(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) publishDelete(id string) error {
	if _, err := m.Table.Get(id).Delete().RunWrite(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) Watch(ctx context.Context) error {
	clients, err := m.Table.Changes().Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	go func() {
		for {
			var update map[string]*fosite.DefaultClient
			for clients.Next(&update) {
				newVal := update["new_val"]
				oldVal := update["old_val"]
				m.Lock()
				if newVal == nil && oldVal != nil {
					delete(m.Clients, oldVal.GetID())
				} else if newVal != nil && oldVal != nil {
					delete(m.Clients, oldVal.GetID())
					m.Clients[newVal.GetID()] = newVal
				} else {
					m.Clients[newVal.GetID()] = newVal
				}
				m.Unlock()

				clients.Close()
				if clients.Err() != nil {
					pkg.LogError(errors.New(clients.Err()))
				}

				clients, err = m.Table.Changes().Run(m.Session)
				if err != nil {
					pkg.LogError(errors.New(clients.Err()))
				}
			}
		}
	}()

	return nil
}