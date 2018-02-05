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

package health

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/metrics"
)

const (
	HealthStatusPath = "/health/status"
)

type Handler struct {
	Metrics        *metrics.MetricsManager
	H              *herodot.JSONWriter
	W              firewall.Firewall
	ResourcePrefix string
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
}

// swagger:route GET /health/status health getInstanceStatus
//
// Check health status of this instance
//
// This endpoint returns `{ "status": "ok" }`. This status let's you know that the HTTP server is up and running. This
// status does currently not include checks whether the database connection is up and running. This endpoint does not
// require the `X-Forwarded-Proto` header when TLS termination is set.
//
//
// Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state,
// only to a single instance.
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) Health(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write([]byte(`{"status": "ok"}`))
}
