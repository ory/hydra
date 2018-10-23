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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/tlsx"
)

const (
	tlsKeyName = "hydra.https-tls"
)

func getOrCreateTLSCertificate(cmd *cobra.Command, c *config.Config) []tls.Certificate {
	cert, err := tlsx.HTTPSCertificate()
	if err == nil {
		return cert
	} else if errors.Cause(err) != tlsx.ErrNoCertificatesConfigured {
		c.GetLogger().WithError(err).Fatalf("Unable to load HTTPS TLS Certificate")
	}

	ctx := c.Context()
	expectDependency(c.GetLogger(), ctx.KeyManager)

	privateKey, err := createOrGetJWK(c, tlsKeyName, "", "private")
	if err != nil {
		c.GetLogger().WithError(err).Fatalf(`Could not fetch TLS keys - did you forget to run "hydra migrate sql" or forget to set the SYSTEM_SECRET?`)
	}

	if len(privateKey.Certificates) == 0 {
		cert, err := tlsx.CreateSelfSignedCertificate(privateKey.Key)
		if err != nil {
			c.GetLogger().WithError(err).Fatalf(`Could not generate a self signed TLS certificate`)
		}

		privateKey.Certificates = []*x509.Certificate{cert}
		if err := ctx.KeyManager.DeleteKey(context.TODO(), tlsKeyName, privateKey.KeyID); err != nil {
			c.GetLogger().WithError(err).Fatal(`Could not update (delete) the self signed TLS certificate`)
		}

		if err := ctx.KeyManager.AddKey(context.TODO(), tlsKeyName, privateKey); err != nil {
			c.GetLogger().WithError(err).Fatal(`Could not update (add) the self signed TLS certificate`)
		}
	}

	block, err := jwk.PEMBlockForKey(privateKey.Key)
	if err != nil {
		c.GetLogger().WithError(err).Fatalf("Could not encode key to PEM")
	}

	if len(privateKey.Certificates) == 0 {
		c.GetLogger().Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: privateKey.Certificates[0].Raw})
	pemKey := pem.EncodeToMemory(block)
	ct, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		c.GetLogger().WithError(err).Fatalf("Could not decode certificate")
	}

	return []tls.Certificate{ct}
}
