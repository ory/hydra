package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	c.BindPort = 13124
}

func TestExecute(t *testing.T) {
	var osArgs = make([]string, len(os.Args))
	copy(osArgs, os.Args)

	for _, c := range []struct {
		args    []string
		timeout time.Duration
	}{
		{
			args:    []string{"host", "--dangerous-auto-logon"},
			timeout: time.Second,
		},
		{args: []string{"clients", "create"}},
		{args: []string{"keys", "create", "foo", "-a", "RS256"}},
		{args: []string{"keys", "create", "foo", "-a", "ES521"}},
		{args: []string{"keys", "get", "foo"}},
		{args: []string{"keys", "delete", "foo"}},
		{args: []string{"policies", "create", "-i", "foobar", "-s", "peter", "max", "-r", "blog", "users", "-a", "post", "ban", "--allow"}},
		{args: []string{"policies", "get", "foobar"}},
		{args: []string{"policies", "delete", "foobar"}},
	} {
		c.args = append(c.args, []string{"--skip-tls-verify", "--config", filepath.Join(os.TempDir(), "hydra.yml")}...)

		if c.timeout > 0 {
			t.Logf("Running async command: %s", c.args)
			go func(args []string) {
				RootCmd.SetArgs(args)
				assert.Nil(t, RootCmd.Execute())
			}(c.args)
			time.Sleep(c.timeout)
			continue
		}
		t.Logf("Running command %s", c.args)

		RootCmd.SetArgs(c.args)
		assert.Nil(t, RootCmd.Execute())
	}
}
