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
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"

	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/corsx"
	"github.com/ory/x/serverx"
)

func injectJWKManager(c *config.Config) {
	ctx := c.Context()

	ctx.KeyManager = ctx.Connection.NewJWKManager(&jwk.AEAD{
		Key: c.GetSystemSecret(),
	})
}

func newJWKHandler(c *config.Config, frontend, backend *httprouter.Router, o fosite.OAuth2Provider, clm client.Manager) *jwk.Handler {
	ctx := c.Context()
	w := herodot.NewJSONWriter(c.GetLogger())
	w.ErrorEnhancer = serverx.ErrorEnhancerRFC6749
	var wellKnown []string

	if c.OAuth2AccessTokenStrategy == "jwt" {
		wellKnown = append(wellKnown, oauth2.OAuth2JWTKeyName)
	}

	expectDependency(c.GetLogger(), ctx.KeyManager)
	h := jwk.NewHandler(
		ctx.KeyManager,
		nil,
		w,
		wellKnown,
	)

	corsMiddleware := newCORSMiddleware(viper.GetString("CORS_ENABLED") == "true", c, corsx.ParseOptions(), o.IntrospectToken, clm.GetConcreteClient)
	h.SetRoutes(frontend, backend, corsMiddleware)
	return h
}
