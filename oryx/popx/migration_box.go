// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"

	"github.com/ory/x/logrusx"
)

type (
	// MigrationBox is a embed migration box.
	MigrationBox struct {
		c                   *pop.Connection
		migrationsUp        Migrations
		migrationsDown      Migrations
		perMigrationTimeout time.Duration
		dumpMigrations      bool
		l                   *logrusx.Logger
		migrationContent    MigrationContent
	}
	MigrationContent   func(mf Migration, c *pop.Connection, r []byte, usingTemplate bool) (string, error)
	MigrationBoxOption func(*MigrationBox)
)

func WithTemplateValues(v map[string]interface{}) MigrationBoxOption {
	return func(m *MigrationBox) {
		m.migrationContent = ParameterizedMigrationContent(v)
	}
}

func WithMigrationContentMiddleware(middleware func(content string, err error) (string, error)) MigrationBoxOption {
	return func(m *MigrationBox) {
		prev := m.migrationContent
		m.migrationContent = func(mf Migration, c *pop.Connection, r []byte, usingTemplate bool) (string, error) {
			return middleware(prev(mf, c, r, usingTemplate))
		}
	}
}

// WithGoMigrations adds migrations that have a custom migration runner.
// TEST THEM THOROUGHLY!
// It will be very hard to fix a buggy migration.
func WithGoMigrations(migrations Migrations) MigrationBoxOption {
	return func(mb *MigrationBox) {
		for _, m := range migrations {
			switch m.Direction {
			case "up":
				mb.migrationsUp = append(mb.migrationsUp, m)
			case "down":
				mb.migrationsDown = append(mb.migrationsDown, m)
			default:
				panic(fmt.Sprintf("unknown migration direction %q for %q", m.Direction, m.Version))
			}
		}
	}
}

func WithPerMigrationTimeout(timeout time.Duration) MigrationBoxOption {
	return func(m *MigrationBox) {
		m.perMigrationTimeout = timeout
	}
}

func WithDumpMigrations() MigrationBoxOption {
	return func(m *MigrationBox) {
		m.dumpMigrations = true
	}
}

// WithTestdata adds testdata to the migration box.
func WithTestdata(t *testing.T, testdata fs.FS) MigrationBoxOption {
	testdataPattern := regexp.MustCompile(`^(\d+)_testdata(|\.[a-zA-Z0-9]+).sql$`)
	return func(m *MigrationBox) {
		require.NoError(t, fs.WalkDir(testdata, ".", func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return errors.WithStack(err)
			}
			if !info.Type().IsRegular() {
				t.Logf("skipping testdata entry that is not a file: %s", path)
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

			m.migrationsUp = append(m.migrationsUp, Migration{
				Version:   version + "9", // run testdata after version
				Path:      path,
				Name:      info.Name(),
				DBType:    flavor,
				Direction: "up",
				Type:      "sql",
				Runner: func(m Migration, c *pop.Connection) error {
					b, err := fs.ReadFile(testdata, m.Path)
					if err != nil {
						return errors.WithStack(err)
					}
					if isMigrationEmpty(string(b)) {
						return nil
					}
					_, err = c.Store.SQLDB().Exec(string(b))
					return errors.WithStack(err)
				},
			})

			m.migrationsDown = append(m.migrationsDown, Migration{
				Version:   version + "9", // run testdata after version
				Path:      path,
				Name:      info.Name(),
				DBType:    flavor,
				Direction: "down",
				Type:      "sql",
				Runner:    func(m Migration, _ *pop.Connection) error { return nil },
			})

			return nil
		}))
	}
}

var emptySQLReplace = regexp.MustCompile(`(?m)^(\s*--.*|\s*)$`)

func isMigrationEmpty(content string) bool {
	return len(strings.ReplaceAll(emptySQLReplace.ReplaceAllString(content, ""), "\n", "")) == 0
}

type queryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// NewMigrationBox creates a new migration box.
func NewMigrationBox(dir fs.FS, c *pop.Connection, l *logrusx.Logger, opts ...MigrationBoxOption) (*MigrationBox, error) {
	mb := &MigrationBox{
		c:                c,
		l:                l,
		migrationContent: ParameterizedMigrationContent(nil),
	}

	for _, o := range opts {
		o(mb)
	}

	txRunner := func(b []byte) func(Migration, *pop.Connection) error {
		return func(mf Migration, c *pop.Connection) error {
			content, err := mb.migrationContent(mf, c, b, true)
			if err != nil {
				return errors.Wrapf(err, "error processing %s", mf.Path)
			}
			if isMigrationEmpty(content) {
				l.WithField("migration", mf.Path).Trace("This is usually ok - ignoring migration because content is empty. This is ok!")
				return nil
			}

			var q queryExecutor = c.Store.SQLDB()
			if c.TX != nil {
				q = c.TX
			}

			if _, err = q.Exec(content); err != nil {
				return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
			}
			return nil
		}
	}

	err := mb.findMigrations(dir, txRunner)
	if err != nil {
		return mb, err
	}

	if err := mb.check(); err != nil {
		return nil, err
	}
	return mb, nil
}

func (mb *MigrationBox) findMigrations(
	dir fs.FS,
	runner func([]byte) func(m Migration, c *pop.Connection) error,
) error {
	err := fs.WalkDir(dir, ".", func(p string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if !info.Type().IsRegular() {
			mb.l.Tracef("ignoring non file: %s", info.Name())
			return nil
		}

		if path.Ext(info.Name()) != ".sql" {
			mb.l.Tracef("ignoring non SQL file: %s", info.Name())
			return nil
		}

		details, err := parseMigrationFilename(info.Name())
		if err != nil {
			if strings.HasPrefix(err.Error(), "unsupported dialect") {
				mb.l.Tracef("This is usually ok - ignoring migration file %s because dialect is not supported: %s", info.Name(), err.Error())
				return nil
			}
			return errors.WithStack(err)
		}

		if details == nil {
			return errors.Errorf("Found a migration file that does not match the file pattern: filename=%s pattern=%s", info.Name(), MigrationFileRegexp)
		}

		content, err := fs.ReadFile(dir, p)
		if err != nil {
			return errors.WithStack(err)
		}

		mf := Migration{
			Path:       p,
			Version:    details.Version,
			Name:       details.Name,
			DBType:     details.DBType,
			Direction:  details.Direction,
			Type:       details.Type,
			Content:    string(content),
			Autocommit: details.Autocommit,
		}

		mf.Runner = runner(content)

		switch details.Direction {
		case "up":
			mb.migrationsUp = append(mb.migrationsUp, mf)
		case "down":
			mb.migrationsDown = append(mb.migrationsDown, mf)
		default:
			return errors.Errorf("unknown migration direction %q for %q", details.Direction, info.Name())
		}
		return nil
	})

	// Sort descending.
	sort.Sort(mb.migrationsDown)
	slices.Reverse(mb.migrationsDown)

	// Sort ascending.
	sort.Sort(mb.migrationsUp)

	return errors.WithStack(err)
}

// hasDownMigrationWithVersion checks if there is a migration with the given
// version.
func (mb *MigrationBox) hasDownMigrationWithVersion(version string) bool {
	for _, down := range mb.migrationsDown {
		if version == down.Version {
			return true
		}
	}
	return false
}

// check checks that every "up" migration has a corresponding "down" migration.
func (mb *MigrationBox) check() error {
	for _, up := range mb.migrationsUp {
		if !mb.hasDownMigrationWithVersion(up.Version) {
			return errors.Errorf("migration %s has no corresponding down migration", up.Version)
		}
	}

	for _, n := range mb.migrationsUp {
		if err := n.Valid(); err != nil {
			return errors.WithStack(err)
		}
	}
	for _, n := range mb.migrationsDown {
		if err := n.Valid(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
