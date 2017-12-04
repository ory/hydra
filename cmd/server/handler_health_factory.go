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
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/health"
)

func newHealthHandler(c *config.Config, router *httprouter.Router) *health.Handler {
	h := &health.Handler{
		Metrics:        c.GetMetrics(),
		H:              herodot.NewJSONWriter(c.GetLogger()),
		W:              c.Context().Warden,
		ResourcePrefix: c.AccessControlResourcePrefix,
	}
	h.SetRoutes(router)
	return h
}
