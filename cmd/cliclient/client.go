package cliclient

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"

	hydra "github.com/ory/hydra-client-go"
)

const (
	envKeyEndpoint = "HYDRA_ADMIN_URL"
	FlagEndpoint   = "endpoint"
)

type ContextKey int

const (
	ClientContextKey ContextKey = iota + 1
)

func NewClient(cmd *cobra.Command) (*hydra.APIClient, error) {
	if f, ok := cmd.Context().Value(ClientContextKey).(func(cmd *cobra.Command) (*hydra.APIClient, error)); ok {
		return f(cmd)
	} else if f != nil {
		return nil, errors.Errorf("ClientContextKey was expected to be *client.OryHydra but it contained an invalid type %T ", f)
	}

	endpoint, err := cmd.Flags().GetString(FlagEndpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if endpoint == "" {
		endpoint = os.Getenv(envKeyEndpoint)
	}

	if endpoint == "" {
		return nil, errors.Errorf("you have to set the remote endpoint, try --help for details")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, `could not parse the endpoint URL "%s"`, endpoint)
	}

	conf := hydra.NewConfiguration()
	conf.HTTPClient = retryablehttp.NewClient().StandardClient()
	conf.HTTPClient.Timeout = time.Second * 10
	conf.Servers = hydra.ServerConfigurations{{URL: u.String()}}
	return hydra.NewAPIClient(conf), nil
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagEndpoint, FlagEndpoint[:1], "", fmt.Sprintf("The URL of Ory Hydra' Admin API. Alternatively set using the %s environmental variable.", envKeyEndpoint))
}
