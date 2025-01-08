// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/stringsx"
)

const (
	sqlTableDeviceAuthCodes tableName = "hydra_oauth2_device_auth_codes"
)

type DeviceRequestSQL struct {
	ID                string               `db:"device_code_signature"`
	UserCodeID        string               `db:"user_code_signature"`
	NID               uuid.UUID            `db:"nid"`
	Request           string               `db:"request_id"`
	ConsentChallenge  sql.NullString       `db:"challenge_id"`
	RequestedAt       time.Time            `db:"requested_at"`
	Client            string               `db:"client_id"`
	Scopes            string               `db:"scope"`
	GrantedScope      string               `db:"granted_scope"`
	RequestedAudience string               `db:"requested_audience"`
	GrantedAudience   string               `db:"granted_audience"`
	Form              string               `db:"form_data"`
	Subject           string               `db:"subject"`
	DeviceCodeActive  bool                 `db:"device_code_active"`
	UserCodeState     fosite.UserCodeState `db:"user_code_state"`
	Session           []byte               `db:"session_data"`
	// InternalExpiresAt denormalizes the expiry from the session to additionally store it as a row.
	InternalExpiresAt sqlxx.NullTime `db:"expires_at" json:"-"`
}

func (r DeviceRequestSQL) TableName() string {
	return string(sqlTableDeviceAuthCodes)
}

func (r *DeviceRequestSQL) toRequest(ctx context.Context, session fosite.Session, p *Persister) (_ *fosite.DeviceRequest, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeviceRequestSQL.toRequest")
	defer otelx.End(span, &err)

	sess := r.Session
	if !gjson.ValidBytes(sess) {
		var err error
		sess, err = p.r.KeyCipher().Decrypt(ctx, string(sess), nil)
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

	return &fosite.DeviceRequest{
		UserCodeState: fosite.UserCodeState(r.UserCodeState),
		Request: fosite.Request{
			ID:          r.Request,
			RequestedAt: r.RequestedAt,
			// ExpiresAt does not need to be populated as we get the expiry time from the session.
			Client:            c,
			RequestedScope:    stringsx.Splitx(r.Scopes, "|"),
			GrantedScope:      stringsx.Splitx(r.GrantedScope, "|"),
			RequestedAudience: stringsx.Splitx(r.RequestedAudience, "|"),
			GrantedAudience:   stringsx.Splitx(r.GrantedAudience, "|"),
			Form:              val,
			Session:           session,
		},
	}, nil
}

func (p *Persister) sqlDeviceSchemaFromRequest(ctx context.Context, deviceCodeSignature, userCodeSignature string, r fosite.DeviceRequester, expiresAt time.Time) (*DeviceRequestSQL, error) {
	subject := ""
	if r.GetSession() == nil {
		p.l.Debugf("Got an empty session in sqlSchemaFromRequest")
	} else {
		subject = r.GetSession().GetSubject()
	}

	session, err := json.Marshal(r.GetSession())
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	if p.config.EncryptSessionData(ctx) {
		ciphertext, err := p.r.KeyCipher().Encrypt(ctx, session, nil)
		if err != nil {
			return nil, errorsx.WithStack(err)
		}
		session = []byte(ciphertext)
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

	return &DeviceRequestSQL{
		Request:           r.GetID(),
		ConsentChallenge:  challenge,
		ID:                deviceCodeSignature,
		UserCodeID:        userCodeSignature,
		RequestedAt:       r.GetRequestedAt(),
		InternalExpiresAt: sqlxx.NullTime(expiresAt),
		Client:            r.GetClient().GetID(),
		Scopes:            strings.Join(r.GetRequestedScopes(), "|"),
		GrantedScope:      strings.Join(r.GetGrantedScopes(), "|"),
		GrantedAudience:   strings.Join(r.GetGrantedAudience(), "|"),
		RequestedAudience: strings.Join(r.GetRequestedAudience(), "|"),
		Form:              r.GetRequestForm().Encode(),
		Session:           session,
		Subject:           subject,
		DeviceCodeActive:  true,
		UserCodeState:     r.GetUserCodeState(),
	}, nil
}

func (p *Persister) createDeviceAuthSession(ctx context.Context, deviceCodeSignature, userCodeSignature string, requester fosite.DeviceRequester, expiresAt time.Time) error {
	req, err := p.sqlDeviceSchemaFromRequest(ctx, deviceCodeSignature, userCodeSignature, requester, expiresAt)
	if err != nil {
		return err
	}

	if err = sqlcon.HandleError(p.CreateWithNetwork(ctx, req)); errors.Is(err, sqlcon.ErrConcurrentUpdate) {
		return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	} else if errors.Is(err, sqlcon.ErrUniqueViolation) {
		return errors.Wrap(fosite.ErrExistingUserCodeSignature, err.Error())
	} else if err != nil {
		return err
	}
	return nil
}

// CreateDeviceCodeSession creates a new device code session and stores it in the database
func (p *Persister) CreateDeviceAuthSession(ctx context.Context, deviceCodeSignature, userCodeSignature string, requester fosite.DeviceRequester) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateDeviceCodeSession")
	defer otelx.End(span, &err)
	return p.createDeviceAuthSession(ctx, deviceCodeSignature, userCodeSignature, requester, requester.GetSession().GetExpiresAt(fosite.DeviceCode).UTC())
}

// UpdateDeviceCodeSessionBySignature updates a device code session by the device_code signature
func (p *Persister) UpdateDeviceCodeSessionBySignature(ctx context.Context, signature string, requester fosite.DeviceRequester) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateDeviceCodeSessionBySignature")
	defer otelx.End(span, &err)

	req, err := p.sqlDeviceSchemaFromRequest(ctx, signature, "", requester, requester.GetSession().GetExpiresAt(fosite.DeviceCode).UTC())
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(
		"UPDATE %s SET granted_scope=?, granted_audience=?, session_data=?, user_code_state=? WHERE device_code_signature=? AND nid = ?",
		sqlTableDeviceAuthCodes,
	)

	/* #nosec G201 table is static */
	err = p.Connection(ctx).RawQuery(stmt, req.GrantedScope, req.GrantedAudience, req.Session, req.UserCodeState, signature, p.NetworkID(ctx)).Exec()
	if err != nil {
		return sqlcon.HandleError(err)
	}

	return nil
}

