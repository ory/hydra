package config

import (
	"net/url"

	"github.com/ory/x/logrusx"
)

func MustValidate(l *logrusx.Logger, p *Provider) {
	if publicTLS := p.TLS(PublicInterface); publicTLS.Enabled() {
		if p.IssuerURL().String() == "" {
			l.Fatalf(`Configuration key "%s" must be set unless flag "--dangerous-force-http" is set. To find out more, use "hydra help serve".`, KeyIssuerURL)
		}

		if p.IssuerURL().Scheme != "https" {
			l.Fatalf(`Scheme from configuration key "%s" must be "https" unless --dangerous-force-http is passed but got scheme in value "%s" is "%s". To find out more, use "hydra help serve".`, KeyIssuerURL, p.IssuerURL().String(), p.IssuerURL().Scheme)
		}

		if len(p.InsecureRedirects()) > 0 {
			l.Fatal(`Flag --dangerous-allow-insecure-redirect-urls can only be used in combination with flag --dangerous-force-http`)
		}
	}
}

func urlRoot(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}
