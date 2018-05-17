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
	"github.com/ory/go-convenience/urlx"
	"github.com/ory/herodot"
	"github.com/pkg/errors"
)

type Handler struct {
	H             herodot.Writer
	M             Manager
	RequestMaxAge time.Duration
}

func NewHandler(
	h herodot.Writer,
	m Manager,
) *Handler {
	return &Handler{
		H: h,
		M: m,
	}
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET("/oauth2/auth/requests/login/:challenge", h.GetLoginRequest)
	r.PUT("/oauth2/auth/requests/login/:challenge/accept", h.AcceptLoginRequest)
	r.PUT("/oauth2/auth/requests/login/:challenge/reject", h.RejectLoginRequest)

	r.GET("/oauth2/auth/requests/consent/:challenge", h.GetConsentRequest)
	r.PUT("/oauth2/auth/requests/consent/:challenge/accept", h.AcceptConsentRequest)
	r.PUT("/oauth2/auth/requests/consent/:challenge/reject", h.RejectConsentRequest)
}

// swagger:route GET /oauth2/auth/requests/login/{challenge} oAuth2 getLoginRequest
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
//       500: genericError
func (h *Handler) GetLoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request, err := h.M.GetAuthenticationRequest(ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/login/{challenge}/accept oAuth2 acceptLoginRequest
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
	p.RequestedAt = time.Now().UTC()

	ar, err := h.M.GetAuthenticationRequest(ps.ByName("challenge"))
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

	request, err := h.M.HandleAuthenticationRequest(ps.ByName("challenge"), &p)
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

// swagger:route PUT /oauth2/auth/requests/login/{challenge}/reject oAuth2 rejectLoginRequest
//
// Reject an logout request
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

	request, err := h.M.HandleAuthenticationRequest(ps.ByName("challenge"), &HandledAuthenticationRequest{
		Error:       &p,
		Challenge:   ps.ByName("challenge"),
		RequestedAt: time.Now().UTC(),
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

// swagger:route GET /oauth2/auth/requests/consent/{challenge} oAuth2 getConsentRequest
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
//       500: genericError
func (h *Handler) GetConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request, err := h.M.GetConsentRequest(ps.ByName("challenge"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.Write(w, r, request)
}

// swagger:route PUT /oauth2/auth/requests/consent/{challenge}/accept oAuth2 acceptConsentRequest
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

	p.Challenge = ps.ByName("challenge")
	p.RequestedAt = time.Now().UTC()
	hr, err := h.M.HandleConsentRequest(ps.ByName("challenge"), &p)
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

// swagger:route PUT /oauth2/auth/requests/consent/{challenge}/reject oAuth2 rejectConsentRequest
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

	request, err := h.M.HandleConsentRequest(ps.ByName("challenge"), &HandledConsentRequest{
		Error:       &p,
		Challenge:   ps.ByName("challenge"),
		RequestedAt: time.Now().UTC(),
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
