// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/x/flagx"
	"github.com/ory/x/pointerx"
)

func clientFromFlags(cmd *cobra.Command) (hydra.OAuth2Client, error) {
	if filename := flagx.MustGetString(cmd, flagFile); filename != "" {
		src := cmd.InOrStdin()
		if filename != "-" {
			f, err := os.Open(filename) // #nosec G304
			if err != nil {
				return hydra.OAuth2Client{}, fmt.Errorf("unable to open file %q: %w", filename, err)
			}
			defer f.Close() //nolint:errcheck
			src = f
		}
		client := hydra.OAuth2Client{}
		if err := json.NewDecoder(src).Decode(&client); err != nil {
			return hydra.OAuth2Client{}, fmt.Errorf("unable to decode JSON: %w", err)
		}
		return client, nil
	}

	return hydra.OAuth2Client{
		AccessTokenStrategy:               pointerx.Ptr(flagx.MustGetString(cmd, flagClientAccessTokenStrategy)),
		AllowedCorsOrigins:                flagx.MustGetStringSlice(cmd, flagClientAllowedCORSOrigin),
		Audience:                          flagx.MustGetStringSlice(cmd, flagClientAudience),
		BackchannelLogoutSessionRequired:  pointerx.Ptr(flagx.MustGetBool(cmd, flagClientBackChannelLogoutSessionRequired)),
		BackchannelLogoutUri:              pointerx.Ptr(flagx.MustGetString(cmd, flagClientBackchannelLogoutCallback)),
		ClientName:                        pointerx.Ptr(flagx.MustGetString(cmd, flagClientName)),
		ClientSecret:                      pointerx.Ptr(flagx.MustGetString(cmd, flagClientSecret)),
		ClientUri:                         pointerx.Ptr(flagx.MustGetString(cmd, flagClientClientURI)),
		Contacts:                          flagx.MustGetStringSlice(cmd, flagClientContact),
		FrontchannelLogoutSessionRequired: pointerx.Ptr(flagx.MustGetBool(cmd, flagClientFrontChannelLogoutSessionRequired)),
		FrontchannelLogoutUri:             pointerx.Ptr(flagx.MustGetString(cmd, flagClientFrontChannelLogoutCallback)),
		GrantTypes:                        flagx.MustGetStringSlice(cmd, flagClientGrantType),
		JwksUri:                           pointerx.Ptr(flagx.MustGetString(cmd, flagClientJWKSURI)),
		LogoUri:                           pointerx.Ptr(flagx.MustGetString(cmd, flagClientLogoURI)),
		Metadata:                          json.RawMessage(flagx.MustGetString(cmd, flagClientMetadata)),
		Owner:                             pointerx.Ptr(flagx.MustGetString(cmd, flagClientOwner)),
		PolicyUri:                         pointerx.Ptr(flagx.MustGetString(cmd, flagClientPolicyURI)),
		PostLogoutRedirectUris:            flagx.MustGetStringSlice(cmd, flagClientPostLogoutCallback),
		RedirectUris:                      flagx.MustGetStringSlice(cmd, flagClientRedirectURI),
		RequestObjectSigningAlg:           pointerx.Ptr(flagx.MustGetString(cmd, flagClientRequestObjectSigningAlg)),
		RequestUris:                       flagx.MustGetStringSlice(cmd, flagClientRequestURI),
		ResponseTypes:                     flagx.MustGetStringSlice(cmd, flagClientResponseType),
		Scope:                             pointerx.Ptr(strings.Join(flagx.MustGetStringSlice(cmd, flagClientScope), " ")),
		SkipConsent:                       pointerx.Ptr(flagx.MustGetBool(cmd, flagClientSkipConsent)),
		SkipLogoutConsent:                 pointerx.Ptr(flagx.MustGetBool(cmd, flagClientLogoutSkipConsent)),
		SectorIdentifierUri:               pointerx.Ptr(flagx.MustGetString(cmd, flagClientSectorIdentifierURI)),
		SubjectType:                       pointerx.Ptr(flagx.MustGetString(cmd, flagClientSubjectType)),
		TokenEndpointAuthMethod:           pointerx.Ptr(flagx.MustGetString(cmd, flagClientTokenEndpointAuthMethod)),
		TosUri:                            pointerx.Ptr(flagx.MustGetString(cmd, flagClientTOSURI)),
	}, nil
}

