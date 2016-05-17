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
			timeout: 10 * time.Second,
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
		RootCmd.SetArgs(c.args)

		t.Logf("Running command: %s", c.args)
		if c.timeout > 0 {
			go func(args []string) {
				assert.Nil(t, RootCmd.Execute())
			}(c.args)
			time.Sleep(c.timeout)
		} else {
			assert.Nil(t, RootCmd.Execute())
		}
	}
}
