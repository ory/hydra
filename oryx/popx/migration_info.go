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
	Runner func(Migration, *pop.Connection, *pop.Tx) error
	// RunnerNoTx function to run/execute the migration. NOT wrapped in a
	// database transaction. Mutually exclusive with Runner.
	RunnerNoTx func(Migration, *pop.Connection) error
	// Content is the raw content of the migration file
	Content string
	// Autocommit is true if the migration should be run outside of a transaction
	Autocommit bool
}

func (m Migration) Valid() error {
	if m.Runner == nil && m.RunnerNoTx == nil {
		return errors.Errorf("no runner defined for %s", m.Path)
	}
	if m.Runner != nil && m.RunnerNoTx != nil {
		return errors.Errorf("incompatible transaction and non-transaction runners defined for %s", m.Path)
	}
	return nil
}

// Migrations is a collection of Migration
type Migrations []Migration

func (mfs Migrations) Len() int {
	return len(mfs)
}

func (mfs Migrations) Less(i, j int) bool {
	return CompareMigration(mfs[i], mfs[j]) < 0
}

func CompareMigration(a, b Migration) int {
	if a.Version == b.Version {
		// Force "all" to be greater.
		if a.DBType == "all" && b.DBType != "all" {
			return 1
		} else if a.DBType != "all" && b.DBType == "all" {
			return -1
		} else {
			return strings.Compare(a.DBType, b.DBType)
		}
	}
	return strings.Compare(a.Version, b.Version)
}

func (mfs Migrations) Swap(i, j int) {
	mfs[i], mfs[j] = mfs[j], mfs[i]
}

func (mfs Migrations) SortAndFilter(dialect string, modifiers ...func(sort.Interface) sort.Interface) Migrations {
	// We need to sort mfs in order to push the dbType=="all" migrations
	// to the back.
	m := make(Migrations, len(mfs))
	copy(m, mfs)
	sort.Sort(m)

	vsf := make(Migrations, 0, len(m))
	for k, v := range m {
		if v.DBType == "all" {
			// Add "all" only if we can not find a more specific migration for the dialect.
			var hasSpecific bool
			for kk, vv := range m {
				if v.Version == vv.Version && kk != k && vv.DBType == dialect {
					hasSpecific = true
					break
				}
			}

			if !hasSpecific {
				vsf = append(vsf, v)
			}
		} else if v.DBType == dialect {
			vsf = append(vsf, v)
		}
	}

	mod := sort.Interface(vsf)
	for _, m := range modifiers {
		mod = m(mod)
	}

	sort.Sort(mod)
	return vsf
}
