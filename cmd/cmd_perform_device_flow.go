// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ory/hydra/v2/cmd/cliclient"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/urlx"
)

func NewPerformDeviceCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device-code",
		Example: "{{ .CommandPath }} --client-id ...",
		Short:   "An exemplary OAuth 2.0 Client performing the OAuth 2.0 Device Code Flow",
		Long: `Performs the device code flow. Useful for getting an access token and an ID token in machines without a browser.
The client that will be used MUST use the "none" or "client_secret_post" token-endpoint-auth-method.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, endpoint, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			endpoint = cliclient.GetOAuth2URLOverride(cmd, endpoint)

			ctx := context.WithValue(cmd.Context(), oauth2.HTTPClient, client)
			scopes := flagx.MustGetStringSlice(cmd, "scope")
			deviceAuthUrl := flagx.MustGetString(cmd, "device-auth-url")
			tokenUrl := flagx.MustGetString(cmd, "token-url")
			audience := flagx.MustGetStringSlice(cmd, "audience")

			clientID := flagx.MustGetString(cmd, "client-id")
			if clientID == "" {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), cmd.UsageString())
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Please provide a Client ID using --client-id flag, or OAUTH2_CLIENT_ID environment variable.")
				return cmdx.FailSilently(cmd)
			}

			clientSecret := flagx.MustGetString(cmd, "client-secret")

			if deviceAuthUrl == "" {
				deviceAuthUrl = urlx.AppendPaths(endpoint, "/oauth2/device/auth").String()
			}

			if tokenUrl == "" {
				tokenUrl = urlx.AppendPaths(endpoint, "/oauth2/token").String()
			}

			conf := oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint: oauth2.Endpoint{
					DeviceAuthURL: deviceAuthUrl,
					TokenURL:      tokenUrl,
				},
				Scopes: scopes,
			}

			params := []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("audience", strings.Join(audience, "+"))}
			if clientSecret != "" {
				params = append(params, oauth2.SetAuthURLParam("client_secret", clientSecret))
			}

			deviceAuthResponse, err := conf.DeviceAuth(
				ctx,
				params...,
			)
			if err != nil {
				_, _ = fmt.Fprintf(
					cmd.ErrOrStderr(), "Failed to perform the device authorization request: %s\n", err)
				return cmdx.FailSilently(cmd)
			}

			_, _ = fmt.Fprintln(
				cmd.ErrOrStderr(),
				"To login please go to:\n\t",
				deviceAuthResponse.VerificationURIComplete,
			)

			token, err := conf.DeviceAccessToken(ctx, deviceAuthResponse)
			if err != nil {
				_, _ = fmt.Fprintf(
					cmd.ErrOrStderr(), "Failed to perform the device token request: %s\n", err)
				return cmdx.FailSilently(cmd)
			}

			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Successfully signed in!")

			cmdx.PrintRow(cmd, outputOAuth2Token(*token))
			return nil
		},
	}

	cmd.Flags().StringSlice("scope", []string{"offline", "openid"}, "Request OAuth2 scope")

	cmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	cmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")

	cmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience")
	cmd.Flags().String("device-auth-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the device authorization url, use this flag")
	cmd.Flags().String("token-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the token url, use this flag")

	return cmd
}
