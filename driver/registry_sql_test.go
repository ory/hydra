// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/persistence/sql"
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
	if err := r.Init(context.Background(), true, false, &contextx.Default{}, nil); err != nil {
		t.Fatalf("unable to init registry: %s", err)
	}
	assert.NoError(t, err)
	assert.IsType(t, &sql.Persister{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}

func TestDbUnknownTableColumns(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	c := config.MustNew(ctx, l, configx.SkipValidation())
	postgresDsn := dockertest.RunTestPostgreSQL(t)
	c.MustSet(ctx, config.KeyDSN, postgresDsn)
	reg, err := NewRegistryFromDSN(ctx, c, l, false, true, &contextx.Default{})
	require.NoError(t, err)

	statement := "ALTER TABLE \"hydra_client\" ADD COLUMN \"temp_column\" VARCHAR(128) NOT NULL DEFAULT '';"
	require.NoError(t, reg.Persister().Connection(ctx).RawQuery(statement).Exec())

	cl := &client.Client{
		LegacyClientID: strconv.Itoa(rand.Int()),
	}
	require.NoError(t, reg.Persister().CreateClient(ctx, cl))
	getClients := func(reg Registry) ([]client.Client, error) {
		readClients := make([]client.Client, 0)
		return readClients, reg.Persister().Connection(ctx).RawQuery("SELECT * FROM \"hydra_client\"").All(&readClients)
	}

	t.Run("with ignore disabled (default behavior)", func(t *testing.T) {
		_, err := getClients(reg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing destination name temp_column")
	})

	t.Run("with ignore enabled", func(t *testing.T) {
		c.MustSet(ctx, config.KeyDBIgnoreUnknownTableColumns, true)
		reg, err := NewRegistryFromDSN(ctx, c, l, false, true, &contextx.Default{})
		require.NoError(t, err)

		actual, err := getClients(reg)
		require.NoError(t, err)
		assert.Len(t, actual, 1)
	})
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
