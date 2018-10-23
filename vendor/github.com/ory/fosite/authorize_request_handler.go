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
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

package fosite

import (
	"net/http"
	"strings"

	"context"

	"fmt"

	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/ory/go-convenience/stringslice"
	"github.com/ory/go-convenience/stringsx"
	"github.com/pkg/errors"
)

func (f *Fosite) authorizeRequestParametersFromOpenIDConnectRequest(request *AuthorizeRequest) error {
	var scope Arguments = stringsx.Splitx(request.Form.Get("scope"), " ")

	// Even if a scope parameter is present in the Request Object value, a scope parameter MUST always be passed using
	// the OAuth 2.0 request syntax containing the openid scope value to indicate to the underlying OAuth 2.0 logic that this is an OpenID Connect request.
	// Source: http://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth
	if !scope.Has("openid") {
		return nil
	}

	if len(request.Form.Get("request")+request.Form.Get("request_uri")) == 0 {
		return nil
	} else if len(request.Form.Get("request")) > 0 && len(request.Form.Get("request_uri")) > 0 {
		return errors.WithStack(ErrInvalidRequest.WithHint(`OpenID Connect parameters "request" and "request_uri" were both given, but you can use at most one.`))
	}

	oidcClient, ok := request.Client.(OpenIDConnectClient)
	if !ok {
		if len(request.Form.Get("request_uri")) > 0 {
			return errors.WithStack(ErrRequestURINotSupported.WithHint(`OpenID Connect "request_uri" context was given, but the OAuth 2.0 Client does not implement advanced OpenID Connect capabilities.`))
		}
		return errors.WithStack(ErrRequestNotSupported.WithHint(`OpenID Connect "request" context was given, but the OAuth 2.0 Client does not implement advanced OpenID Connect capabilities.`))
	}

	if oidcClient.GetJSONWebKeys() == nil && len(oidcClient.GetJSONWebKeysURI()) == 0 {
		return errors.WithStack(ErrInvalidRequest.WithHint(`OpenID Connect "request" or "request_uri" context was given, but the OAuth 2.0 Client does not have any JSON Web Keys registered.`))
	}

	assertion := request.Form.Get("request")
	if location := request.Form.Get("request_uri"); len(location) > 0 {
		if !stringslice.Has(oidcClient.GetRequestURIs(), location) {
			return errors.WithStack(ErrInvalidRequestURI.WithHint(fmt.Sprintf("Request URI \"%s\" is not whitelisted by the OAuth 2.0 Client.", location)))
		}

		hc := f.HTTPClient
		if hc == nil {
			hc = http.DefaultClient
		}

		response, err := hc.Get(location)
		if err != nil {
			return errors.WithStack(ErrInvalidRequestURI.WithHintf(`Unable to fetch OpenID Connect request parameters from "request_uri" because %s.`, err.Error()))
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return errors.WithStack(ErrInvalidRequestURI.WithHintf(`Unable to fetch OpenID Connect request parameters from "request_uri" because status code "%d" was expected, but got "%d".`, http.StatusOK, response.StatusCode))
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.WithStack(ErrInvalidRequestURI.WithHintf(`Unable to fetch OpenID Connect request parameters from "request_uri" because error %s occurred during body parsing.`, err))
		}

		assertion = string(body)
	}

	token, err := jwt.ParseWithClaims(assertion, new(jwt.MapClaims), func(t *jwt.Token) (interface{}, error) {
		if oidcClient.GetRequestObjectSigningAlgorithm() != fmt.Sprintf("%s", t.Header["alg"]) {
			return nil, errors.WithStack(ErrInvalidRequestObject.WithHintf(`The request object uses signing algorithm %s, but the requested OAuth 2.0 Client enforces signing algorithm %s.`, t.Header["alg"], oidcClient.GetRequestObjectSigningAlgorithm()))
		}

		if t.Method == jwt.SigningMethodNone {
			return jwt.UnsafeAllowNoneSignatureType, nil
		}

		switch t.Method.(type) {
		case *jwt.SigningMethodRSA:
			key, err := f.findClientPublicJWK(oidcClient, t)
			if err != nil {
				return nil, errors.WithStack(ErrInvalidRequestObject.WithHintf("Unable to retrieve signing key from OAuth 2.0 Client because %s.", err))
			}
			return key, nil
		case *jwt.SigningMethodECDSA:
			key, err := f.findClientPublicJWK(oidcClient, t)
			if err != nil {
				return nil, errors.WithStack(ErrInvalidRequestObject.WithHintf("Unable to retrieve signing key from OAuth 2.0 Client because %s.", err))
			}
			return key, nil
		case *jwt.SigningMethodRSAPSS:
			key, err := f.findClientPublicJWK(oidcClient, t)
			if err != nil {
				return nil, errors.WithStack(ErrInvalidRequestObject.WithHintf("Unable to retrieve signing key from OAuth 2.0 Client because %s.", err))
			}
			return key, nil
		default:
			return nil, errors.WithStack(ErrInvalidRequestObject.WithHintf(`This request object uses unsupported signing algorithm "%s"."`, t.Header["alg"]))
		}
	})
	if err != nil {
		// Do not re-process already enhanced errors
		if e, ok := errors.Cause(err).(*jwt.ValidationError); ok {
			if e.Inner != nil {
				return e.Inner
			}
			return errors.WithStack(ErrInvalidRequestObject.WithHintf("Unable to verify the request object's signature.").WithDebug(err.Error()))
		}
		return err
	} else if err := token.Claims.Valid(); err != nil {
		return errors.WithStack(ErrInvalidRequestObject.WithHint("Unable to verify the request object because its claims could not be validated, check if the expiry time is set correctly.").WithDebug(err.Error()))
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return errors.WithStack(ErrInvalidRequestObject.WithHint("Unable to type assert claims from request object.").WithDebugf(`Got claims of type %T but expected type "*jwt.MapClaims".`, token.Claims))
	}

	for k, v := range *claims {
		request.Form.Set(k, fmt.Sprintf("%s", v))
	}

	claimScope := stringsx.Splitx(request.Form.Get("scope"), " ")
	for _, s := range scope {
		if !stringslice.Has(claimScope, s) {
			claimScope = append(claimScope, s)
		}
	}

	request.Form.Set("scope", strings.Join(claimScope, " "))
	return nil
}

