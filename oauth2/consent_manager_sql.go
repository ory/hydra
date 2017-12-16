// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	"sort"
)

var sqlConsentParams = []string{
	"id", "client_id", "expires_at", "redirect_url", "requested_scopes",
	"csrf", "granted_scopes", "access_token_extra", "id_token_extra",
	"consent", "deny_reason", "subject", "client", "oidc_context", "requested_at",
}

var consentMigrations = func(db string) *migrate.MemoryMigrationSource {
	return &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "1",
				Up: []string{`CREATE TABLE IF NOT EXISTS hydra_consent_request (
	id      			varchar(36) NOT NULL PRIMARY KEY,
	requested_scopes 	text NOT NULL,
	client_id 			text NOT NULL,
	expires_at 			timestamp NOT NULL,
	redirect_url 		text NOT NULL,
	csrf 				text NOT NULL,
	granted_scopes		text NOT NULL,
	access_token_extra	text NOT NULL,
	id_token_extra		text NOT NULL,
	consent				text NOT NULL,
	deny_reason			text NOT NULL,
	subject				text NOT NULL
)`},
				Down: []string{
					"DROP TABLE hydra_consent_request",
				},
			},
			{
				Id: "2",
				Up: func() []string {
					if db == "mysql" {
						return []string{
							"ALTER TABLE hydra_consent_request ADD client text",
							"ALTER TABLE hydra_consent_request ADD oidc_context text",
							"UPDATE hydra_consent_request SET client='{}'",
							"UPDATE hydra_consent_request SET oidc_context='{}'",
							"ALTER TABLE hydra_consent_request MODIFY client text NOT NULL",
							"ALTER TABLE hydra_consent_request MODIFY oidc_context text NOT NULL",
							"ALTER TABLE hydra_consent_request ADD requested_at timestamp NOT NULL DEFAULT '1990-1-1 00:00:00'",
							"ALTER TABLE hydra_consent_request MODIFY subject varchar(255)",
							"ALTER TABLE hydra_consent_request MODIFY client_id varchar(255)",
							"ALTER TABLE hydra_consent_request ADD INDEX hydra_consent_request_subject_idx (subject)",
							"ALTER TABLE hydra_consent_request ADD INDEX hydra_consent_request_client_id_idx (client_id)",
						}
					} else {
						return []string{
							"ALTER TABLE hydra_consent_request ADD client text",
							"ALTER TABLE hydra_consent_request ADD oidc_context text",
							"UPDATE hydra_consent_request SET client='{}'",
							"UPDATE hydra_consent_request SET oidc_context='{}'",
							"ALTER TABLE hydra_consent_request ALTER COLUMN client SET NOT NULL",
							"ALTER TABLE hydra_consent_request ALTER COLUMN oidc_context SET NOT NULL",
							"ALTER TABLE hydra_consent_request ADD requested_at timestamp NOT NULL DEFAULT '1990-1-1 00:00:00'",
							"ALTER TABLE hydra_consent_request ALTER COLUMN subject TYPE varchar(255)",
							"ALTER TABLE hydra_consent_request ALTER COLUMN client_id TYPE varchar(255)",
							"CREATE INDEX hydra_consent_request_subject_idx ON hydra_consent_request (subject)",
							"CREATE INDEX hydra_consent_request_client_id_idx ON hydra_consent_request (client_id)",
						}
					}
				}(),
				Down: []string{
					"ALTER TABLE hydra_consent_request DROP COLUMN client",
					"ALTER TABLE hydra_consent_request DROP COLUMN oidc_context",
					"ALTER TABLE hydra_consent_request DROP COLUMN requested_at",
				},
			},
		},
	}
}

type consentRequestSqlData struct {
	ID               string    `db:"id"`
	RequestedScopes  string    `db:"requested_scopes"`
	ClientID         string    `db:"client_id"`
	ExpiresAt        time.Time `db:"expires_at"`
	RedirectURL      string    `db:"redirect_url"`
	CSRF             string    `db:"csrf"`
	GrantedScopes    string    `db:"granted_scopes"`
	AccessTokenExtra string    `db:"access_token_extra"`
	IDTokenExtra     string    `db:"id_token_extra"`
	Consent          string    `db:"consent"`
	DenyReason       string    `db:"deny_reason"`
	Subject          string    `db:"subject"`
	Client           string    `db:"client"`
	OIDCContext      string    `db:"oidc_context"`
	RequestedAt      time.Time `db:"requested_at"`
}

func newConsentRequestSqlData(request *ConsentRequest) (*consentRequestSqlData, error) {
	for k, scope := range request.RequestedScopes {
		request.RequestedScopes[k] = strings.Replace(scope, " ", "", -1)
	}
	for k, scope := range request.GrantedScopes {
		request.GrantedScopes[k] = strings.Replace(scope, " ", "", -1)
	}

	atext := ""
	idtext := ""

	if request.AccessTokenExtra != nil {
		if out, err := json.Marshal(request.AccessTokenExtra); err != nil {
			return nil, errors.WithStack(err)
		} else {
			atext = string(out)
		}
	}

	if request.IDTokenExtra != nil {
		if out, err := json.Marshal(request.IDTokenExtra); err != nil {
			return nil, errors.WithStack(err)
		} else {
			idtext = string(out)
		}
	}

	cl := "{}"
	oidcContext := "{}"
	if out, err := json.Marshal(request.Client); err != nil {
		return nil, errors.WithStack(err)
	} else {
		cl = string(out)
	}

	if out, err := json.Marshal(request.OpenIDConnectContext); err != nil {
		return nil, errors.WithStack(err)
	} else {
		oidcContext = string(out)
	}

	return &consentRequestSqlData{
		ID:               request.ID,
		RequestedScopes:  strings.Join(request.RequestedScopes, " "),
		GrantedScopes:    strings.Join(request.GrantedScopes, " "),
		ClientID:         request.ClientID,
		ExpiresAt:        request.ExpiresAt,
		RedirectURL:      request.RedirectURL,
		CSRF:             request.CSRF,
		AccessTokenExtra: atext,
		IDTokenExtra:     idtext,
		Consent:          request.Consent,
		DenyReason:       request.DenyReason,
		Subject:          request.Subject,
		Client:           cl,
		OIDCContext:      oidcContext,
		RequestedAt:      request.RequestedAt,
	}, nil
}

