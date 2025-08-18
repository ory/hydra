// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/server"
	"github.com/ory/hydra/v2/driver"
)

func NewServePublicCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	return &cobra.Command{
		Use:   "public",
		Short: "Serves Public HTTP/2 APIs",
		Long: `This command opens one port and listens to HTTP/2 API requests. The exposed API handles requests coming from
the public internet, like OAuth 2.0 Authorization and Token requests, OpenID Connect UserInfo, OAuth 2.0 Token Revokation,
and OpenID Connect Discovery.

This command is configurable using the same options available to "serve admin" and "serve all".

It is generally recommended to use this command only if you require granular control over the privileged and public APIs.
For example, you might want to run different TLS certificates or CORS settings on the public and privileged API.

This command does not work with the "memory" database. Both services (privileged, public) MUST use the same database
connection to be able to synchronize.

` + serveControls,
		RunE: server.RunServePublic(dOpts),
	}
}
