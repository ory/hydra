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
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/sqlcon"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

type SQLManager struct {
	db *sqlx.DB
	c  client.Manager
}

func NewSQLManager(db *sqlx.DB, c client.Manager) *SQLManager {
	return &SQLManager{
		db: db,
		c:  c,
	}
}

func (m *SQLManager) CreateSchemas() (int, error) {
	migrate.SetTable("hydra_oauth2_authentication_consent_migration")
	n, err := migrate.Exec(m.db.DB, m.db.DriverName(), migrations, migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *SQLManager) CreateConsentRequest(c *ConsentRequest) error {
	d, err := newSQLConsentRequest(c)
	if err != nil {
		return err
	}

	if _, err := m.db.NamedExec(fmt.Sprintf(
		"INSERT INTO hydra_oauth2_consent_request (%s) VALUES (%s)",
		strings.Join(sqlParamsRequest, ", "),
		":"+strings.Join(sqlParamsRequest, ", :"),
	), d); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) GetConsentRequest(challenge string) (*ConsentRequest, error) {
	var d sqlRequest

	if err := m.db.Get(&d, m.db.Rebind("SELECT * FROM hydra_oauth2_consent_request WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	c, err := m.c.GetConcreteClient(d.Client)
	if err != nil {
		return nil, err
	}

	return d.toConsentRequest(c)
}

func (m *SQLManager) CreateAuthenticationRequest(c *AuthenticationRequest) error {
	d, err := newSQLAuthenticationRequest(c)
	if err != nil {
		return err
	}

	if _, err := m.db.NamedExec(fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_request (%s) VALUES (%s)",
		strings.Join(sqlParamsRequest, ", "),
		":"+strings.Join(sqlParamsRequest, ", :"),
	), d); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) GetAuthenticationRequest(challenge string) (*AuthenticationRequest, error) {
	var d sqlRequest

	if err := m.db.Get(&d, m.db.Rebind("SELECT * FROM hydra_oauth2_authentication_request WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	c, err := m.c.GetConcreteClient(d.Client)
	if err != nil {
		return nil, err
	}

	return d.toAuthenticationRequest(c)
}

func (m *SQLManager) HandleConsentRequest(challenge string, r *HandledConsentRequest) (*ConsentRequest, error) {
	d, err := newSQLHandledConsentRequest(r)
	if err != nil {
		return nil, err
	}

	if _, err := m.db.NamedExec(fmt.Sprintf(
		"INSERT INTO hydra_oauth2_consent_request_handled (%s) VALUES (%s)",
		strings.Join(sqlParamsConsentRequestHandled, ", "),
		":"+strings.Join(sqlParamsConsentRequestHandled, ", :"),
	), d); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.GetConsentRequest(challenge)
}

func (m *SQLManager) VerifyAndInvalidateConsentRequest(verifier string) (*HandledConsentRequest, error) {
	var d sqlHandledConsentRequest
	var challenge string

	// This can be solved more elegantly with a join statement, but it works for now

	if err := m.db.Get(&challenge, m.db.Rebind("SELECT challenge FROM hydra_oauth2_consent_request WHERE verifier=?"), verifier); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if err := m.db.Get(&d, m.db.Rebind("SELECT * FROM hydra_oauth2_consent_request_handled WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if d.WasUsed {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
	}

	r, err := m.GetConsentRequest(challenge)
	if err != nil {
		return nil, err
	}

	if _, err := m.db.Exec(m.db.Rebind("UPDATE hydra_oauth2_consent_request_handled SET was_used=true WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return d.toHandledConsentRequest(r)
}

func (m *SQLManager) HandleAuthenticationRequest(challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error) {
	d, err := newSQLHandledAuthenticationRequest(r)
	if err != nil {
		return nil, err
	}

	if _, err := m.db.NamedExec(fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_request_handled (%s) VALUES (%s)",
		strings.Join(sqlParamsAuthenticationRequestHandled, ", "),
		":"+strings.Join(sqlParamsAuthenticationRequestHandled, ", :"),
	), d); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.GetAuthenticationRequest(challenge)
}

func (m *SQLManager) VerifyAndInvalidateAuthenticationRequest(verifier string) (*HandledAuthenticationRequest, error) {
	var d sqlHandledAuthenticationRequest
	var challenge string

	// This can be solved more elegantly with a join statement, but it works for now

	if err := m.db.Get(&challenge, m.db.Rebind("SELECT challenge FROM hydra_oauth2_authentication_request WHERE verifier=?"), verifier); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if err := m.db.Get(&d, m.db.Rebind("SELECT * FROM hydra_oauth2_authentication_request_handled WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if d.WasUsed {
		return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
	}

	if _, err := m.db.Exec(m.db.Rebind("UPDATE hydra_oauth2_authentication_request_handled SET was_used=true WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	r, err := m.GetAuthenticationRequest(challenge)
	if err != nil {
		return nil, err
	}

	return d.toHandledAuthenticationRequest(r)
}

func (m *SQLManager) GetAuthenticationSession(id string) (*AuthenticationSession, error) {
	var a AuthenticationSession
	if err := m.db.Get(&a, m.db.Rebind("SELECT * FROM hydra_oauth2_authentication_session WHERE id=?"), id); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &a, nil
}

func (m *SQLManager) CreateAuthenticationSession(a *AuthenticationSession) error {
	if _, err := m.db.NamedExec(fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_session (%s) VALUES (%s)",
		strings.Join(sqlParamsAuthSession, ", "),
		":"+strings.Join(sqlParamsAuthSession, ", :"),
	), a); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) DeleteAuthenticationSession(id string) error {
	if _, err := m.db.Exec(m.db.Rebind("DELETE FROM hydra_oauth2_authentication_session WHERE id=?"), id); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) FindPreviouslyGrantedConsentRequests(client string, subject string) ([]HandledConsentRequest, error) {
	var a []sqlHandledConsentRequest

	if err := m.db.Select(&a, m.db.Rebind(`SELECT h.* FROM
	hydra_oauth2_consent_request_handled as h
JOIN
	hydra_oauth2_consent_request as r ON (h.challenge = r.challenge)
WHERE
		r.subject=? AND r.client_id=? AND r.skip=FALSE
	AND
		(h.error='{}' AND h.remember=TRUE)
`), subject, client); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(errNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	var aa []HandledConsentRequest
	for _, v := range a {
		r, err := m.GetConsentRequest(v.Challenge)
		if err != nil {
			return nil, err
		} else if errors.Cause(err) == sqlcon.ErrNoRows {
			return nil, errors.WithStack(errNoPreviousConsentFound)
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
		return []HandledConsentRequest{}, nil
	}

	return aa, nil
}
