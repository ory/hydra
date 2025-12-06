// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"time"

	"github.com/ory/hydra/v2/fosite"
)

type HandleHelperConfigProvider interface {
	fosite.AccessTokenLifespanProvider
	fosite.RefreshTokenLifespanProvider
}

type HandleHelper struct {
	AccessTokenStrategy AccessTokenStrategy
	Storage             AccessTokenStorageProvider
	Config              HandleHelperConfigProvider
}

func (h *HandleHelper) IssueAccessToken(ctx context.Context, defaultLifespan time.Duration, requester fosite.AccessRequester, responder fosite.AccessResponder) (signature string, err error) {
	token, signature, err := h.AccessTokenStrategy.GenerateAccessToken(ctx, requester)
	if err != nil {
		return "", err
	} else if err := h.Storage.AccessTokenStorage().CreateAccessTokenSession(ctx, signature, requester.Sanitize([]string{})); err != nil {
		return "", err
	}

	responder.SetAccessToken(token)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(getExpiresIn(requester, fosite.AccessToken, defaultLifespan, time.Now().UTC()))
	responder.SetScopes(requester.GetGrantedScopes())
	return signature, nil
}

func getExpiresIn(r fosite.Requester, key fosite.TokenType, defaultLifespan time.Duration, now time.Time) time.Duration {
	if r.GetSession().GetExpiresAt(key).IsZero() {
		return defaultLifespan
	}
	return time.Duration(r.GetSession().GetExpiresAt(key).UnixNano() - now.UnixNano())
}
