// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"
)

var MigrationFileRegexp = regexp.MustCompile(
	`^(\d+)_([^.]+)(\.[a-z0-9]+)?(\.autocommit)?\.(up|down)\.(sql)$`,
)

const (
	// Human-readable constants for the regex capture groups
	versionIdx = iota + 1
	nameIdx
	dbTypeIdx
	autocommitIdx
	directionIdx
	typeIdx
)

// match holds the information parsed from a migration filename.
type match struct {
	Version    string
	Name       string
	DBType     string
	Direction  string
	Type       string
	Autocommit bool
}

// parseMigrationFilename parses a migration filename.
func parseMigrationFilename(filename string) (*match, error) {
	matches := MigrationFileRegexp.FindAllStringSubmatch(filename, -1)
	if len(matches) == 0 {
		return nil, nil
	}
	m := matches[0]

	var (
		autocommit bool
		dbType     string
	)

	if m[dbTypeIdx] == ".autocommit" {
		// A special case where autocommit group moves forward to the 3rd index.
		autocommit = true
		dbType = "all"
	} else if m[dbTypeIdx] == "" {
		dbType = "all"
	} else {
		dbType = pop.CanonicalDialect(m[dbTypeIdx][1:])
		if !pop.DialectSupported(dbType) {
			return nil, errors.Errorf("unsupported dialect %s", dbType)
		}
	}

	if m[typeIdx] == "fizz" && dbType != "all" {
		return nil, errors.Errorf("invalid database type %q, expected \"all\" because fizz is database type independent", dbType)
	}

	if m[autocommitIdx] == ".autocommit" {
		autocommit = true
	} else if m[autocommitIdx] != "" {
		return nil, errors.Errorf("invalid autocommit flag %q", m[autocommitIdx])
	}

	return &match{
		Version:    m[versionIdx],
		Name:       m[nameIdx],
		DBType:     dbType,
		Autocommit: autocommit,
		Direction:  m[directionIdx],
		Type:       m[typeIdx],
	}, nil
}
