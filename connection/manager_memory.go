package connection

import (
	"sync"

	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
)

type MemoryManager struct {
	Connections map[string]Connection
	sync.RWMutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		Connections: make(map[string]Connection),
	}
}

func (m *MemoryManager) Create(c *Connection) error {
	m.Lock()
	defer m.Unlock()

	m.Connections[c.GetID()] = *c
	return nil
}

func (m *MemoryManager) Delete(id string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Connections, id)
	return nil
}

func (m *MemoryManager) Get(id string) (*Connection, error) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Connections[id]
	if !ok {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}
	return &c, nil
}

func (m *MemoryManager) FindAllByLocalSubject(subject string) ([]Connection, error) {
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

func (m *MemoryManager) FindByRemoteSubject(provider, subject string) (*Connection, error) {
	m.RLock()
	defer m.RUnlock()

	for _, c := range m.Connections {
		if c.GetProvider() == provider && c.GetRemoteSubject() == subject {
			return &c, nil
		}
	}
	return nil, errors.Wrap(pkg.ErrNotFound, "")
}
