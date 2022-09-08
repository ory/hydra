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

package jwk

import (
	"encoding/json"
	"net/http"

	"github.com/ory/x/httprouterx"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/x/urlx"

	"github.com/ory/x/errorsx"

	"github.com/ory/x/stringslice"

	"github.com/ory/hydra/x"

	"github.com/julienschmidt/httprouter"
	jose "gopkg.in/square/go-jose.v2"
)

const (
	KeyHandlerPath    = "/keys"
	WellKnownKeysPath = "/.well-known/jwks.json"
)

type Handler struct {
	r InternalRegistry
}

// It is important that this model object is named JSONWebKeySet for
// "swagger generate spec" to generate only on definition of a
// JSONWebKeySet. Since one with the same name is previously defined as
// client.Client.JSONWebKeys and this one is last, this one will be
// effectively written in the swagger spec.
//
// swagger:model jsonWebKeySet
type jsonWebKeySet struct {
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

	admin.GET(KeyHandlerPath+"/:set/:key", h.adminGetJsonWebKey)
	admin.GET(KeyHandlerPath+"/:set", h.adminGetJsonWebKeySet)

	admin.POST(KeyHandlerPath+"/:set", h.Create)

	admin.PUT(KeyHandlerPath+"/:set/:key", h.adminUpdateJsonWebKey)
	admin.PUT(KeyHandlerPath+"/:set", h.adminUpdateJsonWebKeySet)

	admin.DELETE(KeyHandlerPath+"/:set/:key", h.adminDeleteJsonWebKey)
	admin.DELETE(KeyHandlerPath+"/:set", h.adminDeleteJsonWebKeySet)
}

// swagger:route GET /.well-known/jwks.json v0alpha2 discoverJsonWebKeys
//
// Discover JSON Web Keys
//
// This endpoint returns JSON Web Keys required to verifying OpenID Connect ID Tokens and,
// if enabled, OAuth 2.0 JWT Access Tokens. This endpoint can be used with client libraries like
// [node-jwks-rsa](https://github.com/auth0/node-jwks-rsa) among others.
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
//       200: jsonWebKeySet
//       default: oAuth2ApiError
func (h *Handler) discoverJsonWebKeys(w http.ResponseWriter, r *http.Request) {
	var jwks jose.JSONWebKeySet

	ctx := r.Context()
	for _, set := range stringslice.Unique(h.r.Config().WellKnownKeys(ctx)) {
		keys, err := h.r.KeyManager().GetKeySet(ctx, set)
		if errors.Is(err, x.ErrNotFound) {
			h.r.Logger().Warnf("JSON Web Key Set \"%s\" does not exist yet, generating new key pair...", set)
			keys, err = h.r.KeyManager().GenerateAndPersistKeySet(ctx, set, uuid.Must(uuid.NewV4()).String(), string(jose.RS256), "sig")
			if err != nil {
				h.r.Writer().WriteError(w, r, err)
				return
			}
		} else if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}

		keys = ExcludePrivateKeys(keys)
		jwks.Keys = append(jwks.Keys, keys.Keys...)
	}

	h.r.Writer().Write(w, r, &jwks)
}

// swagger:parameters adminGetJsonWebKey
type adminGetJsonWebKey struct {
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

// swagger:route GET /admin/keys/{set}/{kid} v0alpha2 adminGetJsonWebKey
//
// Fetch a JSON Web Key
//
// This endpoint returns a singular JSON Web Key. It is identified by the set and the specific key ID (kid).
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
//       200: jsonWebKeySet
//       default: oAuth2ApiError
func (h *Handler) adminGetJsonWebKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// swagger:parameters adminGetJsonWebKeySet
type adminGetJsonWebKeySet struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:route GET /admin/keys/{set} v0alpha2 adminGetJsonWebKeySet
//
// Retrieve a JSON Web Key Set
//
// This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//       200: jsonWebKeySet
//       default: oAuth2ApiError
func (h *Handler) adminGetJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")

	keys, err := h.r.KeyManager().GetKeySet(r.Context(), setName)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	keys = ExcludeOpaquePrivateKeys(keys)

	h.r.Writer().Write(w, r, keys)
}

