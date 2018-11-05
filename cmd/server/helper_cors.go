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

package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gobwas/glob"
	"github.com/rs/cors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/oauth2"
)

func newCORSMiddleware(
	enable bool, c *config.Config, po cors.Options,
	o func(ctx context.Context, token string, tokenType fosite.TokenType, session fosite.Session, scope ...string) (fosite.TokenType, fosite.AccessRequester, error),
	clm func(ctx context.Context, id string) (*client.Client, error),
) func(h http.Handler) http.Handler {
	if !enable {
		return func(h http.Handler) http.Handler {
			return h
		}
	}

	c.GetLogger().Info("Enabled CORS")
	var patterns []glob.Glob
	for _, o := range po.AllowedOrigins {
		g, err := glob.Compile(strings.ToLower(o), '.')
		if err != nil {
			c.GetLogger().WithError(err).Fatalf("Unable to parse cors origin: %s", o)
		}
		patterns = append(patterns, g)
	}

	var alwaysAllow bool
	for _, o := range po.AllowedOrigins {
		if o == "*" {
			alwaysAllow = true
			break
		}
	}

	if enable && len(po.AllowedOrigins) == 0 {
		alwaysAllow = true
	}

	options := cors.Options{
		AllowedOrigins:     po.AllowedOrigins,
		AllowedMethods:     po.AllowedMethods,
		AllowedHeaders:     po.AllowedHeaders,
		ExposedHeaders:     po.ExposedHeaders,
		MaxAge:             po.MaxAge,
		AllowCredentials:   po.AllowCredentials,
		OptionsPassthrough: po.OptionsPassthrough,
		Debug:              po.Debug,
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
				_, ar, err := o(context.Background(), token, fosite.AccessToken, session)
				if err != nil {
					return false
				}

				username = ar.GetClient().GetID()
			}

			cl, err := clm(r.Context(), username)
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
