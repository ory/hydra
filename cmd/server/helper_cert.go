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
	"context"
	"crypto/sha1" // #nosec G505 - This is required for certificate chains alongside sha256
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/tlsx"

	"github.com/ory/hydra/jwk"
)

const (
	TlsKeyName = "hydra.https-tls"
)

func AttachCertificate(priv *jose.JSONWebKey, cert *x509.Certificate) {
	priv.Certificates = []*x509.Certificate{cert}
	sig256 := sha256.Sum256(cert.Raw)
	// #nosec G401 - This is required for certificate chains alongside sha256
	sig1 := sha1.Sum(cert.Raw)
	priv.CertificateThumbprintSHA256 = sig256[:]
	priv.CertificateThumbprintSHA1 = sig1[:]
}

func GetOrCreateTLSCertificate(cmd *cobra.Command, d driver.Registry, iface config.ServeInterface) []tls.Certificate {
	cert, err := d.Config().TLS(iface).Certificate()

	if err == nil {
		return cert
	} else if !errors.Is(err, tlsx.ErrNoCertificatesConfigured) {
		d.Logger().WithError(err).Fatalf("Unable to load HTTPS TLS Certificate")
	}

	_, priv, err := jwk.GetOrGenerateKeys(context.Background(), d, d.SoftwareKeyManager(), TlsKeyName, TlsKeyName, "RS256")
	if err != nil {
		d.Logger().WithError(err).Fatal("Unable to fetch or generate HTTPS TLS key pair")
	}

	if len(priv.Certificates) == 0 {
		cert, err := tlsx.CreateSelfSignedCertificate(priv.Key)
		if err != nil {
			d.Logger().WithError(err).Fatalf(`Could not generate a self signed TLS certificate`)
		}

		AttachCertificate(priv, cert)
		if err := d.SoftwareKeyManager().DeleteKey(context.TODO(), TlsKeyName, priv.KeyID); err != nil {
			d.Logger().WithError(err).Fatal(`Could not update (delete) the self signed TLS certificate`)
		}

		if err := d.SoftwareKeyManager().AddKey(context.TODO(), TlsKeyName, priv); err != nil {
			d.Logger().WithError(err).Fatalf(`Could not update (add) the self signed TLS certificate: %s %x %d`, cert.SignatureAlgorithm, cert.Signature, len(cert.Signature))
		}
	}

	block, err := jwk.PEMBlockForKey(priv.Key)
	if err != nil {
		d.Logger().WithError(err).Fatalf("Could not encode key to PEM")
	}

	if len(priv.Certificates) == 0 {
		d.Logger().Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: priv.Certificates[0].Raw})
	pemKey := pem.EncodeToMemory(block)
	ct, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		d.Logger().WithError(err).Fatalf("Could not decode certificate")
	}

	return []tls.Certificate{ct}
}
