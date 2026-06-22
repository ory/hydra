// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
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
	"github.com/ory/x/popx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

var (
	_ consent.Manager                  = (*ConsentPersister)(nil)
	_ consent.LoginManager             = (*Persister)(nil)
	_ consent.LogoutManager            = (*Persister)(nil)
	_ consent.ObfuscatedSubjectManager = (*Persister)(nil)
)

type ConsentPersister struct {
	*BasePersister
}

func (p *ConsentPersister) RevokeSubjectConsentSession(ctx context.Context, user string) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectConsentSession")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ?", user))
}

func (p *ConsentPersister) RevokeSubjectClientConsentSession(ctx context.Context, user, clientID string) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectClientConsentSession", trace.WithAttributes(attribute.String("client.id", clientID)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ? AND client_id = ?", user, clientID))
}

func (p *ConsentPersister) RevokeConsentSessionByID(ctx context.Context, consentRequestID string) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeConsentSessionByID",
		trace.WithAttributes(attribute.String("consent_challenge_id", consentRequestID)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, p.revokeConsentSession("consent_challenge_id = ?", consentRequestID))
}

func (p *ConsentPersister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
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

// RevokeSubjectLoginSessionBatchSize bounds how many login sessions a single
// RevokeSubjectLoginSession DELETE removes per statement, so the ON DELETE SET
// NULL cascade into hydra_oauth2_flow stays well under CockroachDB's
// per-transaction intent budget (kv.transaction.max_intents_bytes). It is a var
// (not a const) only so tests can lower it to exercise multi-batch behavior.
var RevokeSubjectLoginSessionBatchSize = 100

func (p *Persister) RevokeSubjectLoginSession(ctx context.Context, subject string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectLoginSession")
	defer otelx.End(span, &err)

	nid := p.NetworkID(ctx)
	total := 0
	batches := 0
	for {
		// Honor the caller's deadline/cancellation between batches. Already
		// committed batches stay deleted; the caller can retry to finish.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return errors.WithStack(ctxErr)
		}

		// The nested SELECT is required because our MySQL version does not
		// support LIMIT in an IN/ANY subquery. This form is portable across all
		// supported dialects; the inner SELECT is served by the (subject, nid)
		// index. The outer DELETE repeats the nid filter as defense in depth so
		// the network boundary is enforced on the delete target itself, matching
		// QueryWithNetwork.
		var deleted int
		/* #nosec G201 - RevokeSubjectLoginSessionBatchSize is a package-level integer variable, not user input. */
		deleted, err = p.Connection(ctx).RawQuery(
			fmt.Sprintf(`DELETE FROM hydra_oauth2_authentication_session WHERE nid = ? AND id IN (
				SELECT id FROM (
					SELECT id FROM hydra_oauth2_authentication_session WHERE nid = ? AND subject = ? LIMIT %d
				) AS s
			)`, RevokeSubjectLoginSessionBatchSize),
			nid, nid, subject,
		).ExecWithCount()
		if err != nil {
			return sqlcon.HandleError(err)
		}

		total += deleted
		batches++
		p.l.Debugf("Revoking subject login sessions: %d deleted so far", total)

		if deleted < RevokeSubjectLoginSessionBatchSize {
			break // Last partial batch -> the subject is fully drained.
		}
	}

	span.SetAttributes(
		attribute.Int("rows_deleted", total),
		attribute.Int("batches", batches),
	)

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

func (p *Persister) GetForcedObfuscatedLoginSession(ctx context.Context, clientID, obfuscated string) (_ *consent.ForcedObfuscatedLoginSession, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetForcedObfuscatedLoginSession", trace.WithAttributes(attribute.String("client.id", clientID)))
	defer otelx.End(span, &err)

	var s consent.ForcedObfuscatedLoginSession

	if err := p.Connection(ctx).Where(
		"client_id = ? AND subject_obfuscated = ? AND nid = ?",
		clientID,
		obfuscated,
		p.NetworkID(ctx),
	).First(&s); errors.Is(err, sql.ErrNoRows) {
		return nil, errors.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &s, nil
}

func (p *ConsentPersister) CreateConsentSession(ctx context.Context, f *flow.Flow) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateConsentSession")
	defer otelx.End(span, &err)

	if f.NID != p.NetworkID(ctx) {
		return errors.WithStack(sqlcon.ErrNoRows())
	}
	return sqlcon.HandleError(p.Connection(ctx).Create(f))
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

		count, err := tx.RawQuery(
			`DELETE FROM hydra_oauth2_authentication_session WHERE id = ? AND nid = ?`,
			id, p.NetworkID(ctx),
		).ExecWithCount()
		if err != nil {
			return err
		}
		if count == 0 {
			// A concurrent transaction deleted the session between our
			// consistent read above and this delete. Surface ErrNoRows so
			// callers relying on the delete as an at-most-once gate (e.g.
			// back-channel logout execution) treat this call as the loser,
			// matching the DELETE ... RETURNING behavior on other dialects.
			return sql.ErrNoRows
		}
		return nil
	}); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &session, nil
}

