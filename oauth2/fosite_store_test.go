// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/contextx"
	"github.com/ory/x/sqlcon/dockertest"
)

func TestMain(m *testing.M) {
	flag.Parse()

	defer dockertest.KillAllTestDatabases()
	m.Run()
}

var registries = make(map[string]driver.Registry)
var cleanRegistries = func(t *testing.T) {
	registries["memory"] = testhelpers.NewRegistryMemory(t, testhelpers.NewConfigurationWithDefaults(), &contextx.Default{})
}

// returns clean registries that can safely be used for one test
// to reuse call cleanRegistries
func setupRegistries(t *testing.T) {
	if len(registries) == 0 && !testing.Short() {
		// first time called and sql tests
		var cleanSQL func(*testing.T)
		registries["postgres"], registries["mysql"], registries["cockroach"], cleanSQL = testhelpers.ConnectDatabases(t, false, &contextx.Default{})
		cleanMem := cleanRegistries
		cleanMem(t)
		cleanRegistries = func(t *testing.T) {
			cleanMem(t)
			cleanSQL(t)
		}
	} else {
		// reset all/init mem
		cleanRegistries(t)
	}
}

func TestManagers(t *testing.T) {
	setupRegistries(t)

	ctx := context.Background()
	tests := []struct {
		name                   string
		enableSessionEncrypted bool
	}{
		{
			name:                   "DisableSessionEncrypted",
			enableSessionEncrypted: false,
		},
		{
			name:                   "EnableSessionEncrypted",
			enableSessionEncrypted: true,
		},
	}
	for _, tc := range tests {
		t.Run("suite="+tc.name, func(t *testing.T) {
			for k, r := range registries {
				t.Run("database="+k, func(t *testing.T) {
					store := testhelpers.NewRegistrySQLFromURL(t, r.Config().DSN(), true, &contextx.Default{})
					store.Config().MustSet(ctx, config.KeyEncryptSessionData, tc.enableSessionEncrypted)

					if k != "memory" {
						t.Run("testHelperUniqueConstraints", testHelperRequestIDMultiples(store, k))
						t.Run("case=testFositeSqlStoreTransactionsCommitAccessToken", testFositeSqlStoreTransactionCommitAccessToken(store))
						t.Run("case=testFositeSqlStoreTransactionsRollbackAccessToken", testFositeSqlStoreTransactionRollbackAccessToken(store))
						t.Run("case=testFositeSqlStoreTransactionCommitRefreshToken", testFositeSqlStoreTransactionCommitRefreshToken(store))
						t.Run("case=testFositeSqlStoreTransactionRollbackRefreshToken", testFositeSqlStoreTransactionRollbackRefreshToken(store))
						t.Run("case=testFositeSqlStoreTransactionCommitAuthorizeCode", testFositeSqlStoreTransactionCommitAuthorizeCode(store))
						t.Run("case=testFositeSqlStoreTransactionRollbackAuthorizeCode", testFositeSqlStoreTransactionRollbackAuthorizeCode(store))
						t.Run("case=testFositeSqlStoreTransactionCommitPKCERequest", testFositeSqlStoreTransactionCommitPKCERequest(store))
						t.Run("case=testFositeSqlStoreTransactionRollbackPKCERequest", testFositeSqlStoreTransactionRollbackPKCERequest(store))
						t.Run("case=testFositeSqlStoreTransactionCommitOpenIdConnectSession", testFositeSqlStoreTransactionCommitOpenIdConnectSession(store))
						t.Run("case=testFositeSqlStoreTransactionRollbackOpenIdConnectSession", testFositeSqlStoreTransactionRollbackOpenIdConnectSession(store))
					}

					t.Run("testHelperCreateGetDeleteAuthorizeCodes", testHelperCreateGetDeleteAuthorizeCodes(store))
					t.Run("testHelperExpiryFields", testHelperExpiryFields(store))
					t.Run("testHelperCreateGetDeleteAccessTokenSession", testHelperCreateGetDeleteAccessTokenSession(store))
					t.Run("testHelperNilAccessToken", testHelperNilAccessToken(store))
					t.Run("testHelperCreateGetDeleteOpenIDConnectSession", testHelperCreateGetDeleteOpenIDConnectSession(store))
					t.Run("testHelperCreateGetDeleteRefreshTokenSession", testHelperCreateGetDeleteRefreshTokenSession(store))
					t.Run("testHelperRevokeRefreshToken", testHelperRevokeRefreshToken(store))
					t.Run("testHelperCreateGetDeletePKCERequestSession", testHelperCreateGetDeletePKCERequestSession(store))
					t.Run("testHelperFlushTokens", testHelperFlushTokens(store, time.Hour))
					t.Run("testHelperFlushTokensWithLimitAndBatchSize", testHelperFlushTokensWithLimitAndBatchSize(store, 3, 2))
					t.Run("testFositeStoreSetClientAssertionJWT", testFositeStoreSetClientAssertionJWT(store))
					t.Run("testFositeStoreClientAssertionJWTValid", testFositeStoreClientAssertionJWTValid(store))
					t.Run("testHelperDeleteAccessTokens", testHelperDeleteAccessTokens(store))
					t.Run("testHelperRevokeAccessToken", testHelperRevokeAccessToken(store))
					t.Run("testFositeJWTBearerGrantStorage", testFositeJWTBearerGrantStorage(store))
					t.Run("testHelperRotateRefreshToken", testHelperRotateRefreshToken(store))
				})
			}
		})
	}
}
