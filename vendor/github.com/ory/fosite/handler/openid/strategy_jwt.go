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

package openid

import (
	"context"
	"fmt"
	"strconv"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/go-convenience/stringslice"
)

const defaultExpiryTime = time.Hour

type Session interface {
	IDTokenClaims() *jwt.IDTokenClaims
	IDTokenHeaders() *jwt.Headers

	fosite.Session
}

// IDTokenSession is a session container for the id token
type DefaultSession struct {
	Claims    *jwt.IDTokenClaims
	Headers   *jwt.Headers
	ExpiresAt map[fosite.TokenType]time.Time
	Username  string
	Subject   string
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
	jwt.JWTStrategy

	Expiry time.Duration
	Issuer string
}

func (h DefaultStrategy) GenerateIDToken(ctx context.Context, requester fosite.Requester) (token string, err error) {
	if h.Expiry == 0 {
		h.Expiry = defaultExpiryTime
	}

	sess, ok := requester.GetSession().(Session)
	if !ok {
		return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because session must be of type fosite/handler/openid.Session."))
	}

	claims := sess.IDTokenClaims()
	if claims.Subject == "" {
		return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because subject is an empty string."))
	}

	if requester.GetRequestForm().Get("grant_type") != "refresh_token" {
		maxAge, err := strconv.ParseInt(requester.GetRequestForm().Get("max_age"), 10, 64)
		if err != nil {
			maxAge = 0
		}

		// Adds a bit of wiggle room for timing issues
		if claims.AuthTime.After(time.Now().UTC().Add(time.Second * 5)) {
			return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to validate OpenID Connect request because authentication time is in the future."))
		}

		if maxAge > 0 {
			if claims.AuthTime.IsZero() {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because authentication time claim is required when max_age is set."))
			} else if claims.RequestedAt.IsZero() {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because requested at claim is required when max_age is set."))
			} else if claims.AuthTime.Add(time.Second * time.Duration(maxAge)).Before(claims.RequestedAt) {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because authentication time does not satisfy max_age time."))
			}
		}

		prompt := requester.GetRequestForm().Get("prompt")
		if prompt != "" {
			if claims.AuthTime.IsZero() {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Unable to determine validity of prompt parameter because auth_time is missing in id token claims."))
			}
		}

		switch prompt {
		case "none":
			if claims.AuthTime.After(claims.RequestedAt) {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because prompt was set to \"none\" but auth_time happened after the authorization request was registered, indicating that the user was logged in during this request which is not allowed."))
			}
			break
		case "login":
			if claims.AuthTime.Before(claims.RequestedAt) {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because prompt was set to \"login\" but auth_time happened before the authorization request was registered, indicating that the user was not re-authenticated which is forbidden."))
			}
			break
		}

		// If acr_values was requested but no acr value was provided in the ID token, fall back to level 0 which means least
		// confidence in authentication.
		if requester.GetRequestForm().Get("acr_values") != "" && claims.AuthenticationContextClassReference == "" {
			claims.AuthenticationContextClassReference = "0"
		}

		if tokenHintString := requester.GetRequestForm().Get("id_token_hint"); tokenHintString != "" {
			tokenHint, err := h.JWTStrategy.Decode(ctx, tokenHintString)
			if ve, ok := errors.Cause(err).(*jwtgo.ValidationError); ok && ve.Errors == jwtgo.ValidationErrorExpired {
				// Expired ID Tokens are allowed as values to id_token_hint
			} else if err != nil {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug(fmt.Sprintf("Unable to decode id token from id_token_hint parameter because %s.", err.Error())))
			}

			if hintClaims, ok := tokenHint.Claims.(jwtgo.MapClaims); !ok {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Unable to decode id token from id_token_hint to *jwt.StandardClaims."))
			} else if hintSub, _ := hintClaims["sub"].(string); hintSub == "" {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug("Provided id token from id_token_hint does not have a subject."))
			} else if hintSub != claims.Subject {
				return "", errors.WithStack(fosite.ErrServerError.WithDebug(fmt.Sprintf("Subject from authorization mismatches id token subject from id_token_hint.")))
			}
		}
	}

	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = time.Now().UTC().Add(h.Expiry)
	}

	if claims.ExpiresAt.Before(time.Now().UTC()) {
		return "", errors.WithStack(fosite.ErrServerError.WithDebug("Failed to generate id token because expiry claim can not be in the past."))
	}

	if claims.AuthTime.IsZero() {
		claims.AuthTime = time.Now().UTC()
	}

	if claims.Issuer == "" {
		claims.Issuer = h.Issuer
	}

	nonce := requester.GetRequestForm().Get("nonce")
	// OPTIONAL. String value used to associate a Client session with an ID Token, and to mitigate replay attacks.
	if len(nonce) == 0 {
	} else if len(nonce) > 0 && len(nonce) < fosite.MinParameterEntropy {
		// We're assuming that using less then 8 characters for the state can not be considered "unguessable"
		return "", errors.WithStack(fosite.ErrInsufficientEntropy.WithHintf("Parameter \"nonce\" is set but does not satisfy the minimum entropy of %d characters.", fosite.MinParameterEntropy))
	}

	claims.Nonce = nonce
	claims.Audience = stringslice.Unique(append(claims.Audience, requester.GetClient().GetID()))
	claims.IssuedAt = time.Now().UTC()

	token, _, err = h.JWTStrategy.Generate(ctx, claims.ToMapClaims(), sess.IDTokenHeaders())
	return token, err
}
