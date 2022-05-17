package driver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
)

func TestDefaultKeyManager_HsmDisabled(t *testing.T) {
	l := logrusx.New("", "")
	c := config.MustNew(context.Background(), l, configx.SkipValidation())
	c.MustSet(context.Background(), config.KeyDSN, "postgres://user:password@127.0.0.1:9999/postgres")
	c.MustSet(context.Background(), config.HsmEnabled, "false")
	reg, err := NewRegistryFromDSN(context.Background(), c, l, true, false, &contextx.Default{})
	assert.NoError(t, err)
	assert.IsType(t, &sql.Persister{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}
