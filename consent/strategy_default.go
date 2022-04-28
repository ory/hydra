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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ory/hydra/driver/config"

	"github.com/ory/x/errorsx"

	"github.com/ory/x/sqlcon"

	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	jwtgo "github.com/ory/fosite/token/jwt"

	"github.com/ory/x/sqlxx"

	"github.com/ory/x/httpx"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/x/mapx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/stringsx"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/x"
)

const (
	CookieAuthenticationName    = "oauth2_authentication_session"
	CookieAuthenticationSIDName = "sid"

	cookieAuthenticationCSRFName = "oauth2_authentication_csrf"
	cookieConsentCSRFName        = "oauth2_consent_csrf"
)

type DefaultStrategy struct {
	c *config.Provider
	r InternalRegistry
}

func NewStrategy(
	r InternalRegistry,
	c *config.Provider,
) *DefaultStrategy {
	return &DefaultStrategy{
		c: c,
		r: r,
	}
}

var ErrAbortOAuth2Request = errors.New("the OAuth 2.0 Authorization request must be aborted")
var ErrNoPreviousConsentFound = errors.New("no previous OAuth 2.0 Consent could be found for this access request")
var ErrNoAuthenticationSessionFound = errors.New("no previous login session was found")
var ErrHintDoesNotMatchAuthentication = errors.New("subject from hint does not match subject from session")

func (s *DefaultStrategy) matchesValueFromSession(ctx context.Context, c fosite.Client, hintSubject string, sessionSubject string) error {
	obfuscatedUserID, err := s.obfuscateSubjectIdentifier(c, sessionSubject, "")
	if err != nil {
		return err
	}

	var forcedObfuscatedUserID string
	if s, err := s.r.ConsentManager().GetForcedObfuscatedLoginSession(ctx, c.GetID(), hintSubject); errors.Is(err, x.ErrNotFound) {
		// do nothing
	} else if err != nil {
		return err
	} else {
		forcedObfuscatedUserID = s.SubjectObfuscated
	}

	if hintSubject != sessionSubject && hintSubject != obfuscatedUserID && hintSubject != forcedObfuscatedUserID {
		return ErrHintDoesNotMatchAuthentication
	}

	return nil
}

func (s *DefaultStrategy) authenticationSession(w http.ResponseWriter, r *http.Request) (*LoginSession, error) {
	// We try to open the session cookie. If it does not exist (indicated by the error), we must authenticate the user.
	cookie, err := s.r.CookieStore().Get(r, CookieName(s.c.TLS(config.PublicInterface).Enabled(), CookieAuthenticationName))
	if err != nil {
		s.r.Logger().
			WithRequest(r).
			WithError(err).Debug("User logout skipped because cookie store returned an error.")
		return nil, errorsx.WithStack(ErrNoAuthenticationSessionFound)
	}

	sessionID := mapx.GetStringDefault(cookie.Values, CookieAuthenticationSIDName, "")
	if sessionID == "" {
		s.r.Logger().
			WithRequest(r).
			Debug("User logout skipped because cookie exists but session value is empty.")
		return nil, errorsx.WithStack(ErrNoAuthenticationSessionFound)
	}

	session, err := s.r.ConsentManager().GetRememberedLoginSession(r.Context(), sessionID)
	if errors.Is(err, x.ErrNotFound) {
		s.r.Logger().WithRequest(r).WithError(err).
			Debug("User logout skipped because cookie exists and session value exist but are not remembered any more.")
		return nil, errorsx.WithStack(ErrNoAuthenticationSessionFound)
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *DefaultStrategy) requestAuthentication(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester) error {
	prompt := stringsx.Splitx(ar.GetRequestForm().Get("prompt"), " ")
	if stringslice.Has(prompt, "login") {
		return s.forwardAuthenticationRequest(w, r, ar, "", time.Time{}, nil)
	}

	session, err := s.authenticationSession(w, r)
	if errors.Is(err, ErrNoAuthenticationSessionFound) {
		return s.forwardAuthenticationRequest(w, r, ar, "", time.Time{}, nil)
	} else if err != nil {
		return err
	}

	maxAge := int64(0)
	if ma := ar.GetRequestForm().Get("max_age"); len(ma) > 0 {
		var err error
		maxAge, err = strconv.ParseInt(ma, 10, 64)
		if err != nil {
			return err
		}
	}

	if maxAge > 0 && time.Time(session.AuthenticatedAt).UTC().Add(time.Second*time.Duration(maxAge)).Before(time.Now().UTC()) {
		if stringslice.Has(prompt, "none") {
			return errorsx.WithStack(fosite.ErrLoginRequired.WithHint("Request failed because prompt is set to 'none' and authentication time reached 'max_age'."))
		}
		return s.forwardAuthenticationRequest(w, r, ar, "", time.Time{}, nil)
	}

	idTokenHint := ar.GetRequestForm().Get("id_token_hint")
	if idTokenHint == "" {
		return s.forwardAuthenticationRequest(w, r, ar, session.Subject, time.Time(session.AuthenticatedAt), session)
	}

	hintSub, err := s.getSubjectFromIDTokenHint(r.Context(), idTokenHint)
	if err != nil {
		return err
	}

	if err := s.matchesValueFromSession(r.Context(), ar.GetClient(), hintSub, session.Subject); errors.Is(err, ErrHintDoesNotMatchAuthentication) {
		return errorsx.WithStack(fosite.ErrLoginRequired.WithHint("Request failed because subject claim from id_token_hint does not match subject from authentication session."))
	}

	return s.forwardAuthenticationRequest(w, r, ar, session.Subject, time.Time(session.AuthenticatedAt), session)
}

func (s *DefaultStrategy) getIDTokenHintClaims(ctx context.Context, idTokenHint string) (jwtgo.MapClaims, error) {
	token, err := s.r.OpenIDJWTStrategy().Decode(ctx, idTokenHint)
	if ve := new(jwtgo.ValidationError); errors.As(err, &ve) && ve.Errors == jwtgo.ValidationErrorExpired {
		// Expired is ok
	} else if err != nil {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(err.Error()))
	}
	return token.Claims, nil
}

