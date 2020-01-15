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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"

	"github.com/ory/x/httpx"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/x"
	"github.com/ory/x/mapx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/stringsx"
	"github.com/ory/x/urlx"
)

const (
	CookieAuthenticationName    = "oauth2_authentication_session"
	CookieAuthenticationSIDName = "sid"

	cookieAuthenticationCSRFName = "oauth2_authentication_csrf"
	cookieConsentCSRFName        = "oauth2_consent_csrf"
)

type DefaultStrategy struct {
	c Configuration
	r InternalRegistry
}

func NewStrategy(
	r InternalRegistry,
	c Configuration,
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
	if s, err := s.r.ConsentManager().GetForcedObfuscatedLoginSession(ctx, c.GetID(), hintSubject); errors.Cause(err) == x.ErrNotFound {
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
	cookie, err := s.r.CookieStore().Get(r, CookieAuthenticationName)
	if err != nil {
		return nil, errors.WithStack(ErrNoAuthenticationSessionFound)
	}

	sessionID := mapx.GetStringDefault(cookie.Values, CookieAuthenticationSIDName, "")
	if sessionID == "" {
		return nil, errors.WithStack(ErrNoAuthenticationSessionFound)
	}

	session, err := s.r.ConsentManager().GetRememberedLoginSession(r.Context(), sessionID)
	if errors.Cause(err) == x.ErrNotFound {
		return nil, errors.WithStack(ErrNoAuthenticationSessionFound)
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
	if errors.Cause(err) == ErrNoAuthenticationSessionFound {
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

	if maxAge > 0 && session.AuthenticatedAt.UTC().Add(time.Second*time.Duration(maxAge)).Before(time.Now().UTC()) {
		if stringslice.Has(prompt, "none") {
			return errors.WithStack(fosite.ErrLoginRequired.WithDebug("Request failed because prompt is set to \"none\" and authentication time reached max_age"))
		}
		return s.forwardAuthenticationRequest(w, r, ar, "", time.Time{}, nil)
	}

	idTokenHint := ar.GetRequestForm().Get("id_token_hint")
	if idTokenHint == "" {
		return s.forwardAuthenticationRequest(w, r, ar, session.Subject, session.AuthenticatedAt, session)
	}

	hintSub, err := s.getSubjectFromIDTokenHint(r.Context(), idTokenHint)
	if err != nil {
		return err
	}

	if err := s.matchesValueFromSession(r.Context(), ar.GetClient(), hintSub, session.Subject); errors.Cause(err) == ErrHintDoesNotMatchAuthentication {
		return errors.WithStack(fosite.ErrLoginRequired.WithDebug("Request failed because subject claim from id_token_hint does not match subject from authentication session"))
	}

	return s.forwardAuthenticationRequest(w, r, ar, session.Subject, session.AuthenticatedAt, session)
}

func (s *DefaultStrategy) getIDTokenHintClaims(ctx context.Context, idTokenHint string) (jwtgo.MapClaims, error) {
	token, err := s.r.OpenIDJWTStrategy().Decode(ctx, idTokenHint)
	if ve, ok := errors.Cause(err).(*jwtgo.ValidationError); errors.Cause(err) == nil || (ok && ve.Errors == jwtgo.ValidationErrorExpired) {
		// Expired is ok
	} else {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithHint(err.Error()))
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Failed to validate OpenID Connect request as decoding id token from id_token_hint to jwt.MapClaims failed"))
	}

	return claims, nil
}

func (s *DefaultStrategy) getSubjectFromIDTokenHint(ctx context.Context, idTokenHint string) (string, error) {
	claims, err := s.getIDTokenHintClaims(ctx, idTokenHint)
	if err != nil {
		return "", err
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Failed to validate OpenID Connect request because provided id token from id_token_hint does not have a subject"))
	}

	return sub, nil
}

func (s *DefaultStrategy) forwardAuthenticationRequest(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, subject string, authenticatedAt time.Time, session *LoginSession) error {
	if (subject != "" && authenticatedAt.IsZero()) || (subject == "" && !authenticatedAt.IsZero()) {
		return errors.WithStack(fosite.ErrServerError.WithDebug("Consent strategy returned a non-empty subject with an empty auth date, or an empty subject with a non-empty auth date"))
	}

	skip := false
	if subject != "" {
		skip = true
	}

	// Let'id validate that prompt is actually not "none" if we can't skip authentication
	prompt := stringsx.Splitx(ar.GetRequestForm().Get("prompt"), " ")
	if stringslice.Has(prompt, "none") && !skip {
		return errors.WithStack(fosite.ErrLoginRequired.WithDebug(`Prompt "none" was requested, but no existing login session was found`))
	}

	// Set up csrf/challenge/verifier values
	verifier := strings.Replace(uuid.New(), "-", "", -1)
	challenge := strings.Replace(uuid.New(), "-", "", -1)
	csrf := strings.Replace(uuid.New(), "-", "", -1)

	// Generate the request URL
	iu := urlx.AppendPaths(s.c.IssuerURL(), s.c.OAuth2AuthURL())
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
		if err := s.r.ConsentManager().CreateLoginSession(r.Context(), &LoginSession{
			ID:              sessionID,
			Subject:         "",
			AuthenticatedAt: time.Now().UTC(),
			Remember:        false,
		}); err != nil {
			return err
		}
	}

	// Set the session
	if err := s.r.ConsentManager().CreateLoginRequest(
		r.Context(),
		&LoginRequest{
			Challenge:         challenge,
			Verifier:          verifier,
			CSRF:              csrf,
			Skip:              skip,
			RequestedScope:    []string(ar.GetRequestedScopes()),
			RequestedAudience: []string(ar.GetRequestedAudience()),
			Subject:           subject,
			Client:            sanitizeClientFromRequest(ar),
			RequestURL:        iu.String(),
			AuthenticatedAt:   authenticatedAt,
			RequestedAt:       time.Now().UTC(),
			SessionID:         sessionID,
			OpenIDConnectContext: &OpenIDConnectContext{
				IDTokenHintClaims: idTokenHintClaims,
				ACRValues:         stringsx.Splitx(ar.GetRequestForm().Get("acr_values"), " "),
				UILocales:         stringsx.Splitx(ar.GetRequestForm().Get("ui_locales"), " "),
				Display:           ar.GetRequestForm().Get("display"),
				LoginHint:         ar.GetRequestForm().Get("login_hint"),
			},
		},
	); err != nil {
		return errors.WithStack(err)
	}

	if err := createCsrfSession(w, r, s.r.CookieStore(), cookieAuthenticationCSRFName, csrf, s.c.ServesHTTPS()); err != nil {
		return errors.WithStack(err)
	}

	http.Redirect(w, r, urlx.SetQuery(s.c.LoginURL(), url.Values{"login_challenge": {challenge}}).String(), http.StatusFound)

	// generate the verifier
	return errors.WithStack(ErrAbortOAuth2Request)
}

func (s *DefaultStrategy) revokeAuthenticationSession(w http.ResponseWriter, r *http.Request) error {
	sid, err := revokeAuthenticationCookie(w, r, s.r.CookieStore())
	if err != nil {
		return err
	}

	if sid == "" {
		return nil
	}

	return s.r.ConsentManager().DeleteLoginSession(r.Context(), sid)
}

func revokeAuthenticationCookie(w http.ResponseWriter, r *http.Request, s sessions.Store) (string, error) {
	cookie, _ := s.Get(r, CookieAuthenticationName)
	sid, _ := mapx.GetString(cookie.Values, CookieAuthenticationSIDName)

	cookie.Options.MaxAge = -1
	cookie.Values[CookieAuthenticationSIDName] = ""

	if err := cookie.Save(r, w); err != nil {
		return "", errors.WithStack(err)
	}

	return sid, nil
}

func (s *DefaultStrategy) obfuscateSubjectIdentifier(cl fosite.Client, subject, forcedIdentifier string) (string, error) {
	if c, ok := cl.(*client.Client); ok && c.SubjectType == "pairwise" {
		algorithm, ok := s.r.SubjectIdentifierAlgorithm()[c.SubjectType]
		if !ok {
			return "", errors.WithStack(fosite.ErrInvalidRequest.WithHint(fmt.Sprintf(`Subject Identifier Algorithm "%s" was requested by OAuth 2.0 Client "%s", but is not configured.`, c.SubjectType, c.ClientID)))
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
	if errors.Cause(err) == x.ErrNotFound {
		return nil, errors.WithStack(fosite.ErrAccessDenied.WithDebug("The login verifier has already been used, has not been granted, or is invalid."))
	} else if err != nil {
		return nil, err
	}

	if session.Error != nil {
		return nil, errors.WithStack(session.Error.toRFCError())
	}

	if session.RequestedAt.Add(s.c.ConsentRequestMaxAge()).Before(time.Now()) {
		return nil, errors.WithStack(fosite.ErrRequestUnauthorized.WithDebug("The login request has expired, please try again."))
	}

	if err := validateCsrfSession(r, s.r.CookieStore(), cookieAuthenticationCSRFName, session.LoginRequest.CSRF); err != nil {
		return nil, err
	}

	if session.LoginRequest.Skip && !session.Remember {
		return nil, errors.WithStack(fosite.ErrServerError.WithDebug("The login request was previously remembered and can only be forgotten using the reject feature."))
	}

	if session.LoginRequest.Skip && session.Subject != session.LoginRequest.Subject {
		// Revoke the session because there's clearly a mix up wrt the subject that's being authenticated
		if err := s.revokeAuthenticationSession(w, r); err != nil {
			return nil, err
		}

		return nil, errors.WithStack(fosite.ErrServerError.WithDebug("The login request is marked as remember, but the subject from the login confirmation does not match the original subject from the cookie."))
	}

	subjectIdentifier, err := s.obfuscateSubjectIdentifier(req.GetClient(), session.Subject, session.ForceSubjectIdentifier)
	if err != nil {
		return nil, err
	}

	sessionID := session.LoginRequest.SessionID

	if err := s.r.OpenIDConnectRequestValidator().ValidatePrompt(ctx, &fosite.AuthorizeRequest{
		ResponseTypes: req.GetResponseTypes(),
		RedirectURI:   req.GetRedirectURI(),
		State:         req.GetState(),
		// HandledResponseTypes, this can be safely ignored because it's not being used by validation
		Request: fosite.Request{
			ID:                req.GetID(),
			RequestedAt:       req.GetRequestedAt(),
			Client:            req.GetClient(),
			RequestedAudience: []string(req.GetRequestedAudience()),
			GrantedAudience:   []string(req.GetGrantedAudience()),
			RequestedScope:    req.GetRequestedScopes(),
			GrantedScope:      req.GetGrantedScopes(),
			Form:              req.GetRequestForm(),
			Session: &openid.DefaultSession{
				Claims: &jwt.IDTokenClaims{
					Subject:     subjectIdentifier,
					IssuedAt:    time.Now().UTC(),                // doesn't matter
					ExpiresAt:   time.Now().Add(time.Hour).UTC(), // doesn't matter
					AuthTime:    session.AuthenticatedAt,
					RequestedAt: session.RequestedAt,
				},
				Headers: &jwt.Headers{},
				Subject: session.Subject,
			},
		},
	}); errors.Cause(err) == fosite.ErrLoginRequired {
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
		if err := s.r.ConsentManager().ConfirmLoginSession(r.Context(), sessionID, session.Subject, session.Remember); err != nil {
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
	cookie, _ := s.r.CookieStore().Get(r, CookieAuthenticationName)
	cookie.Values[CookieAuthenticationSIDName] = sessionID
	if session.RememberFor >= 0 {
		cookie.Options.MaxAge = session.RememberFor
	}
	cookie.Options.HttpOnly = true

	if s.c.ServesHTTPS() {
		cookie.Options.Secure = true
	}

	if err := cookie.Save(r, w); err != nil {
		return nil, errors.WithStack(err)
	}
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
	if errors.Cause(err) == ErrNoPreviousConsentFound {
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
		return errors.WithStack(fosite.ErrConsentRequired.WithDebug(`Prompt "none" was requested, but no previous consent was found`))
	}

	// Set up csrf/challenge/verifier values
	verifier := strings.Replace(uuid.New(), "-", "", -1)
	challenge := strings.Replace(uuid.New(), "-", "", -1)
	csrf := strings.Replace(uuid.New(), "-", "", -1)

	if err := s.r.ConsentManager().CreateConsentRequest(
		r.Context(),
		&ConsentRequest{
			Challenge:              challenge,
			ACR:                    as.ACR,
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
			LoginChallenge:         as.LoginRequest.Challenge,
			Context:                as.Context,
		},
	); err != nil {
		return errors.WithStack(err)
	}

	if err := createCsrfSession(w, r, s.r.CookieStore(), cookieConsentCSRFName, csrf, s.c.ServesHTTPS()); err != nil {
		return errors.WithStack(err)
	}

	http.Redirect(
		w, r,
		urlx.SetQuery(s.c.ConsentURL(), url.Values{"consent_challenge": {challenge}}).String(),
		http.StatusFound,
	)

	// generate the verifier
	return errors.WithStack(ErrAbortOAuth2Request)
}

func (s *DefaultStrategy) verifyConsent(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester, verifier string) (*HandledConsentRequest, error) {
	session, err := s.r.ConsentManager().VerifyAndInvalidateConsentRequest(r.Context(), verifier)
	if errors.Cause(err) == x.ErrNotFound {
		return nil, errors.WithStack(fosite.ErrAccessDenied.WithDebug("The consent verifier has already been used, has not been granted, or is invalid."))
	} else if err != nil {
		return nil, err
	}

	if session.RequestedAt.Add(s.c.ConsentRequestMaxAge()).Before(time.Now()) {
		return nil, errors.WithStack(fosite.ErrRequestUnauthorized.WithDebug("The consent request has expired, please try again."))
	}

	if session.Error != nil {
		return nil, errors.WithStack(session.Error.toRFCError())
	}

	if session.ConsentRequest.AuthenticatedAt.IsZero() {
		return nil, errors.WithStack(fosite.ErrServerError.WithDebug("The authenticatedAt value was not set."))
	}

	if err := validateCsrfSession(r, s.r.CookieStore(), cookieConsentCSRFName, session.ConsentRequest.CSRF); err != nil {
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
			return nil, errors.WithStack(fosite.ErrServerError.WithHint(fmt.Sprintf("Unable to parse frontchannel_logout_uri: %s", c.FrontChannelLogoutURI)).WithDebug(err.Error()))
		}

		urls = append(urls, urlx.SetQuery(u, url.Values{
			"iss": {s.c.IssuerURL().String()},
			"sid": {sid},
		}).String())
	}

	return urls, nil
}

func (s *DefaultStrategy) executeBackChannelLogout(ctx context.Context, subject, sid string) error {
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
			"aud":    []string{c.ClientID},
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

		tasks = append(tasks, task{url: c.BackChannelLogoutURI, clientID: c.ClientID, token: t})
	}

	var wg sync.WaitGroup
	hc := http.Client{
		Timeout:   time.Second * 5,
		Transport: httpx.NewDefaultResilientRoundTripper(time.Second, time.Second*5),
	}
	wg.Add(len(tasks))

	var execute = func(t task) {
		defer wg.Done()

		res, err := hc.PostForm(t.url, url.Values{"logout_token": {t.token}})
		if err != nil {
			s.r.Logger().WithError(err).
				WithField("client_id", t.clientID).
				WithField("backchannel_logout_url", t.url).
				Warnf("Unable to execute OpenID Connect Back-Channel Logout Request")
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			s.r.Logger().WithError(errors.Errorf("expected HTTP status code %d but got %d", http.StatusOK, res.StatusCode)).
				WithField("client_id", t.clientID).
				WithField("backchannel_logout_url", t.url).
				Warnf("Unable to execute OpenID Connect Back-Channel Logout Request")
			return
		}

		return
	}

	for _, t := range tasks {
		go execute(t)
	}

	wg.Wait()

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

	// The hint must be set if it's an RP-initiated logout flow.
	hint := r.URL.Query().Get("id_token_hint")
	state := r.URL.Query().Get("state")
	requestedRedir := r.URL.Query().Get("post_logout_redirect_uri")

	if len(hint) == 0 {
		// hint is not set, so this is an OP initiated logout

		if len(state) > 0 {
			// state can only be set if it's an RP-initiated logout flow. If not, we should throw an error.
			return nil, errors.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter post_logout_redirect_uri is set but id_token_hint is missing"))
		}

		if len(requestedRedir) > 0 {
			// post_logout_redirect_uri can only be set if it's an RP-initiated logout flow. If not, we should throw an error.
			return nil, errors.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter post_logout_redirect_uri is set but id_token_hint is missing"))
		}

		session, err := s.authenticationSession(w, r)
		if errors.Cause(err) == ErrNoAuthenticationSessionFound {
			// OP initiated log out but no session was found. Since we can not identify the user we can not call
			// any RPs.
			http.Redirect(w, r, redir, http.StatusFound)
			return nil, errors.WithStack(ErrAbortOAuth2Request)
		} else if err != nil {
			return nil, err
		}

		challenge := uuid.New()
		if err := s.r.ConsentManager().CreateLogoutRequest(r.Context(), &LogoutRequest{
			RequestURL:  r.URL.String(),
			Challenge:   challenge,
			Subject:     session.Subject,
			SessionID:   session.ID,
			Verifier:    uuid.New(),
			RPInitiated: false,

			// PostLogoutRedirectURI is set to the value from config.Provider().LogoutRedirectURL()
			PostLogoutRedirectURI: redir,
		}); err != nil {
			return nil, err
		}

		http.Redirect(w, r, urlx.SetQuery(s.c.LogoutURL(), url.Values{"logout_challenge": {challenge}}).String(), http.StatusFound)
		return nil, errors.WithStack(ErrAbortOAuth2Request)
	}

	claims, err := s.getIDTokenHintClaims(r.Context(), hint)
	if err != nil {
		return nil, err
	}

	mksi := mapx.KeyStringToInterface(claims)
	if !claims.VerifyIssuer(s.c.IssuerURL().String(), true) {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.
			WithHint(
				fmt.Sprintf(
					`Logout failed because issuer claim value "%s" from query parameter id_token_hint does not match with issuer value from configuration "%s"`,
					mapx.GetStringDefault(mksi, "iss", ""),
					s.c.IssuerURL().String(),
				),
			),
		)
	}

	now := time.Now().UTC().Unix()
	if !claims.VerifyIssuedAt(now, true) {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.
			WithHint(
				fmt.Sprintf(
					`Logout failed because iat claim value "%.0f" from query parameter id_token_hint is before now ("%d")`,
					mapx.GetFloat64Default(mksi, "iat", float64(0)),
					now,
				),
			),
		)
	}

	hintSid := mapx.GetStringDefault(mksi, "sid", "")
	if len(hintSid) == 0 {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter id_token_hint is missing sid claim"))
	}

	// It doesn't really make sense to use the subject value from the ID Token because it might be obfuscated.
	if hintSub := mapx.GetStringDefault(mksi, "sub", ""); len(hintSub) == 0 {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithHint("Logout failed because query parameter id_token_hint is missing sub claim"))
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
		if errors.Cause(err) == x.ErrNotFound {
			continue
		} else if err != nil {
			return nil, err
		}
		cl = c
		break
	}

	if cl == nil {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.
			WithHint("Logout failed because none of the listed audiences is a registered OAuth 2.0 Client"))
	}

	if len(requestedRedir) > 0 {
		var f *url.URL
		for _, w := range cl.PostLogoutRedirectURIs {
			if w == requestedRedir {
				u, err := url.Parse(w)
				if err != nil {
					return nil, errors.WithStack(fosite.ErrServerError.WithHint(fmt.Sprintf("Unable to parse post_logout_redirect_uri: %s", w)).WithDebug(err.Error()))
				}

				f = u
			}
		}

		if f == nil {
			return nil, errors.WithStack(fosite.ErrInvalidRequest.
				WithHint("Logout failed because query parameter post_logout_redirect_uri is not a whitelisted as a post_logout_redirect_uri for the client"),
			)
		}

		redir = urlx.SetQuery(f, url.Values{
			"state": {r.URL.Query().Get("state")},
		}).String()
	}

	// We do not really want to verify if the user (from id token hint) has a session here because it doesn't really matter.
	// Instead, we'll check this when we're actually revoking the cookie!
	session, err := s.r.ConsentManager().GetRememberedLoginSession(r.Context(), hintSid)
	if errors.Cause(err) == x.ErrNotFound {
		// Such a session does not exist - maybe it has already been revoked? In any case, we can't do much except
		// leaning back and redirecting back.
		http.Redirect(w, r, redir, http.StatusFound)
		return nil, errors.WithStack(ErrAbortOAuth2Request)
	} else if err != nil {
		return nil, err
	}

	challenge := uuid.New()
	if err := s.r.ConsentManager().CreateLogoutRequest(r.Context(), &LogoutRequest{
		RequestURL:  r.URL.String(),
		Challenge:   challenge,
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
	return nil, errors.WithStack(ErrAbortOAuth2Request)
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
		if errors.Cause(err) == ErrNoAuthenticationSessionFound {
			// If we end up here it means that the cookie was revoked between the initial logout request
			// and ending up here - possibly due to a duplicate submit. In that case, we really have nothing to
			// do because the logout was already completed, apparently!

			// We also won't call any front- or back-channel logouts because that would mean we had called them twice!

			// OP initiated log out but no session was found. So let's just redirect back...
			http.Redirect(w, r, lr.PostLogoutRedirectURI, http.StatusFound)
			return nil, errors.WithStack(ErrAbortOAuth2Request)
		} else if err != nil {
			return nil, err
		}

		if session.Subject != lr.Subject {
			// If we end up here it means that the authentication cookie changed between the initial logout request
			// and landing here. That could happen because the user signed in in another browser window. In that
			// case there isn't really a lot to do because we don't want to sign out a different ID, so let's just
			// go to the post redirect uri without actually doing anything!
			http.Redirect(w, r, lr.PostLogoutRedirectURI, http.StatusFound)
			return nil, errors.WithStack(ErrAbortOAuth2Request)
		}
	}

	urls, err := s.generateFrontChannelLogoutURLs(r.Context(), lr.Subject, lr.SessionID)
	if err != nil {
		return nil, err
	}

	if err := s.executeBackChannelLogout(r.Context(), lr.Subject, lr.SessionID); err != nil {
		return nil, err
	}

	if err := s.revokeAuthenticationSession(w, r); err != nil {
		return nil, err
	}

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
