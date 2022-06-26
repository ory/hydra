package cliclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/httpx"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"

	hydra "github.com/ory/hydra-client-go"
)

const (
	envKeyEndpoint    = "HYDRA_ADMIN_URL"
	FlagEndpoint      = "endpoint"
	FlagSkipTLSVerify = "skip-tls-verify"
	FlagHeaders       = "http-header"
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

	rawHeaders := flagx.MustGetStringSlice(cmd, FlagHeaders)
	var header http.Header
	for _, h := range rawHeaders {
		parts := strings.Split(h, ":")
		if len(parts) != 2 {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Unable to parse `--http-header` flag. Format of flag value is a `: ` delimited string like `--http-header 'Some-Header: some-values; other values`. Received: %v", rawHeaders)
			return nil, cmdx.FailSilently(cmd)
		}

		for k := range parts {
			parts[k] = strings.TrimSpace(parts[k])
		}

		header.Add(parts[0], parts[1])
	}

	rt := httpx.NewTransportWithHeader(header)
	rt.RoundTripper = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: flagx.MustGetBool(cmd, FlagSkipTLSVerify)}}
	conf.HTTPClient.Transport = rt

	conf.Servers = hydra.ServerConfigurations{{URL: u.String()}}
	return hydra.NewAPIClient(conf), nil
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagEndpoint, FlagEndpoint[:1], "", fmt.Sprintf("The URL of Ory Hydra' Admin API. Alternatively set using the %s environmental variable.", envKeyEndpoint))
	flags.Bool(FlagSkipTLSVerify, false, "Do not verify TLS certificates. Useful when dealing with self-signed certificates. Do not use in production!")
	flags.StringSliceP(FlagHeaders, "H", []string{}, "A list of additional HTTP headers to set. HTTP headers is separated by a `: `, for example: `-H 'Authorization: bearer some-token'`.")
}
