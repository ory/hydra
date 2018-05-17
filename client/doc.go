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

// Package client implements OAuth 2.0 client management capabilities
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are granted
// to applications that want to use OAuth 2.0 access and refresh tokens.
//
// In ORY Hydra, OAuth 2.0 clients are used to manage ORY Hydra itself. These clients may gain highly privileged access
// if configured that way. This endpoint should be well protected and only called by code you trust.
//

package client

// swagger:parameters createOAuth2Client
type swaggerCreateClientPayload struct {
	// in: body
	// required: true
	Body Client
}

// swagger:parameters updateOAuth2Client
type swaggerUpdateClientPayload struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body Client
}

// swagger:parameters listOAuth2Clients
type swaggerListClientsParameter struct {
	// The maximum amount of policies returned.
	// in: query
	Limit int `json:"limit"`

	// The offset from where to start looking.
	// in: query
	Offset int `json:"offset"`
}

// A list of clients.
// swagger:response oAuth2ClientList
type swaggerListClientsResult struct {
	// in: body
	// type: array
	Body []Client
}

// swagger:parameters getOAuth2Client deleteOAuth2Client
type swaggerQueryClientPayload struct {
	// The id of the OAuth 2.0 Client.
	//
	// unique: true
	// in: path
	ID string `json:"id"`
}
