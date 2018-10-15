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
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/consent"
)

func injectConsentManager(c *config.Config, cm client.Manager) {
	var ctx = c.Context()
	ctx.ConsentManager = ctx.Connection.NewConsentManager(cm, ctx.FositeStore)
}

func newConsentHandler(c *config.Config, frontend, backend *httprouter.Router) *consent.Handler {
	var ctx = c.Context()
	expectDependency(c.GetLogger(), ctx.ConsentManager)

	w := newJSONWriter(c.GetLogger())
	h := consent.NewHandler(w, ctx.ConsentManager, sessions.NewCookieStore(c.GetCookieSecret()), c.LogoutRedirectURL)
	h.SetRoutes(frontend, backend)
	return h
}
