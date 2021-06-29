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

package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/ory/herodot"
	"github.com/ory/x/sqlcon"

	"github.com/ory/hydra/x"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/ory/x/pagination"
)

type Handler struct {
	r InternalRegistry
}

const (
	ClientsHandlerPath    = "/clients"
	DynClientsHandlerPath = "/connect/register"
)

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{
		r: r,
	}
}

func (h *Handler) SetRoutes(admin *x.RouterAdmin, public *x.RouterPublic, dynamicRegistration bool) {
	admin.GET(ClientsHandlerPath, h.List)
	admin.POST(ClientsHandlerPath, h.Create)
	admin.GET(ClientsHandlerPath+"/:id", h.Get)
	admin.PUT(ClientsHandlerPath+"/:id", h.Update)
	admin.PATCH(ClientsHandlerPath+"/:id", h.Patch)
	admin.DELETE(ClientsHandlerPath+"/:id", h.Delete)

	if dynamicRegistration {
		public.POST(DynClientsHandlerPath, h.CreateDynamicRegistration)
		public.GET(DynClientsHandlerPath, h.GetDynamicRegistration)
		public.PUT(DynClientsHandlerPath, h.UpdateDynamicRegistration)
		public.DELETE(DynClientsHandlerPath, h.DeleteDynamicRegistration)
	}
}

// swagger:route POST /clients admin createOAuth2Client
//
// Create an OAuth 2.0 Client
//
// Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       201: oAuth2Client
//       400: genericError
//       409: genericError
//       500: genericError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := h.create(w, r, h.r.ClientValidator().Validate)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	h.r.Writer().WriteCreated(w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

// swagger:route POST /connect/register public createOAuth2ClientPublic
//
// Create an OAuth 2.0 Client
//
// Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somewhere safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       201: oAuth2Client
//       400: genericError
//       409: genericError
//       500: genericError
func (h *Handler) CreateDynamicRegistration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := h.create(w, r, h.r.ClientValidator().ValidateDynamicRegistration)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	h.r.Writer().WriteCreated(w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request, f func(*Client) error) (*Client, error) {
	var c Client

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return nil, err
	}

	if len(c.Secret) == 0 {
		secretb, err := x.GenerateSecret(26)
		if err != nil {
			return nil, err
		}
		c.Secret = string(secretb)
	}

	if err := f(&c); err != nil {
		return nil, err
	}

	secret := c.Secret
	c.CreatedAt = time.Now().UTC().Round(time.Second)
	c.UpdatedAt = c.CreatedAt
	if err := h.r.ClientManager().CreateClient(r.Context(), &c); err != nil {
		return nil, err
	}

	c.Secret = ""
	if !c.IsPublic() {
		c.Secret = secret
	}
	return &c, nil
}

// swagger:route PUT /clients/{id} admin updateOAuth2Client
//
// Update an OAuth 2.0 Client
//
// Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2Client
//       500: genericError
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c Client

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	c.OutfacingID = ps.ByName("id")
	if err := h.updateClient(r.Context(), &c); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &c)
}

// swagger:route PATCH /clients/{id} admin patchOAuth2Client
//
// Patch an OAuth 2.0 Client
//
// Patch an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2Client
//       500: genericError
func (h *Handler) Patch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	patchJSON, err := io.ReadAll(r.Body)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	id := ps.ByName("id")
	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), id)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if err := x.ApplyJSONPatch(patchJSON, c, "/id"); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if err := h.updateClient(r.Context(), c); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, c)
}

func (h *Handler) updateClient(ctx context.Context, c *Client) error {
	var secret string
	if len(c.Secret) > 0 {
		secret = c.Secret
	}
	if err := h.r.ClientValidator().Validate(c); err != nil {
		return err
	}

	c.UpdatedAt = time.Now().UTC().Round(time.Second)
	if err := h.r.ClientManager().UpdateClient(ctx, c); err != nil {
		return err
	}
	c.Secret = secret
	return nil
}

// swagger:route PUT/connect/register public updateOAuth2ClientPublic
//
// Update an OAuth 2.0 Client
//
// Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2Client
//       500: genericError
func (h *Handler) UpdateDynamicRegistration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c Client

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	var secret string
	if len(c.Secret) > 0 {
		secret = c.Secret
	}

	c.OutfacingID = r.URL.Query().Get("client_id")
	if err := h.validateDynClientRegistrationAuthorization(r, c); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.r.ClientValidator().ValidateDynamicRegistration(&c); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	c.UpdatedAt = time.Now().UTC().Round(time.Second)
	if err := h.r.ClientManager().UpdateClient(r.Context(), &c); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	c.Secret = secret
	h.r.Writer().Write(w, r, &c)
}

