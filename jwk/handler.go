// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"encoding/json"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/ory/herodot"
	"github.com/ory/x/httprouterx"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/x/urlx"

	"github.com/ory/x/errorsx"

	"github.com/ory/x/stringslice"

	"github.com/ory/hydra/v2/x"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/julienschmidt/httprouter"
)

const (
	KeyHandlerPath    = "/keys"
	WellKnownKeysPath = "/.well-known/jwks.json"
)

type Handler struct {
	r InternalRegistry
}

// JSON Web Key Set
//
// swagger:model jsonWebKeySet
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type jsonWebKeySet struct {
	// List of JSON Web Keys
	//
	// The value of the "keys" parameter is an array of JSON Web Key (JWK)
	// values. By default, the order of the JWK values within the array does
	// not imply an order of preference among them, although applications
	// of JWK Sets can choose to assign a meaning to the order for their
	// purposes, if desired.
	Keys []x.JSONWebKey `json:"keys"`
}

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{r: r}
}

func (h *Handler) SetRoutes(admin *httprouterx.RouterAdmin, public *httprouterx.RouterPublic, corsMiddleware func(http.Handler) http.Handler) {
	public.Handler("OPTIONS", WellKnownKeysPath, corsMiddleware(http.HandlerFunc(h.handleOptions)))
	public.Handler("GET", WellKnownKeysPath, corsMiddleware(http.HandlerFunc(h.discoverJsonWebKeys)))

	admin.GET(KeyHandlerPath+"/:set/:key", h.getJsonWebKey)
	admin.GET(KeyHandlerPath+"/:set", h.getJsonWebKeySet)

	admin.POST(KeyHandlerPath+"/:set", h.createJsonWebKeySet)

	admin.PUT(KeyHandlerPath+"/:set/:key", h.adminUpdateJsonWebKey)
	admin.PUT(KeyHandlerPath+"/:set", h.setJsonWebKeySet)

	admin.DELETE(KeyHandlerPath+"/:set/:key", h.deleteJsonWebKey)
	admin.DELETE(KeyHandlerPath+"/:set", h.adminDeleteJsonWebKeySet)
}

// swagger:route GET /.well-known/jwks.json wellknown discoverJsonWebKeys
//
// # Discover Well-Known JSON Web Keys
//
// This endpoint returns JSON Web Keys required to verifying OpenID Connect ID Tokens and,
// if enabled, OAuth 2.0 JWT Access Tokens. This endpoint can be used with client libraries like
// [node-jwks-rsa](https://github.com/auth0/node-jwks-rsa) among others.
//
// Adding custom keys requires first creating a keyset via the createJsonWebKeySet operation,
// and then configuring the webfinger.jwks.broadcast_keys configuration value to include the keyset name.
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
//	  200: jsonWebKeySet
//	  default: errorOAuth2
func (h *Handler) discoverJsonWebKeys(w http.ResponseWriter, r *http.Request) {
	eg, ctx := errgroup.WithContext(r.Context())
	wellKnownKeys := stringslice.Unique(h.r.Config().WellKnownKeys(ctx))
	keys := make(chan *jose.JSONWebKeySet, len(wellKnownKeys))
	for _, set := range wellKnownKeys {
		set := set
		eg.Go(func() error {
			k, err := h.r.KeyManager().GetKeySet(ctx, set)
			if errors.Is(err, x.ErrNotFound) {
				h.r.Logger().Warnf("JSON Web Key Set %q does not exist yet, generating new key pair...", set)
				k, err = h.r.KeyManager().GenerateAndPersistKeySet(ctx, set, uuid.Must(uuid.NewV4()).String(), string(jose.RS256), "sig")
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			keys <- ExcludePrivateKeys(k)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	close(keys)
	var jwks jose.JSONWebKeySet
	for k := range keys {
		jwks.Keys = append(jwks.Keys, k.Keys...)
	}

	h.r.Writer().Write(w, r, &jwks)
}

// Get JSON Web Key Request
//
// swagger:parameters getJsonWebKey
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type getJsonWebKey struct {
	// JSON Web Key Set ID
	//
	// in: path
	// required: true
	Set string `json:"set"`

	// JSON Web Key ID
	//
	// in: path
	// required: true
	KID string `json:"kid"`
}

// swagger:route GET /admin/keys/{set}/{kid} jwk getJsonWebKey
//
// # Get JSON Web Key
//
// This endpoint returns a singular JSON Web Key contained in a set. It is identified by the set and the specific key ID (kid).
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
//	  200: jsonWebKeySet
//	  default: errorOAuth2
func (h *Handler) getJsonWebKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")
	var keyName = ps.ByName("key")

	keys, err := h.r.KeyManager().GetKey(r.Context(), setName, keyName)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	keys = ExcludeOpaquePrivateKeys(keys)

	h.r.Writer().Write(w, r, keys)
}

// Get JSON Web Key Set Parameters
//
// swagger:parameters getJsonWebKeySet
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type getJsonWebKeySet struct {
	// JSON Web Key Set ID
	//
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:route GET /admin/keys/{set} jwk getJsonWebKeySet
//
// # Retrieve a JSON Web Key Set
//
// This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//	  200: jsonWebKeySet
//	  default: errorOAuth2
func (h *Handler) getJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")

	keys, err := h.r.KeyManager().GetKeySet(r.Context(), setName)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	keys = ExcludeOpaquePrivateKeys(keys)

	h.r.Writer().Write(w, r, keys)
}

// Create JSON Web Key Set Request
//
// swagger:parameters createJsonWebKeySet
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type adminCreateJsonWebKeySet struct {
	// The JSON Web Key Set ID
	//
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	// required: true
	Body createJsonWebKeySetBody
}

// Create JSON Web Key Set Request Body
//
// swagger:model createJsonWebKeySet
type createJsonWebKeySetBody struct {
	// JSON Web Key Algorithm
	//
	// The algorithm to be used for creating the key. Supports `RS256`, `ES256`, `ES512`, `HS512`, and `HS256`.
	//
	// required: true
	Algorithm string `json:"alg"`

	// JSON Web Key Use
	//
	// The "use" (public key use) parameter identifies the intended use of
	// the public key. The "use" parameter is employed to indicate whether
	// a public key is used for encrypting data or verifying the signature
	// on data. Valid values are "enc" and "sig".
	// required: true
	Use string `json:"use"`

	// JSON Web Key ID
	//
	// The Key ID of the key to be created.
	//
	// required: true
	KeyID string `json:"kid"`
}

// swagger:route POST /admin/keys/{set} jwk createJsonWebKeySet
//
// # Create JSON Web Key
//
// This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//	  201: jsonWebKeySet
//	  default: errorOAuth2
func (h *Handler) createJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keyRequest createJsonWebKeySetBody
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&keyRequest); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode the request body: %s", err)))
		return
	}

	if keys, err := h.r.KeyManager().GenerateAndPersistKeySet(r.Context(), set, keyRequest.KeyID, keyRequest.Algorithm, keyRequest.Use); err == nil {
		keys = ExcludeOpaquePrivateKeys(keys)
		h.r.Writer().WriteCreated(w, r, urlx.AppendPaths(h.r.Config().IssuerURL(r.Context()), "/keys/"+set).String(), keys)
	} else {
		h.r.Writer().WriteError(w, r, err)
	}
}

