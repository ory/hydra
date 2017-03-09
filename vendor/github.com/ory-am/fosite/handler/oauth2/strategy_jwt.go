package oauth2

import (
	"strings"
	"time"

	jwtx "github.com/dgrijalva/jwt-go"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// RS256JWTStrategy is a JWT RS256 strategy.
type RS256JWTStrategy struct {
	*jwt.RS256JWTStrategy
	Issuer string
}

func (h RS256JWTStrategy) signature(token string) string {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return ""
	}

	return split[2]
}

func (h RS256JWTStrategy) AccessTokenSignature(token string) string {
	return h.signature(token)
}

func (h RS256JWTStrategy) RefreshTokenSignature(token string) string {
	return h.signature(token)
}

func (h RS256JWTStrategy) AuthorizeCodeSignature(token string) string {
	return h.signature(token)
}

func (h *RS256JWTStrategy) ValidateJWT(tokenType fosite.TokenType, token string) (requester fosite.Requester, err error) {
	t, err := h.validate(token)
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

func (h *RS256JWTStrategy) GenerateAccessToken(_ context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(fosite.AccessToken, requester)
}

func (h *RS256JWTStrategy) ValidateAccessToken(_ context.Context, _ fosite.Requester, token string) error {
	_, err := h.validate(token)
	return err
}

func (h *RS256JWTStrategy) GenerateRefreshToken(_ context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(fosite.RefreshToken, requester)
}

func (h *RS256JWTStrategy) ValidateRefreshToken(_ context.Context, _ fosite.Requester, token string) error {
	_, err := h.validate(token)
	return err
}

func (h *RS256JWTStrategy) GenerateAuthorizeCode(_ context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(fosite.AuthorizeCode, requester)
}

func (h *RS256JWTStrategy) ValidateAuthorizeCode(_ context.Context, requester fosite.Requester, token string) error {
	_, err := h.validate(token)
	return err
}

func (h *RS256JWTStrategy) validate(token string) (t *jwtx.Token, err error) {
	t, err = h.RS256JWTStrategy.Decode(token)

	if err == nil {
		err = t.Claims.Valid()
	}

	if err != nil {
		if e, ok := errors.Cause(err).(*jwtx.ValidationError); ok {
			switch e.Errors {
			case jwtx.ValidationErrorMalformed:
				err = errors.Wrap(fosite.ErrInvalidTokenFormat, err.Error())
			case jwtx.ValidationErrorUnverifiable:
				err = errors.Wrap(fosite.ErrTokenSignatureMismatch, err.Error())
			case jwtx.ValidationErrorSignatureInvalid:
				err = errors.Wrap(fosite.ErrTokenSignatureMismatch, err.Error())
			case jwtx.ValidationErrorAudience:
				err = errors.Wrap(fosite.ErrTokenClaim, err.Error())
			case jwtx.ValidationErrorExpired:
				err = errors.Wrap(fosite.ErrTokenExpired, err.Error())
			case jwtx.ValidationErrorIssuedAt:
				err = errors.Wrap(fosite.ErrTokenClaim, err.Error())
			case jwtx.ValidationErrorIssuer:
				err = errors.Wrap(fosite.ErrTokenClaim, err.Error())
			case jwtx.ValidationErrorNotValidYet:
				err = errors.Wrap(fosite.ErrTokenClaim, err.Error())
			case jwtx.ValidationErrorId:
				err = errors.Wrap(fosite.ErrTokenClaim, err.Error())
			case jwtx.ValidationErrorClaimsInvalid:
				err = errors.Wrap(fosite.ErrTokenClaim, err.Error())
			default:
				err = errors.Wrap(fosite.ErrRequestUnauthorized, err.Error())
			}
		}
	}

	return
}

func (h *RS256JWTStrategy) generate(tokenType fosite.TokenType, requester fosite.Requester) (string, string, error) {
	if jwtSession, ok := requester.GetSession().(JWTSessionContainer); !ok {
		return "", "", errors.New("Session must be of type JWTSessionContainer")
	} else if jwtSession.GetJWTClaims() == nil {
		return "", "", errors.New("GetTokenClaims() must not be nil")
	} else {
		claims := jwtSession.GetJWTClaims()
		claims.ExpiresAt = jwtSession.GetExpiresAt(tokenType)

		if claims.IssuedAt.IsZero() {
			claims.IssuedAt = time.Now()
		}

		if claims.Issuer == "" {
			claims.Issuer = h.Issuer
		}

		claims.Scope = requester.GetGrantedScopes()

		return h.RS256JWTStrategy.Generate(claims.ToMapClaims(), jwtSession.GetJWTHeader())
	}
}
