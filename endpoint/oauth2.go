package auth

import (
	"github.com/ory-am/fosite"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/ory-am/hydra/pkg"
	"github.com/go-errors/errors"
	"net/url"
	"github.com/ory-am/hydra/consent"
)

type OAuth2Handler struct {
	f fosite.OAuth2Provider

	cv consent.Validator

	SelfURL *url.URL
	SignInURL *url.URL
}

func (h *OAuth2Handler) SetRoutes(r *httprouter.Router) {

}

func (o *OAuth2Handler) TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := fosite.NewContext()

	accessRequest, err := o.f.NewAccessRequest(ctx, r, session)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.f.WriteAccessError(w, accessRequest, err)
		return
	}

	// nothing to do, really

	accessResponse, err := o.f.NewAccessResponse(ctx, r, accessRequest)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.f.WriteAccessError(w, accessRequest, err)
		return
	}

	o.f.WriteAccessResponse(w, accessRequest, accessResponse)
}

func (o *OAuth2Handler) AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := fosite.NewContext()

	authorizeRequest, err := o.f.NewAuthorizeRequest(ctx, r)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.f.WriteAuthorizeError(w, authorizeRequest, err)
		return
	}

	// A session_token will be available if the user was authenticated an gave consent
	consentToken := authorizeRequest.GetRequestForm().Get("consent_token")
	if consentToken == "" {
		// otherwise redirect to log in endpoint
		o.RedirectToSignIn(authorizeRequest)
	}

	// decode consent_token claims
	// verify anti-CSRF (inject state) and anti-replay token (expiry time, good value would be 10 seconds)

	session, err := o.cv.ValidateConsentToken(authorizeRequest, consentToken)
	if err != nil {
		o.f.WriteAuthorizeError(w, authorizeRequest, errors.New(fosite.ErrAccessDenied))
		return
	}

	// done
	response, err := o.f.NewAuthorizeResponse(ctx, r, authorizeRequest, session)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.f.WriteAuthorizeError(w, authorizeRequest, err)
		return
	}

	o.f.WriteAuthorizeResponse(w, authorizeRequest, response)
}

func (o *OAuth2Handler) ConnectorCallbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//
}

func (o *OAuth2Handler) RedirectToSignIn(authorizeRequest fosite.AuthorizeRequester) error {
	p, err := url.Parse(o.SignInURL)
	if err != nil {
		return errors.New(err)
	}

	p.Query().Set("client_id", authorizeRequest.GetClient().GetID())
	p.Query().Set("state", authorizeRequest.GetState())
	for _, scope := range authorizeRequest.GetScopes() {
		p.Query().Add("scope", scope)
	}

	var selfURL *url.URL
	*selfURL = *o.SelfURL
	selfURL.Path = "/auth"
	p.Query().Set("redirect_uri", selfURL.String())

	return nil
}