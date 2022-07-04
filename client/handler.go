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
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ory/x/httprouterx"

	"github.com/ory/x/openapix"

	"github.com/ory/x/uuidx"

	"github.com/ory/x/jsonx"
	"github.com/ory/x/urlx"

	"github.com/ory/fosite"

	"github.com/ory/x/errorsx"

	"github.com/ory/herodot"
	"github.com/ory/hydra/x"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Handler struct {
	r InternalRegistry
}

const (
	ClientsHandlerPath    = "/clients"
	DynClientsHandlerPath = "/oauth2/register"
)

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{
		r: r,
	}
}

func (h *Handler) SetRoutes(admin *httprouterx.RouterAdmin, public *httprouterx.RouterPublic) {
	admin.GET(ClientsHandlerPath, h.adminListOAuth2Clients)
	admin.POST(ClientsHandlerPath, h.adminCreateOAuth2Client)
	admin.GET(ClientsHandlerPath+"/:id", h.Get)
	admin.PUT(ClientsHandlerPath+"/:id", h.adminUpdateOAuth2Client)
	admin.PATCH(ClientsHandlerPath+"/:id", h.adminPatchOAuth2Client)
	admin.DELETE(ClientsHandlerPath+"/:id", h.Delete)
	admin.PUT(ClientsHandlerPath+"/:id/lifespans", h.UpdateLifespans)

	public.POST(DynClientsHandlerPath, h.dynamicClientRegistrationCreateOAuth2Client)
	public.GET(DynClientsHandlerPath+"/:id", h.GetDynamicRegistration)
	public.PUT(DynClientsHandlerPath+"/:id", h.dynamicClientRegistrationUpdateOAuth2Client)
	public.DELETE(DynClientsHandlerPath+"/:id", h.DeleteDynamicRegistration)
}

// swagger:parameters adminCreateOAuth2Client
type adminCreateOAuth2Client struct {
	// in: body
	// required: true
	Body Client
}

// swagger:route POST /admin/clients v1 adminCreateOAuth2Client
//
// Create an OAuth 2.0 Client
//
// Create a new OAuth 2.0 client. If you pass `client_secret` the secret is used, otherwise a random secret
// is generated. The secret is echoed in the response. It is not possible to retrieve it later on.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
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
//       default: oAuth2ApiError
func (h *Handler) adminCreateOAuth2Client(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := h.CreateClient(r, h.r.ClientValidator().Validate, false)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	h.r.Writer().WriteCreated(w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

// swagger:parameters dynamicClientRegistrationCreateOAuth2Client
type dynamicClientRegistrationCreateOAuth2Client struct {
	// in: body
	// required: true
	Body Client
}

// swagger:route POST /oauth2/register v1 dynamicClientRegistrationCreateOAuth2Client
//
// Register an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
//
// This endpoint behaves like the administrative counterpart (`createOAuth2Client`) but is capable of facing the
// public internet directly and can be used in self-service. It implements the OpenID Connect
// Dynamic Client Registration Protocol. This feature needs to be enabled in the configuration. This endpoint
// is disabled by default. It can be enabled by an administrator.
//
// Please note that using this endpoint you are not able to choose the `client_secret` nor the `client_id` as those
// values will be server generated when specifying `token_endpoint_auth_method` as `client_secret_basic` or
// `client_secret_post`.
//
// The `client_secret` will be returned in the response and you will not be able to retrieve it later on.
// Write the secret down and keep it somewhere safe.
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
//       default: oAuth2ApiError
func (h *Handler) dynamicClientRegistrationCreateOAuth2Client(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := h.requireDynamicAuth(r); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	c, err := h.CreateClient(r, h.r.ClientValidator().ValidateDynamicRegistration, true)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	h.r.Writer().WriteCreated(w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

func (h *Handler) CreateClient(r *http.Request, validator func(context.Context, *Client) error, isDynamic bool) (*Client, error) {
	var c Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return nil, err
	}

	if isDynamic {
		if c.Secret != "" {
			return nil, errorsx.WithStack(herodot.ErrBadRequest.WithReasonf("It is not allowed to choose your own OAuth2 Client secret."))
		}
	}

	if len(c.LegacyClientID) > 0 {
		return nil, errorsx.WithStack(herodot.ErrBadRequest.WithReason("It is no longer possible to set an OAuth2 Client ID as a user. The system will generate a unique ID for you."))
	}

	c.ID = uuidx.NewV4()
	c.LegacyClientID = c.ID.String()

	if len(c.Secret) == 0 {
		secretb, err := x.GenerateSecret(26)
		if err != nil {
			return nil, err
		}
		c.Secret = string(secretb)
	}

	if err := validator(r.Context(), &c); err != nil {
		return nil, err
	}

	secret := c.Secret
	c.CreatedAt = time.Now().UTC().Round(time.Second)
	c.UpdatedAt = c.CreatedAt

	token, signature, err := h.r.OAuth2HMACStrategy().GenerateAccessToken(r.Context(), nil)
	if err != nil {
		return nil, err
	}

	c.RegistrationAccessToken = token
	c.RegistrationAccessTokenSignature = signature
	c.RegistrationClientURI = urlx.AppendPaths(h.r.Config().PublicURL(r.Context()), DynClientsHandlerPath+"/"+c.GetID()).String()

	if err := h.r.ClientManager().CreateClient(r.Context(), &c); err != nil {
		return nil, err
	}
	c.Secret = ""
	if !c.IsPublic() {
		c.Secret = secret
	}
	return &c, nil
}

// swagger:parameters adminUpdateOAuth2Client
type adminUpdateOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body Client
}

// swagger:route PUT /admin/clients/{id} v1 adminUpdateOAuth2Client
//
// Update an OAuth 2.0 Client
//
// Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret is used, otherwise a random secret
// is generated. The secret is echoed in the response. It is not possible to retrieve it later on.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
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
//       default: oAuth2ApiError
func (h *Handler) adminUpdateOAuth2Client(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode the request body: %s", err)))
		return
	}

	c.LegacyClientID = ps.ByName("id")
	if err := h.updateClient(r.Context(), &c, h.r.ClientValidator().Validate); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &c)
}

