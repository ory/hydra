/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package client

import (
	"context"
	"sync"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/x/pagination"
	"github.com/ory/x/sqlcon"
)

type MemoryManager struct {
	Clients []Client
	Hasher  fosite.Hasher
	sync.RWMutex
}

func NewMemoryManager(hasher fosite.Hasher) *MemoryManager {
	if hasher == nil {
		hasher = new(fosite.BCrypt)
	}

	return &MemoryManager{
		Clients: []Client{},
		Hasher:  hasher,
	}
}

func (m *MemoryManager) GetConcreteClient(ctx context.Context, id string) (*Client, error) {
	m.RLock()
	defer m.RUnlock()

	for _, c := range m.Clients {
		if c.GetID() == id {
			return &c, nil
		}
	}

	return nil, errors.WithStack(sqlcon.ErrNoRows)
}

func (m *MemoryManager) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return m.GetConcreteClient(ctx, id)
}

func (m *MemoryManager) UpdateClient(ctx context.Context, c *Client) error {
	o, err := m.GetClient(ctx, c.GetID())
	if err != nil {
		return err
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash(ctx, []byte(c.Secret))
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = string(h)
	}
	if err := mergo.Merge(c, o); err != nil {
		return errors.WithStack(err)
	}

	m.Lock()
	defer m.Unlock()
	for k, f := range m.Clients {
		if f.GetID() == c.GetID() {
			m.Clients[k] = *c
		}
	}

	return nil
}

func (m *MemoryManager) Authenticate(ctx context.Context, id string, secret []byte) (*Client, error) {
	m.RLock()
	defer m.RUnlock()

	c, err := m.GetConcreteClient(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := m.Hasher.Compare(ctx, c.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (m *MemoryManager) CreateClient(ctx context.Context, c *Client) error {
	if _, err := m.GetConcreteClient(ctx, c.GetID()); err == nil {
		return sqlcon.ErrUniqueViolation
	}

	m.Lock()
	defer m.Unlock()

	hash, err := m.Hasher.Hash(ctx, []byte(c.Secret))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = string(hash)

	m.Clients = append(m.Clients, *c)
	return nil
}

func (m *MemoryManager) DeleteClient(ctx context.Context, id string) error {
	m.Lock()
	defer m.Unlock()

	for k, f := range m.Clients {
		if f.GetID() == id {
			m.Clients = append(m.Clients[:k], m.Clients[k+1:]...)
			return nil
		}
	}

	return nil
}

func (m *MemoryManager) GetClients(ctx context.Context, limit, offset int) (clients map[string]Client, err error) {
	m.RLock()
	defer m.RUnlock()
	clients = make(map[string]Client)

	start, end := pagination.Index(limit, offset, len(m.Clients))
	for _, c := range m.Clients[start:end] {
		clients[c.GetID()] = c
	}

	return clients, nil
}
