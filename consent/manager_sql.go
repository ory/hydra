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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
)

type SQLManager struct {
	DB    *sqlx.DB
	c     client.Manager
	store pkg.FositeStorer
}

func NewSQLManager(db *sqlx.DB, c client.Manager, store pkg.FositeStorer) *SQLManager {
	return &SQLManager{
		DB:    db,
		c:     c,
		store: store,
	}
}

func (m *SQLManager) CreateSchemas() (int, error) {
	migrate.SetTable("hydra_oauth2_authentication_consent_migration")
	n, err := migrate.Exec(m.DB.DB, m.DB.DriverName(), Migrations[dbal.Canonicalize(m.DB.DriverName())], migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *SQLManager) RevokeUserConsentSession(ctx context.Context, user string) error {
	return m.revokeConsentSession(ctx, user, "")
}

func (m *SQLManager) RevokeUserClientConsentSession(ctx context.Context, user, client string) error {
	return m.revokeConsentSession(ctx, user, client)
}

func (m *SQLManager) revokeConsentSession(ctx context.Context, user, client string) error {
	args := []interface{}{user}
	part := "r.subject=?"
	if client != "" {
		part += " AND r.client_id=?"
		args = append(args, client)
	}

	var challenges = make([]string, 0)
	if err := m.DB.SelectContext(ctx, &challenges, m.DB.Rebind(fmt.Sprintf(
		`SELECT r.challenge FROM hydra_oauth2_consent_request_handled as h 
JOIN hydra_oauth2_consent_request as r ON r.challenge = h.challenge WHERE %s`,
		part,
	)), args...); err != nil {
		if err == sql.ErrNoRows {
			return errors.WithStack(pkg.ErrNotFound)
		}
		return sqlcon.HandleError(err)
	}

	for _, challenge := range challenges {
		if err := m.store.RevokeAccessToken(ctx, challenge); errors.Cause(err) == fosite.ErrNotFound {
			// do nothing
		} else if err != nil {
			return err
		}
		if err := m.store.RevokeRefreshToken(ctx, challenge); errors.Cause(err) == fosite.ErrNotFound {
			// do nothing
		} else if err != nil {
			return err
		}
	}

	var queries []string
	switch m.DB.DriverName() {
	case "mysql":
		queries = append(queries,
			fmt.Sprintf(`DELETE h, r FROM hydra_oauth2_consent_request_handled as h 
JOIN hydra_oauth2_consent_request as r ON r.challenge = h.challenge
WHERE %s`, part),
		)
	default:
		queries = append(queries,
			fmt.Sprintf(`DELETE FROM hydra_oauth2_consent_request_handled 
WHERE challenge IN (SELECT r.challenge FROM hydra_oauth2_consent_request as r WHERE %s)`, part),
			fmt.Sprintf(`DELETE FROM hydra_oauth2_consent_request as r WHERE %s`, part),
		)
	}

	for _, q := range queries {
		rows, err := m.DB.ExecContext(ctx, m.DB.Rebind(q), args...)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.WithStack(pkg.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}

		if count, _ := rows.RowsAffected(); count == 0 {
			return errors.WithStack(pkg.ErrNotFound)
		}
	}
	return nil
}

func (m *SQLManager) RevokeUserAuthenticationSession(ctx context.Context, user string) error {
	_, err := m.DB.ExecContext(
		ctx,
		m.DB.Rebind("DELETE FROM hydra_oauth2_authentication_session WHERE subject=?"),
		user,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.WithStack(pkg.ErrNotFound)
		}
		return sqlcon.HandleError(err)
	}

	// This confuses people, see https://github.com/ory/hydra/issues/1168
	//
	// count, _ := rows.RowsAffected()
	// if count == 0 {
	// 	 return errors.WithStack(pkg.ErrNotFound)
	// }

	return nil
}

