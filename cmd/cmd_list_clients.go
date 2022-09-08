package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewListClientsCmd(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clients",
		Short:   "List OAuth 2.0 Clients",
		Long:    `This command list an OAuth 2.0 Clients.`,
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s ls identities --%s eyJwYWdlIjoxfQ --%s 10", root.Use, cmdx.FlagPageToken, cmdx.FlagPageSize),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			pageToken, pageSize, err := cmdx.ParseTokenPaginationArgs(cmd)
			if err != nil {
				return err
			}

			list, resp, err := m.V0alpha2Api.AdminListOAuth2Clients(cmd.Context()).PageSize(int64(pageSize)).PageToken(pageToken).Execute()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			var collection outputOAuth2ClientCollection
			for k := range list {
				collection.clients = append(collection.clients, list[k])
			}

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
