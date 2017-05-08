package oauth2

import (
	"net/http"
	"net/url"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"strings"
	"time"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

const (
	OpenIDConnectKeyName = "hydra.openid.id-token"

	ConsentPath = "/oauth2/consent"
	TokenPath   = "/oauth2/token"
	AuthPath    = "/oauth2/auth"

	WellKnownPath = "/.well-known/openid-configuration"
	JWKPath = "/.well-known/jwks.json"

	// IntrospectPath points to the OAuth2 introspection endpoint.
	IntrospectPath = "/oauth2/introspect"
	RevocationPath = "/oauth2/revoke"

	consentCookieName = "consent_session"
)

type Handler struct {
	OAuth2  fosite.OAuth2Provider
	Consent ConsentStrategy

	H herodot.Writer

	ForcedHTTP bool
	ConsentURL url.URL

	AccessTokenLifespan time.Duration
	CookieStore         sessions.Store
	
	L logrus.FieldLogger

	Issuer string
}

// swagger:model WellKnown
type WellKnown struct {
	Issuer        string   `json:"issuer"`
	AuthURL       string   `json:"authorization_endpoint"`
	TokenURL      string   `json:"token_endpoint"`
	SubjectTypes  []string `json:"subject_types_supported"`
	JWKsURL       string   `json:"jwks_uri"`
	SigningAlgs   []string `json:"id_token_signing_alg_values_supported"`
	ResponseTypes []string `json:"response_types_supported"`
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(TokenPath, h.TokenHandler)
	r.GET(AuthPath, h.AuthHandler)
	r.POST(AuthPath, h.AuthHandler)
	r.GET(ConsentPath, h.DefaultConsentHandler)
	r.POST(IntrospectPath, h.IntrospectHandler)
	r.POST(RevocationPath, h.RevocationHandler)
	r.GET(WellKnownPath, h.WellKnownHandler)
}

// swagger:route GET /.well-known/openid-configuration oauth2 WellKnownHandler
//
// Server well known configuration
//
// For more information, please refer to https://openid.net/specs/openid-connect-discovery-1_0.html
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2:
//
//     Responses:
//       200: WellKnown
//       401: genericError
//       500: genericError
func (h *Handler) WellKnownHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	wellKnown := WellKnown{
		Issuer:        h.Issuer,
		AuthURL:       h.Issuer + AuthPath,
		TokenURL:      h.Issuer + TokenPath,
		SubjectTypes:  []string{"pairwise", "public"},
		JWKsURL:       h.Issuer + JWKPath,
		SigningAlgs:   []string{"RS256"},
		ResponseTypes: []string{"code", "code id_token", "id_token", "token id_token", "token"},
	}
	h.H.Write( w, r, wellKnown)
}

// swagger:route POST /oauth2/revoke oauth2 revokeOAuthToken
//
// Revoke an OAuth2 access token
//
// For more information, please refer to https://tools.ietf.org/html/rfc7009
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2:
//
//     Responses:
//       200:
//       401: genericError
//       500: genericError
func (h *Handler) RevocationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = fosite.NewContext()

	err := h.OAuth2.NewRevocationRequest(ctx, r)
	if err != nil {
		pkg.LogError(err, h.L)
	}

	h.OAuth2.WriteRevocationResponse(w, err)
}

// swagger:route POST /oauth2/introspect oauth2 introspectOAuthToken
//
// Introspect an OAuth2 access token
//
// For more information, please refer to https://tools.ietf.org/html/rfc7662
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2:
//
//     Responses:
//       200: introspectOAuthTokenResponse
//       401: genericError
//       500: genericError
func (h *Handler) IntrospectHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var session = NewSession("")

	var ctx = fosite.NewContext()
	resp, err := h.OAuth2.NewIntrospectionRequest(ctx, r, session)
	if err != nil {
		pkg.LogError(err, h.L)
		h.OAuth2.WriteIntrospectionError(w, err)
		return
	}

	exp := resp.GetAccessRequester().GetSession().GetExpiresAt(fosite.AccessToken)
	if exp.IsZero() {
		exp = resp.GetAccessRequester().GetRequestedAt().Add(h.AccessTokenLifespan)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(&Introspection{
		Active:    true,
		ClientID:  resp.GetAccessRequester().GetClient().GetID(),
		Scope:     strings.Join(resp.GetAccessRequester().GetGrantedScopes(), " "),
		ExpiresAt: exp.Unix(),
		IssuedAt:  resp.GetAccessRequester().GetRequestedAt().Unix(),
		Subject:   resp.GetAccessRequester().GetSession().GetSubject(),
		Username:  resp.GetAccessRequester().GetSession().GetUsername(),
		Extra:     resp.GetAccessRequester().GetSession().(*Session).Extra,
		Audience:  resp.GetAccessRequester().GetClient().GetID(),
	})
	if err != nil {
		pkg.LogError(err, h.L)
	}
}