func (p *ConsentPersister) FindGrantedAndRememberedConsentRequest(ctx context.Context, client, subject string) (_ *flow.Flow, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindGrantedAndRememberedConsentRequest")
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
AND (state = ? OR state IS NULL)
AND subject = ?
AND client_id = ?
AND consent_skip = FALSE
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

func (p *ConsentPersister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, pageOpts ...keysetpagination.Option) (_ []flow.Flow, _ *keysetpagination.Paginator, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindSubjectsGrantedConsentRequests")
	defer otelx.End(span, &err)

	paginator, err := keysetpagination.NewPaginator(append(pageOpts,
		keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(keysetpagination.Column{Name: "login_challenge", Value: ""})),
	)...)
	if err != nil {
		return nil, nil, err
	}

	var fs []flow.Flow
	err = p.QueryWithNetwork(ctx).
		Where("(state IN (?, ?) OR state IS NULL)", flow.FlowStateConsentUsed, flow.FlowStateConsentUnused).
		Where("subject = ?", subject).
		Where("consent_skip = FALSE").
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

func (p *ConsentPersister) FindSubjectsSessionGrantedConsentRequests(ctx context.Context, subject, sid string, pageOpts ...keysetpagination.Option) (_ []flow.Flow, _ *keysetpagination.Paginator, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindSubjectsSessionGrantedConsentRequests", trace.WithAttributes(attribute.String("sid", sid)))
	defer otelx.End(span, &err)

	paginator, err := keysetpagination.NewPaginator(append(pageOpts,
		keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(keysetpagination.Column{Name: "login_challenge", Value: ""})),
	)...)
	if err != nil {
		return nil, nil, err
	}

	var fs []flow.Flow
	err = p.QueryWithNetwork(ctx).
		Where("(state IN (?, ?) OR state IS NULL)", flow.FlowStateConsentUsed, flow.FlowStateConsentUnused).
		Where("subject = ?", subject).
		Where("login_session_id = ?", sid).
		Where("consent_skip = FALSE").
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

func (p *ConsentPersister) ListClientsWithLogoutURLsForSubjectAndSID(ctx context.Context, subject, sid string) (withFrontChannelURL, withBackChannelURL []client.Client, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ListClientsWithLogoutURLsForSubjectAndSID",
		trace.WithAttributes(attribute.String("sid", sid)))
	defer otelx.End(span, &err)

	defer func() {
		span.SetAttributes(
			attribute.Int("withFrontChannelURL", len(withFrontChannelURL)),
			attribute.Int("withBackChannelURL", len(withBackChannelURL)))
	}()

	var (
		cols                   = pop.NewModel(new(client.Client), ctx).Columns().Readable()
		clientTable, flowTable = clientFlowTableNamesWithQueryHint(p.Connection(ctx).Dialect.Name())

		q = fmt.Sprintf(`
		SELECT %s FROM %s c
		WHERE id IN (
			SELECT client_id
			FROM %s f
			WHERE
				f.nid = ?
				AND f.login_session_id = ?
				AND f.subject = ?
		)
		AND	c.nid = ?
		AND (
			c.frontchannel_logout_uri != '' OR c.backchannel_logout_uri != ''
		)`,
			cols.QuotedString(p.Connection(ctx).Dialect),
			clientTable,
			flowTable)

		nid = p.NetworkID(ctx)
		cs  []client.Client
	)

	err = p.Connection(ctx).RawQuery(q, nid, sid, subject, nid).All(&cs)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, sqlcon.HandleError(err)
	}

	for _, c := range cs {
		if c.FrontChannelLogoutURI != "" {
			withFrontChannelURL = append(withFrontChannelURL, c)
		}
		if c.BackChannelLogoutURI != "" {
			withBackChannelURL = append(withBackChannelURL, c)
		}
	}

	return withFrontChannelURL, withBackChannelURL, nil
}

