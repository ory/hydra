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

package healthx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ory/herodot"
)

const (
	AliveCheckPath = "/health/alive"
	ReadyCheckPath = "/health/ready"
	VersionPath    = "/version"
)

func RoutesToObserve() []string {
	return []string{
		AliveCheckPath,
		ReadyCheckPath,
		VersionPath,
	}
}

type ReadyChecker func() error
type ReadyCheckers map[string]ReadyChecker

func NoopReadyChecker() error {
	return nil
}

type Handler struct {
	H             herodot.Writer
	VersionString string
	ReadyChecks   ReadyCheckers
}

func NewHandler(
	h herodot.Writer,
	version string,
	readyChecks ReadyCheckers,
) *Handler {
	return &Handler{
		H:             h,
		VersionString: version,
		ReadyChecks:   readyChecks,
	}
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET(AliveCheckPath, h.Alive)
	r.GET(ReadyCheckPath, h.Ready)
	r.GET(VersionPath, h.Version)
}

// swagger:route GET /health/alive health isInstanceAlive
//
// Check alive status
//
// This endpoint returns a 200 status code when the HTTP server is up running.
// This status does currently not include checks whether the database connection is working.
//
// If the service supports TLS Edge Termination, this endpoint does not require the
// `X-Forwarded-Proto` header to be set.
//
// Be aware that if you are running multiple nodes of this service, the health status will never
// refer to the cluster state, only to a single instance.
//
//     Produces:
//     - application/json
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
// Check readiness status
//
// This endpoint returns a 200 status code when the HTTP server is up running and the environment dependencies (e.g.
// the database) are responsive as well.
//
// If the service supports TLS Edge Termination, this endpoint does not require the
// `X-Forwarded-Proto` header to be set.
//
// Be aware that if you are running multiple nodes of this service, the health status will never
// refer to the cluster state, only to a single instance.
//
//     Produces:
//     - application/json
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
// Get service version
//
// This endpoint returns the service version typically notated using semantic versioning.
//
// If the service supports TLS Edge Termination, this endpoint does not require the
// `X-Forwarded-Proto` header to be set.
//
// Be aware that if you are running multiple nodes of this service, the health status will never
// refer to the cluster state, only to a single instance.
//
//     Produces:
//     - application/json
//
//	   Responses:
// 			200: version
func (h *Handler) Version(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.H.Write(rw, r, &swaggerVersion{
		Version: h.VersionString,
	})
}
