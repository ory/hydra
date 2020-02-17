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
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ory/x/dbal"
	"github.com/ory/x/stringsx"

	"github.com/ory/hydra/client"
)

var Migrations = map[string]*dbal.PackrMigrationSource{
	"mysql": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/mysql",
	}, true),
	"postgres": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/postgres",
	}, true),
	"cockroach": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/cockroach",
	}, true),
}

var sqlParamsAuthenticationRequestHandled = []string{
	"challenge",
	"subject",
	"remember",
	"remember_for",
	"error",
	"requested_at",
	"authenticated_at",
	"acr",
	"was_used",
	"context",
	"forced_subject_identifier",
}

var sqlParamsAuthenticationRequest = []string{
	"challenge",
	"verifier",
	"client_id",
	"subject",
	"request_url",
	"skip",
	"requested_scope",
	"requested_at_audience",
	"authenticated_at",
	"requested_at",
	"csrf",
	"oidc_context",
	"login_session_id",
}

var sqlParamsConsentRequest = append(sqlParamsAuthenticationRequest,
	"forced_subject_identifier",
	"login_challenge",
	"acr",
	"context",
)

var sqlParamsConsentRequestHandled = []string{
	"challenge",
	"granted_scope",
	"granted_at_audience",
	"remember",
	"remember_for",
	"authenticated_at",
	"error",
	"requested_at",
	"session_access_token",
	"session_id_token",
	"was_used",
	"handled_at",
}
var sqlParamsConsentRequestHandledUpdate = func() []string {
	p := make([]string, len(sqlParamsConsentRequestHandled))
	for i, v := range sqlParamsConsentRequestHandled {
		p[i] = fmt.Sprintf("%s=:%s", v, v)
	}
	return p
}()

var sqlParamsAuthSession = []string{
	"id",
	"authenticated_at",
	"subject",
	"remember",
}

var sqlParamsLogoutRequest = []string{
	"challenge",
	"verifier",
	"subject",
	"sid",
	"request_url",
	"redir_url",
	"was_used",
	"accepted",
	"rejected",
	"client_id",
	"rp_initiated",
}

type sqlLogoutRequest struct {
	Challenge             string         `db:"challenge"`
	Verifier              string         `db:"verifier"`
	Subject               string         `db:"subject"`
	SessionID             string         `db:"sid"`
	RequestURL            string         `db:"request_url"`
	PostLogoutRedirectURI string         `db:"redir_url"`
	WasUsed               bool           `db:"was_used"`
	Accepted              bool           `db:"accepted"`
	Rejected              bool           `db:"rejected"`
	Client                sql.NullString `db:"client_id"`
	RPInitiated           bool           `db:"rp_initiated"`
}

func newSQLLogoutRequest(c *LogoutRequest) *sqlLogoutRequest {
	var clientID sql.NullString
	if c.Client != nil {
		clientID = sql.NullString{
			Valid:  true,
			String: c.Client.ClientID,
		}
	}

	return &sqlLogoutRequest{
		Challenge:             c.Challenge,
		Verifier:              c.Verifier,
		Subject:               c.Subject,
		SessionID:             c.SessionID,
		RequestURL:            c.RequestURL,
		PostLogoutRedirectURI: c.PostLogoutRedirectURI,
		WasUsed:               c.WasUsed,
		Accepted:              c.Accepted,
		Client:                clientID,
		RPInitiated:           c.RPInitiated,
	}
}

func (r *sqlLogoutRequest) ToLogoutRequest(c *client.Client) *LogoutRequest {
	return &LogoutRequest{
		Challenge:             r.Challenge,
		Verifier:              r.Verifier,
		Subject:               r.Subject,
		SessionID:             r.SessionID,
		RequestURL:            r.RequestURL,
		PostLogoutRedirectURI: r.PostLogoutRedirectURI,
		WasUsed:               r.WasUsed,
		Accepted:              r.Accepted,
		Client:                c,
		RPInitiated:           r.RPInitiated,
	}
}

