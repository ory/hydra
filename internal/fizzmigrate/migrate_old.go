package fizzmigrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/internal/fizzmigrate/client"
	"github.com/ory/hydra/internal/fizzmigrate/consent"
	"github.com/ory/hydra/internal/fizzmigrate/jwk"
	"github.com/ory/hydra/internal/fizzmigrate/oauth2"
)

type OldMigrationRunner struct {
	l  logrus.FieldLogger
	db *sqlx.DB
}

type migrator interface {
	CreateSchemas(dbName string) (int, error)
	PlanMigration(dbName string) ([]*migrate.PlannedMigration, error)
}

func (m *OldMigrationRunner) SchemaMigrationPlan(dbName string) (*tablewriter.Table, error) {
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

	for component, s := range []migrator{
		client.NewMigrator(m.db),
		consent.NewMigrator(m.db),
		oauth2.NewMigrator(m.db),
		jwk.NewMigrator(m.db),
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

func (m *OldMigrationRunner) CreateSchemas(dbName string) (int, error) {
	var total int

	m.l.Debugf("Applying %s SQL migrations...", dbName)
	for k, s := range []migrator{
		client.NewMigrator(m.db),
		consent.NewMigrator(m.db),
		oauth2.NewMigrator(m.db),
		jwk.NewMigrator(m.db),
	} {
		m.l.Debugf("Applying %s SQL migrations for manager: %T (%d)", dbName, s, k)
		if c, err := s.CreateSchemas(dbName); err != nil {
			return c, err
		} else {
			m.l.Debugf("Successfully applied %d %s SQL migrations from manager: %T (%d)", c, dbName, s, k)
			total += c
		}
	}
	m.l.Debugf("Applied %d %s SQL migrations", total, dbName)

	return total, nil
}
