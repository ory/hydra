// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package batch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/x/dbal"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/snapshotx"
	"github.com/ory/x/sqlxx"
)

type (
	testModel struct {
		ID          uuid.UUID       `db:"id"`
		NID         uuid.UUID       `db:"nid"`
		String      string          `db:"string"`
		Int         int             `db:"int"`
		NullTimePtr *sqlxx.NullTime `db:"null_time_ptr"`
		CreatedAt   time.Time       `json:"created_at" db:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	}
	testQuoter struct{}
)

func (i testModel) TableName(ctx context.Context) string {
	return "test_models"
}

func (tq testQuoter) Quote(s string) string { return fmt.Sprintf("%q", s) }

func makeModels[T any]() []*T {
	models := make([]*T, 10)
	for k := range models {
		models[k] = new(T)
	}
	return models
}

func Test_buildInsertQueryArgs(t *testing.T) {
	ctx := context.Background()
	t.Run("case=testModel", func(t *testing.T) {
		models := makeModels[testModel]()
		mapper := reflectx.NewMapper("db")
		args := buildInsertQueryArgs(ctx, "other", mapper, testQuoter{}, models)
		snapshotx.SnapshotT(t, args)

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES\n%s", args.TableName, args.ColumnsDecl, args.Placeholders)
		assert.Equal(t, `INSERT INTO "test_models" ("created_at", "id", "int", "nid", "null_time_ptr", "string", "updated_at") VALUES
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?),
(?, ?, ?, ?, ?, ?, ?)`, query)
	})

	t.Run("case=cockroach", func(t *testing.T) {
		models := makeModels[testModel]()
		for k := range models {
			if k%3 == 0 {
				models[k].ID = uuid.FromStringOrNil(fmt.Sprintf("ae0125a9-2786-4ada-82d2-d169cf75047%d", k))
			}
		}
		mapper := reflectx.NewMapper("db")
		args := buildInsertQueryArgs(ctx, "cockroach", mapper, testQuoter{}, models)
		snapshotx.SnapshotT(t, args)
	})
}

func Test_buildInsertQueryValues(t *testing.T) {
	t.Run("case=testModel", func(t *testing.T) {
		model := &testModel{
			String: "string",
			Int:    42,
		}
		mapper := reflectx.NewMapper("db")

		nowFunc := func() time.Time {
			return time.Time{}
		}
		t.Run("case=cockroach", func(t *testing.T) {
			values, err := buildInsertQueryValues(dbal.DriverCockroachDB, mapper, []string{"created_at", "updated_at", "id", "string", "int", "null_time_ptr", "traits"}, []*testModel{model}, nowFunc)
			require.NoError(t, err)
			snapshotx.SnapshotT(t, values)
		})

		t.Run("case=others", func(t *testing.T) {
			values, err := buildInsertQueryValues("other", mapper, []string{"created_at", "updated_at", "id", "string", "int", "null_time_ptr", "traits"}, []*testModel{model}, nowFunc)
			require.NoError(t, err)

			assert.NotNil(t, model.CreatedAt)
			assert.Equal(t, model.CreatedAt, values[0])

			assert.NotNil(t, model.UpdatedAt)
			assert.Equal(t, model.UpdatedAt, values[1])

			assert.NotZero(t, model.ID)
			assert.Equal(t, model.ID, values[2])

			assert.Equal(t, model.String, values[3])
			assert.Equal(t, model.Int, values[4])

			assert.Nil(t, model.NullTimePtr)

		})
	})
}
