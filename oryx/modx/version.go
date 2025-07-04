// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package modx

import (
	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

// FindVersion returns the version for a module given the contents of a go.mod file.
func FindVersion(gomod []byte, module string) (string, error) {
	m, err := modfile.Parse("go.mod", gomod, nil)
	if err != nil {
		return "", err
	}

	for _, r := range m.Require {
		if r.Mod.Path == module {
			return r.Mod.Version, nil
		}
	}

	return "", errors.Errorf("no go.mod entry found for: %s", module)
}

// MustFindVersion returns the version for a module given the contents of a go.mod file or panics.
func MustFindVersion(gomod []byte, module string) string {
	v, err := FindVersion(gomod, module)
	if err != nil {
		panic(err)
	}
	return v
}
