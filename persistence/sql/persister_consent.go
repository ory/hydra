package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/sqlxx"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon"
)

var _ consent.Manager = &Persister{}

func (p *Persister) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	return p.transaction(ctx, p.revokeConsentSession("r.subject = ?", user))
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	return p.transaction(ctx, p.revokeConsentSession("r.subject = ? AND r.client_id = ?", user, client))
}

func (p *Persister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
	return func(ctx context.Context, c *pop.Connection) error {
		hrs := make([]*consent.HandledConsentRequest, 0)
		if err := c.
			Where(whereStmt, whereArgs...).
			Select("r.challenge").
			Join(
				fmt.Sprintf("%s AS r", consent.ConsentRequest{}.TableName()),
				fmt.Sprintf("r.challenge = %s.challenge", consent.HandledConsentRequest{}.TableName())).
			All(&hrs); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}

			return sqlcon.HandleError(err)
		}

		var count int
		for _, hr := range hrs {
			if err := p.RevokeAccessToken(ctx, hr.ID); errors.Is(err, fosite.ErrNotFound) {
				// do nothing
			} else if err != nil {
				return err
			}

			if err := p.RevokeRefreshToken(ctx, hr.ID); errors.Is(err, fosite.ErrNotFound) {
				// do nothing
			} else if err != nil {
				return err
			}

			// Since we ON DELETE CASCADE, hydra_oauth2_consent_request_handled will be removed automagically.
			localCount, err := c.RawQuery("DELETE FROM hydra_oauth2_consent_request WHERE challenge = ?", hr.ID).ExecWithCount()
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

func (p *Persister) CreateConsentRequest(ctx context.Context, req *consent.ConsentRequest) error {
	return errorsx.WithStack(p.Connection(ctx).Create(req))
}

func (p *Persister) GetConsentRequest(ctx context.Context, challenge string) (*consent.ConsentRequest, error) {
	r := &consent.ConsentRequest{}

	if err := r.FindInDB(p.Connection(ctx), challenge); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.WithStack(x.ErrNotFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return r, nil
}

func (p *Persister) CreateLoginRequest(ctx context.Context, req *consent.LoginRequest) error {
	return errorsx.WithStack(p.Connection(ctx).Create(req))
}

func (p *Persister) GetLoginRequest(ctx context.Context, challenge string) (*consent.LoginRequest, error) {
	var lr consent.LoginRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := (&lr).FindInDB(c, challenge); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}

		return nil
	})
}

func (p *Persister) HandleConsentRequest(ctx context.Context, challenge string, r *consent.HandledConsentRequest) (*consent.ConsentRequest, error) {
	c := p.Connection(ctx)

	if err := sqlcon.HandleError(c.Create(r)); errors.Is(err, sqlcon.ErrUniqueViolation) {
		hr := &consent.HandledConsentRequest{}
		if err := c.Find(hr, r.ID); err != nil {
			return nil, sqlcon.HandleError(err)
		}

		if hr.WasHandled {
			return nil, errorsx.WithStack(x.ErrConflict.WithHint("The consent request was already used and can no longer be changed."))
		}

		if err := c.Update(r); err != nil {
			return nil, sqlcon.HandleError(err)
		}
	} else if err != nil {
		return nil, err
	}

	return p.GetConsentRequest(ctx, challenge)
}

func (p *Persister) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.HandledConsentRequest, error) {
	var r consent.HandledConsentRequest
	return &r, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		var cr consent.ConsentRequest
		if err := c.Where("verifier = ?", verifier).Select("challenge", "client_id").First(&cr); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := c.Find(&r, cr.ID); err != nil {
			return sqlcon.HandleError(err)
		}

		if r.WasHandled {
			return errorsx.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already."))
		}

		r.WasHandled = true
		return c.Update(&r)
	})
}

