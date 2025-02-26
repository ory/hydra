// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/pkce"
	"github.com/ory/fosite/handler/rfc7523"
	"github.com/ory/fosite/handler/rfc8628"
	"github.com/ory/fosite/handler/verifiable"
)

type FositeStorer interface {
	fosite.Storage
	oauth2.CoreStorage
	oauth2.TokenRevocationStorage
	openid.OpenIDConnectRequestStorage
	pkce.PKCERequestStorage
	rfc7523.RFC7523KeyStorage
	rfc8628.RFC8628CoreStorage
	verifiable.NonceManager
	oauth2.ResourceOwnerPasswordCredentialsGrantStorage

	// flush the access token requests from the database.
	// no data will be deleted after the 'notAfter' timeframe.
	FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

	// flush the login requests from the database.
	// this will address the database long-term growth issues discussed in https://github.com/ory/hydra/issues/1574.
	// no data will be deleted after the 'notAfter' timeframe.
	FlushInactiveLoginConsentRequests(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

	DeleteAccessTokens(ctx context.Context, clientID string) error

	FlushInactiveRefreshTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

	// DeleteOpenIDConnectSession deletes an OpenID Connect session.
	// This is duplicated from Ory Fosite to help against deprecation linting errors.
	DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error

	GetUserCodeSession(context.Context, string, fosite.Session) (fosite.DeviceRequester, error)
	GetDeviceCodeSessionByRequestID(ctx context.Context, requestID string, requester fosite.Session) (fosite.DeviceRequester, string, error)
	UpdateDeviceCodeSessionBySignature(ctx context.Context, requestID string, requester fosite.DeviceRequester) error
}
