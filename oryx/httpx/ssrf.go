// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/netip"
	"time"

	"code.dny.dev/ssrf"
	"github.com/gobwas/glob"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var _ http.RoundTripper = (*noInternalIPRoundTripper)(nil)

type noInternalIPRoundTripper struct {
	onWhitelist, notOnWhitelist http.RoundTripper
	internalIPExceptions        []string
}

// RoundTrip implements http.RoundTripper.
func (n noInternalIPRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	incoming := IncomingRequestURL(request)
	incoming.RawQuery = ""
	incoming.RawFragment = ""
	for _, exception := range n.internalIPExceptions {
		compiled, err := glob.Compile(exception, '.', '/')
		if err != nil {
			return nil, err
		}
		if compiled.Match(incoming.String()) {
			return n.onWhitelist.RoundTrip(request)
		}
	}

	return n.notOnWhitelist.RoundTrip(request)
}

var (
	prohibitInternalAllowIPv6 http.RoundTripper = OTELTraceTransport(ssrfTransport(
		ssrf.WithAnyPort(),
		ssrf.WithNetworks("tcp4", "tcp6"),
	))

	allowInternalAllowIPv6 http.RoundTripper = OTELTraceTransport(ssrfTransport(
		ssrf.WithAnyPort(),
		ssrf.WithNetworks("tcp4", "tcp6"),
		ssrf.WithAllowedV4Prefixes(
			netip.MustParsePrefix("10.0.0.0/8"),     // Private-Use (RFC 1918)
			netip.MustParsePrefix("127.0.0.0/8"),    // Loopback (RFC 1122, Section 3.2.1.3))
			netip.MustParsePrefix("169.254.0.0/16"), // Link Local (RFC 3927)
			netip.MustParsePrefix("172.16.0.0/12"),  // Private-Use (RFC 1918)
			netip.MustParsePrefix("192.168.0.0/16"), // Private-Use (RFC 1918)
		),
		ssrf.WithAllowedV6Prefixes(
			netip.MustParsePrefix("::1/128"),  // Loopback (RFC 4193)
			netip.MustParsePrefix("fc00::/7"), // Unique Local (RFC 4193)
		),
	))
)

func ssrfTransport(opt ...ssrf.Option) *http.Transport {
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	dialer.Control = ssrf.New(opt...).Safe
	dial := func(ctx context.Context, network string, address string) (net.Conn, error) {
		c, err := dialer.DialContext(ctx, network, address)
		if err == nil {
			return c, nil
		}

		if dnsErr, ok := errors.AsType[*net.DNSError](err); ok {
			// Copy the `*net.DNSError` before masking `Server` to avoid a data
			// race: the DNS resolver uses `singleflight` to deduplicate
			// concurrent lookups, so multiple goroutines may receive the same
			// `*net.DNSError` pointer. Mutating it in place races with concurrent
			// readers (e.g. the `otelhttp` `dnsDone` trace hook).
			maskedDNS := *dnsErr
			maskedDNS.Server = "" // Mask our DNS server's IP address.
			// Also copy the outer `*net.OpError` (if present) to preserve the
			// full error chain that callers depend on, replacing its `Err` field
			// with the masked copy.
			// Surprisingly Go does not have a good way of manipulating error chains.
			if opErr, ok := errors.AsType[*net.OpError](err); ok {
				maskedOp := *opErr
				maskedOp.Err = &maskedDNS
				return nil, &maskedOp
			}
			return nil, &maskedDNS
		}

		if !errors.Is(err, ssrf.ErrProhibitedIP) {
			return nil, err
		}

		host, _, _ := net.SplitHostPort(address)
		_, addrErr := netip.ParseAddrPort(address)
		if addrErr != nil {
			// We were given a DNS name: the error we return must look like a DNS error.
			return nil, &net.OpError{
				Op:   "dial",
				Net:  network,
				Addr: nil,
				Err: &net.DNSError{
					Err:         "no such host",
					Name:        host,
					Server:      "",
					IsTimeout:   false,
					IsTemporary: false,
					IsNotFound:  true,
				},
			}
		}
		return nil, &net.OpError{
			Op:   "dial",
			Net:  network,
			Addr: nil,
			Err:  errors.New("no route to host"),
		}
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dial,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// OTELTraceTransport wraps the given http.Transport with OpenTelemetry instrumentation.
func OTELTraceTransport(t *http.Transport) http.RoundTripper {
	return otelhttp.NewTransport(t, otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
		return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutHeaders(), otelhttptrace.WithoutSubSpans())
	}))
}
