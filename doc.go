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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

// Package main ORY Hydra - Cloud Native OAuth 2.0 and OpenID Connect Server
//
// Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here. Keep in mind that this document reflects the latest branch, always. Support for versioned documentation is coming in the future.
//
//     Schemes: http, https
//     Host:
//     BasePath: /
//     Version: Latest
//     License: Apache 2.0 https://github.com/ory/hydra/blob/master/LICENSE
//     Contact: ORY <hi@ory.am> https://www.ory.sh
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
//         authorizationUrl: https://your-hydra-instance.com/oauth2/auth
//         tokenUrl: https://your-hydra-instance.com/oauth2/token
//         flow: accessCode
//         scopes:
//           hydra.clients: "A scope required to manage OAuth 2.0 Clients"
//           hydra.policies: "A scope required to manage access control policies"
//           hydra.warden: "A scope required to make access control inquiries"
//           hydra.warden.groups: "A scope required to manage warden groups"
//           hydra.keys.get: "A scope required to fetch JSON Web Keys"
//           hydra.keys.create: "A scope required to create JSON Web Keys"
//           hydra.keys.delete: "A scope required to delete JSON Web Keys"
//           hydra.keys.update: "A scope required to get JSON Web Keys"
//           hydra.consent: "A scope required to fetch and modify consent requests"
//           offline: "A scope required when requesting refresh tokens"
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
package main
