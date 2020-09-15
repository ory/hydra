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
	var r consent.ConsentRequest
	return &r, p.transaction(ctx, func(ctx context.Context, connection *pop.Connection) error {
		return p.Connection(ctx).Where("challenge = ?", challenge).LeftJoin("hydra_oauth2_consent_request_handled hr", "hr.challenge = challenge").First(&r)
	})
}

func (p *Persister) HandleConsentRequest(ctx context.Context, challenge string, r *consent.HandledConsentRequest) (cr *consent.ConsentRequest, err error) {
	err = p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := c.Create(r)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		cr, err = p.GetConsentRequest(ctx, challenge)
		return err
	})
	return
}

func (p *Persister) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	return p.transaction(ctx, p.revokeConsentSession("subject = ?", user))
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	return p.transaction(ctx, p.revokeConsentSession("subject = ? AND client_id = ?", user, client))
}

func (p *Persister) revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error {
	return func(ctx context.Context, c *pop.Connection) error {
		challenges := popableStringSlice{
			values: make([]string, 0),
			from:   "",
		}
		if err := c.
			Where(whereStmt, whereArgs...).
			Select("h.challenge AS values").
			Join(fmt.Sprintf("%s AS r", consent.ConsentRequest{}.TableName()), "r.challenge = h.challenge").
			All(&pop.Model{
				Value: &challenges,
				As:    "h",
			}); err != nil {
			return sqlcon.HandleError(err)
		}

		for _, challenge := range challenges.values {
			if err := p.RevokeAccessToken(ctx, challenge); errors.Cause(err) == fosite.ErrNotFound {
				// do nothing
			} else if err != nil {
				return err
			}
			if err := p.RevokeRefreshToken(ctx, challenge); errors.Cause(err) == fosite.ErrNotFound {
				// do nothing
			} else if err != nil {
				return err
			}

			if err := c.RawQuery("DELETE FROM hydra_oauth2_consent_request_handled WHERE challenge = ?", challenge).Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
			if err := c.RawQuery("DELETE FROM hydra_oauth2_consent_request WHERE challenge = ?", challenge).Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
		}
		return nil
	}
}

func (p *Persister) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.HandledConsentRequest, error) {
	var r consent.HandledConsentRequest
	return &r, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		var cr consent.ConsentRequest
		if err := c.Where("verifier = ?", verifier).First(&cr); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := c.Find(&r, cr.ID); err != nil {
			return sqlcon.HandleError(err)
		}

		if r.WasUsed {
			return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
		}

		r.WasUsed = true
		return sqlcon.HandleError(c.Update(&r))
	})
}

func (p *Persister) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, subject string) ([]consent.HandledConsentRequest, error) {
	rs := make([]consent.HandledConsentRequest, 0)

	if err := p.Connection(ctx).
		Where("r.subject = ? AND r.client_id = ? AND r.skip=FALSE", subject, client).
		Where("h.error='{}' AND h.remember=TRUE").
		Join("hydra_oauth2_consent_request AS r", "h.challenge = r.challenge").
		Order("h.requested_at DESC").
		Limit(1).
		All(&pop.Model{
			Value: &rs,
			As:    "h",
		}); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(consent.ErrNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return p.resolveHandledConsentRequests(rs)
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]consent.HandledConsentRequest, error) {
	rs := make([]consent.HandledConsentRequest, 0)

	if err := p.Connection(ctx).
		Where("r.subject = ? AND r.skip=FALSE", subject).
		Where("h.error='{}'").
		Join("hydra_oauth2_consent_request AS r", "h.challenge = r.challenge").
		Order("h.requested_at DESC").
		Paginate(offset/limit+1, limit).
		All(&pop.Model{
			Value: &rs,
			As:    "h",
		}); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(consent.ErrNoPreviousConsentFound)
		}
		return nil, sqlcon.HandleError(err)
	}

	return p.resolveHandledConsentRequests(rs)
}

func (p *Persister) resolveHandledConsentRequests(requests []consent.HandledConsentRequest) ([]consent.HandledConsentRequest, error) {
	var result []consent.HandledConsentRequest
	for _, v := range requests {
		if v.RememberFor > 0 && v.RequestedAt.Add(time.Duration(v.RememberFor)*time.Second).Before(time.Now().UTC()) {
			continue
		}

		result = append(result, v)
	}

	if len(result) == 0 {
		return nil, errors.WithStack(consent.ErrNoPreviousConsentFound)
	}

	return result, nil
}

func (p *Persister) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	n, err := p.Connection(ctx).
		Where("r.subject = ? AND r.skip=FALSE", subject).
		Where("h.error='{}'").
		Join("hydra_oauth2_consent_request as r", "h.challenge = r.challenge").
		Count(&pop.Model{
			Value: &consent.HandledConsentRequest{},
			As:    "h",
		})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) GetRememberedLoginSession(ctx context.Context, id string) (*consent.LoginSession, error) {
	var s consent.LoginSession
	err := p.Connection(ctx).Find(&s, id)
	if err == sql.ErrNoRows {
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
	if err == sql.ErrNoRows {
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
		err := c.Select("r.*", "COALESCE(hr.was_used, false) as was_handled").
			LeftJoin("hydra_oauth2_authentication_request_handled as hr", "r.challenge = hr.challenge").
			Find(&lr, challenge)
		if err != nil {
			if err == sql.ErrNoRows {
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
		if err := c.Where("verifier = ?", verifier).Select("challenge").First(&ar); err != nil {
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
	return &s, errors.WithStack(
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
	c := p.Connection(ctx)
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
		return nil, sqlcon.HandleError(err)
	}

	cs := make([]client.Client, len(ids))
	for k, id := range ids {
		c, err := p.GetConcreteClient(ctx, id)
		if err != nil {
			return nil, err
		}
		cs[k] = *c
	}

	return cs, nil
}

func (p *Persister) CreateLogoutRequest(ctx context.Context, request *consent.LogoutRequest) error {
	return errors.WithStack(p.Connection(ctx).Create(request))
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	var lr consent.LogoutRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.Where("challenge = ? AND rejected = FALSE", challenge).First(&lr); err != nil {
			return sqlcon.HandleError(err)
		}

		if lr.ClientID.Valid {
			var err error
			lr.Client, err = p.GetConcreteClient(ctx, lr.ClientID.String)
			return err
		}
		return nil
	})
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	var lr consent.LogoutRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := p.Connection(ctx).RawQuery("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=?", challenge).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}

		r, err := p.GetLogoutRequest(ctx, lr.Challenge)
		if err != nil {
			return err
		}
		lr = *r
		return nil
	})
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) error {
	return errors.WithStack(
		p.Connection(ctx).
			RawQuery("UPDATE hydra_oauth2_logout_request SET rejected=true, accepted=false WHERE challenge=?", challenge).
			Exec())
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*consent.LogoutRequest, error) {
	var lr consent.LogoutRequest
	return &lr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := c.Where("verifier=? AND was_used=FALSE AND accepted=TRUE AND rejected=FALSE", verifier).Select("challenge").First(&lr); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := c.RawQuery("UPDATE hydra_oauth2_logout_request SET was_used=TRUE WHERE verifier=?", verifier).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}

		r, err := p.GetLogoutRequest(ctx, lr.Challenge)
		if err != nil {
			return err
		}
		lr = *r
		return nil
	})
}
