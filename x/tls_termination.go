// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"net"
	"net/http"
	"strings"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
	"github.com/urfave/negroni"

	"github.com/ory/x/healthx"
	prometheus "github.com/ory/x/prometheusx"
	"github.com/ory/x/stringsx"
)

func MatchesRange(r *http.Request, ranges []string) error {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return errorsx.WithStack(err)
	}

	check := []string{remoteIP}
	for _, fwd := range stringsx.Splitx(r.Header.Get("X-Forwarded-For"), ",") {
		check = append(check, strings.TrimSpace(fwd))
	}

	for _, rn := range ranges {
		_, cidr, err := net.ParseCIDR(rn)
		if err != nil {
			return errorsx.WithStack(err)
		}

		for _, ip := range check {
			addr := net.ParseIP(ip)
			if cidr.Contains(addr) {
				return nil
			}
		}
	}
	return errors.Errorf("neither remote address nor any x-forwarded-for values match CIDR ranges %v: %v, ranges, check)", ranges, check)
}

type tlsRegistry interface {
	RegistryLogger
	RegistryWriter
}

type tlsConfig interface {
	Enabled() bool
	AllowTerminationFrom() []string
}

func RejectInsecureRequests(reg tlsRegistry, c tlsConfig) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.TLS != nil ||
			!c.Enabled() ||
			r.URL.Path == healthx.AliveCheckPath ||
			r.URL.Path == healthx.ReadyCheckPath ||
			r.URL.Path == prometheus.MetricsPrometheusPath {
			next.ServeHTTP(rw, r)
			return
		}

		if len(c.AllowTerminationFrom()) == 0 {
			reg.Logger().WithRequest(r).WithError(errors.New("TLS termination is not enabled")).Error("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		}

		ranges := c.AllowTerminationFrom()
		if err := MatchesRange(r, ranges); err != nil {
			reg.Logger().WithRequest(r).WithError(err).Warnln("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		}

		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			reg.Logger().WithRequest(r).WithError(errors.New("X-Forwarded-Proto header is missing")).Error("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		} else if proto != "https" {
			reg.Logger().WithRequest(r).WithError(errors.New("X-Forwarded-Proto header is missing")).Error("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.Errorf("expected X-Forwarded-Proto header to be https but got: %s", proto))
			return
		}

		next.ServeHTTP(rw, r)
	}
}
