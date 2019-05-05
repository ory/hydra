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

package jwk

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/hydra/x"
)

type MemoryManager struct {
	Keys map[string]*jose.JSONWebKeySet
	sync.RWMutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		Keys: map[string]*jose.JSONWebKeySet{},
	}
}

func (m *MemoryManager) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	m.Lock()
	defer m.Unlock()

	m.alloc()
	if m.Keys[set] == nil {
		m.Keys[set] = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	}

	for _, k := range m.Keys[set].Keys {
		if k.KeyID == key.KeyID {
			return errors.WithStack(&fosite.RFC6749Error{
				Code:        http.StatusConflict,
				Name:        http.StatusText(http.StatusConflict),
				Description: fmt.Sprintf("Unable to create key with kid \"%s\" in set \"%s\" because that kid already exists in the set.", key.KeyID, set),
			})
		}
	}

	m.Keys[set].Keys = append([]jose.JSONWebKey{*key}, m.Keys[set].Keys...)
	return nil
}

func (m *MemoryManager) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	for _, key := range keys.Keys {
		if err := m.AddKey(ctx, set, &key); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryManager) GetKey(ctx context.Context, set, kid string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	result := keys.Key(kid)
	if len(result) == 0 {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	return &jose.JSONWebKeySet{
		Keys: result,
	}, nil
}

func (m *MemoryManager) GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	m.alloc()
	keys, found := m.Keys[set]
	if !found {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	if len(keys.Keys) == 0 {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	return keys, nil
}

func (m *MemoryManager) DeleteKey(ctx context.Context, set, kid string) error {
	keys, err := m.GetKeySet(ctx, set)
	if err != nil {
		return err
	}

	m.Lock()
	var results []jose.JSONWebKey
	for _, key := range keys.Keys {
		if key.KeyID != kid {
			results = append(results)
		}
	}
	m.Keys[set].Keys = results
	defer m.Unlock()

	return nil
}

func (m *MemoryManager) DeleteKeySet(ctx context.Context, set string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.Keys, set)
	return nil
}

func (m *MemoryManager) alloc() {
	if m.Keys == nil {
		m.Keys = make(map[string]*jose.JSONWebKeySet)
	}
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
	metricJWKs.WithLabelValues().Set(float64(len(m.Keys)))

	metricJWKs.Collect(c)
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
	metricJWKs.Describe(c)
}