func (m *SQLManager) CreateForcedObfuscatedAuthenticationSession(ctx context.Context, s *ForcedObfuscatedAuthenticationSession) error {
	tx, err := m.DB.BeginTxx(ctx, nil)
	if err != nil {
		return sqlcon.HandleError(err)
	}

	if _, err := tx.ExecContext(
		ctx,
		m.DB.Rebind("DELETE FROM hydra_oauth2_obfuscated_authentication_session WHERE client_id=? AND subject=?"),
		s.ClientID,
		s.Subject,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return sqlcon.HandleError(err)
		}
		return sqlcon.HandleError(err)
	}

	if _, err := tx.NamedExec(
		"INSERT INTO hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated) VALUES (:subject, :client_id, :subject_obfuscated)",
		s,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return sqlcon.HandleError(err)
		}
		return sqlcon.HandleError(err)
	}

	if err := tx.Commit(); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) GetForcedObfuscatedAuthenticationSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedAuthenticationSession, error) {
	var d ForcedObfuscatedAuthenticationSession

	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_oauth2_obfuscated_authentication_session WHERE client_id=? AND subject_obfuscated=?"), client, obfuscated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(pkg.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return &d, nil
}

func (m *SQLManager) CreateConsentRequest(ctx context.Context, c *ConsentRequest) error {
	d, err := newSQLConsentRequest(c)
	if err != nil {
		return err
	}

	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_consent_request (%s) VALUES (%s)",
		strings.Join(sqlParamsConsentRequest, ", "),
		":"+strings.Join(sqlParamsConsentRequest, ", :"),
	), d); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) GetConsentRequest(ctx context.Context, challenge string) (*ConsentRequest, error) {
	var d sqlConsentRequest
	err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT r.*, COALESCE(hr.was_used, false) as was_handled FROM hydra_oauth2_consent_request r "+
		"LEFT JOIN hydra_oauth2_consent_request_handled hr ON r.challenge = hr.challenge WHERE r.challenge=?"), challenge)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(pkg.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	c, err := m.c.GetConcreteClient(ctx, d.Client)
	if err != nil {
		return nil, err
	}

	return d.toConsentRequest(c)
}

func (m *SQLManager) CreateAuthenticationRequest(ctx context.Context, c *AuthenticationRequest) error {
	d, err := newSQLAuthenticationRequest(c)
	if err != nil {
		return err
	}

	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_request (%s) VALUES (%s)",
		strings.Join(sqlParamsAuthenticationRequest, ", "),
		":"+strings.Join(sqlParamsAuthenticationRequest, ", :"),
	), d); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) GetAuthenticationRequest(ctx context.Context, challenge string) (*AuthenticationRequest, error) {
	var d sqlAuthenticationRequest
	err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT r.*, COALESCE(hr.was_used, false) as was_handled FROM hydra_oauth2_authentication_request r "+
		"LEFT JOIN hydra_oauth2_authentication_request_handled hr ON r.challenge = hr.challenge WHERE r.challenge=?"), challenge)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(pkg.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	c, err := m.c.GetConcreteClient(ctx, d.Client)
	if err != nil {
		return nil, err
	}

	return d.toAuthenticationRequest(c)
}

func (m *SQLManager) HandleConsentRequest(ctx context.Context, challenge string, r *HandledConsentRequest) (*ConsentRequest, error) {
	d, err := newSQLHandledConsentRequest(r)
	if err != nil {
		return nil, err
	}

	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_consent_request_handled (%s) VALUES (%s)",
		strings.Join(sqlParamsConsentRequestHandled, ", "),
		":"+strings.Join(sqlParamsConsentRequestHandled, ", :"),
	), d); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.GetConsentRequest(ctx, challenge)
}

func (m *SQLManager) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*HandledConsentRequest, error) {
	var d sqlHandledConsentRequest
	var challenge string

	// This can be solved more elegantly with a join statement, but it works for now

	if err := m.DB.GetContext(ctx, &challenge, m.DB.Rebind("SELECT challenge FROM hydra_oauth2_consent_request WHERE verifier=?"), verifier); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_oauth2_consent_request_handled WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if d.WasUsed {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
	}

	r, err := m.GetConsentRequest(ctx, challenge)
	if err != nil {
		return nil, err
	}

	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("UPDATE hydra_oauth2_consent_request_handled SET was_used=true WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return d.toHandledConsentRequest(r)
}

func (m *SQLManager) HandleAuthenticationRequest(ctx context.Context, challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error) {
	d, err := newSQLHandledAuthenticationRequest(r)
	if err != nil {
		return nil, err
	}

	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_request_handled (%s) VALUES (%s)",
		strings.Join(sqlParamsAuthenticationRequestHandled, ", "),
		":"+strings.Join(sqlParamsAuthenticationRequestHandled, ", :"),
	), d); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.GetAuthenticationRequest(ctx, challenge)
}

