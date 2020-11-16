package migratest

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	sqlPersister "github.com/ory/hydra/persistence/sql"
	"github.com/ory/x/sqlxx"
)

func assertEqualClients(t *testing.T, expected, actual *client.Client) {
	now := time.Now()
	expected.CreatedAt = now
	expected.UpdatedAt = now
	actual.CreatedAt = now
	actual.UpdatedAt = now

	assert.Equal(t, expected, actual)
}

func assertEqualJWKs(t *testing.T, expected, actual *jwk.SQLData) {
	now := time.Now()
	expected.CreatedAt = now
	actual.CreatedAt = now

	assert.Equal(t, expected, actual)
}

func assertEqualConsentRequests(t *testing.T, expected, actual *consent.ConsentRequest) {
	now := time.Now()
	expected.AuthenticatedAt = sqlxx.NullTime(now)
	expected.RequestedAt = now
	actual.AuthenticatedAt = sqlxx.NullTime(now)
	actual.RequestedAt = now

	assert.NotZero(t, actual.ClientID)
	actual.ClientID = ""
	assert.NotNil(t, actual.Client)
	actual.Client = nil

	assert.Equal(t, expected, actual)
}

func assertEqualLoginRequests(t *testing.T, expected, actual *consent.LoginRequest) {
	now := time.Now()
	expected.AuthenticatedAt = sqlxx.NullTime(now)
	expected.RequestedAt = now
	actual.AuthenticatedAt = sqlxx.NullTime(now)
	actual.RequestedAt = now

	assert.NotZero(t, actual.ClientID)
	actual.ClientID = ""
	assert.NotNil(t, actual.Client)
	actual.Client = nil

	assert.Equal(t, expected, actual)
}

func assertEqualLoginSessions(t *testing.T, expected, actual *consent.LoginSession) {
	now := time.Now()
	expected.AuthenticatedAt = sqlxx.NullTime(now)
	actual.AuthenticatedAt = sqlxx.NullTime(now)

	assert.Equal(t, expected, actual)
}

func assertEqualHandledConsentRequests(t *testing.T, expected, actual *consent.HandledConsentRequest) {
	now := time.Now()
	expected.AuthenticatedAt = sqlxx.NullTime(now)
	expected.RequestedAt = now
	actual.AuthenticatedAt = sqlxx.NullTime(now)
	actual.RequestedAt = now
	actual.HandledAt = sqlxx.NullTime{}

	assert.Equal(t, expected, actual)
}

func assertEqualHandledLoginRequests(t *testing.T, expected, actual *consent.HandledLoginRequest) {
	now := time.Now()
	expected.AuthenticatedAt = sqlxx.NullTime(now)
	expected.RequestedAt = now
	actual.AuthenticatedAt = sqlxx.NullTime(now)
	actual.RequestedAt = now

	assert.Equal(t, expected, actual)
}

func assertEqualLogoutRequests(t *testing.T, expected, actual *consent.LogoutRequest) {
	assert.NotZero(t, actual.ClientID)
	actual.ClientID = sql.NullString{}

	assert.Equal(t, expected, actual)
}

func assertEqualForcedObfucscatedLoginSessions(t *testing.T, expected, actual *consent.ForcedObfuscatedLoginSession) {
	assert.NotNil(t, actual.ClientID)
	actual.ClientID = ""
	assert.Equal(t, expected, actual)
}

func assertEqualOauth2Data(t *testing.T, expected, actual *sqlPersister.OAuth2RequestSQL) {
	now := time.Now()
	expected.RequestedAt = now
	actual.RequestedAt = now

	assert.NotZero(t, actual.Client)
	actual.Client = ""
	if expected.ConsentChallenge.Valid {
		assert.NotZero(t, actual.ConsentChallenge, "%+v", actual)
	}
	expected.ConsentChallenge = sql.NullString{}
	actual.ConsentChallenge = sql.NullString{}

	assert.Equal(t, expected, actual)
}

func assertEqualOauth2BlacklistedJTIs(t *testing.T, expected, actual *oauth2.BlacklistedJTI) {
	now := time.Now()
	expected.Expiry = now
	actual.Expiry = now

	assert.Equal(t, expected, actual)
}
