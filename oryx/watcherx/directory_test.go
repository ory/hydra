// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestWatchDirectory(t *testing.T) {
	t.Run("case=notifies about file creation in directory", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		_, err := WatchDirectory(ctx, dir, c)
		require.NoError(t, err)
		fileName := filepath.Join(dir, "example")
		f, err := os.Create(fileName) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		assertChange(t, <-c, "", fileName)
	})

	t.Run("case=notifies about file write in directory", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		fileName := filepath.Join(dir, "example")
		f, err := os.Create(fileName) //#nosec:G304
		require.NoError(t, err)
		_, err = WatchDirectory(ctx, dir, c)
		require.NoError(t, err)

		_, err = fmt.Fprintf(f, "content")
		require.NoError(t, err)
		require.NoError(t, f.Close())

		assertChange(t, <-c, "content", fileName)
	})

	t.Run("case=nofifies about file delete in directory", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		fileName := filepath.Join(dir, "example")
		f, err := os.Create(fileName) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		_, err = WatchDirectory(ctx, dir, c)
		require.NoError(t, err)
		require.NoError(t, os.Remove(fileName))

		assertRemove(t, <-c, fileName)
	})

	t.Run("case=notifies about file in child directory", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		childDir := filepath.Join(dir, "child")
		require.NoError(t, os.Mkdir(childDir, 0777))

		_, err := WatchDirectory(ctx, dir, c)
		require.NoError(t, err)

		fileName := filepath.Join(childDir, "example")
		f, err := os.Create(fileName) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		assertChange(t, <-c, "", fileName)
	})

	t.Run("case=watches new child directory", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		_, err := WatchDirectory(ctx, dir, c)
		require.NoError(t, err)

		childDir := filepath.Join(dir, "child")
		require.NoError(t, os.Mkdir(childDir, 0777))
		fileName := filepath.Join(childDir, "example")
		// there's not much we can do about this timeout as it takes some time until the new watcher is created
		time.Sleep(time.Millisecond)
		f, err := os.Create(fileName) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		assertChange(t, <-c, "", fileName)
	})

	t.Run("case=does not notify on directory deletion", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("skipping test because IN_DELETE_SELF is unreliable on windows and macOS")
		}

		ctx, c, dir, cancel := setup(t)
		defer cancel()

		childDir := filepath.Join(dir, "child")
		require.NoError(t, os.Mkdir(childDir, 0777))

		_, err := WatchDirectory(ctx, dir, c)
		require.NoError(t, err)

		require.NoError(t, os.Remove(childDir))

		select {
		case e := <-c:
			t.Logf("got unexpected event %T: %+v", e, e)
			t.FailNow()
		case <-time.After(2 * time.Millisecond):
			// expected to not receive an event (1ms is what the watcher waits for the second event)
		}
	})

	t.Run("case=notifies only for files on batch delete", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		childDir := filepath.Join(dir, "child")
		subChildDir := filepath.Join(childDir, "subchild")
		require.NoError(t, os.MkdirAll(subChildDir, 0777))
		f1 := filepath.Join(subChildDir, "f1")
		f, err := os.Create(f1) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())
		f2 := filepath.Join(childDir, "f2")
		f, err = os.Create(f2) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		_, err = WatchDirectory(ctx, dir, c)
		require.NoError(t, err)

		require.NoError(t, os.RemoveAll(childDir))

		events := []Event{<-c, <-c}
		if events[0].Source() > events[1].Source() {
			events[1], events[0] = events[0], events[1]
		}
		assertRemove(t, events[0], f2)
		assertRemove(t, events[1], f1)
	})

	t.Run("case=sends event when requested", func(t *testing.T) {
		ctx, _, dir, cancel := setup(t)
		defer cancel()

		// buffered channel to allow usage of DispatchNow().done
		c := make(EventChannel, 4)

		files := map[string]string{
			"a":                     "foo",
			"b":                     "bar",
			"c":                     "baz",
			filepath.Join("d", "a"): "sub dir content",
		}
		for fn, fc := range files {
			fp := filepath.Join(dir, fn)
			require.NoError(t, os.MkdirAll(filepath.Dir(fp), 0700))
			require.NoError(t, os.WriteFile(fp, []byte(fc), 0600))
		}

		d, err := WatchDirectory(ctx, dir, c)
		require.NoError(t, err)
		done, err := d.DispatchNow()
		require.NoError(t, err)

		// wait for d.DispatchNow to be done
		select {
		case <-time.After(time.Second):
			t.Log("Waiting for done timed out.")
			t.FailNow()
		case eventsSend := <-done:
			assert.Equal(t, 4, eventsSend)
		}

		// because filepath.WalkDir walks lexicographically, we can assume the events come in lex order
		assertChange(t, <-c, files["a"], filepath.Join(dir, "a"))
		assertChange(t, <-c, files["b"], filepath.Join(dir, "b"))
		assertChange(t, <-c, files["c"], filepath.Join(dir, "c"))
		assertChange(t, <-c, files[filepath.Join("d", "a")], filepath.Join(dir, "d", "a"))
	})
}
