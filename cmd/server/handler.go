// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/ory/x/otelx/semconv"

	"github.com/ory/x/servicelocatorx"

	"github.com/ory/x/httprouterx"

	"github.com/ory/analytics-go/v5"
	"github.com/ory/x/configx"

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
	"github.com/ory/x/otelx"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	prometheus "github.com/ory/x/prometheusx"
)

var _ = &consent.Handler{}

func EnhanceMiddleware(ctx context.Context, sl *servicelocatorx.Options, d driver.Registry, n *negroni.Negroni, address string, router *httprouter.Router, iface config.ServeInterface) http.Handler {
	if !networkx.AddressIsUnixSocket(address) {
		n.UseFunc(x.RejectInsecureRequests(d, d.Config().TLS(ctx, iface)))
	}

	for _, mw := range sl.HTTPMiddlewares() {
		n.UseFunc(mw)
	}
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cfg, enabled := d.Config().CORS(r.Context(), iface)
		if !enabled {
			next(w, r)
			return
		}
		cors.New(cfg).ServeHTTP(w, r, next)
	})

	n.UseHandler(router)

	return n
}

func ensureNoMemoryDSN(r driver.Registry) {
	if r.Config().DSN() == "memory" {
		r.Logger().Fatalf(`When using "hydra serve admin" or "hydra serve public" the DSN can not be set to "memory".`)
	}
}

func RunServeAdmin(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		sl := servicelocatorx.NewOptions(slOpts...)

		d, err := driver.New(cmd.Context(), sl, append(dOpts, driver.WithOptions(append(cOpts, configx.WithFlags(cmd.Flags()))...)))
		if err != nil {
			return err
		}
		ensureNoMemoryDSN(d)

		admin, _, adminmw, _ := setup(ctx, d, cmd)
		d.PrometheusManager().RegisterRouter(admin.Router)

		var wg sync.WaitGroup
		wg.Add(1)

		go serve(
			ctx,
			d,
			cmd,
			&wg,
			config.AdminInterface,
			EnhanceMiddleware(ctx, sl, d, adminmw, d.Config().ListenOn(config.AdminInterface), admin.Router, config.AdminInterface),
			d.Config().ListenOn(config.AdminInterface),
			d.Config().SocketPermission(config.AdminInterface),
		)

		wg.Wait()
		return nil
	}
}

func RunServePublic(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		sl := servicelocatorx.NewOptions(slOpts...)

		d, err := driver.New(cmd.Context(), sl, append(dOpts, driver.WithOptions(append(cOpts, configx.WithFlags(cmd.Flags()))...)))
		if err != nil {
			return err
		}
		ensureNoMemoryDSN(d)

		_, public, _, publicmw := setup(ctx, d, cmd)
		d.PrometheusManager().RegisterRouter(public.Router)

		var wg sync.WaitGroup
		wg.Add(1)

		go serve(
			ctx,
			d,
			cmd,
			&wg,
			config.PublicInterface,
			EnhanceMiddleware(ctx, sl, d, publicmw, d.Config().ListenOn(config.PublicInterface), public.Router, config.PublicInterface),
			d.Config().ListenOn(config.PublicInterface),
			d.Config().SocketPermission(config.PublicInterface),
		)

		wg.Wait()
		return nil
	}
}

func RunServeAll(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		sl := servicelocatorx.NewOptions(slOpts...)

		d, err := driver.New(cmd.Context(), sl, append(dOpts, driver.WithOptions(append(cOpts, configx.WithFlags(cmd.Flags()))...)))
		if err != nil {
			return err
		}

		admin, public, adminmw, publicmw := setup(ctx, d, cmd)

		d.PrometheusManager().RegisterRouter(admin.Router)
		d.PrometheusManager().RegisterRouter(public.Router)

		var wg sync.WaitGroup
		wg.Add(2)

		go serve(
			ctx,
			d,
			cmd,
			&wg,
			config.PublicInterface,
			EnhanceMiddleware(ctx, sl, d, publicmw, d.Config().ListenOn(config.PublicInterface), public.Router, config.PublicInterface),
			d.Config().ListenOn(config.PublicInterface),
			d.Config().SocketPermission(config.PublicInterface),
		)

		go serve(
			ctx,
			d,
			cmd,
			&wg,
			config.AdminInterface,
			EnhanceMiddleware(ctx, sl, d, adminmw, d.Config().ListenOn(config.AdminInterface), admin.Router, config.AdminInterface),
			d.Config().ListenOn(config.AdminInterface),
			d.Config().SocketPermission(config.AdminInterface),
		)

		wg.Wait()
		return nil
	}
}

