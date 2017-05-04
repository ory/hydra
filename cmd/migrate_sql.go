// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	ladon "github.com/ory/ladon/manager/sql"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "migrate sql <database-url>",
	Short: "Create and migrate SQL schemas to work with this version",
	Long: `WARNING: Before running this command on an existing database, create a back up!`,
	Run: func(cmd *cobra.Command, args []string) {
		var db *sqlx.DB
		var logger = logrus.New()

		u, err := url.Parse(args[0])
		if err != nil {
			log.Fatalf("Could not parse DATABASE_URL: %s", err)
		}

		if err := pkg.Retry(logger, time.Second*15, time.Minute*2, func() error {
			if u.Scheme == "mysql" {
				args[0] = strings.Replace(args[0], "mysql://", "", -1)
			}

			if db, err = sqlx.Open(u.Scheme, u); err != nil {
				return errors.Errorf("Could not connect to SQL: %s", err)
			} else if err := db.Ping(); err != nil {
				return errors.Errorf("Could not connect to SQL: %s", err)
			}

			return nil
		}); err != nil {
			log.Fatalf("Could not connect to SQL: %s", err)
		}

		if err := (&client.SQLManager{DB: db}).CreateSchemas(); err != nil {
			c.GetLogger().Fatalf("Could not create client schema: %s", err)
		}

		if err := (&oauth2.FositeSQLStore{DB: db}).CreateSchemas(); err != nil {
			c.GetLogger().Fatalf("Could not create oauth2 schema: %s", err)
		}

		if err := (&jwk.SQLManager{DB: db}).CreateSchemas(); err != nil {
			c.GetLogger().Fatalf("Could not create jwk schema: %s", err)
		}

		if err := (&oauth2.FositeSQLStore{DB: db}).CreateSchemas(); err != nil {
			c.GetLogger().Fatalf("Could not create oauth2 schema: %s", err)
		}

		if err := ladon.NewSQLManager(db, nil).CreateSchemas(); err != nil {
			logrus.Fatalf("Could not create policy schema: %s", err)
		}

		if err := (&group.SQLManager{DB: db}).CreateSchemas(); err != nil {
			logrus.Fatalf("Could not create group schema: %s", err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(sqlCmd)
}
