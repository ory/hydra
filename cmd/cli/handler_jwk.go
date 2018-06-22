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
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"os"

	"github.com/ory/hydra/config"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type JWKHandler struct {
	Config *config.Config
}

func (h *JWKHandler) newJwkManager(cmd *cobra.Command) *hydra.JsonWebKeyApi {
	c := hydra.NewJsonWebKeyApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlash(cmd))

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

func (h *JWKHandler) RotateKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 1 || len(args) > 2 {
		fmt.Println(cmd.UsageString())
		return
	}

	setID := args[0]
	kid := ""
	if len(args) == 2 {
		kid = args[1]
	}

	set, response, err := m.GetJsonWebKeySet(setID)
	checkResponse(response, err, http.StatusOK)

	var toCreate = map[string]hydra.JsonWebKeySetGeneratorRequest{}
	var found bool

	var f = func(s string) string {
		return strings.Replace(strings.Replace(s, "public:", "", -1), "private:", "", -1)
	}

	shouldDelete, _ := cmd.Flags().GetBool("delete")

	if len(kid) == 0 {
		for _, key := range set.Keys {
			found = true

			if shouldDelete {
				response, err := m.DeleteJsonWebKey(key.Kid, setID)
				checkResponse(response, err, http.StatusNoContent)
			}

			k := f(key.Kid)
			toCreate[k] = hydra.JsonWebKeySetGeneratorRequest{Use: key.Use, Alg: key.Alg}
		}

		for _, k := range toCreate {
			_, response, err = m.CreateJsonWebKeySet(setID, k)
			checkResponse(response, err, http.StatusCreated)
		}
	} else if len(kid) > 0 {
		var tc hydra.JsonWebKeySetGeneratorRequest
		for _, key := range set.Keys {
			if f(kid) == f(key.Kid) {
				found = true

				if shouldDelete {
					response, err := m.DeleteJsonWebKey(key.Kid, setID)
					checkResponse(response, err, http.StatusNoContent)
				}

				tc = hydra.JsonWebKeySetGeneratorRequest{Alg: key.Alg, Use: key.Use}
			}
		}

		if found {
			_, response, err = m.CreateJsonWebKeySet(setID, tc)
			checkResponse(response, err, http.StatusCreated)
		}
	}

	if !found {
		if kid == "" {
			fmt.Fprintln(os.Stderr, "The JSON Web Key Set does not contain any keys, thus no keys could be rotated.")
		} else {
			fmt.Fprintf(os.Stderr, "The JSON Web Key Set does not contain key with kid \"%s\" keys, thus the key could be rotated.\n", kid)
		}
		os.Exit(1)
	}
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
