package oauth2

import (
	"time"

	"fmt"
	"github.com/ory-am/fosite"
	enigma "github.com/ory-am/fosite/token/hmac"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type HMACSHAStrategy struct {
	Enigma                *enigma.HMACStrategy
	AccessTokenLifespan   time.Duration
	AuthorizeCodeLifespan time.Duration
}

func (h HMACSHAStrategy) AccessTokenSignature(token string) string {
	return h.Enigma.Signature(token)
}
func (h HMACSHAStrategy) RefreshTokenSignature(token string) string {
	return h.Enigma.Signature(token)
}
func (h HMACSHAStrategy) AuthorizeCodeSignature(token string) string {
	return h.Enigma.Signature(token)
}

func (h HMACSHAStrategy) GenerateAccessToken(_ context.Context, _ fosite.Requester) (token string, signature string, err error) {
	return h.Enigma.Generate()
}

func (h HMACSHAStrategy) ValidateAccessToken(_ context.Context, r fosite.Requester, token string) (err error) {
	var exp = r.GetSession().GetExpiresAt(fosite.AccessToken)
	if exp.IsZero() && r.GetRequestedAt().Add(h.AccessTokenLifespan).Before(time.Now()) {
		return errors.Wrap(fosite.ErrTokenExpired, fmt.Sprintf("Access token expired at %s", r.GetRequestedAt().Add(h.AccessTokenLifespan)))
	}
	if !exp.IsZero() && exp.Before(time.Now()) {
		return errors.Wrap(fosite.ErrTokenExpired, fmt.Sprintf("Access token expired at %s", exp))
	}
	return h.Enigma.Validate(token)
}

func (h HMACSHAStrategy) GenerateRefreshToken(_ context.Context, _ fosite.Requester) (token string, signature string, err error) {
	return h.Enigma.Generate()
}

func (h HMACSHAStrategy) ValidateRefreshToken(_ context.Context, _ fosite.Requester, token string) (err error) {
	return h.Enigma.Validate(token)
}

func (h HMACSHAStrategy) GenerateAuthorizeCode(_ context.Context, _ fosite.Requester) (token string, signature string, err error) {
	return h.Enigma.Generate()
}

func (h HMACSHAStrategy) ValidateAuthorizeCode(_ context.Context, r fosite.Requester, token string) (err error) {
	var exp = r.GetSession().GetExpiresAt(fosite.AuthorizeCode)
	if exp.IsZero() && r.GetRequestedAt().Add(h.AuthorizeCodeLifespan).Before(time.Now()) {
		return errors.Wrap(fosite.ErrTokenExpired, fmt.Sprintf("Authorize code expired at %s", r.GetRequestedAt().Add(h.AuthorizeCodeLifespan)))
	}
	if !exp.IsZero() && exp.Before(time.Now()) {
		return errors.Wrap(fosite.ErrTokenExpired, fmt.Sprintf("Authorize code expired at %s", exp))
	}

	return h.Enigma.Validate(token)
}
