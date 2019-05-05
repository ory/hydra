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

	"github.com/ory/hydra/client"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/ory/fosite"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
)

type SQLManager struct {
	DB *sqlx.DB
	r  InternalRegistry
}

func NewSQLManager(db *sqlx.DB, r InternalRegistry) *SQLManager {
	return &SQLManager{
		DB: db,
		r:  r,
	}
}

func (m *SQLManager) PlanMigration() ([]*migrate.PlannedMigration, error) {
	migrate.SetTable("hydra_oauth2_authentication_consent_migration")
	plan, _, err := migrate.PlanMigration(m.DB.DB, m.DB.DriverName(), Migrations[dbal.Canonicalize(m.DB.DriverName())], migrate.Up, 0)
	return plan, errors.WithStack(err)
}

func (m *SQLManager) CreateSchemas() (int, error) {
	migrate.SetTable("hydra_oauth2_authentication_consent_migration")
	n, err := migrate.Exec(m.DB.DB, m.DB.DriverName(), Migrations[dbal.Canonicalize(m.DB.DriverName())], migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *SQLManager) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	return m.revokeConsentSession(ctx, user, "")
}

func (m *SQLManager) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
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
			return errors.WithStack(x.ErrNotFound)
		}
		return sqlcon.HandleError(err)
	}

	for _, challenge := range challenges {
		if err := m.r.OAuth2Storage().RevokeAccessToken(ctx, challenge); errors.Cause(err) == fosite.ErrNotFound {
			// do nothing
		} else if err != nil {
			return err
		} else {
			metricAccessTokensRevoked.Inc()
		}
		if err := m.r.OAuth2Storage().RevokeRefreshToken(ctx, challenge); errors.Cause(err) == fosite.ErrNotFound {
			// do nothing
		} else if err != nil {
			return err
		} else {
			metricRefreshTokensRevoked.Inc()
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
				return errors.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}

		if count, _ := rows.RowsAffected(); count == 0 {
			return errors.WithStack(x.ErrNotFound)
		}
	}
	return nil
}

func (m *SQLManager) RevokeSubjectLoginSession(ctx context.Context, user string) error {
	_, err := m.DB.ExecContext(
		ctx,
		m.DB.Rebind("DELETE FROM hydra_oauth2_authentication_session WHERE subject=?"),
		user,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.WithStack(x.ErrNotFound)
		}
		return sqlcon.HandleError(err)
	}

	// This confuses people, see https://github.com/ory/hydra/issues/1168
	//
	// count, _ := rows.RowsAffected()
	// if count == 0 {
	// 	 return errors.WithStack(x.ErrNotFound)
	// }

	return nil
}

func (m *SQLManager) CreateForcedObfuscatedLoginSession(ctx context.Context, s *ForcedObfuscatedLoginSession) error {
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

func (m *SQLManager) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedLoginSession, error) {
	var d ForcedObfuscatedLoginSession

	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_oauth2_obfuscated_authentication_session WHERE client_id=? AND subject_obfuscated=?"), client, obfuscated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(x.ErrNotFound)
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
			return nil, errors.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	c, err := m.r.ClientManager().GetConcreteClient(ctx, d.Client)
	if err != nil {
		return nil, err
	}

	return d.toConsentRequest(c)
}

func (m *SQLManager) CreateLoginRequest(ctx context.Context, c *LoginRequest) error {
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

func (m *SQLManager) GetLoginRequest(ctx context.Context, challenge string) (*LoginRequest, error) {
	var d sqlAuthenticationRequest
	err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT r.*, COALESCE(hr.was_used, false) as was_handled FROM hydra_oauth2_authentication_request r "+
		"LEFT JOIN hydra_oauth2_authentication_request_handled hr ON r.challenge = hr.challenge WHERE r.challenge=?"), challenge)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	c, err := m.r.ClientManager().GetConcreteClient(ctx, d.Client)
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
		err = sqlcon.HandleError(err)
		if errors.Cause(err) == sqlcon.ErrUniqueViolation {
			return m.replaceUnusedConsentRequest(ctx, challenge, d)
		}
		return nil, err
	}

	return m.GetConsentRequest(ctx, challenge)
}

