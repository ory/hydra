// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/josex"
)

func NewKeysImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "jwk set-id file-1 [file-2] [file-n]",
		Args: cobra.MinimumNArgs(1),
		Example: `{{ .CommandPath }} my-set ./path/to/jwk.json ./path/to/jwk-2.json --format json
{{ .CommandPath }} my-set ./path/to/rsa.key ./path/to/rsa.pub --use enc`,
		Short: "Imports JSON Web Keys from one or more JSON files.",
		Long: `This command allows you to import JSON Web Keys from one or more JSON files or STDIN to the JSON Web Key Store.

Currently supported formats are raw JSON Web Keys or PEM/DER encoded data. If the JSON Web Key Set exists already,
the imported keys will be added to that set. Otherwise, a new set will be created.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			set := args[0]

			streams := map[string]io.Reader{}
			if len(args) == 1 {
				streams["STDIN"] = cmd.InOrStdin()
			} else {
				for _, path := range args[1:] {
					contents, err := os.ReadFile(path)
					if err != nil {
						_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not open file %s: %s", path, err)
						return cmdx.FailSilently(cmd)
					}
					streams[path] = bytes.NewReader(contents)
				}
			}

			keys := map[string][]hydra.JsonWebKey{}
			for src, stream := range streams {
				content, err := io.ReadAll(stream)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not read from %s: %s", src, err)
					return cmdx.FailSilently(cmd)
				}

				var key interface{}
				if priv, privErr := josex.LoadPrivateKey(content); privErr == nil {
					key = priv
				} else if pub, pubErr := josex.LoadPublicKey(content); pubErr == nil {
					key = pub
				} else {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not decode key from `%s` to public nor private keys: %s; %s", src, privErr, pubErr)
				}

				key = cli.ToSDKFriendlyJSONWebKey(key, "", "")

				type jwk hydra.JsonWebKey // opt out of OpenAPI-generated UnmarshalJSON
				var (
					buf        bytes.Buffer
					jsonWebKey jwk
				)
				if err := json.NewEncoder(&buf).Encode(key); err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not encode key from `%s` to JSON: %s", src, err)
					return cmdx.FailSilently(cmd)
				}

				if err := json.NewDecoder(&buf).Decode(&jsonWebKey); err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not decode key from `%s` from JSON: %s", src, err)
					return cmdx.FailSilently(cmd)
				}

				if len(jsonWebKey.Kid) == 0 {
					jsonWebKey.Kid = uuid.Must(uuid.NewV4()).String()
				}

				if len(jsonWebKey.Alg) == 0 {
					jsonWebKey.Alg = flagx.MustGetString(cmd, "alg")
				}

				if len(jsonWebKey.Alg) == 0 {
					_, _ = fmt.Fprint(cmd.ErrOrStderr(), "Flag `--alg` is required when imported key does not define the `alg` field itself.")
					return cmdx.FailSilently(cmd)
				}

				if len(jsonWebKey.Use) == 0 {
					jsonWebKey.Use = flagx.MustGetString(cmd, "use")
				}

				if len(jsonWebKey.Use) == 0 {
					_, _ = fmt.Fprint(cmd.ErrOrStderr(), "Flag `--use` is required when imported key does not define the `use` field itself.")
					return cmdx.FailSilently(cmd)
				}

				keys[src] = append(keys[src], hydra.JsonWebKey(jsonWebKey))
			}

			imported := make([]hydra.JsonWebKey, 0, len(keys))
			failed := make(map[string]error)
			for src, kk := range keys {
				for _, k := range kk {
					result, _, err := m.JwkAPI.SetJsonWebKey(cmd.Context(), set, k.Kid).JsonWebKey(k).Execute() //nolint:bodyclose
					if err != nil {
						failed[src] = cmdx.PrintOpenAPIError(cmd, err)
						continue
					}

					imported = append(imported, *result)
				}
			}

			cmdx.PrintTable(cmd, &outputJSONWebKeyCollection{Set: set, Keys: imported})
			if len(failed) != 0 {
				cmdx.PrintErrors(cmd, failed)
				return cmdx.FailSilently(cmd)
			}

			return nil
		},
	}

	cmd.Flags().String("use", "sig", "Sets the \"use\" value of the JSON Web Key if no \"use\" value was defined by the key itself. Required when importing PEM/DER encoded data.")
	cmd.Flags().String("alg", "", "Sets the \"alg\" value of the JSON Web Key if not \"alg\" value was defined by the key itself. Required when importing PEM/DER encoded data.")
	return cmd
}
