// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"encoding/json"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/errorsx"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/driver/config"
)

// Requester is a token endpoint's request context.
//
// swagger:ignore
type Requester struct {
	// ClientID is the identifier of the OAuth 2.0 client.
	ClientID string `json:"client_id"`
	// GrantedScopes is the list of scopes granted to the OAuth 2.0 client.
	GrantedScopes []string `json:"granted_scopes"`
	// GrantedAudience is the list of audiences granted to the OAuth 2.0 client.
	GrantedAudience []string `json:"granted_audience"`
	// GrantTypes is the requests grant types.
	GrantTypes []string `json:"grant_types"`
}

// RefreshTokenHookRequest is the request body sent to the refresh token hook.
//
// swagger:ignore
type RefreshTokenHookRequest struct {
	// Subject is the identifier of the authenticated end-user.
	Subject string `json:"subject"`
	// Session is the request's session..
	Session *Session `json:"session"`
	// Requester is a token endpoint's request context.
	Requester Requester `json:"requester"`
	// ClientID is the identifier of the OAuth 2.0 client.
	ClientID string `json:"client_id"`
	// GrantedScopes is the list of scopes granted to the OAuth 2.0 client.
	GrantedScopes []string `json:"granted_scopes"`
	// GrantedAudience is the list of audiences granted to the OAuth 2.0 client.
	GrantedAudience []string `json:"granted_audience"`
}

// RefreshTokenHook is an AccessRequestHook called for `refresh_token` grant type.
func RefreshTokenHook(reg interface {
	config.Provider
	x.HTTPClientProvider
}) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookConfig := reg.Config().TokenRefreshHookConfig(ctx)
		if hookConfig == nil {
			return nil
		}

		if !requester.GetGrantTypes().ExactOne("refresh_token") {
			return nil
		}

		session, ok := requester.GetSession().(*Session)
		if !ok {
			return nil
		}

		requesterInfo := Requester{
			ClientID:        requester.GetClient().GetID(),
			GrantedScopes:   requester.GetGrantedScopes(),
			GrantedAudience: requester.GetGrantedAudience(),
			GrantTypes:      requester.GetGrantTypes(),
		}

		reqBody := RefreshTokenHookRequest{
			Session:         session,
			Requester:       requesterInfo,
			Subject:         session.GetSubject(),
			ClientID:        requester.GetClient().GetID(),
			GrantedScopes:   requester.GetGrantedScopes(),
			GrantedAudience: requester.GetGrantedAudience(),
		}

		reqBodyBytes, err := json.Marshal(&reqBody)
		if err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDescription("An error occurred while encoding the token hook.").
					WithDebugf("Unable to encode the token hook body: %s", err),
			)
		}

		err = executeHookAndUpdateSession(ctx, reg, hookConfig, reqBodyBytes, session)
		if err != nil {
			return err
		}

		return nil
	}
}
