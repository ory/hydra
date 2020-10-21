package clients

import (
	"fmt"
	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/pointerx"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List OAuth 2.0 Clients",
		Long:    "This command lists OAuth 2.0 Clients (paginated).",
		Example: fmt.Sprintf("$ hydra clients list --%s 10 --%s 100", cli.FlagPage, cli.FlagLimit),
		RunE: func(cmd *cobra.Command, _ []string) error {
			m := cli.ConfigureClient(cmd)

			_, limit, offset, err := cli.GetPagination(cmd)
			if err != nil {
				return err
			}

			response, err := m.Admin.ListOAuth2Clients(admin.NewListOAuth2ClientsParams().WithLimit(pointerx.Int64(int64(limit))).WithOffset(pointerx.Int64(int64(offset))))
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "The request failed with the following error message:\n%s", cli.FormatSwaggerError(err))
				return cmdx.FailSilently(cmd)
			}

			cmdx.PrintCollection(cmd, outputOAuth2ClientCollection(response.Payload))
			return nil
		},
	}

	cli.RegisterPaginationFlags(cmd.LocalFlags())

	return cmd
}
