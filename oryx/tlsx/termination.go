// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package tlsx

import (
	"net"
	"net/http"
	"strings"

	"github.com/ory/x/httpx"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"

	"github.com/ory/x/healthx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/prometheusx"
)

type dependencies interface {
	logrusx.Provider
	httpx.WriterProvider
}

// EnforceTLSRequests creates a middleware that enforces TLS for incoming HTTP requests.
// It allows termination (non-HTTPS traffic) from specific CIDR ranges provided in the `allowTerminationFrom` slice.
// If the request is not secure and does not match the allowed CIDR ranges, an error response is returned.
// The middleware also validates the `X-Forwarded-Proto` header to ensure it is set to "https".
func EnforceTLSRequests(d dependencies, allowTerminationFrom []string) (negroni.Handler, error) {
	networks := make([]*net.IPNet, 0, len(allowTerminationFrom))
	for _, rn := range allowTerminationFrom {
		_, network, err := net.ParseCIDR(rn)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		networks = append(networks, network)
	}

	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.TLS != nil ||
			r.URL.Path == healthx.AliveCheckPath ||
			r.URL.Path == healthx.ReadyCheckPath ||
			r.URL.Path == prometheusx.MetricsPrometheusPath {
			next(rw, r)
			return
		}

		if len(networks) == 0 {
			d.Logger().WithRequest(r).WithError(errors.New("TLS termination is not enabled")).Error("Could not serve http connection")
			d.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		}

		if err := matchesRange(r, networks); err != nil {
			d.Logger().WithRequest(r).WithError(err).Warnln("Could not serve http connection")
			d.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		}

		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			d.Logger().WithRequest(r).WithError(errors.New("X-Forwarded-Proto header is missing")).Error("Could not serve http connection")
			d.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		} else if proto != "https" {
			d.Logger().WithRequest(r).WithError(errors.New("X-Forwarded-Proto header is missing")).Error("Could not serve http connection")
			d.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.Errorf("expected X-Forwarded-Proto header to be https but got: %s", proto))
			return
		}

		next(rw, r)
	}), nil
}

func matchesRange(r *http.Request, networks []*net.IPNet) error {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	check := []string{remoteIP}
	for fwd := range strings.SplitSeq(r.Header.Get("X-Forwarded-For"), ",") {
		check = append(check, strings.TrimSpace(fwd))
	}

	for _, ipNet := range networks {
		for _, ip := range check {
			addr := net.ParseIP(ip)
			if ipNet.Contains(addr) {
				return nil
			}
		}
	}
	return errors.Errorf("neither remote address nor any x-forwarded-for values match CIDR ranges %+v: %v, ranges, check)", networks, check)
}
