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

package health

import (
	"net/http"

	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/metrics"
)

const (
	HealthStatusPath  = "/health/status"
	HealthVersionPath = "/health/version"
)

type Handler struct {
	Metrics        *metrics.MetricsManager
	H              *herodot.JSONWriter
	ResourcePrefix string
	VersionString  string
}

func (h *Handler) PrefixResource(resource string) string {
	if h.ResourcePrefix == "" {
		h.ResourcePrefix = "rn:hydra"
	}

	if h.ResourcePrefix[len(h.ResourcePrefix)-1] == ':' {
		h.ResourcePrefix = h.ResourcePrefix[:len(h.ResourcePrefix)-1]
	}

	return h.ResourcePrefix + ":" + resource
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET(HealthStatusPath, h.Health)
	r.GET(HealthVersionPath, h.Version)
}

// swagger:route GET /health/status health getInstanceStatus
//
// Check the Health Status
//
// This endpoint returns a 200 status code when the HTTP server is up running. `{ "status": "ok" }`. This status does currently not include checks whether the database connection is working. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.
//
// Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) Health(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write([]byte(`{"status": "ok"}`))
}

// swagger:route GET /health/version ??
//
// Get the version of Hydra
//
// This endpoint returns the version as `{ "version": "VERSION" }`. The version might only be correct with the prebuilt binary and not custom builds.
//
//		Responses:
// 		200: version
//		500: genericError
func (h *Handler) Version(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write([]byte(fmt.Sprintf(`{"version": "%s"}`, h.VersionString)))
}
