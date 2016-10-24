package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/meatballhat/negroni-logrus"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/policy"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
)

func RunHost(c *config.Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		router := httprouter.New()
		serverHandler := &Handler{Config: c}
		serverHandler.registerRoutes(router)
		c.ForceHTTP, _ = cmd.Flags().GetBool("dangerous-force-http")

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
			logrus.Warnln("Do not use flag --dangerous-auto-logon in production.")
			err := c.Persist()
			pkg.Must(err, "Could not write configuration file: %s", err)
		}

		n := negroni.New()
		n.Use(negronilogrus.NewMiddleware())
		n.UseFunc(serverHandler.rejectInsecureRequests)
		n.UseHandler(router)

		var srv = http.Server{
			Addr:    c.GetAddress(),
			Handler: n,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{
					getOrCreateTLSCertificate(cmd, c),
				},
			},
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 10,
		}

		var err error
		logrus.Infof("Setting up http server on %s", c.GetAddress())
		if c.ForceHTTP {
			logrus.Warnln("HTTPS disabled. Never do this in production.")
			err = srv.ListenAndServe()
		} else if c.AllowTLSTermination != "" {
			logrus.Infoln("TLS termination enabled, disabling https.")
			err = srv.ListenAndServe()
		} else {
			err = srv.ListenAndServeTLS("", "")
		}
		pkg.Must(err, "Could not start server: %s %s.", err)
	}
}

type Handler struct {
	Clients *client.Handler
	Keys    *jwk.Handler
	OAuth2  *oauth2.Handler
	Policy  *policy.Handler
	Warden  *warden.WardenHandler
	Config  *config.Config
}

func (h *Handler) registerRoutes(router *httprouter.Router) {
	c := h.Config
	ctx := c.Context()

	// Set up dependencies
	injectJWKManager(c)
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
	}

	// Set up handlers
	h.Clients = newClientHandler(c, router, clientsManager)
	h.Keys = newJWKHandler(c, router)
	h.Policy = newPolicyHandler(c, router)
	h.OAuth2 = newOAuth2Handler(c, router, ctx.KeyManager, oauth2Provider)
	h.Warden = warden.NewHandler(c, router)

	router.GET("/health", func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rw.WriteHeader(http.StatusNoContent)
	})

	// Create root account if new install
	createRS256KeysIfNotExist(c, oauth2.ConsentEndpointKey, "private", "sig")
	createRS256KeysIfNotExist(c, oauth2.ConsentChallengeKey, "private", "sig")

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
		logrus.WithError(err).Warnln("Could not serve http connection")
	}

	ans := new(herodot.JSON)
	ans.WriteErrorCode(context.Background(), rw, r, http.StatusBadGateway, errors.New("Can not serve request over insecure http"))
}
