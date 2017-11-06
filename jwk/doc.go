// Package jwk implements JSON Web Key management capabilities
//
// A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data
// structure that represents a cryptographic key. A JWK Set is a JSON data structure that
// represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality
// to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens).
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

package jwk

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
	Body swaggerJSONWebKey
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

// swagger:model jsonWebKeySet
type swaggerJSONWebKeySet struct {
	// The value of the "keys" parameter is an array of JWK values.  By
	// default, the order of the JWK values within the array does not imply
	// an order of preference among them, although applications of JWK Sets
	// can choose to assign a meaning to the order for their purposes, if
	// desired.
	Keys []swaggerJSONWebKey `json:"keys"`
}

// swagger:model jsonWebKey
type swaggerJSONWebKey struct {
	//  The "use" (public key use) parameter identifies the intended use of
	// the public key. The "use" parameter is employed to indicate whether
	// a public key is used for encrypting data or verifying the signature
	// on data. Values are commonly "sig" (signature) or "enc" (encryption).
	Use string `json:"use,omitempty"`

	// The "kty" (key type) parameter identifies the cryptographic algorithm
	// family used with the key, such as "RSA" or "EC". "kty" values should
	// either be registered in the IANA "JSON Web Key Types" registry
	// established by [JWA] or be a value that contains a Collision-
	// Resistant Name.  The "kty" value is a case-sensitive string.
	Kty string `json:"kty,omitempty"`

	// The "kid" (key ID) parameter is used to match a specific key.  This
	// is used, for instance, to choose among a set of keys within a JWK Set
	// during key rollover.  The structure of the "kid" value is
	// unspecified.  When "kid" values are used within a JWK Set, different
	// keys within the JWK Set SHOULD use distinct "kid" values.  (One
	// example in which different keys might use the same "kid" value is if
	// they have different "kty" (key type) values but are considered to be
	// equivalent alternatives by the application using them.)  The "kid"
	// value is a case-sensitive string.
	Kid string `json:"kid,omitempty"`

	Crv string `json:"crv,omitempty"`

	//  The "alg" (algorithm) parameter identifies the algorithm intended for
	// use with the key.  The values used should either be registered in the
	// IANA "JSON Web Signature and Encryption Algorithms" registry
	// established by [JWA] or be a value that contains a Collision-
	// Resistant Name.
	Alg string `json:"alg,omitempty"`

	// The "x5c" (X.509 certificate chain) parameter contains a chain of one
	// or more PKIX certificates [RFC5280].  The certificate chain is
	// represented as a JSON array of certificate value strings.  Each
	// string in the array is a base64-encoded (Section 4 of [RFC4648] --
	// not base64url-encoded) DER [ITU.X690.1994] PKIX certificate value.
	// The PKIX certificate containing the key value MUST be the first
	// certificate.
	X5c []string `json:"x5c,omitempty"`

	K string `json:"k,omitempty"`
	X string `json:"x,omitempty"`
	Y string `json:"y,omitempty"`
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`

	D  string `json:"d,omitempty"`
	P  string `json:"p,omitempty"`
	Q  string `json:"q,omitempty"`
	Dp string `json:"dp,omitempty"`
	Dq string `json:"dq,omitempty"`
	Qi string `json:"qi,omitempty"`
}
