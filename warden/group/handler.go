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

package group

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/pkg/errors"
)

// swagger:model groupMembers
type membersRequest struct {
	Members []string `json:"members"`
}

type Handler struct {
	Manager Manager
	H       herodot.Writer
	W       firewall.Firewall

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

const (
	GroupsHandlerPath = "/warden/groups"
)

const (
	GroupsResource = "warden:groups"
	GroupResource  = "warden:groups:%s"
	Scope          = "hydra.warden.groups"
)

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(GroupsHandlerPath, h.CreateGroup)
	r.GET(GroupsHandlerPath, h.FindGroupNames)
	r.GET(GroupsHandlerPath+"/:id", h.GetGroup)
	r.DELETE(GroupsHandlerPath+"/:id", h.DeleteGroup)
	r.POST(GroupsHandlerPath+"/:id/members", h.AddGroupMembers)
	r.DELETE(GroupsHandlerPath+"/:id/members", h.RemoveGroupMembers)
}

// swagger:route GET /warden/groups warden findGroupsByMember
//
// Find groups by member
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:groups"],
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
//       oauth2: hydra.warden.groups
//
//     Responses:
//       200: findGroupsByMemberResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) FindGroupNames(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()
	var member = r.URL.Query().Get("member")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(GroupsResource),
		Action:   "list",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	groups, err := h.Manager.FindGroupsByMember(member)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, groups)
}

// swagger:route POST /warden/groups warden createGroup
//
// Create a group
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:groups"],
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
//       oauth2: hydra.warden.groups
//
//     Responses:
//       201: groupResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var g Group
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(GroupsResource),
		Action:   "create",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.CreateGroup(&g); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.WriteCreated(w, r, GroupsHandlerPath+"/"+g.ID, &g)
}

// swagger:route GET /warden/groups/{id} warden getGroup
//
// Get a group by id
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:groups:<id>"],
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
//       oauth2: hydra.warden.groups
//
//     Responses:
//       201: groupResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	g, err := h.Manager.GetGroup(id)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(GroupResource), id),
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, g)
}

// swagger:route DELETE /warden/groups/{id} warden deleteGroup
//
// Delete a group by id
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:groups:<id>"],
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
//       oauth2: hydra.warden.groups
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(GroupResource), id),
		Action:   "delete",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.DeleteGroup(id); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route POST /warden/groups/{id}/members warden addMembersToGroup
//
// Add members to a group
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:groups:<id>"],
//    "actions": ["members.add"],
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
//       oauth2: hydra.warden.groups
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) AddGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	var m membersRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(GroupResource), id),
		Action:   "members.add",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.AddGroupMembers(id, m.Members); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /warden/groups/{id}/members warden removeMembersFromGroup
//
// Remove members from a group
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:warden:groups:<id>"],
//    "actions": ["members.remove"],
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
//       oauth2: hydra.warden.groups
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) RemoveGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	var m membersRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(GroupResource), id),
		Action:   "members.remove",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.RemoveGroupMembers(id, m.Members); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
