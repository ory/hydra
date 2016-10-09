package connection

import (
	"sync"

	r "gopkg.in/dancannon/gorethink.v2"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type RethinkManager struct {
	Session *r.Session
	Table   r.Term

	Connections map[string]Connection

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
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Connections[id]
	if !ok {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}
	return &c, nil
}

func (m *RethinkManager) FindAllByLocalSubject(subject string) ([]Connection, error) {
	m.RLock()
	defer m.RUnlock()

	var cs []Connection
	for _, c := range m.Connections {
		if c.GetLocalSubject() == subject {
			cs = append(cs, c)
		}
	}
	return cs, nil
}

func (m *RethinkManager) FindByRemoteSubject(provider, subject string) (*Connection, error) {
	m.RLock()
	defer m.RUnlock()

	for _, c := range m.Connections {
		if c.GetProvider() == provider && c.GetRemoteSubject() == subject {
			return &c, nil
		}
	}
	return nil, errors.Wrap(pkg.ErrNotFound, "")
}

func (m *RethinkManager) ColdStart() error {
	m.Connections = map[string]Connection{}
	clients, err := m.Table.Run(m.Session)
	if err != nil {
		return errors.Wrap(err, "")
	}

	var connection Connection
	m.Lock()
	defer m.Unlock()
	for clients.Next(&connection) {
		m.Connections[connection.ID] = connection
	}

	return nil
}

func (m *RethinkManager) publishCreate(c *Connection) error {
	if err := m.Table.Insert(c).Exec(m.Session); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *RethinkManager) publishDelete(id string) error {
	if err := m.Table.Get(id).Delete().Exec(m.Session); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *RethinkManager) Watch(ctx context.Context) {
	go pkg.Retry(time.Second*15, time.Minute, func() error {
		connections, err := m.Table.Changes().Run(m.Session)
		if err != nil {
			return errors.Wrap(err, "")
		}
		defer connections.Close()

		var update map[string]*Connection
		for connections.Next(&update) {
			logrus.Debug("Received update in social connection manager.")
			newVal := update["new_val"]
			oldVal := update["old_val"]
			m.Lock()
			if newVal == nil && oldVal != nil {
				delete(m.Connections, oldVal.GetID())
			} else if newVal != nil && oldVal != nil {
				delete(m.Connections, oldVal.GetID())
				m.Connections[newVal.GetID()] = *newVal
			} else {
				m.Connections[newVal.GetID()] = *newVal
			}
			m.Unlock()
		}

		if connections.Err() != nil {
			err = errors.Wrap(connections.Err(), "")
			pkg.LogError(err)
			return err
		}
		return nil
	})
}