func (s *DefaultStrategy) getSubjectFromIDTokenHint(ctx context.Context, idTokenHint string) (string, error) {
	claims, err := s.getIDTokenHintClaims(ctx, idTokenHint)
	if err != nil {
		return "", err
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Failed to validate OpenID Connect request because provided id token from id_token_hint does not have a subject."))
	}

	return sub, nil
}

func (s *DefaultStrategy) forwardAuthenticationRequest(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, subject string, authenticatedAt time.Time, session *LoginSession) error {
	if (subject != "" && authenticatedAt.IsZero()) || (subject == "" && !authenticatedAt.IsZero()) {
		return errorsx.WithStack(fosite.ErrServerError.WithHint("Consent strategy returned a non-empty subject with an empty auth date, or an empty subject with a non-empty auth date."))
	}

	skip := false
	if subject != "" {
		skip = true
	}

	// Let'id validate that prompt is actually not "none" if we can't skip authentication
	prompt := stringsx.Splitx(ar.GetRequestForm().Get("prompt"), " ")
	if stringslice.Has(prompt, "none") && !skip {
		return errorsx.WithStack(fosite.ErrLoginRequired.WithHint(`Prompt 'none' was requested, but no existing login session was found.`))
	}

	// Set up csrf/challenge/verifier values
	verifier := strings.Replace(uuid.New(), "-", "", -1)
	challenge := strings.Replace(uuid.New(), "-", "", -1)
	csrf := strings.Replace(uuid.New(), "-", "", -1)

	// Generate the request URL
	iu := s.c.OAuth2AuthURL()
	iu.RawQuery = r.URL.RawQuery

	var idTokenHintClaims jwtgo.MapClaims
	if idTokenHint := ar.GetRequestForm().Get("id_token_hint"); len(idTokenHint) > 0 {
		claims, err := s.getIDTokenHintClaims(r.Context(), idTokenHint)
		if err != nil {
			return err
		}

		idTokenHintClaims = claims
	}

	sessionID := uuid.New()
	if session != nil {
		sessionID = session.ID
	} else {
		// Create a stub session so that we can later update it.
		if err := s.r.ConsentManager().CreateLoginSession(r.Context(), &LoginSession{ID: sessionID}); err != nil {
			return err
		}
	}

	// Set the session
	if err := s.r.ConsentManager().CreateLoginRequest(
		r.Context(),
		&LoginRequest{
			ID:                challenge,
			Verifier:          verifier,
			CSRF:              csrf,
			Skip:              skip,
			RequestedScope:    []string(ar.GetRequestedScopes()),
			RequestedAudience: []string(ar.GetRequestedAudience()),
			Subject:           subject,
			Client:            sanitizeClientFromRequest(ar),
			RequestURL:        iu.String(),
			AuthenticatedAt:   sqlxx.NullTime(authenticatedAt),
			RequestedAt:       time.Now().Truncate(time.Second).UTC(),
			SessionID:         sqlxx.NullString(sessionID),
			OpenIDConnectContext: &OpenIDConnectContext{
				IDTokenHintClaims: idTokenHintClaims,
				ACRValues:         stringsx.Splitx(ar.GetRequestForm().Get("acr_values"), " "),
				UILocales:         stringsx.Splitx(ar.GetRequestForm().Get("ui_locales"), " "),
				Display:           ar.GetRequestForm().Get("display"),
				LoginHint:         ar.GetRequestForm().Get("login_hint"),
			},
		},
	); err != nil {
		return errorsx.WithStack(err)
	}

	if err := createCsrfSession(w, r, s.r.CookieStore(), cookieAuthenticationCSRFName, csrf, s.c.TLS(config.PublicInterface).Enabled(), s.c.CookieSameSiteMode(), s.c.CookieSameSiteLegacyWorkaround()); err != nil {
		return errorsx.WithStack(err)
	}

	http.Redirect(w, r, urlx.SetQuery(s.c.LoginURL(), url.Values{"login_challenge": {challenge}}).String(), http.StatusFound)

	// generate the verifier
	return errorsx.WithStack(ErrAbortOAuth2Request)
}

