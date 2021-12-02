//go:build hsm
// +build hsm

package driver

import (
	"context"
	"testing"

	"github.com/ory/hydra/hsm"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
)

func TestDefaultKeyManager_HsmEnabled(t *testing.T) {
	l := logrusx.New("", "")
	c := config.MustNew(l, configx.SkipValidation())
	c.MustSet(config.KeyDSN, "postgres://user:password@127.0.0.1:9999/postgres")
	c.MustSet(config.HsmEnabled, "true")
	reg, err := NewRegistryFromDSN(context.Background(), c, l)
	assert.NoError(t, err)
	assert.IsType(t, &hsm.KeyManager{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}
