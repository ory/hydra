// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"

	"github.com/ory/analytics-go/v5"
	"github.com/ory/graceful"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/healthx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/metricsx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/otelx/semconv"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/reqlog"
	"github.com/ory/x/tlsx"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
)

func ensureNoMemoryDSN(r *driver.RegistrySQL) {
	if r.Config().DSN() == "memory" {
		r.Logger().Fatalf(`When using "hydra serve admin" or "hydra serve public" the DSN can not be set to "memory".`)
	}
}

func RunServeAdmin(dOpts []driver.OptionsModifier) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println(banner(config.Version))

		ctx := cmd.Context()

		d, err := driver.New(ctx, append(dOpts, driver.WithConfigOptions(configx.WithFlags(cmd.Flags())))...)
		if err != nil {
			return err
		}
		ensureNoMemoryDSN(d)

		srv, err := adminServer(ctx, d, sqa(ctx, d, cmd))
		if err != nil {
			return err
		}
		return srv()
	}
}

func RunServePublic(dOpts []driver.OptionsModifier) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println(banner(config.Version))

		ctx := cmd.Context()

		d, err := driver.New(ctx, append(dOpts, driver.WithConfigOptions(configx.WithFlags(cmd.Flags())))...)
		if err != nil {
			return err
		}
		ensureNoMemoryDSN(d)

		srv, err := publicServer(ctx, d, sqa(ctx, d, cmd))
		if err != nil {
			return err
		}
		return srv()
	}
}

func RunServeAll(dOpts []driver.OptionsModifier) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println(banner(config.Version))

		ctx := cmd.Context()

		d, err := driver.New(ctx, append(dOpts, driver.WithConfigOptions(configx.WithFlags(cmd.Flags())))...)
		if err != nil {
			return err
		}

		eg, ctx := errgroup.WithContext(ctx)
		ms := sqa(ctx, d, cmd)

		srvAdmin, err := adminServer(ctx, d, ms)
		if err != nil {
			return err
		}
		srvPublic, err := publicServer(ctx, d, ms)
		if err != nil {
			return err
		}

		eg.Go(srvAdmin)
		eg.Go(srvPublic)
		return eg.Wait()
	}
}

var prometheusManager = prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)

func adminServer(ctx context.Context, d *driver.RegistrySQL, sqaMetrics *metricsx.Service) (func() error, error) {
	cfg := d.Config().ServeAdmin(contextx.RootContext)

	n := negroni.New()

	logger := reqlog.
		NewMiddlewareFromLogger(d.Logger(),
			fmt.Sprintf("hydra/admin: %s", d.Config().IssuerURL(ctx).String()))
	if cfg.RequestLog.DisableHealth {
		logger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath, "/admin"+prometheusx.MetricsPrometheusPath)
	}

	n.UseFunc(httprouterx.TrimTrailingSlashNegroni)
	n.UseFunc(httprouterx.NoCacheNegroni)
	n.UseFunc(httprouterx.AddAdminPrefixIfNotPresentNegroni)
	n.UseFunc(semconv.Middleware)
	n.Use(logger)

	if cfg.TLS.Enabled && !networkx.AddressIsUnixSocket(cfg.Host) {
		mw, err := tlsx.EnforceTLSRequests(d, cfg.TLS.AllowTerminationFrom)
		if err != nil {
			return nil, err
		}
		n.Use(mw)
	}

	for _, mw := range d.HTTPMiddlewares() {
		n.Use(mw)
	}
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cfg, enabled := d.Config().CORSAdmin(r.Context())
		if !enabled {
			next(w, r)
			return
		}
		cors.New(cfg).ServeHTTP(w, r, next)
	})
	n.Use(sqaMetrics)

	router := httprouterx.NewRouterAdminWithPrefix(prometheusManager)
	d.RegisterAdminRoutes(router)

	n.UseHandler(router)

	return func() error {
		return serve(ctx, d, cfg, n, "admin")
	}, nil
}

func publicServer(ctx context.Context, d *driver.RegistrySQL, sqaMetrics *metricsx.Service) (func() error, error) {
	cfg := d.Config().ServePublic(contextx.RootContext)

	n := negroni.New()

	logger := reqlog.NewMiddlewareFromLogger(
		d.Logger(),
		fmt.Sprintf("hydra/public: %s", d.Config().IssuerURL(ctx).String()),
	)
	if cfg.RequestLog.DisableHealth {
		logger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	n.UseFunc(httprouterx.TrimTrailingSlashNegroni)
	n.UseFunc(httprouterx.NoCacheNegroni)
	n.UseFunc(semconv.Middleware)
	n.Use(logger)
	if cfg.TLS.Enabled && !networkx.AddressIsUnixSocket(cfg.Host) {
		mw, err := tlsx.EnforceTLSRequests(d, cfg.TLS.AllowTerminationFrom)
		if err != nil {
			return nil, err
		}
		n.Use(mw)
	}

	for _, mw := range d.HTTPMiddlewares() {
		n.Use(mw)
	}
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cfg, enabled := d.Config().CORSPublic(r.Context())
		if !enabled {
			next(w, r)
			return
		}
		cors.New(cfg).ServeHTTP(w, r, next)
	})
	n.Use(sqaMetrics)

	router := x.NewRouterPublic(prometheusManager)
	d.RegisterPublicRoutes(ctx, router)

	n.UseHandler(router)
	return func() error {
		return serve(ctx, d, cfg, n, "public")
	}, nil
}

