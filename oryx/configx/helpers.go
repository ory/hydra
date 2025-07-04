// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags registers the config file flag.
func RegisterFlags(flags *pflag.FlagSet) {
	flags.StringSliceP("config", "c", []string{}, "Path to one or more .json, .yaml, .yml, .toml config files. Values are loaded in the order provided, meaning that the last config file overwrites values from the previous config file.")
}

// host = unix:/path/to/socket => port is discarded, otherwise format as host:port
func GetAddress(host string, port int) string {
	if strings.HasPrefix(host, "unix:") {
		return host
	}
	return fmt.Sprintf("%s:%d", host, port)
}
