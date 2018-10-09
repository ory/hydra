/*
 * Copyright © 2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
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
 *
 */

package fosite

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

const clientAssertionJWTBearerType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

func (f *Fosite) findClientPublicJWK(oidcClient OpenIDConnectClient, t *jwt.Token) (interface{}, error) {
	if set := oidcClient.GetJSONWebKeys(); set != nil {
		return findPublicKey(t, set)
	}

	if location := oidcClient.GetJSONWebKeysURI(); len(location) > 0 {
		keys, err := f.JWKSFetcherStrategy.Resolve(location, false)
		if err != nil {
			return nil, err
		}

		if key, err := findPublicKey(t, keys); err == nil {
			return key, nil
		}

		keys, err = f.JWKSFetcherStrategy.Resolve(location, true)
		if err != nil {
			return nil, err
		}

		return findPublicKey(t, keys)
	}

	return nil, errors.WithStack(ErrInvalidClient.WithHint("The OAuth 2.0 Client has no JSON Web Keys set registered, but they are needed to complete the request."))
}

func (f *Fosite) AuthenticateClient(ctx context.Context, r *http.Request, form url.Values) (Client, error) {
	if assertionType := form.Get("client_assertion_type"); assertionType == clientAssertionJWTBearerType {
		assertion := form.Get("client_assertion")
		if len(assertion) == 0 {
			return nil, errors.WithStack(ErrInvalidRequest.WithHintf("The client_assertion request parameter must be set when using client_assertion_type of \"%s\".", clientAssertionJWTBearerType))
		}

		var clientID string
		var client Client

		token, err := jwt.ParseWithClaims(assertion, new(jwt.MapClaims), func(t *jwt.Token) (interface{}, error) {
			var err error
			clientID, _, err = clientCredentialsFromRequestBody(form, false)
			if err != nil {
				return nil, err
			}

			if clientID == "" {
				if claims, ok := t.Claims.(*jwt.MapClaims); !ok {
					return nil, errors.WithStack(ErrRequestUnauthorized.WithHint("Unable to type assert claims from client_assertion.").WithDebugf(`Expected claims to be of type "*jwt.MapClaims" but got "%T".`, t.Claims))
				} else if sub, ok := (*claims)["sub"].(string); !ok {
					return nil, errors.WithStack(ErrInvalidClient.WithHint(`The claim "sub" from the client_assertion JSON Web Token is undefined.`))
				} else {
					clientID = sub
				}
			}

			client, err = f.Store.GetClient(ctx, clientID)
			if err != nil {
				return nil, errors.WithStack(ErrInvalidClient.WithDebug(err.Error()))
			}

			oidcClient, ok := client.(OpenIDConnectClient)
			if !ok {
				return nil, errors.WithStack(ErrInvalidRequest.WithHint("The server configuration does not support OpenID Connect specific authentication methods."))
			}

			switch oidcClient.GetTokenEndpointAuthMethod() {
			case "client_secret_post":
				fallthrough
			case "client_secret_basic":
				return nil, errors.WithStack(ErrInvalidClient.WithHintf("The OAuth 2.0 Request uses the \"client_secret_jwt\" authentication method, but the OAuth 2.0 Client only supports the \"%s\" client authentication method. You must configure the OAuth 2.0 client's \"token_endpoint_auth_method\" value to accept \"client_secret_jwt\".", oidcClient.GetTokenEndpointAuthMethod()))
			case "client_secret_jwt":
				return nil, errors.WithStack(ErrInvalidClient.WithHint("This requested OAuth 2.0 client only supports client authentication method \"client_secret_jwt\", however that method is not supported by this server."))
			case "private_key_jwt":
			}

			if oidcClient.GetTokenEndpointAuthSigningAlgorithm() != fmt.Sprintf("%s", t.Header["alg"]) {
				return nil, errors.WithStack(ErrInvalidClient.WithHintf(`The "client_assertion"" uses signing algorithm "%s", but the requested OAuth 2.0 Client enforces signing algorithm "%s".`, t.Header["alg"], oidcClient.GetTokenEndpointAuthSigningAlgorithm()))
			}

			if _, ok := t.Method.(*jwt.SigningMethodRSA); ok {
				return f.findClientPublicJWK(oidcClient, t)
			} else if _, ok := t.Method.(*jwt.SigningMethodECDSA); ok {
				return f.findClientPublicJWK(oidcClient, t)
			} else if _, ok := t.Method.(*jwt.SigningMethodRSAPSS); ok {
				return f.findClientPublicJWK(oidcClient, t)
			} else if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
				return nil, errors.WithStack(ErrInvalidClient.WithHint("This authorization server does not support client authentication method \"client_secret_jwt\"."))
			}

			return nil, errors.WithStack(ErrInvalidClient.WithHintf("The client_assertion request parameter uses unsupported signing algorithm \"%s\".", t.Header["alg"]))
		})
		if err != nil {
			// Do not re-process already enhanced errors
			if e, ok := errors.Cause(err).(*jwt.ValidationError); ok {
				if e.Inner != nil {
					return nil, e.Inner
				}
				return nil, errors.WithStack(ErrInvalidClient.WithHint("Unable to verify the integrity of the \"client_assertion\" value.").WithDebug(err.Error()))
			}
			return nil, err
		} else if err := token.Claims.Valid(); err != nil {
			return nil, errors.WithStack(ErrInvalidClient.WithHint("Unable to verify the request object because its claims could not be validated, check if the expiry time is set correctly.").WithDebug(err.Error()))
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return nil, errors.WithStack(ErrInvalidClient.WithHint("Unable to type assert claims from request parameter \"client_assertion\".").WithDebugf(`Got claims of type %T but expected type "*jwt.MapClaims".`, token.Claims))
		}

		if !claims.VerifyIssuer(clientID, true) {
			return nil, errors.WithStack(ErrInvalidClient.WithHint("Claim \"iss\" from \"client_assertion\" must match the \"client_id\" of the OAuth 2.0 Client."))
		} else if f.TokenURL == "" {
			return nil, errors.WithStack(ErrMisconfiguration.WithHint("The authorization server's token endpoint URL has not been set."))
		} else if sub, ok := (*claims)["sub"].(string); !ok || sub != clientID {
			return nil, errors.WithStack(ErrInvalidClient.WithHint("Claim \"sub\" from \"client_assertion\" must match the \"client_id\" of the OAuth 2.0 Client."))
		} else if jti, ok := (*claims)["jti"].(string); !ok || len(jti) == 0 {
			return nil, errors.WithStack(ErrInvalidClient.WithHint("Claim \"jti\" from \"client_assertion\" must be set but is not."))
		}

		if auds, ok := (*claims)["aud"].([]interface{}); !ok {
			if !claims.VerifyAudience(f.TokenURL, true) {
				return nil, errors.WithStack(ErrInvalidClient.WithHintf("Claim \"audience\" from \"client_assertion\" must match the authorization server's token endpoint \"%s\".", f.TokenURL))
			}
		} else {
			var found bool
			for _, aud := range auds {
				if a, ok := aud.(string); ok && a == f.TokenURL {
					found = true
					break
				}
			}

			if !found {
				return nil, errors.WithStack(ErrInvalidClient.WithHintf("Claim \"audience\" from \"client_assertion\" must match the authorization server's token endpoint \"%s\".", f.TokenURL))
			}
		}

		return client, nil
	} else if len(assertionType) > 0 {
		return nil, errors.WithStack(ErrInvalidRequest.WithHintf("Unknown client_assertion_type \"%s\".", assertionType))
	}

	clientID, clientSecret, err := clientCredentialsFromRequest(r, form)
	if err != nil {
		return nil, err
	}

	client, err := f.Store.GetClient(ctx, clientID)
	if err != nil {
		return nil, errors.WithStack(ErrInvalidClient.WithDebug(err.Error()))
	}

	if oidcClient, ok := client.(OpenIDConnectClient); !ok {
		// If this isn't an OpenID Connect client then we actually don't care about any of this, just continue!
	} else if ok && form.Get("client_id") != "" && form.Get("client_secret") != "" && oidcClient.GetTokenEndpointAuthMethod() != "client_secret_post" {
		return nil, errors.WithStack(ErrInvalidClient.WithHintf("The OAuth 2.0 Client supports client authentication method \"%s\", but method \"client_secret_post\" was requested. You must configure the OAuth 2.0 client's \"token_endpoint_auth_method\" value to accept \"client_secret_post\".", oidcClient.GetTokenEndpointAuthMethod()))
	} else if _, _, basicOk := r.BasicAuth(); basicOk && ok && oidcClient.GetTokenEndpointAuthMethod() != "client_secret_basic" {
		return nil, errors.WithStack(ErrInvalidClient.WithHintf("The OAuth 2.0 Client supports client authentication method \"%s\", but method \"client_secret_basic\" was requested. You must configure the OAuth 2.0 client's \"token_endpoint_auth_method\" value to accept \"client_secret_basic\".", oidcClient.GetTokenEndpointAuthMethod()))
	} else if ok && oidcClient.GetTokenEndpointAuthMethod() != "none" && client.IsPublic() {
		return nil, errors.WithStack(ErrInvalidClient.WithHintf("The OAuth 2.0 Client supports client authentication method \"%s\", but method \"none\" was requested. You must configure the OAuth 2.0 client's \"token_endpoint_auth_method\" value to accept \"none\".", oidcClient.GetTokenEndpointAuthMethod()))
	}

	if client.IsPublic() {
		return client, nil
	}

	// Enforce client authentication
	if err := f.Hasher.Compare(ctx, client.GetHashedSecret(), []byte(clientSecret)); err != nil {
		return nil, errors.WithStack(ErrInvalidClient.WithDebug(err.Error()))
	}

	return client, nil
}

