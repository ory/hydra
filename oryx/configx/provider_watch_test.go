// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/watcherx"
)

func tmpConfigFile(t *testing.T, dsn, foo string) (string, string) {
	config := fmt.Sprintf("dsn: %s\nfoo: %s\n", dsn, foo)

	tdir := t.TempDir()
	fn := "config.yml"
	watcherx.KubernetesAtomicWrite(t, tdir, fn, config)

	return tdir, fn
}

func updateConfigFile(t *testing.T, c <-chan struct{}, dir, name, dsn, foo, bar string) {
	config := fmt.Sprintf(`dsn: %s
foo: %s
bar: %s`, dsn, foo, bar)

	watcherx.KubernetesAtomicWrite(t, dir, name, config)
	<-c // Wait for changes to propagate
	time.Sleep(time.Millisecond)
}

func assertNoOpenFDs(t require.TestingT, dir, name string) {
	if runtime.GOOS == "windows" {
		return
	}
	var b, be bytes.Buffer
	// we are only interested in the file descriptors, so we use the `-F f` option
	c := exec.Command("lsof", "-n", "-F", "f", "--", filepath.Join(dir, name))
	c.Stdout = &b
	c.Stderr = &be
	exitErr := new(exec.ExitError)
	require.ErrorAsf(t, c.Run(), &exitErr, "File %q has open file descriptor.\nGot stout: %s\nstderr: %s", filepath.Join(dir, name), b.String(), be.String())
	assert.Equal(t, 1, exitErr.ExitCode(), "got stout: %s\nstderr: %s", b.String(), be.String())
}

