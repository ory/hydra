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

package fizzmigrate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ory/x/sqlcon"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/driver/configuration"
)

type MigrateHandler struct{}

func NewMigrateHandler() *MigrateHandler {
	return &MigrateHandler{}
}

func (h *MigrateHandler) MigrateSQL(cmd *cobra.Command, args []string) {
	l := logrusx.New()
	c := configuration.NewViperProvider(l, false, nil)
	scheme := sqlcon.GetDriverName(c.DSN())

	connection, err := sqlcon.NewSQLConnection(c.DSN(), l)
	if err != nil {
		fmt.Println(cmd.Usage())
		fmt.Println("")
		fmt.Printf("Could not create database connection: %s", err)
		os.Exit(1)
		return
	}

	db, err := connection.GetDatabaseRetry(time.Second*5, time.Minute*5)
	if err != nil {
		fmt.Println(cmd.Usage())
		fmt.Println("")
		fmt.Printf("Could not connect to DSN \"%s\": %s", c.DSN(), err)
		os.Exit(1)
		return
	}

	m := OldMigrationRunner{
		l,
		db,
	}

	plan, err := m.SchemaMigrationPlan(scheme)
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

	// Make create schemas with:
	// > check if table hydra_client_migration exists
	// > if yes -> run migration to fizz 2019010000000 + id
	// >  ..  INSERT INTO hydra_fizz_migrate VALUES CONCAT() (SELECT * FROM hydra_client_migration);
	n, err := m.CreateSchemas(scheme)
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
