// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package tlsx

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/watcherx"
)

// ErrNoCertificatesConfigured is returned when no TLS configuration was found.
var ErrNoCertificatesConfigured = errors.New("no tls configuration was found")

// ErrInvalidCertificateConfiguration is returned when an invalid TLS configuration was found.
var ErrInvalidCertificateConfiguration = errors.New("tls configuration is invalid")

// HTTPSCertificate returns loads a HTTP over TLS Certificate by looking at environment variables.
func HTTPSCertificate() ([]tls.Certificate, error) {
	prefix := "HTTPS_TLS"
	return Certificate(
		os.Getenv(prefix+"_CERT"), os.Getenv(prefix+"_KEY"),
		os.Getenv(prefix+"_CERT_PATH"), os.Getenv(prefix+"_KEY_PATH"),
	)
}

// HTTPSCertificateHelpMessage returns a help message for configuring HTTP over TLS Certificates.
func HTTPSCertificateHelpMessage() string {
	return CertificateHelpMessage("HTTPS_TLS")
}

// CertificateHelpMessage returns a help message for configuring TLS Certificates.
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

// CertificateFromBase64 loads a TLS certificate from a base64-encoded string of
// the PEM representations of the cert and key.
func CertificateFromBase64(certBase64, keyBase64 string) (tls.Certificate, error) {
	certPEM, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to base64 decode the TLS certificate: %v", err)
	}
	keyPEM, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to base64 decode the TLS private key: %v", err)
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to load X509 key pair: %v", err)
	}
	return cert, nil
}

// [deprecated] Certificate returns a TLS Certificate by looking at its
// arguments. If both certPEMBase64 and keyPEMBase64 are not empty and contain
// base64-encoded PEM representations of a cert and key, respectively, that key
// pair is returned. Otherwise, if certPath and keyPath point to PEM files, the
// key pair is loaded from those. Returns ErrNoCertificatesConfigured if all
// arguments are empty, and ErrInvalidCertificateConfiguration if the arguments
// are inconsistent.
//
// This function is deprecated. Use CertificateFromBase64 or GetCertificate
// instead.
func Certificate(
	certPEMBase64, keyPEMBase64 string,
	certPath, keyPath string,
) ([]tls.Certificate, error) {
	if certPEMBase64 == "" && keyPEMBase64 == "" && certPath == "" && keyPath == "" {
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	}

	if certPEMBase64 != "" && keyPEMBase64 != "" {
		cert, err := CertificateFromBase64(certPEMBase64, keyPEMBase64)
		if err != nil {
			return nil, errors.WithStack(err)
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

// GetCertificate returns a function for use with
// "net/tls".Config.GetCertificate.
//
// The certificate and private key are read from the specified filesystem paths.
// The certificate file is watched for changes, upon which the cert+key are
// reloaded in the background. Errors during reloading are deduplicated and
// reported through the errs channel if it is not nil. When the provided context
// is canceled, background reloading stops and the errs channel is closed.
//
// The returned function always yields the latest successfully loaded
// certificate; ClientHelloInfo is unused.
func GetCertificate(
	ctx context.Context,
	certPath, keyPath string,
	errs chan<- error,
) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), error) {
	if certPath == "" || keyPath == "" {
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("unable to load X509 key pair from files: %v", err))
	}
	var store atomic.Value
	store.Store(&cert)

	events := make(chan watcherx.Event)
	// The cert could change without the key changing, but not the other way around.
	// Hence, we only watch the cert.
	_, err = watcherx.WatchFile(ctx, certPath, events)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	go func() {
		if errs != nil {
			defer close(errs)
		}
		var lastReportedError string
		for {
			select {
			case <-ctx.Done():
				return

			case event := <-events:
				var err error
				switch event := event.(type) {
				case *watcherx.ChangeEvent:
					var cert tls.Certificate
					cert, err = tls.LoadX509KeyPair(certPath, keyPath)
					if err == nil {
						store.Store(&cert)
						lastReportedError = ""
						continue
					}
					err = fmt.Errorf("unable to load X509 key pair from files: %v", err)

				case *watcherx.ErrorEvent:
					err = fmt.Errorf("file watch: %v", event)
				default:
					continue
				}

				if err.Error() == lastReportedError { // same message as before: don't spam the error channel
					continue
				}
				// fresh error
				select {
				case errs <- errors.WithStack(err):
					lastReportedError = err.Error()
				case <-time.After(500 * time.Millisecond):
				}
			}
		}
	}()

	return func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, ok := store.Load().(*tls.Certificate); ok {
			return cert, nil
		}
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	}, nil
}

// PublicKey returns the public key for a given private key, or nil.
func PublicKey(key crypto.PrivateKey) interface{ Equal(x crypto.PublicKey) bool } {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

// CreateSelfSignedTLSCertificate creates a self-signed TLS certificate.
func CreateSelfSignedTLSCertificate(key interface{}) (*tls.Certificate, error) {
	c, err := CreateSelfSignedCertificate(key)
	if err != nil {
		return nil, err
	}

	block, err := PEMBlockForKey(key)
	if err != nil {
		return nil, err
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
	pemKey := pem.EncodeToMemory(block)
	cert, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

// CreateSelfSignedCertificate creates a self-signed x509 certificate.
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

// PEMBlockForKey returns a PEM-encoded block for key.
func PEMBlockForKey(key interface{}) (*pem.Block, error) {
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &pem.Block{Type: "PRIVATE KEY", Bytes: b}, nil
}
