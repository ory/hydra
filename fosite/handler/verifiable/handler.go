// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package verifiable

import (
	"context"
	"time"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/x/errorsx"
)

const (
	draftScope         = "userinfo_credential_draft_00"
	draftNonceField    = "c_nonce_draft_00"
	draftNonceExpField = "c_nonce_expires_in_draft_00"
)

type Handler struct {
	Config interface {
		fosite.VerifiableCredentialsNonceLifespanProvider
	}
	NonceManager
}

var _ fosite.TokenEndpointHandler = (*Handler)(nil)

func (c *Handler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !c.CanHandleTokenEndpointRequest(ctx, request) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	return nil
}

func (c *Handler) PopulateTokenEndpointResponse(
	ctx context.Context,
	request fosite.AccessRequester,
	response fosite.AccessResponder,
) error {
	if !c.CanHandleTokenEndpointRequest(ctx, request) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}

	lifespan := c.Config.GetVerifiableCredentialsNonceLifespan(ctx)
	expiry := time.Now().UTC().Add(lifespan)
	nonce, err := c.NewNonce(ctx, response.GetAccessToken(), expiry)
	if err != nil {
		return err
	}

	response.SetExtra(draftNonceField, nonce)
	response.SetExtra(draftNonceExpField, int64(lifespan.Seconds()))

	return nil
}

func (c *Handler) CanSkipClientAuth(context.Context, fosite.AccessRequester) bool {
	return false
}

func (c *Handler) CanHandleTokenEndpointRequest(_ context.Context, requester fosite.AccessRequester) bool {
	return requester.GetGrantedScopes().Has("openid", draftScope)
}
