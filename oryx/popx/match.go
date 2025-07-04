// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"fmt"
	"regexp"

	"github.com/ory/pop/v6"
)

var mrx = regexp.MustCompile(
	`^(\d+)_([^.]+)(\.[a-z0-9]+)?(\.autocommit)?\.(up|down)\.(sql)$`,
)

// Match holds the information parsed from a migration filename.
type Match struct {
	Version    string
	Name       string
	DBType     string
	Direction  string
	Type       string
	Autocommit bool
}

// ParseMigrationFilename parses a migration filename.
func ParseMigrationFilename(filename string) (*Match, error) {
	matches := mrx.FindAllStringSubmatch(filename, -1)
	if len(matches) == 0 {
		return nil, nil
	}
	m := matches[0]

	var autocommit bool
	var dbType string
	if m[3] == ".autocommit" {
		// A special case where autocommit group moves forward to the 3rd index.
		autocommit = true
		dbType = "all"
	} else if m[3] == "" {
		dbType = "all"
	} else {
		dbType = pop.CanonicalDialect(m[3][1:])
		if !pop.DialectSupported(dbType) {
			return nil, fmt.Errorf("unsupported dialect %s", dbType)
		}
	}

	if m[6] == "fizz" && dbType != "all" {
		return nil, fmt.Errorf("invalid database type %q, expected \"all\" because fizz is database type independent", dbType)
	}

	if m[4] == ".autocommit" {
		autocommit = true
	} else if m[4] != "" {
		return nil, fmt.Errorf("invalid autocommit flag %q", m[4])
	}

	match := &Match{
		Version:    m[1],
		Name:       m[2],
		DBType:     dbType,
		Autocommit: autocommit,
		Direction:  m[5],
		Type:       m[6],
	}

	return match, nil
}
