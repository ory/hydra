package oauth2

import (
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/enigma/jwt"
	"github.com/ory-am/hydra/identity"
	"github.com/ory-am/hydra/pkg"
)

type Handler struct {
	OAuth2     fosite.OAuth2Provider
	Consent    ConsentStrategy

	SelfURL    *url.URL
	ConsentURL *url.URL
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/oauth2/token", h.TokenHandler)

	r.GET("/oauth2/auth", h.AuthHandler)
	r.POST("/oauth2/auth", h.AuthHandler)
}

func (o *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var session Session
	var ctx = fosite.NewContext()

	accessRequest, err := o.OAuth2.NewAccessRequest(ctx, r, &session)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.OAuth2.WriteAccessError(w, accessRequest, err)
		return
	}

	accessResponse, err := o.OAuth2.NewAccessResponse(ctx, r, accessRequest)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.OAuth2.WriteAccessError(w, accessRequest, err)
		return
	}

	o.OAuth2.WriteAccessResponse(w, accessRequest, accessResponse)
}

func (o *Handler) AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = fosite.NewContext()

	authorizeRequest, err := o.OAuth2.NewAuthorizeRequest(ctx, r)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	// A session_token will be available if the user was authenticated an gave consent
	consentToken := authorizeRequest.GetRequestForm().Get("consent_token")
	if consentToken == "" {
		// otherwise redirect to log in endpoint
		o.redirectToConsent(w, r, authorizeRequest, "")
		return
	}

	// decode consent_token claims
	// verify anti-CSRF (inject state) and anti-replay token (expiry time, good value would be 10 seconds)
	session, err := o.Consent.ValidateResponseToken(authorizeRequest, consentToken)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.writeAuthorizeError(w, authorizeRequest, errors.New(fosite.ErrAccessDenied))
		return
	}

	// done
	response, err := o.OAuth2.NewAuthorizeResponse(ctx, r, authorizeRequest, session)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	o.OAuth2.WriteAuthorizeResponse(w, authorizeRequest, response)
}

func (o *Handler) redirectToConsent(w http.ResponseWriter, r *http.Request, authorizeRequest fosite.AuthorizeRequester, authenticationToken string) error {
	var p = new(url.URL)
	*p = *o.ConsentURL

	q := p.Query()
	q.Set("client_id", authorizeRequest.GetClient().GetID())
	q.Set("state", authorizeRequest.GetState())
	for _, scope := range authorizeRequest.GetScopes() {
		q.Add("scope", scope)
	}

	if authenticationToken != "" {
		q.Set("authentication_token", authenticationToken)
	}

	var selfURL = new(url.URL)
	*selfURL = *o.SelfURL
	selfURL.Path = "/oauth2/auth"

	q.Set("redirect_uri", selfURL.String())
	http.Redirect(w, r, p.String(), http.StatusFound)

	return nil
}

func (o *Handler) writeAuthorizeError(w http.ResponseWriter, ar fosite.AuthorizeRequester, err error) {
	if !ar.IsRedirectURIValid() {
		var rfcerr = fosite.ErrorToRFC6749Error(err)
		var redirectURI = new(url.URL)
		*redirectURI = *o.ConsentURL

		query := redirectURI.Query()
		query.Add("error", rfcerr.Name)
		query.Add("error_description", rfcerr.Description)
		redirectURI.RawQuery = query.Encode()

		w.Header().Add("Location", redirectURI.String())
		w.WriteHeader(http.StatusFound)
		return
	}

	o.OAuth2.WriteAuthorizeError(w, ar, err)
}
