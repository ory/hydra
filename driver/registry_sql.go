package driver

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/hydra/hsm"
	"github.com/ory/x/contextx"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/popx"

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

func (m *RegistrySQL) determineNetwork(c *pop.Connection, ctx context.Context) (*networkx.Network, error) {
	mb, err := popx.NewMigrationBox(networkx.Migrations, popx.NewMigrator(c, m.Logger(), m.Tracer(ctx), 0))
	if err != nil {
		return nil, err
	}
	s, err := mb.Status(ctx)
	if err != nil {
		return nil, err
	}
	if s.HasPending() {
		return nil, errors.WithStack(errors.New("some migrations are pending"))
	}

	return networkx.NewManager(c, m.Logger(), m.Tracer(ctx)).Determine(ctx)
}

func (m *RegistrySQL) Init(ctx context.Context, skipNetworkInit bool, migrate bool, ctxer contextx.Contextualizer) error {
	if m.persister == nil {
		m.WithContextualizer(ctxer)
		var opts []instrumentedsql.Opt
		if m.Tracer(ctx).IsLoaded() {
			opts = []instrumentedsql.Opt{
				instrumentedsql.WithTracer(opentracing.NewTracer(true)),
			}
		}

		// new db connection
		pool, idlePool, connMaxLifetime, connMaxIdleTime, cleanedDSN := sqlcon.ParseConnectionOptions(m.l, m.Config().DSN(ctx))
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

		p, err := sql.NewPersister(ctx, c, m, m.Config(), m.l)
		if err != nil {
			return err
		}

		// if dsn is memory we have to run the migrations on every start
		// use case - such as
		// - just in memory
		// - shared connection
		// - shared but unique in the same process
		// see: https://sqlite.org/inmemorydb.html
		if dbal.IsMemorySQLite(m.Config().DSN(ctx)) {
			m.Logger().Print("Hydra is running migrations on every startup as DSN is memory.\n")
			m.Logger().Print("This means your data is lost when Hydra terminates.\n")
			if err := p.MigrateUp(context.Background()); err != nil {
				return err
			}
		} else if migrate {
			if err := p.MigrateUp(context.Background()); err != nil {
				return err
			}
		}

		if skipNetworkInit {
			m.persister = p
		} else {
			net, err := p.DetermineNetwork(ctx)
			if err != nil {
				m.Logger().WithError(err).Warnf("Unable to determine network, retrying.")
				return err
			}

			m.persister = p.WithFallbackNetworkID(net.ID)
		}

		if m.Config().HsmEnabled(ctx) {
			hardwareKeyManager := hsm.NewKeyManager(m.HsmContext())
			m.defaultKeyManager = jwk.NewManagerStrategy(hardwareKeyManager, m.persister)
		} else {
			m.defaultKeyManager = m.persister
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
	return m.Persister().Connection(context.Background()).Open()
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
