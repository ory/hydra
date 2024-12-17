// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/go-jose/go-jose/v3"
	"github.com/tidwall/gjson"

	"github.com/ory/x/pointerx"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/x/contextx"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/jwk"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context.
type HandlerTestSuite struct {
	suite.Suite
	registry    driver.Registry
	server      *httptest.Server
	hydraClient *hydra.APIClient
	publicKey   *rsa.PublicKey
}

// Setup will run before the tests in the suite are run.
func (s *HandlerTestSuite) SetupSuite() {
	conf := testhelpers.NewConfigurationWithDefaults()
	conf.MustSet(context.Background(), config.KeySubjectTypesSupported, []string{"public"})
	conf.MustSet(context.Background(), config.KeyDefaultClientScope, []string{"foo", "bar"})
	s.registry = testhelpers.NewRegistryMemory(s.T(), conf, &contextx.Default{})

	router := x.NewRouterAdmin(conf.AdminURL)
	handler := trust.NewHandler(s.registry)
	handler.SetRoutes(router)
	jwkHandler := jwk.NewHandler(s.registry)
	jwkHandler.SetRoutes(router, x.NewRouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	s.server = httptest.NewServer(router)

	c := hydra.NewAPIClient(hydra.NewConfiguration())
	c.GetConfig().Servers = hydra.ServerConfigurations{{URL: s.server.URL}}
	s.hydraClient = c
	s.publicKey = s.generatePublicKey()
}

// Setup before each test.
func (s *HandlerTestSuite) SetupTest() {
}

// Will run after all the tests in the suite have been run.
func (s *HandlerTestSuite) TearDownSuite() {
}

// Will run after each test in the suite.
func (s *HandlerTestSuite) TearDownTest() {
	testhelpers.CleanAndMigrate(s.registry)(s.T())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestGrantCanBeCreatedAndFetched() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)
	model := createRequestParams

	ctx := context.Background()
	createResult, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(ctx).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().NoError(err, "no errors expected on grant creation")
	s.NotEmpty(createResult.Id, " grant id expected to be non-empty")
	s.Equal(model.Issuer, *createResult.Issuer, "issuer must match")
	s.Equal(*model.Subject, *createResult.Subject, "subject must match")
	s.Equal(model.Scope, createResult.Scope, "scopes must match")
	s.Equal(model.Issuer, *createResult.PublicKey.Set, "public key set must match grant issuer")
	s.Equal(model.Jwk.Kid, *createResult.PublicKey.Kid, "public key id must match")
	s.Equal(model.ExpiresAt.Round(time.Second).UTC().String(), createResult.ExpiresAt.Round(time.Second).UTC().String(), "expiration date must match")

	getResult, _, err := s.hydraClient.OAuth2API.GetTrustedOAuth2JwtGrantIssuer(ctx, *createResult.Id).Execute()
	s.Require().NoError(err, "no errors expected on grant fetching")
	s.Equal(*createResult.Id, *getResult.Id, " grant id must match")
	s.Equal(model.Issuer, *getResult.Issuer, "issuer must match")
	s.Equal(*model.Subject, *getResult.Subject, "subject must match")
	s.Equal(model.Scope, getResult.Scope, "scopes must match")
	s.Equal(model.Issuer, *getResult.PublicKey.Set, "public key set must match grant issuer")
	s.Equal(model.Jwk.Kid, *getResult.PublicKey.Kid, "public key id must match")
	s.Equal(model.ExpiresAt.Round(time.Second).UTC().String(), getResult.ExpiresAt.Round(time.Second).UTC().String(), "expiration date must match")
}

func (s *HandlerTestSuite) TestGrantCanNotBeCreatedWithSameIssuerSubjectKey() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	ctx := context.Background()
	_, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(ctx).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().NoError(err, "no errors expected on grant creation")

	_, _, err = s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(ctx).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().Error(err, "expected error, because grant with same issuer+subject+kid exists")

	kid := uuid.New().String()
	createRequestParams.Jwk.Kid = kid
	_, _, err = s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(ctx).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.NoError(err, "no errors expected on grant creation, because kid is now different")
}

func (s *HandlerTestSuite) TestGrantCanNotBeCreatedWithSubjectAndAnySubject() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		true,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().Error(err, "expected error, because a grant with a subject and allow_any_subject cannot be created")
}

