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
	"encoding/json"
	"fmt"
	"os"

	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

func checkResponse(response *hydra.APIResponse, err error, expectedStatusCode int) {
	if response != nil {
		pkg.Must(err, "Command failed because calling \"%s %s\" resulted in error \"%s\" occurred.\n%s\n", response.Request.Method, response.RequestURL, err, response.Payload)
	} else {
		pkg.Must(err, "Command failed because error \"%s\" occurred and no response is available.\n", err)
	}

	if response.StatusCode != expectedStatusCode {
		fmt.Fprintf(os.Stderr, "Command failed because calling \"%s %s\" resulted in status code \"%d\" but code \"%d\" was expected.\n%s\n", response.Request.Method, response.RequestURL, expectedStatusCode, response.StatusCode, response.Payload)
		os.Exit(1)
		return
	}
}

func formatResponse(response interface{}) string {
	out, err := json.MarshalIndent(response, "", "\t")
	pkg.Must(err, `Command failed because an error ("%s") occurred while prettifying output.`, err)
	return string(out)
}
