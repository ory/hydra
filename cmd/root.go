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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/oauth2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	Version   = "dev-master"
	BuildTime = "undefined"
	GitHash   = "undefined"
)

var c = new(config.Config)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hydra",
	Short: "Hydra is a cloud native high throughput OAuth2 and OpenID Connect provider",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var cmdHandler = cli.NewHandler(c)

// Execute adds all child commands to the root command sets flags appropriately.
// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	c.BuildTime = BuildTime
	c.BuildVersion = Version
	c.BuildHash = GitHash

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.hydra.yaml)")
	RootCmd.PersistentFlags().Bool("skip-tls-verify", false, "Foolishly accept TLS certificates signed by unkown certificate authorities")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	viper.AutomaticEnv() // read in environment variables that match

	viper.BindEnv("PUBLIC_HOST")
	viper.SetDefault("PUBLIC_HOST", "")

	viper.BindEnv("ADMIN_HOST")
	viper.SetDefault("ADMIN_HOST", "")

	viper.BindEnv("PUBLIC_PORT")
	viper.SetDefault("PUBLIC_PORT", 4444)

	viper.BindEnv("ADMIN_PORT")
	viper.SetDefault("ADMIN_PORT", 4445)

	viper.BindEnv("CLIENT_ID")
	viper.SetDefault("CLIENT_ID", "")

	viper.BindEnv("OAUTH2_CONSENT_URL")
	viper.SetDefault("OAUTH2_CONSENT_URL", oauth2.DefaultConsentPath)

	viper.BindEnv("OAUTH2_LOGIN_URL")
	viper.SetDefault("OAUTH2_LOGIN_URL", oauth2.DefaultConsentPath)

	viper.BindEnv("OAUTH2_ERROR_URL")
	viper.SetDefault("OAUTH2_ERROR_URL", oauth2.DefaultErrorPath)

	viper.BindEnv("DATABASE_PLUGIN")
	viper.SetDefault("DATABASE_PLUGIN", "")

	viper.BindEnv("DATABASE_URL")
	viper.SetDefault("DATABASE_URL", "")

	viper.BindEnv("SYSTEM_SECRET")
	viper.SetDefault("SYSTEM_SECRET", "")

	viper.BindEnv("CLIENT_SECRET")
	viper.SetDefault("CLIENT_SECRET", "")

	viper.BindEnv("HTTPS_ALLOW_TERMINATION_FROM")
	viper.SetDefault("HTTPS_ALLOW_TERMINATION_FROM", "")

	viper.BindEnv("CLUSTER_URL")
	viper.SetDefault("CLUSTER_URL", "")

	viper.BindEnv("OAUTH2_ACCESS_TOKEN_STRATEGY")
	viper.SetDefault("OAUTH2_ACCESS_TOKEN_STRATEGY", "opaque")

	viper.BindEnv("OAUTH2_ISSUER_URL")
	viper.SetDefault("OAUTH2_ISSUER_URL", "http://localhost:4444")

	viper.BindEnv("BCRYPT_COST")
	viper.SetDefault("BCRYPT_COST", 10)

	viper.BindEnv("OAUTH2_SHARE_ERROR_DEBUG")
	viper.SetDefault("OAUTH2_SHARE_ERROR_DEBUG", false)

	viper.BindEnv("ACCESS_TOKEN_LIFESPAN")
	viper.SetDefault("ACCESS_TOKEN_LIFESPAN", "1h")

	viper.BindEnv("ID_TOKEN_LIFESPAN")
	viper.SetDefault("ID_TOKEN_LIFESPAN", "1h")

	viper.BindEnv("AUTH_CODE_LIFESPAN")
	viper.SetDefault("AUTH_CODE_LIFESPAN", "10m")

	viper.BindEnv("CHALLENGE_TOKEN_LIFESPAN")
	viper.SetDefault("CHALLENGE_TOKEN_LIFESPAN", "10m")

	viper.BindEnv("LOG_LEVEL")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.BindEnv("LOG_FORMAT")
	viper.SetDefault("LOG_FORMAT", "")

	viper.BindEnv("OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE")
	viper.SetDefault("OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE", "")

	viper.BindEnv("RESOURCE_NAME_PREFIX")
	viper.SetDefault("RESOURCE_NAME_PREFIX", "")

	viper.BindEnv("OIDC_DISCOVERY_CLAIMS_SUPPORTED")
	viper.SetDefault("OIDC_DISCOVERY_CLAIMS_SUPPORTED", "")

	viper.BindEnv("OIDC_DISCOVERY_SCOPES_SUPPORTED")
	viper.SetDefault("OIDC_DISCOVERY_SCOPES_SUPPORTED", "")

	viper.BindEnv("OIDC_DISCOVERY_USERINFO_ENDPOINT")
	viper.SetDefault("OIDC_DISCOVERY_USERINFO_ENDPOINT", "")

	viper.BindEnv("OIDC_SUBJECT_TYPES_SUPPORTED")
	viper.SetDefault("OIDC_SUBJECT_TYPES_SUPPORTED", "public")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		fmt.Println("")
	}

	iss := viper.Get("OAUTH2_ISSUER_URL")
	viper.Set("ISSUER", strings.TrimSuffix(iss.(string), "/"))

	if err := viper.Unmarshal(c); err != nil {
		fatal(fmt.Sprintf("Could not read config because %s.", err))
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
