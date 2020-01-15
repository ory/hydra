package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/spf13/viper"
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
		{args: []string{"connect", "--id", "admin", "--secret", "pw", "--url", "https://127.0.0.1:4444/"}},
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

var secrets map[string]string
var errorManagerErr error

type FakeManager struct {
}

func (m FakeManager) GetSecrets(name string) (map[string]string, []byte, error) {
	return secrets, nil, nil
}

type ErrorManager struct {
}

func (m ErrorManager) GetSecrets(name string) (map[string]string, []byte, error) {
	return nil, nil, errorManagerErr
}

func TestSetupAppCerts(t *testing.T) {
	sManager = FakeManager{}
	os.Setenv("AWS_APP_CERTS_SECRET_NAME", "some")
	secrets = map[string]string{
		pkg.AppSSLCert: "hello",
		pkg.AppSSLKey:  "there",
	}
	setupAppCerts()
	assert.Equal(t, "hello", os.Getenv("HTTPS_TLS_CERT"))
	assert.Equal(t, "there", os.Getenv("HTTPS_TLS_KEY"))

	os.Unsetenv("AWS_APP_CERTS_SECRET_NAME")
	os.Unsetenv("HTTPS_TLS_CERT")
	os.Unsetenv("HTTPS_TLS_KEY")
	viper.Reset()
}

func TestSetupAppCertsWithNoSecretName(t *testing.T) {
	sManager = FakeManager{}
	os.Unsetenv("AWS_APP_CERTS_SECRET_NAME")
	assert.Nil(t, setupAppCerts())
}

func TestSetupAppCertsWithSecretError(t *testing.T) {
	sManager = ErrorManager{}
	os.Setenv("AWS_APP_CERTS_SECRET_NAME", "some")
	errorManagerErr = errors.New("hi")

	err := setupAppCerts()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Error getting app certs from Secrets Manager: hi")

	os.Unsetenv("AWS_APP_CERTS_SECRET_NAME")
}

func TestSetupAppCertsWithNoSSLData(t *testing.T) {
	sManager = FakeManager{}
	os.Setenv("AWS_APP_CERTS_SECRET_NAME", "some")
	secrets = map[string]string{
		pkg.AppSSLKey: "there",
	}
	err := setupAppCerts()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("App certificate (%s) on Secrets Manager (%s) not found", pkg.AppSSLCert, "some"))

	secrets = map[string]string{
		pkg.AppSSLCert: "hello",
	}
	err = setupAppCerts()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("App key (%s) on Secrets Manager (%s) not found", pkg.AppSSLKey, "some"))

	os.Unsetenv("AWS_APP_CERTS_SECRET_NAME")
}

func TestSetupDBCerts(t *testing.T) {
	sManager = FakeManager{}
	os.Setenv("AWS_RDS_CERTS_SECRET_NAME", "some")
	secrets = map[string]string{
		pkg.RdsSSLCert: "how",
	}
	setupDBCerts()
	assert.Equal(t, "how", viper.Get("RdsSSLCert").(string))

	os.Unsetenv("AWS_RDS_CERTS_SECRET_NAME")
	viper.Reset()
}

func TestSetupDBCertsWithNoSecret(t *testing.T) {
	sManager = FakeManager{}
	os.Unsetenv("AWS_RDS_CERTS_SECRET_NAME")
	assert.Nil(t, setupDBCerts())
}

func TestSetupDBCertsWithSecretError(t *testing.T) {
	sManager = ErrorManager{}
	os.Setenv("AWS_RDS_CERTS_SECRET_NAME", "some")
	errorManagerErr = errors.New("hi")

	err := setupDBCerts()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Error getting rds certs from Secrets Manager: hi")

	os.Unsetenv("AWS_RDS_CERTS_SECRET_NAME")
}

func TestSetupDBCertsWithNoSSLData(t *testing.T) {
	sManager = FakeManager{}
	os.Setenv("AWS_RDS_CERTS_SECRET_NAME", "some")
	secrets = map[string]string{}

	err := setupDBCerts()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("RDS certificate (%s) on Secrets Manager (%s) not found", pkg.RdsSSLCert, "some"))

	os.Unsetenv("AWS_APP_CERTS_SECRET_NAME")
}
