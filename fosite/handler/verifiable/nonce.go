// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package verifiable

import (
	"context"
	"time"
)

type NonceManager interface {
	// NewNonce creates a new nonce bound to the access token valid until the given expiry time.
	NewNonce(ctx context.Context, accessToken string, expiresAt time.Time) (string, error)

	// IsNonceValid checks if the given nonce is valid for the given access token and not expired.
	IsNonceValid(ctx context.Context, accessToken string, nonce string) error
}

type NonceManagerProvider interface {
	NonceManager() NonceManager
}
