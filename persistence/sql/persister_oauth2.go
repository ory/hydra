package sql

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/ory/fosite/storage"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/ory/fosite"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/stringsx"
)

var _ oauth2.AssertionJWTReader = &Persister{}
var _ storage.Transactional = &Persister{}

type (
	tableName        string
	OAuth2RequestSQL struct {
		ID                string         `db:"signature"`
		Request           string         `db:"request_id"`
		ConsentChallenge  sql.NullString `db:"challenge_id"`
		RequestedAt       time.Time      `db:"requested_at"`
		Client            string         `db:"client_id"`
		Scopes            string         `db:"scope"`
		GrantedScope      string         `db:"granted_scope"`
		RequestedAudience string         `db:"requested_audience"`
		GrantedAudience   string         `db:"granted_audience"`
		Form              string         `db:"form_data"`
		Subject           string         `db:"subject"`
		Active            bool           `db:"active"`
		Session           []byte         `db:"session_data"`
		Table             tableName      `db:"-"`
	}
)

const (
	sqlTableOpenID  tableName = "oidc"
	sqlTableAccess  tableName = "access"
	sqlTableRefresh tableName = "refresh"
	sqlTableCode    tableName = "code"
	sqlTablePKCE    tableName = "pkce"
)

func (r OAuth2RequestSQL) TableName() string {
	return r.Table.TableName()
}

func (table tableName) TableName() string {
	return "hydra_oauth2_" + string(table)
}

func (p *Persister) sqlSchemaFromRequest(rawSignature string, r fosite.Requester, table tableName) (*OAuth2RequestSQL, error) {
	subject := ""
	if r.GetSession() == nil {
		p.l.Debugf("Got an empty session in sqlSchemaFromRequest")
	} else {
		subject = r.GetSession().GetSubject()
	}

	session, err := p.marshalSession(r.GetSession())
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	var challenge sql.NullString
	rr, ok := r.GetSession().(*oauth2.Session)
	if !ok && r.GetSession() != nil {
		return nil, errors.Errorf("Expected request to be of type *Session, but got: %T", r.GetSession())
	} else if ok {
		if len(rr.ConsentChallenge) > 0 {
			challenge = sql.NullString{Valid: true, String: rr.ConsentChallenge}
		}
	}

	return &OAuth2RequestSQL{
		Request:           r.GetID(),
		ConsentChallenge:  challenge,
		ID:                p.hashSignature(rawSignature, table),
		RequestedAt:       r.GetRequestedAt(),
		Client:            r.GetClient().GetID(),
		Scopes:            strings.Join(r.GetRequestedScopes(), "|"),
		GrantedScope:      strings.Join(r.GetGrantedScopes(), "|"),
		GrantedAudience:   strings.Join(r.GetGrantedAudience(), "|"),
		RequestedAudience: strings.Join(r.GetRequestedAudience(), "|"),
		Form:              r.GetRequestForm().Encode(),
		Session:           session,
		Subject:           subject,
		Active:            true,
		Table:             table,
	}, nil
}

func (p *Persister) marshalSession(session fosite.Session) ([]byte, error) {
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	if !p.config.EncryptSessionData() {
		return sessionBytes, nil
	}

	ciphertext, err := p.r.KeyCipher().Encrypt(sessionBytes)
	if err != nil {
		return nil, err
	}

	return []byte(ciphertext), nil
}

func (r *OAuth2RequestSQL) toRequest(ctx context.Context, session fosite.Session, p *Persister) (*fosite.Request, error) {
	sess := r.Session
	if !gjson.ValidBytes(sess) {
		var err error
		sess, err = p.r.KeyCipher().Decrypt(string(sess))
		if err != nil {
			return nil, errorsx.WithStack(err)
		}
	}

	if session != nil {
		if err := json.Unmarshal(sess, session); err != nil {
			return nil, errorsx.WithStack(err)
		}
	} else {
		p.l.Debugf("Got an empty session in toRequest")
	}

	c, err := p.GetClient(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	val, err := url.ParseQuery(r.Form)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return &fosite.Request{
		ID:                r.Request,
		RequestedAt:       r.RequestedAt,
		Client:            c,
		RequestedScope:    stringsx.Splitx(r.Scopes, "|"),
		GrantedScope:      stringsx.Splitx(r.GrantedScope, "|"),
		RequestedAudience: stringsx.Splitx(r.RequestedAudience, "|"),
		GrantedAudience:   stringsx.Splitx(r.GrantedAudience, "|"),
		Form:              val,
		Session:           session,
	}, nil
}

// hashSignature prevents errors where the signature is longer than 128 characters (and thus doesn't fit into the pk).
func (p *Persister) hashSignature(signature string, table tableName) string {
	if table == sqlTableAccess && p.config.IsUsingJWTAsAccessTokens() {
		return fmt.Sprintf("%x", sha512.Sum384([]byte(signature)))
	}
	return signature
}

func (p *Persister) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	j, err := p.GetClientAssertionJWT(ctx, jti)
	if errors.Is(err, sqlcon.ErrNoRows) {
		// the jti is not known => valid
		return nil
	} else if err != nil {
		return err
	}
	if j.Expiry.After(time.Now()) {
		// the jti is not expired yet => invalid
		return errorsx.WithStack(fosite.ErrJTIKnown)
	}
	// the jti is expired => valid
	return nil
}

