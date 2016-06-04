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
	Short: "Start the hydra host service",
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
	hostCmd.Flags().String("tls-key-path", "", "Path to the key file for HTTP/2 over TLS.")
	hostCmd.Flags().String("tls-cert-path", "", "Path to the certificate file for HTTP/2 over TLS.")
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
				getOrCreateTLSCertificate(),
			},
		},
	}

	logrus.Infof("Starting server on %s", c.GetAddress())
	err := srv.ListenAndServeTLS("", "")
	pkg.Must(err, "Could not start server: %s %s.", err)
}

func loadCertificateFromFile(cmd *cobra.Command) {
	keyPath := viper.Get("TLS_KEY_PATH")
	certPath := viper.Get("TLS_CERT_PATH")
	if kp, _ := cmd.Flags().GetString("tls-key-path"); kp != "" {
		keyPath = kp
	} else if cp, _ := cmd.Flags().GetString("tls-cert-path"); cp != "" {
		certPath = cp
	} else if keyPath == "" || certPath == "" {
		return
	}


}

func getOrCreateTLSCertificate() tls.Certificate {
	ctx := c.Context()
	keys, err := ctx.KeyManager.GetKey(TLSKeyName, "private")
	if errors.Is(err, pkg.ErrNotFound) {
		logrus.Warn("Key for TLS not found. Creating new one.")

		generator := jwk.ECDSA256Generator{}
		var set *jose.JsonWebKeySet
		set, err = generator.Generate("")
		pkg.Must(err, "Could not generate key: %s", err)

		err = ctx.KeyManager.AddKeySet(TLSKeyName, set)
		pkg.Must(err, "Could not persist key: %s", err)

		keys, err = ctx.KeyManager.GetKey(TLSKeyName, "private")
		pkg.Must(err, "Could not retrieve persisted key: %s", err)
		logrus.Warn("Temporary key created.")
	}
	pkg.Must(err, "Could not retrieve key: %s", err)

	var network bytes.Buffer
	certificateJWK, err := ctx.KeyManager.GetKey(TLSKeyName, "certificate")
	if errors.Is(err, pkg.ErrNotFound) {
		pemCert, pemKey, err := jwk.ToX509PEMKeyPair(jwk.First(keys.Keys).Key)
		pkg.Must(err, "Could not create X509 PEM Key Pair: %s", err)

		certificate, err := tls.X509KeyPair(pemCert, pemKey)
		pkg.Must(err, "Could not create TLS Certificate: %s", err)

		err = gob.NewEncoder(&network).Encode(certificate)
		pkg.Must(err, "Could not create TLS Certificate: %s", err)

		err = ctx.KeyManager.AddKey(TLSKeyName, jose.JsonWebKey{
			KeyID: "certificate",
			Key: network.Bytes(),
		})
		pkg.Must(err, "Could not persist certificate: %s", err)
	} else if err == nil {
		certificateBytes, ok := jwk.First(certificateJWK.Keys).Key.([]byte)
		if !ok {
			err =errors.New("Certificate type assertion failed")
			pkg.Must(err, "Could decode certificate: %s", err)
		}
		network = bytes.NewBuffer(certificateBytes)
	}
	pkg.Must(err, "Could not retrieve certificate: %s", err)

	var certificate tls.Certificate
	err = gob.NewDecoder(network).Decode(&certificate)
	pkg.Must(err, "Could not retrieve certificate: %s", err)
	return &certificate
}
