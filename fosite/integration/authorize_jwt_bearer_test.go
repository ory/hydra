// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/integration/clients"
)

type authorizeJWTBearerSuite struct {
	suite.Suite

	client *clients.JWTBearer
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithRequiredParamsOnly() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, []string{"fosite"})

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithMultipleAudienceInAssertion() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL, "https://example.com/oauth"},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, []string{"fosite"})

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithMultipleScopesInRequest() {
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
	}, []string{"fosite", "gitlab"})

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithoutScopes() {
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
	}, nil)

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithExtraClaim() {
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
		PrivateClaims: map[string]interface{}{"extraClaim": "extraClaimValue"},
	}, []string{"fosite"})

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithNotBeforeClaim() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:    firstJWTBearerIssuer,
			Subject:   firstJWTBearerSubject,
			Audience:  []string{tokenURL},
			Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}, []string{"fosite"})

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseWithJTIClaim() {
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

func (s *authorizeJWTBearerSuite) TestSuccessResponse() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:    firstJWTBearerIssuer,
			Subject:   firstJWTBearerSubject,
			Audience:  []string{tokenURL, "example.com"},
			Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			ID:        uuid.New().String(),
		},
		PrivateClaims: map[string]interface{}{"random": "random"},
	}, nil)

	s.assertSuccessResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithExpiredJWT() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, []string{"fosite"})

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithExpiryMaxDuration() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, []string{"fosite"})

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithInvalidPrivateKey() {
	ctx := context.Background()
	client := s.getClient()
	wrongPrivateKey := secondPrivateKey

	if err := client.SetPrivateKey(firstKeyID, wrongPrivateKey); err != nil {
		assert.Nil(s.T(), err)
	}

	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithInvalidKeyID() {
	ctx := context.Background()
	client := s.getClient()

	if err := client.SetPrivateKey("wrongKeyID", firstPrivateKey); err != nil {
		assert.Nil(s.T(), err)
	}

	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithInvalidAudience() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{"https://example.com/oauth"},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseForSecondRequestWithSameJTI() {
	ctx := context.Background()
	client := s.getClient()
	config := &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       uuid.New().String(),
		},
	}

	client.GetToken(ctx, config, nil)
	token2, err := client.GetToken(ctx, config, nil)

	s.assertBadResponse(s.T(), token2, err)
}

func (s *authorizeJWTBearerSuite) TestSuccessResponseForSecondRequestWithSameJTIAfterFirstExpired() {
	ctx := context.Background()
	client := s.getClient()
	config := &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Second)),
			IssuedAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			ID:       uuid.New().String(),
		},
	}

	client.GetToken(ctx, config, nil)

	time.Sleep(time.Second)
	config.Expiry = jwt.NewNumericDate(time.Now().Add(time.Hour))

	token2, err := client.GetToken(ctx, config, nil)

	s.assertSuccessResponse(s.T(), token2, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithNotBeforeLaterThenIssueAt() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:    firstJWTBearerIssuer,
			Subject:   firstJWTBearerSubject,
			Audience:  []string{tokenURL},
			Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithoutSubject() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  "",
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithWrongSubject() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  "wrong_subject",
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithWrongIssuer() {
	ctx := context.Background()
	client := s.getClient()
	token, err := client.GetToken(ctx, &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   "wrong_issuer",
			Subject:  firstJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) TestBadResponseWithWrongScope() {
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
	}, []string{"fosite", "permission"})

	s.assertBadResponse(s.T(), token, err)
}

func (s *authorizeJWTBearerSuite) getClient() *clients.JWTBearer {
	client := *s.client

	return &client
}

func (s *authorizeJWTBearerSuite) assertSuccessResponse(t *testing.T, token *clients.Token, err error) {
	assert.Nil(t, err)
	require.NotNil(t, token)

	assert.Equal(t, token.TokenType, "bearer")
	assert.Empty(t, token.RefreshToken)
	assert.NotEmpty(t, token.ExpiresIn)
	assert.NotEmpty(t, token.AccessToken)
}

func (s *authorizeJWTBearerSuite) assertBadResponse(t *testing.T, token *clients.Token, err error) {
	assert.Nil(t, token)
	assert.NotNil(t, err)

	retrieveError, ok := err.(*clients.RequestError)
	assert.True(t, ok)
	assert.Equal(t, retrieveError.Response.StatusCode, http.StatusBadRequest)
}

func TestAuthorizeJWTBearerSuite(t *testing.T) {
	provider := compose.Compose(
		&fosite.Config{
			GrantTypeJWTBearerCanSkipClientAuth:  true,
			GrantTypeJWTBearerIDOptional:         true,
			GrantTypeJWTBearerIssuedDateOptional: true,
			GrantTypeJWTBearerMaxDuration:        24 * time.Hour,
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

	suite.Run(t, &authorizeJWTBearerSuite{
		client: client,
	})
}
