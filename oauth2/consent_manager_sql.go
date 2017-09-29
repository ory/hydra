package oauth2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

var sqlConsentParams = []string{
	"id", "audience", "expires_at", "redirect_url", "requested_scopes",
	"csrf", "granted_scopes", "access_token_extra", "id_token_extra",
	"consent", "deny_reason", "subject",
}

var consentMigrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "1",
			Up: []string{`CREATE TABLE IF NOT EXISTS hydra_consent_request (
	id      			varchar(36) NOT NULL PRIMARY KEY,
	requested_scopes 	text NOT NULL,
	audience 			text NOT NULL,
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
	},
}

type consentRequestSqlData struct {
	ID               string    `db:"id"`
	RequestedScopes  string    `db:"requested_scopes"`
	Audience         string    `db:"audience"`
	ExpiresAt        time.Time `db:"expires_at"`
	RedirectURL      string    `db:"redirect_url"`
	CSRF             string    `db:"csrf"`
	GrantedScopes    string    `db:"granted_scopes"`
	AccessTokenExtra string    `db:"access_token_extra"`
	IDTokenExtra     string    `db:"id_token_extra"`
	Consent          string    `db:"consent"`
	DenyReason       string    `db:"deny_reason"`
	Subject          string    `db:"subject"`
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

	return &consentRequestSqlData{
		ID:               request.ID,
		RequestedScopes:  strings.Join(request.RequestedScopes, " "),
		GrantedScopes:    strings.Join(request.GrantedScopes, " "),
		Audience:         request.Audience,
		ExpiresAt:        request.ExpiresAt,
		RedirectURL:      request.RedirectURL,
		CSRF:             request.CSRF,
		AccessTokenExtra: atext,
		IDTokenExtra:     idtext,
		Consent:          request.Consent,
		DenyReason:       request.DenyReason,
		Subject:          request.Subject,
	}, nil
}

func (r *consentRequestSqlData) toConsentRequest() (*ConsentRequest, error) {
	var atext, idtext map[string]interface{}

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

	return &ConsentRequest{
		ID:               r.ID,
		Audience:         r.Audience,
		ExpiresAt:        r.ExpiresAt,
		RedirectURL:      r.RedirectURL,
		CSRF:             r.CSRF,
		Consent:          r.Consent,
		DenyReason:       r.DenyReason,
		RequestedScopes:  strings.Split(r.RequestedScopes, " "),
		GrantedScopes:    strings.Split(r.GrantedScopes, " "),
		AccessTokenExtra: atext,
		IDTokenExtra:     idtext,
		Subject:          r.Subject,
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
	n, err := migrate.Exec(m.db.DB, m.db.DriverName(), consentMigrations, migrate.Up)
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
