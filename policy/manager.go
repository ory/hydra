package policy

import (
	"github.com/jmoiron/sqlx"
	"github.com/ory-am/ladon"
	"github.com/rubenv/sql-migrate"
)

// Manager is responsible for managing and persisting policies.
type Manager interface {
	ladon.Manager

	// Update a policy.
	Update(policy ladon.Policy) error
}

func NewSQLManager(db *sqlx.DB) (ladon.Manager, error) {
	m := ladon.NewSQLManager(db, nil)
	migrate.SetTable("hydra_policy_migration")
	if err := m.CreateSchemas(); err != nil {
		return nil, err
	}
	return m, nil
}