// swagger:route POST /oauth2/token oauth2 oauthToken
//
// The OAuth 2.0 Token endpoint
//
// For more information, please refer to https://tools.ietf.org/html/rfc6749#section-4
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       basic:
//
//     Responses:
//       200: oauthTokenResponse
//       401: genericError
//       500: genericError
func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var session = NewSession("")
	var ctx = fosite.NewContext()

	accessRequest, err := h.OAuth2.NewAccessRequest(ctx, r, session)
	if err != nil {
		pkg.LogError(err, h.L)
		h.OAuth2.WriteAccessError(w, accessRequest, err)
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

	accessResponse, err := h.OAuth2.NewAccessResponse(ctx, r, accessRequest)
	if err != nil {
		pkg.LogError(err, h.L)
		h.OAuth2.WriteAccessError(w, accessRequest, err)
		return
	}

	h.OAuth2.WriteAccessResponse(w, accessRequest, accessResponse)
}

// swagger:route GET /oauth2/auth oauth2 oauthAuth
//
// The OAuth 2.0 Auth endpoint
//
// For more information, please refer to https://tools.ietf.org/html/rfc6749#section-4
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Schemes: http, https
//
//     Responses:
//       302:
//       401: genericError
//       500: genericError
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = fosite.NewContext()

	authorizeRequest, err := h.OAuth2.NewAuthorizeRequest(ctx, r)
	if err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	// A session_token will be available if the user was authenticated an gave consent
	consentToken := authorizeRequest.GetRequestForm().Get("consent")
	if consentToken == "" {
		// otherwise redirect to log in endpoint
		if err := h.redirectToConsent(w, r, authorizeRequest); err != nil {
			pkg.LogError(err, h.L)
			h.writeAuthorizeError(w, authorizeRequest, err)
			return
		}
		return
	}

	cookie, err := h.CookieStore.Get(r, consentCookieName)
	if err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, errors.Wrapf(fosite.ErrServerError, "Could not open session: %s", err))
		return
	}

	// decode consent_token claims
	// verify anti-CSRF (inject state) and anti-replay token (expiry time, good value would be 10 seconds)
	session, err := h.Consent.ValidateResponse(authorizeRequest, consentToken, cookie)
	if err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, errors.Wrap(fosite.ErrAccessDenied, ""))
		return
	}

	if err := cookie.Save(r, w); err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, errors.Wrapf(fosite.ErrServerError, "Could not store session cookie: %s", err))
		return
	}

	// done
	response, err := h.OAuth2.NewAuthorizeResponse(ctx, r, authorizeRequest, session)
	if err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	h.OAuth2.WriteAuthorizeResponse(w, authorizeRequest, response)
}

func (h *Handler) redirectToConsent(w http.ResponseWriter, r *http.Request, authorizeRequest fosite.AuthorizeRequester) error {
	schema := "https"
	if h.ForcedHTTP {
		schema = "http"
	}

	// Error can be ignored because a session will always be returned
	cookie, _ := h.CookieStore.Get(r, consentCookieName)

	challenge, err := h.Consent.IssueChallenge(authorizeRequest, schema+"://"+r.Host+r.URL.String(), cookie)
	if err != nil {
		return err
	}

	p := h.ConsentURL
	q := p.Query()
	q.Set("challenge", challenge)
	p.RawQuery = q.Encode()

	if err := cookie.Save(r, w); err != nil {
		return err
	}

	http.Redirect(w, r, p.String(), http.StatusFound)
	return nil
}

func (h *Handler) writeAuthorizeError(w http.ResponseWriter, ar fosite.AuthorizeRequester, err error) {
	if !ar.IsRedirectURIValid() {
		var rfcerr = fosite.ErrorToRFC6749Error(err)

		redirectURI := h.ConsentURL
		query := redirectURI.Query()
		query.Add("error", rfcerr.Name)
		query.Add("error_description", rfcerr.Description)
		redirectURI.RawQuery = query.Encode()

		w.Header().Add("Location", redirectURI.String())
		w.WriteHeader(http.StatusFound)
		return
	}

	h.OAuth2.WriteAuthorizeError(w, ar, err)
}
