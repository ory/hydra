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

	"github.com/ory/hydra/x"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/fosite"
	"github.com/ory/x/pagination"
	"github.com/ory/x/sqlcon"
)

type MemoryManager struct {
	r       InternalRegistry
	Clients []Client
	sync.RWMutex
}

func NewMemoryManager(r InternalRegistry) *MemoryManager {
	return &MemoryManager{
		Clients: []Client{},
		r:       r,
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

	return nil, errors.WithStack(x.ErrNotFound)
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
		h, err := m.r.ClientHasher().Hash(ctx, []byte(c.Secret))
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

	if err := m.r.ClientHasher().Compare(ctx, c.GetHashedSecret(), secret); err != nil {
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

	hash, err := m.r.ClientHasher().Hash(ctx, []byte(c.Secret))
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

func (m *MemoryManager) CountClients(ctx context.Context) (n int, err error) {
	return len(m.Clients), nil
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent. The
// descriptor of each sent metric is one of those returned by Describe
// (unless the Collector is unchecked, see above). Returned metrics that
// share the same descriptor must differ in their variable label
// values.
//
// This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. Blocking occurs at the expense
// of total performance of rendering all registered metrics. Ideally,
// Collector implementations support concurrent readers.
func (m *MemoryManager) Collect(c chan<- prometheus.Metric) {
	metricClients.Set(float64(len(m.Clients)))
	metricClients.Collect(c)
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent. The sent descriptors fulfill the
// consistency and uniqueness requirements described in the Desc
// documentation.
//
// It is valid if one and the same Collector sends duplicate
// descriptors. Those duplicates are simply ignored. However, two
// different Collectors must not send duplicate descriptors.
//
// Sending no descriptor at all marks the Collector as “unchecked”,
// i.e. no checks will be performed at registration time, and the
// Collector may yield any Metric it sees fit in its Collect method.
//
// This method idempotently sends the same descriptors throughout the
// lifetime of the Collector. It may be called concurrently and
// therefore must be implemented in a concurrency safe way.
//
// If a Collector encounters an error while executing this method, it
// must send an invalid descriptor (created with NewInvalidDesc) to
// signal the error to the registry.
func (m *MemoryManager) Describe(c chan<- *prometheus.Desc) {
	metricClients.Describe(c)
}