func (p *Persister) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		// delete expired
		now := "now()"
		if c.Dialect.Name() == "sqlite3" {
			now = "CURRENT_TIMESTAMP"
		}
		/* #nosec G201 table is static */
		if err := c.RawQuery(fmt.Sprintf("DELETE FROM %s WHERE expires_at < %s", oauth2.BlacklistedJTI{}.TableName(), now)).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}

		if err := p.SetClientAssertionJWTRaw(ctx, oauth2.NewBlacklistedJTI(jti, exp)); errors.Is(err, sqlcon.ErrUniqueViolation) {
			// found a jti
			return errorsx.WithStack(fosite.ErrJTIKnown)
		} else if err != nil {
			return err
		}
		// setting worked without a problem
		return nil
	})
}

func (p *Persister) GetClientAssertionJWT(ctx context.Context, j string) (*oauth2.BlacklistedJTI, error) {
	jti := oauth2.NewBlacklistedJTI(j, time.Time{})
	return jti, sqlcon.HandleError(p.Connection(ctx).Find(jti, jti.ID))
}

func (p *Persister) SetClientAssertionJWTRaw(ctx context.Context, jti *oauth2.BlacklistedJTI) error {
	return sqlcon.HandleError(p.Connection(ctx).Create(jti))
}

func (p *Persister) createSession(ctx context.Context, signature string, requester fosite.Requester, table tableName) error {
	req, err := p.sqlSchemaFromRequest(signature, requester, table)
	if err != nil {
		return err
	}

	if err := sqlcon.HandleError(p.Connection(ctx).Create(req)); errors.Is(err, sqlcon.ErrConcurrentUpdate) {
		return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	} else if err != nil {
		return err
	}
	return nil
}

func (p *Persister) updateRefreshSession(ctx context.Context, requestId string, session fosite.Session, used bool) error {
	_, ok := session.(*oauth2.Session)
	if !ok && session != nil {
		return errors.Errorf("expected session to be of type *oauth2.Session but got: %T", session)
	}
	sessionBytes, err := p.marshalSession(session)
	if err != nil {
		return err
	}

	updateSql := fmt.Sprintf("UPDATE %s SET session_data = ?, used = ? WHERE request_id = ?",
		sqlTableRefresh.TableName())

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := p.Connection(ctx).RawQuery(updateSql, sessionBytes, used, requestId).Exec()
		if errors.Is(err, sql.ErrNoRows) {
			return errorsx.WithStack(fosite.ErrNotFound)
		} else if err != nil {
			return sqlcon.HandleError(err)
		}
		return nil
	})
}

func (p *Persister) findSessionBySignature(ctx context.Context, rawSignature string, session fosite.Session, table tableName) (fosite.Requester, error) {
	rawSignature = p.hashSignature(rawSignature, table)

	r := OAuth2RequestSQL{Table: table}
	var fr fosite.Requester

	return fr, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		err := p.Connection(ctx).Where("signature = ?", rawSignature).First(&r)
		if errors.Is(err, sql.ErrNoRows) {
			return errorsx.WithStack(fosite.ErrNotFound)
		} else if err != nil {
			return sqlcon.HandleError(err)
		} else if !r.Active {
			fr, err = r.toRequest(ctx, session, p)
			if err != nil {
				return err
			} else if table == sqlTableCode {
				return errorsx.WithStack(fosite.ErrInvalidatedAuthorizeCode)
			}

			return errorsx.WithStack(fosite.ErrInactiveToken)
		}

		fr, err = r.toRequest(ctx, session, p)
		return err
	})
}

