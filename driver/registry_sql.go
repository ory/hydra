package driver

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/dbal"
)

type RegistrySQL struct {
	*RegistryBase
	db *sqlx.DB
}

var _ Registry = new(RegistrySQL)

func init() {
	dbal.RegisterDriver(NewRegistrySQL())
}

type schemaCreator interface {
	CreateSchemas() (int, error)
}

func NewRegistrySQL() *RegistrySQL {
	r := &RegistrySQL{
		RegistryBase: new(RegistryBase),
	}
	r.RegistryBase.with(r)
	return r
}

func (m *RegistrySQL) Init(url string, l logrus.FieldLogger, opts ...dbal.DriverOptionModifier) error {
	m.l = l
	return nil
}

func (m *RegistrySQL) WithDB(db *sqlx.DB) Registry {
	m.db = db
	return m
}

func (m *RegistrySQL) DB() *sqlx.DB {
	return m.db
}

func (m *RegistrySQL) CreateSchemas() (int, error) {
	var total int

	// Ensure dependencies exist
	_, _, _, _ = m.ClientManager(), m.ConsentManager(), m.KeyManager(), m.OAuth2Storage()

	for _, s := range []schemaCreator{
		m.cm.(schemaCreator),
		m.com.(schemaCreator),
		m.fs.(schemaCreator),
		m.km.(schemaCreator),
	} {
		if c, err := s.CreateSchemas(); err != nil {
			return c, err
		} else {
			total += c
		}
	}

	return total, nil
}

func (m *RegistrySQL) CanHandle(dsn string) bool {
	panic("not implemented")
}

func (m *RegistrySQL) Ping() error {
	return m.db.Ping()
}

func (m *RegistrySQL) ClientManager() client.Manager {
	if m.cm == nil {
		m.cm = client.NewSQLManager(m.db, m)
	}
	return m.cm
}

func (m *RegistrySQL) ConsentManager() consent.Manager {
	if m.com == nil {
		m.com = consent.NewSQLManager(m.db, m)
	}
	return m.com
}

func (m *RegistrySQL) OAuth2Storage() x.FositeStorer {
	if m.fs == nil {
		m.fs = oauth2.NewFositeSQLStore(m.db, m.r, m.c)
	}
	return m.fs
}

func (m *RegistrySQL) KeyManager() jwk.Manager {
	if m.km == nil {
		m.km = jwk.NewSQLManager(m.db, m)
	}
	return m.km
}
