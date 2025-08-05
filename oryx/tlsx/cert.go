// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package tlsx

import (
	"bytes"
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
	"io"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

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

type CertFunc = func(*tls.ClientHelloInfo) (*tls.Certificate, error)

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
) (CertFunc, error) {
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

	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
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
func CreateSelfSignedTLSCertificate(key interface{}, opts ...CertificateOpts) (*tls.Certificate, error) {
	c, err := CreateSelfSignedCertificate(key, opts...)
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
func CreateSelfSignedCertificate(key interface{}, opts ...CertificateOpts) (cert *x509.Certificate, err error) {
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
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
	}
	for _, opt := range opts {
		opt(certificate)
	}

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

// NewClientCert creates a new client TLS certificate signed by the given CA.
func NewClientCert(CAcert *x509.Certificate, CAkey crypto.PrivateKey, opts ...CertificateOpts) (*tls.Certificate, error) {
	if !slices.Contains(CAcert.ExtKeyUsage, x509.ExtKeyUsageClientAuth) {
		return nil, errors.Errorf("the CA certificate does not have the client authentication extended key usage (OID 1.3.6.1.5.5.7.3.2) set")
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Errorf("failed to generate serial number: %s", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return nil, errors.Errorf("failed to generate private key: %s", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Ory GmbH"},
			CommonName:   "ORY",
		},
		Issuer:                CAcert.Subject,
		NotBefore:             time.Now().UTC(),
		NotAfter:              CAcert.NotAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	for _, opt := range opts {
		opt(template)
	}

	der, err := x509.CreateCertificate(rand.Reader, template, CAcert, PublicKey(key), CAkey)
	if err != nil {
		return nil, errors.Errorf("failed to create certificate: %s", err)
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	pemBlock, err := PEMBlockForKey(key)
	if err != nil {
		return nil, err
	}
	pemKey := pem.EncodeToMemory(pemBlock)

	cert, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &cert, nil
}

type CertificateOpts func(*x509.Certificate)

// CreateSelfSignedCertificateForTest writes a new, self-signed TLS
// certificate+key (in PEM format) to a temporary location on disk and returns
// the paths to both, and the respective contents in base64 encoding. The
// files are automatically cleaned up when the given *testing.T concludes its
// tests.
func CreateSelfSignedCertificateForTest(t testing.TB) (certPath, keyPath, certBase64, keyBase64 string) {
	tmpDir := t.TempDir()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cert, err := CreateSelfSignedCertificate(privateKey)
	require.NoError(t, err)

	// write cert
	certFile, err := os.Create(filepath.Join(tmpDir, "cert.pem"))
	require.NoError(t, err)
	certPath = certFile.Name()

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	require.NoErrorf(t, pem.Encode(
		io.MultiWriter(enc, certFile),
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw},
	), "Failed to write data to %q", certPath)
	require.NoError(t, enc.Close())
	require.NoErrorf(t, certFile.Close(), "Error closing %q", certPath)
	certBase64 = buf.String()

	// write key
	keyFile, err := os.Create(filepath.Join(tmpDir, "key.pem"))
	require.NoError(t, err)
	keyPath = keyFile.Name()
	buf.Reset()
	enc = base64.NewEncoder(base64.StdEncoding, &buf)

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)

	require.NoErrorf(t, pem.Encode(
		io.MultiWriter(enc, keyFile),
		&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes},
	), "Failed to write data to %q", keyPath)
	require.NoError(t, enc.Close())
	require.NoErrorf(t, keyFile.Close(), "Error closing %q", keyPath)
	keyBase64 = buf.String()

	return
}