func (m *SQLManager) replaceUnusedConsentRequest(ctx context.Context, challenge string, d *sqlHandledConsentRequest) (*ConsentRequest, error) {
	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"UPDATE hydra_oauth2_consent_request_handled SET %s WHERE challenge=:challenge AND was_used=false",
		strings.Join(sqlParamsConsentRequestHandledUpdate, ", "),
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

func (m *SQLManager) HandleLoginRequest(ctx context.Context, challenge string, r *HandledLoginRequest) (*LoginRequest, error) {
	d, err := newSQLHandledLoginRequest(r)
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

	return m.GetLoginRequest(ctx, challenge)
}

func (m *SQLManager) VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*HandledLoginRequest, error) {
	var d sqlHandledLoginRequest
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

	r, err := m.GetLoginRequest(ctx, challenge)
	if err != nil {
		return nil, err
	}

	return d.toHandledLoginRequest(r)
}

func (m *SQLManager) GetRememberedLoginSession(ctx context.Context, id string) (*LoginSession, error) {
	var a LoginSession
	if err := m.DB.GetContext(ctx, &a, m.DB.Rebind("SELECT * FROM hydra_oauth2_authentication_session WHERE id=? AND remember=TRUE"), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return &a, nil
}

func (m *SQLManager) ConfirmLoginSession(ctx context.Context, id string, subject string, remember bool) error {
	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("UPDATE hydra_oauth2_authentication_session SET remember=?, subject=?, authenticated_at=? WHERE id=?"), remember, subject, time.Now().UTC(), id); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) CreateLoginSession(ctx context.Context, a *LoginSession) error {
	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_authentication_session (%s) VALUES (%s)",
		strings.Join(sqlParamsAuthSession, ", "),
		":"+strings.Join(sqlParamsAuthSession, ", :"),
	), a); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) DeleteLoginSession(ctx context.Context, id string) error {
	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("DELETE FROM hydra_oauth2_authentication_session WHERE id=?"), id); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, subject string) ([]HandledConsentRequest, error) {
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

func (m *SQLManager) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]HandledConsentRequest, error) {
	var a []sqlHandledConsentRequest

	if err := m.DB.SelectContext(ctx, &a, m.DB.Rebind(`SELECT h.* FROM
	hydra_oauth2_consent_request_handled as h
JOIN
	hydra_oauth2_consent_request as r ON (h.challenge = r.challenge)
WHERE
		r.subject=? AND r.skip=FALSE
	AND
		(h.error='{}')
ORDER BY h.requested_at DESC
LIMIT ? OFFSET ?
`), subject, limit, offset); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.resolveHandledConsentRequests(ctx, a)
}

func (m *SQLManager) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	var n int

	if err := m.DB.QueryRowContext(ctx, m.DB.Rebind(`SELECT COUNT(*) FROM
	hydra_oauth2_consent_request_handled as h
JOIN
	hydra_oauth2_consent_request as r ON (h.challenge = r.challenge)
WHERE
		r.subject=? AND r.skip=FALSE
	AND
		(h.error='{}')
`), subject).Scan(&n); err != nil {
		return 0, sqlcon.HandleError(err)
	}

	return n, nil
}

