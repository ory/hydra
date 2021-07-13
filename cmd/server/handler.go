/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ory/x/configx"

	analytics "github.com/ory/analytics-go/v4"

	"github.com/ory/x/reqlog"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/ory/graceful"
	"github.com/ory/x/healthx"
	"github.com/ory/x/metricsx"
	"github.com/ory/x/networkx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	prometheus "github.com/ory/x/prometheusx"
)

var _ = &consent.Handler{}

func EnhanceMiddleware(d driver.Registry, n *negroni.Negroni, address string, router *httprouter.Router, enableCORS bool, iface config.ServeInterface) http.Handler {
	if !networkx.AddressIsUnixSocket(address) {
		n.UseFunc(x.RejectInsecureRequests(d, d.Config().TLS(iface)))
	}
	n.UseHandler(router)

	if !enableCORS {
		return n
	}

	options, enabled := d.Config().CORS(iface)
	if !enabled {
		return n
	}

	if enabled {
		d.Logger().
			WithField("options", fmt.Sprintf("%+v", options)).
			Infof("Enabling CORS on interface: %s", address)
		return cors.New(options).Handler(n)
	}
	return n
}

func isDSNAllowed(r driver.Registry) {
	if r.Config().DSN() == "memory" {
		r.Logger().Fatalf(`When using "hydra serve admin" or "hydra serve public" the DSN can not be set to "memory".`)
	}
}

func RunServeAdmin(cmd *cobra.Command, args []string) {
	d := driver.New(cmd.Context(), driver.WithOptions(configx.WithFlags(cmd.Flags())))
	isDSNAllowed(d)

	admin, _, adminmw, _ := setup(d, cmd)
	cert := GetOrCreateTLSCertificate(cmd, d, config.AdminInterface) // we do not want to run this concurrently.

	d.PrometheusManager().RegisterRouter(admin.Router)

	var wg sync.WaitGroup
	wg.Add(1)

	go serve(
		d,
		cmd,
		&wg,
		config.AdminInterface,
		EnhanceMiddleware(d, adminmw, d.Config().ListenOn(config.AdminInterface), admin.Router, true, config.AdminInterface),
		d.Config().ListenOn(config.AdminInterface),
		d.Config().SocketPermission(config.AdminInterface),
		cert,
	)

	wg.Wait()
}

func RunServePublic(cmd *cobra.Command, args []string) {
	d := driver.New(cmd.Context(), driver.WithOptions(configx.WithFlags(cmd.Flags())))
	isDSNAllowed(d)

	_, public, _, publicmw := setup(d, cmd)
	cert := GetOrCreateTLSCertificate(cmd, d, config.PublicInterface) // we do not want to run this concurrently.

	d.PrometheusManager().RegisterRouter(public.Router)

	var wg sync.WaitGroup
	wg.Add(1)

	go serve(
		d,
		cmd,
		&wg,
		config.PublicInterface,
		EnhanceMiddleware(d, publicmw, d.Config().ListenOn(config.PublicInterface), public.Router, false, config.PublicInterface),
		d.Config().ListenOn(config.PublicInterface),
		d.Config().SocketPermission(config.PublicInterface),
		cert,
	)

	wg.Wait()
}

func RunServeAll(cmd *cobra.Command, args []string) {
	d := driver.New(cmd.Context(), driver.WithOptions(configx.WithFlags(cmd.Flags())))

	admin, public, adminmw, publicmw := setup(d, cmd)

	d.PrometheusManager().RegisterRouter(admin.Router)
	d.PrometheusManager().RegisterRouter(public.Router)

	var wg sync.WaitGroup
	wg.Add(2)

	go serve(
		d,
		cmd,
		&wg,
		config.PublicInterface,
		EnhanceMiddleware(d, publicmw, d.Config().ListenOn(config.PublicInterface), public.Router, false, config.PublicInterface),
		d.Config().ListenOn(config.PublicInterface),
		d.Config().SocketPermission(config.PublicInterface),
		GetOrCreateTLSCertificate(cmd, d, config.PublicInterface),
	)

	go serve(
		d,
		cmd,
		&wg,
		config.AdminInterface,
		EnhanceMiddleware(d, adminmw, d.Config().ListenOn(config.AdminInterface), admin.Router, true, config.AdminInterface),
		d.Config().ListenOn(config.AdminInterface),
		d.Config().SocketPermission(config.AdminInterface),
		GetOrCreateTLSCertificate(cmd, d, config.AdminInterface),
	)

	wg.Wait()
}

