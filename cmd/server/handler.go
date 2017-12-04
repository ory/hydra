// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/meatballhat/negroni-logrus"
	"github.com/ory/graceful"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/policy"
	"github.com/ory/hydra/warden"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

func parseCorsOptions() cors.Options {
	allowCredentials, _ := strconv.ParseBool(viper.GetString("CORS_ALLOWED_CREDENTIALS"))
	debug, _ := strconv.ParseBool(viper.GetString("CORS_DEBUG"))
	maxAge, _ := strconv.Atoi(viper.GetString("CORS_MAX_AGE"))
	return cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ","),
		AllowedHeaders:   strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ","),
		ExposedHeaders:   strings.Split(viper.GetString("CORS_EXPOSED_HEADERS"), ","),
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
		Debug:            debug,
	}
}

func RunHost(c *config.Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		router := httprouter.New()
		logger := c.GetLogger()
		serverHandler := &Handler{
			Config: c,
			H:      herodot.NewJSONWriter(logger),
		}
		serverHandler.registerRoutes(router)
		c.ForceHTTP, _ = cmd.Flags().GetBool("dangerous-force-http")

		if !c.ForceHTTP {
			if c.Issuer == "" {
				logger.Fatalln("Issuer must be explicitly specified unless --dangerous-force-http is passed. To find out more, use `hydra help host`.")
			}
			issuer, err := url.Parse(c.Issuer)
			pkg.Must(err, "Could not parse issuer URL: %s", err)
			if issuer.Scheme != "https" {
				logger.Fatalln("Issuer must use HTTPS unless --dangerous-force-http is passed. To find out more, use `hydra help host`.")
			}
		}

		if c.ClusterURL == "" {
			proto := "https"
			if c.ForceHTTP {
				proto = "http"
			}
			host := "localhost"
			if c.BindHost != "" {
				host = c.BindHost
			}
			c.ClusterURL = fmt.Sprintf("%s://%s:%d", proto, host, c.BindPort)
		}

		if ok, _ := cmd.Flags().GetBool("dangerous-auto-logon"); ok {
			logger.Warnln("Do not use flag --dangerous-auto-logon in production.")
			err := c.Persist()
			pkg.Must(err, "Could not write configuration file: %s", err)
		}

		n := negroni.New()

		metrics := c.GetMetrics()
		if ok, _ := cmd.Flags().GetBool("disable-telemetry"); !ok && os.Getenv("DISABLE_TELEMETRY") != "1" {
			go metrics.RegisterSegment(c.BuildVersion, c.BuildHash, c.BuildTime)
			go metrics.CommitTelemetry()
			go metrics.TickKeepAlive()
			n.Use(metrics)
		}

		n.Use(negronilogrus.NewMiddlewareFromLogger(logger, c.Issuer))
		n.UseFunc(serverHandler.rejectInsecureRequests)
		n.UseHandler(router)
		corsHandler := cors.New(parseCorsOptions()).Handler(n)

		var srv = graceful.WithDefaults(&http.Server{
			Addr:    c.GetAddress(),
			Handler: context.ClearHandler(corsHandler),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{getOrCreateTLSCertificate(cmd, c)},
			},
		})

		err := graceful.Graceful(func() error {
			var err error
			logger.Infof("Setting up http server on %s", c.GetAddress())
			if c.ForceHTTP {
				logger.Warnln("HTTPS disabled. Never do this in production.")
				err = srv.ListenAndServe()
			} else if c.AllowTLSTermination != "" {
				logger.Infoln("TLS termination enabled, disabling https.")
				err = srv.ListenAndServe()
			} else {
				err = srv.ListenAndServeTLS("", "")
			}

			return err
		}, srv.Shutdown)
		logger.WithError(err).Fatal("Could not gracefully run server")
	}
}

type Handler struct {
	Clients *client.Handler
	Keys    *jwk.Handler
	OAuth2  *oauth2.Handler
	Consent *oauth2.ConsentSessionHandler
	Policy  *policy.Handler
	Groups  *group.Handler
	Warden  *warden.WardenHandler
	Config  *config.Config
	H       herodot.Writer
}

func (h *Handler) registerRoutes(router *httprouter.Router) {
	c := h.Config
	ctx := c.Context()

	// Set up dependencies
	injectJWKManager(c)
	injectConsentManager(c)
	clientsManager := newClientManager(c)
	injectFositeStore(c, clientsManager)
	oauth2Provider := newOAuth2Provider(c, ctx.KeyManager)

	// set up warden
	ctx.Warden = &warden.LocalWarden{
		Warden: &ladon.Ladon{
			Manager: ctx.LadonManager,
		},
		OAuth2:              oauth2Provider,
		Issuer:              c.Issuer,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
		Groups:              ctx.GroupManager,
		L:                   c.GetLogger(),
	}

	// Set up handlers
	h.Clients = newClientHandler(c, router, clientsManager)
	h.Keys = newJWKHandler(c, router)
	h.Policy = newPolicyHandler(c, router)
	h.Consent = newConsentHanlder(c, router)
	h.OAuth2 = newOAuth2Handler(c, router, ctx.ConsentManager, oauth2Provider)
	h.Warden = warden.NewHandler(c, router)
	h.Groups = &group.Handler{
		H:              herodot.NewJSONWriter(c.GetLogger()),
		W:              ctx.Warden,
		Manager:        ctx.GroupManager,
		ResourcePrefix: c.AccessControlResourcePrefix,
	}
	h.Groups.SetRoutes(router)
	_ = newHealthHandler(c, router)

	h.createRootIfNewInstall(c)
}

func (h *Handler) rejectInsecureRequests(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
