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
	"strings"
	"time"

	"context"

	jwtx "github.com/dgrijalva/jwt-go"
	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	"github.com/pkg/errors"
)

// DefaultJWTStrategy is a JWT RS256 strategy.
type DefaultJWTStrategy struct {
	jwt.JWTStrategy
	HMACSHAStrategy *HMACSHAStrategy
	Issuer          string
}

func (h DefaultJWTStrategy) signature(token string) string {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return ""
	}

	return split[2]
}

func (h DefaultJWTStrategy) AccessTokenSignature(token string) string {
	return h.signature(token)
}

func (h *DefaultJWTStrategy) GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(ctx, fosite.AccessToken, requester)
}

func (h *DefaultJWTStrategy) ValidateAccessToken(ctx context.Context, _ fosite.Requester, token string) error {
	_, err := h.validate(ctx, token)
	return err
}

func (h *DefaultJWTStrategy) ValidateJWT(ctx context.Context, tokenType fosite.TokenType, token string) (requester fosite.Requester, err error) {
	t, err := h.validate(ctx, token)
	if err != nil {
		return nil, err
	}

	claims := jwt.JWTClaims{}
	claims.FromMapClaims(t.Claims.(jwtx.MapClaims))

	requester = &fosite.Request{
		Client:      &fosite.DefaultClient{},
		RequestedAt: claims.IssuedAt,
		Session: &JWTSession{
			JWTClaims: &claims,
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType: claims.ExpiresAt,
			},
			Subject: claims.Subject,
		},
		Scopes:        claims.Scope,
		GrantedScopes: claims.Scope,
	}

	return
}

func (h DefaultJWTStrategy) RefreshTokenSignature(token string) string {
	return h.HMACSHAStrategy.RefreshTokenSignature(token)
}

func (h DefaultJWTStrategy) AuthorizeCodeSignature(token string) string {
	return h.HMACSHAStrategy.AuthorizeCodeSignature(token)
}

func (h *DefaultJWTStrategy) GenerateRefreshToken(ctx context.Context, req fosite.Requester) (token string, signature string, err error) {
	return h.HMACSHAStrategy.GenerateRefreshToken(ctx, req)
}

func (h *DefaultJWTStrategy) ValidateRefreshToken(ctx context.Context, req fosite.Requester, token string) error {
	return h.HMACSHAStrategy.ValidateRefreshToken(ctx, req, token)
}

func (h *DefaultJWTStrategy) GenerateAuthorizeCode(ctx context.Context, req fosite.Requester) (token string, signature string, err error) {
	return h.HMACSHAStrategy.GenerateAuthorizeCode(ctx, req)
}

func (h *DefaultJWTStrategy) ValidateAuthorizeCode(ctx context.Context, req fosite.Requester, token string) error {
	return h.HMACSHAStrategy.ValidateAuthorizeCode(ctx, req, token)
}

func (h *DefaultJWTStrategy) validate(ctx context.Context, token string) (t *jwtx.Token, err error) {
	t, err = h.JWTStrategy.Decode(ctx, token)

	if err == nil {
		err = t.Claims.Valid()
	}

	if err != nil {
		if e, ok := errors.Cause(err).(*jwtx.ValidationError); ok {
			switch e.Errors {
			case jwtx.ValidationErrorMalformed:
				err = errors.WithStack(fosite.ErrInvalidTokenFormat.WithDebug(err.Error()))
			case jwtx.ValidationErrorUnverifiable:
				err = errors.WithStack(fosite.ErrTokenSignatureMismatch.WithDebug(err.Error()))
			case jwtx.ValidationErrorSignatureInvalid:
				err = errors.WithStack(fosite.ErrTokenSignatureMismatch.WithDebug(err.Error()))
			case jwtx.ValidationErrorAudience:
				err = errors.WithStack(fosite.ErrTokenClaim.WithDebug(err.Error()))
			case jwtx.ValidationErrorExpired:
				err = errors.WithStack(fosite.ErrTokenExpired.WithDebug(err.Error()))
			case jwtx.ValidationErrorIssuedAt:
				err = errors.WithStack(fosite.ErrTokenClaim.WithDebug(err.Error()))
			case jwtx.ValidationErrorIssuer:
				err = errors.WithStack(fosite.ErrTokenClaim.WithDebug(err.Error()))
			case jwtx.ValidationErrorNotValidYet:
				err = errors.WithStack(fosite.ErrTokenClaim.WithDebug(err.Error()))
			case jwtx.ValidationErrorId:
				err = errors.WithStack(fosite.ErrTokenClaim.WithDebug(err.Error()))
			case jwtx.ValidationErrorClaimsInvalid:
				err = errors.WithStack(fosite.ErrTokenClaim.WithDebug(err.Error()))
			default:
				err = errors.WithStack(fosite.ErrRequestUnauthorized.WithDebug(err.Error()))
			}
		}
	}

	return
}

func (h *DefaultJWTStrategy) generate(ctx context.Context, tokenType fosite.TokenType, requester fosite.Requester) (string, string, error) {
	if jwtSession, ok := requester.GetSession().(JWTSessionContainer); !ok {
		return "", "", errors.New("Session must be of type JWTSessionContainer")
	} else if jwtSession.GetJWTClaims() == nil {
		return "", "", errors.New("GetTokenClaims() must not be nil")
	} else {
		claims := jwtSession.GetJWTClaims()
		claims.ExpiresAt = jwtSession.GetExpiresAt(tokenType)

		if claims.IssuedAt.IsZero() {
			claims.IssuedAt = time.Now().UTC()
		}

		if claims.Issuer == "" {
			claims.Issuer = h.Issuer
		}

		claims.Scope = requester.GetGrantedScopes()

		return h.JWTStrategy.Generate(ctx, claims.ToMapClaims(), jwtSession.GetJWTHeader())
	}
}
