package trust

import (
	"time"

	"gopkg.in/square/go-jose.v2"
)

type createGrantRequest struct {
	// Issuer identifies the principal that issued the JWT assertion (same as iss claim in jwt).
	Issuer string `json:"issuer"`

	// Subject identifies the principal that is the subject of the JWT.
	Subject string `json:"subject"`

	// AllowAnySubject indicates that the issuer is allowed to have any principal as the subject of the JWT.
	AllowAnySubject bool `json:"allow_any_subject"`

	// Scope contains list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749])
	Scope []string `json:"scope"`

	// PublicKeyJWK contains public key in JWK format issued by Issuer, that will be used to check JWT assertion signature.
	PublicKeyJWK jose.JSONWebKey `json:"jwk"`

	// ExpiresAt indicates, when grant will expire, so we will reject assertion from Issuer targeting Subject.
	ExpiresAt time.Time `json:"expires_at"`
}
