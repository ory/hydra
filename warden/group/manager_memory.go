package group

import (
	"sort"
	"sync"

	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		Groups: map[string]Group{},
	}
}

type MemoryManager struct {
	Groups map[string]Group
	sync.RWMutex
}

var _ Manager = (*MemoryManager)(nil)

func (m *MemoryManager) CreateGroup(g *Group) error {
	if g.ID == "" {
		g.ID = uuid.New()
	}
	if m.Groups == nil {
		m.Groups = map[string]Group{}
	}

	m.Groups[g.ID] = *g
	return nil
}

func (m *MemoryManager) GetGroup(id string) (*Group, error) {
	if g, ok := m.Groups[id]; !ok {
		return nil, errors.WithStack(pkg.ErrNotFound)
	} else {
		return &g, nil
	}
}

func (m *MemoryManager) DeleteGroup(id string) error {
	delete(m.Groups, id)
	return nil
}

func (m *MemoryManager) AddGroupMembers(group string, subjects []string) error {
	g, err := m.GetGroup(group)
	if err != nil {
		return err
	}
	g.Members = append(g.Members, subjects...)
	return m.CreateGroup(g)
}

func (m *MemoryManager) RemoveGroupMembers(group string, subjects []string) error {
	g, err := m.GetGroup(group)
	if err != nil {
		return err
	}

	var subs []string
	for _, s := range g.Members {
		var remove bool
		for _, f := range subjects {
			if f == s {
				remove = true
				break
			}
		}
		if !remove {
			subs = append(subs, s)
		}
	}

	g.Members = subs
	return m.CreateGroup(g)
}

func (m *MemoryManager) ListGroups(limit, offset int64) ([]Group, error) {
	if limit <= 0 {
		limit = 500
	}

	if offset < 0 {
		offset = 0
	}

	if offset >= int64(len(m.Groups)) {
		return []Group{}, nil
	}

	ids := []string{}
	for id := range m.Groups {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	res := make([]Group, len(ids))
	for i, id := range ids {
		res[i] = m.Groups[id]
	}

	res = res[offset:]

	if limit < int64(len(res)) {
		res = res[:limit]
	}

	return res, nil
}

func (m *MemoryManager) FindGroupsByMember(subject string) ([]Group, error) {
	if m.Groups == nil {
		m.Groups = map[string]Group{}
	}

	var res = []Group{}
	for _, g := range m.Groups {
		for _, s := range g.Members {
			if s == subject {
				res = append(res, g)
				break
			}
		}
	}

	return res, nil
}
