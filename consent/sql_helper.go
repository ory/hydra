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
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ory/go-convenience/stringsx"
	"github.com/ory/hydra/client"
	"github.com/ory/x/dbal"
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

var sqlParamsConsentRequest = append(sqlParamsAuthenticationRequest, "forced_subject_identifier", "login_challenge", "acr")

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
}

var sqlParamsAuthSession = []string{
	"id",
	"authenticated_at",
	"subject",
}

type sqlAuthenticationRequest struct {
	OpenIDConnectContext string         `db:"oidc_context"`
	Client               string         `db:"client_id"`
	Subject              string         `db:"subject"`
	RequestURL           string         `db:"request_url"`
	Skip                 bool           `db:"skip"`
	Challenge            string         `db:"challenge"`
	RequestedScope       string         `db:"requested_scope"`
	RequestedAudience    string         `db:"requested_at_audience"`
	Verifier             string         `db:"verifier"`
	CSRF                 string         `db:"csrf"`
	AuthenticatedAt      *time.Time     `db:"authenticated_at"`
	RequestedAt          time.Time      `db:"requested_at"`
	LoginSessionID       sql.NullString `db:"login_session_id"`
	WasHandled           bool           `db:"was_handled"`
}

type sqlConsentRequest struct {
	sqlAuthenticationRequest
	LoginChallenge          string `db:"login_challenge"`
	ACR                     string `db:"acr"`
	ForcedSubjectIdentifier string `db:"forced_subject_identifier"`
}

func toMySQLDateHack(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func fromMySQLDateHack(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
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

	return &sqlConsentRequest{
		sqlAuthenticationRequest: sqlAuthenticationRequest{
			OpenIDConnectContext: string(oidc),
			Client:               c.Client.GetID(),
			Subject:              c.Subject,
			RequestURL:           c.RequestURL,
			Skip:                 c.Skip,
			Challenge:            c.Challenge,
			RequestedScope:       strings.Join(c.RequestedScope, "|"),
			RequestedAudience:    strings.Join(c.RequestedAudience, "|"),
			Verifier:             c.Verifier,
			CSRF:                 c.CSRF,
			AuthenticatedAt:      toMySQLDateHack(c.AuthenticatedAt),
			RequestedAt:          c.RequestedAt,
			LoginSessionID:       sessionID,
		},
		LoginChallenge:          c.LoginChallenge,
		ForcedSubjectIdentifier: c.ForceSubjectIdentifier,
		ACR:                     c.ACR,
	}, nil
}

func newSQLAuthenticationRequest(c *AuthenticationRequest) (*sqlAuthenticationRequest, error) {
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
		RequestedAudience:    strings.Join(c.RequestedAudience, "|"),
		Verifier:             c.Verifier,
		CSRF:                 c.CSRF,
		AuthenticatedAt:      toMySQLDateHack(c.AuthenticatedAt),
		RequestedAt:          c.RequestedAt,
		LoginSessionID:       sessionID,
	}, nil
}

