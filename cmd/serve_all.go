// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/server"
)

// allCmd represents the all command
func NewServeAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Serves both public and administrative HTTP/2 APIs",
		Long: `Starts a process which listens on two ports for public and administrative HTTP/2 API requests.

If you want more granular control (e.g. different TLS settings) over each API group (administrative, public) you
can run "serve admin" and "serve public" separately.

This command exposes a variety of controls via environment variables. You can
set environments using "export KEY=VALUE" (Linux/macOS) or "set KEY=VALUE" (Windows). On Linux,
you can also set environments by prepending key value pairs: "KEY=VALUE KEY2=VALUE2 hydra"

All possible controls are listed below. This command exposes exposes command line flags, which are listed below
the controls section.

` + serveControls,
		Run: server.RunServeAll,
	}
}
