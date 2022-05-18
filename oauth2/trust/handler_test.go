package trust_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/oauth2/trust"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/jwk"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/x"
	"github.com/ory/x/urlx"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context.
type HandlerTestSuite struct {
	suite.Suite
	registry    driver.Registry
	server      *httptest.Server
	hydraClient *hydra.OryHydra
	publicKey   *rsa.PublicKey
}

// Setup will run before the tests in the suite are run.
func (s *HandlerTestSuite) SetupSuite() {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeySubjectTypesSupported, []string{"public"})
	conf.MustSet(config.KeyDefaultClientScope, []string{"foo", "bar"})
	s.registry = internal.NewRegistryMemory(s.T(), conf)

	router := x.NewRouterAdmin()
	handler := trust.NewHandler(s.registry)
	handler.SetRoutes(router)
	jwkHandler := jwk.NewHandler(s.registry, conf)
	jwkHandler.SetRoutes(router, x.NewRouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	s.server = httptest.NewServer(router)

	s.hydraClient = hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(s.server.URL).Host})
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
	internal.CleanAndMigrate(s.registry)(s.T())
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
	model := createRequestParams.Body

	createResult, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)

	s.Require().NoError(err, "no errors expected on grant creation")
	s.NotEmpty(createResult.Payload.ID, " grant id expected to be non-empty")
	s.Equal(*model.Issuer, createResult.Payload.Issuer, "issuer must match")
	s.Equal(model.Subject, createResult.Payload.Subject, "subject must match")
	s.Equal(model.Scope, createResult.Payload.Scope, "scopes must match")
	s.Equal(*model.Issuer, createResult.Payload.PublicKey.Set, "public key set must match grant issuer")
	s.Equal(*model.Jwk.Kid, createResult.Payload.PublicKey.Kid, "public key id must match")
	s.Equal(model.ExpiresAt.String(), createResult.Payload.ExpiresAt.String(), "expiration date must match")

	getRequestParams := admin.NewGetTrustedJwtGrantIssuerParams()
	getRequestParams.ID = createResult.Payload.ID
	getResult, err := s.hydraClient.Admin.GetTrustedJwtGrantIssuer(getRequestParams)

	s.Require().NoError(err, "no errors expected on grant fetching")
	s.Equal(getRequestParams.ID, getResult.Payload.ID, " grant id must match")
	s.Equal(*model.Issuer, getResult.Payload.Issuer, "issuer must match")
	s.Equal(model.Subject, getResult.Payload.Subject, "subject must match")
	s.Equal(model.Scope, getResult.Payload.Scope, "scopes must match")
	s.Equal(*model.Issuer, getResult.Payload.PublicKey.Set, "public key set must match grant issuer")
	s.Equal(*model.Jwk.Kid, getResult.Payload.PublicKey.Kid, "public key id must match")
	s.Equal(model.ExpiresAt.String(), getResult.Payload.ExpiresAt.String(), "expiration date must match")
}

func (s *HandlerTestSuite) TestGrantCanNotBeCreatedWithSameIssuerSubjectKey() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().NoError(err, "no errors expected on grant creation")

	_, err = s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().Error(err, "expected error, because grant with same issuer+subject+kid exists")

	kid := uuid.New().String()
	createRequestParams.Body.Jwk.Kid = &kid
	_, err = s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
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

	_, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().Error(err, "expected error, because a grant with a subject and allow_any_subject cannot be created")
}

func (s *HandlerTestSuite) TestGrantCanNotBeCreatedWithMissingFields() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().Error(err, "expected error, because grant missing issuer")

	createRequestParams = s.newCreateJwtBearerGrantParams(
		"ory",
		"",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	_, err = s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().Error(err, "expected error, because grant missing subject")

	createRequestParams = s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Time{},
	)

	_, err = s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
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
	model := createRequestParams.Body

	_, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().NoError(err, "no error expected on grant creation")

	getJWKRequestParams := admin.NewGetJSONWebKeyParams()
	getJWKRequestParams.Kid = *model.Jwk.Kid
	getJWKRequestParams.Set = *model.Issuer

	getResult, err := s.hydraClient.Admin.GetJSONWebKey(getJWKRequestParams)

	s.Require().NoError(err, "no error expected on fetching public key")
	s.Equal(*model.Jwk.Kid, *getResult.Payload.Keys[0].Kid)
}

