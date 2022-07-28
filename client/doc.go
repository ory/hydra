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

import (
	"github.com/ory/hydra/x"
)

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

// swagger:parameters UpdateOAuth2ClientLifespans
type swaggerUpdateOAuth2ClientLifespans struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	Body UpdateOAuth2ClientLifespans
}

// UpdateOAuth2ClientLifespans holds default lifespan configuration for the different
// token types that may be issued for the client. This configuration takes
// precedence over fosite's instance-wide default lifespan, but it may be
// overridden by a session's expires_at claim.
//
// The OIDC Hybrid grant type inherits token lifespan configuration from the implicit grant.
//
// swagger:model UpdateOAuth2ClientLifespans
type UpdateOAuth2ClientLifespans struct {
	// AuthorizationCodeGrantAccessTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	AuthorizationCodeGrantAccessTokenLifespan x.NullDuration `json:"authorization_code_grant_access_token_lifespan"`
	// AuthorizationCodeGrantIDTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	AuthorizationCodeGrantIDTokenLifespan x.NullDuration `json:"authorization_code_grant_id_token_lifespan"`
	// AuthorizationCodeGrantRefreshTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	AuthorizationCodeGrantRefreshTokenLifespan x.NullDuration `json:"authorization_code_grant_refresh_token_lifespan"`
	// ClientCredentialsGrantAccessTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	ClientCredentialsGrantAccessTokenLifespan x.NullDuration `json:"client_credentials_grant_access_token_lifespan"`
	// ImplicitGrantAccessTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	ImplicitGrantAccessTokenLifespan x.NullDuration `json:"implicit_grant_access_token_lifespan"`
	// ImplicitGrantIDTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	ImplicitGrantIDTokenLifespan x.NullDuration `json:"implicit_grant_id_token_lifespan"`
	// JwtBearerGrantAccessTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	JwtBearerGrantAccessTokenLifespan x.NullDuration `json:"jwt_bearer_grant_access_token_lifespan"`
	// PasswordGrantAccessTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	PasswordGrantAccessTokenLifespan x.NullDuration `json:"password_grant_access_token_lifespan"`
	// PasswordGrantRefreshTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	PasswordGrantRefreshTokenLifespan x.NullDuration `json:"password_grant_refresh_token_lifespan"`
	// RefreshTokenGrantIDTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	RefreshTokenGrantIDTokenLifespan x.NullDuration `json:"refresh_token_grant_id_token_lifespan"`
	// RefreshTokenGrantAccessTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	RefreshTokenGrantAccessTokenLifespan x.NullDuration `json:"refresh_token_grant_access_token_lifespan"`
	// RefreshTokenGrantRefreshTokenLifespan configures this client's lifespan and takes precedence over instance-wide configuration
	RefreshTokenGrantRefreshTokenLifespan x.NullDuration `json:"refresh_token_grant_refresh_token_lifespan"`
}

// swagger:parameters dynamicClientRegistrationDeleteOAuth2Client
type dynamicClientRegistrationDeleteOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}
