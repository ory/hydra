// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fsx

import (
	"crypto/sha512"
	"io"
	"io/fs"
)

// DirHash computes a directory hash from all files contained in any subdirectories.
func DirHash(dir fs.FS) ([]byte, error) {
	hash := sha512.New()
	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		_, _ = io.WriteString(hash, path) // hash write never errors
		f, err := dir.Open(path)
		if err != nil {
			return err
		}
		_, _ = io.Copy(hash, f) // hash write never errors
		if err = f.Close(); err != nil {
			return err
		}
		return nil
	})
	return hash.Sum(nil), err
}
