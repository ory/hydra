/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package client_test

import (
	"fmt"
	"testing"

	"github.com/ory/hydra/internal"

	"github.com/ory/hydra/driver"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"

	. "github.com/ory/hydra/client"
)

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()

	registries := map[string]driver.Registry{
		"memory": internal.NewRegistryMemory(t, conf),
	}

	if !testing.Short() {
		// registries["postgres"], registries["mysql"], registries["cockroach"], _ = internal.ConnectDatabases(t)
		t.Log("connected")
	}

	for k, m := range registries {
		t.Run("case=create-get-update-delete", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperCreateGetUpdateDeleteClient(k, m.ClientManager()))
		})

		t.Run("case=autogenerate-key", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperClientAutoGenerateKey(k, m.ClientManager()))
		})

		t.Run("case=auth-client", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperClientAuthenticate(k, m.ClientManager()))
		})
	}
}
