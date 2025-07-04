// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwksx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwk"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/snapshotx"
)

const (
	multiKeys = `{
  "keys": [
    {
      "use": "sig",
      "kty": "oct",
      "kid": "7d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8",
      "alg": "HS256",
      "k": "Y2hhbmdlbWVjaGFuZ2VtZWNoYW5nZW1lY2hhbmdlbWU"
    },
    {
      "use": "sig",
      "kty": "oct",
      "kid": "8d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8",
      "alg": "HS256",
      "k": "Y2hhbmdlbWVjaGFuZ2VtZWNoYW5nZW1lY2hhbmdlbWU"
    },
    {
      "use": "sig",
      "kty": "oct",
      "kid": "9d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8",
      "alg": "HS256",
      "k": "Y2hhbmdlbWVjaGFuZ2VtZWNoYW5nZW1lY2hhbmdlbWU"
    }
  ]
}`
)

type brokenTransport struct{}

var _ http.RoundTripper = new(brokenTransport)
var errBroken = errors.New("broken")

func (b brokenTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, errBroken
}

func TestFetcherNext(t *testing.T) {
	ctx := context.Background()
	cache, err := ristretto.NewCache[[]byte, jwk.Set](&ristretto.Config[[]byte, jwk.Set]{
		NumCounters:        100 * 10,
		MaxCost:            100,
		BufferItems:        64,
		Metrics:            true,
		IgnoreInternalCost: true,
		Cost: func(jwk.Set) int64 {
			return 1
		},
	})
	require.NoError(t, err)

	f := NewFetcherNext(cache)

	createRemoteProvider := func(called *int, payload string) *httptest.Server {
		cache.Clear()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*called++
			_, _ = w.Write([]byte(payload))
		}))
		t.Cleanup(ts.Close)
		return ts
	}

	t.Run("case=resolve multiple source urls", func(t *testing.T) {
		t.Run("case=fails without forced kid", func(t *testing.T) {
			var called int
			ts1 := createRemoteProvider(&called, keys)
			ts2 := createRemoteProvider(&called, multiKeys)

			_, err := f.ResolveKeyFromLocations(ctx, []string{ts1.URL, ts2.URL})
			require.Error(t, err)
		})
		t.Run("case=succeeds with forced kid", func(t *testing.T) {
			var called int
			ts1 := createRemoteProvider(&called, keys)
			ts2 := createRemoteProvider(&called, multiKeys)

			k, err := f.ResolveKeyFromLocations(ctx, []string{ts1.URL, ts2.URL}, WithForceKID("8d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8"))
			require.NoError(t, err)
			snapshotx.SnapshotT(t, k)
		})
	})
	t.Run("case=resolve single source url", func(t *testing.T) {
		t.Run("case=with forced key", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, keys)

			k, err := f.ResolveKey(ctx, ts.URL, WithForceKID("7d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8"))
			require.NoError(t, err)
			snapshotx.SnapshotT(t, k)
		})

		t.Run("case=forced key is not found", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, keys)

			_, err := f.ResolveKey(ctx, ts.URL, WithForceKID("not-found"))
			require.Error(t, err)
		})

		t.Run("case=no key in remote", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, "{}")

			_, err := f.ResolveKey(ctx, ts.URL)
			require.Error(t, err)
		})

		t.Run("case=remote not returning JSON", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, "lol")

			_, err := f.ResolveKey(ctx, ts.URL)
			require.Error(t, err)
		})

		t.Run("case=without cache", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, keys)

			k, err := f.ResolveKey(ctx, ts.URL)
			require.NoError(t, err)
			snapshotx.SnapshotT(t, k)
			assert.Equal(t, called, 1)

			cache.Wait()

			_, err = f.ResolveKey(ctx, ts.URL)
			require.NoError(t, err)
			assert.Equal(t, called, 2)
		})

		t.Run("case=with cache", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, keys)

			k, err := f.ResolveKey(ctx, ts.URL, WithCacheEnabled())
			require.NoError(t, err)
			assert.Equal(t, called, 1)

			cache.Wait()

			k, err = f.ResolveKey(ctx, ts.URL, WithCacheEnabled())
			require.NoError(t, err)
			assert.Equal(t, called, 1)

			snapshotx.SnapshotT(t, k)
		})

		t.Run("case=with cache and TTL", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, keys)
			waitTime := time.Millisecond * 100

			k, err := f.ResolveKey(ctx, ts.URL, WithCacheEnabled(), WithCacheTTL(waitTime))
			require.NoError(t, err)
			assert.Equal(t, called, 1)

			cache.Wait()

			k, err = f.ResolveKey(ctx, ts.URL, WithCacheEnabled())
			require.NoError(t, err)
			assert.Equal(t, called, 1)

			time.Sleep(waitTime)

			cache.Wait()

			k, err = f.ResolveKey(ctx, ts.URL, WithCacheEnabled())
			require.NoError(t, err)
			assert.Equal(t, called, 2)

			snapshotx.SnapshotT(t, k)
		})

		t.Run("case=with broken HTTP client", func(t *testing.T) {
			var called int
			ts := createRemoteProvider(&called, keys)

			broken := retryablehttp.NewClient()
			broken.RetryMax = 0
			broken.HTTPClient.Transport = new(brokenTransport)

			_, err := f.ResolveKey(ctx, ts.URL, WithHTTPClient(broken))
			require.ErrorIs(t, err, errBroken)
		})
	})
}
