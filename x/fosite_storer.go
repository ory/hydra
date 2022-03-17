/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package x

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/pkce"
	"github.com/ory/fosite/handler/rfc7523"
)

type FositeStorer interface {
	fosite.Storage
	oauth2.CoreStorage
	openid.OpenIDConnectRequestStorage
	pkce.PKCERequestStorage
	rfc7523.RFC7523KeyStorage

	RevokeRefreshToken(ctx context.Context, requestID string) error

	RevokeAccessToken(ctx context.Context, requestID string) error

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
}
