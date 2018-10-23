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

package fosite

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type TokenIntrospector interface {
	IntrospectToken(ctx context.Context, token string, tokenType TokenType, accessRequest AccessRequester, scopes []string) (TokenType, error)
}

func AccessTokenFromRequest(req *http.Request) string {
	// According to https://tools.ietf.org/html/rfc6750 you can pass tokens through:
	// - Form-Encoded Body Parameter. Recommended, more likely to appear. e.g.: Authorization: Bearer mytoken123
	// - URI Query Parameter e.g. access_token=mytoken123

	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		// Nothing in Authorization header, try access_token
		// Empty string returned if there's no such parameter
		if err := req.ParseMultipartForm(1 << 20); err != nil && err != http.ErrNotMultipart {
			return ""
		}
		return req.Form.Get("access_token")
	}

	return split[1]
}

func (f *Fosite) IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scopes ...string) (TokenType, AccessRequester, error) {
	var found bool = false
	var foundTokenType TokenType = ""

	ar := NewAccessRequest(session)
	for _, validator := range f.TokenIntrospectionHandlers {
		tt, err := validator.IntrospectToken(ctx, token, tokenType, ar, scopes)
		if err := errors.Cause(err); err == nil {
			found = true
			foundTokenType = tt
		} else if err.Error() == ErrUnknownRequest.Error() {
			// Nothing to do
		} else if err != nil {
			rfcerr := ErrorToRFC6749Error(err)
			return "", nil, errors.WithStack(rfcerr)
		}
	}

	if !found {
		return "", nil, errors.WithStack(ErrRequestUnauthorized.WithHint("Unable to find a suitable validation strategy for the token, thus it is invalid."))
	}

	return foundTokenType, ar, nil
}
