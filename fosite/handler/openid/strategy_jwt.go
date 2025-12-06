// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"
	"strconv"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"

	"github.com/ory/go-convenience/stringslice"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

const defaultExpiryTime = time.Hour

type Session interface {
	// IDTokenClaims returns a pointer to claims which will be modified in-place by handlers.
	// Session should store this pointer and return always the same pointer.
	IDTokenClaims() *jwt.IDTokenClaims
	// IDTokenHeaders returns a pointer to header values which will be modified in-place by handlers.
	// Session should store this pointer and return always the same pointer.
	IDTokenHeaders() *jwt.Headers

	fosite.Session
}

// IDTokenSession is a session container for the id token
type DefaultSession struct {
	Claims    *jwt.IDTokenClaims             `json:"id_token_claims"`
	Headers   *jwt.Headers                   `json:"headers"`
	ExpiresAt map[fosite.TokenType]time.Time `json:"expires_at"`
	Username  string                         `json:"username,omitempty"`
	Subject   string                         `json:"subject,omitempty"`
}

func NewDefaultSession() *DefaultSession {
	return &DefaultSession{
		Claims: &jwt.IDTokenClaims{
			RequestedAt: time.Now().UTC(),
		},
		Headers: &jwt.Headers{},
	}
}

func (s *DefaultSession) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	return deepcopy.Copy(s).(fosite.Session)
}

func (s *DefaultSession) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}
	s.ExpiresAt[key] = exp
}

func (s *DefaultSession) GetExpiresAt(key fosite.TokenType) time.Time {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}

	if _, ok := s.ExpiresAt[key]; !ok {
		return time.Time{}
	}
	return s.ExpiresAt[key]
}

func (s *DefaultSession) GetUsername() string {
	if s == nil {
		return ""
	}
	return s.Username
}

func (s *DefaultSession) SetSubject(subject string) {
	s.Subject = subject
}

func (s *DefaultSession) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

func (s *DefaultSession) IDTokenHeaders() *jwt.Headers {
	if s.Headers == nil {
		s.Headers = &jwt.Headers{}
	}
	return s.Headers
}

func (s *DefaultSession) IDTokenClaims() *jwt.IDTokenClaims {
	if s.Claims == nil {
		s.Claims = &jwt.IDTokenClaims{}
	}
	return s.Claims
}

type DefaultStrategy struct {
	jwt.Signer

	Config interface {
		fosite.IDTokenIssuerProvider
		fosite.IDTokenLifespanProvider
		fosite.MinParameterEntropyProvider
	}
}

