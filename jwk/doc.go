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

// Package jwk implements JSON Web Key management capabilities
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data
// structure that represents a cryptographic key. A JWK Set is a JSON data structure that
// represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality
// to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens).
package jwk

import "github.com/ory/hydra/x"

// swagger:model jsonWebKeySetGeneratorRequest
type createRequest struct {
	// The algorithm to be used for creating the key. Supports "RS256", "ES256", "ES512", "HS512", and "HS256"
	// required: true
	Algorithm string `json:"alg"`

	// The "use" (public key use) parameter identifies the intended use of
	// the public key. The "use" parameter is employed to indicate whether
	// a public key is used for encrypting data or verifying the signature
	// on data. Valid values are "enc" and "sig".
	// required: true
	Use string `json:"use"`

	// The kid of the key to be created
	// required: true
	KeyID string `json:"kid"`
}

// swagger:parameters getJsonWebKey deleteJsonWebKey
type swaggerJsonWebKeyQuery struct {
	// The kid of the desired key
	// in: path
	// required: true
	KID string `json:"kid"`

	// The set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:parameters updateJsonWebKeySet
type swaggerJwkUpdateSet struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body swaggerJSONWebKeySet
}

// swagger:parameters updateJsonWebKey
type swaggerJwkUpdateSetKey struct {
	// The kid of the desired key
	// in: path
	// required: true
	KID string `json:"kid"`

	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body x.JSONWebKey
}

// swagger:parameters createJsonWebKeySet
type swaggerJwkCreateSet struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body createRequest
}

// swagger:parameters getJsonWebKeySet deleteJsonWebKeySet
type swaggerJwkSetQuery struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`
}

// It is important that this model object is named JSONWebKeySet for
// "swagger generate spec" to generate only on definition of a
// JSONWebKeySet. Since one with the same name is previously defined as
// client.Client.JSONWebKeys and this one is last, this one will be
// effectively written in the swagger spec.
//
// swagger:model JSONWebKeySet
type swaggerJSONWebKeySet struct {
	// The value of the "keys" parameter is an array of JWK values.  By
	// default, the order of the JWK values within the array does not imply
	// an order of preference among them, although applications of JWK Sets
	// can choose to assign a meaning to the order for their purposes, if
	// desired.
	Keys []x.JSONWebKey `json:"keys"`
}
