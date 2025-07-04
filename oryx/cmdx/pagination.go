// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	FlagPageSize  = "page-size"
	FlagPageToken = "page-token"
)

func RegisterTokenPaginationFlags(cmd *cobra.Command) (pageSize int, pageToken string) {
	cmd.Flags().StringVar(&pageToken, FlagPageToken, "", "page token acquired from a previous response")
	cmd.Flags().IntVar(&pageSize, FlagPageSize, 100, "maximum number of items to return")
	return
}

// ParsePaginationArgs parses pagination arguments from the command line.
func ParsePaginationArgs(cmd *cobra.Command, pageArg, perPageArg string) (page, perPage int64, err error) {
	if len(pageArg+perPageArg) > 0 {
		page, err = strconv.ParseInt(pageArg, 0, 64)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not parse page argument\"%s\": %s", pageArg, err)
			return 0, 0, FailSilently(cmd)
		}

		perPage, err = strconv.ParseInt(perPageArg, 0, 64)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not parse per-page argument\"%s\": %s", perPageArg, err)
			return 0, 0, FailSilently(cmd)
		}
	}
	return
}

// ParseTokenPaginationArgs parses token-based pagination arguments from the command line.
func ParseTokenPaginationArgs(cmd *cobra.Command) (page string, perPage int, err error) {
	pageArg, err := cmd.Flags().GetString(FlagPageToken)
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not parse %s argument \"%s\": %s", FlagPageToken, pageArg, err)
		return "", 0, FailSilently(cmd)
	}

	perPageArg, err := cmd.Flags().GetInt(FlagPageSize)
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not parse %s argument \"%d\": %s", FlagPageSize, perPageArg, err)
		return "", 0, FailSilently(cmd)
	}

	return pageArg, perPageArg, nil
}
