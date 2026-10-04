// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"
)

// Migration handles the data for a given database migration
type Migration struct {
	// Path to the migration (./migrations/123_create_widgets.up.sql)
	Path string
	// Version of the migration (123)
	Version string
	// Name of the migration (create_widgets)
	Name string
	// Direction of the migration (up|down)
	Direction string
	// Type of migration (sql|go)
	Type string
	// DB type (all|postgres|mysql...)
	DBType string
	// Runner function to run/execute the migration. Will be wrapped in a
	// database transaction. Mutually exclusive with RunnerNoTx
	Runner func(Migration, *pop.Connection) error
	// Content is the raw content of the migration file
	Content string
	// Autocommit indicates whether the migration should be run in autocommit mode
	Autocommit bool
}

func (m Migration) Valid() error {
	if m.Runner == nil {
		return errors.Errorf("no runner defined for %s", m.Path)
	}

	return nil
}

// Migrations is a collection of Migration
type Migrations []Migration

func (mfs Migrations) Len() int           { return len(mfs) }
func (mfs Migrations) Less(i, j int) bool { return compareMigration(mfs[i], mfs[j]) < 0 }
func (mfs Migrations) Swap(i, j int)      { mfs[i], mfs[j] = mfs[j], mfs[i] }

func compareMigration(a, b Migration) int {
	if a.Version != b.Version {
		return strings.Compare(a.Version, b.Version)
	}
	// Force "all" to be greater.
	if a.DBType == "all" && b.DBType != "all" {
		return 1
	} else if a.DBType != "all" && b.DBType == "all" {
		return -1
	}
	return strings.Compare(a.DBType, b.DBType)
}

// migrationRank ranks a migration's DBType for the given dialect and its ordered
// fallbacks. Higher numbers indicate a better match:
//
//	-1      not applicable
//	 0      "all" (generic, dialect-independent)
//	 1..n   a fallback dialect (earlier fallbacks rank higher)
//	 n+1    the exact dialect
//
// This lets a migration box prefer a dialect-specific override, then a fallback
// dialect (e.g. postgres for a postgres-wire-compatible database), then the
// generic file.
func migrationRank(dbType, dialect string, fallbacks []string) int {
	switch dbType {
	case dialect:
		return len(fallbacks) + 1
	case "all":
		return 0
	}
	for i, fb := range fallbacks {
		if dbType == fb {
			return len(fallbacks) - i
		}
	}
	return -1
}

func (mfs Migrations) sortAndFilter(dialect string, fallbacks ...string) Migrations {
	byVersion := make(map[string]Migration, len(mfs))
	for _, migration := range mfs {
		rank := migrationRank(migration.DBType, dialect, fallbacks)
		if rank < 0 {
			// Not applicable to this dialect.
			continue
		}
		if previousMigration, ok := byVersion[migration.Version]; ok {
			if rank < migrationRank(previousMigration.DBType, dialect, fallbacks) {
				// Previous migration is a better match, skip this one.
				continue
			}
		}
		byVersion[migration.Version] = migration
	}

	filtered := make(Migrations, 0, len(byVersion))
	for k := range byVersion {
		filtered = append(filtered, byVersion[k])
	}
	sort.Sort(filtered)
	return filtered
}

func (mfs Migrations) find(version, dialect string, fallbacks ...string) *Migration {
	var candidate *Migration
	bestRank := -1
	for i := range mfs {
		if mfs[i].Version != version {
			continue
		}
		if rank := migrationRank(mfs[i].DBType, dialect, fallbacks); rank > bestRank {
			bestRank = rank
			candidate = &mfs[i]
		}
	}
	return candidate
}