func (f *Fosite) validateAuthorizeRedirectURI(r *http.Request, request *AuthorizeRequest) error {
	// Fetch redirect URI from request
	rawRedirURI, err := GetRedirectURIFromRequestValues(request.Form)
	if err != nil {
		return err
	}

	// Validate redirect uri
	redirectURI, err := MatchRedirectURIWithClientRedirectURIs(rawRedirURI, request.Client)
	if err != nil {
		return err
	} else if !IsValidRedirectURI(redirectURI) {
		return errors.WithStack(ErrInvalidRequest.WithHintf(`The redirect URI "%s" contains an illegal character (for example #) or is otherwise invalid.`, redirectURI))
	}
	request.RedirectURI = redirectURI
	return nil
}

func (f *Fosite) validateAuthorizeScope(r *http.Request, request *AuthorizeRequest) error {
	scope := removeEmpty(strings.Split(request.Form.Get("scope"), " "))
	for _, permission := range scope {
		if !f.ScopeStrategy(request.Client.GetScopes(), permission) {
			return errors.WithStack(ErrInvalidScope.WithHintf(`The OAuth 2.0 Client is not allowed to request scope "%s".`, permission))
		}
	}
	request.SetRequestedScopes(scope)

	return nil
}

func (f *Fosite) validateResponseTypes(r *http.Request, request *AuthorizeRequest) error {
	// https://tools.ietf.org/html/rfc6749#section-3.1.1
	// Extension response types MAY contain a space-delimited (%x20) list of
	// values, where the order of values does not matter (e.g., response
	// type "a b" is the same as "b a").  The meaning of such composite
	// response types is defined by their respective specifications.
	responseTypes := removeEmpty(stringsx.Splitx(r.Form.Get("response_type"), " "))
	if len(responseTypes) == 0 {
		return errors.WithStack(ErrUnsupportedResponseType.WithHint(`The request is missing the "response_type"" parameter.`))
	}

	var found bool
	for _, t := range request.GetClient().GetResponseTypes() {
		if Arguments(responseTypes).Matches(removeEmpty(stringsx.Splitx(t, " "))...) {
			found = true
			break
		}
	}

	if !found {
		return errors.WithStack(ErrUnsupportedResponseType.WithHintf("The client is not allowed to request response_type \"%s\".", r.Form.Get("response_type")))
	}

	request.ResponseTypes = responseTypes
	return nil
}

func (f *Fosite) NewAuthorizeRequest(ctx context.Context, r *http.Request) (AuthorizeRequester, error) {
	request := &AuthorizeRequest{
		ResponseTypes:        Arguments{},
		HandledResponseTypes: Arguments{},
		Request:              *NewRequest(),
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil && err != http.ErrNotMultipart {
		return request, errors.WithStack(ErrInvalidRequest.WithHint("Unable to parse HTTP body, make sure to send a properly formatted form request body.").WithDebug(err.Error()))
	}

	request.Form = r.Form
	client, err := f.Store.GetClient(ctx, request.GetRequestForm().Get("client_id"))
	if err != nil {
		return request, errors.WithStack(ErrInvalidClient.WithHint("The requested OAuth 2.0 Client does not exist."))
	}
	request.Client = client

	if err := f.authorizeRequestParametersFromOpenIDConnectRequest(request); err != nil {
		return request, err
	}

	if err := f.validateAuthorizeRedirectURI(r, request); err != nil {
		return request, err
	}

	if err := f.validateAuthorizeScope(r, request); err != nil {
		return request, err
	}

	if len(request.Form.Get("registration")) > 0 {
		return request, errors.WithStack(ErrRegistrationNotSupported)
	}

	if err := f.validateResponseTypes(r, request); err != nil {
		return request, err
	}

	// rfc6819 4.4.1.8.  Threat: CSRF Attack against redirect-uri
	// The "state" parameter should be used to link the authorization
	// request with the redirect URI used to deliver the access token (Section 5.3.5).
	//
	// https://tools.ietf.org/html/rfc6819#section-4.4.1.8
	// The "state" parameter should not	be guessable
	state := request.Form.Get("state")
	if len(state) < MinParameterEntropy {
		// We're assuming that using less then 8 characters for the state can not be considered "unguessable"
		return request, errors.WithStack(ErrInvalidState.WithHintf(`Request parameter "state" must be at least be %d characters long to ensure sufficient entropy.`, MinParameterEntropy))
	}
	request.State = state

	return request, nil
}
