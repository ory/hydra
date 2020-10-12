// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package keys

import (
	"bytes"
	"crypto/tls"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/josex"
	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"net/http"
)

const (
	flagUse = "use"
)

// keysImportCmd represents the import command
var keysImportCmd = &cobra.Command{
	Use:   "import <set> <file-1> [<file-2> [<file-3 [<...>]]]",
	Short: "Imports cryptographic keys of any format to the JSON Web Key Store",
	Long: `This command allows you to import cryptographic keys to the JSON Web Key Store.

Currently supported formats are raw JSON Web Keys or PEM/DER encoded data. If the JSON Web Key Set exists already,
the imported keys will be added to that set. Otherwise, a new set will be created.

Please be aware that importing a private key does not automatically import its public key as well.

Examples:
	hydra keys import my-set ./path/to/jwk.json ./path/to/jwk-2.json
	hydra keys import my-set ./path/to/rsa.key ./path/to/rsa.pub
`,
	RunE: importKeys,
	Args: cobra.MinimumNArgs(2),
}

func init() {
	keysImportCmd.LocalFlags().String("use", "sig", "Sets the \"use\" value of the JSON Web Key if no \"use\" value was defined by the key itself")
}

func importKeys(cmd *cobra.Command, args []string) error {
	id := args[0]
	use := flagx.MustGetString(cmd, "use")
	client := &http.Client{}

	/* #nosec G402 - we want to support dev environments, hence tls trickery */
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: flagx.MustGetBool(cmd, "skip-tls-verify"),
		},
	}

	u := Remote(cmd) + "/keys/" + id
	request, err := http.NewRequest("GET", u, nil)
	cmdx.Must(err, "Unable to initialize HTTP request: %s", err)

	if flagx.MustGetBool(cmd, "fake-tls-termination") {
		request.Header.Set("X-Forwarded-Proto", "https")
	}

	if token := flagx.MustGetString(cmd, "access-token"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := client.Do(request)
	cmdx.Must(err, "Unable to fetch data from %s: %s", u, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		cmdx.Fatalf("Expected status code 200 or 404 but got %d while fetching data from %s", response.StatusCode, u)
	}

	var set jose.JSONWebKeySet
	err = json.NewDecoder(response.Body).Decode(&set)
	cmdx.Must(err, "Unable to decode payload to JSON: %s", err)

	for _, path := range args[1:] {
		file, err := ioutil.ReadFile(path)
		cmdx.Must(err, "Unable to read file %s", path)

		if key, privateErr := josex.LoadPrivateKey(file); privateErr != nil {
			key, publicErr := josex.LoadPublicKey(file)
			cmdx.Must(publicErr, `Unable to read key from file %s. Decoding file to private key failed with reason "%s" and decoding it to public key failed with reason: %s`, path, privateErr, publicErr)

			set.Keys = append(set.Keys, toSDKFriendlyJSONWebKey(key, "public:"+uuid.New(), use))
		} else {
			set.Keys = append(set.Keys, toSDKFriendlyJSONWebKey(key, "private:"+uuid.New(), use))
		}

		fmt.Printf("Successfully loaded key from file: %s\n", path)
	}

	body, err := json.Marshal(&set)
	cmdx.Must(err, "Unable to encode JSON Web Keys to JSON: %s", err)

	request, err = http.NewRequest("PUT", u, bytes.NewReader(body))
	cmdx.Must(err, "Unable to initialize HTTP request: %s", err)

	if flagx.MustGetBool(cmd, "fake-tls-termination") {
		request.Header.Set("X-Forwarded-Proto", "https")
	}

	if token := flagx.MustGetString(cmd, "access-token"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err = client.Do(request)
	cmdx.CheckResponse(err, http.StatusOK, response)
	defer response.Body.Close()

	fmt.Println("JSON Web Key Set successfully imported!")
}
