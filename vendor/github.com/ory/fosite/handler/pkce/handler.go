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

package pkce

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
)

type Handler struct {
	// If set to true, public clients must use PKCE.
	Force bool

	// Whether or not to allow the plain challenge method (S256 should be used whenever possible, plain is really discouraged).
	EnablePlainChallengeMethod bool

	AuthorizeCodeStrategy oauth2.AuthorizeCodeStrategy
	Storage               PKCERequestStorage
}

func (c *Handler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	// This let's us define multiple response types, for example open id connect's id_token
	if !ar.GetResponseTypes().Exact("code") {
		return nil
	}

	if !ar.GetClient().IsPublic() {
		return nil
	}

	challenge := ar.GetRequestForm().Get("code_challenge")
	method := ar.GetRequestForm().Get("code_challenge_method")
	if err := c.validate(challenge, method); err != nil {
		return err
	}

	code := resp.GetCode()
	if len(code) == 0 {
		return errors.WithStack(fosite.ErrServerError.WithDebug("The PKCE handler must be loaded after the authorize code handler."))
	}

	signature := c.AuthorizeCodeStrategy.AuthorizeCodeSignature(code)
	if err := c.Storage.CreatePKCERequestSession(ctx, signature, ar.Sanitize([]string{
		"code_challenge",
		"code_challenge_method",
	})); err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	return nil
}

func (c *Handler) validate(challenge, method string) error {
	if c.Force && challenge == "" {
		//If the server requires Proof Key for Code Exchange (PKCE) by OAuth
		//public clients and the client does not send the "code_challenge" in
		//the request, the authorization endpoint MUST return the authorization
		//error response with the "error" value set to "invalid_request".  The
		//"error_description" or the response of "error_uri" SHOULD explain the
		//nature of error, e.g., code challenge required.

		return errors.WithStack(fosite.ErrInvalidRequest.
			WithHint("Public clients must include a code_challenge when performing the authorize code flow, but it is missing.").
			WithDebug("The server is configured in a way that enforces PKCE for public clients."))
	}

	if !c.Force && challenge == "" {
		return nil
	}

	//If the server supporting PKCE does not support the requested
	//transformation, the authorization endpoint MUST return the
	//authorization error response with "error" value set to
	//"invalid_request".  The "error_description" or the response of
	//"error_uri" SHOULD explain the nature of error, e.g., transform
	//algorithm not supported.
	switch method {
	case "S256":
		break
	case "plain":
		fallthrough
	case "":
		if !c.EnablePlainChallengeMethod {
			return errors.WithStack(fosite.ErrInvalidRequest.
				WithHint("Public clients must use code_challenge_method=S256, plain is not allowed.").
				WithDebug("The server is configured in a way that enforces PKCE S256 as challenge method for public clients."))
		}
		break
	default:
		return errors.WithStack(fosite.ErrInvalidRequest.
			WithHint("The code_challenge_method is not supported, use S256 instead."))
	}
	return nil
}

func (c *Handler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !request.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().IsPublic() {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	code := request.GetRequestForm().Get("code")
	signature := c.AuthorizeCodeStrategy.AuthorizeCodeSignature(code)
	authorizeRequest, err := c.Storage.GetPKCERequestSession(ctx, signature, request.GetSession())
	if errors.Cause(err) == fosite.ErrNotFound {
		return errors.WithStack(fosite.ErrInvalidGrant.WithHint("Unable to find initial PKCE data tied to this request").WithDebug(err.Error()))
	} else if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	if err := c.Storage.DeletePKCERequestSession(ctx, signature); err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}

	//code_verifier
	//REQUIRED.  Code verifier
	//
	//The "code_challenge_method" is bound to the Authorization Code when
	//the Authorization Code is issued.  That is the method that the token
	//endpoint MUST use to verify the "code_verifier".
	verifier := request.GetRequestForm().Get("code_verifier")
	challenge := authorizeRequest.GetRequestForm().Get("code_challenge")
	method := authorizeRequest.GetRequestForm().Get("code_challenge_method")
	if err := c.validate(challenge, method); err != nil {
		return err
	}

	if !c.Force && challenge == "" && verifier == "" {
		return nil
	}

	//Upon receipt of the request at the token endpoint, the server
	//verifies it by calculating the code challenge from the received
	//"code_verifier" and comparing it with the previously associated
	//"code_challenge", after first transforming it according to the
	//"code_challenge_method" method specified by the client.
	//
	//	If the "code_challenge_method" from Section 4.3 was "S256", the
	//received "code_verifier" is hashed by SHA-256, base64url-encoded, and
	//then compared to the "code_challenge", i.e.:
	//
	//BASE64URL-ENCODE(SHA256(ASCII(code_verifier))) == code_challenge
	//
	//If the "code_challenge_method" from Section 4.3 was "plain", they are
	//compared directly, i.e.:
	//
	//code_verifier == code_challenge.
	//
	//	If the values are equal, the token endpoint MUST continue processing
	//as normal (as defined by OAuth 2.0 [RFC6749]).  If the values are not
	//equal, an error response indicating "invalid_grant" as described in
	//Section 5.2 of [RFC6749] MUST be returned.
	switch method {
	case "S256":
		verifierLength := base64.RawURLEncoding.DecodedLen(len(verifier))

		// NOTE: The code verifier SHOULD have enough entropy to make it
		//	impractical to guess the value.  It is RECOMMENDED that the output of
		//	a suitable random number generator be used to create a 32-octet
		//	sequence.  The octet sequence is then base64url-encoded to produce a
		//	43-octet URL safe string to use as the code verifier.
		if verifierLength < 32 {
			return errors.WithStack(fosite.ErrInsufficientEntropy.
				WithHint("The PKCE code verifier must contain at least 32 octets."))
		}

		verifierBytes := make([]byte, verifierLength)
		if _, err := base64.RawURLEncoding.Decode(verifierBytes, []byte(verifier)); err != nil {
			return errors.WithStack(fosite.ErrInvalidGrant.WithHint("Unable to decode code_verifier using base64 url decoding without padding.").WithDebug(err.Error()))
		}

		hash := sha256.New()
		if _, err := hash.Write([]byte(verifier)); err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		}

		if base64.RawURLEncoding.EncodeToString(hash.Sum([]byte{})) != challenge {
			return errors.WithStack(fosite.ErrInvalidGrant.
				WithHint("The PKCE code challenge did not match the code verifier."))
		}
		break
	case "plain":
		fallthrough
	default:
		if verifier != challenge {
			return errors.WithStack(fosite.ErrInvalidGrant.
				WithHint("The PKCE code challenge did not match the code verifier."))
		}
	}

	return nil
}

func (c *Handler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	return nil
}
