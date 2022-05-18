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
	enabled := true
	if p.forcedHTTP(ctx) {
		enabled = false
	} else if iface == AdminInterface {
		// Support `tls.enabled` for admin interface only
		enabled = p.getProvider(ctx).Bool(iface.Key(KeySuffixTLSEnabled))
	}

	return &tlsConfig{
		enabled:              enabled,
		allowTerminationFrom: p.getProvider(ctx).StringsF(iface.Key(KeySuffixTLSAllowTerminationFrom), p.getProvider(ctx).Strings(KeyTLSAllowTerminationFrom)),

		certString: p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSCertString), p.getProvider(ctx).String(KeyTLSCertString)),
		keyString:  p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSKeyString), p.getProvider(ctx).String(KeyTLSKeyString)),
		certPath:   p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSCertPath), p.getProvider(ctx).String(KeyTLSCertPath)),
		keyPath:    p.getProvider(ctx).StringF(iface.Key(KeySuffixTLSKeyPath), p.getProvider(ctx).String(KeyTLSKeyPath)),
	}
}

func (c *tlsConfig) Certificate() ([]tls.Certificate, error) {
	return tlsx.Certificate(c.certString, c.keyString, c.certPath, c.keyPath)
}

func (p *DefaultProvider) forcedHTTP(ctx context.Context) bool {
	return p.getProvider(ctx).Bool("dangerous-force-http")
}
