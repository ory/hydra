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
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
)

func newClientManager(c *config.Config) client.Manager {
	ctx := c.Context()
	return ctx.Connection.NewClientManager(ctx.Hasher)
}

func newClientHandler(c *config.Config, router *httprouter.Router, manager client.Manager) *client.Handler {
	w := herodot.NewJSONWriter(c.GetLogger())
	w.ErrorEnhancer = writerErrorEnhancer

	expectDependency(c.GetLogger(), manager)
	h := client.NewHandler(
		manager,
		w,
		strings.Split(c.DefaultClientScope, ","),
		c.GetSubjectTypesSupported(),
	)
	h.SetRoutes(router)
	return h
}
