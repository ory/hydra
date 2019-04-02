package driver

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/dbal"
)

type RegistrySQL struct {
	*RegistryBase
	db          *sqlx.DB
	dbalOptions []sqlcon.OptionModifier
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

func (m *RegistrySQL) WithDB(db *sqlx.DB) Registry {
	m.db = db
	return m
}

func (m *RegistrySQL) Init() error {
	if m.db != nil {
		return nil
	}

	options := append([]sqlcon.OptionModifier{}, m.dbalOptions...)
	if m.Tracer().IsLoaded() {
		options = append(options, sqlcon.WithDistributedTracing(), sqlcon.WithOmitArgsFromTraceSpans())
	}

	connection, err := sqlcon.NewSQLConnection(m.c.DSN(), m.Logger(), options...)
	if err != nil {
		return err
	}

	m.db, err = connection.GetDatabaseRetry(time.Second*5, time.Minute*5)
	if err != nil {
		return err
	}

	return err
}

func (m *RegistrySQL) DB() *sqlx.DB {
	if m.db == nil {
		if err := m.Init(); err != nil {
			m.Logger().WithError(err).Fatalf("Unable to initialize database.")
		}
	}

	return m.db
}

func (m *RegistrySQL) CreateSchemas() (int, error) {
	var total int

	m.Logger().Debugf("Applying %s SQL migrations...", m.db.DriverName())
	for k, s := range []schemaCreator{
		m.KeyManager().(schemaCreator),
		m.ClientManager().(schemaCreator),
		m.ConsentManager().(schemaCreator),
		m.OAuth2Storage().(schemaCreator),
	} {
		m.Logger().Debugf("Applying %s SQL migrations for manager: %T (%d)", m.db.DriverName(), s, k)
		if c, err := s.CreateSchemas(); err != nil {
			return c, err
		} else {
			m.Logger().Debugf("Successfully applied %d %s SQL migrations from manager: %T (%d)", c, m.db.DriverName(), s, k)
			total += c
		}
	}
	m.Logger().Debugf("Applied %d %s SQL migrations", total, m.db.DriverName())

	return total, nil
}

func (m *RegistrySQL) CanHandle(dsn string) bool {
	s := dbal.Canonicalize(urlx.ParseOrFatal(m.l, dsn).Scheme)
	return s == dbal.DriverMySQL || s == dbal.DriverPostgreSQL
}

func (m *RegistrySQL) Ping() error {
	return m.DB().Ping()
}

func (m *RegistrySQL) ClientManager() client.Manager {
	if m.cm == nil {
		m.cm = client.NewSQLManager(m.DB(), m)
	}
	return m.cm
}

func (m *RegistrySQL) ConsentManager() consent.Manager {
	if m.com == nil {
		m.com = consent.NewSQLManager(m.DB(), m)
	}
	return m.com
}

func (m *RegistrySQL) OAuth2Storage() x.FositeStorer {
	if m.fs == nil {
		m.fs = oauth2.NewFositeSQLStore(m.DB(), m.r, m.c)
	}
	return m.fs
}

func (m *RegistrySQL) KeyManager() jwk.Manager {
	if m.km == nil {
		m.km = jwk.NewSQLManager(m.DB(), m)
	}
	return m.km
}
