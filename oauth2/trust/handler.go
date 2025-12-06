// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/httprouterx"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
	"github.com/ory/x/urlx"
)

const (
	grantJWTBearerPath = "/trust/grants/jwt-bearer/issuers" // #nosec G101
)

type Handler struct {
	registry InternalRegistry
}

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{registry: r}
}

func (h *Handler) SetRoutes(admin *httprouterx.RouterAdmin) {
	admin.GET(grantJWTBearerPath+"/{id}", h.getTrustedOAuth2JwtGrantIssuer)
	admin.GET(grantJWTBearerPath, h.adminListTrustedOAuth2JwtGrantIssuers)
	admin.POST(grantJWTBearerPath, h.trustOAuth2JwtGrantIssuer)
	admin.DELETE(grantJWTBearerPath+"/{id}", h.deleteTrustedOAuth2JwtGrantIssuer)
}

// Trust OAuth2 JWT Bearer Grant Type Issuer Request Body
//
// swagger:model trustOAuth2JwtGrantIssuer
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type trustOAuth2JwtGrantIssuerBody struct {
	// The "issuer" identifies the principal that issued the JWT assertion (same as "iss" claim in JWT).
	//
	// required: true
	// example: https://jwt-idp.example.com
	Issuer string `json:"issuer"`

	// The "subject" identifies the principal that is the subject of the JWT.
	//
	// example: mike@example.com
	Subject string `json:"subject"`

	// The "allow_any_subject" indicates that the issuer is allowed to have any principal as the subject of the JWT.
	AllowAnySubject bool `json:"allow_any_subject"`

	// The "scope" contains list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749])
	//
	// required:true
	// example: ["openid", "offline"]
	Scope []string `json:"scope"`

	// The "jwk" contains public key in JWK format issued by "issuer", that will be used to check JWT assertion signature.
	//
	// required:true
	JWK x.JSONWebKey `json:"jwk"`

	// The "expires_at" indicates, when grant will expire, so we will reject assertion from "issuer" targeting "subject".
	//
	// required:true
	ExpiresAt time.Time `json:"expires_at"`
}

// Trust OAuth2 JWT Bearer Grant Type Issuer Request
//
// swagger:parameters trustOAuth2JwtGrantIssuer
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type trustOAuth2JwtGrantIssuer struct {
	// in: body
	Body trustOAuth2JwtGrantIssuerBody
}

// swagger:route POST /admin/trust/grants/jwt-bearer/issuers oAuth2 trustOAuth2JwtGrantIssuer
//
// # Trust OAuth2 JWT Bearer Grant Type Issuer
//
// Use this endpoint to establish a trust relationship for a JWT issuer
// to perform JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication
// and Authorization Grants [RFC7523](https://datatracker.ietf.org/doc/html/rfc7523).
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  201: trustedOAuth2JwtGrantIssuer
//	  default: genericError
func (h *Handler) trustOAuth2JwtGrantIssuer(w http.ResponseWriter, r *http.Request) {
	var grantRequest createGrantRequest

	if err := json.NewDecoder(r.Body).Decode(&grantRequest); err != nil {
		h.registry.Writer().WriteError(w, r,
			errors.WithStack(&fosite.RFC6749Error{
				ErrorField:       "error",
				DescriptionField: err.Error(),
				CodeField:        http.StatusBadRequest,
			}))
		return
	}

	if err := validateGrant(grantRequest); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	grant := Grant{
		ID:              uuid.Must(uuid.NewV4()),
		Issuer:          grantRequest.Issuer,
		Subject:         grantRequest.Subject,
		AllowAnySubject: grantRequest.AllowAnySubject,
		Scope:           grantRequest.Scope,
		PublicKey: PublicKey{
			Set:   grantRequest.Issuer, // group all keys by issuer, so set=issuer
			KeyID: grantRequest.PublicKeyJWK.KeyID,
		},
		CreatedAt: time.Now().UTC().Round(time.Second),
		ExpiresAt: grantRequest.ExpiresAt.UTC().Round(time.Second),
	}

	if err := h.registry.GrantManager().CreateGrant(r.Context(), grant, grantRequest.PublicKeyJWK); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	h.registry.Writer().WriteCreated(w, r, urlx.MustJoin(grantJWTBearerPath, url.PathEscape(grant.ID.String())), &grant)
}

