package cmd

import (
	"net/http"
	"crypto/tls"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/cmd/server"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"bytes"
	"encoding/gob"
	"github.com/square/go-jose"
)

const (
	TLSKeyName = "hydra.tls"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the HTTP/2 host service",
	Long: `Starts all HTTP/2 APIs and connects to a backend.

This command supports the following environment variables:

- DATABASE_URL: A URL to a persistent backend. Hydra supports various backends:
  - None: If DATABASE_URL is empty, all data will be lost when the command is killed.
  - RethinkDB: If DATABASE_URL is a DSN starting with rethinkdb://, RethinkDB will be used as storage backend.

- SYSTEM_SECRET: A secret that is at least 16 characters long. If none is provided, one will be generated. They key
	is used to encrypt sensitive data using AES-GCM (256 bit) and validate HMAC signatures.

- HTTPS_TLS_CERT_PATH: The path to the TLS certificate (pem encoded).
- HTTPS_TLS_KEY_PATH: The path to the TLS private key (pem encoded).
- HTTPS_TLS_CERT: A pem encoded TLS certificate passed as string. Can be used instead of TLS_CERT_PATH.
- HTTPS_TLS_KEY: A pem encoded TLS key passed as string. Can be used instead of TLS_KEY_PATH.

- HYDRA_PROFILING: Set "HYDRA_PROFILING=1" to enable profiling.
`,
	Run:   runHostCmd,
}

func init() {
	RootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	hostCmd.Flags().Bool("dangerous-auto-logon", false, "Stores the root credentials in ~/.hydra.yml. Do not use in production.")
	hostCmd.Flags().String("https-tls-key-path", "", "Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.")
	hostCmd.Flags().String("https-tls-cert-path", "", "Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.")
}

func runHostCmd(cmd *cobra.Command, args []string) {
	router := httprouter.New()
	serverHandler := &server.Handler{}
	serverHandler.Start(c, router)

	if ok, _ := cmd.Flags().GetBool("dangerous-auto-logon"); ok {
		logrus.Warnln("Do not use flag --dangerous-auto-logon in production.")
		err := c.Persist()
		pkg.Must(err, "Could not write configuration file: %s", err)
	}

	http.Handle("/", router)

	srv := &http.Server{
		Addr: c.GetAddress(),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				getOrCreateTLSCertificate(cmd),
			},
		},
	}

	logrus.Infof("Starting server on %s", c.GetAddress())
	err := srv.ListenAndServeTLS("", "")
	pkg.Must(err, "Could not start server: %s %s.", err)
}

func loadCertificateFromFile(cmd *cobra.Command) *tls.Certificate {
	keyPath := viper.GetString("HTTPS_TLS_KEY_PATH")
	certPath := viper.GetString("HTTPS_TLS_CERT_PATH")
	if kp, _ := cmd.Flags().GetString("https-tls-key-path"); kp != "" {
		keyPath = kp
	} else if cp, _ := cmd.Flags().GetString("https-tls-cert-path"); cp != "" {
		certPath = cp
	} else if keyPath == "" || certPath == "" {
		return nil
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		logrus.Warn("Could not load x509 key pair: %s", cert)
		return nil
	}
	return &cert
}

func loadCertificateFromEnv(cmd *cobra.Command) *tls.Certificate {
	keyString := viper.GetString("HTTPS_TLS_KEY")
	certString := viper.GetString("HTTPS_TLS_CERT")
	if keyString == "" || certString == "" {
		return nil
	}

	var cert tls.Certificate
	var err error
	if cert, err = tls.X509KeyPair([]byte(certString), []byte(keyString)); err != nil {
		logrus.Warn("Could not parse x509 key pair from env: %s", cert)
		return nil
	}

	return &cert
}

func getOrCreateTLSCertificate(cmd *cobra.Command) tls.Certificate {
	if cert := loadCertificateFromFile(cmd); cert != nil {
		return *cert
	} else if cert := loadCertificateFromEnv(cmd); cert != nil {
		return *cert
	}

	ctx := c.Context()
	keys, err := ctx.KeyManager.GetKey(TLSKeyName, "private")
	if errors.Is(err, pkg.ErrNotFound) {
		logrus.Warn("Key for TLS not found. Creating new one.")

		keys, err = new(jwk.ECDSA256Generator).Generate("")
		pkg.Must(err, "Could not generate key: %s", err)

		err = ctx.KeyManager.AddKeySet(TLSKeyName, keys)
		pkg.Must(err, "Could not persist key: %s", err)
	} else {
		pkg.Must(err, "Could not retrieve key: %s", err)
	}

	var network bytes.Buffer
	gob.Register(tls.Certificate{})
	certificateJWK, err := ctx.KeyManager.GetKey(TLSKeyName, "certificate")
	if errors.Is(err, pkg.ErrNotFound) {
		pemCert, pemKey, err := jwk.ToX509PEMKeyPair(jwk.First(keys.Key("private")).Key)
		pkg.Must(err, "Could not create X509 PEM Key Pair: %s", err)

		certificate, err := tls.X509KeyPair(pemCert, pemKey)
		pkg.Must(err, "Could not create TLS Certificate: %s", err)

		certificate.PrivateKey = nil
		certificate.Leaf = nil
		err = gob.NewEncoder(&network).Encode(certificate)
		pkg.Must(err, "Could not create TLS Certificate: %s", err)

		err = ctx.KeyManager.AddKey(TLSKeyName, &jose.JsonWebKey{
			KeyID: "certificate",
			Key: network.Bytes(),
		})
		pkg.Must(err, "Could not persist certificate: %s", err)
	} else if err == nil {
		certificateBytes, ok := jwk.First(certificateJWK.Keys).Key.([]byte)
		if !ok {
			err = errors.New("Certificate type assertion failed")
			pkg.Must(err, "Could decode certificate: %s", err)
		}
		network = *bytes.NewBuffer(certificateBytes)
	} else {
		pkg.Must(err, "Could not retrieve certificate: %s", err)
	}

	var certificate tls.Certificate
	err = gob.NewDecoder(&network).Decode(&certificate)
	certificate.PrivateKey = jwk.First(keys.Key("private")).Key
	pkg.Must(err, "Could not retrieve certificate: %s", err)
	return certificate
}
