package migratest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/stringslice"
)

type TestMigrator struct {
	pop.Migrator
}

func NewTestMigrator(t *testing.T, c *pop.Connection, path string, steps int) *TestMigrator {
	tm := TestMigrator{
		Migrator: pop.NewMigrator(c),
	}
	tm.SchemaPath = path

	runner := func(mf pop.Migration, tx *pop.Connection) error {
		f, err := os.Open(mf.Path)
		require.NoError(t, err)
		defer f.Close()
		content, err := pop.MigrationContent(mf, tx, f, true)
		require.NoError(t, err)
		if content == "" {
			return nil
		}
		fmt.Printf("%s:\n\n%s\n\n\n", mf.Path, content)
		err = tx.RawQuery(content).Exec()
		if err != nil {
			return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
		}
		return nil
	}

	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		// directory doesn't exist
		return nil
	}

	uniqueVersions := map[string]int{}
	require.NoError(t, filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
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
			ms := append(tm.Migrations[mf.Direction], mf)

			var versions []string
			for _, m := range ms {
				versions = append(versions, m.Version)
			}
			uniqueVersions[mf.Direction] = len(stringslice.Unique(versions))
			// only add steps number of migrations
			if i, ok := uniqueVersions[match.Direction]; ok && i <= steps {
				tm.Migrations[mf.Direction] = ms
			}
			return nil
		}
		return nil
	}))

	return &tm
}
