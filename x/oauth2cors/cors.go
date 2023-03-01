// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2cors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"

	"github.com/gobwas/glob"
	"github.com/rs/cors"

	"github.com/ory/fosite"
)

func Middleware(
	ctx context.Context,
	reg interface {
		x.RegistryLogger
		oauth2.Registry
		client.Registry
	}) func(h http.Handler) http.Handler {
	opts, enabled := reg.Config().CORS(ctx, config.PublicInterface)
	if !enabled {
		return func(h http.Handler) http.Handler {
			return h
		}
	}

	var alwaysAllow = len(opts.AllowedOrigins) == 0
	var patterns []glob.Glob
	for _, o := range opts.AllowedOrigins {
		if o == "*" {
			alwaysAllow = true
		}
		// if the protocol (http or https) is specified, but the url is wildcard, use special ** glob, which ignore the '.' separator.
		// This way g := glob.Compile("http://**") g.Match("http://google.com") returns true.
		if splittedO := strings.Split(o, "://"); len(splittedO) != 1 && splittedO[1] == "*" {
			o = fmt.Sprintf("%s://**", splittedO[0])
		}
		g, err := glob.Compile(strings.ToLower(o), '.')
		if err != nil {
			reg.Logger().WithError(err).Fatalf("Unable to parse cors origin: %s", o)
		}

		patterns = append(patterns, g)
	}

	options := cors.Options{
		AllowedOrigins:     opts.AllowedOrigins,
		AllowedMethods:     opts.AllowedMethods,
		AllowedHeaders:     opts.AllowedHeaders,
		ExposedHeaders:     opts.ExposedHeaders,
		MaxAge:             opts.MaxAge,
		AllowCredentials:   opts.AllowCredentials,
		OptionsPassthrough: opts.OptionsPassthrough,
		Debug:              opts.Debug,
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			ctx := r.Context()
			if alwaysAllow {
				return true
			}

			origin = strings.ToLower(origin)
			for _, p := range patterns {
				if p.Match(origin) {
					return true
				}
			}

			// pre-flight requests do not contain credentials (cookies, HTTP authorization)
			// so we return true in all cases here.
			if r.Method == http.MethodOptions {
				return true
			}

			var clientID string

			// if the client uses client_secret_post auth it will provide its client ID in form data
			clientID = r.PostFormValue("client_id")

			// if the client uses client_secret_basic auth the client ID will be the username component
			if clientID == "" {
				clientID, _, _ = r.BasicAuth()
			}

			// otherwise, this may be a bearer auth request, in which case we can introspect the token
			if clientID == "" {
				token := fosite.AccessTokenFromRequest(r)
				if token == "" {
					return false
				}

				session := oauth2.NewSessionWithCustomClaims("", reg.Config().AllowedTopLevelClaims(ctx))
				_, ar, err := reg.OAuth2Provider().IntrospectToken(ctx, token, fosite.AccessToken, session)
				if err != nil {
					return false
				}

				clientID = ar.GetClient().GetID()
			}

			cl, err := reg.ClientManager().GetConcreteClient(ctx, clientID)
			if err != nil {
				return false
			}

			for _, o := range cl.AllowedCORSOrigins {
				if o == "*" {
					return true
				}
				g, err := glob.Compile(strings.ToLower(o), '.')
				if err != nil {
					return false
				}
				if g.Match(origin) {
					return true
				}
			}

			return false
		},
	}

	return cors.New(options).Handler
}
