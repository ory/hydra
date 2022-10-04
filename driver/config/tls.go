package config

import (
	"context"
	"crypto/tls"

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
	Certificate() ([]tls.Certificate, error)
}

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

func (c *tlsConfig) Certificate() ([]tls.Certificate, error) {
	return tlsx.Certificate(c.certString, c.keyString, c.certPath, c.keyPath)
}