func registerEncryptFlags(flags *pflag.FlagSet) {
	// encrypt client secret options
	flags.String(cli.FlagEncryptionPGPKey, "", "Base64 encoded PGP encryption key for encrypting client secret.")
	flags.String(cli.FlagEncryptionPGPKeyURL, "", "PGP encryption key URL for encrypting client secret.")
	flags.String(cli.FlagEncryptionKeybase, "", "Keybase username for encrypting client secret.")
}

func registerClientFlags(flags *pflag.FlagSet) {
	flags.String(flagFile, "", "Read a JSON file representing a client from this location. If set, the other client flags are ignored.")

	flags.String(flagClientMetadata, "{}", "Metadata is an arbitrary JSON String of your choosing.")
	flags.String(flagClientOwner, "", "The owner of this client, typically email addresses or a user ID.")
	flags.StringSlice(flagClientContact, nil, "A list representing ways to contact people responsible for this client, typically email addresses.")
	flags.StringSlice(flagClientRequestURI, nil, "Array of request_uri values that are pre-registered by the RP for use at the OP.")
	flags.String(flagClientRequestObjectSigningAlg, "RS256", "Algorithm that must be used for signing Request Objects sent to the OP.")
	flags.String(flagClientSectorIdentifierURI, "", "URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.")
	flags.StringSlice(flagClientRedirectURI, []string{}, "List of allowed OAuth2 Redirect URIs.")
	flags.StringSlice(flagClientGrantType, []string{"authorization_code"}, "A list of allowed grant types.")
	flags.StringSlice(flagClientResponseType, []string{"code"}, "A list of allowed response types.")
	flags.StringSlice(flagClientScope, []string{}, "The scope the client is allowed to request.")
	flags.StringSlice(flagClientAudience, []string{}, "The audience this client is allowed to request.")
	flags.String(flagClientTokenEndpointAuthMethod, "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are `client_secret_post`, `client_secret_basic`, `private_key_jwt`, and `none`.")
	flags.String(flagClientJWKSURI, "", "Define the URL where the JSON Web Key Set should be fetched from when performing the `private_key_jwt` client authentication method.")
	flags.String(flagClientPolicyURI, "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.")
	flags.String(flagClientTOSURI, "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.")
	flags.String(flagClientClientURI, "", "A URL string of a web page providing information about the client")
	flags.String(flagClientLogoURI, "", "A URL string that references a logo for the client")
	flags.StringSlice(flagClientAllowedCORSOrigin, []string{}, "The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.")
	flags.String(flagClientSubjectType, "public", "A identifier algorithm. Valid values are `public` and `pairwise`.")
	flags.String(flagClientSecret, "", "Provide the client's secret.")
	flags.String(flagClientName, "", "The client's name.")
	flags.StringSlice(flagClientPostLogoutCallback, []string{}, "List of allowed URLs to be redirected to after a logout.")
	flags.Bool(flagClientSkipConsent, false, "Boolean flag specifying whether to skip the consent screen for this client. If omitted, the default value is false.")
	flags.Bool(flagClientLogoutSkipConsent, false, "Boolean flag specifying whether to skip the logout consent screen for this client. If omitted, the default value is false.")
	flags.String(flagClientAccessTokenStrategy, "", "The strategy used to generate access tokens. Valid options are `opaque` and `jwt`.")

	// back-channel logout options
	flags.Bool(flagClientBackChannelLogoutSessionRequired, false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the backchannel-logout-callback is used. If omitted, the default value is false.")
	flags.String(flagClientBackchannelLogoutCallback, "", "Client URL that will cause the client to log itself out when sent a Logout Token by Hydra.")

	// front-channel logout options
	flags.Bool(flagClientFrontChannelLogoutSessionRequired, false, "Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the frontchannel-logout-callback is used. If omitted, the default value is false.")
	flags.String(flagClientFrontChannelLogoutCallback, "", "Client URL that will cause the client to log itself out when rendered in an iframe by Hydra.")

	registerEncryptFlags(flags)
}
