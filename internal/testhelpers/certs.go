// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/tlsx"
)

// GenerateTLSCertificateFilesForTests writes a new, self-signed TLS
// certificate+key (in PEM format) to a temporary location on disk and returns
// the paths to both. The files are automatically cleaned up when the given
// *testing.T concludes its tests.
func GenerateTLSCertificateFilesForTests(t *testing.T) (
	certPath, keyPath string,
	cert *x509.Certificate,
	privateKey interface {
		Public() crypto.PublicKey
		Equal(x crypto.PrivateKey) bool
	},
) {
	tmpDir := t.TempDir()
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	cert, err = tlsx.CreateSelfSignedCertificate(privateKey)
	require.NoError(t, err)

	certOut, err := os.CreateTemp(tmpDir, "test-*-cert.pem")
	require.NoError(t, err, "Failed to create temp file for certificate: %v", err)
	certPath = certOut.Name()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	require.NoError(t, err, "Failed to write data to %q: %v", certPath, err)

	err = certOut.Close()
	require.NoError(t, err, "Error closing %q: %v", certPath, err)

	t.Log("wrote", certPath)

	keyOut, err := os.CreateTemp(tmpDir, "test-*-key.pem")
	require.NoError(t, err, "Failed to create temp file for key: %v", err)
	keyPath = keyOut.Name()

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err, "Failed to marshal private key: %v", err)

	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	require.NoError(t, err, "Failed to write data to %q: %v", keyPath, err)

	err = keyOut.Close()
	require.NoError(t, err, "Error closing %q: %v", keyPath, err)

	t.Log("wrote", keyPath)
	return
}
