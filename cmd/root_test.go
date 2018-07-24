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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
)

var frontendPort, backendPort int

func init() {
	var err error
	frontendPort, err = freeport.GetFreePort()
	if err != nil {
		panic(err.Error())
	}

	backendPort, err = freeport.GetFreePort()
	if err != nil {
		panic(err.Error())
	}

	os.Setenv("FRONTEND_PORT", fmt.Sprintf("%d", frontendPort))
	os.Setenv("BACKEND_PORT", fmt.Sprintf("%d", backendPort))
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
			args: []string{"serve", "--disable-telemetry"},
			wait: func() bool {
				client := &http.Client{
					Transport: &transporter{
						FakeTLSTermination: true,
						Transport: &http.Transport{
							TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						},
					},
				}

				_, err := client.Get(fmt.Sprintf("https://127.0.0.1:%d/health/status", frontendPort))
				if err != nil {
					t.Logf("HTTP request failed: %s", err)
				} else {
					time.Sleep(time.Second * 5)
				}
				return err != nil
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