// GenerateIDToken returns a JWT string.
//
// lifespan is ignored if requester.GetSession().IDTokenClaims().ExpiresAt is not zero.
func (h DefaultStrategy) GenerateIDToken(ctx context.Context, lifespan time.Duration, requester fosite.Requester) (token string, err error) {
	if lifespan == 0 {
		lifespan = defaultExpiryTime
	}

	sess, ok := requester.GetSession().(Session)
	if !ok {
		return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because session must be of type fosite/handler/openid.Session."))
	}

	claims := sess.IDTokenClaims()
	if claims.Subject == "" {
		return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because subject is an empty string."))
	}

	if requester.GetRequestForm().Get("grant_type") != "refresh_token" {
		maxAge, err := strconv.ParseInt(requester.GetRequestForm().Get("max_age"), 10, 64)
		if err != nil {
			maxAge = 0
		}

		// Adds a bit of wiggle room for timing issues
		if claims.AuthTime.After(time.Now().UTC().Add(time.Second * 5)) {
			return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to validate OpenID Connect request because authentication time is in the future."))
		}

		if maxAge > 0 {
			if claims.AuthTime.IsZero() {
				return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because authentication time claim is required when max_age is set."))
			} else if claims.RequestedAt.IsZero() {
				return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because requested at claim is required when max_age is set."))
			} else if claims.AuthTime.Add(time.Second * time.Duration(maxAge)).Before(claims.RequestedAt) {
				return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because authentication time does not satisfy max_age time."))
			}
		}

		prompt := requester.GetRequestForm().Get("prompt")
		if prompt != "" {
			if claims.AuthTime.IsZero() {
				return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Unable to determine validity of prompt parameter because auth_time is missing in id token claims."))
			}
		}

		switch prompt {
		case "none":
			if !claims.AuthTime.Equal(claims.RequestedAt) && claims.AuthTime.After(claims.RequestedAt) {
				return "", errorsx.WithStack(fosite.ErrServerError.
					WithDebugf("Failed to generate id token because prompt was set to 'none' but auth_time ('%s') happened after the authorization request ('%s') was registered, indicating that the user was logged in during this request which is not allowed.", claims.AuthTime, claims.RequestedAt))
			}
		case "login":
			if !claims.AuthTime.Equal(claims.RequestedAt) && claims.AuthTime.Before(claims.RequestedAt) {
				return "", errorsx.WithStack(fosite.ErrServerError.
					WithDebugf("Failed to generate id token because prompt was set to 'login' but auth_time ('%s') happened before the authorization request ('%s') was registered, indicating that the user was not re-authenticated which is forbidden.", claims.AuthTime, claims.RequestedAt))
			}
		}

		// If acr_values was requested but no acr value was provided in the ID token, fall back to level 0 which means least
		// confidence in authentication.
		if requester.GetRequestForm().Get("acr_values") != "" && claims.AuthenticationContextClassReference == "" {
			claims.AuthenticationContextClassReference = "0"
		}

		if tokenHintString := requester.GetRequestForm().Get("id_token_hint"); tokenHintString != "" {
			tokenHint, err := h.Signer.Decode(ctx, tokenHintString)
			var ve *jwt.ValidationError
			if errors.As(err, &ve) && ve.Has(jwt.ValidationErrorExpired) {
				// Expired ID Tokens are allowed as values to id_token_hint
			} else if err != nil {
				return "", errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithDebugf("Unable to decode id token from 'id_token_hint' parameter because %s.", err.Error()))
			}

			if hintSub, _ := tokenHint.Claims["sub"].(string); hintSub == "" {
				return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Provided id token from 'id_token_hint' does not have a subject."))
			} else if hintSub != claims.Subject {
				return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Subject from authorization mismatches id token subject from 'id_token_hint'."))
			}
		}
	}

	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = time.Now().UTC().Add(lifespan)
	}

	if claims.ExpiresAt.Before(time.Now().UTC()) {
		return "", errorsx.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because expiry claim can not be in the past."))
	}

	if claims.AuthTime.IsZero() {
		claims.AuthTime = time.Now().Truncate(time.Second).UTC()
	}

	if claims.Issuer == "" {
		claims.Issuer = h.Config.GetIDTokenIssuer(ctx)
	}

	// OPTIONAL. String value used to associate a Client session with an ID Token, and to mitigate replay attacks.
	if nonce := requester.GetRequestForm().Get("nonce"); len(nonce) == 0 {
	} else if len(nonce) > 0 && len(nonce) < h.Config.GetMinParameterEntropy(ctx) {
		// We're assuming that using less then, by default, 8 characters for the state can not be considered "unguessable"
		return "", errorsx.WithStack(fosite.ErrInsufficientEntropy.WithHintf("Parameter 'nonce' is set but does not satisfy the minimum entropy of %d characters.", h.Config.GetMinParameterEntropy(ctx)))
	} else if len(nonce) > 0 {
		claims.Nonce = nonce
	}

	claims.Audience = stringslice.Unique(append(claims.Audience, requester.GetClient().GetID()))
	claims.IssuedAt = time.Now().UTC()

	token, _, err = h.Signer.Generate(ctx, claims.ToMapClaims(), sess.IDTokenHeaders())
	return token, err
}
