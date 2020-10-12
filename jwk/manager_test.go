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

package jwk_test

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/driver"

	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/jwk"
)

var testGenerator = &RS256Generator{}

func TestManager(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	registries := map[string]driver.Registry{
		"memory": internal.NewRegistryMemory(t, conf),
	}

	if !testing.Short() {
		registries["postgres"], registries["mysql"], registries["cockroach"], _ = internal.ConnectDatabases(t)
	}

	t.Run("TestManagerKey", func(t *testing.T) {
		ks, err := testGenerator.Generate("TestManagerKey", "sig")
		require.NoError(t, err)

		for name, r := range registries {
			t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKey(r.KeyManager(), ks, "TestManagerKey"))
		}
	})

	t.Run("TestManagerKeySet", func(t *testing.T) {
		ks, err := testGenerator.Generate("TestManagerKeySet", "sig")
		require.NoError(t, err)
		ks.Key("private")

		for name, r := range registries {
			t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKeySet(r.KeyManager(), ks, "TestManagerKeySet"))
		}
	})
}