func sqa(ctx context.Context, d *driver.RegistrySQL, cmd *cobra.Command) *metricsx.Service {
	urls := []string{
		d.Config().IssuerURL(ctx).Host,
		d.Config().PublicURL(ctx).Host,
		d.Config().AdminURL(ctx).Host,
		d.Config().ServePublic(ctx).BaseURL.Host,
		d.Config().ServeAdmin(ctx).BaseURL.Host,
		d.Config().LoginURL(ctx).Host,
		d.Config().LogoutURL(ctx).Host,
		d.Config().ConsentURL(ctx).Host,
		d.Config().RegistrationURL(ctx).Host,
	}
	if c, y := d.Config().CORSPublic(ctx); y {
		urls = append(urls, c.AllowedOrigins...)
	}
	if c, y := d.Config().CORSAdmin(ctx); y {
		urls = append(urls, c.AllowedOrigins...)
	}
	host := urlx.ExtractPublicAddress(urls...)

	return metricsx.New(
		cmd,
		d.Logger(),
		d.Config().Source(ctx),
		&metricsx.Options{
			Service:      "hydra",
			DeploymentId: metricsx.Hash(d.Persister().NetworkID(ctx).String()),
			IsDevelopment: d.Config().DSN() == "memory" ||
				d.Config().IssuerURL(ctx).String() == "" ||
				strings.Contains(d.Config().IssuerURL(ctx).String(), "localhost"),
			WriteKey: "h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr",
			WhitelistedPaths: []string{
				"/admin" + jwk.KeyHandlerPath,
				jwk.WellKnownKeysPath,

				urlx.MustJoin("/admin", client.ClientsHandlerPath),
				client.DynClientsHandlerPath,

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
				"/admin" + oauth2.IntrospectPath,
				"/admin" + oauth2.DeleteTokensPath,
				oauth2.RevocationPath,

				"/admin" + consent.ConsentPath,
				"/admin" + consent.ConsentPath + "/accept",
				"/admin" + consent.ConsentPath + "/reject",
				"/admin" + consent.LoginPath,
				"/admin" + consent.LoginPath + "/accept",
				"/admin" + consent.LoginPath + "/reject",
				"/admin" + consent.LogoutPath,
				"/admin" + consent.LogoutPath + "/accept",
				"/admin" + consent.LogoutPath + "/reject",
				"/admin" + consent.SessionsPath + "/login",
				"/admin" + consent.SessionsPath + "/consent",

				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				"/admin" + healthx.AliveCheckPath,
				"/admin" + healthx.ReadyCheckPath,
				healthx.VersionPath,
				"/admin" + healthx.VersionPath,
				prometheusx.MetricsPrometheusPath,
				"/admin" + prometheusx.MetricsPrometheusPath,
				"/",
			},
			BuildVersion: config.Version,
			BuildTime:    config.Date,
			BuildHash:    config.Commit,
			Config: &analytics.Config{
				Endpoint:             "https://sqa.ory.sh",
				GzipCompressionLevel: 6,
				BatchMaxSize:         500 * 1000,
				BatchSize:            1000,
				Interval:             time.Hour * 6,
			},
			Hostname: host,
		},
	)
}

func serve(
	ctx context.Context,
	d *driver.RegistrySQL,
	cfg *configx.Serve,
	handler http.Handler,
	ifaceName string,
) error {
	if tracer := d.Tracer(ctx); tracer.IsLoaded() {
		handler = otelx.TraceHandler(
			handler,
			otelhttp.WithTracerProvider(tracer.Provider()),
			otelhttp.WithFilter(func(r *http.Request) bool {
				return !strings.HasPrefix(r.URL.Path, "/admin/metrics/")
			}),
		)
	}

	var tlsConfig *tls.Config
	if cfg.TLS.Enabled {
		// #nosec G402 - This is a false positive because we use graceful.WithDefaults which sets the correct TLS settings.
		tlsConfig = &tls.Config{GetCertificate: GetOrCreateTLSCertificate(ctx, d, cfg.TLS, ifaceName)}
	}

	srv := graceful.WithDefaults(&http.Server{
		Handler:           handler,
		TLSConfig:         tlsConfig,
		ReadHeaderTimeout: time.Second * 5,
	})

	addr := configx.GetAddress(cfg.Host, cfg.Port)
	return graceful.Graceful(func() error {
		d.Logger().Infof("Setting up http server on %s", addr)
		listener, err := networkx.MakeListener(addr, &cfg.Socket)
		if err != nil {
			return err
		}

		if networkx.AddressIsUnixSocket(addr) {
			return srv.Serve(listener)
		}

		if tlsConfig != nil {
			return srv.ServeTLS(listener, "", "")
		}

		d.Logger().Warnln("HTTPS is disabled. Please ensure that your proxy is configured to provide HTTPS, and that it redirects HTTP to HTTPS.")

		return srv.Serve(listener)
	}, srv.Shutdown)
}
