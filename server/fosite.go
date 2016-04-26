package server

import (
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/enigma/jwt"
	"github.com/ory-am/hydra/identity"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
)

type Handler struct {
	fosite           fosite.OAuth2Provider
	consentValidator oauth2.ConsentValidator
	identities       identity.IdentityProviderRegistry
	jwtGenerator     jwt.Enigma

	SelfURL         *url.URL
	ConsentURL      *url.URL
	ErrorHandlerURL *url.URL
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/oauth2/token", h.TokenHandler)

	r.GET("/oauth2/auth", h.AuthHandler)
	r.POST("/oauth2/auth", h.AuthHandler)
}

func (o *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var session oauth2.Session
	var ctx = fosite.NewContext()

	accessRequest, err := o.fosite.NewAccessRequest(ctx, r, &session)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAccessError(w, accessRequest, err)
		return
	}

	accessResponse, err := o.fosite.NewAccessResponse(ctx, r, accessRequest)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAccessError(w, accessRequest, err)
		return
	}

	o.fosite.WriteAccessResponse(w, accessRequest, accessResponse)
}

func (o *Handler) AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = fosite.NewContext()
	authorizeRequest, err := o.fosite.NewAuthorizeRequest(ctx, r)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, err)
		return
	}

	// A session_token will be available if the user was authenticated an gave consent
	consentToken := authorizeRequest.GetRequestForm().Get("consent_token")
	if consentToken == "" {
		// otherwise redirect to log in endpoint
		o.redirectToConsent(w, r, authorizeRequest, "")
	}

	// decode consent_token claims
	// verify anti-CSRF (inject state) and anti-replay token (expiry time, good value would be 10 seconds)
	session, err := o.consentValidator.ValidateConsentToken(authorizeRequest, consentToken)
	if err != nil {
		o.fosite.WriteAuthorizeError(w, authorizeRequest, errors.New(fosite.ErrAccessDenied))
		return
	}

	// done
	response, err := o.fosite.NewAuthorizeResponse(ctx, r, authorizeRequest, session)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, err)
		return
	}

	o.fosite.WriteAuthorizeResponse(w, authorizeRequest, response)
}

func (o *Handler) redirectToConsent(w http.ResponseWriter, r *http.Request, authorizeRequest fosite.AuthorizeRequester, authenticationToken string) error {
	var p *url.URL
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

	var selfURL *url.URL
	*selfURL = *o.SelfURL
	selfURL.Path = "/auth"

	q.Set("redirect_uri", selfURL.String())
	http.Redirect(w, r, p.String(), http.StatusFound)

	return nil
}
