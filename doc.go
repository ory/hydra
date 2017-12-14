// Package main Hydra OAuth2 & OpenID Connect Server
//
// Please refer to the user guide for in-depth documentation: https://ory.gitbooks.io/hydra/content/
//
//
// Hydra offers OAuth 2.0 and OpenID Connect Core 1.0 capabilities as a service. Hydra is different, because it works with any existing authentication infrastructure, not just LDAP or SAML. By implementing a consent app (works with any programming language) you build a bridge between Hydra and your authentication infrastructure.
// Hydra is able to securely manage JSON Web Keys, and has a sophisticated policy-based access control you can use if you want to.
// Hydra is suitable for green- (new) and brownfield (existing) projects. If you are not familiar with OAuth 2.0 and are working on a greenfield project, we recommend evaluating if OAuth 2.0 really serves your purpose. Knowledge of OAuth 2.0 is imperative in understanding what Hydra does and how it works.
//
//
// The official repository is located at https://github.com/ory/hydra
//
//
// ### Important REST API Documentation Notes
//
// The swagger generator used to create this documentation does currently not support example responses. To see
// request and response payloads click on **"Show JSON schema"**:
// ![Enable JSON Schema on Apiary](https://storage.googleapis.com/ory.am/hydra/json-schema.png)
//
//
// The API documentation always refers to the latest tagged version of ORY Hydra. For previous API documentations, please
// refer to https://github.com/ory/hydra/blob/<tag-id>/docs/api.swagger.yaml - for example:
//
// - 0.9.13: https://github.com/ory/hydra/blob/v0.9.13/docs/api.swagger.yaml
// - 0.8.1: https://github.com/ory/hydra/blob/v0.8.1/docs/api.swagger.yaml
//
//     Schemes: http, https
//     Host:
//     BasePath: /
//     Version: Latest
//     License: Apache 2.0 https://github.com/ory/hydra/blob/master/LICENSE
//     Contact: ORY <hi@ory.am> https://www.ory.am
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
//           hydra.health: "A scope required to get health information"
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