func (s *DefaultStrategy) revokeAuthenticationSession(w http.ResponseWriter, r *http.Request) error {
	sid, err := s.revokeAuthenticationCookie(w, r, s.r.CookieStore())
	if err != nil {
		return err
	}

	if sid == "" {
		return nil
	}

	return s.r.ConsentManager().DeleteLoginSession(r.Context(), sid)
}

func (s *DefaultStrategy) revokeAuthenticationCookie(w http.ResponseWriter, r *http.Request, ss sessions.Store) (string, error) {
	cookie, _ := ss.Get(r, CookieName(s.c.TLS(config.PublicInterface).Enabled(), CookieAuthenticationName))
	sid, _ := mapx.GetString(cookie.Values, CookieAuthenticationSIDName)

	cookie.Values[CookieAuthenticationSIDName] = ""
	cookie.Options.HttpOnly = true
	cookie.Options.SameSite = s.c.CookieSameSiteMode()
	cookie.Options.Secure = s.c.TLS(config.PublicInterface).Enabled()
	cookie.Options.MaxAge = -1

	if err := cookie.Save(r, w); err != nil {
		return "", errorsx.WithStack(err)
	}

	return sid, nil
}

func (s *DefaultStrategy) obfuscateSubjectIdentifier(cl fosite.Client, subject, forcedIdentifier string) (string, error) {
	if c, ok := cl.(*client.Client); ok && c.SubjectType == "pairwise" {
		algorithm, ok := s.r.SubjectIdentifierAlgorithm()[c.SubjectType]
		if !ok {
			return "", errorsx.WithStack(fosite.ErrInvalidRequest.WithHintf(`Subject Identifier Algorithm '%s' was requested by OAuth 2.0 Client '%s' but is not configured.`, c.SubjectType, c.OutfacingID))
		}

		if len(forcedIdentifier) > 0 {
			return forcedIdentifier, nil
		}

		return algorithm.Obfuscate(subject, c)
	} else if !ok {
		return "", errors.New("Unable to type assert OAuth 2.0 Client to *client.Client")
	}
	return subject, nil
}

