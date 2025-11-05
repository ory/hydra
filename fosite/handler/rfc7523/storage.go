// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc7523

import (
	"context"
	"time"

	"github.com/go-jose/go-jose/v3"
)

// RFC7523KeyStorage holds information needed to validate jwt assertion in authorization grants.
type RFC7523KeyStorage interface {
	// GetPublicKey returns public key, issued by 'issuer', and assigned for subject. Public key is used to check
	// signature of jwt assertion in authorization grants.
	GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (*jose.JSONWebKey, error)

	// GetPublicKeys returns public key, set issued by 'issuer', and assigned for subject.
	GetPublicKeys(ctx context.Context, issuer string, subject string) (*jose.JSONWebKeySet, error)

	// GetPublicKeyScopes returns assigned scope for assertion, identified by public key, issued by 'issuer'.
	GetPublicKeyScopes(ctx context.Context, issuer string, subject string, keyId string) ([]string, error)

	// IsJWTUsed returns true, if JWT is not known yet or it can not be considered valid, because it must be already
	// expired.
	IsJWTUsed(ctx context.Context, jti string) (bool, error)

	// MarkJWTUsedForTime marks JWT as used for a time passed in exp parameter. This helps ensure that JWTs are not
	// replayed by maintaining the set of used "jti" values for the length of time for which the JWT would be
	// considered valid based on the applicable "exp" instant. (https://tools.ietf.org/html/rfc7523#section-3)
	MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) error
}

type RFC7523KeyStorageProvider interface {
	RFC7523KeyStorage() RFC7523KeyStorage
}
