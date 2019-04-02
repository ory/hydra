package x

import (
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/negroni"

	"github.com/ory/x/healthx"
	"github.com/ory/x/stringsx"
)

func MatchesRange(r *http.Request, ranges []string) error {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	check := []string{remoteIP}
	for _, fwd := range stringsx.Splitx(r.Header.Get("X-Forwarded-For"), ",") {
		check = append(check, strings.TrimSpace(fwd))
	}

	for _, rn := range ranges {
		_, cidr, err := net.ParseCIDR(rn)
		if err != nil {
			return errors.WithStack(err)
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
	AllowTLSTerminationFrom() []string
	ServesHTTPS() bool
}

func RejectInsecureRequests(reg tlsRegistry, c tlsConfig) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.TLS != nil || !c.ServesHTTPS() || r.URL.Path == healthx.AliveCheckPath || r.URL.Path == healthx.ReadyCheckPath {
			next.ServeHTTP(rw, r)
			return
		}

		if len(c.AllowTLSTerminationFrom()) == 0 {
			reg.Logger().WithError(errors.New("TLS termination is not enabled")).Warnln("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		}

		ranges := c.AllowTLSTerminationFrom()
		if err := MatchesRange(r, ranges); err != nil {
			reg.Logger().WithError(err).Warnln("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		}

		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			reg.Logger().WithError(errors.New("X-Forwarded-Proto header is missing")).Warnln("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("can not serve request over insecure http"))
			return
		} else if proto != "https" {
			reg.Logger().WithError(errors.New("X-Forwarded-Proto header is missing")).Warnln("Could not serve http connection")
			reg.Writer().WriteErrorCode(rw, r, http.StatusBadGateway, errors.Errorf("expected X-Forwarded-Proto header to be https but got: %s", proto))
			return
		}

		next.ServeHTTP(rw, r)
	}
}