func (m *SQLManager) VerifyAndInvalidateAuthenticationRequest(ctx context.Context, verifier string) (*HandledAuthenticationRequest, error) {
	var d sqlHandledAuthenticationRequest
	var challenge string

	// This can be solved more elegantly with a join statement, but it works for now

	if err := m.DB.GetContext(ctx, &challenge, m.DB.Rebind("SELECT challenge FROM hydra_oauth2_authentication_request WHERE verifier=?"), verifier); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_oauth2_authentication_request_handled WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if d.WasUsed {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Authentication verifier has been used already"))
	}

	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("UPDATE hydra_oauth2_authentication_request_handled SET was_used=true WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	r, err := m.GetAuthenticationRequest(ctx, challenge)
	if err != nil {
		return nil, err
	}

	return d.toHandledAuthenticationRequest(r)
}

func (m *SQLManager) GetAuthenticationSession(ctx context.Context, id string) (*AuthenticationSession, error) {
	var a AuthenticationSession
	if err := m.DB.GetContext(ctx, &a, m.DB.Rebind("SELECT * FROM hydra_oauth2_authentication_session WHERE id=?"), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(pkg.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return &a, nil
}

func (m *SQLManager) CreateAuthenticationSession(ctx context.Context, a *AuthenticationSession) error {
	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_session (%s) VALUES (%s)",
		strings.Join(sqlParamsAuthSession, ", "),
		":"+strings.Join(sqlParamsAuthSession, ", :"),
	), a); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) DeleteAuthenticationSession(ctx context.Context, id string) error {
	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("DELETE FROM hydra_oauth2_authentication_session WHERE id=?"), id); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) FindPreviouslyGrantedConsentRequests(ctx context.Context, client, subject string) ([]HandledConsentRequest, error) {
	var a []sqlHandledConsentRequest

	if err := m.DB.SelectContext(ctx, &a, m.DB.Rebind(`SELECT h.* FROM
	hydra_oauth2_consent_request_handled as h
JOIN
	hydra_oauth2_consent_request as r ON (h.challenge = r.challenge)
WHERE
		r.subject=? AND r.client_id=? AND r.skip=FALSE
	AND
		(h.error='{}' AND h.remember=TRUE)
ORDER BY h.requested_at DESC
LIMIT 1`), subject, client); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(ErrNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return m.resolveHandledConsentRequests(ctx, a)
}

func (m *SQLManager) FindPreviouslyGrantedConsentRequestsByUser(ctx context.Context, subject string, limit, offset int) ([]HandledConsentRequest, error) {
	var a []sqlHandledConsentRequest

	if err := m.DB.SelectContext(ctx, &a, m.DB.Rebind(`SELECT h.* FROM
	hydra_oauth2_consent_request_handled as h
JOIN
	hydra_oauth2_consent_request as r ON (h.challenge = r.challenge)
WHERE
		r.subject=? AND r.skip=FALSE
	AND
		(h.error='{}' AND h.remember=TRUE)
ORDER BY h.requested_at DESC
LIMIT ? OFFSET ?
`), subject, limit, offset); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.resolveHandledConsentRequests(ctx, a)
}

func (m *SQLManager) resolveHandledConsentRequests(ctx context.Context, requests []sqlHandledConsentRequest) ([]HandledConsentRequest, error) {
	var aa []HandledConsentRequest
	for _, v := range requests {
		r, err := m.GetConsentRequest(ctx, v.Challenge)
		if err != nil {
			return nil, err
		} else if errors.Cause(err) == sqlcon.ErrNoRows {
			return nil, errors.WithStack(ErrNoPreviousConsentFound)
		}

		if v.RememberFor > 0 && v.RequestedAt.Add(time.Duration(v.RememberFor)*time.Second).Before(time.Now().UTC()) {
			continue
		}

		va, err := v.toHandledConsentRequest(r)
		if err != nil {
			return nil, err
		}

		aa = append(aa, *va)
	}

	if len(aa) == 0 {
		return nil, errors.WithStack(ErrNoPreviousConsentFound)
	}

	return aa, nil
}