// swagger:route GET /clients admin listOAuth2Clients
//
// List OAuth 2.0 Clients
//
// This endpoint lists all clients in the database, and never returns client secrets. As a default it lists the first 100 clients. The `limit` parameter can be used to retrieve more clients, but it has an upper bound at 500 objects. Pagination should be used to retrieve more than 500 objects.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
// The "Link" header is also included in successful responses, which contains one or more links for pagination, formatted like so: '<https://hydra-url/admin/clients?limit={limit}&offset={offset}>; rel="{page}"', where page is one of the following applicable pages: 'first', 'next', 'last', and 'previous'.
// Multiple links can be included in this header, and will be separated by a comma.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2ClientList
//       500: genericError
func (h *Handler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, offset := pagination.Parse(r, 100, 0, 500)

	c, err := h.r.ClientManager().GetClients(r.Context(), limit, offset)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	for k := range c {
		c[k].Secret = ""
	}

	n, err := h.r.ClientManager().CountClients(r.Context())
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	pagination.Header(w, r.URL, n, limit, offset)

	if c == nil {
		c = []Client{}
	}

	h.r.Writer().Write(w, r, c)
}

// swagger:route GET /clients/{id} admin getOAuth2Client
//
// Get an OAuth 2.0 Client.
//
// Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2Client
//       401: genericError
//       500: genericError
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")

	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), id)
	if err != nil {
		if errors.Is(err, sqlcon.ErrNoRows) {
			err = herodot.ErrUnauthorized.WithReason("The requested OAuth 2.0 client does not exist or you did not provide the necessary credentials")
		}
		h.r.Writer().WriteError(w, r, err)
		return
	}

	c.Secret = ""
	h.r.Writer().Write(w, r, c)
}

// swagger:route GET /connect/register public getOAuth2ClientPublic
//
// Get an OAuth 2.0 Client.
//
// Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2Client
//       401: genericError
//       500: genericError
func (h *Handler) GetDynamicRegistration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = r.URL.Query().Get("client_id")

	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), id)
	if err != nil {
		err = herodot.ErrUnauthorized.WithReason("The requested OAuth 2.0 client does not exist or you did not provide the necessary credentials")
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if err := h.validateDynClientRegistrationAuthorization(r, *c); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnauthorized, err)
		return
	}

	c.Secret = ""
	c.Metadata = nil
	h.r.Writer().Write(w, r, c)
}

// swagger:route DELETE /clients/{id} admin deleteOAuth2Client
//
// Deletes an OAuth 2.0 Client
//
// Delete an existing OAuth 2.0 Client by its ID.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       404: genericError
//       500: genericError
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")

	if err := h.r.ClientManager().DeleteClient(r.Context(), id); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /connect/register public deleteOAuth2ClientPublic
//
// Deletes an OAuth 2.0 Client
//
// Delete an existing OAuth 2.0 Client by its ID.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       404: genericError
//       500: genericError
func (h *Handler) DeleteDynamicRegistration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = r.URL.Query().Get("client_id")

	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), id)
	if err != nil {
		err = herodot.ErrUnauthorized.WithReason("The requested OAuth 2.0 client does not exist or you did not provide the necessary credentials")
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if err := h.validateDynClientRegistrationAuthorization(r, *c); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.r.ClientManager().DeleteClient(r.Context(), id); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) validateDynClientRegistrationAuthorization(r *http.Request, c Client) error {

	basicAuth := getBasicAuth(r)
	if basicAuth == "" {
		return herodot.ErrUnauthorized.WithReason("Invalid authorization")
	}
	sDec, err := base64.StdEncoding.DecodeString(basicAuth)
	if err != nil {
		return herodot.ErrUnauthorized.WithReason("Invalid authorization")
	}
	split := strings.SplitN(string(sDec), ":", 2)
	if len(split) != 2 {
		return herodot.ErrUnauthorized.WithReason("Invalid authorization")
	}
	if c.OutfacingID != split[0] {
		return herodot.ErrUnauthorized.WithReason("Invalid authorization")
	}
	_, err = h.r.ClientManager().Authenticate(r.Context(), split[0], []byte(split[1]))
	if err != nil {
		return herodot.ErrUnauthorized.WithReason("Invalid authorization")
	}
	return nil
}

func getBasicAuth(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "Basic") {
		// Nothing in Authorization header
		return ""
	}

	return split[1]
}
