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
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/herodot"
)

const (
	IDTokenKeyName    = "hydra.openid.id-token"
	KeyHandlerPath    = "/keys"
	WellKnownKeysPath = "/.well-known/jwks.json"
)

type Handler struct {
	Manager       Manager
	Generators    map[string]KeyGenerator
	H             herodot.Writer
	WellKnownKeys []string
}

func NewHandler(
	manager Manager,
	generators map[string]KeyGenerator,
	h herodot.Writer,
	wellKnownKeys []string,
) *Handler {
	return &Handler{
		Manager:       manager,
		Generators:    generators,
		H:             h,
		WellKnownKeys: append(wellKnownKeys, IDTokenKeyName),
	}
}

func (h *Handler) GetGenerators() map[string]KeyGenerator {
	if h.Generators == nil || len(h.Generators) == 0 {
		h.Generators = map[string]KeyGenerator{
			"RS256": &RS256Generator{},
			"ES512": &ECDSA512Generator{},
			"HS256": &HS256Generator{},
			"HS512": &HS512Generator{},
		}
	}
	return h.Generators
}

func (h *Handler) SetRoutes(frontend, backend *httprouter.Router, corsMiddleware func(http.Handler) http.Handler) {
	frontend.Handler("OPTIONS", WellKnownKeysPath, corsMiddleware(http.HandlerFunc(h.handleOptions)))
	frontend.Handler("GET", WellKnownKeysPath, corsMiddleware(http.HandlerFunc(h.WellKnown)))

	backend.GET(KeyHandlerPath+"/:set/:key", h.GetKey)
	backend.GET(KeyHandlerPath+"/:set", h.GetKeySet)

	backend.POST(KeyHandlerPath+"/:set", h.Create)

	backend.PUT(KeyHandlerPath+"/:set/:key", h.UpdateKey)
	backend.PUT(KeyHandlerPath+"/:set", h.UpdateKeySet)

	backend.DELETE(KeyHandlerPath+"/:set/:key", h.DeleteKey)
	backend.DELETE(KeyHandlerPath+"/:set", h.DeleteKeySet)
}

// swagger:route GET /.well-known/jwks.json public wellKnown
//
// JSON Web Keys Discovery
//
// This endpoint returns JSON Web Keys to be used as public keys for verifying OpenID Connect ID Tokens and,
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
//       200: JSONWebKeySet
//       500: genericError
func (h *Handler) WellKnown(w http.ResponseWriter, r *http.Request) {
	var jwks jose.JSONWebKeySet

	for _, set := range h.WellKnownKeys {
		keys, err := h.Manager.GetKeySet(r.Context(), set)
		if err != nil {
			h.H.WriteError(w, r, err)
			return
		}

		keys, err = FindKeysByPrefix(keys, "public")
		if err != nil {
			h.H.WriteError(w, r, err)
			return
		}

		jwks.Keys = append(jwks.Keys, keys.Keys...)
	}

	h.H.Write(w, r, &jwks)
}

// swagger:route GET /keys/{set}/{kid} admin getJsonWebKey
//
// Fetch a JSON Web Key
//
// This endpoint returns a singular JSON Web Key, identified by the set and the specific key ID (kid).
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
//       200: JSONWebKeySet
//       404: genericError
//       500: genericError
func (h *Handler) GetKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")
	var keyName = ps.ByName("key")

	keys, err := h.Manager.GetKey(r.Context(), setName, keyName)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, keys)
}

// swagger:route GET /keys/{set} admin getJsonWebKeySet
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
//       200: JSONWebKeySet
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) GetKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")

	keys, err := h.Manager.GetKeySet(r.Context(), setName)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, keys)
}

// swagger:route POST /keys/{set} admin createJsonWebKeySet
//
// Generate a new JSON Web Key
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
//       200: JSONWebKeySet
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keyRequest createRequest
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&keyRequest); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
	}

	generator, found := h.GetGenerators()[keyRequest.Algorithm]
	if !found {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.Errorf("Generator %s unknown", keyRequest.Algorithm))
		return
	}

	keys, err := generator.Generate(keyRequest.KeyID, keyRequest.Use)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.AddKeySet(r.Context(), set, keys); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.WriteCreated(w, r, fmt.Sprintf("%s://%s/keys/%s", r.URL.Scheme, r.URL.Host, set), keys)
}

// swagger:route PUT /keys/{set} admin updateJsonWebKeySet
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
//       200: JSONWebKeySet
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) UpdateKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keySet jose.JSONWebKeySet
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&keySet); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.Manager.AddKeySet(r.Context(), set, &keySet); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, &keySet)
}

// swagger:route PUT /keys/{set}/{kid} admin updateJsonWebKey
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
//       200: JSONWebKey
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) UpdateKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var key jose.JSONWebKey
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.Manager.DeleteKey(r.Context(), set, key.KeyID); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.AddKey(r.Context(), set, &key); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, key)
}

// swagger:route DELETE /keys/{set} admin deleteJsonWebKeySet
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
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) DeleteKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")

	if err := h.Manager.DeleteKeySet(r.Context(), setName); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /keys/{set}/{kid} admin deleteJsonWebKey
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
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) DeleteKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var setName = ps.ByName("set")
	var keyName = ps.ByName("key")

	if err := h.Manager.DeleteKey(r.Context(), setName, keyName); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// This function will not be called, OPTIONS request will be handled by cors
// this is just a placeholder.
func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {}
