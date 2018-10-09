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
 *
 */

package oauth2

import (
	"context"
)

// TokenRevocationStorage provides the storage implementation
// as specified in: https://tools.ietf.org/html/rfc7009
type TokenRevocationStorage interface {
	RefreshTokenStorage
	AccessTokenStorage

	// RevokeRefreshToken revokes a refresh token as specified in:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	// If the particular
	// token is a refresh token and the authorization server supports the
	// revocation of access tokens, then the authorization server SHOULD
	// also invalidate all access tokens based on the same authorization
	// grant (see Implementation Note).
	RevokeRefreshToken(ctx context.Context, requestID string) error

	// RevokeAccessToken revokes an access token as specified in:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	// If the token passed to the request
	// is an access token, the server MAY revoke the respective refresh
	// token as well.
	RevokeAccessToken(ctx context.Context, requestID string) error
}