func (m *SQLManager) resolveHandledConsentRequests(ctx context.Context, requests []sqlHandledConsentRequest) ([]HandledConsentRequest, error) {
	var aa []HandledConsentRequest
	for _, v := range requests {
		r, err := m.GetConsentRequest(ctx, v.Challenge)
		if err != nil {
			return nil, err
		} else if errors.Cause(err) == x.ErrNotFound {
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

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent. The
// descriptor of each sent metric is one of those returned by Describe
// (unless the Collector is unchecked, see above). Returned metrics that
// share the same descriptor must differ in their variable label
// values.
//
// This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. Blocking occurs at the expense
// of total performance of rendering all registered metrics. Ideally,
// Collector implementations support concurrent readers.
func (m *SQLManager) Collect(c chan<- prometheus.Metric) {
	collectCounters(c)

	// get t
	go func(m *SQLManager, c chan<- prometheus.Metric) {
		var n int
		m.DB.QueryRow("SELECT count(*) FROM hydra_oauth2_access").Scan(&n)
		metricAccessTokens.Inc()
		metricAccessTokens.Collect(c)
	}(m, c)

	go func(m *SQLManager, c chan<- prometheus.Metric) {
		var n int
		m.DB.QueryRow("SELECT count(*) FROM hydra_oauth2_refresh").Scan(&n)
		metricRefreshTokens.Inc()
		metricRefreshTokens.Collect(c)
	}(m, c)
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent. The sent descriptors fulfill the
// consistency and uniqueness requirements described in the Desc
// documentation.
//
// It is valid if one and the same Collector sends duplicate
// descriptors. Those duplicates are simply ignored. However, two
// different Collectors must not send duplicate descriptors.
//
// Sending no descriptor at all marks the Collector as “unchecked”,
// i.e. no checks will be performed at registration time, and the
// Collector may yield any Metric it sees fit in its Collect method.
//
// This method idempotently sends the same descriptors throughout the
// lifetime of the Collector. It may be called concurrently and
// therefore must be implemented in a concurrency safe way.
//
// If a Collector encounters an error while executing this method, it
// must send an invalid descriptor (created with NewInvalidDesc) to
// signal the error to the registry.
func (m *SQLManager) Describe(c chan<- *prometheus.Desc) {
	describeCounters(c)

	metricAccessTokens.Describe(c)
	metricRefreshTokens.Describe(c)
}

func (m *SQLManager) ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject string) ([]client.Client, error) {
	return m.listUserAuthenticatedClients(ctx, subject, "front")
}

func (m *SQLManager) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject string) ([]client.Client, error) {
	return m.listUserAuthenticatedClients(ctx, subject, "back")
}

func (m *SQLManager) listUserAuthenticatedClients(ctx context.Context, subject string, channel string) ([]client.Client, error) {
	var ids []string
	if err := m.DB.SelectContext(ctx, &ids, m.DB.Rebind(fmt.Sprintf(`SELECT c.id FROM hydra_client as c JOIN hydra_oauth2_consent_request as r ON (c.id = r.client_id) WHERE r.subject=? AND c.%schannel_logout_uri!='' and c.%schannel_logout_uri IS NOT NULL`, channel, channel)), subject); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	cs := make([]client.Client, len(ids))
	for k, id := range ids {
		c, err := m.r.ClientManager().GetConcreteClient(ctx, id)
		if err != nil {
			return nil, err
		}
		cs[k] = *c
	}

	return cs, nil
}

func (m *SQLManager) CreateLogoutRequest(ctx context.Context, r *LogoutRequest) error {
	d := newSQLLogoutRequest(r)
	if _, err := m.DB.NamedExecContext(ctx, fmt.Sprintf(
		"INSERT INTO hydra_oauth2_logout_request (%s) VALUES (%s)",
		strings.Join(sqlParamsLogoutRequest, ", "),
		":"+strings.Join(sqlParamsLogoutRequest, ", :"),
	), d); err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

func (m *SQLManager) AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("UPDATE hydra_oauth2_logout_request SET accepted=true, rejected=false WHERE challenge=?"), challenge); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.GetLogoutRequest(ctx, challenge)
}

func (m *SQLManager) RejectLogoutRequest(ctx context.Context, challenge string) error {
	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=?"), challenge); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) GetLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	var d sqlLogoutRequest
	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_oauth2_logout_request WHERE challenge=? AND rejected=FALSE"), challenge); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	if d.Client.Valid {
		c, err := m.r.ClientManager().GetConcreteClient(ctx, d.Client.String)
		if err != nil {
			return nil, err
		}

		return d.ToLogoutRequest(c), nil
	}

	return d.ToLogoutRequest(nil), nil
}

func (m *SQLManager) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*LogoutRequest, error) {
	var d sqlLogoutRequest
	if err := m.DB.GetContext(ctx, &d, m.DB.Rebind("SELECT * FROM hydra_oauth2_logout_request WHERE verifier=? AND was_used=FALSE AND accepted=TRUE AND rejected=FALSE"), verifier); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	if _, err := m.DB.ExecContext(ctx, m.DB.Rebind("UPDATE hydra_oauth2_logout_request SET was_used=TRUE WHERE verifier=?"), verifier); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return m.GetLogoutRequest(ctx, d.Challenge)
}