func setup(ctx context.Context, d driver.Registry, cmd *cobra.Command) (admin *httprouterx.RouterAdmin, public *httprouterx.RouterPublic, adminmw, publicmw *negroni.Negroni) {
	fmt.Println(banner(config.Version))

	if d.Config().CGroupsV1AutoMaxProcsEnabled() {
		_, err := maxprocs.Set(maxprocs.Logger(d.Logger().Infof))

		if err != nil {
			d.Logger().WithError(err).Fatal("Couldn't set GOMAXPROCS")
		}
	}

	adminmw = negroni.New()
	publicmw = negroni.New()

	admin = x.NewRouterAdmin(d.Config().AdminURL)
	public = x.NewRouterPublic()

	adminLogger := reqlog.
		NewMiddlewareFromLogger(d.Logger(),
			fmt.Sprintf("hydra/admin: %s", d.Config().IssuerURL(ctx).String()))
	if d.Config().DisableHealthAccessLog(config.AdminInterface) {
		adminLogger = adminLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath, "/admin"+prometheus.MetricsPrometheusPath)
	}

	adminmw.UseFunc(semconv.Middleware)
	adminmw.Use(adminLogger)
	adminmw.Use(d.PrometheusManager())

	publicLogger := reqlog.NewMiddlewareFromLogger(
		d.Logger(),
		fmt.Sprintf("hydra/public: %s", d.Config().IssuerURL(ctx).String()),
	)
	if d.Config().DisableHealthAccessLog(config.PublicInterface) {
		publicLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	publicmw.UseFunc(semconv.Middleware)
	publicmw.Use(publicLogger)
	publicmw.Use(d.PrometheusManager())

	metrics := metricsx.New(
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

				"/admin" + client.ClientsHandlerPath,
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
				prometheus.MetricsPrometheusPath,
				"/admin" + prometheus.MetricsPrometheusPath,
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
		},
	)

	adminmw.Use(metrics)
	publicmw.Use(metrics)

	d.RegisterRoutes(ctx, admin, public)

	return
}

func serve(
	ctx context.Context,
	d driver.Registry,
	cmd *cobra.Command,
	wg *sync.WaitGroup,
	iface config.ServeInterface,
	handler http.Handler,
	address string,
	permission *configx.UnixPermission,
) {
	defer wg.Done()

	if tracer := d.Tracer(cmd.Context()); tracer.IsLoaded() {
		handler = otelx.TraceHandler(
			handler,
			otelhttp.WithTracerProvider(tracer.Provider()),
			otelhttp.WithFilter(func(r *http.Request) bool {
				return !strings.HasPrefix(r.URL.Path, "/admin/metrics/")
			}),
		)
	}

	var tlsConfig *tls.Config
	stopReload := make(chan struct{})
	if tc := d.Config().TLS(ctx, iface); tc.Enabled() {
		// #nosec G402 - This is a false positive because we use graceful.WithDefaults which sets the correct TLS settings.
		tlsConfig = &tls.Config{GetCertificate: GetOrCreateTLSCertificate(ctx, d, iface, stopReload)}
	}

	var srv = graceful.WithDefaults(&http.Server{
		Handler:           handler,
		TLSConfig:         tlsConfig,
		ReadHeaderTimeout: time.Second * 5,
	})

	if err := graceful.Graceful(func() error {
		d.Logger().Infof("Setting up http server on %s", address)
		listener, err := networkx.MakeListener(address, permission)
		if err != nil {
			return err
		}

		if networkx.AddressIsUnixSocket(address) {
			return srv.Serve(listener)
		}

		if tlsConfig != nil {
			return srv.ServeTLS(listener, "", "")
		}

		if iface == config.PublicInterface {
			d.Logger().Warnln("HTTPS is disabled. Please ensure that your proxy is configured to provide HTTPS, and that it redirects HTTP to HTTPS.")
		}

		return srv.Serve(listener)
	}, func(ctx context.Context) error {
		close(stopReload)
		return srv.Shutdown(ctx)
	}); err != nil {
		d.Logger().WithError(err).Fatal("Could not gracefully run server")
	}
}