func (r *consentRequestSqlData) toConsentRequest() (*ConsentRequest, error) {
	var atext, idtext map[string]interface{}
	var cl client.Client
	var oidcContext ConsentRequestOpenIDConnectContext

	if r.IDTokenExtra != "" {
		if err := json.Unmarshal([]byte(r.IDTokenExtra), &idtext); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if r.AccessTokenExtra != "" {
		if err := json.Unmarshal([]byte(r.AccessTokenExtra), &atext); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if err := json.Unmarshal([]byte(r.Client), &cl); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal([]byte(r.OIDCContext), &oidcContext); err != nil {
		return nil, errors.WithStack(err)
	}

	return &ConsentRequest{
		ID:                   r.ID,
		ClientID:             r.ClientID,
		ExpiresAt:            r.ExpiresAt,
		RequestedAt:          r.RequestedAt,
		RedirectURL:          r.RedirectURL,
		CSRF:                 r.CSRF,
		Consent:              r.Consent,
		DenyReason:           r.DenyReason,
		RequestedScopes:      pkg.SplitNonEmpty(r.RequestedScopes, " "),
		GrantedScopes:        pkg.SplitNonEmpty(r.GrantedScopes, " "),
		AccessTokenExtra:     atext,
		IDTokenExtra:         idtext,
		Subject:              r.Subject,
		Client:               &cl,
		OpenIDConnectContext: &oidcContext,
	}, nil
}

type ConsentRequestSQLManager struct {
	db *sqlx.DB
}

func NewConsentRequestSQLManager(db *sqlx.DB) *ConsentRequestSQLManager {
	return &ConsentRequestSQLManager{db: db}
}

func (m *ConsentRequestSQLManager) CreateSchemas() (int, error) {
	migrate.SetTable("hydra_consent_request_migration")
	n, err := migrate.Exec(m.db.DB, m.db.DriverName(), consentMigrations(m.db.DriverName()), migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *ConsentRequestSQLManager) PersistConsentRequest(request *ConsentRequest) error {
	if request.ID == "" {
		request.ID = uuid.New()
	}

	data, err := newConsentRequestSqlData(request)
	if err != nil {
		return errors.WithStack(err)
	}

	query := fmt.Sprintf(
		"INSERT INTO hydra_consent_request (%s) VALUES (%s)",
		strings.Join(sqlConsentParams, ", "),
		":"+strings.Join(sqlConsentParams, ", :"),
	)
	if _, err := m.db.NamedExec(query, data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *ConsentRequestSQLManager) AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error {
	r, err := m.GetConsentRequest(id)
	if err != nil {
		return errors.WithStack(err)
	}

	r.Subject = payload.Subject
	r.AccessTokenExtra = payload.AccessTokenExtra
	r.IDTokenExtra = payload.IDTokenExtra
	r.Consent = ConsentRequestAccepted
	r.GrantedScopes = payload.GrantScopes

	return m.updateConsentRequest(r)
}

func (m *ConsentRequestSQLManager) RejectConsentRequest(id string, payload *RejectConsentRequestPayload) error {
	r, err := m.GetConsentRequest(id)
	if err != nil {
		return errors.WithStack(err)
	}

	r.Consent = ConsentRequestRejected
	r.DenyReason = payload.Reason

	return m.updateConsentRequest(r)
}

func (m *ConsentRequestSQLManager) updateConsentRequest(request *ConsentRequest) error {
	d, err := newConsentRequestSqlData(request)
	if err != nil {
		return errors.WithStack(err)
	}

	var query []string
	for _, param := range sqlConsentParams {
		query = append(query, fmt.Sprintf("%s=:%s", param, param))
	}

	if _, err := m.db.NamedExec(fmt.Sprintf(`UPDATE hydra_consent_request SET %s WHERE id=:id`, strings.Join(query, ", ")), d); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *ConsentRequestSQLManager) GetConsentRequest(id string) (*ConsentRequest, error) {
	var d consentRequestSqlData
	if err := m.db.Get(&d, m.db.Rebind("SELECT * FROM hydra_consent_request WHERE id=?"), id); err == sql.ErrNoRows {
		return nil, errors.WithStack(pkg.ErrNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	r, err := d.toConsentRequest()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r, nil
}

func (m *ConsentRequestSQLManager) GetPreviouslyGrantedConsent(subject string, client string, scopes []string) (*ConsentRequest, error) {
	var d []consentRequestSqlData
	if err := m.db.Select(&d, m.db.Rebind("SELECT * FROM hydra_consent_request WHERE subject=? AND client_id=? AND consent=?"), subject, client, ConsentRequestAccepted); err == sql.ErrNoRows {
		return nil, errors.WithStack(pkg.ErrNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	var dd []ConsentRequest
	for _, v := range d {
		vd, err := v.toConsentRequest()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if isSubset(scopes, vd.GrantedScopes) {
			dd = append(dd, *vd)
		}
	}

	if len(dd) == 0 {
		return nil, errors.WithStack(pkg.ErrNotFound)
	}

	toSort := byTime(dd)
	sort.Sort(toSort)

	return &dd[0], nil
}
