// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/ory/pop/v6"
)

func DeleteHydraRows(t *testing.T, c *pop.Connection) {
	for _, tb := range []string{
		"hydra_oauth2_access",
		"hydra_oauth2_refresh",
		"hydra_oauth2_code",
		"hydra_oauth2_oidc",
		"hydra_oauth2_pkce",
		"hydra_oauth2_device_auth_codes",
		"hydra_oauth2_flow",
		"hydra_oauth2_authentication_session",
		"hydra_oauth2_obfuscated_authentication_session",
		"hydra_oauth2_logout_request",
		"hydra_oauth2_jti_blacklist",
		"hydra_oauth2_trusted_jwt_bearer_issuer",
		"hydra_jwk",
		"hydra_client",
	} {
		if err := c.RawQuery("DELETE FROM " + tb).Exec(); err != nil {
			t.Logf(`Unable to delete rows in table "%s": %s`, tb, err)
		}
	}
}
