package cmd

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/cmd/server"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
	"crypto/tls"
	"github.com/ory-am/hydra/jwk"
	"github.com/go-errors/errors"
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
}

func runHostCmd(cmd *cobra.Command, args []string) {
	router := httprouter.New()
	serverHandler := &server.Handler{}
	serverHandler.Start(c, router)

	if ok, _ := cmd.Flags().GetBool("dangerous-auto-logon"); ok {
		logrus.Warnln("Do not use flag --dangerous-auto-logon in production.")
		err := c.Persist()
		pkg.Must(err, "Could not write configuration file: ", err)
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

func getOrCreateTLSCertificate() tls.Certificate {
	ctx := c.Context()
	key, err := ctx.KeyManager.GetKey(TLSKeyName, "private")
	if errors.Is(err, pkg.ErrNotFound) {
		logrus.Warn("Key for TLS not found. Creating new one.")

		generator := jwk.ECDSA521Generator{}
		keys, err := generator.Generate("")
		pkg.Must(err, "Could not generate key: %s", err)

		err = ctx.KeyManager.AddKeySet(TLSKeyName, keys)
		pkg.Must(err, "Could not persist key: %s", err)

		key, err = ctx.KeyManager.GetKey(TLSKeyName, "private")
		pkg.Must(err, "Could not retrieve persisted key: %s", err)
		logrus.Warn("Temporary key created.")
	} else {
		pkg.Must(err, "Could not retrieve key: %s", err)
	}

	pemCert, pemKey, err := jwk.ToX509PEMKeyPair(key.Key)
	pkg.Must(err, "Could not create X509 PEM Key Pair: %s", err)

	cert, err := tls.X509KeyPair(pemCert, pemKey)
	pkg.Must(err, "Could not create TLS Certificate: %s", err)

	return cert
}