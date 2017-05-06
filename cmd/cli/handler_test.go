package cli

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/integration"
	"github.com/ory/ladon"
	lsql "github.com/ory/ladon/manager/sql"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	db = integration.ConnectToPostgres()

	code := m.Run()
	integration.KillAll()
	os.Exit(code)
}

func TestNewHandler(t *testing.T) {
	_ = NewHandler(&config.Config{})
}

func TestMigrateHandlerSQL(t *testing.T) {
	handler := newMigrateHandler(&config.Config{})

	assert.NoError(t, handler.runMigrateSQL(db))

	// create a few policies
	m := lsql.SQLManagerMigrateFromMajor0Minor6ToMajor0Minor7{
		DB: db,
	}

	// create some dummy policies
	for _, p := range []*ladon.DefaultPolicy{
		{
			ID:          uuid.New(),
			Description: "description",
			Subjects:    []string{"user", "anonymous"},
			Effect:      ladon.AllowAccess,
			Resources:   []string{"article", "user"},
			Actions:     []string{"create", "update"},
			Conditions:  ladon.Conditions{},
		},
		{
			ID:          uuid.New(),
			Description: "description",
			Subjects:    []string{},
			Effect:      ladon.AllowAccess,
			Resources:   []string{"<article|user>"},
			Actions:     []string{"view"},
			Conditions:  ladon.Conditions{},
		},
	} {
		require.NoError(t, m.Create(p))
	}

	assert.NoError(t, handler.runMigrateLadon050To060(db))
}
