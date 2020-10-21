package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagPage = "page"
	FlagLimit = "limit"
)

func GetPagination(cmd *cobra.Command) (page, limit, offset int, err error) {
	page, err = cmd.Flags().GetInt(FlagPage)
	if err != nil {
		return 0, 0, 0, errors.WithStack(err)
	}
	limit, err = cmd.Flags().GetInt(FlagLimit)
	offset = (page - 1) * limit

	return page, limit, offset, errors.WithStack(err)
}

func RegisterPaginationFlags(flags *pflag.FlagSet) {
	flags.Int("limit", 20, "The maximum amount returned per page.")
	flags.Int("page", 1, "The count of the page (one-based numbering).")
}
