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

var port int

func init() {
	var err error
	port, err = freeport.GetFreePort()
	if err != nil {
		panic(err.Error())
	}
	os.Setenv("PORT", fmt.Sprintf("%d", port))
	os.Setenv("DATABASE_URL", "memory")
	os.Setenv("HYDRA_URL", fmt.Sprintf("https://localhost:%d/", port))
	os.Setenv("OAUTH2_ISSUER_URL", fmt.Sprintf("https://localhost:%d/", port))
}

func TestExecute(t *testing.T) {
	var osArgs = make([]string, len(os.Args))
	copy(osArgs, os.Args)

	endpoint := fmt.Sprintf("https://localhost:%d/", port)

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

				_, err := client.Get(fmt.Sprintf("https://127.0.0.1:%d/health/status", port))
				if err != nil {
					t.Logf("HTTP request failed: %s", err)
				} else {
					time.Sleep(time.Second * 5)
				}
				return err != nil
			},
		},
		{args: []string{"clients", "create", "--endpoint", endpoint, "--id", "foobarbaz", "--secret", "foobar", "-g", "client_credentials"}},
		{args: []string{"clients", "get", "--endpoint", endpoint, "foobarbaz"}},
		{args: []string{"clients", "create", "--endpoint", endpoint, "--id", "public-foo", "--is-public"}},
		{args: []string{"clients", "delete", "--endpoint", endpoint, "public-foo"}},
		{args: []string{"keys", "create", "foo", "--endpoint", endpoint, "-a", "HS256"}},
		{args: []string{"keys", "get", "--endpoint", endpoint, "foo"}},
		{args: []string{"keys", "rotate", "--endpoint", endpoint, "foo"}},
		{args: []string{"keys", "get", "--endpoint", endpoint, "foo"}},
		{args: []string{"keys", "delete", "--endpoint", endpoint, "foo"}},
		{args: []string{"token", "revoke", "--endpoint", endpoint, "--client-secret", "foobar", "--client-id", "foobarbaz", "foo"}},
		{args: []string{"token", "client", "--endpoint", endpoint, "--client-secret", "foobar", "--client-id", "foobarbaz"}},
		{args: []string{"help", "migrate", "sql"}},
		{args: []string{"version"}},
		{args: []string{"token", "flush", "--endpoint", endpoint}},
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
