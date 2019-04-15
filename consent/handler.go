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
	"github.com/ory/go-convenience/urlx"
	"github.com/ory/hydra/x"
	"github.com/ory/x/pagination"
)

type Handler struct {
	r InternalRegistry
	c Configuration
}

const (
	LoginPath    = "/oauth2/auth/requests/login"
	ConsentPath  = "/oauth2/auth/requests/consent"
	SessionsPath = "/oauth2/auth/sessions"
)

func NewHandler(
	r InternalRegistry,
	c Configuration,
) *Handler {
	return &Handler{
		c: c,
		r: r,
	}
}

func (h *Handler) SetRoutes(admin *x.RouterAdmin, public *x.RouterPublic) {
	admin.GET(LoginPath, h.GetLoginRequest)
	admin.PUT(LoginPath+"/accept", h.AcceptLoginRequest)
	admin.PUT(LoginPath+"/reject", h.RejectLoginRequest)

	admin.GET(ConsentPath, h.GetConsentRequest)
	admin.PUT(ConsentPath+"/accept", h.AcceptConsentRequest)
	admin.PUT(ConsentPath+"/reject", h.RejectConsentRequest)

	admin.DELETE(SessionsPath+"/login/:user", h.DeleteLoginSession)
	admin.GET(SessionsPath+"/consent/:user", h.GetConsentSessions)
	admin.DELETE(SessionsPath+"/consent/:user", h.DeleteUserConsentSession)
	admin.DELETE(SessionsPath+"/consent/:user/:client", h.DeleteUserClientConsentSession)

	public.GET(SessionsPath+"/login/revoke", h.LogoutUser)
}

