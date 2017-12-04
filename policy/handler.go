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

package policy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/pkg"
	"github.com/ory/ladon"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

const (
	endpoint         = "/policies"
	scope            = "hydra.policies"
	policyResource   = "policies"
	policiesResource = "policies:%s"
)

type Handler struct {
	Manager        ladon.Manager
	H              herodot.Writer
	W              firewall.Firewall
	ResourcePrefix string
}

func (h *Handler) PrefixResource(resource string) string {
	if h.ResourcePrefix == "" {
		h.ResourcePrefix = "rn:hydra"
	}

	if h.ResourcePrefix[len(h.ResourcePrefix)-1] == ':' {
		h.ResourcePrefix = h.ResourcePrefix[:len(h.ResourcePrefix)-1]
	}

	return h.ResourcePrefix + ":" + resource
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(endpoint, h.Create)
	r.GET(endpoint, h.List)
	r.GET(endpoint+"/:id", h.Get)
	r.PUT(endpoint+"/:id", h.Update)
	r.DELETE(endpoint+"/:id", h.Delete)
}

// swagger:route GET /policies policy listPolicies
//
// List Access Control Policies
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:policies"],
//    "actions": ["list"],
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
//       oauth2: hydra.policies
//
//     Responses:
//       200: policyList
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()
	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(policyResource),
		Action:   "list",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	val := r.URL.Query().Get("offset")
	if val == "" {
		val = "0"
	}

	offset, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	val = r.URL.Query().Get("limit")
	if val == "" {
		val = "500"
	}

	limit, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	policies, err := h.Manager.GetAll(limit, offset)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	h.H.Write(w, r, policies)
}

// swagger:route POST /policies policy createPolicy
//
// Create an Access Control Policy
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:policies"],
//    "actions": ["create"],
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
//       oauth2: hydra.policies
//
//     Responses:
//       201: policy
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p = ladon.DefaultPolicy{
		Conditions: ladon.Conditions{},
	}
	ctx := r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(policyResource),
		Action:   "create",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if p.ID == "" {
		p.ID = uuid.New()
	}

	if err := h.Manager.Create(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	h.H.WriteCreated(w, r, "/policies/"+p.ID, &p)
}

// swagger:route GET /policies/{id} policy getPolicy
//
// Get an Access Control Policy
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:policies:<id>"],
//    "actions": ["get"],
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
//       oauth2: hydra.policies
//
//     Responses:
//       200: policy
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(policiesResource), ps.ByName("id")),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	policy, err := h.Manager.Get(ps.ByName("id"))
	if err != nil {
		if err.Error() == "Not found" {
			h.H.WriteError(w, r, errors.WithStack(pkg.ErrNotFound))
			return
		}
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	h.H.Write(w, r, policy)
}

// swagger:route DELETE /policies/{id} policy deletePolicy
//
// Delete an Access Control Policy
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:policies:<id>"],
//    "actions": ["delete"],
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
//       oauth2: hydra.policies
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(policiesResource), id),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.Manager.Delete(id); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route PUT /policies/{id} policy updatePolicy
//
// Update an Access Control Polic
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:policies"],
//    "actions": ["update"],
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
//       oauth2: hydra.policies
//
//     Responses:
//       200: policy
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	var p = ladon.DefaultPolicy{Conditions: ladon.Conditions{}}
	var ctx = r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(policiesResource), id),
		Action:   "update",
	}, scope); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if p.ID != id {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Payload ID does not match ID from URL"))
		return
	}

	if err := h.Manager.Update(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	h.H.Write(w, r, p)
}
