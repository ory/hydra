// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/integration/clients"
)

type authorizeJWTBearerRequiredJtiSuite struct {
	suite.Suite

	client *clients.JWTBearer
}

func (s *authorizeJWTBearerRequiredJtiSuite) TestBadResponseWithoutJTI() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, []string{"fosite"})

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerRequiredJtiSuite) TestSuccessResponseWithJTI() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       uuid.New().String(),
		},
	}, []string{"fosite"})

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerRequiredJtiSuite) getClient() *clients.JWTBearer {
	client := *s.client

	return &client
}

func (s *authorizeJWTBearerRequiredJtiSuite) assertSuccessResponse(t *testing.T, token *clients.Token, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, token)

	assert.Equal(t, token.TokenType, "bearer")
	assert.Empty(t, token.RefreshToken)
	assert.NotEmpty(t, token.ExpiresIn)
	assert.NotEmpty(t, token.AccessToken)
}

func (s *authorizeJWTBearerRequiredJtiSuite) assertBadResponse(t *testing.T, token *clients.Token, err error) {
	assert.Nil(t, token)
	assert.NotNil(t, err)

	retrieveError, ok := err.(*clients.RequestError)
	assert.True(t, ok)
	assert.Equal(t, retrieveError.Response.StatusCode, http.StatusBadRequest)
}

func TestAuthorizeJWTBearerRequiredJtiSuite(t *testing.T) {
	provider := compose.Compose(
		&fosite.Config{
			GrantTypeJWTBearerCanSkipClientAuth:  true,
			GrantTypeJWTBearerIDOptional:         false,
			GrantTypeJWTBearerIssuedDateOptional: true,
			TokenURL:                             tokenURL,
		},
		fositeStore,
		jwtStrategy,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.RFC7523AssertionGrantFactory,
	)
	testServer := mockServer(t, provider, &fosite.DefaultSession{})
	defer testServer.Close()

	client := newJWTBearerAppClient(testServer)
	if err := client.SetPrivateKey(firstKeyID, firstPrivateKey); err != nil {
		assert.Nil(t, err)
	}

	suite.Run(t, &authorizeJWTBearerRequiredJtiSuite{
		client: client,
	})
}