type sqlAuthenticationRequest struct {
	OpenIDConnectContext string         `db:"oidc_context"`
	Client               string         `db:"client_id"`
	Subject              string         `db:"subject"`
	RequestURL           string         `db:"request_url"`
	Skip                 bool           `db:"skip"`
	Challenge            string         `db:"challenge"`
	RequestedScope       string         `db:"requested_scope"`
	RequestedAudience    sql.NullString `db:"requested_at_audience"`
	Verifier             string         `db:"verifier"`
	CSRF                 string         `db:"csrf"`
	AuthenticatedAt      sql.NullTime   `db:"authenticated_at"`
	RequestedAt          time.Time      `db:"requested_at"`
	LoginSessionID       sql.NullString `db:"login_session_id"`
	Context              string         `db:"context"`
	WasHandled           bool           `db:"was_handled"`
}

type sqlConsentRequest struct {
	sqlAuthenticationRequest
	LoginChallenge          sql.NullString `db:"login_challenge"`
	ACR                     string         `db:"acr"`
	ForcedSubjectIdentifier string         `db:"forced_subject_identifier"`
}

func toMySQLDateHack(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func fromMySQLDateHack(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

func newSQLConsentRequest(c *ConsentRequest) (*sqlConsentRequest, error) {
	oidc, err := json.Marshal(c.OpenIDConnectContext)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var sessionID sql.NullString
	if len(c.LoginSessionID) > 0 {
		sessionID = sql.NullString{
			Valid:  true,
			String: c.LoginSessionID,
		}
	}

	if c.Context == nil {
		c.Context = map[string]interface{}{}
	}

	context, err := json.Marshal(c.Context)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &sqlConsentRequest{
		sqlAuthenticationRequest: sqlAuthenticationRequest{
			OpenIDConnectContext: string(oidc),
			Client:               c.Client.GetID(),
			Subject:              c.Subject,
			RequestURL:           c.RequestURL,
			Skip:                 c.Skip,
			Challenge:            c.Challenge,
			RequestedScope:       strings.Join(c.RequestedScope, "|"),
			RequestedAudience:    sql.NullString{Valid: true, String: strings.Join(c.RequestedAudience, "|")},
			Verifier:             c.Verifier,
			CSRF:                 c.CSRF,
			AuthenticatedAt:      toMySQLDateHack(c.AuthenticatedAt),
			RequestedAt:          c.RequestedAt,
			LoginSessionID:       sessionID,
			Context:              string(context),
		},
		LoginChallenge:          sql.NullString{Valid: true, String: c.LoginChallenge},
		ForcedSubjectIdentifier: c.ForceSubjectIdentifier,
		ACR:                     c.ACR,
	}, nil
}

func newSQLAuthenticationRequest(c *LoginRequest) (*sqlAuthenticationRequest, error) {
	oidc, err := json.Marshal(c.OpenIDConnectContext)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var sessionID sql.NullString
	if len(c.SessionID) > 0 {
		sessionID = sql.NullString{
			Valid:  true,
			String: c.SessionID,
		}
	}

	return &sqlAuthenticationRequest{
		OpenIDConnectContext: string(oidc),
		Client:               c.Client.GetID(),
		Subject:              c.Subject,
		RequestURL:           c.RequestURL,
		Skip:                 c.Skip,
		Challenge:            c.Challenge,
		RequestedScope:       strings.Join(c.RequestedScope, "|"),
		RequestedAudience:    sql.NullString{Valid: true, String: strings.Join(c.RequestedAudience, "|")},
		Verifier:             c.Verifier,
		CSRF:                 c.CSRF,
		AuthenticatedAt:      toMySQLDateHack(c.AuthenticatedAt),
		RequestedAt:          c.RequestedAt,
		LoginSessionID:       sessionID,
	}, nil
}

func (s *sqlAuthenticationRequest) toAuthenticationRequest(client *client.Client) (*LoginRequest, error) {
	var oidc OpenIDConnectContext
	if err := json.Unmarshal([]byte(s.OpenIDConnectContext), &oidc); err != nil {
		return nil, errors.WithStack(err)
	}

	return &LoginRequest{
		OpenIDConnectContext: &oidc,
		Client:               client,
		Subject:              s.Subject,
		RequestURL:           s.RequestURL,
		Skip:                 s.Skip,
		Challenge:            s.Challenge,
		RequestedScope:       stringsx.Splitx(s.RequestedScope, "|"),
		RequestedAudience:    stringsx.Splitx(s.RequestedAudience.String, "|"),
		Verifier:             s.Verifier,
		CSRF:                 s.CSRF,
		AuthenticatedAt:      fromMySQLDateHack(s.AuthenticatedAt),
		RequestedAt:          s.RequestedAt,
		WasHandled:           s.WasHandled,
		SessionID:            s.LoginSessionID.String,
	}, nil
}

func (s *sqlConsentRequest) toConsentRequest(client *client.Client) (*ConsentRequest, error) {
	var oidc OpenIDConnectContext
	var context map[string]interface{}
	if err := json.Unmarshal([]byte(s.OpenIDConnectContext), &oidc); err != nil {
		return nil, errors.WithStack(err)
	}

	if s.Context == "" {
		s.Context = "{}"
	}

	if err := json.Unmarshal([]byte(s.Context), &context); err != nil {
		return nil, errors.WithStack(err)
	}

	return &ConsentRequest{
		OpenIDConnectContext:   &oidc,
		Client:                 client,
		Subject:                s.Subject,
		RequestURL:             s.RequestURL,
		Skip:                   s.Skip,
		Challenge:              s.Challenge,
		RequestedScope:         stringsx.Splitx(s.RequestedScope, "|"),
		RequestedAudience:      stringsx.Splitx(s.RequestedAudience.String, "|"),
		Verifier:               s.Verifier,
		CSRF:                   s.CSRF,
		AuthenticatedAt:        fromMySQLDateHack(s.AuthenticatedAt),
		ForceSubjectIdentifier: s.ForcedSubjectIdentifier,
		RequestedAt:            s.RequestedAt,
		WasHandled:             s.WasHandled,
		LoginSessionID:         s.LoginSessionID.String,
		LoginChallenge:         s.LoginChallenge.String,
		Context:                context,
		ACR:                    s.ACR,
	}, nil
}

type sqlHandledConsentRequest struct {
	GrantedScope       string         `db:"granted_scope"`
	GrantedAudience    sql.NullString `db:"granted_at_audience"`
	SessionIDToken     string         `db:"session_id_token"`
	SessionAccessToken string         `db:"session_access_token"`
	Remember           bool           `db:"remember"`
	RememberFor        int            `db:"remember_for"`
	Error              string         `db:"error"`
	Challenge          string         `db:"challenge"`
	RequestedAt        time.Time      `db:"requested_at"`
	WasUsed            bool           `db:"was_used"`
	AuthenticatedAt    sql.NullTime   `db:"authenticated_at"`
	HandledAt          sql.NullTime   `db:"handled_at"`
}

func newSQLHandledConsentRequest(c *HandledConsentRequest) (*sqlHandledConsentRequest, error) {
	sidt := "{}"
	sat := "{}"
	e := "{}"

	if c.Session != nil {
		if len(c.Session.IDToken) > 0 {
			if out, err := json.Marshal(c.Session.IDToken); err != nil {
				return nil, errors.WithStack(err)
			} else {
				sidt = string(out)
			}
		}

		if len(c.Session.AccessToken) > 0 {
			if out, err := json.Marshal(c.Session.AccessToken); err != nil {
				return nil, errors.WithStack(err)
			} else {
				sat = string(out)
			}
		}
	}

	if c.Error != nil {
		if out, err := json.Marshal(c.Error); err != nil {
			return nil, errors.WithStack(err)
		} else {
			e = string(out)
		}
	}

	return &sqlHandledConsentRequest{
		GrantedScope:       strings.Join(c.GrantedScope, "|"),
		GrantedAudience:    sql.NullString{Valid: true, String: strings.Join(c.GrantedAudience, "|")},
		SessionIDToken:     sidt,
		SessionAccessToken: sat,
		Remember:           c.Remember,
		RememberFor:        c.RememberFor,
		Error:              e,
		Challenge:          c.Challenge,
		RequestedAt:        c.RequestedAt,
		WasUsed:            c.WasUsed,
		AuthenticatedAt:    toMySQLDateHack(c.AuthenticatedAt),
		HandledAt:          sql.NullTime{Time: c.HandledAt, Valid: true},
	}, nil
}

func (s *sqlHandledConsentRequest) toHandledConsentRequest(r *ConsentRequest) (*HandledConsentRequest, error) {
	var idt map[string]interface{}
	var at map[string]interface{}
	var e *RequestDeniedError

	if err := json.Unmarshal([]byte(s.SessionIDToken), &idt); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := json.Unmarshal([]byte(s.SessionAccessToken), &at); err != nil {
		return nil, errors.WithStack(err)
	}

	if len(s.Error) > 0 && s.Error != "{}" {
		e = new(RequestDeniedError)
		if err := json.Unmarshal([]byte(s.Error), &e); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &HandledConsentRequest{
		GrantedScope:    stringsx.Splitx(s.GrantedScope, "|"),
		GrantedAudience: stringsx.Splitx(s.GrantedAudience.String, "|"),
		RememberFor:     s.RememberFor,
		Remember:        s.Remember,
		Challenge:       s.Challenge,
		RequestedAt:     s.RequestedAt,
		WasUsed:         s.WasUsed,
		Session: &ConsentRequestSessionData{
			IDToken:     idt,
			AccessToken: at,
		},
		Error:           e,
		ConsentRequest:  r,
		AuthenticatedAt: fromMySQLDateHack(s.AuthenticatedAt),
		HandledAt:       s.HandledAt.Time,
	}, nil
}

type sqlHandledLoginRequest struct {
	Remember               bool         `db:"remember"`
	RememberFor            int          `db:"remember_for"`
	ACR                    string       `db:"acr"`
	Subject                string       `db:"subject"`
	Error                  string       `db:"error"`
	Challenge              string       `db:"challenge"`
	RequestedAt            time.Time    `db:"requested_at"`
	WasUsed                bool         `db:"was_used"`
	AuthenticatedAt        sql.NullTime `db:"authenticated_at"`
	Context                string       `db:"context"`
	ForceSubjectIdentifier string       `db:"forced_subject_identifier"`
}

func newSQLHandledLoginRequest(c *HandledLoginRequest) (*sqlHandledLoginRequest, error) {
	e := "{}"

	if c.Error != nil {
		if out, err := json.Marshal(c.Error); err != nil {
			return nil, errors.WithStack(err)
		} else {
			e = string(out)
		}
	}

	ctx := "{}"
	if c.Context != nil {
		if out, err := json.Marshal(c.Context); err != nil {
			return nil, errors.WithStack(err)
		} else {
			ctx = string(out)
		}
	}

	return &sqlHandledLoginRequest{
		ACR:                    c.ACR,
		Subject:                c.Subject,
		Remember:               c.Remember,
		RememberFor:            c.RememberFor,
		Error:                  e,
		Challenge:              c.Challenge,
		Context:                ctx,
		RequestedAt:            c.RequestedAt,
		WasUsed:                c.WasUsed,
		AuthenticatedAt:        toMySQLDateHack(c.AuthenticatedAt),
		ForceSubjectIdentifier: c.ForceSubjectIdentifier,
	}, nil
}

func (s *sqlHandledLoginRequest) toHandledLoginRequest(a *LoginRequest) (*HandledLoginRequest, error) {
	var e *RequestDeniedError

	if len(s.Error) > 0 && s.Error != "{}" {
		e = new(RequestDeniedError)
		if err := json.Unmarshal([]byte(s.Error), &e); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var context map[string]interface{}
	if err := json.Unmarshal([]byte(s.Context), &context); err != nil {
		return nil, errors.WithStack(err)
	}

	return &HandledLoginRequest{
		ForceSubjectIdentifier: s.ForceSubjectIdentifier,
		RememberFor:            s.RememberFor,
		Remember:               s.Remember,
		Challenge:              s.Challenge,
		RequestedAt:            s.RequestedAt,
		WasUsed:                s.WasUsed,
		ACR:                    s.ACR,
		Error:                  e,
		LoginRequest:           a,
		Context:                context,
		Subject:                s.Subject,
		AuthenticatedAt:        fromMySQLDateHack(s.AuthenticatedAt),
	}, nil
}
