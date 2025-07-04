// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fetcher

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	stderrors "errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/ory/x/httpx"
	"github.com/ory/x/stringsx"
)

// Fetcher is able to load file contents from http, https, file, and base64 locations.
type Fetcher struct {
	hc    *retryablehttp.Client
	limit int64
	cache *ristretto.Cache[[]byte, []byte]
	ttl   time.Duration
}

type opts struct {
	hc    *retryablehttp.Client
	limit int64
	cache *ristretto.Cache[[]byte, []byte]
	ttl   time.Duration
}

var ErrUnknownScheme = stderrors.New("unknown scheme")

// WithClient sets the http.Client the fetcher uses.
func WithClient(hc *retryablehttp.Client) Modifier {
	return func(o *opts) {
		o.hc = hc
	}
}

// WithMaxHTTPMaxBytes reads at most limit bytes from the HTTP response body,
// returning bytes.ErrToLarge if the limit would be exceeded.
func WithMaxHTTPMaxBytes(limit int64) Modifier {
	return func(o *opts) {
		o.limit = limit
	}
}

func WithCache(cache *ristretto.Cache[[]byte, []byte], ttl time.Duration) Modifier {
	return func(o *opts) {
		if ttl < 0 {
			return
		}
		o.cache = cache
		o.ttl = ttl
	}
}

func newOpts() *opts {
	return &opts{
		hc: httpx.NewResilientClient(),
	}
}

type Modifier func(*opts)

// NewFetcher creates a new fetcher instance.
func NewFetcher(opts ...Modifier) *Fetcher {
	o := newOpts()
	for _, f := range opts {
		f(o)
	}
	return &Fetcher{hc: o.hc, limit: o.limit, cache: o.cache, ttl: o.ttl}
}

// Fetch fetches the file contents from the source.
func (f *Fetcher) Fetch(source string) (*bytes.Buffer, error) {
	return f.FetchContext(context.Background(), source)
}

// FetchContext fetches the file contents from the source and allows to pass a
// context that is used for HTTP requests.
func (f *Fetcher) FetchContext(ctx context.Context, source string) (*bytes.Buffer, error) {
	b, err := f.FetchBytes(ctx, source)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

// FetchBytes fetches the file contents from the source and allows to pass a
// context that is used for HTTP requests.
func (f *Fetcher) FetchBytes(ctx context.Context, source string) ([]byte, error) {
	switch s := stringsx.SwitchPrefix(source); {
	case s.HasPrefix("http://", "https://"):
		return f.fetchRemote(ctx, source)
	case s.HasPrefix("file://"):
		return f.fetchFile(strings.TrimPrefix(source, "file://"))
	case s.HasPrefix("base64://"):
		src, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(source, "base64://"))
		if err != nil {
			return nil, errors.Wrapf(err, "base64decode: %s", source)
		}
		return src, nil
	default:
		return nil, errors.Wrap(ErrUnknownScheme, s.ToUnknownPrefixErr().Error())
	}
}

func (f *Fetcher) fetchRemote(ctx context.Context, source string) (b []byte, err error) {
	if f.cache != nil {
		cacheKey := sha256.Sum256([]byte(source))
		if v, ok := f.cache.Get(cacheKey[:]); ok {
			b = make([]byte, len(v))
			copy(b, v)
			return b, nil
		}
		defer func() {
			if err == nil && len(b) > 0 {
				toCache := make([]byte, len(b))
				copy(toCache, b)
				f.cache.SetWithTTL(cacheKey[:], toCache, int64(len(toCache)), f.ttl)
			}
		}()
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "new request: %s", source)
	}
	res, err := f.hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, source)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("expected http response status code 200 but got %d when fetching: %s", res.StatusCode, source)
	}

	if f.limit > 0 {
		var buf bytes.Buffer
		n, err := io.Copy(&buf, io.LimitReader(res.Body, f.limit+1))
		if n > f.limit {
			return nil, bytes.ErrTooLarge
		}
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return io.ReadAll(res.Body)
}

func (f *Fetcher) fetchFile(source string) ([]byte, error) {
	fp, err := os.Open(source) // #nosec:G304
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open file: %s", source)
	}
	defer fp.Close()
	b, err := io.ReadAll(fp)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file: %s", source)
	}
	if err := fp.Close(); err != nil {
		return nil, errors.Wrapf(err, "unable to close file: %s", source)
	}
	return b, nil
}
