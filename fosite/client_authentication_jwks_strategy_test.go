// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initServerWithKey(t *testing.T) *httptest.Server {
	var set *jose.JSONWebKeySet
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewEncoder(w).Encode(set))
	}
	ts := httptest.NewServer(h)

	set = &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				KeyID: "bar",
				Use:   "sig",
				Key:   &gen.MustRSAKey().PublicKey,
			},
		},
	}

	t.Cleanup(ts.Close)
	return ts
}

var errRoundTrip = errors.New("roundtrip error")

type failingTripper struct{}

func (r *failingTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errRoundTrip
}

func TestDefaultJWKSFetcherStrategy(t *testing.T) {
	ctx := context.Background()
	var h http.HandlerFunc

	s := NewDefaultJWKSFetcherStrategy()
	t.Run("case=fetching", func(t *testing.T) {
		var set *jose.JSONWebKeySet
		h = func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(set))
		}
		ts := httptest.NewServer(h)
		defer ts.Close()

		set = &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					KeyID: "foo",
					Use:   "sig",
					Key:   &gen.MustRSAKey().PublicKey,
				},
			},
		}

		keys, err := s.Resolve(ctx, ts.URL, false)
		require.NoError(t, err)
		assert.True(t, len(keys.Key("foo")) == 1)

		set = &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					KeyID: "bar",
					Use:   "sig",
					Key:   &gen.MustRSAKey().PublicKey,
				},
			},
		}

		keys, err = s.Resolve(ctx, ts.URL, false)
		require.NoError(t, err)
		assert.Len(t, keys.Keys, 1, "%+v", keys)
		assert.True(t, len(keys.Key("foo")) == 1)
		assert.True(t, len(keys.Key("bar")) == 0)

		keys, err = s.Resolve(ctx, ts.URL, true)
		require.NoError(t, err)
		assert.True(t, len(keys.Key("foo")) == 0)
		assert.True(t, len(keys.Key("bar")) == 1)
	})

	t.Run("JWKSFetcherWithCache", func(t *testing.T) {
		ts := initServerWithKey(t)

		cache, _ := ristretto.NewCache(&ristretto.Config[string, *jose.JSONWebKeySet]{NumCounters: 10 * 1000, MaxCost: 1000, BufferItems: 64})
		location := ts.URL
		expected := &jose.JSONWebKeySet{}
		require.True(t, cache.Set(defaultJWKSFetcherStrategyCachePrefix+location, expected, 1))
		cache.Wait()

		s := NewDefaultJWKSFetcherStrategy(JWKSFetcherWithCache(cache))
		actual, err := s.Resolve(ctx, location, false)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("JWKSFetcherWithTTL", func(t *testing.T) {
		ts := initServerWithKey(t)

		s := NewDefaultJWKSFetcherStrategy(JKWKSFetcherWithDefaultTTL(time.Nanosecond))
		_, err := s.Resolve(ctx, ts.URL, false)
		require.NoError(t, err)
		s.(*DefaultJWKSFetcherStrategy).cache.Wait()

		_, ok := s.(*DefaultJWKSFetcherStrategy).cache.Get(defaultJWKSFetcherStrategyCachePrefix + ts.URL)
		assert.Falsef(t, ok, "expected cache to be empty")
	})

	t.Run("JWKSFetcherWithHTTPClient", func(t *testing.T) {
		rt := retryablehttp.NewClient()
		rt.RetryMax = 0
		rt.HTTPClient = &http.Client{Transport: new(failingTripper)}
		s := NewDefaultJWKSFetcherStrategy(JWKSFetcherWithHTTPClient(rt))
		_, err := s.Resolve(ctx, "https://google.com", false)
		require.ErrorIs(t, err, errRoundTrip)
	})

	t.Run("JWKSFetcherWithHTTPClientSource", func(t *testing.T) {
		rt := retryablehttp.NewClient()
		rt.RetryMax = 0
		rt.HTTPClient = &http.Client{Transport: new(failingTripper)}
		s := NewDefaultJWKSFetcherStrategy(
			JWKSFetcherWithHTTPClient(retryablehttp.NewClient()),
			JWKSFetcherWithHTTPClientSource(func(ctx context.Context) *retryablehttp.Client {
				return rt
			}))
		_, err := s.Resolve(ctx, "https://www.google.com", false)
		require.ErrorIs(t, err, errRoundTrip)
	})

	t.Run("case=error_network", func(t *testing.T) {
		s := NewDefaultJWKSFetcherStrategy()
		h = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
		}
		ts := httptest.NewServer(h)
		defer ts.Close()

		_, err := s.Resolve(context.Background(), ts.URL, true)
		require.Error(t, err)

		_, err = s.Resolve(context.Background(), "$%/19", true)
		require.Error(t, err)
	})

	t.Run("case=error_encoding", func(t *testing.T) {
		s := NewDefaultJWKSFetcherStrategy()
		h = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("[]"))
		}
		ts := httptest.NewServer(h)
		defer ts.Close()

		_, err := s.Resolve(context.Background(), ts.URL, true)
		require.Error(t, err)
	})
}