func (s *sqlAuthenticationRequest) toAuthenticationRequest(client *client.Client) (*AuthenticationRequest, error) {
	var oidc OpenIDConnectContext
	if err := json.Unmarshal([]byte(s.OpenIDConnectContext), &oidc); err != nil {
		return nil, errors.WithStack(err)
	}

	return &AuthenticationRequest{
		OpenIDConnectContext: &oidc,
		Client:               client,
		Subject:              s.Subject,
		RequestURL:           s.RequestURL,
		Skip:                 s.Skip,
		Challenge:            s.Challenge,
		RequestedScope:       stringsx.Splitx(s.RequestedScope, "|"),
		RequestedAudience:    stringsx.Splitx(s.RequestedAudience, "|"),
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
	if err := json.Unmarshal([]byte(s.OpenIDConnectContext), &oidc); err != nil {
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
		RequestedAudience:      stringsx.Splitx(s.RequestedAudience, "|"),
		Verifier:               s.Verifier,
		CSRF:                   s.CSRF,
		AuthenticatedAt:        fromMySQLDateHack(s.AuthenticatedAt),
		ForceSubjectIdentifier: s.ForcedSubjectIdentifier,
		RequestedAt:            s.RequestedAt,
		WasHandled:             s.WasHandled,
		LoginSessionID:         s.LoginSessionID.String,
		LoginChallenge:         s.LoginChallenge,
		ACR:                    s.ACR,
	}, nil
}

type sqlHandledConsentRequest struct {
	GrantedScope       string     `db:"granted_scope"`
	GrantedAudience    string     `db:"granted_at_audience"`
	SessionIDToken     string     `db:"session_id_token"`
	SessionAccessToken string     `db:"session_access_token"`
	Remember           bool       `db:"remember"`
	RememberFor        int        `db:"remember_for"`
	Error              string     `db:"error"`
	Challenge          string     `db:"challenge"`
	RequestedAt        time.Time  `db:"requested_at"`
	WasUsed            bool       `db:"was_used"`
	AuthenticatedAt    *time.Time `db:"authenticated_at"`
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
		GrantedAudience:    strings.Join(c.GrantedAudience, "|"),
		SessionIDToken:     sidt,
		SessionAccessToken: sat,
		Remember:           c.Remember,
		RememberFor:        c.RememberFor,
		Error:              e,
		Challenge:          c.Challenge,
		RequestedAt:        c.RequestedAt,
		WasUsed:            c.WasUsed,
		AuthenticatedAt:    toMySQLDateHack(c.AuthenticatedAt),
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
		GrantedAudience: stringsx.Splitx(s.GrantedAudience, "|"),
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
	}, nil
}

type sqlHandledAuthenticationRequest struct {
	Remember               bool       `db:"remember"`
	RememberFor            int        `db:"remember_for"`
	ACR                    string     `db:"acr"`
	Subject                string     `db:"subject"`
	Error                  string     `db:"error"`
	Challenge              string     `db:"challenge"`
	RequestedAt            time.Time  `db:"requested_at"`
	WasUsed                bool       `db:"was_used"`
	AuthenticatedAt        *time.Time `db:"authenticated_at"`
	ForceSubjectIdentifier string     `db:"forced_subject_identifier"`
}

func newSQLHandledAuthenticationRequest(c *HandledAuthenticationRequest) (*sqlHandledAuthenticationRequest, error) {
	e := "{}"

	if c.Error != nil {
		if out, err := json.Marshal(c.Error); err != nil {
			return nil, errors.WithStack(err)
		} else {
			e = string(out)
		}
	}

	return &sqlHandledAuthenticationRequest{
		ACR:                    c.ACR,
		Subject:                c.Subject,
		Remember:               c.Remember,
		RememberFor:            c.RememberFor,
		Error:                  e,
		Challenge:              c.Challenge,
		RequestedAt:            c.RequestedAt,
		WasUsed:                c.WasUsed,
		AuthenticatedAt:        toMySQLDateHack(c.AuthenticatedAt),
		ForceSubjectIdentifier: c.ForceSubjectIdentifier,
	}, nil
}

func (s *sqlHandledAuthenticationRequest) toHandledAuthenticationRequest(a *AuthenticationRequest) (*HandledAuthenticationRequest, error) {
	var e *RequestDeniedError

	if len(s.Error) > 0 && s.Error != "{}" {
		e = new(RequestDeniedError)
		if err := json.Unmarshal([]byte(s.Error), &e); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &HandledAuthenticationRequest{
		ForceSubjectIdentifier: s.ForceSubjectIdentifier,
		RememberFor:            s.RememberFor,
		Remember:               s.Remember,
		Challenge:              s.Challenge,
		RequestedAt:            s.RequestedAt,
		WasUsed:                s.WasUsed,
		ACR:                    s.ACR,
		Error:                  e,
		AuthenticationRequest:  a,
		Subject:                s.Subject,
		AuthenticatedAt:        fromMySQLDateHack(s.AuthenticatedAt),
	}, nil
}
