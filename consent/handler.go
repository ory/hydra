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

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/go-convenience/urlx"
	"github.com/ory/herodot"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/pagination"
)

type Handler struct {
	H                 herodot.Writer
	M                 Manager
	LogoutRedirectURL string
	RequestMaxAge     time.Duration
	CookieStore       sessions.Store
}

const (
	LoginPath    = "/oauth2/auth/requests/login"
	ConsentPath  = "/oauth2/auth/requests/consent"
	SessionsPath = "/oauth2/auth/sessions"
)

func NewHandler(
	h herodot.Writer,
	m Manager,
	c sessions.Store,
	u string,
) *Handler {
	return &Handler{
		H:                 h,
		M:                 m,
		LogoutRedirectURL: u,
		CookieStore:       c,
	}
}

func (h *Handler) SetRoutes(frontend, backend *httprouter.Router) {
	backend.GET(LoginPath+"/:challenge", h.GetLoginRequest)
	backend.PUT(LoginPath+"/:challenge/accept", h.AcceptLoginRequest)
	backend.PUT(LoginPath+"/:challenge/reject", h.RejectLoginRequest)

	backend.GET(ConsentPath+"/:challenge", h.GetConsentRequest)
	backend.PUT(ConsentPath+"/:challenge/accept", h.AcceptConsentRequest)
	backend.PUT(ConsentPath+"/:challenge/reject", h.RejectConsentRequest)

	backend.DELETE(SessionsPath+"/login/:user", h.DeleteLoginSession)
	backend.GET(SessionsPath+"/consent/:user", h.GetConsentSessions)
	backend.DELETE(SessionsPath+"/consent/:user", h.DeleteUserConsentSession)
	backend.DELETE(SessionsPath+"/consent/:user/:client", h.DeleteUserClientConsentSession)

	frontend.GET(SessionsPath+"/login/revoke", h.LogoutUser)
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
	if err := h.M.RevokeUserConsentSession(r.Context(), user); err != nil {
		h.H.WriteError(w, r, err)
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
		h.H.WriteError(w, r, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Parameter client is not defined")))
		return
	}

	if err := h.M.RevokeUserClientConsentSession(r.Context(), user, client); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/sessions/consent/{user} admin listUserConsentSessions
//
// Lists all consent sessions of a user
//
// This endpoint lists all user's granted consent sessions, including client and granted scope
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
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) GetConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("user")
	if user == "" {
		h.H.WriteError(w, r, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Parameter user is not defined")))
		return
	}

	limit, offset := pagination.Parse(r, 100, 0, 500)
	s, err := h.M.FindPreviouslyGrantedConsentRequestsByUser(r.Context(), user, limit, offset)
	if errors.Cause(err) == ErrNoPreviousConsentFound {
		h.H.Write(w, r, []PreviousConsentSession{})
		return
	} else if err != nil {
		h.H.WriteError(w, r, err)
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

	h.H.Write(w, r, a)
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

	if err := h.M.RevokeUserAuthenticationSession(r.Context(), user); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /oauth2/auth/requests/login/{challenge} admin getLoginRequest
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
//       401: genericError
//       409: genericError
//       500: genericError
func (h *Handler) GetLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request, err := h.M.GetAuthenticationRequest(r.Context(), ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.H.WriteError(w, r, pkg.ErrConflict.WithDebug("Login request has been handled already"))
		return
	}

	request.Client = sanitizeClient(request.Client)

	h.H.Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/login/{challenge}/accept admin acceptLoginRequest
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
//       401: genericError
//       500: genericError
func (h *Handler) AcceptLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p HandledAuthenticationRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	p.Challenge = ps.ByName("challenge")
	ar, err := h.M.GetAuthenticationRequest(r.Context(), ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	} else if ar.Subject != "" && p.Subject != ar.Subject {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Subject from payload does not match subject from previous authentication"))
		return
	} else if ar.Skip && p.Remember {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Can not remember authentication because no user interaction was required"))
		return
	}

	if !ar.Skip {
		p.AuthenticatedAt = time.Now().UTC()
	} else {
		p.AuthenticatedAt = ar.AuthenticatedAt
	}
	p.RequestedAt = ar.RequestedAt

	request, err := h.M.HandleAuthenticationRequest(r.Context(), ps.ByName("challenge"), &p)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"login_verifier": {request.Verifier}}).String(),
	})
}

// swagger:route PUT /oauth2/auth/requests/login/{challenge}/reject admin rejectLoginRequest
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
//       500: genericError
func (h *Handler) RejectLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	ar, err := h.M.GetAuthenticationRequest(r.Context(), ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	request, err := h.M.HandleAuthenticationRequest(r.Context(), ps.ByName("challenge"), &HandledAuthenticationRequest{
		Error:       &p,
		Challenge:   ps.ByName("challenge"),
		RequestedAt: ar.RequestedAt,
	})
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"login_verifier": {request.Verifier}}).String(),
	})
}

// swagger:route GET /oauth2/auth/requests/consent/{challenge} admin getConsentRequest
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
//       401: genericError
//       409: genericError
//       500: genericError
func (h *Handler) GetConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request, err := h.M.GetConsentRequest(r.Context(), ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.H.WriteError(w, r, pkg.ErrConflict.WithDebug("Consent request has been handled already"))
		return
	}

	request.Client = sanitizeClient(request.Client)

	h.H.Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/consent/{challenge}/accept admin acceptConsentRequest
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
//       401: genericError
//       500: genericError
func (h *Handler) AcceptConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p HandledConsentRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	cr, err := h.M.GetConsentRequest(r.Context(), ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	p.Challenge = ps.ByName("challenge")
	p.RequestedAt = cr.RequestedAt

	hr, err := h.M.HandleConsentRequest(r.Context(), ps.ByName("challenge"), &p)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	} else if hr.Skip && p.Remember {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Can not remember consent because no user interaction was required"))
		return
	}

	ru, err := url.Parse(hr.RequestURL)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"consent_verifier": {hr.Verifier}}).String(),
	})
}

// swagger:route PUT /oauth2/auth/requests/consent/{challenge}/reject admin rejectConsentRequest
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
//       401: genericError
//       500: genericError
func (h *Handler) RejectConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	hr, err := h.M.GetConsentRequest(r.Context(), ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	request, err := h.M.HandleConsentRequest(r.Context(), ps.ByName("challenge"), &HandledConsentRequest{
		Error:       &p,
		Challenge:   ps.ByName("challenge"),
		RequestedAt: hr.RequestedAt,
	})
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, &RequestHandlerResponse{
		RedirectTo: urlx.SetQuery(ru, url.Values{"consent_verifier": {request.Verifier}}).String(),
	})
}

// swagger:route GET /oauth2/auth/sessions/login/revoke admin revokeUserLoginCookie
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
	sid, err := revokeAuthenticationCookie(w, r, h.CookieStore)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if sid != "" {
		if err := h.M.DeleteAuthenticationSession(r.Context(), sid); err != nil {
			h.H.WriteError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, h.LogoutRedirectURL, 302)
}