// GetDeviceCodeSession returns a device code session from the database
func (p *Persister) GetDeviceCodeSession(ctx context.Context, signature string, session fosite.Session) (_ fosite.DeviceRequester, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetDeviceCodeSession")
	defer otelx.End(span, &err)

	r := DeviceRequestSQL{}
	err = p.QueryWithNetwork(ctx).Where("device_code_signature = ?", signature).First(&r)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errorsx.WithStack(fosite.ErrNotFound)
	}
	if err != nil {
		return nil, sqlcon.HandleError(err)
	}
	if !r.DeviceCodeActive {
		fr, err := r.toRequest(ctx, session, p)
		if err != nil {
			return nil, err
		}
		return fr, errorsx.WithStack(fosite.ErrInactiveToken)
	}

	return r.toRequest(ctx, session, p)
}

// GetDeviceCodeSessionByRequestID returns a device code session from the database
func (p *Persister) GetDeviceCodeSessionByRequestID(ctx context.Context, requestID string, session fosite.Session) (_ fosite.DeviceRequester, deviceCodeSignature string, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetDeviceCodeSessionByRequestID")
	defer otelx.End(span, &err)

	r := DeviceRequestSQL{}
	err = p.QueryWithNetwork(ctx).Where("request_id = ?", requestID).First(&r)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", errorsx.WithStack(fosite.ErrNotFound)
	}
	if err != nil {
		return nil, "", sqlcon.HandleError(err)
	}
	if !r.DeviceCodeActive {
		fr, err := r.toRequest(ctx, session, p)
		if err != nil {
			return nil, "", err
		}
		return fr, r.ID, errorsx.WithStack(fosite.ErrInactiveToken)
	}

	fr, err := r.toRequest(ctx, session, p)
	if err != nil {
		return nil, "", err
	}
	return fr, r.ID, nil
}

// InvalidateDeviceCodeSession invalidates a device code session
func (p *Persister) InvalidateDeviceCodeSession(ctx context.Context, signature string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.InvalidateDeviceCodeSession")
	defer otelx.End(span, &err)

	/* #nosec G201 table is static */
	return sqlcon.HandleError(
		p.Connection(ctx).
			RawQuery(
				fmt.Sprintf("UPDATE %s SET device_code_active=false WHERE device_code_signature=? AND nid = ?", sqlTableDeviceAuthCodes),
				signature,
				p.NetworkID(ctx),
			).
			Exec(),
	)
}

// GetUserCodeSession returns a user code session from the database
func (p *Persister) GetUserCodeSession(ctx context.Context, signature string, session fosite.Session) (_ fosite.DeviceRequester, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetUserCodeSession")
	defer otelx.End(span, &err)

	r := DeviceRequestSQL{}
	if session == nil {
		session = oauth2.NewSession("")
	}
	err = p.QueryWithNetwork(ctx).Where("user_code_signature = ?", signature).First(&r)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errorsx.WithStack(fosite.ErrNotFound)
	}
	if err != nil {
		return nil, sqlcon.HandleError(err)
	}

	fr, err := r.toRequest(ctx, session, p)
	if err != nil {
		return nil, err
	}
	if r.UserCodeState != fosite.UserCodeUnused {
		return fr, errorsx.WithStack(fosite.ErrInactiveToken)
	}

	return fr, err
}
