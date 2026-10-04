// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package testingx contains helper functions and extensions used when writing tests in Ory.
package testingx

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/stretchr/testify/require"
)

// RepoRootPath returns the absolute path of the closest parent directory that has a go.mod file relative to the caller.
func RepoRootPath(t require.TestingT) (repoRoot string) {
	_, fpath, _, _ := runtime.Caller(1)
	for dir := filepath.Dir(filepath.FromSlash(fpath)); dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			repoRoot = dir
			break
		}
	}
	require.NotEmptyf(t, repoRoot, "could not determine repo root using path: %q", fpath)
	return repoRoot
}
