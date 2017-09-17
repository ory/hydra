package oauth2

import (
	"time"

	"net/http"

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

func (s *DefaultConsentStrategy) ValidateConsentRequest(req fosite.AuthorizeRequester, session string, cookie *sessions.Session) (claims *Session, err error) {
	consent, err := s.ConsentManager.GetConsentRequest(session)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if time.Now().After(consent.ExpiresAt) {
		return nil, errors.Errorf("Token expired")
	}

	if consent.Audience != req.GetClient().GetID() {
		return nil, errors.Errorf("Audience mismatch")
	}

	if consent.Subject == "" {
		return nil, errors.Errorf("Subject key is empty or undefined in consent response, check your payload.")
	}

	if j, ok := cookie.Values[CookieCSRFKey]; !ok {
		return nil, errors.Errorf("Session cookie is missing anti-replay token")
	} else if js, ok := j.(string); !ok {
		return nil, errors.Errorf("Session cookie anti-replay value is not a string")
	} else if js != consent.CSRF {
		return nil, errors.Errorf("Session cookie anti-replay value does not match value from consent response")
	}

	if !consent.IsConsentGranted() {
		err := errors.New("The resource owner denied consent for this request")
		return nil, &fosite.RFC6749Error{
			Name:        "Resource owner denied consent",
			Description: err.Error(),
			Debug:       err.Error(),
			Hint:        "Token validation failed.",
			Code:        http.StatusUnauthorized,
		}
	}

	for _, scope := range consent.GrantedScopes {
		req.GrantScope(scope)
	}

	delete(cookie.Values, CookieCSRFKey)

	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims: &ejwt.IDTokenClaims{
				Audience:  req.GetClient().GetID(),
				Subject:   consent.Subject,
				Issuer:    s.Issuer,
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(s.DefaultIDTokenLifespan),
				Extra:     consent.IDTokenExtra,
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
		RequestedScope:   req.GetRequestedScopes(),
		Audience:         req.GetClient().GetID(),
		ExpiresAt:        time.Now().Add(s.DefaultChallengeLifespan),
		RedirectURL:      redirectURL + "&consent=" + id,
		AccessTokenExtra: map[string]interface{}{},
		IDTokenExtra:     map[string]interface{}{},
	}

	if err := s.ConsentManager.PersistConsentRequest(consent); err != nil {
		return "", errors.WithStack(err)
	}

	return id, nil
}
