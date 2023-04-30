// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/tlsutil"
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

	KeyClientTLSInsecureSkipVerify = "tls.insecure_skip_verify"
	KeySuffixClientTLSCipherSuites = "tls.cipher_suites"
	KeySuffixClientTLSMinVer       = "tls.min_version"
	KeySuffixClientTLSMaxVer       = "tls.max_version"
)

type ClientInterface interface {
	Key(suffix string) string
}

func (iface *clientPrefix) Key(suffix string) string {
	return fmt.Sprintf("%s.%s", iface.prefix, suffix)
}

type clientPrefix struct {
	prefix string
}

var (
	KeyPrefixClientDefault ClientInterface = &clientPrefix{
		prefix: "client.default",
	}
	KeyPrefixClientBackChannelLogout ClientInterface = &clientPrefix{
		prefix: "client.back_channel_logout",
	}
	KeyPrefixClientRefreshTokenHook ClientInterface = &clientPrefix{
		prefix: "client.refresh_token_hook",
	}
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

func (p *DefaultProvider) TLSClientConfigDefault() (*tls.Config, error) {
	return p.TLSClientConfigWithDefaultFallback(KeyPrefixClientDefault)
}

func (p *DefaultProvider) TLSClientConfigWithDefaultFallback(iface ClientInterface) (*tls.Config, error) {
	tlsClientConfig := new(tls.Config)
	tlsClientConfig.InsecureSkipVerify = p.p.BoolF(KeyClientTLSInsecureSkipVerify, false)

	if p.p.Exists(KeyPrefixClientDefault.Key(KeySuffixClientTLSCipherSuites)) || p.p.Exists(iface.Key(KeySuffixClientTLSCipherSuites)) {
		keyCipherSuites := p.p.StringsF(iface.Key(KeySuffixClientTLSCipherSuites), p.p.Strings(KeyPrefixClientDefault.Key(KeySuffixClientTLSCipherSuites)))
		cipherSuites, err := tlsutil.ParseCiphers(strings.Join(keyCipherSuites[:], ","))
		if err != nil {
			return nil, errors.WithMessage(err, "Unable to setup client TLS configuration")
		}
		tlsClientConfig.CipherSuites = cipherSuites
	}

	if p.p.Exists(KeyPrefixClientDefault.Key(KeySuffixClientTLSMinVer)) || p.p.Exists(iface.Key(KeySuffixClientTLSMinVer)) {
		keyMinVer := p.p.StringF(iface.Key(KeySuffixClientTLSMinVer), p.p.String(KeyPrefixClientDefault.Key(KeySuffixClientTLSMinVer)))
		if tlsMinVer, found := tlsutil.TLSLookup[keyMinVer]; !found {
			return nil, errors.Errorf("Unable to setup client TLS configuration. Invalid minimum TLS version: %s", keyMinVer)
		} else {
			tlsClientConfig.MinVersion = tlsMinVer
		}
	}

	if p.p.Exists(KeyPrefixClientDefault.Key(KeySuffixClientTLSMaxVer)) || p.p.Exists(iface.Key(KeySuffixClientTLSMaxVer)) {
		keyMaxVer := p.p.StringF(iface.Key(KeySuffixClientTLSMaxVer), p.p.String(KeyPrefixClientDefault.Key(KeySuffixClientTLSMaxVer)))
		if tlsMaxVer, found := tlsutil.TLSLookup[keyMaxVer]; !found {
			return nil, errors.Errorf("Unable to setup client TLS configuration. Invalid maximum TLS version: %s", keyMaxVer)
		} else {
			tlsClientConfig.MaxVersion = tlsMaxVer
		}
	}

	return tlsClientConfig, nil
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
