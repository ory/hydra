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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mendsley/gojwk"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
)

type JWKHandler struct {
	Config *config.Config
}

func (h *JWKHandler) newJwkManager(cmd *cobra.Command) *hydra.JsonWebKeyApi {
	c := hydra.NewJsonWebKeyApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))

	skipTLSTermination, _ := cmd.Flags().GetBool("skip-tls-verify")
	c.Configuration.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSTermination},
	}

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	if token, _ := cmd.Flags().GetString("access-token"); token != "" {
		c.Configuration.DefaultHeader["Authorization"] = "Bearer " + token
	}

	return c
}

func newJWKHandler(c *config.Config) *JWKHandler {
	return &JWKHandler{Config: c}
}

func (h *JWKHandler) CreateKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 1 || len(args) > 2 {
		fmt.Println(cmd.UsageString())
		return
	}

	kid := ""
	if len(args) == 2 {
		kid = args[1]
	}

	alg, _ := cmd.Flags().GetString("alg")
	use, _ := cmd.Flags().GetString("use")
	keys, response, err := m.CreateJsonWebKeySet(args[0], hydra.JsonWebKeySetGeneratorRequest{Alg: alg, Kid: kid, Use: use})
	checkResponse(response, err, http.StatusCreated)
	fmt.Printf("%s\n", formatResponse(keys))
}

func toSDKFriendlyJSONWebKey(key interface{}, kid string, use string, public bool) jose.JSONWebKey {
	if jwk, ok := key.(*jose.JSONWebKey); ok {
		key = jwk.Key
		if jwk.KeyID != "" {
			kid = jwk.KeyID
		}
		if jwk.Use != "" {
			use = jwk.Use
		}
	}

	var err error
	var jwk *gojwk.Key
	if public {
		jwk, err = gojwk.PublicKey(key)
		pkg.Must(err, "Unable to convert public key to JSON Web Key because %s", err)
	} else {
		jwk, err = gojwk.PrivateKey(key)
		pkg.Must(err, "Unable to convert private key to JSON Web Key because %s", err)
	}

	return jose.JSONWebKey{
		KeyID:     kid,
		Use:       use,
		Algorithm: jwk.Alg,
		Key:       key,
	}
}

func (h *JWKHandler) ImportKeys(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println(cmd.UsageString())
		return
	}

	id := args[0]
	use, _ := cmd.Flags().GetString("use")
	client := &http.Client{}

	if skipTLSTermination, _ := cmd.Flags().GetBool("skip-tls-verify"); skipTLSTermination {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSTermination}}
	}

	u := h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd) + "/keys/" + id
	request, err := http.NewRequest("GET", u, nil)
	pkg.Must(err, "Unable to initialize HTTP request")

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		request.Header.Set("X-Forwarded-Proto", "https")
	}

	if token, _ := cmd.Flags().GetString("access-token"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := client.Do(request)
	pkg.Must(err, "Unable to fetch data from %s because %s", u, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		fmt.Printf("Expected status code 200 or 404 but got %d while fetching data from %s.\n", response.StatusCode, u)
		os.Exit(1)
	}

	var set jose.JSONWebKeySet
	pkg.Must(json.NewDecoder(response.Body).Decode(&set), "Unable to decode payload to JSON")

	for _, path := range args[1:] {
		file, err := ioutil.ReadFile(path)
		pkg.Must(err, "Unable to read file %s", path)

		if key, privateErr := pkg.LoadPrivateKey(file); privateErr != nil {
			key, publicErr := pkg.LoadPublicKey(file)
			if publicErr != nil {
				fmt.Printf("Unable to read key from file %s. Decoding file to private key failed with reason \"%s\" and decoding it to public key failed with reason \"%s\".\n", path, privateErr, publicErr)
				os.Exit(1)
			}

			set.Keys = append(set.Keys, toSDKFriendlyJSONWebKey(key, "public:"+uuid.New(), use, true))
		} else {
			set.Keys = append(set.Keys, toSDKFriendlyJSONWebKey(key, "private:"+uuid.New(), use, false))
		}

		fmt.Printf("Successfully loaded key from file %s\n", path)
	}

	body, err := json.Marshal(&set)
	pkg.Must(err, "Unable to encode JSON Web Keys to JSON")

	request, err = http.NewRequest("PUT", u, bytes.NewReader(body))
	pkg.Must(err, "Unable to initialize HTTP request")

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		request.Header.Set("X-Forwarded-Proto", "https")
	}

	if token, _ := cmd.Flags().GetString("access-token"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err = client.Do(request)
	pkg.Must(err, "Unable to post data to %s because %s", u, err)
	defer response.Body.Close()

	fmt.Println("Keys successfully imported!")
}

func (h *JWKHandler) GetKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) != 1 {
		fmt.Println(cmd.UsageString())
		return
	}

	keys, response, err := m.GetJsonWebKeySet(args[0])
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("%s\n", formatResponse(keys))
}

func (h *JWKHandler) DeleteKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) != 1 {
		fmt.Println(cmd.UsageString())
		return
	}

	response, err := m.DeleteJsonWebKeySet(args[0])
	checkResponse(response, err, http.StatusNoContent)
	fmt.Printf("Key set %s deleted.\n", args[0])
}
