// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package group

import (
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
