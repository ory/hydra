// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	OpenIDConnectKeyName = "hydra.openid.id-token"

	DefaultConsentPath = "/oauth2/consent-fallback"
	TokenPath          = "/oauth2/token"
	AuthPath           = "/oauth2/auth"

	UserinfoPath  = "/userinfo"
	WellKnownPath = "/.well-known/openid-configuration"
	JWKPath       = "/.well-known/jwks.json"

	// IntrospectPath points to the OAuth2 introspection endpoint.
	IntrospectPath = "/oauth2/introspect"
	RevocationPath = "/oauth2/revoke"

	IntrospectScope = "hydra.introspect"

	consentCookieName = "consent_session"
)

// swagger:model wellKnown
type WellKnown struct {
	// URL using the https scheme with no query or fragment component that the OP asserts as its Issuer Identifier.
	// If Issuer discovery is supported , this value MUST be identical to the issuer value returned
	// by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this Issuer.
	//
	// required: true
	Issuer string `json:"issuer"`

	// URL of the OP's OAuth 2.0 Authorization Endpoint
	//
	// required: true
	AuthURL string `json:"authorization_endpoint"`

	// URL of the OP's OAuth 2.0 Token Endpoint
	//
	// required: true
	TokenURL string `json:"token_endpoint"`

	// URL of the OP's JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate
	// signatures from the OP. The JWK Set MAY also contain the Server's encryption key(s), which are used by RPs
	// to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use)
	// parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage.
	// Although some algorithms allow the same key to be used for both signatures and encryption, doing so is
	// NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of
	// keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.
	//
	// required: true
	JWKsURI string `json:"jwks_uri"`

	// JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include
	// pairwise and public.
	//
	// required: true
	SubjectTypes []string `json:"subject_types_supported"`

	// JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID
	// Providers MUST support the code, id_token, and the token id_token Response Type values.
	//
	// required: true
	ResponseTypes []string `json:"response_types_supported"`

	// JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply
	// values for. Note that for privacy or other reasons, this might not be an exhaustive list.
	ClaimsSupported []string `json:"claims_supported"`

	// URL of the OP's UserInfo Endpoint.
	UserinfoEndpoint string `json:"userinfo_endpoint"`

	// SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST
	// support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used
	ScopesSupported []string `json:"scopes_supported"`

	// JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are
	// client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`

	// JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token
	// to encode the Claims in a JWT.
	//
	// required: true
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(TokenPath, h.TokenHandler)
	r.GET(AuthPath, h.AuthHandler)
	r.POST(AuthPath, h.AuthHandler)
	r.GET(DefaultConsentPath, h.DefaultConsentHandler)
	r.POST(IntrospectPath, h.IntrospectHandler)
	r.POST(RevocationPath, h.RevocationHandler)
	r.GET(WellKnownPath, h.WellKnownHandler)
	r.GET(UserinfoPath, h.UserinfoHandler)
	r.POST(UserinfoPath, h.UserinfoHandler)
}

// swagger:route GET /.well-known/openid-configuration oAuth2 getWellKnown
//
// Server well known configuration
//
// The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll
// your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this
// flow at https://openid.net/specs/openid-connect-discovery-1_0.html
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: wellKnown
//       401: genericError
//       500: genericError
func (h *Handler) WellKnownHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userInfoEndpoint := h.Issuer + UserinfoPath
	if h.UserinfoEndpoint != "" {
		userInfoEndpoint = h.UserinfoEndpoint
	}

	claimsSupported := []string{"sub"}
	if h.ClaimsSupported != "" {
		claimsSupported = append(claimsSupported, strings.Split(h.ClaimsSupported, ",")...)
	}

	scopesSupported := []string{"offline", "openid"}
	if h.ScopesSupported != "" {
		scopesSupported = append(scopesSupported, strings.Split(h.ScopesSupported, ",")...)
	}

	h.H.Write(w, r, &WellKnown{
		Issuer:                            h.Issuer,
		AuthURL:                           h.Issuer + AuthPath,
		TokenURL:                          h.Issuer + TokenPath,
		JWKsURI:                           h.Issuer + JWKPath,
		SubjectTypes:                      []string{"pairwise", "public"},
		ResponseTypes:                     []string{"code", "code id_token", "id_token", "token id_token", "token", "token id_token code"},
		ClaimsSupported:                   claimsSupported,
		ScopesSupported:                   scopesSupported,
		UserinfoEndpoint:                  userInfoEndpoint,
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
	})
}

// swagger:route POST /userinfo oAuth2 userinfo
//
// OpenID Connect Userinfo
//
// This endpoint returns the payload of the ID Token, including the idTokenExtra values, of the provided OAuth 2.0 access token.
// The endpoint implements http://openid.net/specs/openid-connect-core-1_0.html#UserInfo .
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
//       200: userinfoResponse
//       401: genericError
//       500: genericError
func (h *Handler) UserinfoHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := NewSession("")
	ar, err := h.OAuth2.IntrospectToken(r.Context(), fosite.AccessTokenFromRequest(r), fosite.AccessToken, session)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	interim := ar.GetSession().(*Session).IDTokenClaims().ToMap()
	delete(interim, "aud")
	delete(interim, "iss")
	delete(interim, "nonce")
	delete(interim, "at_hash")
	delete(interim, "c_hash")
	delete(interim, "auth_time")
	delete(interim, "iat")
	delete(interim, "exp")

	h.H.Write(w, r, interim)
}

