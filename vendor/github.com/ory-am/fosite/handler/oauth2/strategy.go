package oauth2

import (
	"github.com/ory-am/fosite"
	"golang.org/x/net/context"
)

type CoreStrategy interface {
	AccessTokenStrategy
	RefreshTokenStrategy
	AuthorizeCodeStrategy
}

type JWTStrategy interface {
	ValidateJWT(tokenType fosite.TokenType, token string) (requester fosite.Requester, err error)
}

type AccessTokenStrategy interface {
	AccessTokenSignature(token string) string
	GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error)
	ValidateAccessToken(ctx context.Context, requester fosite.Requester, token string) (err error)
}

type RefreshTokenStrategy interface {
	RefreshTokenSignature(token string) string
	GenerateRefreshToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error)
	ValidateRefreshToken(ctx context.Context, requester fosite.Requester, token string) (err error)
}

type AuthorizeCodeStrategy interface {
	AuthorizeCodeSignature(token string) string
	GenerateAuthorizeCode(ctx context.Context, requester fosite.Requester) (token string, signature string, err error)
	ValidateAuthorizeCode(ctx context.Context, requester fosite.Requester, token string) (err error)
}