// Set JSON Web Key Set Request
//
// swagger:parameters setJsonWebKeySet
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type setJsonWebKeySet struct {
	// The JSON Web Key Set ID
	//
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body jsonWebKeySet
}

// swagger:route PUT /admin/keys/{set} jwk setJsonWebKeySet
//
// # Update a JSON Web Key Set
//
// Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//	  200: jsonWebKeySet
//	  default: errorOAuth2
func (h *Handler) setJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keySet jose.JSONWebKeySet
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&keySet); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode the request body: %s", err)))
		return
	}

	if err := h.r.KeyManager().UpdateKeySet(r.Context(), set, &keySet); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &keySet)
}

// Set JSON Web Key Request
//
// swagger:parameters setJsonWebKey
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type setJsonWebKey struct {
	// The JSON Web Key Set ID
	//
	// in: path
	// required: true
	Set string `json:"set"`

	// JSON Web Key ID
	//
	// in: path
	// required: true
	KID string `json:"kid"`

	// in: body
	Body x.JSONWebKey
}

// swagger:route PUT /admin/keys/{set}/{kid} jwk setJsonWebKey
//
// # Set JSON Web Key
//
// Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//	  200: jsonWebKey
//	  default: errorOAuth2
func (h *Handler) adminUpdateJsonWebKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var key jose.JSONWebKey
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(herodot.ErrBadRequest.WithReasonf("Unable to decode the request body: %s", err)))
		return
	}

	if err := h.r.KeyManager().UpdateKey(r.Context(), set, &key); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, key)
}

// Delete JSON Web Key Set Parameters
//
// swagger:parameters deleteJsonWebKeySet
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type deleteJsonWebKeySet struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:route DELETE /admin/keys/{set} jwk deleteJsonWebKeySet
//
// # Delete JSON Web Key Set
//
// Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//	  default: errorOAuth2
func (h *Handler) adminDeleteJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")

	if err := h.r.KeyManager().DeleteKeySet(r.Context(), setName); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete JSON Web Key Parameters
//
// swagger:parameters deleteJsonWebKey
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type deleteJsonWebKey struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`

	// The JSON Web Key ID (kid)
	//
	// in: path
	// required: true
	KID string `json:"kid"`
}

// swagger:route DELETE /admin/keys/{set}/{kid} jwk deleteJsonWebKey
//
// # Delete JSON Web Key
//
// Use this endpoint to delete a single JSON Web Key.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A
// JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses
// this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens),
// and allows storing user-defined keys as well.
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
//	  default: errorOAuth2
func (h *Handler) deleteJsonWebKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")
	var keyName = ps.ByName("key")

	if err := h.r.KeyManager().DeleteKey(r.Context(), setName, keyName); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// This function will not be called, OPTIONS request will be handled by cors
// this is just a placeholder.
func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {}