// swagger:route POST /oauth2/revoke oAuth2 revokeOAuth2Token
//
// Revoke OAuth2 tokens
//
// Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no
// longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token.
// Revoking a refresh token also invalidates the access token that was created with it.
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Schemes: http, https
//
//     Security:
//       basic:
//
//     Responses:
//       200: emptyResponse
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

// swagger:route POST /oauth2/introspect oAuth2 introspectOAuth2Token
//
// Introspect OAuth2 tokens
//
// The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token
// is neither expired nor revoked. If a token is active, additional information on the token will be included. You can
// set additional data for a token by setting `accessTokenExtra` during the consent flow.
//
//  ```
//  {
//    "resources": ["rn:hydra:oauth2:tokens"],
//    "actions": ["introspect"],
//    "effect": "allow"
//  }
//  ```
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
//       oauth2: hydra.introspect
//
//     Responses:
//       200: introspectOAuth2TokenResponse
//       401: genericError
//       500: genericError
func (h *Handler) IntrospectHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if token := h.W.TokenFromRequest(r); token != "" {
		if _, err := h.W.TokenAllowed(r.Context(), token, &firewall.TokenAccessRequest{
			Resource: fmt.Sprintf(h.PrefixResource("oauth2:tokens")),
			Action:   "introspect",
		}, IntrospectScope); err != nil {
			h.H.WriteError(w, r, err)
			return
		}
	} else if client, _, ok := r.BasicAuth(); ok {
		// If no token is given, we do not need a scope.
		if err := h.W.IsAllowed(r.Context(), &firewall.AccessRequest{
			Subject:  client,
			Resource: fmt.Sprintf(h.PrefixResource("oauth2:tokens")),
			Action:   "introspect",
		}); err != nil {
			h.H.WriteError(w, r, err)
			return
		}
	} else {
		h.H.WriteError(w, r, errors.WithStack(fosite.ErrRequestUnauthorized))
		return
	}

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
	if err = json.NewEncoder(w).Encode(&Introspection{
		Active:    true,
		ClientID:  resp.GetAccessRequester().GetClient().GetID(),
		Scope:     strings.Join(resp.GetAccessRequester().GetGrantedScopes(), " "),
		ExpiresAt: exp.Unix(),
		IssuedAt:  resp.GetAccessRequester().GetRequestedAt().Unix(),
		Subject:   resp.GetAccessRequester().GetSession().GetSubject(),
		Username:  resp.GetAccessRequester().GetSession().GetUsername(),
		Extra:     resp.GetAccessRequester().GetSession().(*Session).Extra,
		Issuer:    h.Issuer,
	}); err != nil {
		pkg.LogError(err, h.L)
	}
}

// swagger:route POST /oauth2/token oAuth2 oauthToken
//
// The OAuth 2.0 token endpoint
//
// This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows.
// OAuth2 is a very popular protocol and a library for your programming language will exists.
//
// To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749
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
//       oauth2:
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
			if h.ScopeStrategy(accessRequest.GetClient().GetScopes(), scope) {
				accessRequest.GrantScope(scope)
			}
		}
	}

	accessResponse, err := h.OAuth2.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		pkg.LogError(err, h.L)
		h.OAuth2.WriteAccessError(w, accessRequest, err)
		return
	}

	h.OAuth2.WriteAccessResponse(w, accessRequest, accessResponse)
}

// swagger:route GET /oauth2/auth oAuth2 oauthAuth
//
// The OAuth 2.0 authorize endpoint
//
// This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows.
// OAuth2 is a very popular protocol and a library for your programming language will exists.
//
// To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Schemes: http, https
//
//     Responses:
//       302: emptyResponse
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
	consent := authorizeRequest.GetRequestForm().Get("consent")
	if consent == "" {
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
	session, err := h.Consent.ValidateConsentRequest(authorizeRequest, consent, cookie)
	if err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	if err := cookie.Save(r, w); err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, errors.Wrapf(fosite.ErrServerError, "Could not store session cookie: %s", err))
		return
	}

	// done
	response, err := h.OAuth2.NewAuthorizeResponse(ctx, authorizeRequest, session)
	if err != nil {
		pkg.LogError(err, h.L)
		h.writeAuthorizeError(w, authorizeRequest, err)
		return
	}

	h.OAuth2.WriteAuthorizeResponse(w, authorizeRequest, response)
}

func (h *Handler) redirectToConsent(w http.ResponseWriter, r *http.Request, authorizeRequest fosite.AuthorizeRequester) error {
	// Error can be ignored because a session will always be returned
	cookie, _ := h.CookieStore.Get(r, consentCookieName)

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	authUrl, err := url.Parse(h.Issuer + AuthPath)
	if err != nil {
		return err
	}
	authHost, _, err := net.SplitHostPort(authUrl.Host)
	if err != nil {
		authHost = authUrl.Host
	}
	if authHost != host {
		h.L.WithFields(logrus.Fields{
			"request_host": host,
			"issuer_host":  authHost,
		}).Warnln("Host from auth request does not match issuer host. The consent return redirect may fail.")
	}
	authUrl.RawQuery = r.URL.RawQuery

	challenge, err := h.Consent.CreateConsentRequest(authorizeRequest, authUrl.String(), cookie)
	if err != nil {
		return err
	}

	p := h.ConsentURL
	q := p.Query()
	q.Set("consent", challenge)
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
