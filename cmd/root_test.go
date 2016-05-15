package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"time"
	"path/filepath"
)

func TestExecute(t *testing.T) {
	var osArgs = make([]string, len(os.Args))
	copy(osArgs, os.Args)

	for _, c := range []struct {
		args    []string
		timeout time.Duration
	}{
		{
			args:[]string{"host"},
			timeout: time.Second,
		},
		{
			args:[]string{"clients", "create"},
		},
		{
			args:[]string{"keys", "create", "foo", "-a", "RS256"},
		},
		{
			args:[]string{"keys", "create", "foo", "-a", "ES521"},
		},
		{
			args:[]string{"keys", "get", "foo"},
		},
		{
			args:[]string{"keys", "delete", "foo"},
		},
	} {
		c.args = append(c.args, []string{"--skip-ca-check", "--config", filepath.Join(os.TempDir(), "hydra.yml")}...)

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