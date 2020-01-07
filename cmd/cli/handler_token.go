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
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/client/public"
	"github.com/ory/hydra/internal/httpclient/models"

	"github.com/spf13/cobra"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

type TokenHandler struct{}

func newTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

func (h *TokenHandler) RevokeToken(cmd *cobra.Command, args []string) {
	cmdx.ExactArgs(cmd, args, 1)

	handler := configureClientWithoutAuth(cmd)

	clientID, clientSecret := flagx.MustGetString(cmd, "client-id"), flagx.MustGetString(cmd, "client-secret")
	if clientID == "" || clientSecret == "" {
		cmdx.Fatalf(`%s

Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET
`, cmd.UsageString())
	}

	token := args[0]
	_, err := handler.Public.RevokeOAuth2Token(public.NewRevokeOAuth2TokenParams().WithToken(args[0]), httptransport.BasicAuth(clientID, clientSecret))
	cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))

	fmt.Printf("Revoked OAuth 2.0 Access Token: %s\n", token)
}

func (h *TokenHandler) FlushTokens(cmd *cobra.Command, args []string) {
	handler := configureClient(cmd)
	_, err := handler.Admin.FlushInactiveOAuth2Tokens(admin.NewFlushInactiveOAuth2TokensParams().WithBody(&models.FlushInactiveOAuth2TokensRequest{
		NotAfter: strfmt.DateTime(time.Now().Add(-flagx.MustGetDuration(cmd, "min-age"))),
	}))
	cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
	fmt.Println("Successfully flushed inactive access tokens")
}
