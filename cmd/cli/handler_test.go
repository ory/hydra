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

package cli

import (
	"flag"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/integration"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		db = integration.ConnectToPostgres()
	}

	code := m.Run()
	integration.KillAll()
	os.Exit(code)
}

func TestNewHandler(t *testing.T) {
	_ = NewHandler(&config.Config{})
}
