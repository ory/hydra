// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"crypto/tls"

	"github.com/pkg/errors"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/tlsx"
)

const (
	KeySuffixTLSEnabled              = "tls.enabled"
	KeySuffixTLSAllowTerminationFrom = "tls.allow_termination_from"
	KeySuffixTLSCertString           = "tls.cert.base64"
	KeySuffixTLSKeyString            = "tls.key.base64"
	KeySuffixTLSCertPath             = "tls.cert.path"
	KeySuffixTLSKeyPath              = "tls.key.path"

	KeyTLSAllowTerminationFrom = "serve." + KeySuffixTLSAllowTerminationFrom
	KeyTLSCertString           = "serve." + KeySuffixTLSCertString
	KeyTLSKeyString            = "serve." + KeySuffixTLSKeyString
	KeyTLSCertPath             = "serve." + KeySuffixTLSCertPath
	KeyTLSKeyPath              = "serve." + KeySuffixTLSKeyPath
	KeyTLSEnabled              = "serve." + KeySuffixTLSEnabled
)

type TLSConfig interface {
	Enabled() bool
	AllowTerminationFrom() []string
	GetCertificateFunc(stopReload <-chan struct{}, _ *logrusx.Logger) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), error)
}

var _ TLSConfig = (*tlsConfig)(nil)

type tlsConfig struct {
	enabled              bool
	allowTerminationFrom []string

	certString string
	keyString  string
	certPath   string
	keyPath    string
}

func (c *tlsConfig) Enabled() bool {
	return c.enabled
}

func (c *tlsConfig) AllowTerminationFrom() []string {
	return c.allowTerminationFrom
}

func (p *DefaultProvider) TLS(ctx context.Context, iface ServeInterface) TLSConfig {
	return &tlsConfig{
		enabled:              p.getProvider(ctx).BoolF(iface.Key(KeySuffixTLSEnabled), p.getProvider(ctx).Bool(KeyTLSEnabled)),
		allowTerminationFrom: p.getProvider(ctx).StringsF(iface.Key(KeySuffixTLSAllowTerminationFrom), p.getProvider(ctx).Strings(KeyTLSAllowTerminationFrom)),
		certString:           p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSCertString), p.getProvider(ctx).String(KeyTLSCertString)),
		keyString:            p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSKeyString), p.getProvider(ctx).String(KeyTLSKeyString)),
		certPath:             p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSCertPath), p.getProvider(ctx).String(KeyTLSCertPath)),
		keyPath:              p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSKeyPath), p.getProvider(ctx).String(KeyTLSKeyPath)),
	}
}

func (c *tlsConfig) GetCertificateFunc(stopReload <-chan struct{}, log *logrusx.Logger) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), error) {
	if c.certPath != "" && c.keyPath != "" { // attempt to load from disk first (enables hot-reloading)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-stopReload
			cancel()
		}()
		errs := make(chan error, 1)
		getCert, err := tlsx.GetCertificate(ctx, c.certPath, c.keyPath, errs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		go func() {
			for err := range errs {
				log.WithError(err).Error("Failed to reload TLS certificates. Using the previously loaded certificates.")
			}
		}()
		return getCert, nil
	}
	if c.certString != "" && c.keyString != "" { // base64-encoded directly in config
		cert, err := tlsx.CertificateFromBase64(c.certString, c.keyString)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			return &cert, nil
		}, nil
	}
	return nil, tlsx.ErrNoCertificatesConfigured
}