func findPublicKey(t *jwt.Token, set *jose.JSONWebKeySet) (*rsa.PublicKey, error) {
	kid, ok := t.Header["kid"].(string)
	if !ok {
		return nil, errors.WithStack(ErrInvalidRequest.WithHint("The JSON Web Token must contain a kid header value but did not."))
	}

	keys := set.Key(kid)
	if len(keys) == 0 {
		return nil, errors.WithStack(ErrInvalidRequest.WithHintf("The JSON Web Token uses signing key with kid \"%s\", which could not be found.", kid))
	}

	for _, key := range keys {
		if key.Use != "sig" {
			continue
		}
		if k, ok := key.Key.(*rsa.PublicKey); ok {
			return k, nil
		}
	}

	return nil, errors.WithStack(ErrInvalidRequest.WithHintf("Unable to find RSA public key with use=\"sig\" for kid \"%s\" in JSON Web Key Set.", kid))
}

func clientCredentialsFromRequest(r *http.Request, form url.Values) (clientID, clientSecret string, err error) {
	if id, secret, ok := r.BasicAuth(); !ok {
		return clientCredentialsFromRequestBody(form, true)
	} else if clientID, err = url.QueryUnescape(id); err != nil {
		return "", "", errors.WithStack(ErrInvalidRequest.WithHint(`The client id in the HTTP authorization header could not be decoded from "application/x-www-form-urlencoded".`).WithDebug(err.Error()))
	} else if clientSecret, err = url.QueryUnescape(secret); err != nil {
		return "", "", errors.WithStack(ErrInvalidRequest.WithHint(`The client secret in the HTTP authorization header could not be decoded from "application/x-www-form-urlencoded".`).WithDebug(err.Error()))
	}

	return clientID, clientSecret, nil
}

func clientCredentialsFromRequestBody(form url.Values, forceID bool) (clientID, clientSecret string, err error) {
	clientID = form.Get("client_id")
	clientSecret = form.Get("client_secret")

	if clientID == "" && forceID {
		return "", "", errors.WithStack(ErrInvalidRequest.WithHint("Client credentials missing or malformed in both HTTP Authorization header and HTTP POST body."))
	}

	if clientID, err = url.QueryUnescape(clientID); err != nil {
		return "", "", errors.WithStack(ErrInvalidRequest.WithHint(`The client id in the HTTP authorization header could not be decoded from "application/x-www-form-urlencoded".`).WithDebug(err.Error()))
	} else if clientSecret, err = url.QueryUnescape(clientSecret); err != nil {
		return "", "", errors.WithStack(ErrInvalidRequest.WithHint(`The client secret in the HTTP authorization header could not be decoded from "application/x-www-form-urlencoded".`).WithDebug(err.Error()))
	}

	return clientID, clientSecret, nil
}
