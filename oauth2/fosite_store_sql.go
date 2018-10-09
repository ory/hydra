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

package oauth2

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/x/sqlcon"
)

type FositeSQLStore struct {
	client.Manager
	DB                  *sqlx.DB
	L                   logrus.FieldLogger
	AccessTokenLifespan time.Duration
	HashSignature       bool
}

func NewFositeSQLStore(m client.Manager,
	db *sqlx.DB,
	l logrus.FieldLogger,
	accessTokenLifespan time.Duration,
	hashSignature bool,
) *FositeSQLStore {
	return &FositeSQLStore{
		Manager:             m,
		L:                   l,
		DB:                  db,
		AccessTokenLifespan: accessTokenLifespan,
		HashSignature:       hashSignature,
	}
}

func sqlSchemaUp(db, table, id string) string {
	shared := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS hydra_oauth2_%s (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		text NOT NULL,
	scope  			text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		text NOT NULL,
	session_data  	text NOT NULL
)`, table),
		fmt.Sprintf("ALTER TABLE hydra_oauth2_%s ADD subject varchar(255) NOT NULL DEFAULT ''", table),
		`CREATE TABLE IF NOT EXISTS hydra_oauth2_pkce (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		text NOT NULL,
	scope  			text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		text NOT NULL,
	session_data  	text NOT NULL,
	subject 		varchar(255) NOT NULL
)`,
		fmt.Sprintf("ALTER TABLE hydra_oauth2_%s ADD active BOOL NOT NULL DEFAULT TRUE", table),
		fmt.Sprintf("CREATE UNIQUE INDEX hydra_oauth2_%s_request_id_idx ON hydra_oauth2_%s (request_id)", table, table),
		fmt.Sprintf("CREATE INDEX hydra_oauth2_%s_requested_at_idx ON hydra_oauth2_%s (requested_at)", table, table),
	}

	m := map[string]map[string]string{
		"mysql": {
			"1": shared[0],
			"2": shared[1],
			"3": shared[2],
			"4": shared[3],
			"5": shared[4],
			"6": shared[5],
		},
		"postgres": {
			"1": shared[0],
			"2": shared[1],
			"3": shared[2],
			"4": shared[3],
			"5": shared[4],
			"6": shared[5],
		},
	}

	return m[db][id]
}

func sqlSchemaDown(db, table, id string) string {
	shared := []string{
		fmt.Sprintf(`DROP TABLE %s)`, table),
		fmt.Sprintf("ALTER TABLE hydra_oauth2_%s DROP COLUMN subject", table),
		"DROP TABLE hydra_oauth2_pkce",
		fmt.Sprintf("ALTER TABLE hydra_oauth2_%s DROP COLUMN active", table),
		fmt.Sprintf("DROP INDEX hydra_oauth2_%s_request_id_idx ON hydra_oauth2_%s", table, table),
		fmt.Sprintf("DROP INDEX hydra_oauth2_%s_requested_at_idx ON hydra_oauth2_%s", table, table),
	}

	m := map[string]map[string]string{
		"mysql": {
			"1": shared[0],
			"2": shared[1],
			"3": shared[2],
			"4": shared[3],
			"5": shared[4],
			"6": shared[5],
		},
		"postgres": {
			"1": shared[0],
			"2": shared[1],
			"3": shared[2],
			"4": shared[3],
			"5": shared[4],
			"6": shared[5],
		},
	}

	return m[db][id]
}

const (
	sqlTableOpenID  = "oidc"
	sqlTableAccess  = "access"
	sqlTableRefresh = "refresh"
	sqlTableCode    = "code"
	sqlTablePKCE    = "pkce"
)

func createMigrationSource(db string) *migrate.MemoryMigrationSource {
	return &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "1",
				Up: []string{
					sqlSchemaUp(db, sqlTableAccess, "1"),
					sqlSchemaUp(db, sqlTableRefresh, "1"),
					sqlSchemaUp(db, sqlTableCode, "1"),
					sqlSchemaUp(db, sqlTableOpenID, "1"),
				},
				Down: []string{
					sqlSchemaDown(db, sqlTableAccess, "1"),
					sqlSchemaDown(db, sqlTableRefresh, "1"),
					sqlSchemaDown(db, sqlTableCode, "1"),
					sqlSchemaDown(db, sqlTableOpenID, "1"),
				},
			},
			{
				Id: "2",
				Up: []string{
					sqlSchemaUp(db, sqlTableAccess, "2"),
					sqlSchemaUp(db, sqlTableRefresh, "2"),
					sqlSchemaUp(db, sqlTableCode, "2"),
					sqlSchemaUp(db, sqlTableOpenID, "2"),
				},
				Down: []string{
					sqlSchemaDown(db, sqlTableAccess, "2"),
					sqlSchemaDown(db, sqlTableRefresh, "2"),
					sqlSchemaDown(db, sqlTableCode, "2"),
					sqlSchemaDown(db, sqlTableOpenID, "2"),
				},
			},
			{
				Id: "3",
				Up: []string{
					sqlSchemaUp(db, sqlTablePKCE, "3"),
				},
				Down: []string{
					sqlSchemaDown(db, sqlTablePKCE, "3"),
				},
			},
			{
				Id: "4",
				Up: []string{
					sqlSchemaUp(db, sqlTableAccess, "4"),
					sqlSchemaUp(db, sqlTableRefresh, "4"),
					sqlSchemaUp(db, sqlTableCode, "4"),
					sqlSchemaUp(db, sqlTableOpenID, "4"),
					sqlSchemaUp(db, sqlTablePKCE, "4"),
				},
				Down: []string{
					sqlSchemaDown(db, sqlTableAccess, "4"),
					sqlSchemaDown(db, sqlTableRefresh, "4"),
					sqlSchemaDown(db, sqlTableCode, "4"),
					sqlSchemaDown(db, sqlTableOpenID, "4"),
					sqlSchemaDown(db, sqlTablePKCE, "4"),
				},
			},
			{
				Id: "5",
				Up: []string{
					sqlSchemaUp(db, sqlTableAccess, "5"),
					sqlSchemaUp(db, sqlTableRefresh, "5"),
				},
				Down: []string{
					sqlSchemaDown(db, sqlTableAccess, "5"),
					sqlSchemaDown(db, sqlTableRefresh, "5"),
				},
			},
			{
				Id: "6",
				Up: []string{
					sqlSchemaUp(db, sqlTableAccess, "6"),
				},
				Down: []string{
					sqlSchemaDown(db, sqlTableAccess, "6"),
				},
			},
		},
	}
}

var migrations = map[string]*migrate.MemoryMigrationSource{
	"mysql":    createMigrationSource("mysql"),
	"postgres": createMigrationSource("postgres"),
}

var sqlParams = []string{
	"signature",
	"request_id",
	"requested_at",
	"client_id",
	"scope",
	"granted_scope",
	"form_data",
	"session_data",
	"subject",
	"active",
}

type sqlData struct {
	PK            int       `db:"pk"`
	Signature     string    `db:"signature"`
	Request       string    `db:"request_id"`
	RequestedAt   time.Time `db:"requested_at"`
	Client        string    `db:"client_id"`
	Scopes        string    `db:"scope"`
	GrantedScopes string    `db:"granted_scope"`
	Form          string    `db:"form_data"`
	Subject       string    `db:"subject"`
	Active        bool      `db:"active"`
	Session       []byte    `db:"session_data"`
}

func sqlSchemaFromRequest(signature string, r fosite.Requester, logger logrus.FieldLogger) (*sqlData, error) {
	subject := ""
	if r.GetSession() == nil {
		logger.Debugf("Got an empty session in sqlSchemaFromRequest")
	} else {
		subject = r.GetSession().GetSubject()
	}

	session, err := json.Marshal(r.GetSession())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &sqlData{
		Request:       r.GetID(),
		Signature:     signature,
		RequestedAt:   r.GetRequestedAt(),
		Client:        r.GetClient().GetID(),
		Scopes:        strings.Join([]string(r.GetRequestedScopes()), "|"),
		GrantedScopes: strings.Join([]string(r.GetGrantedScopes()), "|"),
		Form:          r.GetRequestForm().Encode(),
		Session:       session,
		Subject:       subject,
		Active:        true,
	}, nil
}

func (s *sqlData) toRequest(session fosite.Session, cm client.Manager, logger logrus.FieldLogger) (*fosite.Request, error) {
	if session != nil {
		if err := json.Unmarshal(s.Session, session); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		logger.Debugf("Got an empty session in toRequest")
	}

	c, err := cm.GetClient(context.Background(), s.Client)
	if err != nil {
		return nil, err
	}

	val, err := url.ParseQuery(s.Form)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := &fosite.Request{
		ID:            s.Request,
		RequestedAt:   s.RequestedAt,
		Client:        c,
		Scopes:        fosite.Arguments(strings.Split(s.Scopes, "|")),
		GrantedScopes: fosite.Arguments(strings.Split(s.GrantedScopes, "|")),
		Form:          val,
		Session:       session,
	}

	return r, nil
}

// hashSignature prevents errors where the signature is longer than 128 characters (and thus doesn't fit into the pk).
func (s *FositeSQLStore) hashSignature(signature, table string) string {
	if table == sqlTableAccess && s.HashSignature {
		return fmt.Sprintf("%x", sha512.Sum384([]byte(signature)))
	}
	return signature
}

func (s *FositeSQLStore) createSession(ctx context.Context, signature string, requester fosite.Requester, table string) error {
	signature = s.hashSignature(signature, table)

	data, err := sqlSchemaFromRequest(signature, requester, s.L)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO hydra_oauth2_%s (%s) VALUES (%s)",
		table,
		strings.Join(sqlParams, ", "),
		":"+strings.Join(sqlParams, ", :"),
	)
	if _, err := s.DB.NamedExecContext(ctx, query, data); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (s *FositeSQLStore) findSessionBySignature(ctx context.Context, signature string, session fosite.Session, table string) (fosite.Requester, error) {
	signature = s.hashSignature(signature, table)

	var d sqlData
	if err := s.DB.GetContext(ctx, &d, s.DB.Rebind(fmt.Sprintf("SELECT * FROM hydra_oauth2_%s WHERE signature=?", table)), signature); err == sql.ErrNoRows {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	} else if !d.Active && table == sqlTableCode {
		if r, err := d.toRequest(session, s.Manager, s.L); err != nil {
			return nil, err
		} else {
			return r, errors.WithStack(fosite.ErrInvalidatedAuthorizeCode)
		}
	} else if !d.Active {
		return nil, errors.WithStack(fosite.ErrInactiveToken)
	}

	return d.toRequest(session, s.Manager, s.L)
}

func (s *FositeSQLStore) deleteSession(ctx context.Context, signature string, table string) error {
	signature = s.hashSignature(signature, table)

	if _, err := s.DB.ExecContext(ctx, s.DB.Rebind(fmt.Sprintf("DELETE FROM hydra_oauth2_%s WHERE signature=?", table)), signature); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (s *FositeSQLStore) CreateSchemas() (int, error) {
	database := s.DB.DriverName()
	switch database {
	case "pgx", "pq":
		database = "postgres"
	}

	migrate.SetTable("hydra_oauth2_migration")
	n, err := migrate.Exec(s.DB.DB, s.DB.DriverName(), migrations[database], migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (s *FositeSQLStore) CreateOpenIDConnectSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(ctx, signature, requester, sqlTableOpenID)
}

func (s *FositeSQLStore) GetOpenIDConnectSession(ctx context.Context, signature string, requester fosite.Requester) (fosite.Requester, error) {
	return s.findSessionBySignature(ctx, signature, requester.GetSession(), sqlTableOpenID)
}

func (s *FositeSQLStore) DeleteOpenIDConnectSession(ctx context.Context, signature string) error {
	return s.deleteSession(ctx, signature, sqlTableOpenID)
}

func (s *FositeSQLStore) CreateAuthorizeCodeSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(ctx, signature, requester, sqlTableCode)
}

func (s *FositeSQLStore) GetAuthorizeCodeSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(ctx, signature, session, sqlTableCode)
}

func (s *FositeSQLStore) InvalidateAuthorizeCodeSession(ctx context.Context, signature string) error {
	if _, err := s.DB.ExecContext(ctx, s.DB.Rebind(fmt.Sprintf(
		"UPDATE hydra_oauth2_%s SET active=false WHERE signature=?",
		sqlTableCode,
	)), signature); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (s *FositeSQLStore) CreateAccessTokenSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(ctx, signature, requester, sqlTableAccess)
}

func (s *FositeSQLStore) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(ctx, signature, session, sqlTableAccess)
}

func (s *FositeSQLStore) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	return s.deleteSession(ctx, signature, sqlTableAccess)
}

func (s *FositeSQLStore) CreateRefreshTokenSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(ctx, signature, requester, sqlTableRefresh)
}

func (s *FositeSQLStore) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(ctx, signature, session, sqlTableRefresh)
}

func (s *FositeSQLStore) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	return s.deleteSession(ctx, signature, sqlTableRefresh)
}

func (s *FositeSQLStore) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(ctx, signature, requester, sqlTablePKCE)
}

func (s *FositeSQLStore) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(ctx, signature, session, sqlTablePKCE)
}

func (s *FositeSQLStore) DeletePKCERequestSession(ctx context.Context, signature string) error {
	return s.deleteSession(ctx, signature, sqlTablePKCE)
}

func (s *FositeSQLStore) CreateImplicitAccessTokenSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.CreateAccessTokenSession(ctx, signature, requester)
}

func (s *FositeSQLStore) RevokeRefreshToken(ctx context.Context, id string) error {
	return s.revokeSession(ctx, id, sqlTableRefresh)
}

func (s *FositeSQLStore) RevokeAccessToken(ctx context.Context, id string) error {
	return s.revokeSession(ctx, id, sqlTableAccess)
}

func (s *FositeSQLStore) revokeSession(ctx context.Context, id string, table string) error {
	if _, err := s.DB.ExecContext(ctx, s.DB.Rebind(fmt.Sprintf("DELETE FROM hydra_oauth2_%s WHERE request_id=?", table)), id); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (s *FositeSQLStore) FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time) error {
	if _, err := s.DB.ExecContext(ctx, s.DB.Rebind(fmt.Sprintf("DELETE FROM hydra_oauth2_%s WHERE requested_at < ? AND requested_at < ?", sqlTableAccess)), time.Now().Add(-s.AccessTokenLifespan), notAfter); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}
