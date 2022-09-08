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

	"github.com/ory/x/httprouterx"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/x"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/stringsx"
	"github.com/ory/x/urlx"
)

type Handler struct {
	r InternalRegistry
	c *config.DefaultProvider
}

const (
	LoginPath    = "/oauth2/auth/requests/login"
	ConsentPath  = "/oauth2/auth/requests/consent"
	LogoutPath   = "/oauth2/auth/requests/logout"
	SessionsPath = "/oauth2/auth/sessions"
)

func NewHandler(
	r InternalRegistry,
	c *config.DefaultProvider,
) *Handler {
	return &Handler{
		c: c,
		r: r,
	}
}

func (h *Handler) SetRoutes(admin *httprouterx.RouterAdmin) {
	admin.GET(LoginPath, h.adminGetOAuth2LoginRequest)
	admin.PUT(LoginPath+"/accept", h.adminAcceptOAuth2LoginRequest)
	admin.PUT(LoginPath+"/reject", h.adminRejectOAuth2LoginRequest)

	admin.GET(ConsentPath, h.adminGetOAuth2ConsentRequest)
	admin.PUT(ConsentPath+"/accept", h.adminAcceptOAuth2ConsentRequest)
	admin.PUT(ConsentPath+"/reject", h.adminRejectOAuth2ConsentRequest)

	admin.DELETE(SessionsPath+"/login", h.adminRevokeOAuth2LoginSessions)
	admin.GET(SessionsPath+"/consent", h.adminListOAuth2SubjectConsentSessions)
	admin.DELETE(SessionsPath+"/consent", h.adminRevokeOAuth2ConsentSessions)

	admin.GET(LogoutPath, h.adminGetOAuth2LogoutRequest)
	admin.PUT(LogoutPath+"/accept", h.adminAcceptOAuth2LogoutRequest)
	admin.PUT(LogoutPath+"/reject", h.adminRejectOAuth2LogoutRequest)
}

// swagger:parameters adminRevokeOAuth2ConsentSessions
type adminRevokeOAuth2ConsentSessions struct {
	// The subject (Subject) whose consent sessions should be deleted.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`

	// If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID
	//
	// in: query
	Client string `json:"client"`

	// If set to `true` deletes all consent sessions by the Subject that have been granted.
	//
	// in: query
	All bool `json:"all"`
}

// swagger:route DELETE /admin/oauth2/auth/sessions/consent v0alpha2 adminRevokeOAuth2ConsentSessions
//
// Revokes OAuth 2.0 Consent Sessions of a Subject for a Specific OAuth 2.0 Client
//
// This endpoint revokes a subject's granted consent sessions for a specific OAuth 2.0 Client and invalidates all
// associated OAuth 2.0 Access Tokens.
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
//       default: oAuth2ApiError
func (h *Handler) adminRevokeOAuth2ConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// swagger:parameters adminListOAuth2SubjectConsentSessions
type adminListOAuth2SubjectConsentSessions struct {
	x.PaginationHeaders

	// The subject to list the consent sessions for.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`
}

// swagger:route GET /admin/oauth2/auth/sessions/consent v0alpha2 adminListOAuth2SubjectConsentSessions
//
// List OAuth 2.0 Consent Sessions of a Subject
//
// This endpoint lists all subject's granted consent sessions, including client and granted scope.
// If the subject is unknown or has not granted any consent sessions yet, the endpoint returns an
// empty JSON array with status code 200 OK.
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
//       200: previousOAuth2ConsentSessions
//       default: oAuth2ApiError
func (h *Handler) adminListOAuth2SubjectConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	page, itemsPerPage := x.ParsePagination(r)
	s, err := h.r.ConsentManager().FindSubjectsGrantedConsentRequests(r.Context(), subject, itemsPerPage, itemsPerPage*page)
	if errors.Is(err, ErrNoPreviousConsentFound) {
		h.r.Writer().Write(w, r, []PreviousOAuth2ConsentSession{})
		return
	} else if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	var a []PreviousOAuth2ConsentSession
	for _, session := range s {
		session.ConsentRequest.Client = sanitizeClient(session.ConsentRequest.Client)
		a = append(a, PreviousOAuth2ConsentSession(session))
	}

	if len(a) == 0 {
		a = []PreviousOAuth2ConsentSession{}
	}

	n, err := h.r.ConsentManager().CountSubjectsGrantedConsentRequests(r.Context(), subject)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	x.PaginationHeader(w, r.URL, int64(n), itemsPerPage, itemsPerPage*page)
	h.r.Writer().Write(w, r, a)
}

// swagger:parameters adminRevokeOAuth2LoginSessions
type adminRevokeOAuth2LoginSessions struct {
	// The subject to revoke authentication sessions for.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`
}

