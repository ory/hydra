// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"
	"net/http"

	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/fosite"
)

// VerifyConsentForTest exposes the unexported verifyConsent for testing.
func VerifyConsentForTest(s Strategy, ctx context.Context, r *http.Request, req fosite.AuthorizeRequester, verifier string) (*flow.Flow, error) {
	return s.(*defaultStrategy).verifyConsent(ctx, nil, r, req, verifier)
}
