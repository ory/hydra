// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/otelx"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/popx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

var _ consent.Manager = &Persister{}

func (p *Persister) RevokeSubjectConsentSession(ctx context.Context, user string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectConsentSession")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ?", user))
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectClientConsentSession", trace.WithAttributes(attribute.String("client", client)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ? AND client_id = ?", user, client))
}

func (p *Persister) RevokeConsentSessionByID(ctx context.Context, consentRequestID string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeConsentSessionByID",
		trace.WithAttributes(attribute.String("consent_challenge_id", consentRequestID)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, p.revokeConsentSession("consent_challenge_id = ?", consentRequestID))
}

func (p *Persister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
	return func(ctx context.Context, c *pop.Connection) error {
		fs := make([]*flow.Flow, 0)
		if err := p.QueryWithNetwork(ctx).
			Where(whereStmt, whereArgs...).
			Select("consent_challenge_id").
			All(&fs); errors.Is(err, sql.ErrNoRows) {
			return errors.WithStack(x.ErrNotFound)
		} else if err != nil {
			return sqlcon.HandleError(err)
		}

		ids := make([]interface{}, 0, len(fs))
		nid := p.NetworkID(ctx)
		for _, f := range fs {
			ids = append(ids, f.ConsentRequestID.String())
		}

		if len(ids) == 0 {
			return nil
		}

		if err := p.QueryWithNetwork(ctx).
			Where("nid = ?", nid).
			Where("request_id IN (?)", ids...).
			Delete(OAuth2RequestSQL{Table: sqlTableAccess}.TableName()); errors.Is(err, fosite.ErrNotFound) {
			// do nothing
		} else if err != nil {
			return err
		}

		if err := p.QueryWithNetwork(ctx).
			Where("nid = ?", nid).
			Where("request_id IN (?)", ids...).
			Delete(OAuth2RefreshTable{}.TableName()); errors.Is(err, fosite.ErrNotFound) {
			// do nothing
		} else if err != nil {
			return err
		}

		if err := p.QueryWithNetwork(ctx).
			Where("nid = ?", nid).
			Where("consent_challenge_id IN (?)", ids...).
			Delete(new(flow.Flow)); errors.Is(err, sql.ErrNoRows) {
			return errors.WithStack(x.ErrNotFound)
		} else if err != nil {
			return sqlcon.HandleError(err)
		}

		return nil
	}
}

func (p *Persister) RevokeSubjectLoginSession(ctx context.Context, subject string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectLoginSession")
	defer otelx.End(span, &err)

	err = p.QueryWithNetwork(ctx).Where("subject = ?", subject).Delete(&flow.LoginSession{})
	if err != nil {
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

func (p *Persister) CreateForcedObfuscatedLoginSession(ctx context.Context, session *consent.ForcedObfuscatedLoginSession) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateForcedObfuscatedLoginSession")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		nid := p.NetworkID(ctx)
		if err := c.RawQuery(
			"DELETE FROM hydra_oauth2_obfuscated_authentication_session WHERE nid = ? AND client_id = ? AND subject = ?",
			nid,
			session.ClientID,
			session.Subject,
		).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}

		return sqlcon.HandleError(c.RawQuery(
			"INSERT INTO hydra_oauth2_obfuscated_authentication_session (nid, subject, client_id, subject_obfuscated) VALUES (?, ?, ?, ?)",
			nid,
			session.Subject,
			session.ClientID,
			session.SubjectObfuscated,
		).Exec())
	})
}