func (p *Persister) HandleLoginRequest(ctx context.Context, challenge string, r *consent.HandledLoginRequest) (lr *consent.LoginRequest, err error) {
	return lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := c.Create(r)
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
		var ar consent.LoginRequest
		if err := c.Where("verifier = ?", verifier).Select("challenge", "client_id").First(&ar); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := c.Find(&d, ar.ID); err != nil {
			return sqlcon.HandleError(err)
		}

		if d.WasHandled {
			return errorsx.WithStack(fosite.ErrInvalidRequest.WithDebug("Login verifier has been used already."))
		}

		d.WasHandled = true
		return sqlcon.HandleError(c.Update(&d))
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
		var cr consent.HandledConsentRequest
		tn := consent.HandledConsentRequest{}.TableName()
		if err := c.
			Where(fmt.Sprintf("r.subject = ? AND r.client_id = ? AND r.skip=FALSE AND (%s.error='{}' AND %s.remember=TRUE)", tn, tn), subject, client).
			Join("hydra_oauth2_consent_request AS r", fmt.Sprintf("%s.challenge = r.challenge", tn)).
			Order(fmt.Sprintf("%s.requested_at DESC", tn)).
			Limit(1).
			First(&cr); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorsx.WithStack(consent.ErrNoPreviousConsentFound)
			}
			return sqlcon.HandleError(err)
		}

		var err error
		rs, err = p.resolveHandledConsentRequests(ctx, []consent.HandledConsentRequest{cr})
		return err
	})
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]consent.HandledConsentRequest, error) {
	var rs []consent.HandledConsentRequest
	c := p.Connection(ctx)
	tn := consent.HandledConsentRequest{}.TableName()

	if err := c.
		Where(fmt.Sprintf("r.subject = ? AND r.skip=FALSE AND %s.error='{}'", tn), subject).
		Join("hydra_oauth2_consent_request AS r", fmt.Sprintf("%s.challenge = r.challenge", tn)).
		Order(fmt.Sprintf("%s.requested_at DESC", tn)).
		Paginate(offset/limit+1, limit).
		All(&rs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.WithStack(consent.ErrNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return p.resolveHandledConsentRequests(ctx, rs)
}

func (p *Persister) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	tn := consent.HandledConsentRequest{}.TableName()

	n, err := p.Connection(ctx).
		Where(fmt.Sprintf("r.subject = ? AND r.skip=FALSE AND %s.error='{}'", tn), subject).
		Join("hydra_oauth2_consent_request as r", fmt.Sprintf("%s.challenge = r.challenge", tn)).
		Count(&consent.HandledConsentRequest{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) resolveHandledConsentRequests(ctx context.Context, requests []consent.HandledConsentRequest) ([]consent.HandledConsentRequest, error) {
	var result []consent.HandledConsentRequest

	for _, v := range requests {
		_, err := p.GetConsentRequest(ctx, v.ID)
		if errors.Is(err, sqlcon.ErrNoRows) || errors.Is(err, x.ErrNotFound) {
			return nil, errorsx.WithStack(consent.ErrNoPreviousConsentFound)
		} else if err != nil {
			return nil, err
		}

		// This will probably never error because we first check if the consent request actually exists
		if err := v.AfterFind(p.Connection(ctx)); err != nil {
			return nil, err
		}
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
			fmt.Sprintf(`SELECT DISTINCT c.* FROM hydra_client as c JOIN hydra_oauth2_consent_request as r ON (c.id = r.client_id) WHERE r.subject=? AND c.%schannel_logout_uri!='' AND c.%schannel_logout_uri IS NOT NULL AND r.login_session_id = ?`,
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
	var lr consent.LoginRequest
	var lrh consent.HandledLoginRequest

	var cr consent.ConsentRequest
	var crh consent.HandledConsentRequest

	// The value of notAfter should be the minimum between input parameter and request max expire based on its configured age
	requestMaxExpire := time.Now().Add(-p.config.ConsentRequestMaxAge())
	if requestMaxExpire.Before(notAfter) {
		notAfter = requestMaxExpire
	}

	challenges := []string{}
	queryFormat := `
	SELECT %[1]s.challenge
	FROM %[1]s
	LEFT JOIN %[2]s ON %[1]s.challenge = %[2]s.challenge
	WHERE (
		(%[2]s.challenge IS NULL)
		OR (%[2]s.error IS NOT NULL AND %[2]s.error <> '{}' AND %[2]s.error <> '')
	)
	AND %[1]s.requested_at < ?
	ORDER BY %[1]s.challenge
	LIMIT %[3]d
	`

	// Select challenges from all consent requests that can be safely deleted with limit
	// where hydra_oauth2_consent_request were unhandled or rejected, so either of these is true
	// - hydra_oauth2_authentication_request_handled does not exist (unhandled)
	// - hydra_oauth2_consent_request_handled has valid error (rejected)
	// AND timed-out
	// - hydra_oauth2_consent_request.requested_at < minimum between ttl.login_consent_request and notAfter
	q := p.Connection(ctx).RawQuery(fmt.Sprintf(queryFormat, (&cr).TableName(), (&crh).TableName(), limit), notAfter)

	if err := q.All(&challenges); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	}

	// Delete in batch consent requests and their references in cascade
	for i := 0; i < len(challenges); i += batchSize {
		j := i + batchSize
		if j > len(challenges) {
			j = len(challenges)
		}

		if i != j {
			q := p.Connection(ctx).RawQuery(
				fmt.Sprintf("DELETE FROM %s WHERE challenge in (?)", (&cr).TableName()),
				challenges[i:j],
			)

			if err := q.Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
		}
	}

	// Select challenges from all authentication requests that can be safely deleted with limit
	// where hydra_oauth2_authentication_request were unhandled or rejected, so either of these is true
	// - hydra_oauth2_authentication_request_handled does not exist (unhandled)
	// - hydra_oauth2_authentication_request_handled has valid error (rejected)
	// AND timed-out
	// - hydra_oauth2_authentication_request.requested_at < minimum between ttl.login_consent_request and notAfter
	q = p.Connection(ctx).RawQuery(fmt.Sprintf(queryFormat, (&lr).TableName(), (&lrh).TableName(), limit), notAfter)

	if err := q.All(&challenges); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	}

	// Delete in batch authentication requests
	for i := 0; i < len(challenges); i += batchSize {
		j := i + batchSize
		if j > len(challenges) {
			j = len(challenges)
		}

		if i != j {
			q := p.Connection(ctx).RawQuery(
				fmt.Sprintf("DELETE FROM %s WHERE challenge in (?)", (&lr).TableName()),
				challenges[i:j],
			)

			if err := q.Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
		}
	}

	return nil
}
