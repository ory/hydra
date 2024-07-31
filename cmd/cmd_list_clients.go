// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewListClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oauth2-clients",
		Aliases: []string{"clients"},
		Short:   "List OAuth 2.0 Clients",
		Long:    `This command list an OAuth 2.0 Clients.`,
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("{{ .CommandPath }} --%s eyJwYWdlIjoxfQ --%s 10", cmdx.FlagPageToken, cmdx.FlagPageSize),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			pageToken, pageSize, err := cmdx.ParseTokenPaginationArgs(cmd)
			if err != nil {
				return err
			}

			// nolint:bodyclose
			list, resp, err := m.OAuth2API.ListOAuth2Clients(cmd.Context()).PageSize(int64(pageSize)).PageToken(pageToken).Execute()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}
			defer resp.Body.Close()

			var collection outputOAuth2ClientCollection
			collection.clients = append(collection.clients, list...)

			interfaceList := make([]interface{}, len(list))
			for k := range collection.clients {
				interfaceList[k] = interface{}(&list[k])
			}

			result := &cmdx.PaginatedList{Items: interfaceList, Collection: collection}
			result.NextPageToken = getPageToken(resp)
			result.IsLastPage = result.NextPageToken == ""
			cmdx.PrintTable(cmd, result)
			return nil
		},
	}
	cmdx.RegisterTokenPaginationFlags(cmd)
	return cmd
}
