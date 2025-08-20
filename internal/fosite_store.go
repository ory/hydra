// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
)

func AddFositeExamples(t *testing.T, r *driver.RegistrySQL) {
	for _, c := range []client.Client{
		{
			ID:            "my-client",
			Secret:        "foobar",
			RedirectURIs:  []string{"http://localhost:3846/callback"},
			ResponseTypes: []string{"id_token", "code", "token"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			Scope:         "fosite,openid,photos,offline",
		},
		{
			ID:            "encoded:client",
			Secret:        "encoded&password",
			RedirectURIs:  []string{"http://localhost:3846/callback"},
			ResponseTypes: []string{"id_token", "code", "token"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			Scope:         "fosite,openid,photos,offline",
		},
	} {
		// #nosec G601
		require.NoError(t, r.ClientManager().CreateClient(t.Context(), &c))
	}
}