func clientFlowTableNamesWithQueryHint(dialect string) (clientTable, flowTable string) {
	switch dialect {
	case "cockroach":
		return "hydra_client@primary", "hydra_oauth2_flow@hydra_oauth2_flow_login_session_subject_idx"
	// TODO: more
	default:
		return "hydra_client", "hydra_oauth2_flow"
	}
}

// logoutPayload is the wire format used when AEAD-encoding logout challenges
// and verifiers. It carries internal fields that must not appear in the public
// oAuth2LogoutRequest HTTP response (PostLogoutRedirectURI) and the network ID
// used to enforce tenant isolation across the encrypted blob.
type logoutPayload struct {
	*flow.LogoutRequest
	PostLogoutRedirectURI string    `json:"redir_url,omitempty"`
	NID                   uuid.UUID `json:"nid"`
}

func newLogoutPayload(req *flow.LogoutRequest, nid uuid.UUID) *logoutPayload {
	return &logoutPayload{
		LogoutRequest:         req,
		PostLogoutRedirectURI: req.PostLogoutRedirectURI,
		NID:                   nid,
	}
}

// ToLogoutRequest validates the payload's NID against currentNID and returns
// the embedded LogoutRequest with PostLogoutRedirectURI restored. Returns
// x.ErrNotFound if the NID does not match, mirroring how flow.decodeFlow
// rejects cross-tenant challenges for login, consent, and device flows, and
// ErrorLogoutFlowExpired if the embedded expiry has passed. Stateless blobs
// cannot be revoked, so the expiry is their only time bound and every decode
// path must enforce it.
func (payload *logoutPayload) ToLogoutRequest(currentNID uuid.UUID) (*flow.LogoutRequest, error) {
	if payload.NID != currentNID {
		return nil, errors.WithStack(x.ErrNotFound.WithDescription("Network IDs are not matching."))
	}
	if payload.LogoutRequest == nil {
		payload.LogoutRequest = &flow.LogoutRequest{}
	}
	if payload.ExpiresAt == nil || payload.ExpiresAt.Before(time.Now()) {
		return nil, errors.WithStack(flow.ErrorLogoutFlowExpired)
	}
	payload.LogoutRequest.PostLogoutRedirectURI = payload.PostLogoutRedirectURI
	return payload.LogoutRequest, nil
}

// looksLikeLegacyUUID reports whether s parses as a UUID. Before the
// stateless logout flow shipped, logout challenges and verifiers were
// generated via uuid.New() and stored in the hydra_oauth2_logout_request
// table. During the rollout window — until any such legacy
// challenge/verifier has aged past the configured logout request TTL — we
// surface UUID-shaped inputs as "expired" so clients retry through their
// normal logout-flow restart path, rather than as a generic decryption
// failure.
func looksLikeLegacyUUID(s string) bool {
	_, err := uuid.FromString(s)
	return err == nil
}

