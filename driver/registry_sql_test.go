package driver

import (
	"context"
	"github.com/ory/x/dbal"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
)

func TestDefaultKeyManager_HsmDisabled(t *testing.T) {
	l := logrusx.New("", "")
	c := config.MustNew(context.Background(), l, configx.SkipValidation())
	c.MustSet(config.KeyDSN, dbal.SQLiteInMemory)
	c.MustSet(config.HsmEnabled, "false")
	reg, err := NewRegistryFromDSN(context.Background(), c, l)
	assert.NoError(t, err)
	assert.IsType(t, &sql.Persister{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}
