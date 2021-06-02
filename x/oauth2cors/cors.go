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

package oauth2cors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/x"

	"github.com/ory/hydra/oauth2"

	"github.com/gobwas/glob"
	"github.com/rs/cors"

	"github.com/ory/fosite"
)

func Middleware(reg interface {
	Config() *config.Provider
	x.RegistryLogger
	oauth2.Registry
	client.Registry
}) func(h http.Handler) http.Handler {
	opts, enabled := reg.Config().CORS(config.PublicInterface)
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

				session := oauth2.NewSessionWithCustomClaims("", reg.Config().AllowedTopLevelClaims())
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