// swagger:parameters adminCreateJsonWebKeySet
type adminCreateJsonWebKeySet struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	// required: true
	Body adminCreateJsonWebKeySetBody
}

// swagger:model adminCreateJsonWebKeySetBody
type adminCreateJsonWebKeySetBody struct {
	// The algorithm to be used for creating the key. Supports "RS256", "ES256", "ES512", "HS512", and "HS256"
	//
	// required: true
	Algorithm string `json:"alg"`

	// The "use" (public key use) parameter identifies the intended use of
	// the public key. The "use" parameter is employed to indicate whether
	// a public key is used for encrypting data or verifying the signature
	// on data. Valid values are "enc" and "sig".
	// required: true
	Use string `json:"use"`

	// The kid of the key to be created
	//
	// required: true
	KeyID string `json:"kid"`
}

// swagger:route POST /admin/keys/{set} v0alpha2 adminCreateJsonWebKeySet
//
// Generate a New JSON Web Key
//
// This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//       201: jsonWebKeySet
//       default: oAuth2ApiError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keyRequest adminCreateJsonWebKeySetBody
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&keyRequest); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
	}

	if keys, err := h.r.KeyManager().GenerateAndPersistKeySet(r.Context(), set, keyRequest.KeyID, keyRequest.Algorithm, keyRequest.Use); err == nil {
		keys = ExcludeOpaquePrivateKeys(keys)
		h.r.Writer().WriteCreated(w, r, urlx.AppendPaths(h.r.Config().IssuerURL(r.Context()), "/keys/"+set).String(), keys)
	} else {
		h.r.Writer().WriteError(w, r, err)
	}
}

// swagger:parameters adminUpdateJsonWebKeySet
type adminUpdateJsonWebKeySet struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body jsonWebKeySet
}

// swagger:route PUT /admin/keys/{set} v0alpha2 adminUpdateJsonWebKeySet
//
// Update a JSON Web Key Set
//
// Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//       200: jsonWebKeySet
//       default: oAuth2ApiError
func (h *Handler) adminUpdateJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keySet jose.JSONWebKeySet
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&keySet); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	if err := h.r.KeyManager().UpdateKeySet(r.Context(), set, &keySet); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &keySet)
}

// swagger:parameters adminUpdateJsonWebKey
type adminUpdateJsonWebKey struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`

	// The JSON Web Key ID (kid)
	//
	// in: path
	// required: true
	KID string `json:"kid"`

	// in: body
	Body x.JSONWebKey
}

// swagger:route PUT /admin/keys/{set}/{kid} v0alpha2 adminUpdateJsonWebKey
//
// Update a JSON Web Key
//
// Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
//       200: jsonWebKey
//       default: oAuth2ApiError
func (h *Handler) adminUpdateJsonWebKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var key jose.JSONWebKey
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	if err := h.r.KeyManager().UpdateKey(r.Context(), set, &key); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, key)
}

// swagger:parameters adminDeleteJsonWebKeySet
type adminDeleteJsonWebKeySet struct {
	// The JSON Web Key Set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:route DELETE /admin/keys/{set} v0alpha2 adminDeleteJsonWebKeySet
//
// Delete a JSON Web Key Set
//
// Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
func (h *Handler) adminDeleteJsonWebKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")

	if err := h.r.KeyManager().DeleteKeySet(r.Context(), setName); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:parameters adminDeleteJsonWebKey
type adminDeleteJsonWebKey struct {
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

// swagger:route DELETE /admin/keys/{set}/{kid} v0alpha2 adminDeleteJsonWebKey
//
// Delete a JSON Web Key
//
// Use this endpoint to delete a single JSON Web Key.
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
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
func (h *Handler) adminDeleteJsonWebKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
