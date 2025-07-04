// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"net/netip"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoPrivateIPs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	}))
	t.Cleanup(ts.Close)

	target, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	_, port, err := net.SplitHostPort(target.Host)
	require.NoError(t, err)

	allowedURL := "http://localhost:" + port + "/foobar"
	allowedGlob := "http://localhost:" + port + "/glob/*"

	c := NewResilientClient(
		ResilientClientWithMaxRetry(1),
		ResilientClientDisallowInternalIPs(),
		ResilientClientAllowInternalIPRequestsTo(allowedURL, allowedGlob),
	)

	for i := 0; i < 10; i++ {
		for destination, passes := range map[string]bool{
			"http://127.0.0.1:" + port:                   false,
			"http://localhost:" + port:                   false,
			"http://192.168.178.5:" + port:               false,
			allowedURL:                                   true,
			"http://localhost:" + port + "/glob/bar":     true,
			"http://localhost:" + port + "/glob/bar/baz": false,
			"http://localhost:" + port + "/FOOBAR":       false,
		} {
			_, err := c.Get(destination)
			if !passes {
				require.Errorf(t, err, "dest = %s", destination)
				assert.Containsf(t, err.Error(), "is not a permitted destination", "dest = %s", destination)
			} else {
				require.NoErrorf(t, err, "dest = %s", destination)
			}
		}
	}
}

func TestNoIPV6(t *testing.T) {
	for _, tc := range []struct {
		name string
		c    *retryablehttp.Client
	}{
		{
			"internal IPs allowed",
			NewResilientClient(
				ResilientClientWithMaxRetry(1),
				ResilientClientNoIPv6(),
			),
		}, {
			"internal IPs disallowed",
			NewResilientClient(
				ResilientClientWithMaxRetry(1),
				ResilientClientDisallowInternalIPs(),
				ResilientClientNoIPv6(),
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var connectDone int32
			ctx := httptrace.WithClientTrace(context.Background(), &httptrace.ClientTrace{
				DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
					for _, ip := range dnsInfo.Addrs {
						netIP, ok := netip.AddrFromSlice(ip.IP)
						assert.True(t, ok)
						assert.Truef(t, netIP.Is4(), "ip = %s", ip)
					}
				},
				ConnectDone: func(network, addr string, err error) {
					atomic.AddInt32(&connectDone, 1)
					assert.NoError(t, err)
					assert.Equalf(t, "tcp4", network, "network = %s addr = %s", network, addr)
				},
			})

			// Dual stack
			req, err := retryablehttp.NewRequestWithContext(ctx, "GET", "http://dual.tlund.se/", nil)
			require.NoError(t, err)
			atomic.StoreInt32(&connectDone, 0)
			res, err := tc.c.Do(req)
			require.GreaterOrEqual(t, int32(1), atomic.LoadInt32(&connectDone))
			require.NoError(t, err)
			t.Cleanup(func() { _ = res.Body.Close() })
			require.EqualValues(t, http.StatusOK, res.StatusCode)

			// IPv4 only
			req, err = retryablehttp.NewRequestWithContext(ctx, "GET", "http://ipv4.tlund.se/", nil)
			require.NoError(t, err)
			atomic.StoreInt32(&connectDone, 0)
			res, err = tc.c.Do(req)
			require.EqualValues(t, 1, atomic.LoadInt32(&connectDone))
			require.NoError(t, err)
			t.Cleanup(func() { _ = res.Body.Close() })
			require.EqualValues(t, http.StatusOK, res.StatusCode)

			// IPv6 only
			req, err = retryablehttp.NewRequestWithContext(ctx, "GET", "http://ipv6.tlund.se/", nil)
			require.NoError(t, err)
			atomic.StoreInt32(&connectDone, 0)
			_, err = tc.c.Do(req)
			require.EqualValues(t, 0, atomic.LoadInt32(&connectDone))
			require.ErrorContains(t, err, "no such host")
		})
	}
}