func (s *HandlerTestSuite) TestGrantWithAnySubjectCanBeCreated() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"",
		true,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	grant, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().NoError(err, "no error expected on grant creation")

	assert.Empty(s.T(), grant.Payload.Subject)
	assert.Truef(s.T(), grant.Payload.AllowAnySubject, "grant with any subject must be true")
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

	_, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().NoError(err, "no errors expected on grant creation")

	_, err = s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams2)
	s.Require().NoError(err, "no errors expected on grant creation")

	getRequestParams := admin.NewListTrustedJwtGrantIssuersParams()
	getResult, err := s.hydraClient.Admin.ListTrustedJwtGrantIssuers(getRequestParams)

	s.Require().NoError(err, "no errors expected on grant list fetching")
	s.Len(getResult.Payload, 2, "expected to get list of 2 grants")

	getRequestParams.Issuer = createRequestParams2.Body.Issuer
	getResult, err = s.hydraClient.Admin.ListTrustedJwtGrantIssuers(getRequestParams)

	s.Require().NoError(err, "no errors expected on grant list fetching")
	s.Len(getResult.Payload, 1, "expected to get list of 1 grant, when filtering by issuer")
	s.Equal(*createRequestParams2.Body.Issuer, getResult.Payload[0].Issuer, "issuer must match")
}

func (s *HandlerTestSuite) TestGrantCanBeDeleted() {
	createRequestParams := s.newCreateJwtBearerGrantParams(
		"ory",
		"hackerman@example.com",
		false,
		[]string{"openid", "offline", "profile"},
		time.Now().Add(time.Hour),
	)

	createResult, err := s.hydraClient.Admin.TrustJwtGrantIssuer(createRequestParams)
	s.Require().NoError(err, "no errors expected on grant creation")

	deleteRequestParams := admin.NewDeleteTrustedJwtGrantIssuerParams()
	deleteRequestParams.ID = createResult.Payload.ID
	_, err = s.hydraClient.Admin.DeleteTrustedJwtGrantIssuer(deleteRequestParams)

	s.Require().NoError(err, "no errors expected on grant deletion")

	_, err = s.hydraClient.Admin.DeleteTrustedJwtGrantIssuer(deleteRequestParams)
	s.Error(err, "expected error, because grant has been already deleted")
}

func (s *HandlerTestSuite) generateJWK(publicKey *rsa.PublicKey) *models.JSONWebKey {
	jwk := jose.JSONWebKey{
		Key:       publicKey,
		KeyID:     uuid.New().String(),
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	b, err := jwk.MarshalJSON()
	s.Require().NoError(err)

	mJWK := &models.JSONWebKey{}
	err = mJWK.UnmarshalBinary(b)
	s.Require().NoError(err)

	return mJWK
}

func (s *HandlerTestSuite) newCreateJwtBearerGrantParams(
	issuer, subject string, allowAnySubject bool, scope []string, expiresAt time.Time,
) *admin.TrustJwtGrantIssuerParams {
	createRequestParams := admin.NewTrustJwtGrantIssuerParams()
	exp := strfmt.DateTime(expiresAt.UTC().Round(time.Second))
	model := &models.TrustJwtGrantIssuerBody{
		ExpiresAt:       &exp,
		Issuer:          &issuer,
		Jwk:             s.generateJWK(s.publicKey),
		Scope:           scope,
		Subject:         subject,
		AllowAnySubject: allowAnySubject,
	}
	createRequestParams.SetBody(model)

	return createRequestParams
}

func (s *HandlerTestSuite) generatePublicKey() *rsa.PublicKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	return &privateKey.PublicKey
}
