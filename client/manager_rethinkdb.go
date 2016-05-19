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
	"fmt"
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

	if err := m.publishCreate(c); err != nil {
		return err
	}

	// m.Clients[c.GetID()] = c
	return nil
}

func (m *RethinkManager) DeleteClient(id string) error {
	m.Lock()
	defer m.Unlock()

	if err := m.publishDelete(id); err != nil {
		return err
	}

	// delete(m.Clients, id)
	return nil
}

func (m *RethinkManager) GetClients() (map[string]*fosite.DefaultClient, error) {
	m.Lock()
	defer m.Unlock()

	return m.Clients, nil
}

func (m *RethinkManager) fetch() error {
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
	if err := m.Table.Insert(client).Exec(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) publishDelete(id string) error {
	if err := m.Table.Get(id).Delete().Exec(m.Session); err != nil {
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
		var update map[string]*fosite.DefaultClient
		defer clients.Close()
		for clients.Next(&update) {
			newVal := update["new_val"]
			oldVal := update["old_val"]
			m.Lock()
			if newVal == nil && oldVal != nil {
				fmt.Println("delete")
				delete(m.Clients, oldVal.GetID())
			} else if newVal != nil && oldVal != nil {
				fmt.Println("update")
				delete(m.Clients, oldVal.GetID())
				m.Clients[newVal.GetID()] = newVal
			} else {
				fmt.Println("create")
				m.Clients[newVal.GetID()] = newVal
			}
			m.Unlock()
		}
	}()

	return nil
}