func (h *Handler) updateClient(ctx context.Context, c *Client, validator func(context.Context, *Client) error) error {
	var secret string
	if len(c.Secret) > 0 {
		secret = c.Secret
	}

	if err := validator(ctx, c); err != nil {
		return err
	}

	c.UpdatedAt = time.Now().UTC().Round(time.Second)
	if err := h.r.ClientManager().UpdateClient(ctx, c); err != nil {
		return err
	}
	c.Secret = secret
	return nil
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

// swagger:route PUT /oauth2/register/{id} v1 dynamicClientRegistrationUpdateOAuth2Client
//
// Update an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
//
// This endpoint behaves like the administrative counterpart (`updateOAuth2Client`) but is capable of facing the
// public internet directly and can be used in self-service. It implements the OpenID Connect
// Dynamic Client Registration Protocol. This feature needs to be enabled in the configuration. This endpoint
// is disabled by default. It can be enabled by an administrator.
//
// If you pass `client_secret` the secret is used, otherwise a random secret
// is generated. The secret is echoed in the response. It is not possible to retrieve it later on.
//
// To use this endpoint, you will need to present the client's authentication credentials. If the OAuth2 Client
// uses the Token Endpoint Authentication Method `client_secret_post`, you need to present the client secret in the URL query.
// If it uses `client_secret_basic`, present the Client ID and the Client Secret in the Authorization header.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//       bearer:
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2Client
//       default: oAuth2ApiError
//
func (h *Handler) dynamicClientRegistrationUpdateOAuth2Client(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := h.requireDynamicAuth(r); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	client, err := h.ValidDynamicAuth(r, ps)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	var c Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode the request body. Is it valid JSON?").WithDebug(err.Error())))
		return
	}

	if c.Secret != "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(herodot.ErrForbidden.WithReasonf("It is not allowed to choose your own OAuth2 Client secret.")))
		return
	}

	// Regenerate the registration access token
	token, signature, err := h.r.OAuth2HMACStrategy().GenerateAccessToken(r.Context(), nil)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	c.RegistrationAccessToken = token
	c.RegistrationAccessTokenSignature = signature

	c.LegacyClientID = client.GetID()
	if err := h.updateClient(r.Context(), &c, h.r.ClientValidator().ValidateDynamicRegistration); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &c)
}

