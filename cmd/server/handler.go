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
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"

	"github.com/ory/graceful"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/corsx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/healthx"
	"github.com/ory/x/metricsx"
)

var _ = &consent.Handler{}

func EnhanceRouter(c *config.Config, cmd *cobra.Command, serverHandler *Handler, router *httprouter.Router, middlewares []negroni.Handler, enableCors bool) http.Handler {
	n := negroni.New()
	for _, m := range middlewares {
		n.Use(m)
	}
	n.UseFunc(serverHandler.RejectInsecureRequests)
	n.UseHandler(router)
	if enableCors {
		c.GetLogger().Info("Enabled CORS")
		options := corsx.ParseOptions()
		return context.ClearHandler(cors.New(options).Handler(n))
	} else {
		return context.ClearHandler(n)
	}
}

func RunServeAdmin(c *config.Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		c.MustValidate()
		checkDatabaseAllowed(c)
		serverHandler, _, backend, mws := setup(c, cmd, args, "admin")

		var wg sync.WaitGroup
		wg.Add(2)

		cert := getOrCreateTLSCertificate(cmd, c)
		// go serve(c, cmd, enhanceRouter(c, cmd, serverHandler, frontend), c.GetFrontendAddress(), &wg)
		go serve(c, cmd, EnhanceRouter(c, cmd, serverHandler, backend, mws, viper.GetString("CORS_ENABLED") == "true"), c.GetBackendAddress(), &wg, cert)

		wg.Wait()
	}
}

func RunServePublic(c *config.Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		c.MustValidate()
		checkDatabaseAllowed(c)
		serverHandler, frontend, _, mws := setup(c, cmd, args, "public")

		var wg sync.WaitGroup
		wg.Add(2)

		cert := getOrCreateTLSCertificate(cmd, c)
		go serve(c, cmd, EnhanceRouter(c, cmd, serverHandler, frontend, mws, false), c.GetFrontendAddress(), &wg, cert)
		// go serve(c, cmd, enhanceRouter(c, cmd, serverHandler, backend), c.GetBackendAddress(), &wg)

		wg.Wait()
	}
}

func RunServeAll(c *config.Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		c.MustValidate()
		serverHandler, frontend, backend, mws := setup(c, cmd, args, "all")

		var wg sync.WaitGroup
		wg.Add(2)

		cert := getOrCreateTLSCertificate(cmd, c)
		go serve(c, cmd, EnhanceRouter(c, cmd, serverHandler, frontend, mws, false), c.GetFrontendAddress(), &wg, cert)
		go serve(c, cmd, EnhanceRouter(c, cmd, serverHandler, backend, mws, viper.GetString("CORS_ENABLED") == "true"), c.GetBackendAddress(), &wg, cert)

		wg.Wait()
	}
}

func setup(c *config.Config, cmd *cobra.Command, args []string, name string) (handler *Handler, frontend, backend *httprouter.Router, middlewares []negroni.Handler) {
	fmt.Println(banner(c.BuildVersion))

	frontend = httprouter.New()
	backend = httprouter.New()

	logger := c.GetLogger()
	w := newJSONWriter(logger)

	if tracer, err := c.GetTracer(); err != nil {
		c.GetLogger().Fatalf("Failed to initialize tracer: %s", err)
	} else if tracer.IsLoaded() {
		middlewares = append(middlewares, tracer)
	}

	handler = NewHandler(c, w)
	handler.RegisterRoutes(frontend, backend)
	c.ForceHTTP = flagx.MustGetBool(cmd, "dangerous-force-http")

	if !c.ForceHTTP {
		if c.Issuer == "" {
			logger.Fatalln("OAUTH2_ISSUER_URL must be explicitly specified unless --dangerous-force-http is passed. To find out more, use `hydra help serve`.")
		}
		issuer, err := url.Parse(c.Issuer)
		cmdx.Must(err, "Could not parse issuer URL: %s", err)
		if issuer.Scheme != "https" {
			logger.Fatalln("OAUTH2_ISSUER_URL must use HTTPS unless --dangerous-force-http is passed. To find out more, use `hydra help serve`.")
		}
	}

	middlewares = append(
		middlewares,
		negronilogrus.NewMiddlewareFromLogger(c.GetLogger(), c.Issuer),
		c.GetPrometheusMetrics(),
	)

	if flagx.MustGetBool(cmd, "disable-telemetry") {
		c.GetLogger().Println("Transmission of telemetry data is enabled, to learn more go to: https://www.ory.sh/docs/guides/latest/telemetry/")

		enable := !(c.DatabaseURL == "" || c.DatabaseURL == "memory" || c.Issuer == "" || strings.Contains(c.Issuer, "localhost"))
		m := metricsx.NewMetricsManager(
			metricsx.Hash(c.Issuer+"|"+c.DatabaseURL),
			enable,
			"h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr",
			[]string{
				client.ClientsHandlerPath,
				jwk.KeyHandlerPath,
				jwk.WellKnownKeysPath,
				oauth2.DefaultConsentPath,
				oauth2.TokenPath,
				oauth2.AuthPath,
				oauth2.UserinfoPath,
				oauth2.WellKnownPath,
				oauth2.IntrospectPath,
				oauth2.RevocationPath,
				consent.ConsentPath,
				consent.LoginPath,
				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,
				metricsPrometheusPath,
				"/oauth2/auth/sessions/login",
				"/oauth2/auth/sessions/consent",
				"/",
			},
			c.GetLogger(),
			"ory-hydra",
			100,
			"",
		)

		go m.RegisterSegment(c.BuildVersion, c.BuildHash, c.BuildTime)
		go m.CommitMemoryStatistics()

		middlewares = append(middlewares, m)
	}

	return
}

