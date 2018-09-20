/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"time"

	"context"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tlsKeyName = "hydra.https-tls"
)

func loadCertificateFromFile(cmd *cobra.Command, c *config.Config) *tls.Certificate {
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
		c.GetLogger().WithError(errors.WithStack(err)).Warn("Could not load x509 key pair, will use self-issued certificate instead.")
		return nil
	}
	return &cert
}

func loadCertificateFromEnv(c *config.Config) *tls.Certificate {
	keyString := viper.GetString("HTTPS_TLS_KEY")
	certString := viper.GetString("HTTPS_TLS_CERT")
	if keyString == "" || certString == "" {
		return nil
	}

	keyString = strings.Replace(keyString, "\\n", "\n", -1)
	certString = strings.Replace(certString, "\\n", "\n", -1)

	var cert tls.Certificate
	var err error
	if cert, err = tls.X509KeyPair([]byte(certString), []byte(keyString)); err != nil {
		c.GetLogger().Warningf("Could not parse x509 key pair from env: %s", err)
		return nil
	}

	return &cert
}

func getOrCreateTLSCertificate(cmd *cobra.Command, c *config.Config) tls.Certificate {
	if cert := loadCertificateFromFile(cmd, c); cert != nil {
		c.GetLogger().Info("Loaded tls certificate from file")
		return *cert
	} else if cert := loadCertificateFromEnv(c); cert != nil {
		c.GetLogger().Info("Loaded certificate from environment variable")
		return *cert
	}

	ctx := c.Context()
	expectDependency(c.GetLogger(), ctx.KeyManager)

	privateKey, err := createOrGetJWK(c, tlsKeyName, "", "private")
	if err != nil {
		c.GetLogger().WithError(err).Fatalf(`Could not fetch TLS keys - did you forget to run "hydra migrate sql" or forget to set the SYSTEM_SECRET?`)
	}

	if len(privateKey.Certificates) == 0 {
		cert, err := createSelfSignedCertificate(privateKey.Key)
		if err != nil {
			c.GetLogger().WithError(err).Fatalf(`Could not generate a self signed TLS certificate.`)
		}

		privateKey.Certificates = []*x509.Certificate{cert}
		if err := ctx.KeyManager.DeleteKey(context.TODO(), tlsKeyName, privateKey.KeyID); err != nil {
			c.GetLogger().WithError(err).Fatalf(`Could not update (delete) the self signed TLS certificate.`)
		}
		if err := ctx.KeyManager.AddKey(context.TODO(), tlsKeyName, privateKey); err != nil {
			c.GetLogger().WithError(err).Fatalf(`Could not update (add) the self signed TLS certificate.`)
		}
	}

	block, err := jwk.PEMBlockForKey(privateKey.Key)
	if err != nil {
		pkg.Must(err, "Could not encode key to PEM: %s", err)
	}

	if len(privateKey.Certificates) == 0 {
		c.GetLogger().Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: privateKey.Certificates[0].Raw})
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
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().UTC().Add(time.Hour * 24 * 7),
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
