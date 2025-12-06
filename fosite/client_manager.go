// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"time"
)

// ClientManager defines the (persistent) manager interface for clients.
type ClientManager interface {
	// GetClient loads the client by its ID or returns an error
	// if the client does not exist or another error occurred.
	GetClient(ctx context.Context, id string) (Client, error)
	// ClientAssertionJWTValid returns an error if the JTI is
	// known or the DB check failed and nil if the JTI is not known.
	ClientAssertionJWTValid(ctx context.Context, jti string) error
	// SetClientAssertionJWT marks a JTI as known for the given
	// expiry time. Before inserting the new JTI, it will clean
	// up any existing JTIs that have expired as those tokens can
	// not be replayed due to the expiry.
	SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error
}