func TestReload(t *testing.T) {
	setup := func(t *testing.T, dir, name string, c chan<- struct{}, modifiers ...OptionModifier) (*Provider, *logrusx.Logger) {
		l := logrusx.New("configx", "test")
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		modifiers = append(modifiers,
			WithLogrusWatcher(l),
			WithLogger(l),
			AttachWatcher(func(event watcherx.Event, err error) {
				fmt.Printf("Received event: %+v error: %+v\n", event, err)
				c <- struct{}{}
			}),
			WithContext(ctx),
		)
		p, err := newKoanf(ctx, "./stub/watch/config.schema.json", []string{filepath.Join(dir, name)}, modifiers...)
		require.NoError(t, err)
		return p, l
	}

	t.Run("case=rejects not validating changes", func(t *testing.T) {
		t.Parallel()
		dir, name := tmpConfigFile(t, "memory", "bar")
		c := make(chan struct{})
		p, l := setup(t, dir, name, c)
		hook := test.NewLocal(l.Entry.Logger)

		assertNoOpenFDs(t, dir, name)

		assert.Equal(t, []*logrus.Entry{}, hook.AllEntries())
		assert.Equal(t, "memory", p.String("dsn"))
		assert.Equal(t, "bar", p.String("foo"))

		updateConfigFile(t, c, dir, name, "memory", "not bar", "bar")

		entries := hook.AllEntries()
		require.False(t, len(entries) > 4, "%+v", entries) // should be 2 but addresses flake https://github.com/ory/x/runs/2332130952

		assert.Equal(t, "A change to a configuration file was detected.", entries[0].Message)
		assert.Equal(t, "The changed configuration is invalid and could not be loaded. Rolling back to the last working configuration revision. Please address the validation errors before restarting the process.", entries[1].Message)

		assert.Equal(t, "memory", p.String("dsn"))
		assert.Equal(t, "bar", p.String("foo"))

		// but it is still watching the files
		updateConfigFile(t, c, dir, name, "memory", "bar", "baz")
		assert.Equal(t, "baz", p.String("bar"))

		time.Sleep(time.Millisecond * 250)

		assertNoOpenFDs(t, dir, name)
	})

	t.Run("case=rejects to update immutable", func(t *testing.T) {
		t.Parallel()
		dir, name := tmpConfigFile(t, "memory", "bar")
		c := make(chan struct{})
		p, l := setup(t, dir, name, c,
			WithImmutables("dsn"))
		hook := test.NewLocal(l.Entry.Logger)

		assertNoOpenFDs(t, dir, name)

		assert.Equal(t, []*logrus.Entry{}, hook.AllEntries())
		assert.Equal(t, "memory", p.String("dsn"))
		assert.Equal(t, "bar", p.String("foo"))

		updateConfigFile(t, c, dir, name, "some db", "bar", "baz")

		entries := hook.AllEntries()
		require.False(t, len(entries) > 4, "%+v", entries) // should be 2 but addresses flake https://github.com/ory/x/runs/2332130952
		assert.Equal(t, "A change to a configuration file was detected.", entries[0].Message)
		assert.Equal(t, "A configuration value marked as immutable has changed. Rolling back to the last working configuration revision. To reload the values please restart the process.", entries[1].Message)
		assert.Equal(t, "memory", p.String("dsn"))
		assert.Equal(t, "bar", p.String("foo"))

		// but it is still watching the files
		updateConfigFile(t, c, dir, name, "memory", "bar", "baz")
		assert.Equal(t, "baz", p.String("bar"))

		assertNoOpenFDs(t, dir, name)
	})

	t.Run("case=allows to update excepted immutable", func(t *testing.T) {
		t.Parallel()
		config := `{"foo": {"bar": "a", "baz": "b"}}`

		dir := t.TempDir()
		name := "config.json"
		watcherx.KubernetesAtomicWrite(t, dir, name, config)

		c := make(chan struct{})
		p, _ := setup(t, dir, name, c,
			WithImmutables("foo"),
			WithExceptImmutables("foo.baz"),
			SkipValidation())

		assert.Equal(t, "a", p.String("foo.bar"))
		assert.Equal(t, "b", p.String("foo.baz"))

		config = `{"foo": {"bar": "a", "baz": "x"}}`
		watcherx.KubernetesAtomicWrite(t, dir, name, config)
		<-c
		time.Sleep(time.Millisecond)

		assert.Equal(t, "x", p.String("foo.baz"))
	})

	t.Run("case=runs without validation errors", func(t *testing.T) {
		t.Parallel()
		dir, name := tmpConfigFile(t, "some string", "bar")
		c := make(chan struct{})
		p, l := setup(t, dir, name, c)
		hook := test.NewLocal(l.Entry.Logger)

		assert.Equal(t, []*logrus.Entry{}, hook.AllEntries())
		assert.Equal(t, "some string", p.String("dsn"))
		assert.Equal(t, "bar", p.String("foo"))
	})

	t.Run("case=runs and reloads", func(t *testing.T) {
		t.Parallel()
		dir, name := tmpConfigFile(t, "some string", "bar")
		c := make(chan struct{})
		p, l := setup(t, dir, name, c)
		hook := test.NewLocal(l.Entry.Logger)

		assert.Equal(t, []*logrus.Entry{}, hook.AllEntries())
		assert.Equal(t, "some string", p.String("dsn"))
		assert.Equal(t, "bar", p.String("foo"))

		updateConfigFile(t, c, dir, name, "memory", "bar", "baz")
		assert.Equal(t, "baz", p.String("bar"))
	})

	t.Run("case=has with validation errors", func(t *testing.T) {
		t.Parallel()
		dir, name := tmpConfigFile(t, "some string", "not bar")
		l := logrusx.New("", "")
		hook := test.NewLocal(l.Entry.Logger)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var b bytes.Buffer
		_, err := newKoanf(ctx, "./stub/watch/config.schema.json", []string{filepath.Join(dir, name)},
			WithStandardValidationReporter(&b),
			WithLogrusWatcher(l),
		)
		require.Error(t, err)

		entries := hook.AllEntries()
		require.Equal(t, 0, len(entries))
		assert.Equal(t, "The configuration contains values or keys which are invalid:\nfoo: not bar\n     ^-- value must be \"bar\"\n\n", b.String())
	})

	t.Run("case=is not leaking open files", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip()
		}

		dir, name := tmpConfigFile(t, "some string", "bar")
		c := make(chan struct{})
		p, _ := setup(t, dir, name, c)

		assertNoOpenFDs(t, dir, name)

		for i := range 30 {
			t.Run(fmt.Sprintf("iteration=%d", i), func(t *testing.T) {
				expected := []string{"foo", "bar", "baz"}[i%3]
				updateConfigFile(t, c, dir, name, "memory", "bar", expected)
				assertNoOpenFDs(t, dir, name)
				require.EqualValues(t, expected, p.String("bar"))
			})
		}

		assertNoOpenFDs(t, dir, name)
	})

	t.Run("case=callback can use the provider to get the new value", func(t *testing.T) {
		t.Parallel()
		dsn := "old"

		dir, name := tmpConfigFile(t, dsn, "bar")
		c := make(chan struct{})

		var p *Provider
		p, _ = setup(t, dir, name, c, AttachWatcher(func(watcherx.Event, error) {
			dsn = p.String("dsn")
		}))

		// change dsn
		updateConfigFile(t, c, dir, name, "new", "bar", "bar")

		assert.Equal(t, "new", dsn)
	})
}

type mockTestingT struct {
	failed bool
}

func (m *mockTestingT) FailNow() {
	m.failed = true
}

func (m *mockTestingT) Errorf(string, ...interface{}) {}

var _ require.TestingT = (*mockTestingT)(nil)

func TestAssertNoOpenFDs(t *testing.T) {
	t.Parallel()

	mt := &mockTestingT{}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "foo"))
	require.NoError(t, err)

	assertNoOpenFDs(mt, dir, "foo")
	assert.True(t, mt.failed)

	mt = &mockTestingT{}
	require.NoError(t, f.Close())
	assertNoOpenFDs(mt, dir, "foo")
	assert.False(t, mt.failed)
}
