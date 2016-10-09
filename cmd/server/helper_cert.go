package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/square/go-jose"
)

const (
	tlsKeyName = "hydra.https-tls"
)

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

func loadCertificateFromEnv() *tls.Certificate {
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

func getOrCreateTLSCertificate(cmd *cobra.Command, c *config.Config) tls.Certificate {
	if cert := loadCertificateFromFile(cmd); cert != nil {
		return *cert
	} else if cert := loadCertificateFromEnv(); cert != nil {
		return *cert
	}

	ctx := c.Context()
	keys, err := ctx.KeyManager.GetKey(tlsKeyName, "private")
	if errors.Cause(err) == pkg.ErrNotFound {
		logrus.Warn("No TLS Key / Certificate for HTTPS found. Generating self-signed certificate.")

		keys, err = new(jwk.ECDSA256Generator).Generate("")
		pkg.Must(err, "Could not generate key: %s", err)

		cert, err := createSelfSignedCertificate(jwk.First(keys.Key("private")).Key)
		pkg.Must(err, "Could not create X509 PEM Key Pair: %s", err)

		private := jwk.First(keys.Key("private"))
		private.Certificates = []*x509.Certificate{cert}
		keys = &jose.JsonWebKeySet{
			Keys: []jose.JsonWebKey{
				*private,
				*jwk.First(keys.Key("public")),
			},
		}

		err = ctx.KeyManager.AddKeySet(tlsKeyName, keys)
		pkg.Must(err, "Could not persist key: %s", err)
	} else {
		pkg.Must(err, "Could not retrieve key: %s", err)
	}

	private := jwk.First(keys.Key("private"))
	block, err := jwk.PEMBlockForKey(private.Key)
	if err != nil {
		pkg.Must(err, "Could not encode key to PEM: %s", err)
	}

	if len(private.Certificates) == 0 {
		logrus.Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: private.Certificates[0].Raw})
	pemKey := pem.EncodeToMemory(block)
	cert, err := tls.X509KeyPair(pemCert, pemKey)
	pkg.Must(err, "Could not decode certificate: %s", err)

	return cert
}

func createSelfSignedCertificate(key interface{}) (cert *x509.Certificate, err error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return cert, errors.Errorf("Failed to generate serial number: %s", err)
	}

	certificate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Hydra"},
			CommonName:   "Hydra",
		},
		Issuer: pkix.Name{
			Organization: []string{"Hydra"},
			CommonName:   "Hydra",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 7),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certificate.IsCA = true
	certificate.KeyUsage |= x509.KeyUsageCertSign
	certificate.DNSNames = append(certificate.DNSNames, "localhost")
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, publicKey(key), key)
	if err != nil {
		return cert, errors.Errorf("Failed to create certificate: %s", err)
	}

	cert, err = x509.ParseCertificate(der)
	if err != nil {
		return cert, errors.Errorf("Failed to encode private key: %s", err)
	}
	return cert, nil
}
