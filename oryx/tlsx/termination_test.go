// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package tlsx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/herodot"
	"github.com/ory/x/healthx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/prometheusx"
)

func failHandler(t *testing.T) http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		t.Fatal("handler should not have been called")
	}
}

func noopHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

type dependencyProvider struct {
	l *logrusx.Logger
	w herodot.Writer
}

func (d *dependencyProvider) Logger() *logrusx.Logger { return d.l }
func (d *dependencyProvider) Writer() herodot.Writer  { return d.w }

func TestRejectInsecureRequests(t *testing.T) {
	d := &dependencyProvider{
		l: logrusx.New("", ""),
		w: herodot.NewJSONWriter(logrusx.New("", "")),
	}

	allowedRanges := []string{"126.0.0.1/24", "127.0.0.1/24"}

	const (
		addrInRange          = "127.0.0.1"
		remoteAddrInRange    = "127.0.0.1:123"
		addrNotInRange       = "227.0.0.1"
		remoteAddrNotInRange = "227.0.0.1:123s"
	)

	t.Run("no allowTerminationFrom set", func(t *testing.T) {
		res := httptest.NewRecorder()
		h, err := EnforceTLSRequests(d, nil)
		require.NoError(t, err)
		h.ServeHTTP(res, &http.Request{RemoteAddr: remoteAddrNotInRange, Header: http.Header{}, URL: new(url.URL)}, failHandler(t))
		assert.EqualValues(t, http.StatusBadGateway, res.Code)

		res = httptest.NewRecorder()
		h, err = EnforceTLSRequests(d, []string{})
		require.NoError(t, err)
		h.ServeHTTP(res, &http.Request{RemoteAddr: remoteAddrNotInRange, Header: http.Header{}, URL: new(url.URL)}, failHandler(t))
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	t.Run("invalid CIDR", func(t *testing.T) {
		_, err := EnforceTLSRequests(d, []string{"invalidCIDR"})
		assert.ErrorContains(t, err, "invalid CIDR address")
	})

	for _, tc := range []struct {
		name          string
		req           *http.Request
		expectBlocked bool
	}{{
		name: "missing x-forwarded-proto",
		req: &http.Request{
			RemoteAddr: remoteAddrInRange,
			Header:     http.Header{},
			URL:        new(url.URL),
		},
		expectBlocked: true,
	}, {
		name: "x-forwarded-proto is http",
		req: &http.Request{
			RemoteAddr: remoteAddrInRange,
			Header:     http.Header{"X-Forwarded-Proto": []string{"http"}},
			URL:        new(url.URL),
		},
		expectBlocked: true,
	}, {
		name: "missing x-forwarded-for",
		req: &http.Request{
			Header: http.Header{"X-Forwarded-Proto": []string{"https"}},
			URL:    new(url.URL),
		},
		expectBlocked: true,
	}, {
		name: "remote not in any range",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
			URL:        new(url.URL),
		},
		expectBlocked: true,
	}, {
		name: "remote and forwarded not in any range",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header: http.Header{
				"X-Forwarded-Proto": []string{"https"},
				"X-Forwarded-For":   []string{addrNotInRange},
			},
			URL: new(url.URL),
		},
		expectBlocked: true,
	}, {
		name: "remote is in some range",
		req: &http.Request{
			RemoteAddr: remoteAddrInRange,
			Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
			URL:        new(url.URL),
		},
		expectBlocked: false,
	}, {
		name: "one of x-forwarded-for is in some range",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header: http.Header{
				"X-Forwarded-For":   []string{fmt.Sprintf("%s, %s, %s", addrNotInRange, addrInRange, addrNotInRange)},
				"X-Forwarded-Proto": []string{"https"},
			},
			URL: new(url.URL),
		},
		expectBlocked: false,
	}, {
		name: "health alive check is exempted",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header:     http.Header{},
			URL:        &url.URL{Path: healthx.AliveCheckPath},
		},
		expectBlocked: false,
	}, {
		name: "health ready check is exempted",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header:     http.Header{},
			URL:        &url.URL{Path: healthx.ReadyCheckPath},
		},
		expectBlocked: false,
	}, {
		name: "metrics prometheus check is exempted",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header:     http.Header{},
			URL:        &url.URL{Path: prometheusx.MetricsPrometheusPath},
		},
	}, {
		name: "x-forwarded-for without spaces",
		req: &http.Request{
			RemoteAddr: remoteAddrNotInRange,
			Header: http.Header{
				"X-Forwarded-For":   []string{fmt.Sprintf("%s,%s,%s", addrNotInRange, addrInRange, addrNotInRange)},
				"X-Forwarded-Proto": []string{"https"},
			},
			URL: new(url.URL),
		},
		expectBlocked: false,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			handler := noopHandler
			expectedStatus := http.StatusNoContent
			if tc.expectBlocked {
				handler = failHandler(t)
				expectedStatus = http.StatusBadGateway
			}
			h, err := EnforceTLSRequests(d, allowedRanges)
			require.NoError(t, err)
			h.ServeHTTP(res, tc.req, handler)
			assert.EqualValues(t, expectedStatus, res.Code)
		})
	}
}
