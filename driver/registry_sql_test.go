package driver

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlcon/dockertest"
)

func TestDefaultKeyManager_HsmDisabled(t *testing.T) {
	l := logrusx.New("", "")
	c := config.MustNew(context.Background(), l, configx.SkipValidation())
	c.MustSet(context.Background(), config.KeyDSN, "postgres://user:password@127.0.0.1:9999/postgres")
	c.MustSet(context.Background(), config.HSMEnabled, "false")
	reg, err := NewRegistryWithoutInit(c, l)
	r := reg.(*RegistrySQL)
	r.initialPing = sussessfulPing()
	if err := r.Init(context.Background(), true, false, &contextx.Default{}); err != nil {
		t.Fatalf("unable to init registry: %s", err)
	}
	assert.NoError(t, err)
	assert.IsType(t, &sql.Persister{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}

func TestDbUnknownTableColumns(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    string
		expectError  bool
		expectedSize int
	}{
		{name: "with unsafe", flagValue: "true", expectError: false, expectedSize: 1},
		{name: "without unsafe", flagValue: "false", expectError: true, expectedSize: 0},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				ctx := context.Background()
				l := logrusx.New("", "")
				c := config.MustNew(ctx, l, configx.SkipValidation())
				postgresDsn := dockertest.RunTestPostgreSQL(t)
				c.MustSet(ctx, config.KeyDSN, postgresDsn)
				c.MustSet(ctx, config.KeyDbIgnoreUnknownTableColumns, test.flagValue)
				reg, err := NewRegistryFromDSN(ctx, c, l, false, true, &contextx.Default{})
				assert.NoError(t, err)

				statement := "ALTER TABLE \"hydra_client\" ADD COLUMN \"temp_column\" VARCHAR(128) NOT NULL DEFAULT '';"
				err = reg.Persister().Connection(ctx).RawQuery(statement).Exec()
				assert.NoError(t, err)

				cl := &client.Client{
					LegacyClientID: strconv.Itoa(rand.Int()),
				}

				err = reg.Persister().CreateClient(ctx, cl)
				assert.NoError(t, err)

				readClients := make([]client.Client, 0)
				err = reg.Persister().Connection(ctx).RawQuery("SELECT * FROM \"hydra_client\"").All(&readClients)
				if test.expectError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "missing destination name temp_column")
				} else {
					assert.NoError(t, err)
				}
				assert.Len(t, readClients, test.expectedSize)
			},
		)
	}
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