func (s *DefaultStrategy) verifyAuthentication(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester, verifier string) (*HandledLoginRequest, error) {
	ctx := r.Context()
	session, err := s.r.ConsentManager().VerifyAndInvalidateLoginRequest(ctx, verifier)
	if errors.Is(err, sqlcon.ErrNoRows) {
		return nil, errorsx.WithStack(fosite.ErrAccessDenied.WithHint("The login verifier has already been used, has not been granted, or is invalid."))
	} else if err != nil {
		return nil, err
	}

	if session.HasError() {
		session.Error.SetDefaults(loginRequestDeniedErrorName)
		return nil, errorsx.WithStack(session.Error.toRFCError())
	}

	if session.RequestedAt.Add(s.c.ConsentRequestMaxAge()).Before(time.Now()) {
		return nil, errorsx.WithStack(fosite.ErrRequestUnauthorized.WithHint("The login request has expired. Please try again."))
	}

	if err := validateCsrfSession(r, s.r.CookieStore(), cookieAuthenticationCSRFName, session.LoginRequest.CSRF, s.c.CookieSameSiteLegacyWorkaround(), s.c.TLS(config.PublicInterface).Enabled()); err != nil {
		return nil, err
	}

	if session.LoginRequest.Skip && !session.Remember {
		return nil, errorsx.WithStack(fosite.ErrServerError.WithHint("The login request was previously remembered and can only be forgotten using the reject feature."))
	}

	if session.LoginRequest.Skip && session.Subject != session.LoginRequest.Subject {
		// Revoke the session because there's clearly a mix up wrt the subject that's being authenticated
		if err := s.revokeAuthenticationSession(w, r); err != nil {
			return nil, err
		}

		return nil, errorsx.WithStack(fosite.ErrServerError.WithHint("The login request is marked as remember, but the subject from the login confirmation does not match the original subject from the cookie."))
	}

	subjectIdentifier, err := s.obfuscateSubjectIdentifier(req.GetClient(), session.Subject, session.ForceSubjectIdentifier)
	if err != nil {
		return nil, err
	}

	sessionID := session.LoginRequest.SessionID.String()

	if err := s.r.OpenIDConnectRequestValidator().ValidatePrompt(ctx, &fosite.AuthorizeRequest{
		ResponseTypes: req.GetResponseTypes(),
		RedirectURI:   req.GetRedirectURI(),
		State:         req.GetState(),
		// HandledResponseTypes, this can be safely ignored because it's not being used by validation
		Request: fosite.Request{
			ID:                req.GetID(),
			RequestedAt:       req.GetRequestedAt(),
			Client:            req.GetClient(),
			RequestedAudience: req.GetRequestedAudience(),
			GrantedAudience:   req.GetGrantedAudience(),
			RequestedScope:    req.GetRequestedScopes(),
			GrantedScope:      req.GetGrantedScopes(),
			Form:              req.GetRequestForm(),
			Session: &openid.DefaultSession{
				Claims: &jwt.IDTokenClaims{
					Subject:     subjectIdentifier,
					IssuedAt:    time.Now().UTC(),                // doesn't matter
					ExpiresAt:   time.Now().Add(time.Hour).UTC(), // doesn't matter
					AuthTime:    time.Time(session.AuthenticatedAt),
					RequestedAt: session.RequestedAt,
				},
				Headers: &jwt.Headers{},
				Subject: session.Subject,
			},
		},
	}); errors.Is(err, fosite.ErrLoginRequired) {
		// This indicates that something went wrong with checking the subject id - let's destroy the session to be safe
		if err := s.revokeAuthenticationSession(w, r); err != nil {
			return nil, err
		}

		return nil, err
	} else if err != nil {
		return nil, err
	}

	if session.ForceSubjectIdentifier != "" {
		if err := s.r.ConsentManager().CreateForcedObfuscatedLoginSession(r.Context(), &ForcedObfuscatedLoginSession{
			Subject:           session.Subject,
			ClientID:          req.GetClient().GetID(),
			SubjectObfuscated: session.ForceSubjectIdentifier,
		}); err != nil {
			return nil, err
		}
	}

	if !session.LoginRequest.Skip {
		if time.Time(session.AuthenticatedAt).IsZero() {
			return nil, errorsx.WithStack(fosite.ErrServerError.WithHint("Expected the handled login request to contain a valid authenticated_at value but it was zero. This is a bug which should be reported to https://github.com/ory/hydra."))
		}

		if err := s.r.ConsentManager().ConfirmLoginSession(r.Context(), sessionID, time.Time(session.AuthenticatedAt), session.Subject, session.Remember); err != nil {
			return nil, err
		}
	}

	if !session.Remember && !session.LoginRequest.Skip {
		// If the session should not be remembered (and we're actually not skipping), than the user clearly don't
		// wants us to store a cookie. So let's bust the authentication session (if one exists).
		if err := s.revokeAuthenticationSession(w, r); err != nil {
			return nil, err
		}
	}

	if !session.Remember || session.LoginRequest.Skip {
		// If the user doesn't want to remember the session, we do not store a cookie.
		// If login was skipped, it means an authentication cookie was present and
		// we don't want to touch it (in order to preserve its original expiry date)
		return session, nil
	}

	// Not a skipped login and the user asked to remember its session, store a cookie
	cookie, _ := s.r.CookieStore().Get(r, CookieName(s.c.TLS(config.PublicInterface).Enabled(), CookieAuthenticationName))
	cookie.Values[CookieAuthenticationSIDName] = sessionID
	if session.RememberFor >= 0 {
		cookie.Options.MaxAge = session.RememberFor
	}
	cookie.Options.HttpOnly = true
	cookie.Options.SameSite = s.c.CookieSameSiteMode()
	cookie.Options.Secure = s.c.TLS(config.PublicInterface).Enabled()
	if err := cookie.Save(r, w); err != nil {
		return nil, errorsx.WithStack(err)
	}

	s.r.Logger().WithRequest(r).
		WithFields(logrus.Fields{
			"cookie_name":      CookieName(s.c.TLS(config.PublicInterface).Enabled(), CookieAuthenticationName),
			"cookie_http_only": true,
			"cookie_same_site": s.c.CookieSameSiteMode(),
			"cookie_secure":    s.c.TLS(config.PublicInterface).Enabled(),
		}).Debug("Authentication session cookie was set.")
	return session, nil
}

func (s *DefaultStrategy) requestConsent(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, authenticationSession *HandledLoginRequest) error {
	prompt := stringsx.Splitx(ar.GetRequestForm().Get("prompt"), " ")
	if stringslice.Has(prompt, "consent") {
		return s.forwardConsentRequest(w, r, ar, authenticationSession, nil)
	}

	// https://tools.ietf.org/html/rfc6749
	//
	// As stated in Section 10.2 of OAuth 2.0 [RFC6749], the authorization
	// server SHOULD NOT process authorization requests automatically
	// without user consent or interaction, except when the identity of the
	// client can be assured.  This includes the case where the user has
	// previously approved an authorization request for a given client id --
	// unless the identity of the client can be proven, the request SHOULD
	// be processed as if no previous request had been approved.
	//
	// Measures such as claimed "https" scheme redirects MAY be accepted by
	// authorization servers as identity proof.  Some operating systems may
	// offer alternative platform-specific identity features that MAY be
	// accepted, as appropriate.
	if ar.GetClient().IsPublic() {
		// The OpenID Connect Test Tool fails if this returns `consent_required` when `prompt=none` is used.
		// According to the quote above, it should be ok to allow https to skip consent.
		//
		// This is tracked as issue: https://github.com/ory/hydra/issues/866
		// This is also tracked as upstream issue: https://github.com/openid-certification/oidctest/issues/97
		if !(ar.GetRedirectURI().Scheme == "https" || (fosite.IsLocalhost(ar.GetRedirectURI()) && ar.GetRedirectURI().Scheme == "http")) {
			return s.forwardConsentRequest(w, r, ar, authenticationSession, nil)
		}
	}

	// This breaks OIDC Conformity Tests and is probably a bit paranoid.
	//
	// if ar.GetResponseTypes().Has("token") {
	// 	 // We're probably requesting the implicit or hybrid flow in which case we MUST authenticate and authorize the request
	// 	 return s.forwardConsentRequest(w, r, ar, authenticationSession, nil)
	// }

	consentSessions, err := s.r.ConsentManager().FindGrantedAndRememberedConsentRequests(r.Context(), ar.GetClient().GetID(), authenticationSession.Subject)
	if errors.Is(err, ErrNoPreviousConsentFound) {
		return s.forwardConsentRequest(w, r, ar, authenticationSession, nil)
	} else if err != nil {
		return err
	}

	if found := matchScopes(s.r.ScopeStrategy(), consentSessions, ar.GetRequestedScopes()); found != nil {
		return s.forwardConsentRequest(w, r, ar, authenticationSession, found)
	}

	return s.forwardConsentRequest(w, r, ar, authenticationSession, nil)
}

