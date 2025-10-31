// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

// #nosec:gosec G101 - False Positive
const jwtBearerGrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"

type JWTBearer struct {
	tokenURL string
	client   *http.Client

	Signer jose.Signer
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

type Header struct {
	Algorithm string `json:"alg"`
	Typ       string `json:"typ"`
	KeyID     string `json:"kid,omitempty"`
}

type JWTBearerPayload struct {
	*jwt.Claims

	PrivateClaims map[string]interface{}
}

func (c *JWTBearer) SetPrivateKey(keyID string, privateKey *rsa.PrivateKey) error {
	jwk := jose.JSONWebKey{Key: privateKey, KeyID: keyID, Algorithm: string(jose.RS256)}
	signingKey := jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       jwk,
	}
	signerOptions := &jose.SignerOptions{}
	signerOptions.WithType("JWT")

	sig, err := jose.NewSigner(signingKey, signerOptions)
	if err != nil {
		return err
	}

	c.Signer = sig

	return nil
}

func (c *JWTBearer) GetToken(ctx context.Context, payloadData *JWTBearerPayload, scope []string) (*Token, error) {
	builder := jwt.Signed(c.Signer).
		Claims(payloadData.Claims).
		Claims(payloadData.PrivateClaims)

	assertion, err := builder.CompactSerialize()
	if err != nil {
		return nil, err
	}

	requestBodyReader, err := c.getRequestBodyReader(assertion, scope)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", c.tokenURL, requestBodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if c := response.StatusCode; c < 200 || c > 299 {
		return nil, &RequestError{
			Response: response,
			Body:     body,
		}
	}

	token := &Token{}

	if err := json.Unmarshal(body, token); err != nil {
		return nil, err
	}

	return token, err
}

func (c *JWTBearer) getRequestBodyReader(assertion string, scope []string) (io.Reader, error) {
	data := url.Values{}
	data.Set("grant_type", jwtBearerGrantType)
	data.Set("assertion", string(assertion))

	if len(scope) != 0 {
		data.Set("scope", strings.Join(scope, " "))
	}

	return strings.NewReader(data.Encode()), nil
}

func NewJWTBearer(tokenURL string) *JWTBearer {
	return &JWTBearer{
		client:   &http.Client{},
		tokenURL: tokenURL,
	}
}
