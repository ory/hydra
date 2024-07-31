// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/pointerx"
)

func NewImportClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oauth2-client <file-1.json> [<file-2.json> ...]",
		Aliases: []string{"client", "clients", "oauth2-clients"},
		Short:   "Import one or more OAuth 2.0 Clients from files or STDIN",
		Example: `Import an example OAuth2 Client:
	cat > ./file.json <<EOF
	[
      {
	    "grant_types": ["implicit"],
	    "scope": "openid"
	  },
      {
	    "grant_types": ["authorize_code"],
	    "scope": "openid"
	  }
    ]
	EOF

	{{ .CommandPath }} file.json

Alternatively:

	cat file.json | {{ .CommandPath }}

To encrypt an auto-generated OAuth2 Client Secret, use flags ` + "`--pgp-key`" + `, ` + "`--pgp-key-url`" + ` or ` + "`--keybase`" + ` flag, for example:

  {{ .CommandPath }} -n "my app" -g client_credentials -r token -a core,foobar --keybase keybase_username
`,
		Long: `This command reads in each listed JSON file and imports their contents as a list of OAuth 2.0 Clients.

The format for the JSON file is:

[
  {
    "client_secret": "...",
    // ... all other fields of the OAuth 2.0 Client model are allowed here
  }
]

Please be aware that this command does not update existing clients. If the client exists already, this command will fail.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			ek, encryptSecret, err := cli.NewEncryptionKey(cmd, nil)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to load encryption key: %s", err)
				return cmdx.FailSilently(cmd)
			}

			streams := map[string]io.Reader{"STDIN": cmd.InOrStdin()}
			for _, path := range args {
				contents, err := os.ReadFile(path)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not open file %s: %s", path, err)
					return cmdx.FailSilently(cmd)
				}
				streams[path] = bytes.NewReader(contents)
			}

			clients := map[string][]hydra.OAuth2Client{}
			for src, stream := range streams {
				var current []hydra.OAuth2Client
				if err := json.NewDecoder(stream).Decode(&current); err != nil {
					if errors.Is(err, io.EOF) {
						continue
					}
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not decode JSON: %s", err)
					return cmdx.FailSilently(cmd)
				}
				clients[src] = append(clients[src], current...)
			}

			imported := make([]hydra.OAuth2Client, 0, len(clients))
			failed := make(map[string]error)

			for src, cc := range clients {
				for _, c := range cc {
					result, _, err := m.OAuth2API.CreateOAuth2Client(cmd.Context()).OAuth2Client(c).Execute() //nolint:bodyclose
					if err != nil {
						failed[src] = cmdx.PrintOpenAPIError(cmd, err)
						continue
					}

					if result.ClientSecret == nil {
						result.ClientSecret = c.ClientSecret
					}

					if encryptSecret && result.ClientSecret != nil {
						enc, err := ek.Encrypt([]byte(*result.ClientSecret))
						if err != nil {
							_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to encrypt client secret: %s", err)
							return cmdx.FailSilently(cmd)
						}

						result.ClientSecret = pointerx.String(enc.Base64Encode())
					}

					imported = append(imported, *result)
				}
			}

			if len(imported) == 1 {
				cmdx.PrintRow(cmd, (*outputOAuth2Client)(&imported[0]))
			} else {
				cmdx.PrintTable(cmd, &outputOAuth2ClientCollection{clients: imported})
			}

			if len(failed) != 0 {
				cmdx.PrintErrors(cmd, failed)
				return cmdx.FailSilently(cmd)
			}

			return nil
		},
	}

	registerEncryptFlags(cmd.Flags())
	return cmd
}
