package x

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func CleanSQL(t *testing.T, db *sqlx.DB) {
	t.Logf("Cleaning up database: %s", db.DriverName())
	for _, tb := range []string{
		"hydra_oauth2_access",
		"hydra_oauth2_refresh",
		"hydra_oauth2_code",
		"hydra_oauth2_oidc",
		"hydra_oauth2_pkce",
		"hydra_oauth2_consent_request_handled",
		"hydra_oauth2_consent_request",
		"hydra_oauth2_authentication_request_handled",
		"hydra_oauth2_authentication_request",
		"hydra_oauth2_authentication_session",
		"hydra_oauth2_obfuscated_authentication_session",
		"hydra_oauth2_logout_request",
		"hydra_jwk",
		"hydra_client",
		// Migrations
		"hydra_oauth2_authentication_consent_migration",
		"hydra_client_migration",
		"hydra_oauth2_migration",
		"hydra_jwk_migration",
	} {
		if _, err := db.Exec("DROP TABLE IF EXISTS " + tb); err != nil {
			t.Logf(`Unable to clean up table "%s": %s`, tb, err)
		}
	}
	t.Logf("Successfully cleaned up database: %s", db.DriverName())
}