// Get Trusted OAuth2 JWT Bearer Grant Type Issuer Request
//
// swagger:parameters getTrustedOAuth2JwtGrantIssuer
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type getTrustedOAuth2JwtGrantIssuer struct {
	// The id of the desired grant
	//
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:route GET /admin/trust/grants/jwt-bearer/issuers/{id} oAuth2 getTrustedOAuth2JwtGrantIssuer
//
// # Get Trusted OAuth2 JWT Bearer Grant Type Issuer
//
// Use this endpoint to get a trusted JWT Bearer Grant Type Issuer. The ID is the one returned when you
// created the trust relationship.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: trustedOAuth2JwtGrantIssuer
//	  default: genericError
func (h *Handler) getTrustedOAuth2JwtGrantIssuer(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	id, err := uuid.FromString(rawID)
	if err != nil {
		h.registry.Writer().WriteError(w, r,
			errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to parse parameter id: %v", err)))
		return
	}

	grant, err := h.registry.GrantManager().GetConcreteGrant(r.Context(), id)
	if err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	h.registry.Writer().Write(w, r, grant)
}

// Delete Trusted OAuth2 JWT Bearer Grant Type Issuer Request
//
// swagger:parameters deleteTrustedOAuth2JwtGrantIssuer
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type deleteTrustedOAuth2JwtGrantIssuer struct {
	// The id of the desired grant
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:route DELETE /admin/trust/grants/jwt-bearer/issuers/{id} oAuth2 deleteTrustedOAuth2JwtGrantIssuer
//
// # Delete Trusted OAuth2 JWT Bearer Grant Type Issuer
//
// Use this endpoint to delete trusted JWT Bearer Grant Type Issuer. The ID is the one returned when you
// created the trust relationship.
//
// Once deleted, the associated issuer will no longer be able to perform the JSON Web Token (JWT) Profile
// for OAuth 2.0 Client Authentication and Authorization Grant.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  204: emptyResponse
//	  default: genericError
func (h *Handler) deleteTrustedOAuth2JwtGrantIssuer(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	id, err := uuid.FromString(rawID)
	if err != nil {
		h.registry.Writer().WriteError(w, r,
			errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to parse parameter id: %v", err)))
		return
	}

	if err := h.registry.GrantManager().DeleteGrant(r.Context(), id); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List Trusted OAuth2 JWT Bearer Grant Type Issuers Request
//
// swagger:parameters listTrustedOAuth2JwtGrantIssuers
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type listTrustedOAuth2JwtGrantIssuers struct {
	// If optional "issuer" is supplied, only jwt-bearer grants with this issuer will be returned.
	//
	// in: query
	// required: false
	Issuer string `json:"issuer"`

	keysetpagination.RequestParameters
}

// swagger:route GET /admin/trust/grants/jwt-bearer/issuers oAuth2 listTrustedOAuth2JwtGrantIssuers
//
// # List Trusted OAuth2 JWT Bearer Grant Type Issuers
//
// Use this endpoint to list all trusted JWT Bearer Grant Type Issuers.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: trustedOAuth2JwtGrantIssuers
//	  default: genericError
func (h *Handler) adminListTrustedOAuth2JwtGrantIssuers(w http.ResponseWriter, r *http.Request) {
	optionalIssuer := r.URL.Query().Get("issuer")

	pageKeys := h.registry.Config().GetPaginationEncryptionKeys(r.Context())
	pageOpts, err := keysetpagination.ParseQueryParams(pageKeys, r.URL.Query())
	if err != nil {
		h.registry.Writer().WriteError(w, r,
			errors.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to parse pagination parameters: %v", err)))
		return
	}

	grants, nextPage, err := h.registry.GrantManager().GetGrants(r.Context(), optionalIssuer, pageOpts...)
	if err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}
	if grants == nil {
		grants = []Grant{}
	}

	keysetpagination.SetLinkHeader(w, pageKeys, r.URL, nextPage)
	h.registry.Writer().Write(w, r, grants)
}
