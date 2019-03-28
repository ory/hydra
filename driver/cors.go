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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package driver

import (
	"context"
	"net/http"
	"strings"

	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/oauth2"

	"github.com/gobwas/glob"
	"github.com/rs/cors"

	"github.com/ory/fosite"
)

func OAuth2AwareCORSMiddleware(iface string, reg Registry, conf configuration.Provider) func(h http.Handler) http.Handler {
	if !conf.CORSEnabled(iface) {
		return func(h http.Handler) http.Handler {
			return h
		}
	}

	corsOptions := conf.CORSOptions(iface)
	var patterns []glob.Glob
	for _, o := range corsOptions.AllowedOrigins {
		g, err := glob.Compile(strings.ToLower(o), '.')
		if err != nil {
			reg.Logger().WithError(err).Fatalf("Unable to parse cors origin: %s", o)
		}
		patterns = append(patterns, g)
	}

	var alwaysAllow bool
	for _, o := range corsOptions.AllowedOrigins {
		if o == "*" {
			alwaysAllow = true
			break
		}
	}

	if len(corsOptions.AllowedOrigins) == 0 {
		alwaysAllow = true
	}

	options := cors.Options{
		AllowedOrigins:     corsOptions.AllowedOrigins,
		AllowedMethods:     corsOptions.AllowedMethods,
		AllowedHeaders:     corsOptions.AllowedHeaders,
		ExposedHeaders:     corsOptions.ExposedHeaders,
		MaxAge:             corsOptions.MaxAge,
		AllowCredentials:   corsOptions.AllowCredentials,
		OptionsPassthrough: corsOptions.OptionsPassthrough,
		Debug:              corsOptions.Debug,
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			if alwaysAllow {
				return true
			}

			origin = strings.ToLower(origin)
			for _, p := range patterns {
				if p.Match(origin) {
					return true
				}
			}

			username, _, ok := r.BasicAuth()
			if !ok || username == "" {
				token := fosite.AccessTokenFromRequest(r)
				if token == "" {
					return false
				}

				session := oauth2.NewSession("")
				_, ar, err := reg.OAuth2Provider().IntrospectToken(context.Background(), token, fosite.AccessToken, session)
				if err != nil {
					return false
				}

				username = ar.GetClient().GetID()
			}

			cl, err := reg.ClientManager().GetConcreteClient(r.Context(), username)
			if err != nil {
				return false
			}

			if alwaysAllow {
				return true
			}

			for _, p := range cl.AllowedCORSOrigins {
				if p == "*" {
					return true
				}
			}

			var clientPatterns []glob.Glob
			for _, o := range cl.AllowedCORSOrigins {
				g, err := glob.Compile(strings.ToLower(o), '.')
				if err != nil {
					return false
				}
				clientPatterns = append(patterns, g)
			}

			for _, p := range clientPatterns {
				if p.Match(origin) {
					return true
				}
			}

			return false
		},
	}

	return cors.New(options).Handler
}