func checkDatabaseAllowed(c *config.Config) {
	if c.DatabaseURL == "memory" {
		c.GetLogger().Fatalf(`When using "hydra serve admin" or "hydra serve public" the DATABASE_URL can not be set to "memory".`)
	}
}

func serve(c *config.Config, cmd *cobra.Command, handler http.Handler, address string, wg *sync.WaitGroup, cert []tls.Certificate) {
	defer wg.Done()

	var srv = graceful.WithDefaults(&http.Server{
		Addr:    address,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: cert,
		},
	})

	if tracer, err := c.GetTracer(); err == nil && tracer.IsLoaded() {
		srv.RegisterOnShutdown(tracer.Close)
	}

	err := graceful.Graceful(func() error {
		var err error
		c.GetLogger().Infof("Setting up http server on %s", address)
		if c.ForceHTTP {
			c.GetLogger().Warnln("HTTPS disabled. Never do this in production.")
			err = srv.ListenAndServe()
		} else if c.AllowTLSTermination != "" {
			c.GetLogger().Infoln("TLS termination enabled, disabling https.")
			err = srv.ListenAndServe()
		} else {
			err = srv.ListenAndServeTLS("", "")
		}

		return err
	}, srv.Shutdown)
	if err != nil {
		c.GetLogger().WithError(err).Fatal("Could not gracefully run server")
	}
}

type Handler struct {
	Clients *client.Handler
	Keys    *jwk.Handler
	OAuth2  *oauth2.Handler
	Consent *consent.Handler
	Config  *config.Config
	H       herodot.Writer
}

func NewHandler(c *config.Config, h herodot.Writer) *Handler {
	return &Handler{Config: c, H: h}
}

func (h *Handler) RegisterRoutes(frontend, backend *httprouter.Router) {
	c := h.Config
	ctx := c.Context()

	// Set up dependencies
	injectJWKManager(c)
	clientsManager := newClientManager(c)

	injectFositeStore(c, clientsManager)
	injectConsentManager(c, clientsManager)

	oauth2Provider := newOAuth2Provider(c)

	// Set up handlers
	h.Clients = newClientHandler(c, backend, clientsManager)
	h.Keys = newJWKHandler(c, frontend, backend, oauth2Provider, clientsManager)
	h.Consent = newConsentHandler(c, frontend, backend)
	h.OAuth2 = newOAuth2Handler(c, frontend, backend, ctx.ConsentManager, oauth2Provider, clientsManager)
	healthHandler := newHealthHandler(c, backend)

	// Register AliveCheckPath for frontend too
	frontend.GET(healthx.AliveCheckPath, healthHandler.Alive)
}

func (h *Handler) RejectInsecureRequests(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.TLS != nil || h.Config.ForceHTTP {
		next.ServeHTTP(rw, r)
		return
	}

	if err := h.Config.DoesRequestSatisfyTermination(r); err == nil {
		next.ServeHTTP(rw, r)
		return
	} else {
		h.Config.GetLogger().WithError(err).Warnln("Could not serve http connection")
	}

	h.H.WriteErrorCode(rw, r, http.StatusBadGateway, errors.New("Can not serve request over insecure http"))
}