func (p *Persister) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (_ *consent.ForcedObfuscatedLoginSession, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetForcedObfuscatedLoginSession", trace.WithAttributes(attribute.String("client", client)))
	defer otelx.End(span, &err)

	var s consent.ForcedObfuscatedLoginSession

	if err := p.Connection(ctx).Where(
		"client_id = ? AND subject_obfuscated = ? AND nid = ?",
		client,
		obfuscated,
		p.NetworkID(ctx),
	).First(&s); errors.Is(err, sql.ErrNoRows) {
		return nil, errors.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &s, nil
}

type FlowWithConstantColumns struct {
	*flow.Flow

	State flow.State `db:"state"`

	// we need to write these columns because of the check constraint, but we will soon switch to a new table anyway that will not have them at all
	LoginRemember    bool   `db:"login_remember"`
	LoginRememberFor int    `db:"login_remember_for"`
	LoginError       string `db:"login_error"`
	LoginUsed        bool   `db:"login_was_used"`
	ConsentVerifier  string `db:"consent_verifier"`
	ConsentCSRF      string `db:"consent_csrf"`
	ConsentError     string `db:"consent_error"`
	ConsentUsed      bool   `db:"consent_was_used"`

	// these columns have NOT NULL constraints, but are not required to be stored
	LoginVerifier string `db:"login_verifier"`
	LoginCSRF     string `db:"login_csrf"`
	LoginSkip     bool   `db:"login_skip"`
}

func (p *Persister) CreateConsentSession(ctx context.Context, f *flow.Flow) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateConsentSession")
	defer otelx.End(span, &err)

	if f.NID != p.NetworkID(ctx) {
		return errors.WithStack(sqlcon.ErrNoRows)
	}
	if f.ConsentRememberFor == nil {
		// This is really stupid: we treat 0 the same as NULL, which means the flow does not expire.
		// However, for some reason it is part of the check constraint and required to be NOT NULL.
		f.ConsentRememberFor = pointerx.Ptr(0)
	}

	fx := &FlowWithConstantColumns{
		Flow:         f,
		State:        flow.FlowStateConsentUsed, // if this was another state, we'd not store it in the DB
		ConsentUsed:  true,
		LoginUsed:    true,
		LoginError:   "{}",
		ConsentError: "{}",
	}
	return sqlcon.HandleError(p.Connection(ctx).Create(fx))
}

func (p *Persister) GetRememberedLoginSession(ctx context.Context, id string) (_ *flow.LoginSession, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetRememberedLoginSession")
	defer otelx.End(span, &err)

	var s flow.LoginSession
	if err := p.QueryWithNetwork(ctx).Where("remember = TRUE").Find(&s, id); errors.Is(err, sql.ErrNoRows) {
		return nil, errors.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &s, nil
}

// ConfirmLoginSession creates or updates the login session. The NID will be set to the network ID of the context.
func (p *Persister) ConfirmLoginSession(ctx context.Context, loginSession *flow.LoginSession) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ConfirmLoginSession")
	defer otelx.End(span, &err)

	loginSession.NID = p.NetworkID(ctx)
	loginSession.AuthenticatedAt = sqlxx.NullTime(time.Time(loginSession.AuthenticatedAt).Truncate(time.Second))
	loginSession.ExpiresAt = sqlxx.NullTime(time.Now().Truncate(time.Second).Add(p.r.Config().GetAuthenticationSessionLifespan(ctx)).UTC())

	if p.Connection(ctx).Dialect.Name() == "mysql" {
		// MySQL does not support UPSERT.
		return p.mySQLConfirmLoginSession(ctx, loginSession)
	}

	res, err := p.Connection(ctx).Store.NamedExecContext(ctx, `
INSERT INTO hydra_oauth2_authentication_session (id, nid, authenticated_at, subject, remember, identity_provider_session_id, expires_at)
VALUES (:id, :nid, :authenticated_at, :subject, :remember, :identity_provider_session_id, :expires_at)
ON CONFLICT(id) DO
UPDATE SET
	authenticated_at = :authenticated_at,
	subject = :subject,
	remember = :remember,
	identity_provider_session_id = :identity_provider_session_id,
	expires_at = :expires_at
WHERE hydra_oauth2_authentication_session.id = :id AND hydra_oauth2_authentication_session.nid = :nid
`, loginSession)
	if err != nil {
		return sqlcon.HandleError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return sqlcon.HandleError(err)
	}
	if n == 0 {
		return errors.WithStack(x.ErrNotFound)
	}
	return nil
}

