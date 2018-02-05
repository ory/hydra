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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	ejwt "github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

const (
	cookieCSRFKey      = "consent_csrf"
	sessionUserKey     = "consent_session_user"
	sessionAuthTimeKey = "consent_session_auth_time"
	jwtIdTokenKey      = "public"
)

type DefaultConsentStrategy struct {
	Issuer string

	KeyID                    string
	DefaultIDTokenLifespan   time.Duration
	DefaultChallengeLifespan time.Duration
	ConsentManager           ConsentRequestManager
}

var errRequiresAuthentication = errors.New("Requires authentication")

func checkAntiReplayToken(consent *ConsentRequest, cookie *sessions.Session) error {
	if j, ok := cookie.Values[cookieCSRFKey]; !ok {
		return errors.Errorf("Session cookie is missing anti-replay token")
	} else if js, ok := j.(string); !ok {
		return errors.Errorf("Session cookie anti-replay value is not a string")
	} else if js != consent.CSRF {
		return errors.Errorf("Session cookie anti-replay value does not match value from consent response")
	}
	return nil
}

func (s *DefaultConsentStrategy) ValidateConsentRequest(req fosite.AuthorizeRequester, session string, cookie *sessions.Session) (claims *Session, err error) {
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

	if time.Now().After(consent.ExpiresAt) {
		return nil, errors.Errorf("Token expired")
	}

	if consent.ClientID != req.GetClient().GetID() {
		return nil, errors.Errorf("ClientID mismatch")
	}

	if consent.Subject == "" {
		return nil, errors.Errorf("Subject key is empty or undefined in consent response, check your payload.")
	}

	if err := checkAntiReplayToken(consent, cookie); err != nil {
		if err := s.ConsentManager.RejectConsentRequest(session, &RejectConsentRequestPayload{
			Reason: "Session cookie is missing anti-replay token",
		}); err != nil {
			return nil, err
		}
		return nil, err
	}

	for _, scope := range consent.GrantedScopes {
		req.GrantScope(scope)
	}

	now := time.Now().UTC()

	delete(cookie.Values, cookieCSRFKey)
	cookie.Values[sessionUserKey] = consent.Subject
	cookie.Values[sessionAuthTimeKey] = now.Format(time.RFC3339)

	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims: &ejwt.IDTokenClaims{
				Audience:    req.GetClient().GetID(),
				Subject:     consent.Subject,
				Issuer:      s.Issuer,
				IssuedAt:    now,
				ExpiresAt:   now.Add(s.DefaultIDTokenLifespan),
				AuthTime:    now,
				RequestedAt: consent.RequestedAt,
				Extra:       consent.IDTokenExtra,
			},
			// required for lookup on jwk endpoint
			Headers: &ejwt.Headers{Extra: map[string]interface{}{"kid": s.KeyID}},
			Subject: consent.Subject,
		},
		Extra: consent.AccessTokenExtra,
	}, err
}

func userFromCookie(cookie *sessions.Session) (string, bool) {
	if cookieValue, ok := cookie.Values[sessionUserKey]; !ok {
		return "", false
	} else if user, ok := cookieValue.(string); !ok {
		return "", false
	} else if len(user) == 0 {
		return "", false
	} else {
		return user, true
	}
}

func authTimeFromCookie(cookie *sessions.Session) (time.Time, bool) {
	if cookieValue, ok := cookie.Values[sessionAuthTimeKey]; !ok {
		return time.Time{}, false
	} else if value, ok := cookieValue.(string); !ok {
		return time.Time{}, false
	} else if authTime, err := time.Parse(time.RFC3339, value); err != nil {
		return time.Time{}, false
	} else {
		return authTime, true
	}
}

func contains(s []string, n string) bool {
	for _, v := range s {
		if v == n {
			return true
		}
	}
	return false
}

func containsWhiteListedOnly(items []string, whiteList []string) bool {
	for _, item := range items {
		if !contains(whiteList, item) {
			return false
		}
	}
	return true
}

