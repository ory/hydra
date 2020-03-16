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

package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
)

var Migrations = map[string]*dbal.PackrMigrationSource{
	dbal.DriverMySQL:       dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/shared", "migrations/sql/mysql"}, true),
	dbal.DriverPostgreSQL:  dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/shared", "migrations/sql/postgres"}, true),
	dbal.DriverCockroachDB: dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/cockroach"}, true),
}

func NewSQLManager(db *sqlx.DB, r InternalRegistry) *SQLManager {
	return &SQLManager{
		r:  r,
		DB: db,
	}
}

type SQLManager struct {
	r  InternalRegistry
	DB *sqlx.DB
}

var sqlParams = []string{
	"id",
	"client_name",
	"client_secret",
	"redirect_uris",
	"grant_types",
	"response_types",
	"scope",
	"owner",
	"policy_uri",
	"tos_uri",
	"client_uri",
	"subject_type",
	"logo_uri",
	"contacts",
	"client_secret_expires_at",
	"sector_identifier_uri",
	"jwks",
	"jwks_uri",
	"token_endpoint_auth_method",
	"request_uris",
	"request_object_signing_alg",
	"userinfo_signed_response_alg",
	"allowed_cors_origins",
	"audience",
	"updated_at",
	"created_at",
	"frontchannel_logout_uri",
	"frontchannel_logout_session_required",
	"post_logout_redirect_uris",
	"backchannel_logout_uri",
	"backchannel_logout_session_required",
	"metadata",
}

func (m *SQLManager) PlanMigration(dbName string) ([]*migrate.PlannedMigration, error) {
	migrate.SetTable("hydra_client_migration")
	plan, _, err := migrate.PlanMigration(m.DB.DB, dbal.Canonicalize(m.DB.DriverName()), Migrations[dbName], migrate.Up, 0)
	return plan, errors.WithStack(err)
}

func (m *SQLManager) CreateSchemas(dbName string) (int, error) {
	migrate.SetTable("hydra_client_migration")
	n, err := migrate.Exec(m.DB.DB, dbal.Canonicalize(m.DB.DriverName()), Migrations[dbName], migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d Migrations", n)
	}
	return n, nil
}

func (m *SQLManager) GetConcreteClient(ctx context.Context, id string) (*Client, error) {
	var d Client
	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_client WHERE id=?"), id); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &d, nil
}

func (m *SQLManager) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return m.GetConcreteClient(ctx, id)
}

func (m *SQLManager) UpdateClient(ctx context.Context, c *Client) error {
	o, err := m.GetClient(ctx, c.GetID())
	if err != nil {
		return errors.WithStack(err)
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.r.ClientHasher().Hash(ctx, []byte(c.Secret))
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = string(h)
	}

	var query []string
	for _, param := range sqlParams {
		query = append(query, fmt.Sprintf("%s=:%s", param, param))
	}

	/* #nosec G201 - query is constructed using predefined variables only that are never modified */
	if _, err := m.DB.NamedExecContext(
		ctx,
		fmt.Sprintf(`UPDATE hydra_client SET %s WHERE id=:id`, strings.Join(query, ", ")),
		setDefaults(c),
	); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) Authenticate(ctx context.Context, id string, secret []byte) (*Client, error) {
	c, err := m.GetConcreteClient(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.r.ClientHasher().Compare(ctx, c.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (m *SQLManager) CreateClient(ctx context.Context, c *Client) error {
	h, err := m.r.ClientHasher().Hash(ctx, []byte(c.Secret))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = string(h)

	/* #nosec G201 - sqlParams is a "constant" array */
	if _, err := m.DB.NamedExecContext(ctx,
		fmt.Sprintf(
			"INSERT INTO hydra_client (%s) VALUES (%s)",
			strings.Join(sqlParams, ", "),
			":"+strings.Join(sqlParams, ", :"),
		),
		setDefaults(c),
	); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) DeleteClient(ctx context.Context, id string) error {
	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind(`DELETE FROM hydra_client WHERE id=?`), id); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) GetClients(ctx context.Context, limit, offset int) (clients []Client, err error) {
	clients = []Client{}

	if err := m.DB.SelectContext(ctx, &clients, m.DB.Rebind("SELECT * FROM hydra_client ORDER BY id LIMIT ? OFFSET ?"), limit, offset); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return clients, nil
}

func (m *SQLManager) CountClients(ctx context.Context) (int, error) {
	var n int
	if err := m.DB.QueryRow("SELECT count(*) FROM hydra_client").Scan(&n); err != nil {
		fmt.Println(err.Error())
		return 0, sqlcon.HandleError(err)
	}

	return n, nil
}

func setDefaults(c *Client) *Client {
	if c.JSONWebKeys == nil {
		c.JSONWebKeys = new(x.JoseJSONWebKeySet)
	}

	if c.Metadata == nil {
		c.Metadata = []byte("{}")
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}

	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now().UTC()
	}

	return c
}
