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

package hydra

import (
	"github.com/ory/hydra/sdk/go/hydra/swagger"
)

// SDK helps developers interact with ORY Hydra using a Go API.
type SDK interface {
	JWKApi
	OAuth2API
}

type JWKApi interface {
	CreateJsonWebKeySet(set string, body swagger.JsonWebKeySetGeneratorRequest) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
	DeleteJsonWebKey(kid string, set string) (*swagger.APIResponse, error)
	DeleteJsonWebKeySet(set string) (*swagger.APIResponse, error)
	GetJsonWebKey(kid string, set string) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
	GetJsonWebKeySet(set string) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
	UpdateJsonWebKey(kid string, set string, body swagger.JsonWebKey) (*swagger.JsonWebKey, *swagger.APIResponse, error)
	UpdateJsonWebKeySet(set string, body swagger.JsonWebKeySet) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
}

type OAuth2API interface {
	AcceptConsentRequest(challenge string, body swagger.AcceptConsentRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)
	AcceptLoginRequest(challenge string, body swagger.AcceptLoginRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)
	RejectConsentRequest(challenge string, body swagger.RejectRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)
	RejectLoginRequest(challenge string, body swagger.RejectRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)
	GetLoginRequest(challenge string) (*swagger.LoginRequest, *swagger.APIResponse, error)
	GetConsentRequest(challenge string) (*swagger.ConsentRequest, *swagger.APIResponse, error)

	CreateOAuth2Client(body swagger.OAuth2Client) (*swagger.OAuth2Client, *swagger.APIResponse, error)
	DeleteOAuth2Client(id string) (*swagger.APIResponse, error)
	GetOAuth2Client(id string) (*swagger.OAuth2Client, *swagger.APIResponse, error)
	DiscoverOpenIDConfiguration() (*swagger.WellKnown, *swagger.APIResponse, error)
	IntrospectOAuth2Token(token string, scope string) (*swagger.OAuth2TokenIntrospection, *swagger.APIResponse, error)
	ListOAuth2Clients(limit int64, offset int64) ([]swagger.OAuth2Client, *swagger.APIResponse, error)
	RevokeOAuth2Token(token string) (*swagger.APIResponse, error)
	RevokeAllUserConsentSessions(user string) (*swagger.APIResponse, error)
	UpdateOAuth2Client(id string, body swagger.OAuth2Client) (*swagger.OAuth2Client, *swagger.APIResponse, error)
	RevokeAuthenticationSession(user string) (*swagger.APIResponse, error)
	RevokeUserClientConsentSessions(user string, client string) (*swagger.APIResponse, error)

	ListUserConsentSessions(user string) ([]swagger.PreviousConsentSession, *swagger.APIResponse, error)
	FlushInactiveOAuth2Tokens(body swagger.FlushInactiveOAuth2TokensRequest) (*swagger.APIResponse, error)
}
