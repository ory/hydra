// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/pointerx"
)

func NewUpdateClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oauth2-client [id]",
		Aliases: []string{"client"},
		Short:   "Update an OAuth 2.0 Client",
		Args:    cobra.ExactArgs(1),
		Example: `{{ .CommandPath }} <client-id-here> --redirect-uri http://localhost/cb --grant-type authorization_code --response-type code --scope core,foobar

To encrypt an auto-generated OAuth2 Client Secret, use flags ` + "`--pgp-key`" + `, ` + "`--pgp-key-url`" + ` or ` + "`--keybase`" + ` flag, for example:

  {{ .CommandPath }} e6e96aa5-9cd2-4a70-bf56-ad6434c8aaa2 --name "my app" --grant-type client_credentials --response-type token --scope core,foobar --keybase keybase_username
`,
		Long: `This command replaces an OAuth 2.0 Client by its ID. Please be aware that this command replaces the entire client. If only the name flag (-n "my updated app") is provided, the all other fields are updated to their default values.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			secret := flagx.MustGetString(cmd, "secret")
			ek, encryptSecret, err := cli.NewEncryptionKey(cmd, nil)
			if err != nil {
				return err
			}

			id := args[0]
			cc, err := clientFromFlags(cmd)
			if err != nil {
				return err
			}

			client, _, err := m.OAuth2API.SetOAuth2Client(context.Background(), id).OAuth2Client(cc).Execute() //nolint:bodyclose
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			if client.ClientSecret == nil && len(secret) > 0 {
				client.ClientSecret = pointerx.Ptr(secret)
			}

			if encryptSecret && client.ClientSecret != nil {
				enc, err := ek.Encrypt([]byte(*client.ClientSecret))
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to encrypt client secret: %s", err)
					return cmdx.FailSilently(cmd)
				}

				client.ClientSecret = pointerx.Ptr(enc.Base64Encode())
			}

			cmdx.PrintRow(cmd, (*outputOAuth2Client)(client))
			return nil
		},
	}

	registerClientFlags(cmd.Flags())

	return cmd
}
