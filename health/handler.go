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

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	AliveCheckPath        = "/health/alive"
	ReadyCheckPath        = "/health/ready"
	VersionPath           = "/version"
	MetricsPrometheusPath = "/metrics/prometheus"
)

type ReadyChecker func() error

type Handler struct {
	H             *herodot.JSONWriter
	VersionString string
	ReadyChecks   map[string]ReadyChecker
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET(AliveCheckPath, h.Alive)
	r.GET(ReadyCheckPath, h.Ready)
	r.GET(VersionPath, h.Version)

	// using r.Handler because promhttp.Handler() returns http.Handler
	r.Handler("GET", MetricsPrometheusPath, promhttp.Handler())

	// BC compatible health check
	r.GET("/health/status", pkg.PermanentRedirect(AliveCheckPath))

}

// swagger:route GET /health/alive health isInstanceAlive
//
// Check the Alive Status
//
// This endpoint returns a 200 status code when the HTTP server is up running.
// This status does currently not include checks whether the database connection is working.
// This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.
//
// Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) Alive(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.H.Write(rw, r, &swaggerHealthStatus{
		Status: "ok",
	})
}

// swagger:route GET /health/ready health isInstanceReady
//
// Check the Readiness Status
//
// This endpoint returns a 200 status code when the HTTP server is up running and the environment dependencies (e.g.
// the database) are responsive as well.
//
// This status does currently not include checks whether the database connection is working.
// This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.
//
// Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.
//
//     Responses:
//       200: healthStatus
//       503: healthNotReadyStatus
func (h *Handler) Ready(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var notReady = swaggerNotReadyStatus{
		Errors: map[string]string{},
	}

	for n, c := range h.ReadyChecks {
		if err := c(); err != nil {
			notReady.Errors[n] = err.Error()
		}
	}

	if len(notReady.Errors) > 0 {
		h.H.WriteCode(rw, r, http.StatusServiceUnavailable, notReady)
		return
	}

	h.H.Write(rw, r, &swaggerHealthStatus{
		Status: "ok",
	})
}

// swagger:route GET /version version getVersion
//
// Get the version of Hydra
//
// This endpoint returns the version as `{ "version": "VERSION" }`. The version is only correct with the prebuilt binary and not custom builds.
//
//		Responses:
// 		200: version
func (h *Handler) Version(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.H.Write(rw, r, &swaggerVersion{
		Version: h.VersionString,
	})
}