func (p *Persister) deleteSessionBySignature(ctx context.Context, signature string, table tableName) error {
	signature = p.hashSignature(signature, table)

	/* #nosec G201 table is static */
	return sqlcon.HandleError(
		p.Connection(ctx).
			RawQuery(fmt.Sprintf("DELETE FROM %s WHERE signature=?", OAuth2RequestSQL{Table: table}.TableName()), signature).
			Exec())
}

func (p *Persister) deleteSessionByRequestID(ctx context.Context, id string, table tableName) error {
	/* #nosec G201 table is static */
	if err := p.Connection(ctx).RawQuery(
		fmt.Sprintf("DELETE FROM %s WHERE request_id=?", OAuth2RequestSQL{Table: table}.TableName()),
		id,
	).Exec(); errors.Is(err, sql.ErrNoRows) {
		return errorsx.WithStack(fosite.ErrNotFound)
	} else if err := sqlcon.HandleError(err); err != nil {
		if errors.Is(err, sqlcon.ErrConcurrentUpdate) {
			return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
		} else if strings.Contains(err.Error(), "Error 1213") { // InnoDB Deadlock?
			return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
		}
		return err
	}
	return nil
}

func (p *Persister) deactivateSessionByRequestID(ctx context.Context, id string, table tableName) error {
	/* #nosec G201 table is static */
	return sqlcon.HandleError(
		p.Connection(ctx).
			RawQuery(
				fmt.Sprintf("UPDATE %s SET active=false, used=true WHERE request_id=?", OAuth2RequestSQL{Table: table}.TableName()),
				id,
			).
			Exec(),
	)
}

func (p *Persister) getRefreshTokenUsedStatusBySignature(ctx context.Context, signature string) (bool, error) {
	var used bool
	return used, p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		query := fmt.Sprintf("SELECT used FROM %s WHERE signature = ?", sqlTableRefresh.TableName())
		err := p.Connection(ctx).RawQuery(query, signature).First(&used)
		if errors.Is(err, sql.ErrNoRows) {
			return errorsx.WithStack(fosite.ErrNotFound)
		} else if err != nil {
			return sqlcon.HandleError(err)
		}
		return err
	})
}

func (p *Persister) CreateAuthorizeCodeSession(ctx context.Context, signature string, requester fosite.Requester) (err error) {
	return p.createSession(ctx, signature, requester, sqlTableCode)
}

func (p *Persister) GetAuthorizeCodeSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return p.findSessionBySignature(ctx, signature, session, sqlTableCode)
}

func (p *Persister) InvalidateAuthorizeCodeSession(ctx context.Context, signature string) (err error) {
	/* #nosec G201 table is static */
	return sqlcon.HandleError(
		p.Connection(ctx).
			RawQuery(
				fmt.Sprintf("UPDATE %s SET active=false WHERE signature=?", OAuth2RequestSQL{Table: sqlTableCode}.TableName()),
				signature).
			Exec())
}

func (p *Persister) CreateAccessTokenSession(ctx context.Context, signature string, requester fosite.Requester) (err error) {
	return p.createSession(ctx, signature, requester, sqlTableAccess)
}

func (p *Persister) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return p.findSessionBySignature(ctx, signature, session, sqlTableAccess)
}

func (p *Persister) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	return p.deleteSessionBySignature(ctx, signature, sqlTableAccess)
}

func (p *Persister) CreateRefreshTokenSession(ctx context.Context, signature string, requester fosite.Requester) (err error) {
	return p.createSession(ctx, signature, requester, sqlTableRefresh)
}

func (p *Persister) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return p.findSessionBySignature(ctx, signature, session, sqlTableRefresh)
}

func (p *Persister) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return p.deleteSessionBySignature(ctx, signature, sqlTableRefresh)
}

func (p *Persister) CreateOpenIDConnectSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return p.createSession(ctx, signature, requester, sqlTableOpenID)
}

func (p *Persister) GetOpenIDConnectSession(ctx context.Context, signature string, requester fosite.Requester) (fosite.Requester, error) {
	return p.findSessionBySignature(ctx, signature, requester.GetSession(), sqlTableOpenID)
}

func (p *Persister) DeleteOpenIDConnectSession(ctx context.Context, signature string) error {
	return p.deleteSessionBySignature(ctx, signature, sqlTableOpenID)
}

func (p *Persister) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return p.findSessionBySignature(ctx, signature, session, sqlTablePKCE)
}

func (p *Persister) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return p.createSession(ctx, signature, requester, sqlTablePKCE)
}

func (p *Persister) DeletePKCERequestSession(ctx context.Context, signature string) error {
	return p.deleteSessionBySignature(ctx, signature, sqlTablePKCE)
}

