// Package firewall defines an API for validating access requests.
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

package firewall

import (
	"context"
	"net/http"
	"time"
)

// Context contains an access token's session data
type Context struct {
	// Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app.
	// This is usually a uuid but you can choose a urn or some other id too.
	Subject string `json:"subject"`

	// GrantedScopes is a list of scopes that the subject authorized when asked for consent.
	GrantedScopes []string `json:"grantedScopes"`

	// Issuer is the id of the issuer, typically an hydra instance.
	Issuer string `json:"issuer"`

	// ClientID is id of the client the token was issued for..
	ClientID string `json:"clientId"`

	// IssuedAt is the token creation time stamp.
	IssuedAt time.Time `json:"issuedAt"`

	// ExpiresAt is the expiry timestamp.
	ExpiresAt time.Time `json:"expiresAt"`

	// Extra represents arbitrary session data.
	Extra map[string]interface{} `json:"accessTokenExtra"`
}

// AccessRequest is the warden's request object.
//
// swagger:model wardenAccessRequest
type AccessRequest struct {
	// Resource is the resource that access is requested to.
	Resource string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action string `json:"action"`

	// Subejct is the subject that is requesting access.
	Subject string `json:"subject"`

	// Context is the request's environmental context.
	Context map[string]interface{} `json:"context"`
}

// swagger:model tokenAllowedRequest
type TokenAccessRequest struct {
	// Resource is the resource that access is requested to.
	Resource string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action string `json:"action"`

	// Context is the request's environmental context.
	Context map[string]interface{} `json:"context"`
}

// Firewall offers various validation strategies for access tokens.
type Firewall interface {
	// IsAllowed uses policies to return nil if the access request can be fulfilled or an error if not.
	//
	//  ctx, err := firewall.IsAllowed(context.Background(), &AccessRequest{
	//    Subject:  "alice",
	//    Resource: "matrix",
	//    Action:   "create",
	//    Context:  ladon.Context{},
	//  }, "photos", "files")
	//
	//  fmt.Sprintf("%s", ctx.Subject)
	IsAllowed(ctx context.Context, accessRequest *AccessRequest) error

	// TokenAllowed uses policies and a token to return a context and no error if the access request can be fulfilled or an error if not.
	//
	//  ctx, err := firewall.TokenAllowed(context.Background(), "access-token", &TokenAccessRequest{
	//    Resource: "matrix",
	//    Action:   "create",
	//    Context:  ladon.Context{},
	//  }, "photos", "files")
	//
	//  fmt.Sprintf("%s", ctx.Subject)
	TokenAllowed(ctx context.Context, token string, accessRequest *TokenAccessRequest, scopes ...string) (*Context, error)

	// TokenFromRequest returns an access token from the HTTP Authorization header.
	//
	//  func anyHttpHandler(w http.ResponseWriter, r *http.Request) {
	//    ctx, err := firewall.TokenAllowed(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
	//    fmt.Sprintf("%s", ctx.Subject)
	//  }
	TokenFromRequest(r *http.Request) string
}
