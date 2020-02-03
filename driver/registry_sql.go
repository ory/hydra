package driver

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"

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

type schemaCreator interface {
	CreateSchemas(dbName string) (int, error)
	PlanMigration(dbName string) ([]*migrate.PlannedMigration, error)
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

func (m *RegistrySQL) SchemaMigrationPlan(dbName string) (*tablewriter.Table, error) {
	names := map[int]string{
		0: "JSON Web Keys",
		1: "OAuth 2.0 Clients",
		2: "Login &Consent",
		3: "OAuth 2.0",
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetColMinWidth(4, 20)
	table.SetHeader([]string{
		"Driver",
		"Module",
		"ID",
		"#",
		"Query",
	})

	for component, s := range []schemaCreator{
		m.KeyManager().(schemaCreator),
		m.ClientManager().(schemaCreator),
		m.ConsentManager().(schemaCreator),
		m.OAuth2Storage().(schemaCreator),
	} {
		plans, err := s.PlanMigration(dbName)
		if err != nil {
			return nil, err
		}

		for _, plan := range plans {
			for k, up := range plan.Up {
				up = strings.Replace(strings.TrimSpace(up), "\n", "", -1)
				up = strings.Join(strings.Fields(up), " ")
				if len(up) > 0 {
					table.Append([]string{m.db.DriverName(), names[component], plan.Id + ".sql", fmt.Sprintf("%d", k), up})
				}
			}
		}
	}

	return table, nil
}

func (m *RegistrySQL) CreateSchemas(dbName string) (int, error) {
	var total int

	m.Logger().Debugf("Applying %s SQL migrations...", dbName)
	for k, s := range []schemaCreator{
		m.KeyManager().(schemaCreator),
		m.ClientManager().(schemaCreator),
		m.ConsentManager().(schemaCreator),
		m.OAuth2Storage().(schemaCreator),
	} {
		m.Logger().Debugf("Applying %s SQL migrations for manager: %T (%d)", dbName, s, k)
		if c, err := s.CreateSchemas(dbName); err != nil {
			return c, err
		} else {
			m.Logger().Debugf("Successfully applied %d %s SQL migrations from manager: %T (%d)", c, dbName, s, k)
			total += c
		}
	}
	m.Logger().Debugf("Applied %d %s SQL migrations", total, dbName)

	return total, nil
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
