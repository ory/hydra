/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package cmd

import (
	"errors"
	"fmt"
	"github.com/ory/hydra/cmd/clients"
	"github.com/ory/hydra/cmd/keys"
	"github.com/ory/hydra/cmd/migrate"
	"github.com/ory/hydra/cmd/token"
	"github.com/ory/x/cmdx"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ory/hydra/driver/configuration"

	"github.com/spf13/cobra"

	"github.com/ory/viper"

	"github.com/ory/hydra/cmd/cli"
)

var cfgFile string

var (
	Version = "master"
	Date    = "undefined"
	Commit  = "undefined"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hydra",
	Short: "Run and manage ORY Hydra",
}

var cmdHandler = cli.NewHandler()

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); errors.Is(err, cmdx.ErrNoPrintButFail) {
		os.Exit(-1)
	} else if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.hydra.yaml)")

	clients.RegisterCommandRecursive(RootCmd)
	keys.RegisterCommandRecursive(RootCmd)
	migrate.RegisterCommandRecursive(RootCmd)
	token.RegisterCommandRecursive(RootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		path := absPathify("$HOME")
		if _, err := os.Stat(filepath.Join(path, ".hydra.yml")); err != nil {
			_, _ = os.Create(filepath.Join(path, ".hydra.yml"))
		}

		viper.SetConfigType("yaml")
		viper.SetConfigName(".hydra") // name of config file (without extension)
		viper.AddConfigPath("$HOME")  // adding home directory as first search path
	}

	viper.SetDefault(configuration.ViperKeyLogLevel, "info")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		fmt.Println("")
	}
}

func absPathify(inPath string) string {
	if strings.HasPrefix(inPath, "$HOME") {
		inPath = userHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
