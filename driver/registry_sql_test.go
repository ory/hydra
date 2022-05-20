package driver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/errorsx"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
)

func TestDefaultKeyManager_HsmDisabled(t *testing.T) {
	l := logrusx.New("", "")
	c := config.MustNew(context.Background(), l, configx.SkipValidation())
	c.MustSet(config.KeyDSN, "postgres://user:password@127.0.0.1:9999/postgres")
	c.MustSet(config.HsmEnabled, "false")
	reg, err := NewRegistryWithoutInit(c, l)
	r := reg.(*RegistrySQL)
	r.initialPing = sussessfulPing()
	if err := r.Init(context.Background()); err != nil {
		t.Fatalf("unable to init registry: %s", err)
	}
	assert.NoError(t, err)
	assert.IsType(t, &sql.Persister{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}

func sussessfulPing() func(r *RegistrySQL) error {
	return func(r *RegistrySQL) error {
		// fake that ping is successful
		return nil
	}
}

func failedPing(err error) func(r *RegistrySQL) error {
	return func(r *RegistrySQL) error {
		r.Logger().Fatalf(err.Error())
		return errorsx.WithStack(err)
	}
}
