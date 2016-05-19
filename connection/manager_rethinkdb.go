package connection

import (
	r "github.com/dancannon/gorethink"
	"sync"

	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
	"golang.org/x/net/context"
)

type RethinkManager struct {
	Session     *r.Session
	Table       r.Term

	Connections map[string]*Connection

	sync.RWMutex
}

func (m *RethinkManager) Create(c *Connection) error {
	if err := m.publishCreate(c); err != nil {
		return err
	}
	return nil
}

func (m *RethinkManager) Delete(id string) error {
	if err := m.publishDelete(id); err != nil {
		return err
	}

	return nil
}

func (m *RethinkManager) Get(id string) (*Connection, error) {
	m.Lock()
	defer m.Unlock()

	c, ok := m.Connections[id]
	if !ok {
		return nil, errors.New(pkg.ErrNotFound)
	}
	return c, nil
}

func (m *RethinkManager) FindAllByLocalSubject(subject string) ([]*Connection, error) {
	m.Lock()
	defer m.Unlock()

	var cs []*Connection
	for _, c := range m.Connections {
		if c.GetLocalSubject() == subject {
			cs = append(cs, c)
		}
	}
	return cs, nil
}

func (m *RethinkManager) FindByRemoteSubject(provider, subject string) (*Connection, error) {
	m.Lock()
	defer m.Unlock()

	for _, c := range m.Connections {
		if c.GetProvider() == provider && c.GetRemoteSubject() == subject {
			return c, nil
		}
	}
	return nil, errors.New(pkg.ErrNotFound)
}

func (m *RethinkManager) fetch() error {
	m.Connections = map[string]*Connection{}
	clients, err := m.Table.Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	var connection Connection
	m.Lock()
	defer m.Unlock()
	for clients.Next(&connection) {
		m.Connections[connection.ID] = &connection
	}

	return nil
}

func (m *RethinkManager) publishCreate(c *Connection) error {
	if err := m.Table.Insert(c).Exec(m.Session); err != nil {
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
	connections, err := m.Table.Changes().Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	go func() {
		var update map[string]*Connection
		defer connections.Close()
		for connections.Next(&update) {
			newVal := update["new_val"]
			oldVal := update["old_val"]
			m.Lock()
			if newVal == nil && oldVal != nil {
				delete(m.Connections, oldVal.GetID())
			} else if newVal != nil && oldVal != nil {
				delete(m.Connections, oldVal.GetID())
				m.Connections[newVal.GetID()] = newVal
			} else {
				m.Connections[newVal.GetID()] = newVal
			}
			m.Unlock()
		}
	}()

	return nil
}