// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net/http"
	"time"

	"github.com/ory/x/reqlog"
)

// MeasureExternalLatencyTransport is an http.RoundTripper that measures the latency of all requests as external latency.
type MeasureExternalLatencyTransport struct {
	Transport http.RoundTripper
}

var _ http.RoundTripper = (*MeasureExternalLatencyTransport)(nil)

func (m *MeasureExternalLatencyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	upstreamHostPath := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
	defer reqlog.StartMeasureExternalCall(req.Context(), "http_request", upstreamHostPath, time.Now())

	t := m.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	return t.RoundTrip(req)
}
