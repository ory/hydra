package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/sqlxx"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/flow"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon"
)

var _ consent.Manager = &Persister{}

func (p *Persister) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectConsentSession")
	defer span.End()

	return p.transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ?", user))
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectClientConsentSession")
	defer span.End()

	return p.transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ? AND client_id = ?", user, client))
}

func (p *Persister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
	return func(ctx context.Context, c *pop.Connection) error {
		fs := make([]*flow.Flow, 0)
		if err := p.QueryWithNetwork(ctx).
			Where(whereStmt, whereArgs...).
			Select("consent_challenge_id").
			All(&fs); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}

			return sqlcon.HandleError(err)
		}

		var count int
		for _, f := range fs {
			if err := p.RevokeAccessToken(ctx, f.ConsentChallengeID.String()); errors.Is(err, fosite.ErrNotFound) {
				// do nothing
			} else if err != nil {
				return err
			}

			if err := p.RevokeRefreshToken(ctx, f.ConsentChallengeID.String()); errors.Is(err, fosite.ErrNotFound) {
				// do nothing
			} else if err != nil {
				return err
			}

			localCount, err := c.RawQuery("DELETE FROM hydra_oauth2_flow WHERE consent_challenge_id = ? AND nid = ?", f.ConsentChallengeID, p.NetworkID(ctx)).ExecWithCount()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return errorsx.WithStack(x.ErrNotFound)
				}
				return sqlcon.HandleError(err)
			}

			// If there are no sessions to revoke we should return an error to indicate to the caller
			// that the request failed.
			count += localCount
		}

		if count == 0 {
			return errorsx.WithStack(x.ErrNotFound)
		}

		return nil
	}
}

func (p *Persister) RevokeSubjectLoginSession(ctx context.Context, subject string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RevokeSubjectLoginSession")
	defer span.End()

	err := p.QueryWithNetwork(ctx).Where("subject = ?", subject).Delete(&consent.LoginSession{})
	if err != nil {
		return sqlcon.HandleError(err)
	}

	// This confuses people, see https://github.com/ory/hydra/issues/1168
	//
	// count, _ := rows.RowsAffected()
	// if count == 0 {
	// 	 return errorsx.WithStack(x.ErrNotFound)
	// }

	return nil
}

func (p *Persister) CreateForcedObfuscatedLoginSession(ctx context.Context, session *consent.ForcedObfuscatedLoginSession) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateForcedObfuscatedLoginSession")
	defer span.End()

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
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

func (p *Persister) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*consent.ForcedObfuscatedLoginSession, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetForcedObfuscatedLoginSession")
	defer span.End()

	var s consent.ForcedObfuscatedLoginSession

	if err := p.Connection(ctx).Where(
		"client_id = ? AND subject_obfuscated = ? AND nid = ?",
		client,
		obfuscated,
		p.NetworkID(ctx),
	).First(&s); errors.Is(err, sql.ErrNoRows) {
		return nil, errorsx.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &s, nil
}

// CreateConsentRequest configures fields that are introduced or changed in the
// consent request. It doesn't touch fields that would be copied from the login
// request.
func (p *Persister) CreateConsentRequest(ctx context.Context, req *consent.OAuth2ConsentRequest) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateConsentRequest")
	defer span.End()

	c, err := p.Connection(ctx).RawQuery(`
UPDATE hydra_oauth2_flow
SET
	state = ?,
	consent_challenge_id = ?,
	consent_skip = ?,
	consent_verifier = ?,
	consent_csrf = ?
WHERE login_challenge = ? AND nid = ?;
`,
		flow.FlowStateConsentInitialized,
		sqlxx.NullString(req.ID),
		req.Skip,
		req.Verifier,
		req.CSRF,
		req.LoginChallenge.String(),
		p.NetworkID(ctx),
	).ExecWithCount()
	if err != nil {
		return sqlcon.HandleError(err)
	}
	if c != 1 {
		return errorsx.WithStack(x.ErrNotFound)
	}
	return nil
}

func (p *Persister) GetFlowByConsentChallenge(ctx context.Context, challenge string) (*flow.Flow, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetFlowByConsentChallenge")
	defer span.End()

	f := &flow.Flow{}

	if err := sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("consent_challenge_id = ?", challenge).First(f)); err != nil {
		return nil, err
	}

	return f, nil
}

