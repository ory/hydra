// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"io"
	"io/fs"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"

	"github.com/ory/x/logrusx"
)

type (
	// MigrationBox is a embed migration box.
	MigrationBox struct {
		*Migrator

		Dir              fs.FS
		l                *logrusx.Logger
		migrationContent MigrationContent
		goMigrations     Migrations
	}
	MigrationContent   func(mf Migration, c *pop.Connection, r []byte, usingTemplate bool) (string, error)
	MigrationBoxOption func(*MigrationBox) *MigrationBox
)

func WithTemplateValues(v map[string]interface{}) MigrationBoxOption {
	return func(m *MigrationBox) *MigrationBox {
		m.migrationContent = ParameterizedMigrationContent(v)
		return m
	}
}

func WithMigrationContentMiddleware(middleware func(content string, err error) (string, error)) MigrationBoxOption {
	return func(m *MigrationBox) *MigrationBox {
		prev := m.migrationContent
		m.migrationContent = func(mf Migration, c *pop.Connection, r []byte, usingTemplate bool) (string, error) {
			return middleware(prev(mf, c, r, usingTemplate))
		}
		return m
	}
}

// WithGoMigrations adds migrations that have a custom migration runner.
// TEST THEM THOROUGHLY!
// It will be very hard to fix a buggy migration.
func WithGoMigrations(migrations Migrations) MigrationBoxOption {
	return func(m *MigrationBox) *MigrationBox {
		m.goMigrations = migrations
		return m
	}
}

// WithTestdata adds testdata to the migration box.
func WithTestdata(t *testing.T, testdata fs.FS) MigrationBoxOption {
	testdataPattern := regexp.MustCompile(`^(\d+)_testdata(|\.[a-zA-Z0-9]+).sql$`)
	return func(m *MigrationBox) *MigrationBox {
		require.NoError(t, fs.WalkDir(testdata, ".", func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			match := testdataPattern.FindStringSubmatch(info.Name())
			if len(match) != 2 && len(match) != 3 {
				t.Logf(`WARNING! Found a test migration which does not match the test data pattern: %s`, info.Name())
				return nil
			}

			version := match[1]
			flavor := "all"
			if len(match) == 3 && len(match[2]) > 0 {
				flavor = pop.CanonicalDialect(strings.TrimPrefix(match[2], "."))
			}

			//t.Logf("Found test migration \"%s\" (%s, %+v): %s", flavor, match, err, info.Name())

			m.Migrations["up"] = append(m.Migrations["up"], Migration{
				Version:   version + "9", // run testdata after version
				Path:      path,
				Name:      info.Name(),
				DBType:    flavor,
				Direction: "up",
				Type:      "sql",
				Runner: func(m Migration, _ *pop.Connection, tx *pop.Tx) error {
					b, err := fs.ReadFile(testdata, m.Path)
					if err != nil {
						return err
					}
					if isMigrationEmpty(string(b)) {
						return nil
					}
					_, err = tx.Exec(string(b))
					//match := match
					//t.Logf("Ran test migration \"%s\" (%s, %+v) with error \"%v\" and content:\n %s", m.Path, m.DBType, match, err, string(b))
					return err
				},
			})

			m.Migrations["down"] = append(m.Migrations["down"], Migration{
				Version:   version + "9", // run testdata after version
				Path:      path,
				Name:      info.Name(),
				DBType:    flavor,
				Direction: "down",
				Type:      "sql",
				Runner: func(m Migration, _ *pop.Connection, tx *pop.Tx) error {
					return nil
				},
			})

			sort.Sort(m.Migrations["up"])
			sort.Sort(sort.Reverse(m.Migrations["down"]))
			return nil
		}))
		return m
	}
}

var emptySQLReplace = regexp.MustCompile(`(?m)^(\s*--.*|\s*)$`)

func isMigrationEmpty(content string) bool {
	return len(strings.ReplaceAll(emptySQLReplace.ReplaceAllString(content, ""), "\n", "")) == 0
}

