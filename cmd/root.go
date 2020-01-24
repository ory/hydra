package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"time"

	conf "github.com/coupa/foundation-go/config"
	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	Version   = "dev-master"
	BuildTime = time.Now().String()
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

var (
	cmdHandler                = cli.NewHandler(c)
	sManager   SecretsManager = AwsSecretsManager{}
)

type SecretsManager interface {
	GetSecrets(string) (map[string]string, []byte, error)
}

type AwsSecretsManager struct {
}

func (m AwsSecretsManager) GetSecrets(name string) (map[string]string, []byte, error) {
	return conf.GetSecrets(name)
}

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

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hydra.yaml)")
	RootCmd.PersistentFlags().Bool("skip-tls-verify", false, "foolishly accept TLS certificates signed by unkown certificate authorities")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	smName := os.Getenv("AWSSM_NAME")
	if smName != "" {
		if os.Getenv("AWS_REGION") == "" {
			//Default region to us-east-1
			if err := os.Setenv("AWS_REGION", "us-east-1"); err != nil {
				log.Fatalf("Error setting AWS_REGION: %v", err)
			}
		}
		if err := conf.WriteSecretsToENV(smName); err != nil {
			log.Fatalf("Error reading from Secrets Manager: %v", err)
		}
	}

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

	viper.BindEnv("HOST")
	viper.SetDefault("HOST", "")

	viper.BindEnv("CLIENT_ID")
	viper.SetDefault("CLIENT_ID", "")

	viper.BindEnv("CONSENT_URL")
	viper.SetDefault("CONSENT_URL", "")

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

	viper.BindEnv("PORT")
	viper.SetDefault("PORT", 4444)

	viper.BindEnv("ISSUER")
	viper.SetDefault("ISSUER", "http://localhost:4444")

	viper.BindEnv("BCRYPT_COST")
	viper.SetDefault("BCRYPT_COST", 10)

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

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		fmt.Println("")
	}
	if smName != "" {
		if err := setupAppCerts(); err != nil {
			log.Fatal(err)
		}
		if err := setupDBCerts(); err != nil {
			log.Fatal(err)
		}
	}

	iss := viper.Get("ISSUER")
	viper.Set("ISSUER", strings.TrimSuffix(iss.(string), "/"))

	if err := viper.Unmarshal(c); err != nil {
		fatal(fmt.Sprintf("Could not read config because %s.", err))
	}
}

func setupAppCerts() error {
	appcertsSecretName := os.Getenv("AWS_APP_CERTS_SECRET_NAME")
	if appcertsSecretName == "" {
		return nil
	}
	secrets, _, err := sManager.GetSecrets(appcertsSecretName)
	if err != nil {
		return fmt.Errorf("Error getting app certs from Secrets Manager: %v", err)
	}

	if viper.GetString("HTTPS_TLS_CERT") == "" {
		pem := secrets[pkg.AppSSLCert]
		if pem == "" {
			return fmt.Errorf("App certificate (%s) on Secrets Manager (%s) not found", pkg.AppSSLCert, appcertsSecretName)
		}
		if err = os.Setenv("HTTPS_TLS_CERT", pkg.FixPemFormat(pem)); err != nil {
			return fmt.Errorf("Error setting HTTPS_TLS_CERT with cert from secrets manager: %v", err)
		}
		log.Infof("Successfully set TLS cert from Secrets Manager")
	}
	if viper.GetString("HTTPS_TLS_KEY") == "" {
		key := secrets[pkg.AppSSLKey]
		if key == "" {
			return fmt.Errorf("App key (%s) on Secrets Manager (%s) not found", pkg.AppSSLKey, appcertsSecretName)
		}
		if err = os.Setenv("HTTPS_TLS_KEY", pkg.FixPemFormat(key)); err != nil {
			return fmt.Errorf("Error setting HTTPS_TLS_KEY with TLS key from secrets manager: %v", err)
		}
		log.Infof("Successfully set TLS key from Secrets Manager")
	}
	return nil
}

func setupDBCerts() error {
	rdscertsSecretName := os.Getenv("AWS_RDS_CERTS_SECRET_NAME")
	if rdscertsSecretName == "" {
		return nil
	}
	secrets, _, err := sManager.GetSecrets(rdscertsSecretName)
	if err != nil {
		return fmt.Errorf("Error getting rds certs from Secrets Manager: %v", err)
	}
	pem := secrets[pkg.RdsSSLCert]
	if pem == "" {
		return fmt.Errorf("RDS certificate (%s) on Secrets Manager (%s) not found", pkg.RdsSSLCert, rdscertsSecretName)
	}
	pem = pkg.FixPemFormat(pem)
	viper.Set("RdsSSLCert", pem)
	log.Infof("Successfully set RDS cert from Secrets Manager")
	return nil
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
