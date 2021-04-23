package config

import (
	"crypto/tls"

	"github.com/ory/x/tlsx"
)

const (
	KeySuffixTLSStrict               = "tls.strict"
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
	Strict() bool
	AllowTerminationFrom() []string
	Certificate() ([]tls.Certificate, error)
}

func (p *Provider) TLS(iface ServeInterface) TLSConfig {
	strict := false
	if !p.forcedHTTP() {
		strict = p.p.Bool(iface.Key(KeySuffixTLSStrict))
	}
	return &tlsConfig{
		strict:               strict,
		allowTerminationFrom: p.p.StringsF(iface.Key(KeySuffixTLSAllowTerminationFrom), p.p.Strings(KeyTLSAllowTerminationFrom)),

		certString: p.p.StringF(iface.Key(KeySuffixTLSCertString), p.p.String(KeyTLSCertString)),
		keyString:  p.p.StringF(iface.Key(KeySuffixTLSKeyString), p.p.String(KeyTLSKeyString)),
		certPath:   p.p.StringF(iface.Key(KeySuffixTLSCertPath), p.p.String(KeyTLSCertPath)),
		keyPath:    p.p.StringF(iface.Key(KeySuffixTLSKeyPath), p.p.String(KeyTLSKeyPath)),
	}
}

type tlsConfig struct {
	strict               bool
	allowTerminationFrom []string

	certString string
	keyString  string
	certPath   string
	keyPath    string
}

func (c *tlsConfig) Strict() bool {
	return c.strict
}

func (c *tlsConfig) AllowTerminationFrom() []string {
	return c.allowTerminationFrom
}

func (c *tlsConfig) Certificate() ([]tls.Certificate, error) {
	return tlsx.Certificate(c.certString, c.keyString, c.certPath, c.keyPath)
}

func (p *Provider) forcedHTTP() bool {
	return p.p.Bool("dangerous-force-http")
}
