package sql

import (
	"context"
	"fmt"
	"github.com/ory/x/sqlcon"
	"io"
	"strconv"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/ory/hydra/persistence"
)

var _ persistence.Persister = new(Persister)
var migrations = packr.New("migrations", "migrations")

type Persister struct {
	c  *pop.Connection
	mb pop.MigrationBox
}

func NewPersister(c *pop.Connection) (*Persister, error) {
	mb, err := pop.NewMigrationBox(migrations, c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Persister{
		c,
		mb,
	}, nil
}

func (p *Persister) MigrationStatus(_ context.Context, w io.Writer) error {
	return errors.WithStack(p.mb.Status(w))
}

func (p *Persister) MigrateDown(_ context.Context, steps int) error {
	return errors.WithStack(p.mb.Down(steps))
}

func (p *Persister) MigrateUp(_ context.Context) error {
	if err := p.migrateOldMigrationTables(); err != nil {
		return err
	}
	return errors.WithStack(p.mb.Up())
}

func (p *Persister) MigrateUpTo(_ context.Context, steps int) (int, error) {
	if err := p.migrateOldMigrationTables(); err != nil {
		return 0, err
	}
	n, err := p.mb.UpTo(steps)
	return n, errors.WithStack(err)
}

func (p *Persister) PrepareMigration(_ context.Context) error {
	return p.migrateOldMigrationTables()
}

func (p *Persister) Connection(_ context.Context) *pop.Connection {
	return p.c
}

type oldTableName string

const (
	clientMigrationTableName  oldTableName = "hydra_client_migration"
	jwkMigrationTableName     oldTableName = "hydra_jwk_migration"
	consentMigrationTableName oldTableName = "hydra_oauth2_authentication_consent_migration"
	oauth2MigrationTableName  oldTableName = "hydra_oauth2_migration"
)

// this type is copied from sql-migrate to remove the dependency
type OldMigrationRecord struct {
	Id        string    `db:"id"`
	AppliedAt time.Time `db:"applied_at"`
}

// this function is idempotent
func (p *Persister) migrateOldMigrationTables() error {
	if err := p.c.RawQuery(fmt.Sprintf("SELECT * FROM %s", clientMigrationTableName)).Exec(); err != nil {
		// assume there are no old migration tables => done
		return nil
	}

	return errors.WithStack(p.c.Transaction(func(tx *pop.Connection) error {
		if err := pop.CreateSchemaMigrations(p.c); err != nil {
			return errors.WithStack(err)
		}

		// in this order the migrations only depend on already done ones
		for i, table := range []oldTableName{clientMigrationTableName, jwkMigrationTableName, consentMigrationTableName, oauth2MigrationTableName} {
			// in some cases the tables might not exist, so we just add empty ones
			err := errors.WithStack(
				tx.RawQuery(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ();", table)).
					Exec())
			if err != nil {
				return err
			}

			// get old migrations
			var migrations []OldMigrationRecord

			/* #nosec G201 TableName is static */
			err = errors.WithStack(
				tx.RawQuery(fmt.Sprintf("SELECT * FROM %s", table)).
					Eager().
					All(&migrations))
			if err != nil {
				return err
			}

			// translate migrations
			for _, m := range migrations {
				if m.AppliedAt.Before(time.Now()) {
					// the migration was run already -> set it as run for fizz
					// fizz standard version pattern: YYYYMMDDhhmmss
					migrationNumber, err := strconv.ParseInt(m.Id, 10, 0)
					if err != nil {
						return errors.WithStack(err)
					}

					/* #nosec G201 - i is static (0..3) and migrationNumber is from the database */
					if err := tx.RawQuery(
						fmt.Sprintf("INSERT INTO schema_migration (version) VALUES (2019%02d%08d)", i+1, migrationNumber)).
						Exec(); err != nil {
						return sqlcon.HandleError(err)
					}
				}
			}

			// delete old migration table
			if err := tx.RawQuery(fmt.Sprintf("DROP TABLE %s", table)).Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
		}

		return nil
	}))
}

