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
	"fmt"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/x/logrusx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

type MigrateHandler struct {
	d driver.Driver
}

func newMigrateHandler() *MigrateHandler {
	l := logrusx.New()
	c := configuration.NewViperProvider(l, false)
	return &MigrateHandler{
		d: driver.NewDefaultDriverWithoutValidation(
			driver.MustNewRegistry(c),
			c,
			l,
		),
	}
}

func (h *MigrateHandler) MigrateSQL(cmd *cobra.Command, args []string) {
	var dbu string
	if flagx.MustGetBool(cmd, "read-from-env") {
		dbu := h.d.Configuration().DSN()
		if len(dbu) == 0 {
			fmt.Println(cmd.UsageString())
			fmt.Println("")
			fmt.Println("When using flag -e, environment variable DATABASE_URL must be set")
			os.Exit(1)
			return
		}
	} else {
		if len(args) != 1 {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
			return
		}
		dbu = args[1]
	}

	viper.Set(configuration.ViperKeyDSN, dbu)

	reg, ok := h.d.Registry().(*driver.RegistrySQL)
	if !ok {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Printf("Migrations can only be executed against a SQL-compatible driver but DSN %s is not a SQL source.\n", dbu)
		os.Exit(1)
		return
	}

	n, err := reg.CreateSchemas()
	cmdx.Must(err, "An error occurred while connecting to SQL: %s", err)
	fmt.Printf("Successfully applied %d SQL migrations!\n", n)
}
