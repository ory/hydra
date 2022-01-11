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

// swagger:parameters createOAuth2Client dynamicClientRegistrationCreateOAuth2Client
type dynamicClientRegistrationCreateOAuth2Client struct {
	// in: body
	// required: true
	Body Client
}

// swagger:parameters updateOAuth2Client
type swaggerUpdateClientPayload struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body Client
}

// swagger:parameters dynamicClientRegistrationUpdateOAuth2Client
type dynamicClientRegistrationUpdateOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body Client
}

// swagger:parameters patchOAuth2Client
type swaggerPatchClientPayload struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body patchRequest
}

// A JSONPatch request
//
// swagger:model patchRequest
type patchRequest []patchDocument

// A JSONPatch document as defined by RFC 6902
//
// swagger:model patchDocument
type patchDocument struct {
	// The operation to be performed
	//
	// required: true
	// example: "replace"
	Op string `json:"op"`

	// A JSON-pointer
	//
	// required: true
	// example: "/name"
	Path string `json:"path"`

	// The value to be used within the operations
	Value interface{} `json:"value"`

	// A JSON-pointer
	From string `json:"from"`
}

// A list of clients.
// swagger:response oAuth2ClientList
type swaggerListClientsResult struct {
	// in: body
	// type: array
	Body []Client
}

// swagger:parameters getOAuth2Client
type swaggerGetOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters dynamicClientRegistrationGetOAuth2Client
type dynamicClientRegistrationGetOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters deleteOAuth2Client
type swaggerDeleteOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters dynamicClientRegistrationDeleteOAuth2Client
type dynamicClientRegistrationDeleteOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}