// swagger:route DELETE /oauth2/auth/sessions/consent/{user} admin revokeAllUserConsentSessions
//
// Revokes all previous consent sessions of a user
//
// This endpoint revokes a user's granted consent sessions and invalidates all associated OAuth 2.0 Access Tokens.
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
//       404: genericError
//       500: genericError
func (h *Handler) DeleteUserConsentSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	if err := h.r.ConsentManager().RevokeUserConsentSession(r.Context(), user); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /oauth2/auth/sessions/consent/{user}/{client} admin revokeUserClientConsentSessions
//
// Revokes consent sessions of a user for a specific OAuth 2.0 Client
//
// This endpoint revokes a user's granted consent sessions for a specific OAuth 2.0 Client and invalidates all
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
//       404: genericError
//       500: genericError
func (h *Handler) DeleteUserClientConsentSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client := ps.ByName("client")
	user := ps.ByName("user")
	if client == "" {
		h.r.Writer().WriteError(w, r, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Parameter client is not defined")))
		return
	}

	if err := h.r.ConsentManager().RevokeUserClientConsentSession(r.Context(), user, client); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/sessions/consent/{user} admin listUserConsentSessions
//
// Lists all consent sessions of a user
//
// This endpoint lists all user's granted consent sessions, including client and granted scope.
// The "Link" header is also included in successful responses, which contains one or more links for pagination, formatted like so: '<https://hydra-url/admin/oauth2/auth/sessions/consent/{user}?limit={limit}&offset={offset}>; rel="{page}"', where page is one of the following applicable pages: 'first', 'next', 'last', and 'previous'.
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
//       404: genericError
//       500: genericError
func (h *Handler) GetConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	if user == "" {
		h.r.Writer().WriteError(w, r, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Parameter user is not defined")))
		return
	}

	limit, offset := pagination.Parse(r, 100, 0, 500)
	s, err := h.r.ConsentManager().FindSubjectsGrantedConsentRequests(r.Context(), user, limit, offset)
	if errors.Cause(err) == ErrNoPreviousConsentFound {
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

	n, err := h.r.ConsentManager().CountSubjectsGrantedConsentRequests(r.Context(), user)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, a)
	pagination.Header(r.URL, n, limit, offset).Write(w)
}

// swagger:route DELETE /oauth2/auth/sessions/login/{user} admin revokeAuthenticationSession
//
// Invalidates a user's authentication session
//
// This endpoint invalidates a user's authentication session. After revoking the authentication session, the user
// has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens.
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
//       404: genericError
//       500: genericError
func (h *Handler) DeleteLoginSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")

	if err := h.r.ConsentManager().RevokeUserAuthenticationSession(r.Context(), user); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/requests/login admin getLoginRequest
//
// Get an login request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the user and then tell ORY Hydra now about it. The login
// provider is an web-app you write and host, and it must be able to authenticate ("show the user a login screen")
// a user (in OAuth2 the proper name for user is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login
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
//       404: genericError
//       409: genericError
//       500: genericError
func (h *Handler) GetLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := r.URL.Query().Get("challenge")
	request, err := h.r.ConsentManager().GetAuthenticationRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.r.Writer().WriteError(w, r, x.ErrConflict.WithDebug("Login request has been handled already"))
		return
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/login/accept admin acceptLoginRequest
//
// Accept an login request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the user and then tell ORY Hydra now about it. The login
// provider is an web-app you write and host, and it must be able to authenticate ("show the user a login screen")
// a user (in OAuth2 the proper name for user is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
// This endpoint tells ORY Hydra that the user has successfully authenticated and includes additional information such as
// the user's ID and if ORY Hydra should remember the user's user agent for future authentication attempts by setting
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
//       404: genericError
//       401: genericError
//       500: genericError
func (h *Handler) AcceptLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := r.URL.Query().Get("challenge")

	var p HandledLoginRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}
	if p.Subject == "" {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Subject from payload can not be empty"))
	}

	p.Challenge = challenge
	ar, err := h.r.ConsentManager().GetAuthenticationRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	} else if ar.Subject != "" && p.Subject != ar.Subject {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Subject from payload does not match subject from previous authentication"))
		return
	}

	if !ar.Skip {
		p.AuthenticatedAt = time.Now().UTC()
	} else {
		p.Remember = false
		p.AuthenticatedAt = ar.AuthenticatedAt
	}
	p.RequestedAt = ar.RequestedAt

	request, err := h.r.ConsentManager().HandleAuthenticationRequest(r.Context(), challenge, &p)
	if err != nil {
		h.r.Writer().WriteError(w, r, errors.WithStack(err))
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
// Reject a login request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the user and then tell ORY Hydra now about it. The login
// provider is an web-app you write and host, and it must be able to authenticate ("show the user a login screen")
// a user (in OAuth2 the proper name for user is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
// This endpoint tells ORY Hydra that the user has not authenticated and includes a reason why the authentication
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
//       401: genericError
//       404: genericError
//       500: genericError
func (h *Handler) RejectLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := r.URL.Query().Get("challenge")

	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	ar, err := h.r.ConsentManager().GetAuthenticationRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	request, err := h.r.ConsentManager().HandleAuthenticationRequest(r.Context(), challenge, &HandledLoginRequest{
		Error:       &p,
		Challenge:   challenge,
		RequestedAt: ar.RequestedAt,
	})
	if err != nil {
		h.r.Writer().WriteError(w, r, errors.WithStack(err))
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
// Get consent request information
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.
//
// The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to
// grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").
//
// The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted
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
//       404: genericError
//       409: genericError
//       500: genericError
func (h *Handler) GetConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := r.URL.Query().Get("challenge")

	request, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.r.Writer().WriteError(w, r, x.ErrConflict.WithDebug("Consent request has been handled already"))
		return
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/consent/accept admin acceptConsentRequest
//
// Accept an consent request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.
//
// The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to
// grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").
//
// The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted
// or rejected the request.
//
// This endpoint tells ORY Hydra that the user has authorized the OAuth 2.0 client to access resources on his/her behalf.
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
//       404: genericError
//       500: genericError
func (h *Handler) AcceptConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := r.URL.Query().Get("challenge")

	var p HandledConsentRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	cr, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, errors.WithStack(err))
		return
	}

	p.Challenge = challenge
	p.RequestedAt = cr.RequestedAt

	hr, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), challenge, &p)
	if err != nil {
		h.r.Writer().WriteError(w, r, errors.WithStack(err))
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
// Reject an consent request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
// to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.
//
// The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to
// grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").
//
// The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted
// or rejected the request.
//
// This endpoint tells ORY Hydra that the user has not authorized the OAuth 2.0 client to access resources on his/her behalf.
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
//       404: genericError
//       500: genericError
func (h *Handler) RejectConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := r.URL.Query().Get("challenge")

	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	hr, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, errors.WithStack(err))
		return
	}

	request, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), challenge, &HandledConsentRequest{
		Error:       &p,
		Challenge:   challenge,
		RequestedAt: hr.RequestedAt,
	})
	if err != nil {
		h.r.Writer().WriteError(w, r, errors.WithStack(err))
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

// swagger:route GET /oauth2/auth/sessions/login/revoke public revokeUserLoginCookie
//
// Logs user out by deleting the session cookie
//
// This endpoint deletes ths user's login session cookie and redirects the browser to the url
// listed in `LOGOUT_REDIRECT_URL` environment variable. This endpoint does not work as an API but has to
// be called from the user's browser.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       302: emptyResponse
//       404: genericError
//       500: genericError
func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sid, err := revokeAuthenticationCookie(w, r, h.r.CookieStore())
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	if sid != "" {
		if err := h.r.ConsentManager().DeleteAuthenticationSession(r.Context(), sid); err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, h.c.LogoutRedirectURL().String(), 302)
}
