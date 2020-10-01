package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon"
)

var _ consent.Manager = &Persister{}

func (p *Persister) CreateConsentRequest(ctx context.Context, req *consent.ConsentRequest) error {
	return errors.WithStack(p.Connection(ctx).Create(req))
}

func (p *Persister) GetConsentRequest(ctx context.Context, challenge string) (*consent.ConsentRequest, error) {
	r := &consent.ConsentRequest{}
	return r, sqlcon.HandleError(r.FindInDB(p.Connection(ctx), challenge))
}

func (p *Persister) HandleConsentRequest(ctx context.Context, challenge string, r *consent.HandledConsentRequest) (*consent.ConsentRequest, error) {
	c := p.Connection(ctx)
	err := sqlcon.HandleError(c.Create(r))
	if err != nil {
		if errors.Is(err, sqlcon.ErrUniqueViolation) {
			hr := &consent.HandledConsentRequest{}
			if err := c.Find(hr, r.ID); err != nil {
				return nil, sqlcon.HandleError(err)
			}
			if hr.WasUsed {
				return nil, sqlcon.ErrNoRows
			}
			prevErrNil := false
			if r.Error != nil {
				prevErrNil = true
				func() {}()
			}
			if err := c.Update(r); err != nil {
				return nil, sqlcon.HandleError(err)
			}
			if r.Error != nil && !prevErrNil {
				prevErrNil = false
			}
		}
	}
	return p.GetConsentRequest(ctx, challenge)
}

func (p *Persister) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	return p.transaction(ctx, p.revokeConsentSession("subject = ?", user))
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	return p.transaction(ctx, p.revokeConsentSession("subject = ? AND client_id = ?", user, client))
}

func (p *Persister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
	return func(ctx context.Context, c *pop.Connection) error {
		hrs := make([]*consent.HandledConsentRequest, 0)
		if err := c.
			Where(whereStmt, whereArgs...).
			Select("r.challenge").
			Join(fmt.Sprintf("%s AS r", consent.ConsentRequest{}.TableName()), fmt.Sprintf("r.challenge = %s.challenge", consent.HandledConsentRequest{}.TableName())).
			All(&hrs); err != nil {
			return sqlcon.HandleError(err)
		}

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

			if err := c.RawQuery("DELETE FROM hydra_oauth2_consent_request_handled WHERE challenge = ?", hr.ID).Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
			if err := c.RawQuery("DELETE FROM hydra_oauth2_consent_request WHERE challenge = ?", hr.ID).Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
		}
		return nil
	}
}

func (p *Persister) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.HandledConsentRequest, error) {
	var r consent.HandledConsentRequest
	var cr consent.ConsentRequest
	return &r, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.Where("verifier = ?", verifier).Select("challenge", "client_id").First(&cr); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := c.Find(&r, cr.ID); err != nil {
			return sqlcon.HandleError(err)
		}

		if r.WasUsed {
			return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
		}

		r.WasUsed = true
		return c.Update(&r)
	})
}

func (p *Persister) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, subject string) ([]consent.HandledConsentRequest, error) {
	rs := make([]consent.HandledConsentRequest, 0, 1)
	tn := consent.HandledConsentRequest{}.TableName()

	return rs, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := c.
			Where("r.subject = ? AND r.client_id = ? AND r.skip=FALSE", subject, client).
			Where(fmt.Sprintf("%s.error='{}' AND %s.remember=TRUE", tn, tn)).
			Join("hydra_oauth2_consent_request AS r", fmt.Sprintf("%s.challenge = r.challenge", tn)).
			Order(fmt.Sprintf("%s.requested_at DESC", tn)).
			Limit(1).
			All(&rs)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.WithStack(consent.ErrNoPreviousConsentFound)
			}
			return sqlcon.HandleError(err)
		}

		rs, err = p.resolveHandledConsentRequests(ctx, rs)
		return err
	})
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]consent.HandledConsentRequest, error) {
	rs := make([]consent.HandledConsentRequest, 0)
	tn := consent.HandledConsentRequest{}.TableName()

	return rs, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := c.
			Where("r.subject = ? AND r.skip=FALSE", subject).
			Where(fmt.Sprintf("%s.error='{}'", tn)).
			Join("hydra_oauth2_consent_request AS r", fmt.Sprintf("%s.challenge = r.challenge", tn)).
			Order(fmt.Sprintf("%s.requested_at DESC", tn)).
			Paginate(offset/limit+1, limit).
			All(&rs)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.WithStack(consent.ErrNoPreviousConsentFound)
			}
			return sqlcon.HandleError(err)
		}

		rs, err = p.resolveHandledConsentRequests(ctx, rs)
		return err
	})
}

func (p *Persister) resolveHandledConsentRequests(ctx context.Context, requests []consent.HandledConsentRequest) ([]consent.HandledConsentRequest, error) {
	var result []consent.HandledConsentRequest
	return result, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		for _, v := range requests {
			_, err := p.GetConsentRequest(ctx, v.ID)
			if errors.Is(err, sqlcon.ErrNoRows) {
				return errors.WithStack(consent.ErrNoPreviousConsentFound)
			} else if err != nil {
				return err
			}

			// this will probably never error because we first check if the consent request actually exists
			if err := v.AfterFind(p.Connection(ctx)); err != nil {
				return err
			}
			if v.RememberFor > 0 && v.RequestedAt.Add(time.Duration(v.RememberFor)*time.Second).Before(time.Now().UTC()) {
				continue
			}

			result = append(result, v)
		}

		if len(result) == 0 {
			return errors.WithStack(consent.ErrNoPreviousConsentFound)
		}
		return nil
	})
}