func (p *Persister) CreateLogoutChallenge(ctx context.Context, request *flow.LogoutRequest) (challenge string, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateLogoutChallenge")
	defer otelx.End(span, &err)

	// The challenge IS the AEAD blob we're about to produce, so encoding it
	// into the payload would be self-referential and bloat the size. Blank it
	// before encoding; GetLogoutRequest will populate it again on decode.
	request.Challenge = ""
	challenge, err = flow.Encode(ctx, p.r.FlowCipher(), newLogoutPayload(request, p.NetworkID(ctx)), flow.AsLogoutChallenge)
	if err != nil {
		return "", errors.WithStack(fosite.ErrServerError.WithWrap(err).WithHintf("Failed to encrypt the logout challenge."))
	}
	return challenge, nil
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (verifier string, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AcceptLogoutRequest")
	defer otelx.End(span, &err)

	payload, err := flow.Decode[logoutPayload](ctx, p.r.FlowCipher(), challenge, flow.AsLogoutChallenge)
	if err != nil {
		if looksLikeLegacyUUID(challenge) {
			return "", errors.WithStack(flow.ErrorLogoutFlowExpired)
		}
		return "", errors.WithStack(x.ErrNotFound.WithWrap(err).WithHintf("Failed to decrypt the logout challenge."))
	}
	nid := p.NetworkID(ctx)
	req, err := payload.ToLogoutRequest(nid)
	if err != nil {
		return "", err
	}
	// The verifier IS the AEAD blob we're about to produce. Avoid embedding the
	// caller-provided challenge inside the verifier payload (self-referential
	// and bloats size).
	req.Challenge = ""
	// The verifier is consumed machine-to-machine by the logout completion
	// handler, which only needs the subject, session ID, RP flag, expiry, and
	// post-logout redirect URI. Drop the client metadata and the original
	// request URL (which can embed a full id_token_hint JWT) to keep the
	// verifier — and thus the redirect URL it travels in — short.
	req.Client = nil
	req.RequestURL = ""

	verifier, err = flow.Encode(ctx, p.r.FlowCipher(), newLogoutPayload(req, nid), flow.AsLogoutVerifier)
	if err != nil {
		return "", errors.WithStack(fosite.ErrServerError.WithWrap(err).WithHintf("Failed to encrypt the logout verifier."))
	}

	return verifier, nil
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RejectLogoutRequest")
	defer otelx.End(span, &err)

	payload, err := flow.Decode[logoutPayload](ctx, p.r.FlowCipher(), challenge, flow.AsLogoutChallenge)
	if err != nil {
		if looksLikeLegacyUUID(challenge) {
			return errors.WithStack(flow.ErrorLogoutFlowExpired)
		}
		return errors.WithStack(x.ErrNotFound.WithWrap(err).WithHintf("Failed to decrypt the logout challenge."))
	}
	if _, err := payload.ToLogoutRequest(p.NetworkID(ctx)); err != nil {
		return err
	}
	return nil
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (_ *flow.LogoutRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetLogoutRequest")
	defer otelx.End(span, &err)

	payload, err := flow.Decode[logoutPayload](ctx, p.r.FlowCipher(), challenge, flow.AsLogoutChallenge)
	if err != nil {
		if looksLikeLegacyUUID(challenge) {
			return nil, errors.WithStack(flow.ErrorLogoutFlowExpired)
		}
		return nil, errors.WithStack(x.ErrNotFound.WithWrap(err).WithHintf("Failed to decrypt the logout challenge."))
	}
	req, err := payload.ToLogoutRequest(p.NetworkID(ctx))
	if err != nil {
		return nil, err
	}
	req.Challenge = challenge
	return req, nil
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (_ *flow.LogoutRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.VerifyAndInvalidateLogoutRequest")
	defer otelx.End(span, &err)

	payload, err := flow.Decode[logoutPayload](ctx, p.r.FlowCipher(), verifier, flow.AsLogoutVerifier)
	if err != nil {
		if looksLikeLegacyUUID(verifier) {
			return nil, errors.WithStack(flow.ErrorLogoutFlowExpired)
		}
		return nil, errors.WithStack(x.ErrNotFound.WithWrap(err).WithHintf("Failed to decrypt the logout verifier."))
	}
	lr, err := payload.ToLogoutRequest(p.NetworkID(ctx))
	if err != nil {
		return nil, err
	}

	return lr, nil
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
			(state != ? AND state IS NOT NULL)
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

		if !errors.Is(err, sqlcon.ErrUniqueViolation()) {
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
