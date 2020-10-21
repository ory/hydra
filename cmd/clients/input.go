package clients

import (
	"encoding/json"
	"fmt"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/x/cmdx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"strings"
)

const (
	FlagClientID                 = "id"
	FlagCallbacks                = "callbacks"
	FlagGrantTypes               = "grant-types"
	FlagResponseTypes            = "response-types"
	FlagScope                    = "scope"
	FlagAudience                 = "audience"
	FlagTokentEndpointAuthMethod = "token-endpoint-auth-method"
	FlagJWKsURI                  = "jwks-uri"
	FlagPolicyURI                = "policy-uri"
	FlagTOSURI                   = "tos-uri"
	FlagClientURI                = "client-uri"
	FlagLogoURI                  = "logo-uri"
	FlagAllowedCORSOrigins       = "allowed-cors-origins"
	FlagSubjectType              = "subject-type"
	FlagSecret                   = "secret"
	FlagClientName               = "name"
	FlagPostLogoutCallbacks      = "post-logout-callbacks"

	helperStdInFile = "To read JSON from STD_IN, use \"-\" as the filename."
)

func clientFromAllSources(cmd *cobra.Command, fn string) (*models.OAuth2Client, error) {
	var client *models.OAuth2Client

	if fn == "-" {
		var err error
		client, err = clientFromFile(cmd.InOrStdin())
		if err != nil {
			return nil, err
		}
	} else {
		f, err := os.Open(fn)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not open file %s: %s", fn, err)
			return nil, cmdx.FailSilently(cmd)
		}
		client, err = clientFromFile(f)
		if err != nil {
			return nil, err
		}
	}

	cf, err := clientFromFlags(cmd)
	if err != nil {
		return nil, err
	}

	flagData, err := json.Marshal(cf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal(flagData, client); err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

func clientFromFile(f io.Reader) (*models.OAuth2Client, error) {
	var c models.OAuth2Client
	return &c, errors.WithStack(json.NewDecoder(f).Decode(&c))
}

func clientFromFlags(cmd *cobra.Command) (c *models.OAuth2Client, retErr error) {
	var getString = func(name string) string {
		v, err := cmd.Flags().GetString(name)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error getting the flag value for --%s: %+v", name, err)
			retErr = cmdx.FailSilently(cmd)
		}
		return v
	}

	var getStringSlice = func(name string) []string {
		v, err := cmd.Flags().GetStringSlice(name)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error getting the flag value for --%s: %+v", name, err)
			retErr = cmdx.FailSilently(cmd)
		}
		return v
	}

	c = &models.OAuth2Client{
		AllowedCorsOrigins:      getStringSlice(FlagAllowedCORSOrigins),
		Audience:                getStringSlice(FlagAudience),
		ClientID:                getString(FlagClientID),
		ResponseTypes:           getStringSlice(FlagResponseTypes),
		Scope:                   strings.Join(getStringSlice(FlagScope), " "),
		GrantTypes:              getStringSlice(FlagGrantTypes),
		RedirectUris:            getStringSlice(FlagCallbacks),
		ClientName:              getString(FlagClientName),
		TokenEndpointAuthMethod: getString(FlagTokentEndpointAuthMethod),
		JwksURI:                 getString(FlagJWKsURI),
		TosURI:                  getString(FlagTOSURI),
		PolicyURI:               getString(FlagPolicyURI),
		LogoURI:                 getString(FlagLogoURI),
		ClientURI:               getString(FlagClientURI),
		SubjectType:             getString(FlagSubjectType),
		PostLogoutRedirectUris:  getStringSlice(FlagPostLogoutCallbacks),
	}
	return
}

func registerClientFlags(flags *pflag.FlagSet) {
	flags.String(FlagClientID, "", "Give the client this ID")
	flags.StringSliceP(FlagCallbacks, "c", []string{}, "REQUIRED list of allowed callback URLs")
	flags.StringSliceP(FlagGrantTypes, "g", []string{"authorization_code"}, "A list of allowed grant types")
	flags.StringSliceP(FlagResponseTypes, "r", []string{"code"}, "A list of allowed response types")
	flags.StringSliceP(FlagScope, "a", []string{""}, "The scope the client is allowed to request")
	flags.StringSlice(FlagAudience, []string{}, "The audience this client is allowed to request")
	flags.String(FlagTokentEndpointAuthMethod, "client_secret_basic", "Define which authentication method the client may use at the Token Endpoint. Valid values are \"client_secret_post\", \"client_secret_basic\", \"private_key_jwt\", and \"none\"")
	flags.String(FlagJWKsURI, "", "Define the URL where the JSON Web Key Set should be fetched from when performing the \"private_key_jwt\" client authentication method")
	flags.String(FlagPolicyURI, "", "A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data")
	flags.String(FlagTOSURI, "", "A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client")
	flags.String(FlagClientURI, "", "A URL string of a web page providing information about the client")
	flags.String(FlagLogoURI, "", "A URL string that references a logo for the client")
	flags.StringSlice(FlagAllowedCORSOrigins, []string{}, "The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.")
	flags.String(FlagSubjectType, "public", "A identifier algorithm. Valid values are \"public\" and \"pairwise\"")
	flags.String(FlagSecret, "", "Provide the client's secret")
	flags.StringP(FlagClientName, "n", "", "The client's name")
	flags.StringSlice(FlagPostLogoutCallbacks, []string{}, "List of allowed URLs to be redirected to after a logout")
}

func warnClientSecretFlag(cmd *cobra.Command) {
	// Printing this warning to stderr so that stdout can still be used for piping
	_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "You should not provide secrets using command line flags, the secret might leak to bash history and similar systems")
}