// swagger:parameters adminPatchOAuth2Client
type adminPatchOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body openapix.JSONPatchDocument
}

// swagger:route PATCH /clients/{id} v1 adminPatchOAuth2Client
//
// Patch an OAuth 2.0 Client
//
// Patch an existing OAuth 2.0 Client. If you pass `client_secret`
// the secret will be updated and returned via the API. This is the
// only time you will be able to retrieve the client secret, so write it down and keep it safe.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
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
//       default: oAuth2ApiError
func (h *Handler) adminPatchOAuth2Client(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	oldSecret := c.Secret

	if err := jsonx.ApplyJSONPatch(patchJSON, c, "/id"); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	// fix for #2869
	// GetConcreteClient returns a client with the hashed secret, however updateClient expects
	// an empty secret if the secret hasn't changed. As such we need to check if the patch has
	// updated the secret or not
	if oldSecret == c.Secret {
		c.Secret = ""
	}

	if err := h.updateClient(r.Context(), c, h.r.ClientValidator().Validate); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, c)
}

// The list of clients and pagination information.
//
// swagger:response adminListOAuth2ClientsResponse
type adminListOAuth2ClientsResponse struct {
	x.PaginationHeaders

	// in:body
	Body []Client
}

// swagger:parameters adminListOAuth2Clients
type adminListOAuth2Clients struct {
	x.PaginationParams

	// The name of the clients to filter by.
	// in: query
	Name string `json:"client_name"`

	// The owner of the clients to filter by.
	// in: query
	Owner string `json:"owner"`
}

// swagger:route GET /clients v1 adminListOAuth2Clients
//
// List OAuth 2.0 Clients
//
// This endpoint lists all clients in the database, and never returns client secrets.
// As a default it lists the first 100 clients. The `limit` parameter can be used to retrieve more clients,
// but it has an upper bound at 500 objects. Pagination should be used to retrieve more than 500 objects.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
//
// The "Link" header is also included in successful responses, which contains one or more links for pagination,
// formatted like so: '<https://project-slug.projects.oryapis.com/admin/clients?limit={limit}&offset={offset}>; rel="{page}"',
// where page is one of the following applicable pages: 'first', 'next', 'last', and 'previous'. Multiple links can
// be included in this header, and will be separated by a comma.
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
//       200: adminListOAuth2ClientsResponse
//       default: oAuth2ApiError
func (h *Handler) adminListOAuth2Clients(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	page, itemsPerPage := x.ParsePagination(r)
	filters := Filter{
		Limit:  itemsPerPage,
		Offset: page * itemsPerPage,
		Name:   r.URL.Query().Get("client_name"),
		Owner:  r.URL.Query().Get("owner"),
	}

	c, err := h.r.ClientManager().GetClients(r.Context(), filters)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if c == nil {
		c = []Client{}
	}

	for k := range c {
		c[k].Secret = ""
	}

	total, err := h.r.ClientManager().CountClients(r.Context())
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	x.PaginationHeader(w, r.URL, int64(total), page, itemsPerPage)
	h.r.Writer().Write(w, r, c)
}

// swagger:parameters adminGetOAuth2Client
type adminGetOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:route GET /clients/{id} v1 adminGetOAuth2Client
//
// Get an OAuth 2.0 Client
//
// Get an OAuth 2.0 client by its ID. This endpoint never returns the client secret.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
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
//       default: oAuth2ApiError
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), id)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	c.Secret = ""
	h.r.Writer().Write(w, r, c)
}

// swagger:parameters dynamicClientRegistrationGetOAuth2Client
type dynamicClientRegistrationGetOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:route GET /oauth2/register/{id} v1 dynamicClientRegistrationGetOAuth2Client
//
// Get an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
//
// This endpoint behaves like the administrative counterpart (`getOAuth2Client`) but is capable of facing the
// public internet directly and can be used in self-service. It implements the OpenID Connect
// Dynamic Client Registration Protocol. This feature needs to be enabled in the configuration. This endpoint
// is disabled by default. It can be enabled by an administrator.
//
// To use this endpoint, you will need to present the client's authentication credentials. If the OAuth2 Client
// uses the Token Endpoint Authentication Method `client_secret_post`, you need to present the client secret in the URL query.
// If it uses `client_secret_basic`, present the Client ID and the Client Secret in the Authorization header.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
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
//       bearer:
//
//     Responses:
//       200: oAuth2Client
//       default: oAuth2ApiError
func (h *Handler) GetDynamicRegistration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := h.requireDynamicAuth(r); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	client, err := h.ValidDynamicAuth(r, ps)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), client.GetID())
	if err != nil {
		err = herodot.ErrUnauthorized.WithReason("The requested OAuth 2.0 client does not exist or you did not provide the necessary credentials")
		h.r.Writer().WriteError(w, r, err)
		return
	}

	c.Secret = ""
	c.Metadata = nil
	h.r.Writer().Write(w, r, c)
}

