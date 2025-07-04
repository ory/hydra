// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMigrationEmpty(t *testing.T) {
	assert.True(t, isMigrationEmpty(""))
	assert.True(t, isMigrationEmpty("-- this is a comment"))
	assert.True(t, isMigrationEmpty(`

-- this is a comment

`))
	assert.False(t, isMigrationEmpty(`SELECT foo`))
	assert.False(t, isMigrationEmpty(`INSERT bar -- test`))
	assert.False(t, isMigrationEmpty(`
--test
INSERT bar -- test

`))
}

func TestMigrationSort(t *testing.T) {

	migrations := []Migration{
		{Version: "99", DBType: "mysql"},
		{Version: "98", DBType: "mysql"},
		{Version: "99", DBType: "sqlite"},
		{Version: "99", DBType: "all"},
		{Version: "97", DBType: "mysql"},
		{Version: "99", DBType: "postgresql"},
		{Version: "97", DBType: ""},
		{Version: "99", DBType: ""},
	}

	slices.SortFunc(migrations, CompareMigration)

	expected := []Migration{
		{Version: "97", DBType: ""},
		{Version: "97", DBType: "mysql"},
		{Version: "98", DBType: "mysql"},
		{Version: "99", DBType: ""},
		{Version: "99", DBType: "mysql"},
		{Version: "99", DBType: "postgresql"},
		{Version: "99", DBType: "sqlite"},
		{Version: "99", DBType: "all"},
	}
	assert.Equal(t, expected, migrations)
}

func isLesserThan(a, b Migration) bool {
	return -1 == CompareMigration(a, b)
}

// `slices.SortFunc` requires that `cmp` is a strict weak ordering: (https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.)
// - Irreflexivity: For all x ∈ S , it is not true that x < x .
// - Transitivity: For all x , y , z ∈ S , if x < y  and  y < z then x < z .
// - Asymmetry: For all x , y ∈ S , if x < y is true then y < x is false.
// - (there is a fourth rule which does not apply to us).
func TestSortStrictWeakOrdering(t *testing.T) {
	m := Migrations{
		{Version: "0", DBType: "b"}, {Version: "0", DBType: "c"}, {Version: "0", DBType: "all"}, {Version: "1", DBType: "d"},
	}

	// Irreflexivity.
	for _, m := range migrations {
		assert.False(t, isLesserThan(m, m))
	}

	// Transitivity.
	// All 3-three_permutations.
	three_permutations := [][3]int{
		{0, 1, 2}, {0, 1, 3}, {0, 2, 1}, {0, 2, 3}, {0, 3, 1}, {0, 3, 2},
		{1, 0, 2}, {1, 0, 3}, {1, 2, 0}, {1, 2, 3}, {1, 3, 0}, {1, 3, 2},
		{2, 0, 1}, {2, 0, 3}, {2, 1, 0}, {2, 1, 3}, {2, 3, 0}, {2, 3, 1},
		{3, 0, 1}, {3, 0, 2}, {3, 1, 0}, {3, 1, 2}, {3, 2, 0}, {3, 2, 1},
	}

	for _, p := range three_permutations {
		x := m[p[0]]
		y := m[p[1]]
		z := m[p[2]]

		if isLesserThan(x, y) && isLesserThan(y, z) {
			assert.True(t, isLesserThan(x, z))
		}
	}

	// Asymmetry.
	// All 2-two_permutations.
	two_permutations := [][2]int{
		{0, 1}, {0, 2}, {0, 3},
		{1, 0}, {1, 2}, {1, 3},
		{2, 0}, {2, 1}, {2, 3},
		{3, 0}, {3, 1}, {3, 2},
	}

	for _, p := range two_permutations {
		x := m[p[0]]
		y := m[p[1]]

		if isLesserThan(x, y) {
			assert.False(t, isLesserThan(y, x))
		}
	}
}
