// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cliclient

import (
	"net/url"

	"github.com/pkg/errors"

	"github.com/ory/x/cmdx"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
)

type ContextKey int

const (
	ClientContextKey ContextKey = iota + 1
	OAuth2URLOverrideContextKey
)

func GetOAuth2URLOverride(cmd *cobra.Command, fallback *url.URL) *url.URL {
	if override, ok := cmd.Context().Value(OAuth2URLOverrideContextKey).(func(cmd *cobra.Command) *url.URL); ok {
		return override(cmd)
	}
	return fallback
}

func NewClient(cmd *cobra.Command) (*hydra.APIClient, *url.URL, error) {
	if f, ok := cmd.Context().Value(ClientContextKey).(func(cmd *cobra.Command) (*hydra.APIClient, *url.URL, error)); ok {
		return f(cmd)
	}

	hc, target, err := cmdx.NewClient(cmd)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	conf := hydra.NewConfiguration()
	conf.HTTPClient = hc
	conf.Servers = hydra.ServerConfigurations{{URL: target.String()}}
	return hydra.NewAPIClient(conf), target, nil
}
