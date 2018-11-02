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

package cli

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func configureClient(cmd *cobra.Command, c *hydra.Configuration) *hydra.Configuration {
	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: flagx.MustGetBool(cmd, "skip-tls-verify")},
	}

	if flagx.MustGetBool(cmd, "fake-tls-termination") {
		c.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	if token := flagx.MustGetString(cmd, "access-token"); token != "" {
		c.DefaultHeader["Authorization"] = "Bearer " + token
	}
	return c
}

func configureClientWithoutAuth(cmd *cobra.Command, c *hydra.Configuration) *hydra.Configuration {
	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: flagx.MustGetBool(cmd, "skip-tls-verify")},
	}

	if flagx.MustGetBool(cmd, "fake-tls-termination") {
		c.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	return c
}

func checkResponse(err error, expectedStatusCode int, response *hydra.APIResponse) {
	r := new(http.Response)
	r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("Response is nil")))
	if response != nil && response.Response != nil {
		r = response.Response
		r.Body = ioutil.NopCloser(bytes.NewBuffer(response.Payload))
	}

	cmdx.CheckResponse(err, expectedStatusCode, r)
}

func formatResponse(response interface{}) string {
	out, err := json.MarshalIndent(response, "", "\t")
	cmdx.Must(err, `Command failed because an error ("%s") occurred while prettifying output`, err)
	return string(out)
}
