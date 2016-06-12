package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/pborman/uuid"
)

func init() {
	c.BindPort = 13124
}

func TestExecute(t *testing.T) {
	var osArgs = make([]string, len(os.Args))
	var path = filepath.Join(os.TempDir(), fmt.Sprintf("hydra-%s.yml", uuid.New()))
	copy(osArgs, os.Args)

	for _, c := range []struct {
		args    []string
		timeout time.Duration
		wait    func() bool
	}{
		{
			args:    []string{"host", "--dangerous-auto-logon"},
			wait: func() bool {
				_, err := os.Stat(path)
				return err != nil
			},
		},
		{
			args:    []string{"token", "user", "--no-open"},
			timeout: time.Second,
		},
		{args: []string{"clients", "create", "--id", "foobarbaz"}},
		{args: []string{"clients", "delete", "foobarbaz"}},
		{args: []string{"keys", "create", "foo", "-a", "RS256"}},
		{args: []string{"keys", "create", "foo", "-a", "ES521"}},
		{args: []string{"keys", "get", "foo"}},
		{args: []string{"keys", "delete", "foo"}},
		{args: []string{"connections", "create", "google", "localuser", "googleuser"}},
		{args: []string{"token", "client"}},
		{args: []string{"policies", "create", "-i", "foobar", "-s", "peter", "max", "-r", "blog", "users", "-a", "post", "ban", "--allow"}},
		{args: []string{"policies", "get", "foobar"}},
		{args: []string{"policies", "delete", "foobar"}},
	} {
		c.args = append(c.args, []string{"--skip-tls-verify", "--config", path}...)
		RootCmd.SetArgs(c.args)

		t.Logf("Running command: %s", c.args)
		if c.wait != nil || c.timeout > 0 {
			go func() {
				assert.Nil(t, RootCmd.Execute())
			}()
		}

		if c.wait != nil {
			for c.wait() {
				time.Sleep(time.Millisecond * 500)
			}
		} else if c.timeout > 0 {
			time.Sleep(c.timeout)
		} else {
			assert.Nil(t, RootCmd.Execute())
		}
	}
}