// swagger:route DELETE /admin/oauth2/auth/sessions/login v0alpha2 adminRevokeOAuth2LoginSessions
//
// Invalidates All OAuth 2.0 Login Sessions of a Certain User
//
// This endpoint invalidates a subject's authentication session. After revoking the authentication session, the subject
// has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens and does not work with OpenID Connect
// Front- or Back-channel logout.
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
//       default: oAuth2ApiError
func (h *Handler) adminRevokeOAuth2LoginSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	if err := h.r.ConsentManager().RevokeSubjectLoginSession(r.Context(), subject); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:parameters adminGetOAuth2LoginRequest
type adminGetOAuth2LoginRequest struct {
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`
}

// swagger:route GET /admin/oauth2/auth/requests/login v0alpha2 adminGetOAuth2LoginRequest
//
// Get an OAuth 2.0 Login Request
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
//       200: oAuth2LoginRequest
//       410: handledOAuth2LoginRequest
//       default: oAuth2ApiError
func (h *Handler) adminGetOAuth2LoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		h.r.Writer().WriteCode(w, r, http.StatusGone, &RequestHandlerResponse{
			RedirectTo: request.RequestURL,
		})
		return
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// swagger:parameters adminAcceptOAuth2LoginRequest
type adminAcceptOAuth2LoginRequest struct {
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`

	// in: body
	Body HandledLoginRequest
}

// swagger:route PUT /admin/oauth2/auth/requests/login/accept v0alpha2 adminAcceptOAuth2LoginRequest
//
// Accept an OAuth 2.0 Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory Hydra asks the login provider
// (sometimes called "identity provider") to authenticate the subject and then tell Ory Hydra now about it. The login
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
//       200: successfulOAuth2RequestResponse
//       default: oAuth2ApiError
func (h *Handler) adminAcceptOAuth2LoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// swagger:parameters adminRejectOAuth2LoginRequest
type adminRejectOAuth2LoginRequest struct {
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:route PUT /admin/oauth2/auth/requests/login/reject v0alpha2 adminRejectOAuth2LoginRequest
//
// Reject an OAuth 2.0 Login Request
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
// was denied.
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
//       200: successfulOAuth2RequestResponse
//       default: oAuth2ApiError
func (h *Handler) adminRejectOAuth2LoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// swagger:parameters adminGetOAuth2ConsentRequest
type adminGetOAuth2ConsentRequest struct {
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`
}

// swagger:route GET /admin/oauth2/auth/requests/consent v0alpha2 adminGetOAuth2ConsentRequest
//
// Get OAuth 2.0 Consent Request Information
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
//       200: oAuth2ConsentRequest
//       410: handledOAuth2ConsentRequest
//       default: oAuth2ApiError
func (h *Handler) adminGetOAuth2ConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		h.r.Writer().WriteCode(w, r, http.StatusGone, &HandledOAuth2ConsentRequest{
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

// swagger:parameters adminAcceptOAuth2ConsentRequest
type adminAcceptOAuth2ConsentRequest struct {
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`

	// in: body
	Body AcceptOAuth2ConsentRequest
}

// swagger:route PUT /admin/oauth2/auth/requests/consent/accept v0alpha2 adminAcceptOAuth2ConsentRequest
//
// Accept an OAuth 2.0 Consent Request
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
//       200: successfulOAuth2RequestResponse
//       default: oAuth2ApiError
func (h *Handler) adminAcceptOAuth2ConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p AcceptOAuth2ConsentRequest
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

	hr, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), &p)
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

// swagger:parameters adminRejectOAuth2ConsentRequest
type adminRejectOAuth2ConsentRequest struct {
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:route PUT /admin/oauth2/auth/requests/consent/reject v0alpha2 adminRejectOAuth2ConsentRequest
//
// Reject an OAuth 2.0 Consent Request
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
//       200: successfulOAuth2RequestResponse
//       default: oAuth2ApiError
func (h *Handler) adminRejectOAuth2ConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	request, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), &AcceptOAuth2ConsentRequest{
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

// swagger:parameters adminAcceptOAuth2LogoutRequest
type adminAcceptOAuth2LogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:route PUT /admin/oauth2/auth/requests/logout/accept v0alpha2 adminAcceptOAuth2LogoutRequest
//
// Accept an OAuth 2.0 Logout Request
//
// When a user or an application requests ORY Hydra to log out a user, this endpoint is used to confirm that logout request.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: successfulOAuth2RequestResponse
//       default: oAuth2ApiError
func (h *Handler) adminAcceptOAuth2LogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		RedirectTo: urlx.SetQuery(urlx.AppendPaths(h.c.PublicURL(r.Context()), "/oauth2/sessions/logout"), url.Values{"logout_verifier": {c.Verifier}}).String(),
	})
}

// swagger:parameters adminRejectOAuth2LogoutRequest
type adminRejectOAuth2LogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:route PUT /admin/oauth2/auth/requests/logout/reject v0alpha2 adminRejectOAuth2LogoutRequest
//
// Reject an OAuth 2.0 Logout Request
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
//       default: oAuth2ApiError
func (h *Handler) adminRejectOAuth2LogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// swagger:parameters adminGetOAuth2LogoutRequest
type adminGetOAuth2LogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:route GET /admin/oauth2/auth/requests/logout v0alpha2 adminGetOAuth2LogoutRequest
//
// Get an OAuth 2.0 Logout Request
//
// Use this endpoint to fetch a logout request.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2LogoutRequest
//       410: handledOAuth2LogoutRequest
//       default: oAuth2ApiError
func (h *Handler) adminGetOAuth2LogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		h.r.Writer().WriteCode(w, r, http.StatusGone, &HandledOAuth2ConsentRequest{
			RedirectTo: request.RequestURL,
		})
		return
	}

	h.r.Writer().Write(w, r, request)
}
