package oauth2

import (
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/pkg"
)

const (
	OpenIDConnectKeyName = "hydra.openid.connect"

	ConsentPath = "/oauth2/consent"
	TokenPath   = "/oauth2/token"
	AuthPath    = "/oauth2/auth"

	// IntrospectPath points to the OAuth2 introspection endpoint.
	IntrospectPath = "/oauth2/introspect"
)

type Handler struct {
	OAuth2     fosite.OAuth2Provider
	Consent    ConsentStrategy
	ForcedHTTP bool

	ConsentURL url.URL
}

func (this *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(TokenPath, this.TokenHandler)
	r.GET(AuthPath, this.AuthHandler)
	r.POST(AuthPath, this.AuthHandler)
	r.GET(ConsentPath, this.DefaultConsentHandler)
	r.POST(IntrospectPath, this.Introspect)
}

func (this *Handler) Introspect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := herodot.NewContext()
	clientCtx, err := this.Warden.TokenValid(ctx, TokenFromRequest(r))
	if err != nil {
		this.H.WriteError(ctx, w, r, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		this.H.WriteError(ctx, w, r, err)
		return
	}

	auth, err := this.Warden.IntrospectToken(ctx, r.PostForm.Get("token"))
	if err != nil {
		this.H.Write(ctx, w, r, &inactive)
		return
	} else if clientCtx.Subject != auth.Audience {
		this.H.Write(ctx, w, r, &inactive)
		return
	}

	this.H.Write(ctx, w, r, auth)
}

func (this *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var session = NewSession("")
	var ctx = fosite.NewContext()

	accessRequest, err := this.OAuth2.NewAccessRequest(ctx, r, session)
	if err != nil {
		pkg.LogError(err)
		this.OAuth2.WriteAccessError(w, accessRequest, err)
		return
	}

	if accessRequest.GetGrantTypes().Exact("client_credentials") {
		session.Subject = accessRequest.GetClient().GetID()
		for _, scope := range accessRequest.GetRequestedScopes() {
			if fosite.HierarchicScopeStrategy(accessRequest.GetClient().GetScopes(), scope) {
				accessRequest.GrantScope(scope)
			}
		}
	}

	accessResponse, err := this.OAuth2.NewAccessResponse(ctx, r, accessRequest)
	if err != nil {
		pkg.LogError(err)
		this.OAuth2.WriteAccessError(w, accessRequest, err)
		return
	}

	this.OAuth2.WriteAccessResponse(w, accessRequest, accessResponse)
}

func (this *Handler) AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = fosite.NewContext()

	authorizeRequest, err := this.OAuth2.NewAuthorizeRequest(ctx, r)
	if err != nil {
		pkg.LogError(err)
		this.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	// A session_token will be available if the user was authenticated an gave consent
	consentToken := authorizeRequest.GetRequestForm().Get("consent")
	if consentToken == "" {
		// otherwise redirect to log in endpoint
		if err := this.redirectToConsent(w, r, authorizeRequest); err != nil {
			pkg.LogError(err)
			this.writeAuthorizeError(w, authorizeRequest, err)
			return
		}
		return
	}

	// decode consent_token claims
	// verify anti-CSRF (inject state) and anti-replay token (expiry time, good value would be 10 seconds)
	session, err := this.Consent.ValidateResponse(authorizeRequest, consentToken)
	if err != nil {
		pkg.LogError(err)
		this.writeAuthorizeError(w, authorizeRequest, errors.New(fosite.ErrAccessDenied))
		return
	}

	// done
	response, err := this.OAuth2.NewAuthorizeResponse(ctx, r, authorizeRequest, session)
	if err != nil {
		pkg.LogError(err)
		this.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	this.OAuth2.WriteAuthorizeResponse(w, authorizeRequest, response)
}

func (this *Handler) redirectToConsent(w http.ResponseWriter, r *http.Request, authorizeRequest fosite.AuthorizeRequester) error {
	schema := "https"
	if this.ForcedHTTP {
		schema = "http"
	}

	challenge, err := this.Consent.IssueChallenge(authorizeRequest, schema+"://"+r.Host+r.URL.String())
	if err != nil {
		return err
	}

	p := this.ConsentURL
	q := p.Query()
	q.Set("challenge", challenge)
	p.RawQuery = q.Encode()
	http.Redirect(w, r, p.String(), http.StatusFound)
	return nil
}

func (this *Handler) writeAuthorizeError(w http.ResponseWriter, ar fosite.AuthorizeRequester, err error) {
	if !ar.IsRedirectURIValid() {
		var rfcerr = fosite.ErrorToRFC6749Error(err)

		redirectURI := this.ConsentURL
		query := redirectURI.Query()
		query.Add("error", rfcerr.Name)
		query.Add("error_description", rfcerr.Description)
		redirectURI.RawQuery = query.Encode()

		w.Header().Add("Location", redirectURI.String())
		w.WriteHeader(http.StatusFound)
		return
	}

	this.OAuth2.WriteAuthorizeError(w, ar, err)
}