func (s *HandlerTestSuite) TestGrantCanNotBeCreatedWithUnknownJWK() {
	createRequestParams := hydra.TrustOAuth2JwtGrantIssuer{
		AllowAnySubject: pointerx.Ptr(true),
		ExpiresAt:       time.Now().Add(1 * time.Hour),
		Issuer:          "ory",
		Jwk: hydra.JsonWebKey{
			Alg: "unknown",
		},
		Scope: []string{"openid", "offline", "profile"},
	}

	_, res, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Assert().Equal(http.StatusBadRequest, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	s.Contains(gjson.GetBytes(body, "error_description").String(), "unknown json web key type")
	s.Require().Error(err, "expected error, because the key type was unknown")
}

func (s *HandlerTestSuite) TestGrantCanNotBeCreatedWithMissingFields() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().Error(err, "expected error, because grant missing issuer")

	createRequestParams = s.newCreateJwtBearerGrantParams(
		"ory",
		"",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, _, err = s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().Error(err, "expected error, because grant missing subject")

	createRequestParams = s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Time{},
	)

	_, _, err = s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Error(err, "expected error, because grant missing expiration date")
}

func (s *HandlerTestSuite) TestGrantPublicCanBeFetched() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().NoError(err, "no error expected on grant creation")

	getResult, _, err := s.hydraClient.JwkAPI.GetJsonWebKey(context.Background(), createRequestParams.Issuer, createRequestParams.Jwk.Kid).Execute()

	s.Require().NoError(err, "no error expected on fetching public key")
	s.Equal(createRequestParams.Jwk.Kid, getResult.Keys[0].Kid)
}

func (s *HandlerTestSuite) TestGrantWithAnySubjectCanBeCreated() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"",
		true,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	grant, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().NoError(err, "no error expected on grant creation")

	assert.Empty(s.T(), grant.Subject)
	assert.Truef(s.T(), *grant.AllowAnySubject, "grant with any subject must be true")
}

func (s *HandlerTestSuite) TestGrantListCanBeFetched() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)
	createRequestParams2 := s.newCreateJwtBearerGrantParams(
		"ory2",
		"safetyman@example.com",
		false,
		[]string{"profile"},
		time.Now().Add(time.Hour),
	)

	_, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().NoError(err, "no errors expected on grant creation")

	_, _, err = s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams2).Execute()
	s.Require().NoError(err, "no errors expected on grant creation")

	getResult, _, err := s.hydraClient.OAuth2API.ListTrustedOAuth2JwtGrantIssuers(context.Background()).Execute()
	s.Require().NoError(err, "no errors expected on grant list fetching")
	s.Len(getResult, 2, "expected to get list of 2 grants")

	getResult, _, err = s.hydraClient.OAuth2API.ListTrustedOAuth2JwtGrantIssuers(context.Background()).Issuer(createRequestParams2.Issuer).Execute()

	s.Require().NoError(err, "no errors expected on grant list fetching")
	s.Len(getResult, 1, "expected to get list of 1 grant, when filtering by issuer")
	s.Equal(createRequestParams2.Issuer, *getResult[0].Issuer, "issuer must match")
}

func (s *HandlerTestSuite) TestGrantCanBeDeleted() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	createResult, _, err := s.hydraClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(createRequestParams).Execute()
	s.Require().NoError(err, "no errors expected on grant creation")

	_, err = s.hydraClient.OAuth2API.DeleteTrustedOAuth2JwtGrantIssuer(context.Background(), *createResult.Id).Execute()
	s.Require().NoError(err, "no errors expected on grant deletion")

	_, err = s.hydraClient.OAuth2API.DeleteTrustedOAuth2JwtGrantIssuer(context.Background(), *createResult.Id).Execute()
	s.Error(err, "expected error, because grant has been already deleted")
}

func (s *HandlerTestSuite) generateJWK(publicKey *rsa.PublicKey) hydra.JsonWebKey {
	var b bytes.Buffer
	s.Require().NoError(json.NewEncoder(&b).Encode(&jose.JSONWebKey{
		Key:       publicKey,
		KeyID:     uuid.New().String(),
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}))

	var mJWK hydra.JsonWebKey
	s.Require().NoError(json.NewDecoder(&b).Decode(&mJWK))
	return mJWK
}

func (s *HandlerTestSuite) newCreateJwtBearerGrantParams(
	issuer, subject string, allowAnySubject bool, scope []string, expiresAt time.Time,
) hydra.TrustOAuth2JwtGrantIssuer {
	return hydra.TrustOAuth2JwtGrantIssuer{
		ExpiresAt:       expiresAt,
		Issuer:          issuer,
		Jwk:             s.generateJWK(s.publicKey),
		Scope:           scope,
		Subject:         pointerx.String(subject),
		AllowAnySubject: pointerx.Bool(allowAnySubject),
	}
}

func (s *HandlerTestSuite) generatePublicKey() *rsa.PublicKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	return &privateKey.PublicKey
}