func (p *Persister) GetConsentRequest(ctx context.Context, challenge string) (*consent.OAuth2ConsentRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetConsentRequest")
	defer span.End()

	f, err := p.GetFlowByConsentChallenge(ctx, challenge)
	if err != nil {
		if errors.Is(err, sqlcon.ErrNoRows) {
			return nil, errorsx.WithStack(x.ErrNotFound)
		}
		return nil, err
	}

	return f.GetConsentRequest(), nil
}

func (p *Persister) CreateLoginRequest(ctx context.Context, req *consent.LoginRequest) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateLoginRequest")
	defer span.End()

	f := flow.NewFlow(req)
	return sqlcon.HandleError(p.CreateWithNetwork(ctx, f))
}

func (p *Persister) GetFlow(ctx context.Context, login_challenge string) (*flow.Flow, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetFlow")
	defer span.End()

	var f flow.Flow
	return &f, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := p.QueryWithNetwork(ctx).Where("login_challenge = ?", login_challenge).First(&f); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}

		return nil
	})
}

func (p *Persister) GetLoginRequest(ctx context.Context, login_challenge string) (*consent.LoginRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetLoginRequest")
	defer span.End()

	var lr *consent.LoginRequest
	return lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		var f flow.Flow
		if err := p.QueryWithNetwork(ctx).Where("login_challenge = ?", login_challenge).First(&f); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}
		lr = f.GetLoginRequest()

		return nil
	})
}

func (p *Persister) HandleConsentRequest(ctx context.Context, r *consent.AcceptOAuth2ConsentRequest) (*consent.OAuth2ConsentRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.HandleConsentRequest")
	defer span.End()

	f := &flow.Flow{}

	if err := sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("consent_challenge_id = ?", r.ID).First(f)); errors.Is(err, sqlcon.ErrNoRows) {
		return nil, err
	}

	if err := f.HandleConsentRequest(r); err != nil {
		return nil, errorsx.WithStack(err)
	}

	_, err := p.UpdateWithNetwork(ctx, f)
	if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetConsentRequest(ctx, r.ID)
}

func (p *Persister) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.AcceptOAuth2ConsentRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.VerifyAndInvalidateConsentRequest")
	defer span.End()

	var r consent.AcceptOAuth2ConsentRequest
	return &r, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		var f flow.Flow
		if err := p.QueryWithNetwork(ctx).Where("consent_verifier = ?", verifier).First(&f); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := f.InvalidateConsentRequest(); err != nil {
			return errorsx.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
		}

		r = *f.GetHandledConsentRequest()
		_, err := p.UpdateWithNetwork(ctx, &f)
		return err
	})
}

func (p *Persister) HandleLoginRequest(ctx context.Context, challenge string, r *consent.HandledLoginRequest) (lr *consent.LoginRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.HandleLoginRequest")
	defer span.End()

	return lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		f, err := p.GetFlow(ctx, challenge)
		if err != nil {
			return sqlcon.HandleError(err)
		}
		err = f.HandleLoginRequest(r)
		if err != nil {
			return err
		}

		_, err = p.UpdateWithNetwork(ctx, f)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		lr, err = p.GetLoginRequest(ctx, challenge)
		return sqlcon.HandleError(err)
	})
}

func (p *Persister) VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*consent.HandledLoginRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.VerifyAndInvalidateLoginRequest")
	defer span.End()

	var d consent.HandledLoginRequest
	return &d, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		var f flow.Flow
		if err := p.QueryWithNetwork(ctx).Where("login_verifier = ?", verifier).First(&f); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := f.InvalidateLoginRequest(); err != nil {
			return errorsx.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
		}

		d = f.GetHandledLoginRequest()
		_, err := p.UpdateWithNetwork(ctx, &f)
		return sqlcon.HandleError(err)
	})
}

func (p *Persister) GetRememberedLoginSession(ctx context.Context, id string) (*consent.LoginSession, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetRememberedLoginSession")
	defer span.End()

	var s consent.LoginSession

	if err := p.QueryWithNetwork(ctx).Where("remember = TRUE").Find(&s, id); errors.Is(err, sql.ErrNoRows) {
		return nil, errorsx.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &s, nil
}

func (p *Persister) ConfirmLoginSession(ctx context.Context, id string, authenticatedAt time.Time, subject string, remember bool) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ConfirmLoginSession")
	defer span.End()

	_, err := p.Connection(ctx).Where("id = ? AND nid = ?", id, p.NetworkID(ctx)).UpdateQuery(&consent.LoginSession{
		AuthenticatedAt: sqlxx.NullTime(authenticatedAt),
		Subject:         subject,
		Remember:        remember,
	}, "authenticated_at", "subject", "remember")
	return sqlcon.HandleError(err)
}

