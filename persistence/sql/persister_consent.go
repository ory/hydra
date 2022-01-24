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
	return p.transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ?", user))
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	return p.transaction(ctx, p.revokeConsentSession("consent_challenge_id IS NOT NULL AND subject = ? AND client_id = ?", user, client))
}

func (p *Persister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
	return func(ctx context.Context, c *pop.Connection) error {
		fs := make([]*flow.Flow, 0)
		if err := c.
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

			localCount, err := c.RawQuery("DELETE FROM hydra_oauth2_flow WHERE consent_challenge_id = ?", f.ConsentChallengeID).ExecWithCount()
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
	if err := p.Connection(ctx).RawQuery("DELETE FROM hydra_oauth2_authentication_session WHERE subject = ?", subject).Exec(); errors.Is(err, sql.ErrNoRows) {
		return errorsx.WithStack(x.ErrNotFound)
	} else if err != nil {
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
	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.RawQuery(
			"DELETE FROM hydra_oauth2_obfuscated_authentication_session WHERE client_id = ? AND subject = ?",
			session.ClientID,
			session.Subject,
		).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}

		return sqlcon.HandleError(c.RawQuery(
			"INSERT INTO hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated) VALUES (?, ?, ?)",
			session.Subject,
			session.ClientID,
			session.SubjectObfuscated,
		).Exec())
	})
}

func (p *Persister) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*consent.ForcedObfuscatedLoginSession, error) {
	var s consent.ForcedObfuscatedLoginSession

	if err := p.Connection(ctx).Where(
		"client_id = ? AND subject_obfuscated = ?",
		client,
		obfuscated,
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
func (p *Persister) CreateConsentRequest(ctx context.Context, req *consent.ConsentRequest) error {
	c, err := p.Connection(ctx).RawQuery(`
UPDATE hydra_oauth2_flow
SET
	state = ?,
	consent_challenge_id = ?,
	consent_skip = ?,
	consent_verifier = ?,
	consent_csrf = ?
WHERE login_challenge = ?;
`,
		flow.FlowStateConsentInitialized,
		sqlxx.NullString(req.ID),
		req.Skip,
		req.Verifier,
		req.CSRF,
		req.LoginChallenge.String(),
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
	f := &flow.Flow{}

	if err := f.FindByConsentChallengeID(p.Connection(ctx), challenge); err != nil {
		return nil, err
	}

	return f, nil
}

func (p *Persister) GetConsentRequest(ctx context.Context, challenge string) (*consent.ConsentRequest, error) {
	f, err := p.GetFlowByConsentChallenge(ctx, challenge)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return f.GetConsentRequest(), nil
}

func (p *Persister) CreateLoginRequest(ctx context.Context, req *consent.LoginRequest) error {
	f := flow.NewFlow(req)
	return sqlcon.HandleError(p.Connection(ctx).Create(f))
}

func (p *Persister) GetFlow(ctx context.Context, challenge string) (*flow.Flow, error) {
	var f flow.Flow
	return &f, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.Find(&f, challenge); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}

		return nil
	})
}

func (p *Persister) GetLoginRequest(ctx context.Context, challenge string) (*consent.LoginRequest, error) {
	var lr *consent.LoginRequest
	return lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		f := &flow.Flow{}
		if err := c.Find(f, challenge); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}
		lr = f.GetLoginRequest()

		return nil
	})
}

func (p *Persister) HandleConsentRequest(ctx context.Context, challenge string, r *consent.HandledConsentRequest) (*consent.ConsentRequest, error) {
	c := p.Connection(ctx)
	f := &flow.Flow{}

	if err := f.FindByConsentChallengeID(c, r.ID); err == sqlcon.ErrNoRows {
		return nil, sqlcon.HandleError(err)
	}

	if err := f.HandleConsentRequest(r); err != nil {
		return nil, errorsx.WithStack(err)
	}

	if err := c.Update(f); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetConsentRequest(ctx, challenge)
}

func (p *Persister) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.HandledConsentRequest, error) {
	var r consent.HandledConsentRequest
	return &r, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		f := &flow.Flow{}
		if err := c.Where("consent_verifier = ?", verifier).Select("*").First(f); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := f.InvalidateConsentRequest(); err != nil {
			return errorsx.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
		}

		r = *f.GetHandledConsentRequest()
		return c.Update(f)
	})
}

func (p *Persister) HandleLoginRequest(ctx context.Context, challenge string, r *consent.HandledLoginRequest) (lr *consent.LoginRequest, err error) {
	return lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		f, err := p.GetFlow(ctx, challenge)
		if err != nil {
			return sqlcon.HandleError(err)
		}
		err = f.HandleLoginRequest(r)
		if err != nil {
			return err
		}

		err = c.Update(f)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		lr, err = p.GetLoginRequest(ctx, challenge)
		return sqlcon.HandleError(err)
	})
}

func (p *Persister) VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*consent.HandledLoginRequest, error) {
	var d consent.HandledLoginRequest
	return &d, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		var f flow.Flow
		if err := c.Where("login_verifier = ?", verifier).Select("*").First(&f); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := f.InvalidateLoginRequest(); err != nil {
			return errorsx.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
		}

		d = f.GetHandledLoginRequest()
		return sqlcon.HandleError(c.Update(&f))
	})
}

func (p *Persister) GetRememberedLoginSession(ctx context.Context, id string) (*consent.LoginSession, error) {
	var s consent.LoginSession

	if err := p.Connection(ctx).Where("remember = TRUE").Find(&s, id); errors.Is(err, sql.ErrNoRows) {
		return nil, errorsx.WithStack(x.ErrNotFound)
	} else if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return &s, nil
}

