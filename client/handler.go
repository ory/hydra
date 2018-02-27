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

package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/rand/sequence"
	"github.com/ory/ladon"
	"github.com/ory/pagination"
	"github.com/pkg/errors"
)

type Handler struct {
	Manager        Manager
	H              herodot.Writer
	W              firewall.Firewall
	ResourcePrefix string
}

const (
	ClientsHandlerPath = "/clients"
)

const (
	ClientsResource = "clients"
	ClientResource  = "clients:%s"
	Scope           = "hydra.clients"
)

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
	r.GET(ClientsHandlerPath, h.List)
	r.POST(ClientsHandlerPath, h.Create)
	r.GET(ClientsHandlerPath+"/:id", h.Get)
	r.PUT(ClientsHandlerPath+"/:id", h.Update)
	r.DELETE(ClientsHandlerPath+"/:id", h.Delete)
}

// swagger:route POST /clients oAuth2 createOAuth2Client
//
// Create an OAuth 2.0 client
//
// Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients"],
//    "actions": ["create"],
//    "effect": "allow"
//  }
//  ```
//
//  Additionally, the context key "owner" is set to the owner of the client, allowing policies such as:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients"],
//    "actions": ["create"],
//    "effect": "allow",
//    "conditions": { "owner": { "type": "EqualsSubjectCondition" } }
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
//       oauth2: hydra.clients
//
//     Responses:
//       200: oAuth2Client
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var c Client
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(ClientsResource),
		Action:   "create",
		Context: map[string]interface{}{
			"owner": c.Owner,
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if len(c.Secret) == 0 {
		secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-.~"))
		if err != nil {
			h.H.WriteError(w, r, errors.WithStack(err))
			return
		}
		c.Secret = string(secret)
	} else if len(c.Secret) < 6 {
		h.H.WriteError(w, r, errors.New("The client secret must be at least 6 characters long"))
		return
	}

	secret := c.Secret
	if err := h.Manager.CreateClient(&c); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	c.Secret = ""

	if !c.Public {
		c.Secret = secret
	}

	h.H.WriteCreated(w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

// swagger:route PUT /clients/{id} oAuth2 updateOAuth2Client
//
// Update an OAuth 2.0 Client
//
// Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients"],
//    "actions": ["update"],
//    "effect": "allow"
//  }
//  ```
//
//  Additionally, the context key "owner" is set to the owner of the client, allowing policies such as:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients"],
//    "actions": ["update"],
//    "effect": "allow",
//    "conditions": { "owner": { "type": "EqualsSubjectCondition" } }
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
//       oauth2: hydra.clients
//
//     Responses:
//       200: oAuth2Client
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c Client
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	o, err := h.Manager.GetConcreteClient(ps.ByName("id"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(ClientsResource),
		Action:   "update",
		Context: ladon.Context{
			"owner": o.Owner,
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	var secret string
	if len(c.Secret) > 0 && len(c.Secret) < 6 {
		h.H.WriteError(w, r, errors.New("The client secret must be at least 6 characters long"))
		return
	} else {
		secret = c.Secret
	}

	c.ID = ps.ByName("id")
	if err := h.Manager.UpdateClient(&c); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	c.Secret = secret
	h.H.WriteCreated(w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

// swagger:route GET /clients oAuth2 listOAuth2Clients
//
// List OAuth 2.0 Clients
//
// This endpoint lists all clients in the database, and never returns client secrets.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
// The subject making the request needs to be assigned to a policy containing:
//
// ```
// {
//   "resources": ["rn:hydra:clients"],
//   "actions": ["get"],
//   "effect": "allow"
// }
// ```
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
//       oauth2: hydra.clients
//
//     Responses:
//       200: oAuth2ClientList
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: h.PrefixResource(ClientsResource),
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	limit, offset := pagination.Parse(r, 100, 0, 500)
	c, err := h.Manager.GetClients(limit, offset)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	clients := make([]Client, len(c))
	k := 0
	for _, cc := range c {
		clients[k] = cc
		clients[k].Secret = ""
		k++
	}

	h.H.Write(w, r, clients)
}

// swagger:route GET /clients/{id} oAuth2 getOAuth2Client
//
// Get an OAuth 2.0 Client.
//
// Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients:<some-id>"],
//    "actions": ["get"],
//    "effect": "allow"
//  }
//  ```
//
//  Additionally, the context key "owner" is set to the owner of the client, allowing policies such as:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients:<some-id>"],
//    "actions": ["get"],
//    "effect": "allow",
//    "conditions": { "owner": { "type": "EqualsSubjectCondition" } }
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
//       oauth2: hydra.clients
//
//     Responses:
//       200: oAuth2Client
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	c, err := h.Manager.GetConcreteClient(id)
	if err != nil {
		if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
			Resource: fmt.Sprintf(h.PrefixResource(ClientResource), id),
			Action:   "get",
		}, Scope); err != nil {
			h.H.WriteError(w, r, err)
			return
		}

		h.H.WriteError(w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(ClientResource), id),
		Action:   "get",
		Context: ladon.Context{
			"owner": c.GetOwner(),
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	c.Secret = ""
	h.H.Write(w, r, c)
}

// swagger:route DELETE /clients/{id} oAuth2 deleteOAuth2Client
//
// Deletes an OAuth 2.0 Client
//
// Delete an existing OAuth 2.0 Client by its ID.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients:<some-id>"],
//    "actions": ["delete"],
//    "effect": "allow"
//  }
//  ```
//
//  Additionally, the context key "owner" is set to the owner of the client, allowing policies such as:
//
//  ```
//  {
//    "resources": ["rn:hydra:clients:<some-id>"],
//    "actions": ["delete"],
//    "effect": "allow",
//    "conditions": { "owner": { "type": "EqualsSubjectCondition" } }
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
//       oauth2: hydra.clients
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	c, err := h.Manager.GetConcreteClient(id)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(h.PrefixResource(ClientResource), id),
		Action:   "delete",
		Context: ladon.Context{
			"owner": c.GetOwner(),
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.DeleteClient(id); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
