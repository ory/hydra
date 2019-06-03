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
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"
)

type MigrateHandler struct{}

func newMigrateHandler() *MigrateHandler {
	return &MigrateHandler{}
}

func (h *MigrateHandler) MigrateSQL(cmd *cobra.Command, args []string) {
	var d driver.Driver

	if flagx.MustGetBool(cmd, "read-from-env") {
		d = driver.NewDefaultDriver(logrusx.New(), false, nil, "", "", "", false)
		if len(d.Configuration().DSN()) == 0 {
			fmt.Println(cmd.UsageString())
			fmt.Println("")
			fmt.Println("When using flag -e, environment variable DSN must be set")
			os.Exit(1)
			return
		}
	} else {
		if len(args) != 1 {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
			return
		}
		viper.Set(configuration.ViperKeyDSN, args[0])
		d = driver.NewDefaultDriver(logrusx.New(), false, nil, "", "", "", false)
	}

	reg, ok := d.Registry().(*driver.RegistrySQL)
	if !ok {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Printf("Migrations can only be executed against a SQL-compatible driver but DSN is not a SQL source.\n")
		os.Exit(1)
		return
	}

	u, err := url.Parse(d.Configuration().DSN())
	if err != nil {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	plan, err := reg.SchemaMigrationPlan(u.Scheme)
	cmdx.Must(err, "An error occurred planning migrations: %s", err)

	fmt.Println("The following migration is planned:")
	fmt.Println("")
	plan.Render()

	if !flagx.MustGetBool(cmd, "yes") {
		fmt.Println("")
		fmt.Println("To skip the next question use flag --yes (at your own risk).")
		if !askForConfirmation("Do you wish to execute this migration plan?") {
			fmt.Println("Migration aborted.")
			return
		}
	}

	n, err := reg.CreateSchemas(u.Scheme)
	cmdx.Must(err, "An error occurred while connecting to SQL: %s", err)
	fmt.Printf("Successfully applied %d SQL migrations!\n", n)
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		cmdx.Must(err, "%s", err)

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