func (p *Persister) ConfirmLoginSession(ctx context.Context, id string, authenticatedAt time.Time, subject string, remember bool) error {
	return sqlcon.HandleError(
		p.Connection(ctx).Update(&consent.LoginSession{
			ID:              id,
			AuthenticatedAt: sqlxx.NullTime(authenticatedAt),
			Subject:         subject,
			Remember:        remember,
		}))
}

func (p *Persister) CreateLoginSession(ctx context.Context, session *consent.LoginSession) error {
	return sqlcon.HandleError(p.Connection(ctx).Create(session))
}

func (p *Persister) DeleteLoginSession(ctx context.Context, id string) error {
	return sqlcon.HandleError(
		p.Connection(ctx).Destroy(
			&consent.LoginSession{ID: id},
		))
}

func (p *Persister) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, subject string) ([]consent.HandledConsentRequest, error) {
	rs := make([]consent.HandledConsentRequest, 0)

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
consent_remember=TRUE`, flow.FlowStateConsentUsed, flow.FlowStateConsentUnused,
				)),
				subject, client).
			Order("requested_at DESC").
			Limit(1).
			First(f); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(consent.ErrNoPreviousConsentFound)
			}
			return sqlcon.HandleError(err)
		}

		var err error
		rs, err = p.filterExpiredConsentRequests(ctx, []consent.HandledConsentRequest{*f.GetHandledConsentRequest()})
		return err
	})
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]consent.HandledConsentRequest, error) {
	var fs []flow.Flow
	c := p.Connection(ctx)

	if err := c.
		Where(
			strings.TrimSpace(fmt.Sprintf(`
(state = %d OR state = %d) AND
subject = ? AND
consent_skip=FALSE AND
consent_error='{}'`, flow.FlowStateConsentUsed, flow.FlowStateConsentUnused,
			)),
			subject).
		Order("requested_at DESC").
		Paginate(offset/limit+1, limit).
		All(&fs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.WithStack(consent.ErrNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	var rs []consent.HandledConsentRequest
	for _, f := range fs {
		rs = append(rs, *f.GetHandledConsentRequest())
	}

	return p.filterExpiredConsentRequests(ctx, rs)
}

func (p *Persister) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	n, err := p.Connection(ctx).
		Where(
			strings.TrimSpace(fmt.Sprintf(`
(state = %d OR state = %d) AND
subject = ? AND
consent_skip=FALSE AND
consent_error='{}'`, flow.FlowStateConsentUsed, flow.FlowStateConsentUnused,
			)),
			subject).
		Count(&flow.Flow{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) filterExpiredConsentRequests(ctx context.Context, requests []consent.HandledConsentRequest) ([]consent.HandledConsentRequest, error) {
	var result []consent.HandledConsentRequest
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
	return p.listUserAuthenticatedClients(ctx, subject, sid, "front")
}

func (p *Persister) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	return p.listUserAuthenticatedClients(ctx, subject, sid, "back")
}

func (p *Persister) listUserAuthenticatedClients(ctx context.Context, subject, sid, channel string) ([]client.Client, error) {
	var cs []client.Client
	return cs, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.RawQuery(
			/* #nosec G201 - channel can either be "front" or "back" */
			fmt.Sprintf(`SELECT DISTINCT c.* FROM hydra_client as c JOIN hydra_oauth2_flow as f ON (c.id = f.client_id) WHERE f.subject=? AND c.%schannel_logout_uri!='' AND c.%schannel_logout_uri IS NOT NULL AND f.login_session_id = ?`,
				channel,
				channel,
			),
			subject,
			sid,
		).All(&cs); err != nil {
			return sqlcon.HandleError(err)
		}

		return nil
	})
}

func (p *Persister) CreateLogoutRequest(ctx context.Context, request *consent.LogoutRequest) error {
	return errorsx.WithStack(p.Connection(ctx).Create(request))
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	if err := p.Connection(ctx).RawQuery("UPDATE hydra_oauth2_logout_request SET accepted=true, rejected=false WHERE challenge=?", challenge).Exec(); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetLogoutRequest(ctx, challenge)
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) error {
	return errorsx.WithStack(
		p.Connection(ctx).
			RawQuery("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=?", challenge).
			Exec())
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	var lr consent.LogoutRequest
	return &lr, sqlcon.HandleError(p.Connection(ctx).Where("challenge = ? AND rejected = FALSE", challenge).First(&lr))
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*consent.LogoutRequest, error) {
	var lr consent.LogoutRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.Where("verifier=? AND was_used=FALSE AND accepted=TRUE AND rejected=FALSE", verifier).Select("challenge").First(&lr); err != nil {
			if err == sql.ErrNoRows {
				return errorsx.WithStack(x.ErrNotFound)
			}

			return sqlcon.HandleError(err)
		}

		if err := c.RawQuery("UPDATE hydra_oauth2_logout_request SET was_used=TRUE WHERE verifier=?", verifier).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}

		updated, err := p.GetLogoutRequest(ctx, lr.ID)
		if err != nil {
			return err
		}

		lr = *updated
		return nil
	})
}

func (p *Persister) FlushInactiveLoginConsentRequests(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	/* #nosec G201 table is static */
	var f flow.Flow

	// The value of notAfter should be the minimum between input parameter and request max expire based on its configured age
	requestMaxExpire := time.Now().Add(-p.config.ConsentRequestMaxAge())
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
	q := p.Connection(ctx).RawQuery(fmt.Sprintf(queryFormat, limit), flow.FlowStateConsentUsed, notAfter)

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
			fmt.Sprintf("DELETE FROM %s WHERE login_challenge in (?)", (&f).TableName()),
			challenges[i:j],
		)

		if err := q.Exec(); err != nil {
			return sqlcon.HandleError(err)
		}
	}

	return nil
}
