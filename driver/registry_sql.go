package driver

import (
	"context"
	"strings"
	"time"

	"github.com/ory/hydra/hsm"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/x/errorsx"

	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"

	"github.com/ory/x/resilience"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/ory/hydra/persistence/sql"

	"github.com/jmoiron/sqlx"

	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
)

type RegistrySQL struct {
	*RegistryBase
	db                *sqlx.DB
	defaultKeyManager jwk.Manager
	initialPing       func(r *RegistrySQL) error
}

var _ Registry = new(RegistrySQL)

// defaultInitialPing is the default function that will be called within RegistrySQL.Init to make sure
// the database is reachable. It can be injected for test purposes by changing the value
// of RegistrySQL.initialPing.
var defaultInitialPing = func(m *RegistrySQL) error {
	if err := resilience.Retry(m.l, 5*time.Second, 5*time.Minute, m.Ping); err != nil {
		m.Logger().Print("Could not ping database: ", err)
		return errorsx.WithStack(err)
	}
	return nil
}

func init() {
	dbal.RegisterDriver(func() dbal.Driver {
		return NewRegistrySQL()
	})
}

func NewRegistrySQL() *RegistrySQL {
	r := &RegistrySQL{
		RegistryBase: new(RegistryBase),
		initialPing:  defaultInitialPing,
	}
	r.RegistryBase.with(r)
	return r
}

func (m *RegistrySQL) Init(ctx context.Context) error {
	if m.persister == nil {
		var opts []instrumentedsql.Opt
		if m.Tracer(ctx).IsLoaded() {
			opts = []instrumentedsql.Opt{
				instrumentedsql.WithTracer(opentracing.NewTracer(true)),
			}
		}

		// new db connection
		pool, idlePool, connMaxLifetime, connMaxIdleTime, cleanedDSN := sqlcon.ParseConnectionOptions(m.l, m.C.DSN())
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL:                       sqlcon.FinalizeDSN(m.l, cleanedDSN),
			IdlePool:                  idlePool,
			ConnMaxLifetime:           connMaxLifetime,
			ConnMaxIdleTime:           connMaxIdleTime,
			Pool:                      pool,
			UseInstrumentedDriver:     m.Tracer(ctx).IsLoaded(),
			InstrumentedDriverOptions: opts,
		})
		if err != nil {
			return errorsx.WithStack(err)
		}
		if err := resilience.Retry(m.l, 5*time.Second, 5*time.Minute, c.Open); err != nil {
			return errorsx.WithStack(err)
		}
		m.persister, err = sql.NewPersister(ctx, c, m, m.C, m.l)
		if err != nil {
			return err
		}
		if err := m.initialPing(m); err != nil {
			return err
		}

		if m.C.HsmEnabled() {
			hardwareKeyManager := hsm.NewKeyManager(m.HsmContext(), m.C)
			m.defaultKeyManager = jwk.NewManagerStrategy(hardwareKeyManager, m.persister)
		} else {
			m.defaultKeyManager = m.persister
		}

		// if dsn is memory we have to run the migrations on every start
		// use case - such as
		// - just in memory
		// - shared connection
		// - shared but unique in the same process
		// see: https://sqlite.org/inmemorydb.html
		if dbal.IsMemorySQLite(m.C.DSN()) {
			m.Logger().Print("Hydra is running migrations on every startup as DSN is memory.\n")
			m.Logger().Print("This means your data is lost when Hydra terminates.\n")
			if err := m.persister.MigrateUp(context.Background()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *RegistrySQL) alwaysCanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	s := dbal.Canonicalize(scheme)
	return s == dbal.DriverMySQL || s == dbal.DriverPostgreSQL || s == dbal.DriverCockroachDB
}

func (m *RegistrySQL) Ping() error {
	return m.Persister().Ping()
}

func (m *RegistrySQL) ClientManager() client.Manager {
	return m.Persister()
}

func (m *RegistrySQL) ConsentManager() consent.Manager {
	return m.Persister()
}

func (m *RegistrySQL) OAuth2Storage() x.FositeStorer {
	return m.Persister()
}

func (m *RegistrySQL) KeyManager() jwk.Manager {
	return m.defaultKeyManager
}

func (m *RegistrySQL) SoftwareKeyManager() jwk.Manager {
	return m.Persister()
}

func (m *RegistrySQL) GrantManager() trust.GrantManager {
	return m.Persister()
}
