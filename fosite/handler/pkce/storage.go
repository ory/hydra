// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pkce

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
)

type (
	// PKCERequestStorage handles storage of PKCE requests.
	//
	// The /token-side methods receive both the original authorization code and its
	// signature. See AuthorizeCodeStorage.GetAuthorizeCodeSession for the rationale on
	// passing both — commercial Hydra's AEAD storage decodes session state from the code
	// itself, while the default SQL persister keys by signature.
	PKCERequestStorage interface {
		// GetPKCERequestSession returns the PKCE request that was saved during the
		// authorize step. code is the original authorization code; signature is the
		// lookup key.
		GetPKCERequestSession(ctx context.Context, code, signature string, session fosite.Session) (fosite.Requester, error)
		CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error
		// DeletePKCERequestSession removes the PKCE row for an authorization code. code is
		// the original authorization code; signature is the lookup key. Storages that elide
		// the PKCE row entirely for some code formats (e.g., AEAD codes that inline the
		// challenge) use code to dispatch.
		DeletePKCERequestSession(ctx context.Context, code, signature string) error
	}
	PKCERequestStorageProvider interface {
		PKCERequestStorage() PKCERequestStorage
	}
)