func (p *Persister) DeleteLoginSession(ctx context.Context, id string) (_ *flow.LoginSession, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteLoginSession")
	defer otelx.End(span, &err)

	c := p.Connection(ctx)
	if c.Dialect.Name() == "mysql" {
		// MySQL does not support RETURNING.
		return p.mySQLDeleteLoginSession(ctx, id)
	}

	var session flow.LoginSession
	columns := popx.DBColumns[flow.LoginSession](c.Dialect)
	if err := p.Connection(ctx).RawQuery(
		fmt.Sprintf(`DELETE FROM hydra_oauth2_authentication_session WHERE id = ? AND nid = ? RETURNING %s`, columns),
		id,
		p.NetworkID(ctx),
	).First(&session); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &session, nil
}

func (p *Persister) mySQLDeleteLoginSession(ctx context.Context, id string) (_ *flow.LoginSession, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.mySQLDeleteLoginSession")
	defer otelx.End(span, &err)

	var session flow.LoginSession
	if err := p.Connection(ctx).Transaction(func(tx *pop.Connection) error {
		if err := tx.Where("id = ? AND nid = ?", id, p.NetworkID(ctx)).First(&session); err != nil {
			return err
		}

		return tx.RawQuery(
			`DELETE FROM hydra_oauth2_authentication_session WHERE id = ? AND nid = ?`,
			id, p.NetworkID(ctx),
		).Exec()
	}); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &session, nil
}

func (p *Persister) FindGrantedAndRememberedConsentRequest(ctx context.Context, client, subject string) (_ *flow.Flow, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindGrantedAndRememberedConsentRequest")
	defer otelx.End(span, &err)

	f := flow.Flow{}
	conn := p.Connection(ctx)

	// apply index hint
	tableName := applyTableNameWithIndexHint(conn, f.TableName(), "hydra_oauth2_flow_previous_consents_idx")

	// prepare columns
	cols := popx.DBColumns[flow.Flow](conn.Dialect)

	// prepare sql statement
	q := fmt.Sprintf(`
SELECT %s FROM %s
WHERE nid = ?
AND state = ?
AND subject = ?
AND client_id = ?
AND consent_skip = FALSE
AND consent_error = '{}'
AND consent_remember = TRUE
AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
ORDER BY requested_at DESC
LIMIT 1`,
		cols,
		tableName,
	)

	// query first record
	err = conn.RawQuery(q,
		p.NetworkID(ctx),
		flow.FlowStateConsentUsed,
		subject,
		client,
	).First(&f)

	// handle error
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.WithStack(consent.ErrNoPreviousConsentFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &f, nil
}

func applyTableNameWithIndexHint(conn *pop.Connection, table string, index string) string {
	switch conn.Dialect.Name() {
	case "cockroach":
		return table + "@" + index
	case "sqlite3":
		return table + " INDEXED BY " + index
	case "mysql":
		return table + " USE INDEX(" + index + ")"
	default:
		return table
	}
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, pageOpts ...keysetpagination.Option) (_ []flow.Flow, _ *keysetpagination.Paginator, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindSubjectsGrantedConsentRequests")
	defer otelx.End(span, &err)

	paginator := keysetpagination.NewPaginator(append(pageOpts,
		keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(keysetpagination.Column{Name: "login_challenge", Value: ""})),
	)...)

	var fs []flow.Flow
	err = p.QueryWithNetwork(ctx).
		Where("state IN (?, ?)", flow.FlowStateConsentUsed, flow.FlowStateConsentUnused).
		Where("subject = ?", subject).
		Where("consent_skip = FALSE").
		Where("consent_error = '{}'").
		Where("(expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)").
		Scope(keysetpagination.Paginate[flow.Flow](paginator)).
		All(&fs)
	if err != nil {
		return nil, nil, sqlcon.HandleError(err)
	}
	if len(fs) == 0 {
		return nil, nil, errors.WithStack(consent.ErrNoPreviousConsentFound)
	}

	fs, nextPage := keysetpagination.Result(fs, paginator)
	return fs, nextPage, nil
}

