// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"io"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/ory/x/logrusx"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"
)

// TestMigrator is a modified pop.FileMigrator
type TestMigrator struct {
	*Migrator
}

// NewTestMigrator returns a new TestMigrator
// After running each migration it applies it's corresponding testData sql files.
// They are identified by having the same version (= number in the front of the filename).
// The filenames are expected to be of the format ([0-9]+).*(_testdata(\.[dbtype])?.sql
func NewTestMigrator(t *testing.T, c *pop.Connection, migrations, testData fs.FS, l *logrusx.Logger) *TestMigrator {
	tm := TestMigrator{
		Migrator: NewMigrator(c, l, nil, time.Minute),
	}

	runner := func(mf Migration, c *pop.Connection, tx *pop.Tx) error {
		b, err := fs.ReadFile(migrations, mf.Path)
		require.NoError(t, err)

		content, err := ParameterizedMigrationContent(nil)(mf, c, b, true)
		require.NoError(t, err)

		if len(strings.TrimSpace(content)) != 0 {
			_, err = tx.Exec(content)
			if err != nil {
				return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
			}
		}

		t.Logf("Applied: %s", mf.Version)

		if mf.Direction != "up" {
			return nil
		}

		appliedVersion := mf.Version[:14]

		// find migration index
		if len(mf.Version) > 14 {
			upMigrations := tm.Migrations["up"].SortAndFilter(c.Dialect.Name())
			mgs := upMigrations

			require.False(t, len(mgs) == 0)

			var migrationIndex = -1
			for k, m := range mgs {
				if m.Version == mf.Version {
					migrationIndex = k
					break
				}
			}

			require.NotEqual(t, -1, migrationIndex)

			if migrationIndex+1 > len(mgs)-1 {
				//
			} else {
				require.EqualValues(t, mf.Version, mgs[migrationIndex].Version)
				require.NotEqual(t, mf.Version, mgs[migrationIndex+1].Version)

				nextMigration := mgs[migrationIndex+1]
				if nextMigration.Version[:14] > appliedVersion {
					t.Logf("Executing transactional interim version %s (%s) because next is %s (%s)", mf.Version, appliedVersion, nextMigration.Version, nextMigration.Version[:14])
				} else if nextMigration.Version[:14] == appliedVersion {
					t.Logf("Skipping transactional interim version %s (%s) because next is %s (%s)", mf.Version, appliedVersion, nextMigration.Version, nextMigration.Version[:14])
					return nil
				} else {
					panic("asdf")
				}
			}
		}

		t.Logf("Adding migration test data %s (%s)", mf.Version, appliedVersion)

		// exec testdata
		f, err := testData.Open(appliedVersion + "_testdata." + c.Dialect.Name() + ".sql")
		if errors.Is(err, fs.ErrNotExist) {
			// could not find specific test data; try generic
			f, err = testData.Open(appliedVersion + "_testdata.sql")
			if errors.Is(err, fs.ErrNotExist) {
				// found no test data
				t.Logf("Found no test data for migration %s %s", mf.Version, mf.DBType)
				return nil
			} else if err != nil {
				return errors.WithStack(err)
			}
		} else if err != nil {
			return errors.WithStack(err)
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}

		fi, err := f.Stat()
		if err != nil {
			return errors.WithStack(err)
		}
		if len(strings.TrimSpace(string(data))) == 0 {
			t.Logf("data is empty for: %s", fi.Name())
			return nil
		}

		return nil
	}

	require.NoError(t, fs.WalkDir(migrations, ".", func(p string, info fs.DirEntry, err error) error {
		if !info.IsDir() {
			match, err := pop.ParseMigrationFilename(info.Name())
			if err != nil {
				return err
			}
			if match == nil {
				return nil
			}

			mf := Migration{
				Path:      p,
				Version:   match.Version,
				Name:      match.Name,
				DBType:    match.DBType,
				Direction: match.Direction,
				Type:      match.Type,
				Runner:    runner,
			}
			tm.Migrations[mf.Direction] = append(tm.Migrations[mf.Direction], mf)
		}
		return nil
	}))

	return &tm
}
