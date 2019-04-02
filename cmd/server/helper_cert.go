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
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"github.com/spf13/viper"

	"github.com/ory/hydra/driver"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/hydra/jwk"
	"github.com/ory/x/tlsx"
)

const (
	tlsKeyName = "hydra.https-tls"
)

func getOrCreateTLSCertificate(cmd *cobra.Command, d driver.Driver) []tls.Certificate {
	cert, err := tlsx.Certificate(
		viper.GetString("serve.tls.cert.base64"),
		viper.GetString("serve.tls.key.base64"),
		viper.GetString("serve.tls.cert.path"),
		viper.GetString("serve.tls.key.path"),
	)

	if err == nil {
		return cert
	} else if errors.Cause(err) != tlsx.ErrNoCertificatesConfigured {
		d.Registry().Logger().WithError(err).Fatalf("Unable to load HTTPS TLS Certificate")
	}

	_, priv, err := jwk.AsymmetricKeypair(context.Background(), d.Registry(), &jwk.RS256Generator{KeyLength: 4069}, tlsKeyName)
	if err != nil {
		d.Registry().Logger().WithError(err).Fatal("Unable to fetch HTTPS TLS key pairs")
	}

	if len(priv.Certificates) == 0 {
		cert, err := tlsx.CreateSelfSignedCertificate(priv.Key)
		if err != nil {
			d.Registry().Logger().WithError(err).Fatalf(`Could not generate a self signed TLS certificate`)
		}

		priv.Certificates = []*x509.Certificate{cert}
		if err := d.Registry().KeyManager().DeleteKey(context.TODO(), tlsKeyName, priv.KeyID); err != nil {
			d.Registry().Logger().WithError(err).Fatal(`Could not update (delete) the self signed TLS certificate`)
		}

		if err := d.Registry().KeyManager().AddKey(context.TODO(), tlsKeyName, priv); err != nil {
			d.Registry().Logger().WithError(err).Fatal(`Could not update (add) the self signed TLS certificate`)
		}
	}

	block, err := jwk.PEMBlockForKey(priv.Key)
	if err != nil {
		d.Registry().Logger().WithError(err).Fatalf("Could not encode key to PEM")
	}

	if len(priv.Certificates) == 0 {
		d.Registry().Logger().Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: priv.Certificates[0].Raw})
	pemKey := pem.EncodeToMemory(block)
	ct, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		d.Registry().Logger().WithError(err).Fatalf("Could not decode certificate")
	}

	return []tls.Certificate{ct}
}
