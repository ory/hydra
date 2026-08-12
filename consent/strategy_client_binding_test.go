// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal/testhelpers"
)

// TestVerifyConsentUsesParsedClient proves that consent verification binds the
// flow to the parsed request client and not to the raw query parameter. On a
// POST /oauth2/auth the body value resolved by fosite must win over the query
// value, otherwise the two can disagree and a consent granted for one client is
// applied to a request made by another.
func TestVerifyConsentUsesParsedClient(t *testing.T) {
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t)

	f := &flow.Flow{
		ID:          "client-binding-flow",
		NID:         reg.Networker().NetworkID(ctx),
		Client:      &client.Client{ID: "client-a"},
		Subject:     "subject",
		RequestedAt: time.Now().UTC(),
		State:       flow.FlowStateConsentUnused,
	}

	verifier, err := f.ToConsentVerifier(ctx, reg)
	require.NoError(t, err)

	// Query carries the flow's own client id (client-a), but the parsed request
	// resolves to a different client (client-b).
	r := httptest.NewRequest("POST", "/oauth2/auth?client_id=client-a", nil)
	req := fosite.NewAuthorizeRequest()
	req.Client = &client.Client{ID: "client-b"}

	s := consent.NewStrategy(reg)
	_, err = consent.VerifyConsentForTest(s, ctx, r, req, verifier)

	require.Error(t, err)
	assert.ErrorIs(t, err, fosite.ErrInvalidClient)
	assert.Contains(t, fosite.ErrorToRFC6749Error(err).HintField, "does not match")
}