func (p *Persister) FindSubjectsSessionGrantedConsentRequests(ctx context.Context, subject, sid string, pageOpts ...keysetpagination.Option) (_ []flow.Flow, _ *keysetpagination.Paginator, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindSubjectsSessionGrantedConsentRequests", trace.WithAttributes(attribute.String("sid", sid)))
	defer otelx.End(span, &err)

	paginator := keysetpagination.NewPaginator(append(pageOpts,
		keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(keysetpagination.Column{Name: "login_challenge", Value: ""})),
	)...)

	var fs []flow.Flow
	err = p.QueryWithNetwork(ctx).
		Where("state IN (?, ?)", flow.FlowStateConsentUsed, flow.FlowStateConsentUnused).
		Where("subject = ?", subject).
		Where("login_session_id = ?", sid).
		Where("consent_skip = FALSE").
		Where("consent_error = '{}'").
		Where("(expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)").
		Scope(keysetpagination.Paginate[flow.Flow](paginator)).
		All(&fs)
	if err != nil {
		return nil, nil, sqlcon.HandleError(err)
	}
	if len(fs) == 0 {
		return nil, nil, errors.WithStack(consent.ErrNoPreviousConsentFound)
	}

	fs, nextPage := keysetpagination.Result(fs, paginator)
	return fs, nextPage, nil
}

func (p *Persister) ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) (_ []client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ListUserAuthenticatedClientsWithFrontChannelLogout")
	defer otelx.End(span, &err)

	return p.listUserAuthenticatedClients(ctx, subject, sid, "front")
}

func (p *Persister) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) (_ []client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ListUserAuthenticatedClientsWithBackChannelLogout")
	defer otelx.End(span, &err)

	return p.listUserAuthenticatedClients(ctx, subject, sid, "back")
}

func (p *Persister) listUserAuthenticatedClients(ctx context.Context, subject, sid, channel string) (cs []client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.listUserAuthenticatedClients",
		trace.WithAttributes(attribute.String("sid", sid)))
	defer otelx.End(span, &err)

	conn := p.Connection(ctx)
	columns := popx.DBColumns[client.Client](&popx.AliasQuoter{Alias: "c", Quoter: conn.Dialect})

	if err := conn.RawQuery(
		/* #nosec G201 - channel can either be "front" or "back" */
		fmt.Sprintf(`
SELECT DISTINCT %s FROM hydra_client as c
JOIN hydra_oauth2_flow as f ON (c.id = f.client_id AND c.nid = f.nid)
WHERE
	f.subject = ? AND
	c.%schannel_logout_uri != '' AND
	c.%schannel_logout_uri IS NOT NULL AND
	f.login_session_id = ? AND
	f.nid = ? AND
	c.nid = ?`,
			columns,
			channel,
			channel,
		),
		subject,
		sid,
		p.NetworkID(ctx),
		p.NetworkID(ctx),
	).All(&cs); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return cs, nil
}

