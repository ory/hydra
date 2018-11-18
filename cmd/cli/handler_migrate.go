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
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/resilience"
)

type MigrateHandler struct {
	c *config.Config
}

func newMigrateHandler(c *config.Config) *MigrateHandler {
	return &MigrateHandler{c: c}
}

type schemaCreator interface {
	CreateSchemas() (int, error)
}

func (h *MigrateHandler) connectToSql(dsn string) (*sqlx.DB, error) {
	var db *sqlx.DB

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Errorf("could not parse DATABASE_URL: %s", err)
	}

	if err := resilience.Retry(h.c.GetLogger(), time.Second*15, time.Minute*2, func() error {
		if u.Scheme == "mysql" {
			dsn = strings.Replace(dsn, "mysql://", "", -1)
		}

		if db, err = sqlx.Open(u.Scheme, dsn); err != nil {
			return errors.Errorf("could not connect to SQL: %s", err)
		} else if err := db.Ping(); err != nil {
			return errors.Errorf("could not connect to SQL: %s", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func getDBUrl(cmd *cobra.Command, args []string, position int) (dburl string) {
	if flagx.MustGetBool(cmd, "read-from-env") {
		if len(viper.GetString("DATABASE_URL")) == 0 {
			fmt.Println(cmd.UsageString())
			fmt.Println("")
			fmt.Println("When using flag -e, environment variable DATABASE_URL must be set")
			return
		}
		dburl = viper.GetString("DATABASE_URL")
	} else {
		if len(args) <= position {
			fmt.Println(cmd.UsageString())
			return
		}
		dburl = args[position]
	}
	if dburl == "" {
		fmt.Println(cmd.UsageString())
		return
	}
	return
}

func (h *MigrateHandler) MigrateSecret(cmd *cobra.Command, args []string) {
	dburl := getDBUrl(cmd, args, 0)
	if dburl == "" {
		return
	}

	db, err := h.connectToSql(dburl)
	cmdx.Must(err, "An error occurred while connecting to SQL: %s", err)

	oldSecret := viper.GetString("OLD_SYSTEM_SECRET")
	newSecret := viper.GetString("NEW_SYSTEM_SECRET")

	if len(oldSecret) < 16 {
		cmdx.Fatalf("Value of environment variable OLD_SYSTEM_SECRET has to be at least 16 characters long but got: %d", len(oldSecret))
	}

	if len(newSecret) < 16 {
		cmdx.Fatalf("Value of environment variable NEW_SYSTEM_SECRET has to be at least 16 characters long but got: %d", len(oldSecret))
	}

	fmt.Println("Rotating encryption keys for JSON Web Key storage...")

	hashedOldSecret := pkg.HashStringSecret(oldSecret)
	hashedNewSecret := pkg.HashStringSecret(newSecret)
	manager := jwk.NewSQLManager(db, hashedOldSecret)
	err = manager.RotateKeys(context.TODO(), &jwk.AEAD{Key: hashedNewSecret})
	cmdx.Must(err, "Unable to rotate JSON Web Keys: %s\nAll changes have been rolled back.", err)

	fmt.Println("Rotating encryption keys for JSON Web Key storage completed successfully!")
	fmt.Printf(`You may now run ORY Hydra with the new system secret. If you wish that old OAuth 2.0 Access and Refres
tokens stay valid, please set environment variable ROTATED_SYSTEM_SECRET to the new secret:

ROTATED_SYSTEM_SECRET=%s hydra serve ...

If you wish that OAuth 2.0 Access and Refresh Tokens issued with the old secret are revoked, simply omit environment variable
ROTATED_SYSTEM_SECRET. This will NOT affect OpenID Connect ID Tokens!
`, newSecret)
}

func (h *MigrateHandler) MigrateSQL(cmd *cobra.Command, args []string) {
	dburl := getDBUrl(cmd, args, 0)
	if dburl == "" {
		return
	}

	db, err := h.connectToSql(dburl)
	cmdx.Must(err, "An error occurred while connecting to SQL: %s", err)

	err = h.runMigrateSQL(db)
	cmdx.Must(err, "An error occurred while running the migrations: %s", err)

	fmt.Println("Migration successful!")
}

func (h *MigrateHandler) runMigrateSQL(db *sqlx.DB) error {
	var total int
	migrators := map[string]schemaCreator{
		"client":  &client.SQLManager{DB: db},
		"oauth2":  &oauth2.FositeSQLStore{DB: db},
		"jwk":     &jwk.SQLManager{DB: db},
		"consent": consent.NewSQLManager(db, nil, nil),
	}
	for _, k := range []string{"jwk", "client", "consent", "oauth2"} {
		m := migrators[k]
		fmt.Printf("Applying `%s` SQL migrations...\n", k)
		if num, err := m.CreateSchemas(); err != nil {
			return errors.Wrapf(err, "could not apply %s SQL migrations", k)
		} else {
			fmt.Printf("Applied %d `%s` SQL migrations.\n", num, k)
			total += num
		}
	}

	fmt.Printf("Migration successful! Applied a total of %d SQL migrations.\n", total)
	return nil
}
