package tlsx

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var ErrNoCertificatesConfigured = errors.New("no tls configuration was found")
var ErrInvalidCertificateConfiguration = errors.New("tls configuration is invalid")

func HTTPSCertificate() ([]tls.Certificate, error) {
	return Certificate("HTTPS_TLS")
}

func HTTPSCertificateHelpMessage() string {
	return CertificateHelpMessage("HTTPS_TLS")
}

// CertificateHelpMessage returns a help message for configuring TLS Certificates
func CertificateHelpMessage(prefix string) string {
	return `- ` + prefix + `_CERT_PATH: The path to the TLS certificate (pem encoded).
	Example: ` + prefix + `_CERT_PATH=~/cert.pem

- ` + prefix + `_KEY_PATH: The path to the TLS private key (pem encoded).
	Example: ` + prefix + `_KEY_PATH=~/key.pem

- ` + prefix + `_CERT: Base64 encoded (without padding) string of the TLS certificate (PEM encoded) to be used for HTTP over TLS (HTTPS).
	Example: ` + prefix + `_CERT="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."

- ` + prefix + `_KEY: Base64 encoded (without padding) string of the private key (PEM encoded) to be used for HTTP over TLS (HTTPS).
	Example: ` + prefix + `_KEY="-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFDjBABgkqhkiG9w0BBQ0wMzAbBgkqhkiG9w0BBQwwDg..."
`
}

// Certificate returns loads a TLS Certificate by looking at environment variables
func Certificate(prefix string) ([]tls.Certificate, error) {
	certString, keyString := viper.GetString(prefix+"_CERT"), viper.GetString(prefix+"_KEY")
	certPath, keyPath := viper.GetString(prefix+"_CERT_PATH"), viper.GetString(prefix+"_KEY_PATH")

	if certString == "" && keyString == "" && certPath == "" && keyPath == "" {
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	} else if certString != "" && keyString != "" {
		tlsCertBytes, err := base64.StdEncoding.DecodeString(certString)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode the TLS certificate: %v", err)
		}
		tlsKeyBytes, err := base64.StdEncoding.DecodeString(keyString)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode the TLS private key: %v", err)
		}

		cert, err := tls.X509KeyPair(tlsCertBytes, tlsKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to load X509 key pair: %v", err)
		}
		return []tls.Certificate{cert}, nil
	}

	if certPath != "" && keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load X509 key pair from files: %v", err)
		}
		return []tls.Certificate{cert}, nil
	}

	return nil, errors.WithStack(ErrInvalidCertificateConfiguration)
}

func PublicKey(key interface{}) interface{} {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func CreateSelfSignedCertificate(key interface{}) (cert *x509.Certificate, err error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return cert, errors.Errorf("failed to generate serial number: %s", err)
	}

	certificate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ORY GmbH"},
			CommonName:   "ORY",
		},
		Issuer: pkix.Name{
			Organization: []string{"ORY GmbH"},
			CommonName:   "ORY",
		},
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().UTC().Add(time.Hour * 24 * 31),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certificate.IsCA = true
	certificate.KeyUsage |= x509.KeyUsageCertSign
	certificate.DNSNames = append(certificate.DNSNames, "localhost")
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, PublicKey(key), key)
	if err != nil {
		return cert, errors.Errorf("failed to create certificate: %s", err)
	}

	cert, err = x509.ParseCertificate(der)
	if err != nil {
		return cert, errors.Errorf("failed to encode private key: %s", err)
	}
	return cert, nil
}