func (s *DefaultStrategy) forwardConsentRequest(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, as *HandledLoginRequest, cs *HandledConsentRequest) error {
	skip := false
	if cs != nil {
		skip = true
	}

	prompt := stringsx.Splitx(ar.GetRequestForm().Get("prompt"), " ")
	if stringslice.Has(prompt, "none") && !skip {
		return errorsx.WithStack(fosite.ErrConsentRequired.WithHint(`Prompt 'none' was requested, but no previous consent was found.`))
	}

	// Set up csrf/challenge/verifier values
	verifier := strings.Replace(uuid.New(), "-", "", -1)
	challenge := strings.Replace(uuid.New(), "-", "", -1)
	csrf := strings.Replace(uuid.New(), "-", "", -1)

	if err := s.r.ConsentManager().CreateConsentRequest(
		r.Context(),
		&ConsentRequest{
			ID:                     challenge,
			ACR:                    as.ACR,
			AMR:                    as.AMR,
			Verifier:               verifier,
			CSRF:                   csrf,
			Skip:                   skip,
			RequestedScope:         []string(ar.GetRequestedScopes()),
			RequestedAudience:      []string(ar.GetRequestedAudience()),
			Subject:                as.Subject,
			Client:                 sanitizeClientFromRequest(ar),
			RequestURL:             as.LoginRequest.RequestURL,
			AuthenticatedAt:        as.AuthenticatedAt,
			RequestedAt:            as.RequestedAt,
			ForceSubjectIdentifier: as.ForceSubjectIdentifier,
			OpenIDConnectContext:   as.LoginRequest.OpenIDConnectContext,
			LoginSessionID:         as.LoginRequest.SessionID,
			LoginChallenge:         sqlxx.NullString(as.LoginRequest.ID),
			Context:                as.Context,
		},
	); err != nil {
		return errorsx.WithStack(err)
	}

	if err := createCsrfSession(w, r, s.r.CookieStore(), cookieConsentCSRFName, csrf, s.c.TLS(config.PublicInterface).Enabled(), s.c.CookieSameSiteMode(), s.c.CookieSameSiteLegacyWorkaround()); err != nil {
		return errorsx.WithStack(err)
	}

	http.Redirect(
		w, r,
		urlx.SetQuery(s.c.ConsentURL(), url.Values{"consent_challenge": {challenge}}).String(),
		http.StatusFound,
	)

	// generate the verifier
	return errorsx.WithStack(ErrAbortOAuth2Request)
}

func (s *DefaultStrategy) verifyConsent(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester, verifier string) (*HandledConsentRequest, error) {
	session, err := s.r.ConsentManager().VerifyAndInvalidateConsentRequest(r.Context(), verifier)
	if errors.Is(err, sqlcon.ErrNoRows) {
		return nil, errorsx.WithStack(fosite.ErrAccessDenied.WithHint("The consent verifier has already been used, has not been granted, or is invalid."))
	} else if err != nil {
		return nil, err
	}

	if session.RequestedAt.Add(s.c.ConsentRequestMaxAge()).Before(time.Now()) {
		return nil, errorsx.WithStack(fosite.ErrRequestUnauthorized.WithHint("The consent request has expired, please try again."))
	}

	if session.HasError() {
		session.Error.SetDefaults(consentRequestDeniedErrorName)
		return nil, errorsx.WithStack(session.Error.toRFCError())
	}

	if time.Time(session.ConsentRequest.AuthenticatedAt).IsZero() {
		return nil, errorsx.WithStack(fosite.ErrServerError.WithHint("The authenticatedAt value was not set."))
	}

	if err := validateCsrfSession(r, s.r.CookieStore(), cookieConsentCSRFName, session.ConsentRequest.CSRF, s.c.CookieSameSiteLegacyWorkaround(), s.c.TLS(config.PublicInterface).Enabled()); err != nil {
		return nil, err
	}

	pw, err := s.obfuscateSubjectIdentifier(req.GetClient(), session.ConsentRequest.Subject, session.ConsentRequest.ForceSubjectIdentifier)
	if err != nil {
		return nil, err
	}

	if session.Session == nil {
		session.Session = NewConsentRequestSessionData()
	}

	if session.Session.AccessToken == nil {
		session.Session.AccessToken = map[string]interface{}{}
	}

	if session.Session.IDToken == nil {
		session.Session.IDToken = map[string]interface{}{}
	}

	session.ConsentRequest.SubjectIdentifier = pw
	session.AuthenticatedAt = session.ConsentRequest.AuthenticatedAt
	return session, nil
}

