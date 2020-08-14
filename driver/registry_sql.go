package driver

import (
	"strings"
	"time"

	"github.com/ory/x/resilience"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/ory/hydra/persistence/sql"

	"github.com/jmoiron/sqlx"

	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
)

type RegistrySQL struct {
	*RegistryBase
	db          *sqlx.DB
	dbalOptions []sqlcon.OptionModifier
}

var _ Registry = new(RegistrySQL)

func init() {
	dbal.RegisterDriver(func() dbal.Driver {
		return NewRegistrySQL()
	})
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
	if m.db == nil {
		// old db connection
		options := append([]sqlcon.OptionModifier{}, m.dbalOptions...)
		if m.Tracer().IsLoaded() {
			options = append(options, sqlcon.WithDistributedTracing(), sqlcon.WithOmitArgsFromTraceSpans())
		}

		connection, err := sqlcon.NewSQLConnection(m.C.DSN(), m.Logger(), options...)
		if err != nil {
			return err
		}

		m.db, err = connection.GetDatabaseRetry(time.Second*5, time.Minute*5)
		if err != nil {
			return err
		}
	}

	if m.persister == nil {
		// new db connection
		pool, idlePool, connMaxLifetime, cleanedDSN := sqlcon.ParseConnectionOptions(m.l, m.C.DSN())
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL:             sqlcon.FinalizeDSN(m.l, cleanedDSN),
			IdlePool:        idlePool,
			ConnMaxLifetime: connMaxLifetime,
			Pool:            pool,
		})
		if err != nil {
			return errors.WithStack(err)
		}
		if err := resilience.Retry(m.l, 5*time.Second, 5*time.Minute, c.Open); err != nil {
			return errors.WithStack(err)
		}
		m.persister, err = sql.NewPersister(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *RegistrySQL) DB() *sqlx.DB {
	if m.db == nil {
		if err := m.Init(); err != nil {
			m.Logger().WithError(err).Fatalf("Unable to initialize database.")
		}
	}

	return m.db
}

func (m *RegistrySQL) CanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	s := dbal.Canonicalize(scheme)
	return s == dbal.DriverMySQL || s == dbal.DriverPostgreSQL || s == dbal.DriverCockroachDB
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
		m.fs = oauth2.NewFositeSQLStore(m.DB(), m.r, m.C, m.KeyCipher())
	}
	return m.fs
}

func (m *RegistrySQL) KeyManager() jwk.Manager {
	if m.km == nil {
		m.km = jwk.NewSQLManager(m.DB(), m)
	}
	return m.km
}