func (s *DefaultConsentStrategy) HandleConsentRequest(req fosite.AuthorizeRequester, cookie *sessions.Session) (claims *Session, err error) {
	prompt := req.GetRequestForm().Get("prompt")
	prompts := pkg.SplitNonEmpty(prompt, " ")
	maxAge, _ := strconv.ParseInt(req.GetRequestForm().Get("max_age"), 10, 64)

	if !containsWhiteListedOnly(prompts, []string{"login", "none", "consent", "select_account"}) {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug(fmt.Sprintf(`Used unknown value "%s" for prompt parameter`, prompt)))
	}

	if contains(prompts, "none") && len(prompts) > 1 {
		// If this parameter contains none with any other value, an error is returned.
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Parameter prompt was set to none, but contains other values as well which is not allowed"))
	}

	if contains(prompts, "login") || contains(prompts, "consent") || contains(prompts, "select_account") {
		return nil, errors.WithStack(errRequiresAuthentication)
	}

	if req.GetClient().IsPublic() {
		return nil, errors.WithStack(errRequiresAuthentication)
	}

	var user string
	var ok bool
	var authTime time.Time

	if user, ok = userFromCookie(cookie); !ok {
		if contains(prompts, "none") {
			return nil, errors.WithStack(fosite.ErrLoginRequired.WithDebug("Session not set but prompt is set to none"))
		}
		return nil, errors.WithStack(errRequiresAuthentication)
	}

	if authTime, ok = authTimeFromCookie(cookie); !ok {
		if contains(prompts, "none") {
			return nil, errors.WithStack(fosite.ErrLoginRequired.WithDebug("Session not set but prompt is set to none"))
		}
		return nil, errors.WithStack(errRequiresAuthentication)
	}

	if maxAge > 0 && authTime.UTC().Add(time.Second*time.Duration(maxAge)).Before(time.Now().UTC()) {
		if contains(prompts, "none") {
			return nil, errors.WithStack(fosite.ErrLoginRequired.WithDebug("Max age was reached and prompt is set to none"))
		}
		return nil, errors.WithStack(errRequiresAuthentication)
	}

	consent, err := s.ConsentManager.GetPreviouslyGrantedConsent(user, req.GetClient().GetID(), req.GetRequestedScopes())
	if err != nil {
		if contains(prompts, "none") {
			return nil, errors.WithStack(fosite.ErrLoginRequired.WithDebug(err.Error()))
		} else {
			return nil, errors.WithStack(errRequiresAuthentication)
		}
	}

	for _, scope := range req.GetRequestedScopes() {
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
				AuthTime:    consent.RequestedAt,
				RequestedAt: time.Now().UTC(),
				Extra:       consent.IDTokenExtra,
				AuthenticationContextClassReference: "0",
			},
			// required for lookup on jwk endpoint
			Headers: &ejwt.Headers{Extra: map[string]interface{}{"kid": jwtIdTokenKey}},
			Subject: consent.Subject,
		},
		Extra: consent.AccessTokenExtra,
	}, nil

}

func (s *DefaultConsentStrategy) CreateConsentRequest(
	req fosite.AuthorizeRequester,
	redirectURL string,
	cookie *sessions.Session,
) (string, error) {
	csrf := uuid.New()
	id := uuid.New()

	requestClient := req.GetClient().(*client.Client)
	requestClient.Secret = ""

	var prompt string
	if req.GetRequestForm().Get("prompt") == "consent" {
		prompt = "consent"
	}

	cookie.Values[cookieCSRFKey] = csrf
	consent := &ConsentRequest{
		ID:              id,
		CSRF:            csrf,
		GrantedScopes:   []string{},
		RequestedScopes: req.GetRequestedScopes(),
		ClientID:        req.GetClient().GetID(),
		Client:          requestClient,
		ExpiresAt:       time.Now().UTC().Add(s.DefaultChallengeLifespan),
		RequestedAt:     time.Now().UTC(),
		RedirectURL:     redirectURL + "&consent=" + id,
		OpenIDConnectContext: &ConsentRequestOpenIDConnectContext{
			Prompt:    prompt,
			UILocales: req.GetRequestForm().Get("ui_locales"),
			Display:   req.GetRequestForm().Get("display"),
			LoginHint: req.GetRequestForm().Get("login_hint"),
			ACRValues: pkg.SplitNonEmpty(req.GetRequestForm().Get("acr_values"), " "),
		},
		AccessTokenExtra: map[string]interface{}{},
		IDTokenExtra:     map[string]interface{}{},
	}

	if err := s.ConsentManager.PersistConsentRequest(consent); err != nil {
		return "", errors.WithStack(err)
	}

	return id, nil
}