func setup(d driver.Registry, cmd *cobra.Command) (admin *x.RouterAdmin, public *x.RouterPublic, adminmw, publicmw *negroni.Negroni) {
	fmt.Println(banner(config.Version))

	if d.Config().CGroupsV1AutoMaxProcsEnabled() {
		_, err := maxprocs.Set(maxprocs.Logger(d.Logger().Infof))

		if err != nil {
			d.Logger().WithError(err).Fatal("Couldn't set GOMAXPROCS")
		}
	}

	adminmw = negroni.New()
	publicmw = negroni.New()

	admin = x.NewRouterAdmin()
	public = x.NewRouterPublic()

	if tracer := d.Tracer(cmd.Context()); tracer.IsLoaded() {
		adminmw.Use(tracer)
		publicmw.Use(tracer)
	}

	adminLogger := reqlog.
		NewMiddlewareFromLogger(d.Logger(),
			fmt.Sprintf("hydra/admin: %s", d.Config().IssuerURL().String()))
	if d.Config().DisableHealthAccessLog(config.AdminInterface) {
		adminLogger = adminLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	adminmw.Use(adminLogger)
	adminmw.Use(d.PrometheusManager())

	publicLogger := reqlog.NewMiddlewareFromLogger(
		d.Logger(),
		fmt.Sprintf("hydra/public: %s", d.Config().IssuerURL().String()),
	)
	if d.Config().DisableHealthAccessLog(config.PublicInterface) {
		publicLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	publicmw.Use(publicLogger)
	publicmw.Use(d.PrometheusManager())

	metrics := metricsx.New(
		cmd,
		d.Logger(),
		d.Config().Source(),
		&metricsx.Options{
			Service: "ory-hydra",
			ClusterID: metricsx.Hash(fmt.Sprintf("%s|%s",
				d.Config().IssuerURL().String(),
				d.Config().DSN(),
			)),
			IsDevelopment: d.Config().DSN() == "memory" ||
				d.Config().IssuerURL().String() == "" ||
				strings.Contains(d.Config().IssuerURL().String(), "localhost"),
			WriteKey: "h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr",
			WhitelistedPaths: []string{
				jwk.KeyHandlerPath,
				jwk.WellKnownKeysPath,

				client.ClientsHandlerPath,

				oauth2.DefaultConsentPath,
				oauth2.DefaultLoginPath,
				oauth2.DefaultPostLogoutPath,
				oauth2.DefaultLogoutPath,
				oauth2.DefaultErrorPath,
				oauth2.TokenPath,
				oauth2.AuthPath,
				oauth2.LogoutPath,
				oauth2.UserinfoPath,
				oauth2.WellKnownPath,
				oauth2.JWKPath,
				oauth2.IntrospectPath,
				oauth2.RevocationPath,
				oauth2.FlushPath,

				consent.ConsentPath,
				consent.ConsentPath + "/accept",
				consent.ConsentPath + "/reject",
				consent.LoginPath,
				consent.LoginPath + "/accept",
				consent.LoginPath + "/reject",
				consent.LogoutPath,
				consent.LogoutPath + "/accept",
				consent.LogoutPath + "/reject",
				consent.SessionsPath + "/login",
				consent.SessionsPath + "/consent",

				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,
				prometheus.MetricsPrometheusPath,
				"/",
			},
			BuildVersion: config.Version,
			BuildTime:    config.Date,
			BuildHash:    config.Commit,
			Config: &analytics.Config{
				Endpoint:             "https://sqa.ory.sh",
				GzipCompressionLevel: 6,
				BatchMaxSize:         500 * 1000,
				BatchSize:            250,
				Interval:             time.Hour * 24,
			},
		},
	)

	adminmw.Use(metrics)
	publicmw.Use(metrics)

	d.RegisterRoutes(admin, public)

	return
}

func serve(
	d driver.Registry,
	cmd *cobra.Command,
	wg *sync.WaitGroup,
	iface config.ServeInterface,
	handler http.Handler,
	address string,
	permission *configx.UnixPermission,
	cert []tls.Certificate,
) {
	defer wg.Done()

	var srv = graceful.WithDefaults(&http.Server{
		Handler: handler,
		// #nosec G402 - This is a false positive because we use graceful.WithDefaults which sets the correct TLS settings.
		TLSConfig: &tls.Config{
			Certificates: cert,
		},
	})

	if d.Tracer(cmd.Context()).IsLoaded() {
		srv.RegisterOnShutdown(d.Tracer(cmd.Context()).Close)
	}

	if err := graceful.Graceful(func() error {
		d.Logger().Infof("Setting up http server on %s", address)
		listener, err := networkx.MakeListener(address, permission)
		if err != nil {
			return err
		}

		if networkx.AddressIsUnixSocket(address) {
			return srv.Serve(listener)
		} else {
			tls := d.Config().TLS(iface)
			if !tls.Enabled() {
				if iface == config.PublicInterface {
					d.Logger().Warnln("HTTPS disabled. Never do this in production.")
				}
				return srv.Serve(listener)
			} else if len(tls.AllowTerminationFrom()) > 0 {
				d.Logger().Infoln("Upstream TLS termination enabled, disabling https.")
				return srv.Serve(listener)
			} else {
				return srv.ServeTLS(listener, "", "")
			}
		}
	}, srv.Shutdown); err != nil {
		d.Logger().WithError(err).Fatal("Could not gracefully run server")
	}
}
