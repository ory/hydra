// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import "net/http"

// WrapTransportWithHeader wraps a http.Transport to always use the values from the given header.
func WrapTransportWithHeader(parent http.RoundTripper, h http.Header) *TransportWithHeader {
	return &TransportWithHeader{
		RoundTripper: parent,
		h:            h,
	}
}

// NewTransportWithHeader returns a new http.Transport that always uses the values from the given header.
func NewTransportWithHeader(h http.Header) *TransportWithHeader {
	return &TransportWithHeader{
		RoundTripper: http.DefaultTransport,
		h:            h,
	}
}

// TransportWithHeader is an http.RoundTripper that always uses the values from the given header.
type TransportWithHeader struct {
	http.RoundTripper
	h http.Header
}

// RoundTrip implements http.RoundTripper.
func (ct *TransportWithHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k := range ct.h {
		req.Header.Set(k, ct.h.Get(k))
	}
	return ct.RoundTripper.RoundTrip(req)
}

// NewTransportWithHost returns a new http.Transport that always uses the given host.
func NewTransportWithHost(host string) *TransportWithHost {
	return &TransportWithHost{
		RoundTripper: http.DefaultTransport,
		host:         host,
	}
}

// WrapRoundTripperWithHost wraps a http.RoundTripper that always uses the given host.
func WrapRoundTripperWithHost(parent http.RoundTripper, host string) *TransportWithHost {
	return &TransportWithHost{
		RoundTripper: parent,
		host:         host,
	}
}

// TransportWithHost is an http.RoundTripper that always uses the given host.
type TransportWithHost struct {
	http.RoundTripper
	host string
}

// RoundTrip implements http.RoundTripper.
func (ct *TransportWithHost) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = ct.host
	return ct.RoundTripper.RoundTrip(req)
}
