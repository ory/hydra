// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ory/x/httpx"
)

const (
	envKeyEndpoint    = "ORY_SDK_URL"
	FlagEndpoint      = "endpoint"
	FlagSkipTLSVerify = "skip-tls-verify"
	FlagHeaders       = "http-header"
)

// Remote returns the remote endpoint for the given command.
func Remote(cmd *cobra.Command) (string, error) {
	endpoint, err := cmd.Flags().GetString(FlagEndpoint)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if endpoint != "" {
		return strings.TrimRight(endpoint, "/"), nil
	} else if endpoint := os.Getenv("ORY_SDK_URL"); endpoint != "" {
		return strings.TrimRight(endpoint, "/"), nil
	}

	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "To execute this command, the endpoint URL must point to the URL where Ory is located. To set the endpoint URL, use flag `--endpoint` or environment variable `ORY_SDK_URL`.")
	return "", FailSilently(cmd)
}

// RemoteURI returns the remote URI for the given command.
func RemoteURI(cmd *cobra.Command) (*url.URL, error) {
	remote, err := Remote(cmd)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.ParseRequestURI(remote)
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not parse endpoint URL: %s", err)
		return nil, err
	}

	return endpoint, nil
}

// NewClient creates a new HTTP client.
func NewClient(cmd *cobra.Command) (*http.Client, *url.URL, error) {
	endpoint, err := cmd.Flags().GetString(FlagEndpoint)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if endpoint == "" {
		endpoint = os.Getenv(envKeyEndpoint)
	}

	if endpoint == "" {
		return nil, nil, errors.Errorf("you have to set the remote endpoint, try --help for details")
	}

	u, err := url.Parse(strings.TrimRight(endpoint, "/"))
	if err != nil {
		return nil, nil, errors.Wrapf(err, `could not parse the endpoint URL "%s"`, endpoint)
	}

	hc := retryablehttp.NewClient().StandardClient()
	hc.Timeout = time.Second * 10

	rawHeaders, err := cmd.Flags().GetStringSlice(FlagHeaders)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	header := http.Header{}
	for _, h := range rawHeaders {
		parts := strings.Split(h, ":")
		if len(parts) != 2 {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Unable to parse `--http-header` flag. Format of flag value is a `: ` delimited string like `--http-header 'Some-Header: some-values; other values`. Received: %v", rawHeaders)
			return nil, nil, FailSilently(cmd)
		}

		for k := range parts {
			parts[k] = strings.TrimSpace(parts[k])
		}

		header.Add(parts[0], parts[1])
	}

	skipVerify, err := cmd.Flags().GetBool(FlagSkipTLSVerify)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	rt := httpx.NewTransportWithHeader(header)
	rt.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipVerify, //nolint:gosec // This is a false positive
		},
	}
	hc.Transport = rt
	return hc, u, nil
}

// RegisterHTTPClientFlags registers HTTP client configuration flags.
func RegisterHTTPClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagEndpoint, FlagEndpoint[:1], "", fmt.Sprintf("The API URL this command should target. Alternatively set using the %s environmental variable.", envKeyEndpoint))
	flags.Bool(FlagSkipTLSVerify, false, "Do not verify TLS certificates. Useful when dealing with self-signed certificates. Do not use in production!")
	flags.StringSliceP(FlagHeaders, "H", []string{}, "A list of additional HTTP headers to set. HTTP headers is separated by a `: `, for example: `-H 'Authorization: bearer some-token'`.")
}
