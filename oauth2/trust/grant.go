// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"time"

	"github.com/gofrs/uuid"
)

type Grant struct {
	ID uuid.UUID `json:"id"`

	// Issuer identifies the principal that issued the JWT assertion (same as iss claim in jwt).
	Issuer string `json:"issuer"`

	// Subject identifies the principal that is the subject of the JWT.
	Subject string `json:"subject"`

	// AllowAnySubject indicates that the issuer is allowed to have any principal as the subject of the JWT.
	AllowAnySubject bool `json:"allow_any_subject"`

	// Scope contains list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749])
	Scope []string `json:"scope"`

	// PublicKeys contains information about public key issued by Issuer, that will be used to check JWT assertion signature.
	PublicKey PublicKey `json:"public_key"`

	// CreatedAt indicates, when grant was created.
	CreatedAt time.Time `json:"created_at"`

	// ExpiresAt indicates, when grant will expire, so we will reject assertion from Issuer targeting Subject.
	ExpiresAt time.Time `json:"expires_at"`
}

type PublicKey struct {
	// Set is basically a name for a group(set) of keys. Will be the same as Issuer in grant.
	Set string `json:"set"`

	// KeyID is key unique identifier (same as kid header in jws/jwt).
	KeyID string `json:"kid"`
}