func (p *Persister) CreateLoginSession(ctx context.Context, session *consent.LoginSession) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateLoginSession")
	defer span.End()

	return sqlcon.HandleError(p.CreateWithNetwork(ctx, session))
}

func (p *Persister) DeleteLoginSession(ctx context.Context, id string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteLoginSession")
	defer span.End()

	count, err := p.Connection(ctx).RawQuery("DELETE FROM hydra_oauth2_authentication_session WHERE id=? AND nid = ?", id, p.NetworkID(ctx)).ExecWithCount()
	if count == 0 {
		return errorsx.WithStack(x.ErrNotFound)
	} else {
		return sqlcon.HandleError(err)
	}
}

func (p *Persister) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, subject string) ([]consent.AcceptOAuth2ConsentRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindGrantedAndRememberedConsentRequests")
	defer span.End()

	rs := make([]consent.AcceptOAuth2ConsentRequest, 0)

	return rs, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		f := &flow.Flow{}

		if err := c.
			Where(
				strings.TrimSpace(fmt.Sprintf(`
(state = %d OR state = %d) AND
subject = ? AND
client_id = ? AND
consent_skip=FALSE AND
consent_error='{}' AND
consent_remember=TRUE AND
nid = ?`, flow.FlowStateConsentUsed, flow.FlowStateConsentUnused,
				)),
				subject, client, p.NetworkID(ctx)).
			Order("requested_at DESC").
			Limit(1).
			First(f); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(consent.ErrNoPreviousConsentFound)
			}
			return sqlcon.HandleError(err)
		}

		var err error
		rs, err = p.filterExpiredConsentRequests(ctx, []consent.AcceptOAuth2ConsentRequest{*f.GetHandledConsentRequest()})
		return err
	})
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]consent.AcceptOAuth2ConsentRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FindSubjectsGrantedConsentRequests")
	defer span.End()

	var fs []flow.Flow
	c := p.Connection(ctx)

	if err := c.
		Where(
			strings.TrimSpace(fmt.Sprintf(`
(state = %d OR state = %d) AND
subject = ? AND
consent_skip=FALSE AND
consent_error='{}' AND
nid = ?`, flow.FlowStateConsentUsed, flow.FlowStateConsentUnused,
			)),
			subject, p.NetworkID(ctx)).
		Order("requested_at DESC").
		Paginate(offset/limit+1, limit).
		All(&fs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.WithStack(consent.ErrNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	var rs []consent.AcceptOAuth2ConsentRequest
	for _, f := range fs {
		rs = append(rs, *f.GetHandledConsentRequest())
	}

	return p.filterExpiredConsentRequests(ctx, rs)
}

func (p *Persister) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CountSubjectsGrantedConsentRequests")
	defer span.End()

	n, err := p.Connection(ctx).
		Where(
			strings.TrimSpace(fmt.Sprintf(`
(state = %d OR state = %d) AND
subject = ? AND
consent_skip=FALSE AND
consent_error='{}' AND
nid = ?`, flow.FlowStateConsentUsed, flow.FlowStateConsentUnused,
			)),
			subject, p.NetworkID(ctx)).
		Count(&flow.Flow{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) filterExpiredConsentRequests(ctx context.Context, requests []consent.AcceptOAuth2ConsentRequest) ([]consent.AcceptOAuth2ConsentRequest, error) {
	_, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.filterExpiredConsentRequests")
	defer span.End()

	var result []consent.AcceptOAuth2ConsentRequest
	for _, v := range requests {
		if v.RememberFor > 0 && v.RequestedAt.Add(time.Duration(v.RememberFor)*time.Second).Before(time.Now().UTC()) {
			continue
		}
		result = append(result, v)
	}
	if len(result) == 0 {
		return nil, errorsx.WithStack(consent.ErrNoPreviousConsentFound)
	}
	return result, nil
}

func (p *Persister) ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ListUserAuthenticatedClientsWithFrontChannelLogout")
	defer span.End()

	return p.listUserAuthenticatedClients(ctx, subject, sid, "front")
}

func (p *Persister) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.ListUserAuthenticatedClientsWithBackChannelLogout")
	defer span.End()
	return p.listUserAuthenticatedClients(ctx, subject, sid, "back")
}

