// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"
	"github.com/ory/x/fsx"
	"github.com/ory/x/sqlcon/dockertest"
)

var hashDumpRegex = regexp.MustCompile(`-- migrations hash: ([^\n]+)\n`)

func RestoreFromSchemaDump(t testing.TB, c *pop.Connection, migrations fs.FS, writeTo string) func(testing.TB) {
	newHash, err := fsx.DirHash(migrations)
	require.NoError(t, err)

	dumpFilename := filepath.Join(writeTo, c.Dialect.Name()+"_dump.sql")

	updateDump := func(t testing.TB) {
		dump := dockertest.DumpSchema(t, c)
		f, err := os.OpenFile(dumpFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err)
		defer f.Close()

		_, _ = fmt.Fprintf(f, "-- migrations hash: %x\n\n%s", newHash, dump)
		t.Fatal("database schema restored from migrations and dump updated, please re-run the test")
	}

	dump, err := os.ReadFile(dumpFilename)
	if errors.Is(err, fs.ErrNotExist) {
		return updateDump
	}
	require.NoError(t, err)

	matches := hashDumpRegex.FindSubmatch(dump)
	if len(matches) != 2 {
		return updateDump
	}

	currentHash, err := hex.DecodeString(string(matches[1]))
	require.NoError(t, err)

	if !bytes.Equal(newHash, currentHash) {
		return updateDump
	}

	require.NoError(t, c.RawQuery(string(dump)).Exec())
	return nil
}
