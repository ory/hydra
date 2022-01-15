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
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

type Handler struct {
	Clients       *ClientHandler
	Keys          *JWKHandler
	Introspection *IntrospectionHandler
	Token         *TokenHandler
	Migration     *MigrateHandler
	Janitor       *JanitorHandler
}

func Remote(cmd *cobra.Command) string {
	if endpoint := flagx.MustGetString(cmd, "endpoint"); endpoint != "" {
		return strings.TrimRight(endpoint, "/")
	} else if endpoint := os.Getenv("HYDRA_URL"); endpoint != "" {
		return strings.TrimRight(endpoint, "/")
	}

	cmdx.Fatalf("To execute this command, the endpoint URL must point to the URL where Ory Hydra is located. To set the endpoint URL, use flag --endpoint or environment variable HYDRA_URL if an administrative command was used.")
	return ""
}

func RemoteURI(cmd *cobra.Command) *url.URL {
	endpoint, err := url.ParseRequestURI(Remote(cmd))
	cmdx.Must(err, "Unable to parse remote url: %s", err)
	return endpoint
}

func NewHandler() *Handler {
	return &Handler{
		Clients:       newClientHandler(),
		Keys:          newJWKHandler(),
		Introspection: newIntrospectionHandler(),
		Token:         newTokenHandler(),
		Migration:     newMigrateHandler(),
		Janitor:       NewJanitorHandler(),
	}
}
