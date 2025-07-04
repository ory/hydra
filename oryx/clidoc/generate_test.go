package clidoc

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func noopRun(_ *cobra.Command, _ []string) {}

var (
	root = &cobra.Command{Use: "root", Run: noopRun, Long: `A sample text
root

<[some argument]>
`}
	child1 = &cobra.Command{Use: "child1", Run: noopRun, Long: `A sample text
child1

<[some argument]>
`, Example: "{{ .CommandPath }} --whatever"}
	child2 = &cobra.Command{Use: "child2", Run: noopRun, Long: `A sample text
child2

<[some argument]>
`}
	subChild1 = &cobra.Command{Use: "subChild1 <args>", Run: noopRun, Long: `A sample text
subChild1

<[some argument]>
`}
)

func snapshotDir(t *testing.T, path ...string) (assertNoChange func(t *testing.T)) {
	var (
		as  []func(*testing.T)
		fps []string
	)

	require.NoError(t, filepath.WalkDir(filepath.Join(path...), func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err, path)
		if !d.IsDir() {
			fps = append(fps, path)
			as = append(as, snapshotFile(t, path))
		}
		return nil
	}))

	return func(t *testing.T) {
		fileN := 0
		require.NoError(t, filepath.WalkDir(filepath.Join(path...), func(path string, d fs.DirEntry, err error) error {
			require.NoError(t, err)
			if !d.IsDir() {
				assert.Contains(t, fps, path)
				fileN++
			}
			return nil
		}))
		assert.Equal(t, len(fps), fileN)

		for _, a := range as {
			a(t)
		}
	}
}

func snapshotFile(t *testing.T, path ...string) (assertNoChange func(t *testing.T)) {
	pre, err := os.ReadFile(filepath.Join(path...))
	require.NoError(t, err)
	pre = bytes.ReplaceAll(pre, []byte("\r\n"), []byte("\n"))

	return func(t *testing.T) {
		post, err := os.ReadFile(filepath.Join(path...))
		require.NoError(t, err)

		assert.Equal(t, string(pre), string(post), "%s", post)
	}
}

func init() {
	child1.AddCommand(subChild1)
	root.AddCommand(child1, child2)
}

func TestGenerate(t *testing.T) {
	assertNoChange := snapshotDir(t, "testdata")
	require.NoError(t, Generate(root, []string{"testdata"}))
	assertNoChange(t)
}