func (p *Persister) listUserAuthenticatedClients(ctx context.Context, subject, sid, channel string) ([]client.Client, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.listUserAuthenticatedClients")
	defer span.End()

	var cs []client.Client
	return cs, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.RawQuery(
			/* #nosec G201 - channel can either be "front" or "back" */
			fmt.Sprintf(`
SELECT DISTINCT c.* FROM hydra_client as c
JOIN hydra_oauth2_flow as f ON (c.id = f.client_id)
WHERE
	f.subject=? AND
	c.%schannel_logout_uri!='' AND
	c.%schannel_logout_uri IS NOT NULL AND
	f.login_session_id = ? AND
	f.nid = ? AND
	c.nid = ?`,
				channel,
				channel,
			),
			subject,
			sid,
			p.NetworkID(ctx),
			p.NetworkID(ctx),
		).All(&cs); err != nil {
			return sqlcon.HandleError(err)
		}

		return nil
	})
}

func (p *Persister) CreateLogoutRequest(ctx context.Context, request *consent.LogoutRequest) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateLogoutRequest")
	defer span.End()

	return errorsx.WithStack(p.CreateWithNetwork(ctx, request))
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AcceptLogoutRequest")
	defer span.End()

	if err := p.Connection(ctx).RawQuery("UPDATE hydra_oauth2_logout_request SET accepted=true, rejected=false WHERE challenge=? AND nid = ?", challenge, p.NetworkID(ctx)).Exec(); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetLogoutRequest(ctx, challenge)
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.RejectLogoutRequest")
	defer span.End()

	count, err := p.Connection(ctx).
		RawQuery("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=? AND nid = ?", challenge, p.NetworkID(ctx)).
		ExecWithCount()
	if count == 0 {
		return errorsx.WithStack(x.ErrNotFound)
	} else {
		return errorsx.WithStack(err)
	}
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetLogoutRequest")
	defer span.End()

	var lr consent.LogoutRequest
	return &lr, sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("challenge = ? AND rejected = FALSE", challenge).First(&lr))
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*consent.LogoutRequest, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.VerifyAndInvalidateLogoutRequest")
	defer span.End()

	var lr consent.LogoutRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if count, err := c.RawQuery(
			"UPDATE hydra_oauth2_logout_request SET was_used=TRUE WHERE nid = ? AND verifier=? AND was_used=FALSE AND accepted=TRUE AND rejected=FALSE",
			p.NetworkID(ctx),
			verifier,
		).ExecWithCount(); count == 0 && err == nil {
			return errorsx.WithStack(x.ErrNotFound)
		} else if err != nil {
			return sqlcon.HandleError(err)
		}

		err := sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("verifier=?", verifier).First(&lr))
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *Persister) FlushInactiveLoginConsentRequests(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FlushInactiveLoginConsentRequests")
	defer span.End()

	/* #nosec G201 table is static */
	var f flow.Flow

	// The value of notAfter should be the minimum between input parameter and request max expire based on its configured age
	requestMaxExpire := time.Now().Add(-p.config.ConsentRequestMaxAge(ctx))
	if requestMaxExpire.Before(notAfter) {
		notAfter = requestMaxExpire
	}

	challenges := []string{}
	queryFormat := `
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
	LIMIT %[1]d
	`

	// Select up to [limit] flows that can be safely deleted, i.e. flows that meet
	// the following criteria:
	// - flow.state is anything between FlowStateLoginInitialized and FlowStateConsentUnused (unhandled)
	// - flow.login_error has valid error (login rejected)
	// - flow.consent_error has valid error (consent rejected)
	// AND timed-out
	// - flow.requested_at < minimum of ttl.login_consent_request and notAfter
	q := p.Connection(ctx).RawQuery(fmt.Sprintf(queryFormat, limit), flow.FlowStateConsentUsed, notAfter, p.NetworkID(ctx))

	if err := q.All(&challenges); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	}

	// Delete in batch consent requests and their references in cascade
	for i := 0; i < len(challenges); i += batchSize {
		j := i + batchSize
		if j > len(challenges) {
			j = len(challenges)
		}

		q := p.Connection(ctx).RawQuery(
			fmt.Sprintf("DELETE FROM %s WHERE login_challenge in (?) AND nid = ?", (&f).TableName()),
			challenges[i:j],
			p.NetworkID(ctx),
		)

		if err := q.Exec(); err != nil {
			return sqlcon.HandleError(err)
		}
	}

	return nil
}