// swagger:parameters adminDeleteOAuth2Client
type adminDeleteOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:route DELETE /clients/{id} v1 adminDeleteOAuth2Client
//
// Deletes an OAuth 2.0 Client
//
// Delete an existing OAuth 2.0 Client by its ID.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
//
// Make sure that this endpoint is well protected and only callable by first-party components.
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
//       default: oAuth2ApiError
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	if err := h.r.ClientManager().DeleteClient(r.Context(), id); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:parameters dynamicClientRegistrationDeleteOAuth2Client
type dynamicClientRegistrationDeleteOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:route DELETE /oauth2/register/{id} v1 dynamicClientRegistrationDeleteOAuth2Client
//
// Deletes an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
//
// This endpoint behaves like the administrative counterpart (`deleteOAuth2Client`) but is capable of facing the
// public internet directly and can be used in self-service. It implements the OpenID Connect
// Dynamic Client Registration Protocol. This feature needs to be enabled in the configuration. This endpoint
// is disabled by default. It can be enabled by an administrator.
//
// To use this endpoint, you will need to present the client's authentication credentials. If the OAuth2 Client
// uses the Token Endpoint Authentication Method `client_secret_post`, you need to present the client secret in the URL query.
// If it uses `client_secret_basic`, present the Client ID and the Client Secret in the Authorization header.
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       bearer:
//
//     Responses:
//       204: emptyResponse
//       default: oAuth2ApiError
func (h *Handler) DeleteDynamicRegistration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := h.requireDynamicAuth(r); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	client, err := h.ValidDynamicAuth(r, ps)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if err := h.r.ClientManager().DeleteClient(r.Context(), client.GetID()); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ValidDynamicAuth(r *http.Request, ps httprouter.Params) (fosite.Client, error) {
	c, err := h.r.ClientManager().GetConcreteClient(r.Context(), ps.ByName("id"))
	if err != nil {
		return nil, herodot.ErrUnauthorized.
			WithTrace(err).
			WithReason("The requested OAuth 2.0 client does not exist or you provided incorrect credentials.").WithDebug(err.Error())
	}

	if len(c.RegistrationAccessTokenSignature) == 0 {
		return nil, errors.WithStack(herodot.ErrUnauthorized.
			WithReason("The requested OAuth 2.0 client does not exist or you provided incorrect credentials.").WithDebug("The OAuth2 Client does not have a registration access token."))
	}

	token := strings.TrimPrefix(fosite.AccessTokenFromRequest(r), "ory_at_")
	if err := h.r.OAuth2HMACStrategy().Enigma.Validate(r.Context(), token); err != nil {
		return nil, herodot.ErrUnauthorized.
			WithTrace(err).
			WithReason("The requested OAuth 2.0 client does not exist or you provided incorrect credentials.").WithDebug(err.Error())
	}

	signature := h.r.OAuth2HMACStrategy().Enigma.Signature(token)
	if subtle.ConstantTimeCompare([]byte(c.RegistrationAccessTokenSignature), []byte(signature)) == 0 {
		return nil, errors.WithStack(herodot.ErrUnauthorized.
			WithReason("The requested OAuth 2.0 client does not exist or you provided incorrect credentials.").WithDebug("Registration access tokens do not match."))
	}

	return c, nil
}

func (h *Handler) requireDynamicAuth(r *http.Request) *herodot.DefaultError {
	if !h.r.Config().PublicAllowDynamicRegistration(r.Context()) {
		return herodot.ErrNotFound.WithReason("Dynamic registration is not enabled.")
	}
	return nil
}
