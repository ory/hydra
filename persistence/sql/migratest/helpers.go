package migratest

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// TestMigrator is a modified pop.FileMigrator
type TestMigrator struct {
	pop.Migrator
}

func NewTestMigrator(t *testing.T, c *pop.Connection, migrationPath, testDataPath string) *TestMigrator {
	tm := TestMigrator{
		Migrator: pop.NewMigrator(c),
	}
	tm.SchemaPath = migrationPath
	testDataPath = strings.TrimSuffix(testDataPath, "/")

	runner := func(mf pop.Migration, tx *pop.Connection) error {
		f, err := os.Open(mf.Path)
		require.NoError(t, err)
		defer f.Close()
		content, err := pop.MigrationContent(mf, tx, f, true)
		require.NoError(t, err)
		if content == "" {
			return nil
		}
		err = tx.RawQuery(content).Exec()
		if err != nil {
			return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
		}

		if mf.Direction != "up" {
			return nil
		}

		// exec testdata
		var fileName string
		if fi, err := os.Stat(filepath.Join(testDataPath, mf.Version+"_testdata."+mf.DBType+".sql")); err == nil && !fi.IsDir() {
			// found specific test data
			fileName = fi.Name()
		} else if fi, err := os.Stat(filepath.Join(testDataPath, mf.Version+"_testdata.sql")); err == nil && !fi.IsDir() {
			// found generic test data
			fileName = fi.Name()
		} else {
			// found no test data
			log.Printf("Found no test data for migration %s", mf.Version)
			return nil
		}

		// Workaround for https://github.com/cockroachdb/cockroach/issues/42643#issuecomment-611475836
		// This is not a problem as the test should fail anyway if there occurs any error
		// (either within a transaction or on it's own).
		if mf.DBType == "cockroach" && tx.TX != nil {
			if err := tx.TX.Commit(); err != nil {
				return errors.WithStack(err)
			}
			newTx, err := c.NewTransaction()
			if err != nil {
				return errors.WithStack(err)
			}
			*tx = *newTx
		}

		data, err := ioutil.ReadFile(filepath.Join(testDataPath, fileName))
		if err != nil {
			return errors.WithStack(err)
		}

		if err := tx.RawQuery(string(data)).Exec(); err != nil {
			t.Logf(mf.Version)
			return errors.WithStack(err)
		}
		return nil
	}

	if fi, err := os.Stat(migrationPath); err != nil || !fi.IsDir() {
		t.Fatalf("could not find directory %s", migrationPath)
		return nil
	}
	if fi, err := os.Stat(testDataPath); err != nil || !fi.IsDir() {
		t.Fatalf("could not find directory %s", testDataPath)
		return nil
	}

	require.NoError(t, filepath.Walk(migrationPath, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			match, err := pop.ParseMigrationFilename(info.Name())
			if err != nil {
				return err
			}
			if match == nil {
				return nil
			}

			mf := pop.Migration{
				Path:      p,
				Version:   match.Version,
				Name:      match.Name,
				DBType:    match.DBType,
				Direction: match.Direction,
				Type:      match.Type,
				Runner:    runner,
			}
			tm.Migrations[mf.Direction] = append(tm.Migrations[mf.Direction], mf)
		}
		return nil
	}))

	return &tm
}