// NewMigrationBox creates a new migration box.
func NewMigrationBox(dir fs.FS, m *Migrator, opts ...MigrationBoxOption) (*MigrationBox, error) {
	mb := &MigrationBox{
		Migrator:         m,
		Dir:              dir,
		l:                m.l,
		migrationContent: ParameterizedMigrationContent(nil),
	}

	for _, o := range opts {
		mb = o(mb)
	}

	txRunner := func(b []byte) func(Migration, *pop.Connection, *pop.Tx) error {
		return func(mf Migration, c *pop.Connection, tx *pop.Tx) error {
			content, err := mb.migrationContent(mf, c, b, true)
			if err != nil {
				return errors.Wrapf(err, "error processing %s", mf.Path)
			}
			if isMigrationEmpty(content) {
				m.l.WithField("migration", mf.Path).Trace("This is usually ok - ignoring migration because content is empty. This is ok!")
				return nil
			}
			if _, err = tx.Exec(content); err != nil {
				return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
			}
			return nil
		}
	}

	autoCommitRunner := func(b []byte) func(Migration, *pop.Connection) error {
		return func(mf Migration, c *pop.Connection) error {
			content, err := mb.migrationContent(mf, c, b, true)
			if err != nil {
				return errors.Wrapf(err, "error processing %s", mf.Path)
			}
			if isMigrationEmpty(content) {
				m.l.WithField("migration", mf.Path).Trace("This is usually ok - ignoring migration because content is empty. This is ok!")
				return nil
			}
			if _, err = c.RawQuery(content).ExecWithCount(); err != nil {
				return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
			}
			return nil
		}
	}

	err := mb.findMigrations(txRunner, autoCommitRunner)
	if err != nil {
		return mb, err
	}

	for _, migration := range mb.goMigrations {
		mb.Migrations[migration.Direction] = append(mb.Migrations[migration.Direction], migration)
	}

	if err := mb.check(); err != nil {
		return nil, err
	}
	return mb, nil
}

func (fm *MigrationBox) findMigrations(
	runner func([]byte) func(mf Migration, c *pop.Connection, tx *pop.Tx) error,
	runnerNoTx func([]byte) func(mf Migration, c *pop.Connection) error,
) error {
	err := fs.WalkDir(fm.Dir, ".", func(p string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}

		match, err := ParseMigrationFilename(info.Name())
		if err != nil {
			if strings.HasPrefix(err.Error(), "unsupported dialect") {
				fm.l.Tracef("This is usually ok - ignoring migration file %s because dialect is not supported: %s", info.Name(), err.Error())
				return nil
			}
			return errors.WithStack(err)
		}

		if match == nil {
			fm.l.Tracef("This is usually ok - ignoring migration file %s because it does not match the file pattern.", info.Name())
			return nil
		}

		f, err := fm.Dir.Open(p)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}

		mf := Migration{
			Path:       p,
			Version:    match.Version,
			Name:       match.Name,
			DBType:     match.DBType,
			Direction:  match.Direction,
			Type:       match.Type,
			Content:    string(content),
			Autocommit: match.Autocommit,
		}

		if match.Autocommit {
			mf.RunnerNoTx = runnerNoTx(content)
		} else {
			mf.Runner = runner(content)
		}

		fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		return nil
	})

	// Sort descending.
	slices.SortFunc(fm.Migrations["down"], func(a, b Migration) int { return -CompareMigration(a, b) })

	// Sort ascending.
	slices.SortFunc(fm.Migrations["up"], CompareMigration)

	return err
}

// hasDownMigrationWithVersion checks if there is a migration with the given
// version.
func (fm *MigrationBox) hasDownMigrationWithVersion(version string) bool {
	for _, down := range fm.Migrations["down"] {
		if version == down.Version {
			return true
		}
	}
	return false
}

// check checks that every "up" migration has a corresponding "down" migration.
func (fm *MigrationBox) check() error {
	for _, up := range fm.Migrations["up"] {
		if !fm.hasDownMigrationWithVersion(up.Version) {
			return errors.Errorf("migration %s has no corresponding down migration", up.Version)
		}
	}

	for _, m := range fm.Migrations {
		for _, n := range m {
			if err := n.Valid(); err != nil {
				return err
			}
		}
	}
	return nil
}
