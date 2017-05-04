package jwk

// swagger:parameters getJwkSetKey deleteJwkKey
type swaggerJwkSetKeyQuery struct {
	// The kid of the desired key
	// in: path
	// required: true
	KID string `json:"kid"`

	// The set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:parameters updateJwkSet
type swaggerJwkUpdateSet struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body swaggerJSONWebKeySet
}

// swagger:parameters updateJwkKey
type swaggerJwkUpdateKey struct {
	// The kid of the desired key
	// in: path
	// required: true
	KID string `json:"kid"`

	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body swaggerJSONWebKeySet
}

// swagger:parameters createJwkKey
type swaggerJwkCreateKey struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body createRequest
}

// swagger:parameters getJwkSet deleteJwkSet
type swaggerJwkSetQuery struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:model jwkSet
type swaggerJSONWebKeySet struct {
	// The value of the "keys" parameter is an array of JWK values.  By
	// default, the order of the JWK values within the array does not imply
	// an order of preference among them, although applications of JWK Sets
	// can choose to assign a meaning to the order for their purposes, if
	// desired.
	Keys []swaggerJSONWebKey `json:"keys"`
}

// swagger:model jwk
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

	K []byte `json:"k,omitempty"`
	X []byte `json:"x,omitempty"`
	Y []byte `json:"y,omitempty"`
	N []byte `json:"n,omitempty"`
	E []byte `json:"e,omitempty"`

	D  []byte `json:"d,omitempty"`
	P  []byte `json:"p,omitempty"`
	Q  []byte `json:"q,omitempty"`
	Dp []byte `json:"dp,omitempty"`
	Dq []byte `json:"dq,omitempty"`
	Qi []byte `json:"qi,omitempty"`
}
