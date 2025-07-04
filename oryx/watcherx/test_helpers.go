// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func KubernetesAtomicWrite(t *testing.T, dir, fileName, content string) {
	// atomic write according to https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/atomic_writer.go
	const (
		dataDirName    = "..data"
		newDataDirName = "..data_tmp"
	)
	// (2)
	dataDirPath := filepath.Join(dir, dataDirName)
	oldTsDir, err := os.Readlink(dataDirPath)
	if err != nil {
		require.True(t, os.IsNotExist(err), "%+v", err)
		// although Readlink() returns "" on err, don't be fragile by relying on it (since it's not specified in docs)
		// empty oldTsDir indicates that it didn't exist
		oldTsDir = ""
	}
	oldTsPath := filepath.Join(dir, oldTsDir)

	// (3) we are not interested in the case where a file gets deleted as we just operate on one file
	// (4) we assume the file needs an update

	// (5)
	tsDir, err := os.MkdirTemp(dir, time.Now().UTC().Format("..2006_01_02_15_04_05."))
	require.NoError(t, err)
	tsDirName := filepath.Base(tsDir)

	// (6)
	require.NoError(
		t,
		os.WriteFile(path.Join(tsDir, fileName), []byte(content), 0600),
	)

	// (7)
	_, err = os.Readlink(filepath.Join(dir, fileName))
	if err != nil && os.IsNotExist(err) {
		// The link into the data directory for this path doesn't exist; create it
		require.NoError(
			t,
			os.Symlink(filepath.Join(dataDirName, fileName), filepath.Join(dir, fileName)),
		)
	}

	// (8)
	newDataDirPath := filepath.Join(dir, newDataDirName)
	require.NoError(
		t,
		os.Symlink(tsDirName, newDataDirPath),
	)

	// (9)
	if runtime.GOOS == "windows" {
		require.NoError(t, os.Remove(dataDirPath))
		require.NoError(t, os.Symlink(tsDirName, dataDirPath))
		require.NoError(t, os.Remove(newDataDirPath))
	} else {
		require.NoError(t, os.Rename(newDataDirPath, dataDirPath))
	}

	// (10) in our case there is nothing to remove

	// (11)
	if len(oldTsDir) > 0 {
		require.NoError(t, os.RemoveAll(oldTsPath))
	}
}
