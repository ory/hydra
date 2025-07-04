// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx_test

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"

	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/popx"
)

func TestGoMigrations(t *testing.T) {
	var called []time.Time

	goMigrations := popx.Migrations{
		{
			Path:      "gomigration_0",
			Version:   "20000101000000",
			Name:      "gomigration_0",
			Direction: "up",
			Type:      "go",
			DBType:    "all",
			Runner: func(popx.Migration, *pop.Connection, *pop.Tx) error {
				called[0] = time.Now()
				return nil
			},
		},
		{
			Path:      "gomigration_0",
			Version:   "20000101000000",
			Name:      "gomigration_0",
			Direction: "down",
			Type:      "go",
			DBType:    "all",
			Runner: func(_ popx.Migration, _ *pop.Connection, _ *pop.Tx) error {
				called[1] = time.Now()
				return nil
			},
		},
		{
			Path:      "gomigration_1",
			Version:   "20220215110652",
			Name:      "gomigration_1",
			Direction: "up",
			Type:      "go",
			DBType:    "all",
			Runner: func(_ popx.Migration, _ *pop.Connection, _ *pop.Tx) error {
				called[2] = time.Now()
				return nil
			},
		},
		{
			Path:      "gomigration_1",
			Version:   "20220215110652",
			Name:      "gomigration_1",
			Direction: "down",
			Type:      "go",
			DBType:    "all",
			Runner: func(_ popx.Migration, _ *pop.Connection, _ *pop.Tx) error {
				called[3] = time.Now()
				return nil
			},
		},
	}

	t.Run("tc=calls_all_migrations", func(t *testing.T) {
		called = make([]time.Time, len(goMigrations))

		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL: dbal.NewSQLiteTestDatabase(t),
		})
		require.NoError(t, err)
		require.NoError(t, c.Open())

		mb, err := popx.NewMigrationBox(transactionalMigrations, popx.NewMigrator(c, logrusx.New("", ""), nil, 0), popx.WithGoMigrations(goMigrations))
		require.NoError(t, err)
		require.NoError(t, mb.Up(context.Background()))

		assert.Zero(t, called[1])
		assert.Zero(t, called[3])
		assert.NotZero(t, called[0])
		assert.NotZero(t, called[2])
		assert.True(t, called[0].Before(called[2]))

		require.NoError(t, mb.Down(context.Background(), -1))
		assert.NotZero(t, called[1])
		assert.NotZero(t, called[3])
		assert.True(t, called[3].Before(called[1]))
	})

	t.Run("tc=errs_on_missing_down_migration", func(t *testing.T) {
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL: dbal.NewSQLiteTestDatabase(t),
		})
		require.NoError(t, err)
		require.NoError(t, c.Open())

		_, err = popx.NewMigrationBox(transactionalMigrations, popx.NewMigrator(c, logrusx.New("", ""), nil, 0), popx.WithGoMigrations(goMigrations[:1]))
		require.Error(t, err)
	})

	t.Run("tc=runs everything in one transaction", func(t *testing.T) {
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL: dbal.NewSQLiteTestDatabase(t),
		})
		require.NoError(t, err)
		require.NoError(t, c.Open())

		require.NoError(t, c.RawQuery("CREATE TABLE tests (i INTEGER)").Exec())

		errSecondStatement := errors.New("second statement failed as expected")
		mb, err := popx.NewMigrationBox(empty, popx.NewMigrator(c, logrusx.New("", ""), nil, 0), popx.WithGoMigrations(
			popx.Migrations{
				{
					Path:      "gomigration_1",
					Version:   "20220215110652",
					Name:      "gomigration_1",
					Direction: "up",
					Type:      "go",
					DBType:    "all",
					Runner: func(_ popx.Migration, c *pop.Connection, _ *pop.Tx) error {
						if err := c.RawQuery("INSERT INTO tests (i) VALUES (1)").Exec(); err != nil {
							return errors.WithStack(err)
						}
						if err := c.RawQuery("INSERT INTO unknown_table (data) VALUES ('foo')").Exec(); err != nil {
							return errSecondStatement
						}
						return errors.New("this should not be reached")
					},
				},
				{
					Path:      "gomigration_1",
					Version:   "20220215110652",
					Name:      "gomigration_1",
					Direction: "down",
					Type:      "go",
					DBType:    "all",
					Runner: func(_ popx.Migration, c *pop.Connection, _ *pop.Tx) error {
						return nil
					},
				},
			},
		))
		require.NoError(t, err)
		require.ErrorIs(t, mb.Up(context.Background()), errSecondStatement)
		type test struct {
			I int `db:"i"`
		}
		tt := &test{}
		assert.ErrorIs(t, c.Where("i=1").First(tt), sql.ErrNoRows, "%+v", tt)
	})
}

