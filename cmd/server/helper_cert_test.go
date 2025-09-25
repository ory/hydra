// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package server_test

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/servicelocatorx"
	"github.com/ory/x/tlsx"

	"github.com/ory/hydra/v2/cmd/server"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
)

func TestGetOrCreateTLSCertificate(t *testing.T) {
	certPath, keyPath, cert, priv := testhelpers.GenerateTLSCertificateFilesForTests(t)
	logger := logrusx.New("", "")
	logger.Logger.ExitFunc = func(code int) { t.Fatalf("Logger called os.Exit(%v)", code) }
	d, err := driver.New(t.Context(),
		driver.WithConfigOptions(configx.WithValues(map[string]interface{}{
			"dsn":                 config.DSNMemory,
			"serve.tls.enabled":   true,
			"serve.tls.cert.path": certPath,
			"serve.tls.key.path":  keyPath,
		})),
		driver.WithServiceLocatorOptions(servicelocatorx.WithLogger(logger)),
	)
	require.NoError(t, err)
	getCert := server.GetOrCreateTLSCertificate(t.Context(), d, d.Config().ServeAdmin(t.Context()).TLS, "admin")
	require.NotNil(t, getCert)
	tlsCert, err := getCert(nil)
	require.NoError(t, err)
	require.NotNil(t, tlsCert)
	if tlsCert.Leaf == nil {
		tlsCert.Leaf, err = x509.ParseCertificate(tlsCert.Certificate[0])
		require.NoError(t, err)
	}
	require.True(t, tlsCert.Leaf.Equal(cert))
	require.True(t, priv.Equal(tlsCert.PrivateKey))

	// generate new cert+key
	newCertPath, newKeyPath, newCert, newPriv := testhelpers.GenerateTLSCertificateFilesForTests(t)
	require.False(t, cert.Equal(newCert))
	require.False(t, priv.Equal(newPriv))
	require.NotEqual(t, certPath, newCertPath)
	require.NotEqual(t, keyPath, newKeyPath)

	hook := test.NewLocal(logger.Logger)

	// move them into place
	require.NoError(t, os.Rename(newKeyPath, keyPath))
	require.NoError(t, os.Rename(newCertPath, certPath))

	// give it some time and check we're reloaded
	time.Sleep(150 * time.Millisecond)
	require.Nil(t, hook.LastEntry())

	// request another certificate: it should be the new one
	tlsCert, err = getCert(nil)
	require.NoError(t, err)
	if tlsCert.Leaf == nil {
		tlsCert.Leaf, err = x509.ParseCertificate(tlsCert.Certificate[0])
		require.NoError(t, err)
	}
	require.True(t, tlsCert.Leaf.Equal(newCert))
	require.True(t, newPriv.Equal(tlsCert.PrivateKey))

	require.NoError(t, os.WriteFile(certPath, []byte{'j', 'u', 'n', 'k'}, 0))

	timeout := time.After(500 * time.Millisecond)
	for {
		if hook.LastEntry() != nil {
			break
		}
		select {
		case <-timeout:
			require.FailNow(t, "expected error log entry")
		default:
		}
	}
	require.Contains(t, hook.LastEntry().Message, "Failed to reload TLS certificates, using previous certificates")
}

func TestGetOrCreateTLSCertificateBase64(t *testing.T) {
	certPath, keyPath, cert, priv := testhelpers.GenerateTLSCertificateFilesForTests(t)
	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)
	certBase64 := base64.StdEncoding.EncodeToString(certPEM)
	keyPEM, err := os.ReadFile(keyPath)
	require.NoError(t, err)
	keyBase64 := base64.StdEncoding.EncodeToString(keyPEM)

	d, err := driver.New(t.Context(), driver.WithConfigOptions(configx.WithValues(map[string]interface{}{
		"dsn":                   config.DSNMemory,
		"serve.tls.enabled":     true,
		"serve.tls.cert.base64": certBase64,
		"serve.tls.key.base64":  keyBase64,
	})))
	require.NoError(t, err)
	getCert := server.GetOrCreateTLSCertificate(t.Context(), d, d.Config().ServeAdmin(t.Context()).TLS, "admin")
	require.NotNil(t, getCert)
	tlsCert, err := getCert(nil)
	require.NoError(t, err)
	require.NotNil(t, tlsCert)
	if tlsCert.Leaf == nil {
		tlsCert.Leaf, err = x509.ParseCertificate(tlsCert.Certificate[0])
		require.NoError(t, err)
	}
	require.True(t, tlsCert.Leaf.Equal(cert))
	require.True(t, priv.Equal(tlsCert.PrivateKey))
}

func TestCreateSelfSignedCertificate(t *testing.T) {
	keys, err := jwk.GenerateJWK(jose.RS256, uuid.Must(uuid.NewV4()).String(), "sig")
	require.NoError(t, err)

	private := keys.Keys[0]
	cert, err := tlsx.CreateSelfSignedCertificate(private.Key)
	require.NoError(t, err)
	server.AttachCertificate(&private, cert)

	var actual jose.JSONWebKeySet
	var b bytes.Buffer
	require.NoError(t, json.NewEncoder(&b).Encode(keys))
	require.NoError(t, json.NewDecoder(&b).Decode(&actual))
}