func (p *Persister) RevokeRefreshToken(ctx context.Context, requestId string) error {
	return p.deactivateSessionByRequestID(ctx, requestId, sqlTableRefresh)
}

func (p *Persister) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestId string, signature string) error {
	gracePeriod := p.config.RefreshTokenRotationGracePeriod()
	if gracePeriod <= 0 {
		return p.RevokeRefreshToken(ctx, requestId)
	}

	var requester fosite.Requester
	var err error
	session := new(oauth2.Session)
	if requester, err = p.GetRefreshTokenSession(ctx, signature, session); err != nil {
		p.l.Errorf("signature: %s not found. grace period not applied", signature)
		return errors.WithStack(err)
	}

	var used bool
	if used,err = p.getRefreshTokenUsedStatusBySignature(ctx, signature); err != nil {
		p.l.Errorf("signature: %s used status not found. grace period not applied", signature)
		return errors.WithStack(err)
	}

	if ! used {
		session := requester.GetSession()
		session.SetExpiresAt(fosite.RefreshToken, time.Now().UTC().Add(gracePeriod))
		if err = p.updateRefreshSession(ctx, requestId, session, true); err != nil {
			p.l.Errorf("failed to update session with signature: %s", signature)
			return errors.WithStack(err)
		}
	} else {
		p.l.Debugf("request_id: %s has already been used and is in the grace period", requestId)
	}
	return nil
}

func (p *Persister) RevokeAccessToken(ctx context.Context, id string) error {
	return p.deleteSessionByRequestID(ctx, id, sqlTableAccess)
}

func (p *Persister) FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	/* #nosec G201 table is static */
	// The value of notAfter should be the minimum between input parameter and access token max expire based on its configured age
	requestMaxExpire := time.Now().Add(-p.config.AccessTokenLifespan())
	if requestMaxExpire.Before(notAfter) {
		notAfter = requestMaxExpire
	}

	signatures := []string{}

	// Select tokens' signatures with limit
	q := p.Connection(ctx).RawQuery(
		fmt.Sprintf("SELECT signature FROM %s WHERE requested_at < ? ORDER BY signature LIMIT %d",
			OAuth2RequestSQL{Table: sqlTableAccess}.TableName(), limit),
		notAfter,
	)
	if err := q.All(&signatures); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	}

	// Delete tokens in batch
	var err error
	for i := 0; i < len(signatures); i += batchSize {
		j := i + batchSize
		if j > len(signatures) {
			j = len(signatures)
		}

		if i != j {
			err = p.Connection(ctx).RawQuery(
				fmt.Sprintf("DELETE FROM %s WHERE signature in (?)", OAuth2RequestSQL{Table: sqlTableAccess}.TableName()),
				signatures[i:j],
			).Exec()
		}
	}
	return sqlcon.HandleError(err)
}

func (p *Persister) FlushInactiveRefreshTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	/* #nosec G201 table is static */
	// The value of notAfter should be the minimum between input parameter and refresh token max expire based on its configured age
	requestMaxExpire := time.Now().Add(-p.config.RefreshTokenLifespan())
	if requestMaxExpire.Before(notAfter) {
		notAfter = requestMaxExpire
	}

	signatures := []string{}

	// Select tokens' signatures with limit
	q := p.Connection(ctx).RawQuery(
		fmt.Sprintf("SELECT signature FROM %s WHERE requested_at < ? ORDER BY signature LIMIT %d",
			OAuth2RequestSQL{Table: sqlTableRefresh}.TableName(), limit),
		notAfter,
	)
	if err := q.All(&signatures); err == sql.ErrNoRows {
		return errors.Wrap(fosite.ErrNotFound, "")
	}

	// Delete tokens in batch
	var err error
	for i := 0; i < len(signatures); i += batchSize {
		j := i + batchSize
		if j > len(signatures) {
			j = len(signatures)
		}

		if i != j {
			err = p.Connection(ctx).RawQuery(
				fmt.Sprintf("DELETE FROM %s WHERE signature in (?)", OAuth2RequestSQL{Table: sqlTableRefresh}.TableName()),
				signatures[i:j],
			).Exec()
		}
	}
	return sqlcon.HandleError(err)
}

func (p *Persister) DeleteAccessTokens(ctx context.Context, clientID string) error {
	/* #nosec G201 table is static */
	return sqlcon.HandleError(
		p.Connection(ctx).
			RawQuery(fmt.Sprintf("DELETE FROM %s WHERE client_id=?", OAuth2RequestSQL{Table: sqlTableAccess}.TableName()), clientID).
			Exec())
}
