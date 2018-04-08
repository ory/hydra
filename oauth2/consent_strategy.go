// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	ejwt "github.com/ory/fosite/token/jwt"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

const (
	CookieCSRFKey = "consent_csrf"
)

type DefaultConsentStrategy struct {
	Issuer string

	DefaultIDTokenLifespan   time.Duration
	DefaultChallengeLifespan time.Duration
	ConsentManager           ConsentRequestManager
}

func (s *DefaultConsentStrategy) validateSession(req fosite.AuthorizeRequester, consent *ConsentRequest, cookie *sessions.Session) error {
	if j, ok := cookie.Values[CookieCSRFKey]; !ok {
		return errors.Errorf("Session cookie is missing CSRF token")
	} else if js, ok := j.(string); !ok {
		return errors.Errorf("CSRF value in session cookie is not a string")
	} else if js != consent.CSRF {
		return errors.Errorf("CSRF value in session cookie does not match consent CSRF value")
	} else if consent.CSRF != req.GetRequestForm().Get("consent_csrf") {
		return errors.Errorf("CSRF value from query parameters does not match consent CSRF value")
	}

	if time.Now().UTC().After(consent.ExpiresAt) {
		return errors.Errorf("Consent session expired")
	}

	if consent.ClientID != req.GetClient().GetID() {
		return errors.Errorf("ClientID mismatch")
	}

	if consent.Subject == "" {
		return errors.Errorf("Subject key is empty or undefined in consent response, check your payload.")
	}

	return nil
}

func (s *DefaultConsentStrategy) ValidateConsentRequest(req fosite.AuthorizeRequester, session string, cookie *sessions.Session) (claims *Session, err error) {
	defer delete(cookie.Values, CookieCSRFKey)

	consent, err := s.ConsentManager.GetConsentRequest(session)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !consent.IsConsentGranted() {
		err := errors.New("The resource owner denied consent for this request")
		return nil, &fosite.RFC6749Error{
			Name:        "rejected_consent_request",
			Description: consent.DenyReason,
			Debug:       err.Error(),
			Hint:        consent.DenyReason,
			Code:        http.StatusUnauthorized,
		}
	}

	if err := s.validateSession(req, consent, cookie); err != nil {
		if err := s.ConsentManager.RejectConsentRequest(session, &RejectConsentRequestPayload{
			Reason: "Unable to validate consent request",
		}); err != nil {
			return nil, err
		}
		return nil, err
	}

	for _, scope := range consent.GrantedScopes {
		req.GrantScope(scope)
	}

	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims: &ejwt.IDTokenClaims{
				Audience:    req.GetClient().GetID(),
				Subject:     consent.Subject,
				Issuer:      s.Issuer,
				IssuedAt:    time.Now().UTC(),
				ExpiresAt:   time.Now().UTC().Add(s.DefaultIDTokenLifespan),
				AuthTime:    time.Now().UTC(),
				RequestedAt: time.Now().UTC(),
				Extra:       consent.IDTokenExtra,
			},
			// required for lookup on jwk endpoint
			Headers: &ejwt.Headers{Extra: map[string]interface{}{"kid": "public"}},
			Subject: consent.Subject,
		},
		Extra: consent.AccessTokenExtra,
	}, err
}

func (s *DefaultConsentStrategy) CreateConsentRequest(req fosite.AuthorizeRequester, redirectURL string, cookie *sessions.Session) (string, error) {
	csrf := uuid.New()
	id := uuid.New()

	cookie.Values[CookieCSRFKey] = csrf
	consent := &ConsentRequest{
		ID:               id,
		CSRF:             csrf,
		GrantedScopes:    []string{},
		RequestedScopes:  req.GetRequestedScopes(),
		ClientID:         req.GetClient().GetID(),
		ExpiresAt:        time.Now().UTC().Add(s.DefaultChallengeLifespan),
		RedirectURL:      redirectURL + "&consent=" + id + "&consent_csrf=" + csrf,
		AccessTokenExtra: map[string]interface{}{},
		IDTokenExtra:     map[string]interface{}{},
	}

	if err := s.ConsentManager.PersistConsentRequest(consent); err != nil {
		return "", errors.WithStack(err)
	}

	return id, nil
}
