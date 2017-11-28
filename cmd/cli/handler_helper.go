// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

func checkResponse(response *hydra.APIResponse, err error, expectedStatusCode int) {
	pkg.Must(err, "Command failed because error \"%s\" occurred.\n", err)

	if response.StatusCode != expectedStatusCode {
		fmt.Fprintf(os.Stderr, "Command failed because status code %d was expeceted but code %d was received.\n", expectedStatusCode, response.StatusCode)
		os.Exit(1)
		return
	}
}

func formatResponse(response interface{}) string {
	out, err := json.MarshalIndent(response, "", "\t")
	pkg.Must(err, `Command failed because an error ("%s") occurred while prettifying output.`, err)
	return string(out)
}
