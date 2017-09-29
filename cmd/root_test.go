package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	c.BindPort = 13124
}

func TestExecute(t *testing.T) {
	var osArgs = make([]string, len(os.Args))
	var path = filepath.Join(os.TempDir(), fmt.Sprintf("hydra-%s.yml", uuid.New()))
	os.Setenv("DATABASE_URL", "memory")
	os.Setenv("FORCE_ROOT_CLIENT_CREDENTIALS", "admin:pw")
	os.Setenv("ISSUER", "https://localhost:4444/")
	copy(osArgs, os.Args)

	for _, c := range []struct {
		args      []string
		wait      func() bool
		expectErr bool
	}{
		{
			args: []string{"host", "--dangerous-auto-logon", "--disable-telemetry"},
			wait: func() bool {
				_, err := os.Stat(path)
				if err != nil {
					t.Logf("Could not stat path %s because %s", path, err)
				} else {
					time.Sleep(time.Second * 5)
				}
				return err != nil
			},
		},
		{args: []string{"connect", "--id", "admin", "--secret", "pw", "--url", "https://127.0.0.1:4444"}},
		{args: []string{"clients", "create", "--id", "foobarbaz"}},
		{args: []string{"clients", "get", "foobarbaz"}},
		{args: []string{"clients", "create", "--id", "public-foo", "--is-public"}},
		{args: []string{"clients", "delete", "foobarbaz"}},
		{args: []string{"keys", "create", "foo", "-a", "HS256"}},
		{args: []string{"keys", "create", "foo", "-a", "HS256"}},
		{args: []string{"keys", "get", "foo"}},
		{args: []string{"keys", "delete", "foo"}},
		{args: []string{"token", "revoke", "foo"}},
		{args: []string{"token", "client"}},
		{args: []string{"policies", "create", "-i", "foobar", "-s", "peter,max", "-r", "blog,users", "-a", "post,ban", "--allow"}},
		{args: []string{"policies", "actions", "add", "foobar", "update|create"}},
		{args: []string{"policies", "actions", "remove", "foobar", "update|create"}},
		{args: []string{"policies", "resources", "add", "foobar", "printer"}},
		{args: []string{"policies", "resources", "remove", "foobar", "printer"}},
		{args: []string{"policies", "subjects", "add", "foobar", "ken", "tracy"}},
		{args: []string{"policies", "subjects", "remove", "foobar", "ken", "tracy"}},
		{args: []string{"policies", "get", "foobar"}},
		{args: []string{"policies", "delete", "foobar"}},
		{args: []string{"groups", "create", "my-group"}},
		{args: []string{"groups", "members", "add", "my-group", "peter"}},
		{args: []string{"groups", "find", "peter"}},
		{args: []string{"groups", "members", "remove", "my-group", "peter"}},
		{args: []string{"groups", "delete", "my-group"}},
		{args: []string{"help", "migrate", "sql"}},
		{args: []string{"help", "migrate", "ladon", "0.6.0"}},
		{args: []string{"version"}},
		{args: []string{"token", "user", "--no-open"}, wait: func() bool {
			time.Sleep(time.Millisecond * 10)
			return false
		}},
	} {
		c.args = append(c.args, []string{"--skip-tls-verify", "--config", path}...)
		RootCmd.SetArgs(c.args)

		t.Run(fmt.Sprintf("command=%v", c.args), func(t *testing.T) {
			if c.wait != nil {
				go func() {
					assert.Nil(t, RootCmd.Execute())
				}()
			}

			if c.wait != nil {
				var count = 0
				for c.wait() {
					t.Logf("Config file has not been found yet, retrying attempt #%d...", count)
					count++
					if count > 200 {
						t.FailNow()
					}
					time.Sleep(time.Second * 2)
				}
			} else {
				err := RootCmd.Execute()
				if c.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
