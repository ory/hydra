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

// ORY Hydra
//
// Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
//
//     Schemes: http, https
//     Host:
//     BasePath: /
//     Version: latest
//
//     Consumes:
//     - application/json
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     SecurityDefinitions:
//     oauth2:
//         type: oauth2
//         authorizationUrl: https://hydra.demo.ory.sh/oauth2/auth
//         tokenUrl: https://hydra.demo.ory.sh/oauth2/token
//         flow: accessCode
//         scopes:
//           offline: "A scope required when requesting refresh tokens (alias for `offline_access`)"
//           offline_access: "A scope required when requesting refresh tokens"
//           openid: "Request an OpenID Connect ID Token"
//     basic:
//         type: basic
//
//     Extensions:
//     ---
//     x-request-id: string
//     x-forwarded-proto: string
//     ---
//
// swagger:meta
package x
