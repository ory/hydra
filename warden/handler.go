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

package warden

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/firewall"
	"github.com/pkg/errors"
)

const (
	// TokenAllowedHandlerPath points to the token access request validation endpoint.
	TokenAllowedHandlerPath = "/warden/token/allowed"

	// AllowedHandlerPath points to the access request validation endpoint.
	AllowedHandlerPath = "/warden/allowed"
)

type wardenAuthorizedRequest struct {
	// Scopes is an array of scopes that are requried.
	Scopes []string `json:"scopes"`

	// Token is the token to introspect.
	Token string `json:"token"`
}

type wardenAccessRequest struct {
	*firewall.TokenAccessRequest
	*wardenAuthorizedRequest
}

var notAllowed = struct {
	Allowed bool `json:"allowed"`
}{Allowed: false}

// WardenHandler is capable of handling HTTP request and validating access tokens and access requests.
type WardenHandler struct {
	H      herodot.Writer
	Warden firewall.Firewall

	ResourcePrefix string
}

func (h *WardenHandler) PrefixResource(resource string) string {
	if h.ResourcePrefix == "" {
		h.ResourcePrefix = "rn:hydra"
	}

	if h.ResourcePrefix[len(h.ResourcePrefix)-1] == ':' {
		h.ResourcePrefix = h.ResourcePrefix[:len(h.ResourcePrefix)-1]
	}

	return h.ResourcePrefix + ":" + resource
}

func NewHandler(c *config.Config, router *httprouter.Router) *WardenHandler {
	ctx := c.Context()

	h := &WardenHandler{
		H:              herodot.NewJSONWriter(c.GetLogger()),
		Warden:         ctx.Warden,
		ResourcePrefix: c.AccessControlResourcePrefix,
	}
	h.SetRoutes(router)

	return h
}

func (h *WardenHandler) SetRoutes(r *httprouter.Router) {
	r.POST(TokenAllowedHandlerPath, h.TokenAllowed)
	r.POST(AllowedHandlerPath, h.Allowed)
}

// swagger:route POST /warden/allowed warden doesWardenAllowAccessRequest
//
// Check if an access request is valid (without providing an access token)
//
// Checks if a subject (typically a user or a service) is allowed to perform an action on a resource. This endpoint requires a subject,
// a resource name, an action name and a context. If the subject is not allowed to perform the action on the resource,
// this endpoint returns a 200 response with `{ "allowed": false}`, otherwise `{ "allowed": true }` is returned.
//
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:allowed"],
//    "actions": ["decide"],
//    "effect": "allow"
//  }
//  ```
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.warden
//
//     Responses:
//       200: wardenAccessRequestResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *WardenHandler) Allowed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()
	if _, err := h.Warden.TokenAllowed(ctx, h.Warden.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource("warden:allowed"),
		Action:   "decide",
	}, "hydra.warden"); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	var access = new(firewall.AccessRequest)
	if err := json.NewDecoder(r.Body).Decode(access); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	defer r.Body.Close()

	if err := h.Warden.IsAllowed(ctx, access); err != nil {
		h.H.Write(w, r, &notAllowed)
		return
	}

	res := notAllowed
	res.Allowed = true
	h.H.Write(w, r, &res)
}

// swagger:route POST /warden/token/allowed warden doesWardenAllowTokenAccessRequest
//
// Check if an access request is valid (providing an access token)
//
// Checks if a token is valid and if the token subject is allowed to perform an action on a resource.
// This endpoint requires a token, a scope, a resource name, an action name and a context.
//
//
// If a token is expired/invalid, has not been granted the requested scope or the subject is not allowed to
// perform the action on the resource, this endpoint returns a 200 response with `{ "allowed": false}`.
//
//
// Extra data set through the `accessTokenExtra` field in the consent flow will be included in the response.
//
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:token:allowed"],
//    "actions": ["decide"],
//    "effect": "allow"
//  }
//  ```
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.warden
//
//     Responses:
//       200: wardenTokenAccessRequestResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *WardenHandler) TokenAllowed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	_, err := h.Warden.TokenAllowed(ctx, h.Warden.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource("warden:token:allowed"),
		Action:   "decide",
	}, "hydra.warden")
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	var ar = wardenAccessRequest{
		TokenAccessRequest:      new(firewall.TokenAccessRequest),
		wardenAuthorizedRequest: new(wardenAuthorizedRequest),
	}
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	defer r.Body.Close()

	authContext, err := h.Warden.TokenAllowed(ctx, ar.Token, ar.TokenAccessRequest, ar.Scopes...)
	if err != nil {
		h.H.Write(w, r, &notAllowed)
		return
	}

	h.H.Write(w, r, struct {
		*firewall.Context
		Allowed bool `json:"allowed"`
	}{
		Context: authContext,
		Allowed: true,
	})
}

func TokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		return ""
	}

	return split[1]
}
