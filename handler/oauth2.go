package handler

import (
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/enigma/jwt"
	"github.com/ory-am/hydra/connector"
	"github.com/ory-am/hydra/consent"
	"github.com/ory-am/hydra/identity"
	"github.com/ory-am/hydra/pkg"
)

type OAuth2Handler struct {
	fosite           fosite.OAuth2Provider
	consentValidator consent.Validator
	connectors       connector.ConnectorRegistry
	identities       identity.IdentityProviderRegistry
	jwtGenerator     jwt.Enigma

	SelfURL         *url.URL
	ConsentURL      *url.URL
	ErrorHandlerURL *url.URL
}

func (h *OAuth2Handler) SetRoutes(r *httprouter.Router) {

}

func (o *OAuth2Handler) TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var session consent.ConsentClaims
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

func (o *OAuth2Handler) AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = fosite.NewContext()

	authorizeRequest, err := o.fosite.NewAuthorizeRequest(ctx, r)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, err)
		return
	}

	// check if third party thingy
	conor, err := o.connectors.GetConnector(authorizeRequest.GetRequestForm().Get("connector"))
	if err == connector.ConnectorNotFound {
		// nothing to do
	} else if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, errors.New(fosite.ErrServerError))
		return
	} else {
		url, err := conor.PersistAuthorizeSession(authorizeRequest)
		if err != nil {
			pkg.LogError(errors.New(err))
			o.fosite.WriteAuthorizeError(w, authorizeRequest, err)
		}

		http.Redirect(w, r, url.String(), http.StatusFound)
		return
	}

	// A session_token will be available if the user was authenticated an gave consent
	consentToken := authorizeRequest.GetRequestForm().Get("consent_token")
	if consentToken == "" {
		// otherwise redirect to log in endpoint
		o.redirectToConsent(authorizeRequest, "")
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

func (o *OAuth2Handler) ConnectorCallbackHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		pkg.LogError(errors.New(err))
		pkg.ForwardToErrorHandler(w, r, err, o.ErrorHandlerURL)
		return
	}

	conor, err := o.connectors.GetConnector(p.ByName("connector"))
	if err != nil {
		pkg.LogError(errors.New(err))
		pkg.ForwardToErrorHandler(w, r, err, o.ErrorHandlerURL)
		return
	}

	authorizeRequest, err := conor.GetAuthorizeSession(r.Form)
	if err != nil {
		pkg.LogError(errors.New(err))
		pkg.ForwardToErrorHandler(w, r, err, o.ErrorHandlerURL)
		return
	}

	subject, err := conor.Exchange(r.Form)
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, errors.New(err))
		return
	}

	if err := o.identities.IsIdentityAuthenticable(subject); err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, errors.New(err))
		return
	}

	authenticationToken, _, err := o.jwtGenerator.Generate(&jwt.Claims{
		Subject: subject,
	}, &jwt.Header{})
	if err != nil {
		pkg.LogError(errors.New(err))
		o.fosite.WriteAuthorizeError(w, authorizeRequest, errors.New(err))
		return
	}

	o.redirectToConsent(authorizeRequest, authenticationToken)
}

func (o *OAuth2Handler) redirectToConsent(authorizeRequest fosite.AuthorizeRequester, authenticationToken string) error {
	p, err := url.Parse(o.ConsentURL)
	if err != nil {
		return errors.New(err)
	}

	p.Query().Set("client_id", authorizeRequest.GetClient().GetID())
	p.Query().Set("state", authorizeRequest.GetState())
	for _, scope := range authorizeRequest.GetScopes() {
		p.Query().Add("scope", scope)
	}

	if authenticationToken != "" {
		p.Query().Set("authentication_token", authenticationToken)
	}

	var selfURL *url.URL
	*selfURL = *o.SelfURL
	selfURL.Path = "/auth"
	p.Query().Set("redirect_uri", selfURL.String())
	return nil
}
