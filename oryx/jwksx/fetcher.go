// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwksx

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
)

// Fetcher is a small helper for fetching JSON Web Keys from remote endpoints.
//
// DEPRECATED: Use FetcherNext instead.
type Fetcher struct {
	sync.RWMutex
	remote string
	c      *http.Client
	keys   map[string]jose.JSONWebKey
}

// NewFetcher returns a new fetcher that can download JSON Web Keys from remote endpoints.
//
// DEPRECATED: Use FetcherNext instead.
func NewFetcher(remote string) *Fetcher {
	return &Fetcher{
		remote: remote,
		c:      http.DefaultClient,
		keys:   make(map[string]jose.JSONWebKey),
	}
}

// GetKey retrieves a JSON Web Key from the cache, fetches it from a remote if it is not yet cached or returns an error.
//
// DEPRECATED: Use FetcherNext instead.
func (f *Fetcher) GetKey(kid string) (*jose.JSONWebKey, error) {
	f.RLock()
	if k, ok := f.keys[kid]; ok {
		f.RUnlock()
		return &k, nil
	}
	f.RUnlock()

	res, err := f.c.Get(f.remote)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("expected status code 200 but got %d when requesting %s", res.StatusCode, f.remote)
	}

	var set jose.JSONWebKeySet
	if err := json.NewDecoder(res.Body).Decode(&set); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, k := range set.Keys {
		f.Lock()
		f.keys[k.KeyID] = k
		f.Unlock()
	}

	f.RLock()
	defer f.RUnlock()
	if k, ok := f.keys[kid]; ok {
		return &k, nil
	}

	return nil, errors.Errorf("unable to find JSON Web Key with ID: %s", kid)
}