func (s *DefaultStrategy) generateFrontChannelLogoutURLs(ctx context.Context, subject, sid string) ([]string, error) {
	clients, err := s.r.ConsentManager().ListUserAuthenticatedClientsWithFrontChannelLogout(ctx, subject, sid)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, c := range clients {
		u, err := url.Parse(c.FrontChannelLogoutURI)
		if err != nil {
			return nil, errorsx.WithStack(fosite.ErrServerError.WithHintf("Unable to parse frontchannel_logout_uri because %s.", c.FrontChannelLogoutURI).WithDebug(err.Error()))
		}

		urls = append(urls, urlx.SetQuery(u, url.Values{
			"iss": {s.c.IssuerURL().String()},
			"sid": {sid},
		}).String())
	}

	return urls, nil
}

func (s *DefaultStrategy) executeBackChannelLogout(ctx context.Context, r *http.Request, subject, sid string) error {
	clients, err := s.r.ConsentManager().ListUserAuthenticatedClientsWithBackChannelLogout(ctx, subject, sid)
	if err != nil {
		return err
	}

	openIDKeyID, err := s.r.OpenIDJWTStrategy().GetPublicKeyID(ctx)
	if err != nil {
		return err
	}

	type task struct {
		url      string
		token    string
		clientID string
	}

	var tasks []task
	for _, c := range clients {
		// Getting the forced obfuscated login session is tricky because the user id could be obfuscated with a new
		// ID every time the algorithm is used. Thus, we would only get the most recent version. It therefore makes
		// sense to just use the sid.
		//
		// s.r.ConsentManager().GetForcedObfuscatedLoginSession(context.Background(), subject, <missing>)
		// sub := s.obfuscateSubjectIdentifier(c, subject, )

		t, _, err := s.r.OpenIDJWTStrategy().Generate(ctx, jwtgo.MapClaims{
			"iss":    s.c.IssuerURL().String(),
			"aud":    []string{c.OutfacingID},
			"iat":    time.Now().UTC().Unix(),
			"jti":    uuid.New(),
			"events": map[string]struct{}{"http://schemas.openid.net/event/backchannel-logout": {}},
			"sid":    sid,
		}, &jwt.Headers{
			Extra: map[string]interface{}{"kid": openIDKeyID},
		})
		if err != nil {
			return err
		}

		tasks = append(tasks, task{url: c.BackChannelLogoutURI, clientID: c.OutfacingID, token: t})
	}

	hc := httpx.NewResilientClient()

	var execute = func(t task) {
		log := s.r.Logger().WithRequest(r).
			WithField("client_id", t.clientID).
			WithField("backchannel_logout_url", t.url)

		res, err := hc.PostForm(t.url, url.Values{"logout_token": {t.token}})
		if err != nil {
			log.WithError(err).Error("Unable to execute OpenID Connect Back-Channel Logout Request")
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.WithError(errors.Errorf("expected HTTP status code %d but got %d", http.StatusOK, res.StatusCode)).
				Error("Unable to execute OpenID Connect Back-Channel Logout Request")
			return
		} else {
			log.Info("Back-Channel Logout Request")
		}
	}

	for _, t := range tasks {
		go execute(t)
	}

	return nil
}

