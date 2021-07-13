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
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/x"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/pagination"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/stringsx"
	"github.com/ory/x/urlx"
)

type Handler struct {
	r InternalRegistry
	c *config.Provider
}

const (
	LoginPath    = "/oauth2/auth/requests/login"
	ConsentPath  = "/oauth2/auth/requests/consent"
	LogoutPath   = "/oauth2/auth/requests/logout"
	SessionsPath = "/oauth2/auth/sessions"
)

func NewHandler(
	r InternalRegistry,
	c *config.Provider,
) *Handler {
	return &Handler{
		c: c,
		r: r,
	}
}

func (h *Handler) SetRoutes(admin *x.RouterAdmin) {
	admin.GET(LoginPath, h.GetLoginRequest)
	admin.PUT(LoginPath+"/accept", h.AcceptLoginRequest)
	admin.PUT(LoginPath+"/reject", h.RejectLoginRequest)

	admin.GET(ConsentPath, h.GetConsentRequest)
	admin.PUT(ConsentPath+"/accept", h.AcceptConsentRequest)
	admin.PUT(ConsentPath+"/reject", h.RejectConsentRequest)

	admin.DELETE(SessionsPath+"/login", h.DeleteLoginSession)
	admin.GET(SessionsPath+"/consent", h.GetConsentSessions)
	admin.DELETE(SessionsPath+"/consent", h.DeleteConsentSession)

	admin.GET(LogoutPath, h.GetLogoutRequest)
	admin.PUT(LogoutPath+"/accept", h.AcceptLogoutRequest)
	admin.PUT(LogoutPath+"/reject", h.RejectLogoutRequest)
}

