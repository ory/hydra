// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/x/errorsx"
)

// DefaultJWTStrategy is a JWT RS256 strategy.
type DefaultJWTStrategy struct {
	jwt.Signer
	Strategy interface {
		AuthorizeCodeStrategyProvider
		AccessTokenStrategyProvider
		RefreshTokenStrategyProvider
	}
	Config interface {
		fosite.AccessTokenIssuerProvider
		fosite.JWTScopeFieldProvider
	}
}

func (h DefaultJWTStrategy) signature(token string) string {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return ""
	}

	return split[2]
}

func (h *DefaultJWTStrategy) AuthorizeCodeStrategy() AuthorizeCodeStrategy {
	return h
}

func (h *DefaultJWTStrategy) AccessTokenStrategy() AccessTokenStrategy {
	return h
}

func (h *DefaultJWTStrategy) RefreshTokenStrategy() RefreshTokenStrategy {
	return h
}

func (h DefaultJWTStrategy) AccessTokenSignature(ctx context.Context, token string) string {
	return h.signature(token)
}

func (h *DefaultJWTStrategy) GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(ctx, fosite.AccessToken, requester)
}

func (h *DefaultJWTStrategy) ValidateAccessToken(ctx context.Context, _ fosite.Requester, token string) error {
	_, err := validate(ctx, h.Signer, token)
	return err
}

func (h DefaultJWTStrategy) RefreshTokenSignature(ctx context.Context, token string) string {
	return h.Strategy.RefreshTokenStrategy().RefreshTokenSignature(ctx, token)
}

func (h DefaultJWTStrategy) AuthorizeCodeSignature(ctx context.Context, token string) string {
	return h.Strategy.AuthorizeCodeStrategy().AuthorizeCodeSignature(ctx, token)
}

func (h *DefaultJWTStrategy) GenerateRefreshToken(ctx context.Context, req fosite.Requester) (token string, signature string, err error) {
	return h.Strategy.RefreshTokenStrategy().GenerateRefreshToken(ctx, req)
}

func (h *DefaultJWTStrategy) ValidateRefreshToken(ctx context.Context, req fosite.Requester, token string) error {
	return h.Strategy.RefreshTokenStrategy().ValidateRefreshToken(ctx, req, token)
}

func (h *DefaultJWTStrategy) GenerateAuthorizeCode(ctx context.Context, req fosite.Requester) (token string, signature string, err error) {
	return h.Strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(ctx, req)
}

func (h *DefaultJWTStrategy) ValidateAuthorizeCode(ctx context.Context, req fosite.Requester, token string) error {
	return h.Strategy.AuthorizeCodeStrategy().ValidateAuthorizeCode(ctx, req, token)
}

func validate(ctx context.Context, jwtStrategy jwt.Signer, token string) (t *jwt.Token, err error) {
	t, err = jwtStrategy.Decode(ctx, token)
	if err == nil {
		err = t.Claims.Valid()
		return
	}

	var e *jwt.ValidationError
	if errors.As(err, &e) {
		err = errorsx.WithStack(toRFCErr(e).WithWrap(err).WithDebug(err.Error()))
	}

	return
}

func toRFCErr(v *jwt.ValidationError) *fosite.RFC6749Error {
	switch {
	case v == nil:
		return nil
	case v.Has(jwt.ValidationErrorMalformed):
		return fosite.ErrInvalidTokenFormat
	case v.Has(jwt.ValidationErrorUnverifiable | jwt.ValidationErrorSignatureInvalid):
		return fosite.ErrTokenSignatureMismatch
	case v.Has(jwt.ValidationErrorExpired):
		return fosite.ErrTokenExpired
	case v.Has(jwt.ValidationErrorAudience |
		jwt.ValidationErrorIssuedAt |
		jwt.ValidationErrorIssuer |
		jwt.ValidationErrorNotValidYet |
		jwt.ValidationErrorId |
		jwt.ValidationErrorClaimsInvalid):
		return fosite.ErrTokenClaim
	default:
		return fosite.ErrRequestUnauthorized
	}
}

func (h *DefaultJWTStrategy) generate(ctx context.Context, tokenType fosite.TokenType, requester fosite.Requester) (string, string, error) {
	if jwtSession, ok := requester.GetSession().(JWTSessionContainer); !ok {
		return "", "", errors.Errorf("Session must be of type JWTSessionContainer but got type: %T", requester.GetSession())
	} else if claims := jwtSession.GetJWTClaims(); claims == nil {
		return "", "", errors.New("GetTokenClaims() must not be nil")
	} else {
		claims.
			With(
				jwtSession.GetExpiresAt(tokenType),
				requester.GetGrantedScopes(),
				requester.GetGrantedAudience(),
			).
			WithDefaults(
				time.Now().UTC(),
				h.Config.GetAccessTokenIssuer(ctx),
			).
			WithScopeField(
				h.Config.GetJWTScopeField(ctx),
			)

		return h.Signer.Generate(ctx, claims.ToMapClaims(), jwtSession.GetJWTHeader())
	}
}
