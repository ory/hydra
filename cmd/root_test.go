/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package cmd

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ory/hydra/config"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
)

var frontendPort, backendPort int

func freePort() (int, int) {
	var err error
	r := make([]int, 2)

	if r[0], err = freeport.GetFreePort(); err != nil {
		panic(err.Error())
	}

	tries := 0
	for {
		r[1], err = freeport.GetFreePort()
		if r[0] != r[1] {
			break
		}
		tries++
		if tries > 10 {
			panic("Unable to find free port")
		}
	}
	return r[0], r[1]
}

func init() {
	frontendPort, backendPort = freePort()

	os.Setenv("PUBLIC_PORT", fmt.Sprintf("%d", frontendPort))
	os.Setenv("ADMIN_PORT", fmt.Sprintf("%d", backendPort))
	os.Setenv("DATABASE_URL", "memory")
	//os.Setenv("HYDRA_URL", fmt.Sprintf("https://localhost:%d/", frontendPort))
	os.Setenv("OAUTH2_ISSUER_URL", fmt.Sprintf("https://localhost:%d/", frontendPort))
}

func TestExecute(t *testing.T) {
	var osArgs = make([]string, len(os.Args))
	copy(osArgs, os.Args)

	frontend := fmt.Sprintf("https://localhost:%d/", frontendPort)
	backend := fmt.Sprintf("https://localhost:%d/", backendPort)

	for _, c := range []struct {
		args      []string
		wait      func() bool
		expectErr bool
	}{
		{
			args: []string{"serve", "all", "--disable-telemetry"},
			wait: func() bool {
				client := &http.Client{
					Transport: &transporter{
						FakeTLSTermination: true,
						Transport: &http.Transport{
							TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						},
					},
				}

				for _, u := range []string{
					fmt.Sprintf("https://127.0.0.1:%d/.well-known/openid-configuration", frontendPort),
					fmt.Sprintf("https://127.0.0.1:%d/health/status", backendPort),
				} {
					if resp, err := client.Get(u); err != nil {
						t.Logf("HTTP request to %s failed: %s", u, err)
						return true
					} else if resp.StatusCode != http.StatusOK {
						t.Logf("HTTP request to %s got status code %d but expected was 200", u, resp.StatusCode)
						return true
					}
				}

				// Give a bit more time to initialize
				time.Sleep(time.Second * 5)
				return false
			},
		},
		{args: []string{"clients", "create", "--endpoint", backend, "--id", "foobarbaz", "--secret", "foobar", "-g", "client_credentials"}},
		{args: []string{"clients", "get", "--endpoint", backend, "foobarbaz"}},
		{args: []string{"clients", "create", "--endpoint", backend, "--id", "public-foo"}},
		{args: []string{"clients", "delete", "--endpoint", backend, "public-foo"}},
		{args: []string{"keys", "create", "foo", "--endpoint", backend, "-a", "HS256"}},
		{args: []string{"keys", "get", "--endpoint", backend, "foo"}},
		{args: []string{"keys", "rotate", "--endpoint", backend, "foo"}},
		{args: []string{"keys", "get", "--endpoint", backend, "foo"}},
		{args: []string{"keys", "delete", "--endpoint", backend, "foo"}},
		{args: []string{"keys", "import", "--endpoint", backend, "import-1", "../test/stub/ecdh.key", "../test/stub/ecdh.pub"}},
		{args: []string{"keys", "import", "--endpoint", backend, "import-2", "../test/stub/rsa.key", "../test/stub/rsa.pub"}},
		{args: []string{"token", "revoke", "--endpoint", frontend, "--client-secret", "foobar", "--client-id", "foobarbaz", "foo"}},
		{args: []string{"token", "client", "--endpoint", frontend, "--client-secret", "foobar", "--client-id", "foobarbaz"}},
		{args: []string{"help", "migrate", "sql"}},
		{args: []string{"version"}},
		{args: []string{"token", "flush", "--endpoint", backend}},
	} {
		c.args = append(c.args, []string{"--skip-tls-verify"}...)
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
					t.Logf("Ports are not yet open, retrying attempt #%d...", count)
					count++
					if count > 15 {
						t.FailNow()
					}
					time.Sleep(time.Second)
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

func TestInitConfig(t *testing.T) {
	fd1, _ := ioutil.TempFile("", "test_system_secret.txt")
	fd2, _ := ioutil.TempFile("", "test_cookie_secret.txt")
	tss := fd1.Name()
	tcs := fd2.Name()
	defer func() {
		_ = os.Remove(tss)
		_ = os.Remove(tcs)
	}()

	expectedString1 := "ThisIsASystemSecretThatShouldNotBeSharedWithAnyone"
	expectedString2 := "ThisIsACookieSecretThatShouldNotBeSharedWithAnyone"
	expectedString3 := "ThisIsAnotherSystemSecretThatShouldNotBeSharedWithAnyone"
	expectedString4 := "ThisIsAnotherCookieSecretThatShouldNotBeSharedWithAnyone"

	_ = ioutil.WriteFile(tss, []byte(expectedString1+"\n"), 0644)
	_ = ioutil.WriteFile(tcs, []byte(expectedString2), 0644)

	os.Setenv("SYSTEM_SECRET_PATH", tss)
	os.Setenv("COOKIE_SECRET_PATH", tcs)

	c = new(config.Config)
	initConfig()

	assert.Equal(t, expectedString1, c.SystemSecret)
	assert.Equal(t, expectedString2, c.CookieSecret)

	assert.Len(t, c.GetSystemSecret(), 32)
	hash := sha256.Sum256([]byte(expectedString1))
	assert.Equal(t, hash[:], c.GetSystemSecret())

	hashCookie := sha256.Sum256([]byte(expectedString2))
	assert.Equal(t, hashCookie[:], c.GetCookieSecret())

	os.Setenv("SYSTEM_SECRET", expectedString3)

	c = new(config.Config)
	initConfig()

	assert.Equal(t, expectedString3, c.SystemSecret)
	assert.Equal(t, expectedString2, c.CookieSecret)

	assert.Len(t, c.GetSystemSecret(), 32)
	hash = sha256.Sum256([]byte(expectedString3))
	assert.Equal(t, hash[:], c.GetSystemSecret())
	hashCookie = sha256.Sum256([]byte(expectedString2))
	assert.Equal(t, hashCookie[:], c.GetCookieSecret())

	os.Unsetenv("SYSTEM_SECRET")
	os.Setenv("COOKIE_SECRET", expectedString4)

	c = new(config.Config)
	initConfig()

	assert.Equal(t, expectedString1, c.SystemSecret)
	assert.Equal(t, expectedString4[:], c.CookieSecret)

	assert.Len(t, c.GetSystemSecret(), 32)
	hash = sha256.Sum256([]byte(expectedString1))
	assert.Equal(t, hash[:], c.GetSystemSecret())
	hashCookie = sha256.Sum256([]byte(expectedString4))
	assert.Equal(t, hashCookie[:], c.GetCookieSecret())
}
