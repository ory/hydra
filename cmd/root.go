package cmd

import (
	"fmt"
	"os"

	"github.com/ory-am/hydra/cmd/cli"
	"github.com/ory-am/hydra/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"path/filepath"
	"runtime"
	"sync"
)

var cfgFile string

var c = new(config.Config)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hydra",
	Short: "Hydra is a twelve factor OAuth2 and OpenID Connect provider",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var cmdHandler = cli.NewHandler(c)
var mutex = &sync.RWMutex{}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
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

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hydra.yaml)")
	RootCmd.PersistentFlags().Bool("skip-tls-verify", false, "foolishly accept TLS certificates signed by unkown certificate authorities")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	mutex.Lock()
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	path := absPathify("$HOME")
	if _, err := os.Stat(filepath.Join(path, ".hydra.yml")); err != nil {
		_, _ = os.Create(filepath.Join(path, ".hydra.yml"))
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(".hydra") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		fmt.Println("")
	}

	if err := viper.Unmarshal(c); err != nil {
		fatal("Could not read config because %s.", err)
	}

	if consentURL, ok := viper.Get("CONSENT_URL").(string); ok {
		c.ConsentURL = consentURL
	}

	if clientID, ok := viper.Get("CLIENT_ID").(string); ok {
		c.ClientID = clientID
	}

	if clientSecret, ok := viper.Get("CLIENT_SECRET").(string); ok {
		c.ClientSecret = clientSecret
	}

	if c.ClusterURL == "" {
		fmt.Printf("Pointing cluster at %s\n", c.GetClusterURL())
	}
	mutex.Unlock()
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