func TestIncompatibleRunners(t *testing.T) {
	mb, err := popx.NewMigrationBox(empty, popx.NewMigrator(nil, logrusx.New("", ""), nil, 0), popx.WithGoMigrations(
		popx.Migrations{
			{
				Path:      "transactional",
				Version:   "1",
				Name:      "gomigration_tx",
				Direction: "up",
				Type:      "go",
				DBType:    "all",
				RunnerNoTx: func(m popx.Migration, c *pop.Connection) error {
					return nil
				},
				Runner: func(m popx.Migration, c *pop.Connection, tx *pop.Tx) error {
					return nil
				},
			},
			{
				Path:      "transactional",
				Version:   "1",
				Name:      "gomigration_tx",
				Direction: "down",
				Type:      "go",
				DBType:    "all",
				RunnerNoTx: func(m popx.Migration, c *pop.Connection) error {
					return nil
				},
			},
		}))
	require.ErrorContains(t, err, "incompatible transaction and non-transaction runners defined")
	require.Nil(t, mb)

	mb, err = popx.NewMigrationBox(empty, popx.NewMigrator(nil, logrusx.New("", ""), nil, 0), popx.WithGoMigrations(
		popx.Migrations{
			{
				Path:       "transactional",
				Version:    "1",
				Name:       "gomigration_tx",
				Direction:  "up",
				Type:       "go",
				DBType:     "all",
				RunnerNoTx: nil,
				Runner:     nil,
			},
			{
				Path:       "transactional",
				Version:    "1",
				Name:       "gomigration_tx",
				Direction:  "down",
				Type:       "go",
				DBType:     "all",
				RunnerNoTx: nil,
				Runner:     nil,
			},
		}))
	require.ErrorContains(t, err, "no runner defined")
	require.Nil(t, mb)
}

func TestNoTransaction(t *testing.T) {
	c, err := pop.NewConnection(&pop.ConnectionDetails{
		URL: dbal.NewSQLiteTestDatabase(t),
	})
	require.NoError(t, err)
	require.NoError(t, c.Open())

	require.NoError(t, c.RawQuery("CREATE TABLE tests (i INTEGER, j INTEGER)").Exec())

	up1, up2 := make(chan struct{}), make(chan struct{})
	down1, down2 := make(chan struct{}), make(chan struct{})
	rnd := rand.NewSource(time.Now().Unix())
	i1, i2, j1, j2 := rnd.Int63(), rnd.Int63(), rnd.Int63(), rnd.Int63()
	mb, err := popx.NewMigrationBox(empty, popx.NewMigrator(c, logrusx.New("", ""), nil, 0), popx.WithGoMigrations(
		popx.Migrations{
			{
				Path:      "gomigration_notx",
				Version:   "1",
				Name:      "gomigration no transaction",
				Direction: "up",
				Type:      "go",
				DBType:    "all",
				RunnerNoTx: func(m popx.Migration, c *pop.Connection) error {
					if _, err := c.Store.Exec("INSERT INTO tests (i, j) VALUES (?, ?)", i1, j1); err != nil {
						return errors.WithStack(err)
					}
					close(up1)
					<-up2
					return nil
				},
			},
			{
				Path:      "gomigration_notx",
				Version:   "1",
				Name:      "gomigration no transaction",
				Direction: "down",
				Type:      "go",
				DBType:    "all",
				RunnerNoTx: func(m popx.Migration, c *pop.Connection) error {
					if _, err := c.Store.Exec("INSERT INTO tests (i, j) VALUES (?, ?)", i2, j2); err != nil {
						return errors.WithStack(err)
					}
					close(down1)
					<-down2
					return nil
				},
			},
		},
	))
	require.NoError(t, err)
	errs := make(chan error, 10)
	go func() {
		errs <- mb.Up(context.Background())
	}()
	<-up1
	var j int64
	require.NoError(t, c.Store.Get(&j, "SELECT j FROM tests WHERE i = ?", i1))
	assert.Equal(t, j1, j)
	close(up2)
	assert.NoError(t, <-errs)

	go func() {
		errs <- mb.Down(context.Background(), 20)
	}()
	<-down1
	j = 0
	require.NoError(t, c.Store.Get(&j, "SELECT j FROM tests WHERE i = ?", i2))
	assert.Equal(t, j2, j)
	close(down2)
	assert.NoError(t, <-errs)
}