// swagger:route DELETE /oauth2/auth/sessions/consent admin revokeConsentSessions
//
// Revokes Consent Sessions of a Subject for a Specific OAuth 2.0 Client
//
// This endpoint revokes a subject's granted consent sessions for a specific OAuth 2.0 Client and invalidates all
// associated OAuth 2.0 Access Tokens.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       400: jsonError
//       500: jsonError
func (h *Handler) DeleteConsentSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	client := r.URL.Query().Get("client")
	allClients := r.URL.Query().Get("all") == "true"
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	switch {
	case len(client) > 0:
		if err := h.r.ConsentManager().RevokeSubjectClientConsentSession(r.Context(), subject, client); err != nil && !errors.Is(err, x.ErrNotFound) {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	case allClients:
		if err := h.r.ConsentManager().RevokeSubjectConsentSession(r.Context(), subject); err != nil && !errors.Is(err, x.ErrNotFound) {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	default:
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter both 'client' and 'all' is not defined but one of them should have been.`)))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/sessions/consent admin listSubjectConsentSessions
//
// Lists All Consent Sessions of a Subject
//
// This endpoint lists all subject's granted consent sessions, including client and granted scope.
// If the subject is unknown or has not granted any consent sessions yet, the endpoint returns an
// empty JSON array with status code 200 OK.
//
//
// The "Link" header is also included in successful responses, which contains one or more links for pagination, formatted like so: '<https://hydra-url/admin/oauth2/auth/sessions/consent?subject={user}&limit={limit}&offset={offset}>; rel="{page}"', where page is one of the following applicable pages: 'first', 'next', 'last', and 'previous'.
// Multiple links can be included in this header, and will be separated by a comma.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: handledConsentRequestList
//       400: jsonError
//       500: jsonError
func (h *Handler) GetConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	limit, offset := pagination.Parse(r, 100, 0, 500)
	s, err := h.r.ConsentManager().FindSubjectsGrantedConsentRequests(r.Context(), subject, limit, offset)
	if errors.Is(err, ErrNoPreviousConsentFound) {
		h.r.Writer().Write(w, r, []PreviousConsentSession{})
		return
	} else if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	var a []PreviousConsentSession
	for _, session := range s {
		session.ConsentRequest.Client = sanitizeClient(session.ConsentRequest.Client)
		a = append(a, PreviousConsentSession(session))
	}

	if len(a) == 0 {
		a = []PreviousConsentSession{}
	}

	n, err := h.r.ConsentManager().CountSubjectsGrantedConsentRequests(r.Context(), subject)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	pagination.Header(w, r.URL, n, limit, offset)

	h.r.Writer().Write(w, r, a)
}

// swagger:route DELETE /oauth2/auth/sessions/login admin revokeAuthenticationSession
//
// Invalidates All Login Sessions of a Certain User
// Invalidates a Subject's Authentication Session
//
// This endpoint invalidates a subject's authentication session. After revoking the authentication session, the subject
// has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens and does not work with OpenID Connect
// Front- or Back-channel logout.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       400: jsonError
//       500: jsonError
func (h *Handler) DeleteLoginSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	if err := h.r.ConsentManager().RevokeSubjectLoginSession(r.Context(), subject); err != nil && !errors.Is(err, x.ErrNotFound) {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/requests/login admin getLoginRequest
//
// Get a Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the subject and then tell ORY Hydra now about it. The login
// provider is an web-app you write and host, and it must be able to authenticate ("show the subject a login screen")
// a subject (in OAuth2 the proper name for subject is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: loginRequest
//       400: jsonError
//       404: jsonError
//       410: requestWasHandledResponse
//       500: jsonError
func (h *Handler) GetLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("login_challenge"),
		r.URL.Query().Get("challenge"),
	)

	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	request, err := h.r.ConsentManager().GetLoginRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.r.Writer().WriteCode(w, r, http.StatusGone, &RequestWasHandledResponse{
			RedirectTo: request.RequestURL,
		})
		return
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/login/accept admin acceptLoginRequest
//
// Accept a Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the subject and then tell ORY Hydra now about it. The login
// provider is an web-app you write and host, and it must be able to authenticate ("show the subject a login screen")
// a subject (in OAuth2 the proper name for subject is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
// This endpoint tells ORY Hydra that the subject has successfully authenticated and includes additional information such as
// the subject's ID and if ORY Hydra should remember the subject's subject agent for future authentication attempts by setting
// a cookie.
//
// The response contains a redirect URL which the login provider should redirect the user-agent to.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: completedRequest
//       400: jsonError
//       404: jsonError
//       401: jsonError
//       500: jsonError
func (h *Handler) AcceptLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("login_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p HandledLoginRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithWrap(err).WithHintf("Unable to decode body because: %s", err)))
		return
	}

	if p.Subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Field 'subject' must not be empty.")))
		return
	}

	p.ID = challenge
	ar, err := h.r.ConsentManager().GetLoginRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	} else if ar.Subject != "" && p.Subject != ar.Subject {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Field 'subject' does not match subject from previous authentication.")))
		return
	}

	if ar.Skip {
		p.Remember = true // If skip is true remember is also true to allow consecutive calls as the same user!
		p.AuthenticatedAt = ar.AuthenticatedAt
	} else {
		p.AuthenticatedAt = sqlxx.NullTime(time.Now().UTC().
			// Rounding is important to avoid SQL time synchronization issues in e.g. MySQL!
			Truncate(time.Second))
		ar.AuthenticatedAt = p.AuthenticatedAt
	}
	p.RequestedAt = ar.RequestedAt

	request, err := h.r.ConsentManager().HandleLoginRequest(r.Context(), challenge, &p)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"login_verifier": {request.Verifier}}).String(),
	})
}

// swagger:route PUT /oauth2/auth/requests/login/reject admin rejectLoginRequest
//
// Reject a Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the subject and then tell ORY Hydra now about it. The login
// provider is an web-app you write and host, and it must be able to authenticate ("show the subject a login screen")
// a subject (in OAuth2 the proper name for subject is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
// This endpoint tells ORY Hydra that the subject has not authenticated and includes a reason why the authentication
// was be denied.
//
// The response contains a redirect URL which the login provider should redirect the user-agent to.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: completedRequest
//       400: jsonError
//       401: jsonError
//       404: jsonError
//       500: jsonError
func (h *Handler) RejectLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("login_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithWrap(err).WithHintf("Unable to decode body because: %s", err)))
		return
	}

	p.valid = true
	p.SetDefaults(loginRequestDeniedErrorName)
	ar, err := h.r.ConsentManager().GetLoginRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	request, err := h.r.ConsentManager().HandleLoginRequest(r.Context(), challenge, &HandledLoginRequest{
		Error:       &p,
		ID:          challenge,
		RequestedAt: ar.RequestedAt,
	})
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"login_verifier": {request.Verifier}}).String(),
	})
}

// swagger:route GET /oauth2/auth/requests/consent admin getConsentRequest
//
// Get Consent Request Information
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.
//
// The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to
// grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").
//
// The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted
// or rejected the request.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: consentRequest
//       404: jsonError
//       410: requestWasHandledResponse
//       500: jsonError
func (h *Handler) GetConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	request, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.r.Writer().WriteCode(w, r, http.StatusGone, &RequestWasHandledResponse{
			RedirectTo: request.RequestURL,
		})
		return
	}

	if request.RequestedScope == nil {
		request.RequestedScope = []string{}
	}

	if request.RequestedAudience == nil {
		request.RequestedAudience = []string{}
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/consent/accept admin acceptConsentRequest
//
// Accept a Consent Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.
//
// The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to
// grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").
//
// The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted
// or rejected the request.
//
// This endpoint tells ORY Hydra that the subject has authorized the OAuth 2.0 client to access resources on his/her behalf.
// The consent provider includes additional information, such as session data for access and ID tokens, and if the
// consent request should be used as basis for future requests.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: completedRequest
//       404: jsonError
//       500: jsonError
func (h *Handler) AcceptConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p HandledConsentRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errorsx.WithStack(err))
		return
	}

	cr, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	p.ID = challenge
	p.RequestedAt = cr.RequestedAt
	p.HandledAt = sqlxx.NullTime(time.Now().UTC())

	hr, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), challenge, &p)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	} else if hr.Skip {
		p.Remember = false
	}

	ru, err := url.Parse(hr.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"consent_verifier": {hr.Verifier}}).String(),
	})
}

// swagger:route PUT /oauth2/auth/requests/consent/reject admin rejectConsentRequest
//
// Reject a Consent Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.
//
// The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to
// grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").
//
// The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted
// or rejected the request.
//
// This endpoint tells ORY Hydra that the subject has not authorized the OAuth 2.0 client to access resources on his/her behalf.
// The consent provider must include a reason why the consent was not granted.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: completedRequest
//       404: jsonError
//       500: jsonError
func (h *Handler) RejectConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errorsx.WithStack(err))
		return
	}

	p.valid = true
	p.SetDefaults(consentRequestDeniedErrorName)
	hr, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	request, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), challenge, &HandledConsentRequest{
		Error:       &p,
		ID:          challenge,
		RequestedAt: hr.RequestedAt,
		HandledAt:   sqlxx.NullTime(time.Now().UTC()),
	})
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"consent_verifier": {request.Verifier}}).String(),
	})
}

// swagger:route PUT /oauth2/auth/requests/logout/accept admin acceptLogoutRequest
//
// Accept a Logout Request
//
// When a user or an application requests ORY Hydra to log out a user, this endpoint is used to confirm that logout request.
// No body is required.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: completedRequest
//       404: jsonError
//       500: jsonError
func (h *Handler) AcceptLogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("logout_challenge"),
		r.URL.Query().Get("challenge"),
	)

	c, err := h.r.ConsentManager().AcceptLogoutRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(urlx.AppendPaths(h.c.PublicURL(), "/oauth2/sessions/logout"), url.Values{"logout_verifier": {c.Verifier}}).String(),
	})
}

// swagger:route PUT /oauth2/auth/requests/logout/reject admin rejectLogoutRequest
//
// Reject a Logout Request
//
// When a user or an application requests ORY Hydra to log out a user, this endpoint is used to deny that logout request.
// No body is required.
//
// The response is empty as the logout provider has to chose what action to perform next.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       404: jsonError
//       500: jsonError
func (h *Handler) RejectLogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("logout_challenge"),
		r.URL.Query().Get("challenge"),
	)

	if err := h.r.ConsentManager().RejectLogoutRequest(r.Context(), challenge); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/requests/logout admin getLogoutRequest
//
// Get a Logout Request
//
// Use this endpoint to fetch a logout request.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: logoutRequest
//       404: jsonError
//       410: requestWasHandledResponse
//       500: jsonError
func (h *Handler) GetLogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("logout_challenge"),
		r.URL.Query().Get("challenge"),
	)

	request, err := h.r.ConsentManager().GetLogoutRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	// We do not want to share the secret so remove it.
	if request.Client != nil {
		request.Client.Secret = ""
	}

	if request.WasHandled {
		h.r.Writer().WriteCode(w, r, http.StatusGone, &RequestWasHandledResponse{
			RedirectTo: request.RequestURL,
		})
		return
	}

	h.r.Writer().Write(w, r, request)
}
