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
	"path"
	"runtime"
	"sync"
	"time"

	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/tlsx"

	"github.com/ory/hydra/jwk"

	"github.com/fsnotify/fsnotify"
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

func GetOrCreateTLSCertificate(cmd *cobra.Command, d driver.Registry, iface config.ServeInterface) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, location, err := d.Config().TLS(iface).Certificate()

	if err == nil {
		return newCertificatesProvider(cert, location, d, iface).getCertificate
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

	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return &ct, nil
	}
}

type certificatesProvider struct {
	certs        []tls.Certificate
	mu           sync.Mutex
	iface        config.ServeInterface
	certLocation *config.CertLocation
	d            driver.Registry
	watcher      *fsnotify.Watcher
}

func newCertificatesProvider(certs []tls.Certificate, certLocation *config.CertLocation, d driver.Registry, iface config.ServeInterface) *certificatesProvider {
	ret := &certificatesProvider{
		certLocation: certLocation,
		d:            d,
		iface:        iface,
	}
	ret.load(certs)
	if certLocation != nil {
		ret.watchCertificatesChanges()
	}

	runtime.SetFinalizer(ret, func(ret *certificatesProvider) { ret.stop() })

	return ret
}

func (p *certificatesProvider) load(certs []tls.Certificate) {
	for i := range certs {
		tlsCert := &certs[i]
		if tlsCert.Leaf != nil {
			continue
		}
		for _, bCert := range tlsCert.Certificate {
			cert, _ := x509.ParseCertificate(bCert)
			if !cert.IsCA {
				tlsCert.Leaf = cert
			}
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.certs = certs
}

func (p *certificatesProvider) watchCertificatesChanges() {
	var err error
	p.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		p.d.Logger().WithError(err).Fatalf("Could not activate certificate change watcher")
	}

	go func() {
		p.d.Logger().Infof("Starting tls certificate auto-refresh")
		for {
			select {
			case _, ok := <-p.watcher.Events:
				if !ok {
					return
				}

				p.waitForAllFilesChanges()

				p.d.Logger().Infof("TLS certificates changed, updating")
				certs, _, err := p.d.Config().TLS(p.iface).Certificate()
				if err != nil {
					p.d.Logger().WithError(err).Fatalf("Error in the new tls certificates")
					return
				}
				p.load(certs)
			case err, ok := <-p.watcher.Errors:
				if !ok {
					return
				}
				p.d.Logger().WithError(err).Fatalf("Error occured in the tls certificate change watcher")
			}
		}
	}()

	certPath := path.Dir(p.certLocation.CertPath)
	keyPath := path.Dir(p.certLocation.KeyPath)

	err = p.watcher.Add(certPath)
	if err != nil {
		p.d.Logger().WithError(err).Fatalf("Error watching the certFolder for tls certificate change")
	}

	if certPath != keyPath {
		err = p.watcher.Add(keyPath)
		if err != nil {
			p.d.Logger().WithError(err).Fatalf("Error watching the keyFolder for tls certificate change")
		}
	}
}

func (p *certificatesProvider) waitForAllFilesChanges() {
	flushUntil := time.After(2 * time.Second)
	p.d.Logger().Infof("TLS certificates files changed, waiting for changes to finish")
	stop := false
	for {
		select {
		case <-flushUntil:
			stop = true
		case <-p.watcher.Events:
			continue
		}

		if stop {
			break
		}
	}
}

func (p *certificatesProvider) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if hello != nil {
		for _, cert := range p.certs {
			if cert.Leaf != nil && cert.Leaf.VerifyHostname(hello.ServerName) == nil {
				return &cert, nil
			}
		}
	}
	return &p.certs[0], nil
}

func (p *certificatesProvider) stop() {
	if p.watcher != nil {
		p.watcher.Close()
		p.watcher = nil
	}
}
