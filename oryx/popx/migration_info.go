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

func (mfs Migrations) sortAndFilter(dialect string) Migrations {
	usable := make(map[string]Migration, len(mfs))
	for _, v := range mfs {
		if v.DBType == dialect {
			usable[v.Version] = v
		} else if v.DBType == "all" {
			// Add "all" only if we do not have a more specific migration for the dialect.
			// If a more specific migration is found later, it will override this one.
			if _, ok := usable[v.Version]; !ok {
				usable[v.Version] = v
			}
		}
	}

	filtered := make(Migrations, 0, len(usable))
	for k := range usable {
		filtered = append(filtered, usable[k])
	}
	sort.Sort(filtered)
	return filtered
}

func (mfs Migrations) find(version, dbType string) *Migration {
	var candidate *Migration
	for _, m := range mfs {
		if m.Version == version {
			switch m.DBType {
			case "all":
				// there might still be a more specific migration for the dbType
				candidate = &m
			case dbType:
				return &m
			}
		}
	}
	return candidate
}
