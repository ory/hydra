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

// swagger:model healthStatus
type swaggerHealthStatus struct {
	// Status always contains "ok"
	Status string `json:"status"`
}

// swagger:model version
type swaggerVersion struct {
	Version string `json:"version"`
}

// swagger:route GET /metrics/prometheus metrics getPrometheusMetrics
//
// Retrieve Prometheus metrics
//
// This endpoint returns metrics formatted for Prometheus.
//
//     Responses:
//       200
func prometheusDummy() {}