func (p *Persister) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	tn := consent.HandledConsentRequest{}.TableName()

	n, err := p.Connection(ctx).
		Where("r.subject = ? AND r.skip=FALSE", subject).
		Where(fmt.Sprintf("%s.error='{}'", tn)).
		Join("hydra_oauth2_consent_request as r", fmt.Sprintf("%s.challenge = r.challenge", tn)).
		Count(&consent.HandledConsentRequest{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) GetRememberedLoginSession(ctx context.Context, id string) (*consent.LoginSession, error) {
	var s consent.LoginSession
	err := p.Connection(ctx).Where("remember = TRUE").Find(&s, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	return &s, sqlcon.HandleError(err)
}

func (p *Persister) CreateLoginSession(ctx context.Context, session *consent.LoginSession) error {
	return errors.WithStack(p.Connection(ctx).Create(session))
}

func (p *Persister) DeleteLoginSession(ctx context.Context, id string) error {
	return sqlcon.HandleError(
		p.Connection(ctx).Destroy(
			&consent.LoginSession{ID: id},
		))
}

func (p *Persister) RevokeSubjectLoginSession(ctx context.Context, subject string) error {
	err := p.Connection(ctx).RawQuery("DELETE FROM hydra_oauth2_authentication_session WHERE subject = ?", subject).Exec()
	if errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(x.ErrNotFound)
	}
	return sqlcon.HandleError(err)
}

func (p *Persister) ConfirmLoginSession(ctx context.Context, id string, subject string, remember bool) error {
	return sqlcon.HandleError(
		p.Connection(ctx).Update(&consent.LoginSession{
			ID:              id,
			AuthenticatedAt: time.Now().UTC(),
			Subject:         subject,
			Remember:        remember,
		}))
}

func (p *Persister) CreateLoginRequest(ctx context.Context, req *consent.LoginRequest) error {
	return errors.WithStack(p.Connection(ctx).Create(req))
}

func (p *Persister) GetLoginRequest(ctx context.Context, challenge string) (*consent.LoginRequest, error) {
	var lr consent.LoginRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := (&lr).FindInDB(c, challenge)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.WithStack(x.ErrNotFound)
			}
			return sqlcon.HandleError(err)
		}

		lr.Client, err = p.GetConcreteClient(ctx, lr.ClientID)
		return err
	})
}

func (p *Persister) HandleLoginRequest(ctx context.Context, challenge string, r *consent.HandledLoginRequest) (lr *consent.LoginRequest, err error) {
	err = p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := c.Create(r)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		lr, err = p.GetLoginRequest(ctx, challenge)
		return sqlcon.HandleError(err)
	})
	return
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

		if d.WasUsed {
			return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Authentication verifier has been used already"))
		}

		d.WasUsed = true
		return sqlcon.HandleError(c.Update(&d))
	})
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
		err := c.RawQuery(
			"INSERT INTO hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated) VALUES (?, ?, ?)",
			session.Subject,
			session.ClientID,
			session.SubjectObfuscated,
		).Exec()
		return sqlcon.HandleError(err)
	})
}

func (p *Persister) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*consent.ForcedObfuscatedLoginSession, error) {
	var s consent.ForcedObfuscatedLoginSession
	return &s, sqlcon.HandleError(
		p.Connection(ctx).
			Where(
				"client_id = ? AND subject_obfuscated = ?",
				client,
				obfuscated,
			).First(&s),
	)
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
		var ids []string
		/* #nosec G201 - channel can either be "front" or "back" */
		if err := c.RawQuery(
			fmt.Sprintf(`SELECT DISTINCT(c.id) FROM hydra_client as c JOIN hydra_oauth2_consent_request as r ON (c.id = r.client_id) WHERE r.subject=? AND c.%schannel_logout_uri!='' AND c.%schannel_logout_uri IS NOT NULL AND r.login_session_id = ?`,
				channel,
				channel,
			),
			subject,
			sid,
		).All(&ids); err != nil {
			return sqlcon.HandleError(err)
		}

		cs = make([]client.Client, len(ids))
		for k, id := range ids {
			c, err := p.GetConcreteClient(ctx, id)
			if err != nil {
				return err
			}
			cs[k] = *c
		}
		return nil
	})
}

func (p *Persister) CreateLogoutRequest(ctx context.Context, request *consent.LogoutRequest) error {
	return errors.WithStack(p.Connection(ctx).Create(request))
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	var lr consent.LogoutRequest
	return &lr, sqlcon.HandleError(p.Connection(ctx).Where("challenge = ? AND rejected = FALSE", challenge).First(&lr))
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	if err := p.Connection(ctx).RawQuery("UPDATE hydra_oauth2_logout_request SET accepted=true, rejected=false WHERE challenge=?", challenge).Exec(); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetLogoutRequest(ctx, challenge)
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) error {
	return errors.WithStack(
		p.Connection(ctx).
			RawQuery("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=?", challenge).
			Exec())
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*consent.LogoutRequest, error) {
	c := p.Connection(ctx)

	var lr consent.LogoutRequest
	if err := c.Where("verifier=? AND was_used=FALSE AND accepted=TRUE AND rejected=FALSE", verifier).Select("challenge").First(&lr); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if err := c.RawQuery("UPDATE hydra_oauth2_logout_request SET was_used=TRUE WHERE verifier=?", verifier).Exec(); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.GetLogoutRequest(ctx, lr.ID)
}