func (s *DefaultStrategy) issueLogoutVerifier(w http.ResponseWriter, r *http.Request) (*LogoutResult, error) {
	// There are two types of log out flows:
	//
	// - RP initiated logout
	// - OP initiated logout

	// Per default, we're redirecting to the global redirect URL. This is assuming that we're not an RP-initiated
	// logout flow.
	redir := s.c.LogoutRedirectURL().String()

	if err := r.ParseForm(); err != nil {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.
			WithHintf("Logout failed because the '%s' request could not be parsed.", r.Method),
		)
	}

	hint := r.Form.Get("id_token_hint")
	state := r.Form.Get("state")
	requestedRedir := r.Form.Get("post_logout_redirect_uri")

	if len(hint) == 0 {
		// hint is not set, so this is an OP initiated logout

		if len(state) > 0 {
			// state can only be set if it's an RP-initiated logout flow. If not, we should throw an error.
			return nil, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter state is set but id_token_hint is missing."))
		}

		if len(requestedRedir) > 0 {
			// post_logout_redirect_uri can only be set if it's an RP-initiated logout flow. If not, we should throw an error.
			return nil, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter post_logout_redirect_uri is set but id_token_hint is missing."))
		}

		session, err := s.authenticationSession(w, r)
		if errors.Is(err, ErrNoAuthenticationSessionFound) {
			// OP initiated log out but no session was found. Since we can not identify the user we can not call
			// any RPs.
			s.r.AuditLogger().
				WithRequest(r).
				Info("User logout skipped because no authentication session exists.")
			http.Redirect(w, r, redir, http.StatusFound)
			return nil, errorsx.WithStack(ErrAbortOAuth2Request)
		} else if err != nil {
			return nil, err
		}

		challenge := uuid.New()
		if err := s.r.ConsentManager().CreateLogoutRequest(r.Context(), &LogoutRequest{
			RequestURL:  r.URL.String(),
			ID:          challenge,
			Subject:     session.Subject,
			SessionID:   session.ID,
			Verifier:    uuid.New(),
			RPInitiated: false,

			// PostLogoutRedirectURI is set to the value from config.Provider().LogoutRedirectURL()
			PostLogoutRedirectURI: redir,
		}); err != nil {
			return nil, err
		}

		s.r.AuditLogger().
			WithRequest(r).
			Info("User logout requires user confirmation, redirecting to Logout UI.")
		http.Redirect(w, r, urlx.SetQuery(s.c.LogoutURL(), url.Values{"logout_challenge": {challenge}}).String(), http.StatusFound)
		return nil, errorsx.WithStack(ErrAbortOAuth2Request)
	}

	claims, err := s.getIDTokenHintClaims(r.Context(), hint)
	if err != nil {
		return nil, err
	}

	mksi := mapx.KeyStringToInterface(claims)
	if !claims.VerifyIssuer(s.c.IssuerURL().String(), true) {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.
			WithHintf(
				`Logout failed because issuer claim value '%s' from query parameter id_token_hint does not match with issuer value from configuration '%s'.`,
				mapx.GetStringDefault(mksi, "iss", ""),
				s.c.IssuerURL().String(),
			),
		)
	}

	now := time.Now().UTC().Unix()
	if !claims.VerifyIssuedAt(now, true) {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.
			WithHintf(
				`Logout failed because iat claim value '%.0f' from query parameter id_token_hint is before now ('%d').`,
				mapx.GetFloat64Default(mksi, "iat", float64(0)),
				now,
			),
		)
	}

	hintSid := mapx.GetStringDefault(mksi, "sid", "")
	if len(hintSid) == 0 {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter id_token_hint is missing sid claim."))
	}

	// It doesn't really make sense to use the subject value from the ID Token because it might be obfuscated.
	if hintSub := mapx.GetStringDefault(mksi, "sub", ""); len(hintSub) == 0 {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter id_token_hint is missing sub claim."))
	}

	// Let's find the client by cycling through the audiences. Typically, we only have one audience
	var cl *client.Client
	for _, aud := range mapx.GetStringSliceDefault(
		mksi,
		"aud",
		[]string{
			mapx.GetStringDefault(mksi, "aud", ""),
		},
	) {
		c, err := s.r.ClientManager().GetConcreteClient(r.Context(), aud)
		if errors.Is(err, x.ErrNotFound) {
			continue
		} else if err != nil {
			return nil, err
		}
		cl = c
		break
	}

	if cl == nil {
		return nil, errorsx.WithStack(fosite.ErrInvalidRequest.
			WithHint("Logout failed because none of the listed audiences is a registered OAuth 2.0 Client."))
	}

	if len(requestedRedir) > 0 {
		var f *url.URL
		for _, w := range cl.PostLogoutRedirectURIs {
			if w == requestedRedir {
				u, err := url.Parse(w)
				if err != nil {
					return nil, errorsx.WithStack(fosite.ErrServerError.WithHintf("Unable to parse post_logout_redirect_uri '%s'.", w).WithDebug(err.Error()))
				}

				f = u
			}
		}

		if f == nil {
			return nil, errorsx.WithStack(fosite.ErrInvalidRequest.
				WithHint("Logout failed because query parameter post_logout_redirect_uri is not a whitelisted as a post_logout_redirect_uri for the client."),
			)
		}

		params := url.Values{}
		if state != "" {
			params.Add("state", state)
		}

		redir = urlx.SetQuery(f, params).String()
	}

	// We do not really want to verify if the user (from id token hint) has a session here because it doesn't really matter.
	// Instead, we'll check this when we're actually revoking the cookie!
	session, err := s.r.ConsentManager().GetRememberedLoginSession(r.Context(), hintSid)
	if errors.Is(err, x.ErrNotFound) {
		// Such a session does not exist - maybe it has already been revoked? In any case, we can't do much except
		// leaning back and redirecting back.
		http.Redirect(w, r, redir, http.StatusFound)
		return nil, errorsx.WithStack(ErrAbortOAuth2Request)
	} else if err != nil {
		return nil, err
	}

	challenge := uuid.New()
	if err := s.r.ConsentManager().CreateLogoutRequest(r.Context(), &LogoutRequest{
		RequestURL:  r.URL.String(),
		ID:          challenge,
		SessionID:   hintSid,
		Subject:     session.Subject,
		Verifier:    uuid.New(),
		Client:      cl,
		RPInitiated: true,

		// PostLogoutRedirectURI is set to the value from config.Provider().LogoutRedirectURL()
		PostLogoutRedirectURI: redir,
	}); err != nil {
		return nil, err
	}

	http.Redirect(w, r, urlx.SetQuery(s.c.LogoutURL(), url.Values{"logout_challenge": {challenge}}).String(), http.StatusFound)
	return nil, errorsx.WithStack(ErrAbortOAuth2Request)
}