func (p *Persister) CreateLogoutRequest(ctx context.Context, request *flow.LogoutRequest) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateLogoutRequest")
	defer otelx.End(span, &err)

	return errors.WithStack(p.CreateWithNetwork(ctx, request))
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (_ *flow.LogoutRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AcceptLogoutRequest")
	defer otelx.End(span, &err)

	if err := p.Connection(ctx).RawQuery("UPDATE hydra_oauth2_logout_request SET accepted=true, rejected=false WHERE challenge=? AND nid = ?", challenge, p.NetworkID(ctx)).Exec(); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetLogoutRequest(ctx, challenge)
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RejectLogoutRequest")
	defer otelx.End(span, &err)

	count, err := p.Connection(ctx).
		RawQuery("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=? AND nid = ?", challenge, p.NetworkID(ctx)).
		ExecWithCount()
	if count == 0 {
		return errors.WithStack(x.ErrNotFound)
	} else {
		return errors.WithStack(err)
	}
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (_ *flow.LogoutRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetLogoutRequest")
	defer otelx.End(span, &err)

	var lr flow.LogoutRequest
	return &lr, sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("challenge = ? AND rejected = FALSE", challenge).First(&lr))
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (_ *flow.LogoutRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.VerifyAndInvalidateLogoutRequest")
	defer otelx.End(span, &err)

	var lr flow.LogoutRequest
	if count, err := p.Connection(ctx).RawQuery(`
UPDATE hydra_oauth2_logout_request
  SET was_used = TRUE
WHERE nid = ?
  AND verifier = ?
  AND accepted = TRUE
  AND rejected = FALSE`,
		p.NetworkID(ctx),
		verifier,
	).ExecWithCount(); count == 0 && err == nil {
		return nil, errors.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	err = sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("verifier = ?", verifier).First(&lr))
	if err != nil {
		return nil, err
	}

	if expiry := time.Time(lr.ExpiresAt);
	// If the expiry is unset, we are in a legacy use case (allow logout).
	// TODO: Remove this in the future.
	!expiry.IsZero() && expiry.Before(time.Now().UTC()) {
		return nil, errors.WithStack(flow.ErrorLogoutFlowExpired)
	}

	return &lr, nil
}

func (p *Persister) FlushInactiveLoginConsentRequests(ctx context.Context, notAfter time.Time, limit, batchSize int) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FlushInactiveLoginConsentRequests")
	defer otelx.End(span, &err)

	// The value of notAfter should be the minimum between input parameter and request max expire based on its configured age
	requestMaxExpire := time.Now().Add(-p.r.Config().ConsentRequestMaxAge(ctx))
	if requestMaxExpire.Before(notAfter) {
		notAfter = requestMaxExpire
	}

	challenges := make([]string, 0, limit)
	// Select up to [limit] flows that can be safely deleted, i.e. flows that meet
	// the following criteria:
	// - flow.state is anything between FlowStateLoginInitialized and FlowStateConsentUnused (unhandled)
	// - flow.login_error has valid error (login rejected)
	// - flow.consent_error has valid error (consent rejected)
	// AND timed-out
	// - flow.requested_at < minimum of ttl.login_consent_request and notAfter
	q := p.Connection(ctx).RawQuery(`
		SELECT login_challenge
		FROM hydra_oauth2_flow
		WHERE (
			(state != ?)
			OR (login_error IS NOT NULL AND login_error <> '{}' AND login_error <> '')
			OR (consent_error IS NOT NULL AND consent_error <> '{}' AND consent_error <> '')
		)
		AND requested_at < ?
		AND nid = ?
		ORDER BY login_challenge
		LIMIT ?`,
		flow.FlowStateConsentUsed, notAfter, p.NetworkID(ctx), limit)

	if err := q.All(&challenges); err != nil {
		return errors.WithStack(err)
	}

	// Delete in batch consent requests and their references in cascade
	for i := 0; i < len(challenges); i += batchSize {
		j := min(i+batchSize, len(challenges))

		q := p.Connection(ctx).RawQuery(
			"DELETE FROM hydra_oauth2_flow WHERE login_challenge in (?) AND nid = ?",
			challenges[i:j],
			p.NetworkID(ctx),
		)

		if err := q.Exec(); err != nil {
			return sqlcon.HandleError(err)
		}
	}

	return nil
}

func (p *Persister) mySQLConfirmLoginSession(ctx context.Context, session *flow.LoginSession) error {
	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := sqlcon.HandleError(c.Create(session))
		if err == nil {
			return nil
		}

		if !errors.Is(err, sqlcon.ErrUniqueViolation) {
			return err
		}

		n, err := c.
			Where("id = ? and nid = ?", session.ID, session.NID).
			UpdateQuery(session, "authenticated_at", "subject", "identity_provider_session_id", "remember", "expires_at")
		if err != nil {
			return errors.WithStack(sqlcon.HandleError(err))
		}
		if n == 0 {
			return errors.WithStack(x.ErrNotFound)
		}

		return nil
	})
}