func (s *DefaultStrategy) completeLogout(w http.ResponseWriter, r *http.Request) (*LogoutResult, error) {
	verifier := r.URL.Query().Get("logout_verifier")

	lr, err := s.r.ConsentManager().VerifyAndInvalidateLogoutRequest(r.Context(), verifier)
	if err != nil {
		return nil, err
	}

	if !lr.RPInitiated {
		// If this is true it means that no id_token_hint was given, so the session id and subject id
		// came from an original cookie.

		session, err := s.authenticationSession(w, r)
		if errors.Is(err, ErrNoAuthenticationSessionFound) {
			// If we end up here it means that the cookie was revoked between the initial logout request
			// and ending up here - possibly due to a duplicate submit. In that case, we really have nothing to
			// do because the logout was already completed, apparently!

			// We also won't call any front- or back-channel logouts because that would mean we had called them twice!

			// OP initiated log out but no session was found. So let's just redirect back...
			http.Redirect(w, r, lr.PostLogoutRedirectURI, http.StatusFound)
			return nil, errorsx.WithStack(ErrAbortOAuth2Request)
		} else if err != nil {
			return nil, err
		}

		if session.Subject != lr.Subject {
			// If we end up here it means that the authentication cookie changed between the initial logout request
			// and landing here. That could happen because the user signed in in another browser window. In that
			// case there isn't really a lot to do because we don't want to sign out a different ID, so let's just
			// go to the post redirect uri without actually doing anything!
			http.Redirect(w, r, lr.PostLogoutRedirectURI, http.StatusFound)
			return nil, errorsx.WithStack(ErrAbortOAuth2Request)
		}
	}

	_, _ = s.revokeAuthenticationCookie(w, r, s.r.CookieStore()) // Cookie removal is optional

	urls, err := s.generateFrontChannelLogoutURLs(r.Context(), lr.Subject, lr.SessionID)
	if err != nil {
		return nil, err
	}

	if err := s.executeBackChannelLogout(r.Context(), r, lr.Subject, lr.SessionID); err != nil {
		return nil, err
	}

	// We delete the session after back channel log out has worked as the session is otherwise removed
	// from the store which will break the query for finding all the channels.
	//
	// executeBackChannelLogout only fails on system errors so not on URL errors, so this should be fine
	// even if an upstream URL fails!
	if err := s.r.ConsentManager().DeleteLoginSession(r.Context(), lr.SessionID); errors.Is(err, sqlcon.ErrNoRows) {
		// This is ok (session probably already revoked), do nothing!
	} else if err != nil {
		return nil, err
	}

	s.r.AuditLogger().
		WithRequest(r).
		WithField("subject", lr.Subject).
		Info("User logout completed!")

	return &LogoutResult{
		RedirectTo:             lr.PostLogoutRedirectURI,
		FrontChannelLogoutURLs: urls,
	}, nil
}

func (s *DefaultStrategy) HandleOpenIDConnectLogout(w http.ResponseWriter, r *http.Request) (*LogoutResult, error) {
	verifier := r.URL.Query().Get("logout_verifier")
	if verifier == "" {
		return s.issueLogoutVerifier(w, r)
	}

	return s.completeLogout(w, r)
}

func (s *DefaultStrategy) HandleOAuth2AuthorizationRequest(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*HandledConsentRequest, error) {
	authenticationVerifier := strings.TrimSpace(req.GetRequestForm().Get("login_verifier"))
	consentVerifier := strings.TrimSpace(req.GetRequestForm().Get("consent_verifier"))
	if authenticationVerifier == "" && consentVerifier == "" {
		// ok, we need to process this request and redirect to auth endpoint
		return nil, s.requestAuthentication(w, r, req)
	} else if authenticationVerifier != "" {
		authSession, err := s.verifyAuthentication(w, r, req, authenticationVerifier)
		if err != nil {
			return nil, err
		}

		// ok, we need to process this request and redirect to auth endpoint
		return nil, s.requestConsent(w, r, req, authSession)
	}

	consentSession, err := s.verifyConsent(w, r, req, consentVerifier)
	if err != nil {
		return nil, err
	}

	return consentSession, nil
}
