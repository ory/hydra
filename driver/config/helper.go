package config

import (
	"context"
	"net/url"

	"github.com/ory/x/logrusx"
)

func MustValidate(ctx context.Context, l *logrusx.Logger, p *DefaultProvider) {
	if publicTLS := p.TLS(ctx, PublicInterface); publicTLS.Enabled() {
		if p.IssuerURL(ctx).String() == "" {
			l.Fatalf(`Configuration key "%s" must be set unless flag "--dangerous-force-http" is set. To find out more, use "hydra help serve".`, KeyIssuerURL)
		}

		if p.IssuerURL(ctx).Scheme != "https" {
			l.Fatalf(`Scheme from configuration key "%s" must be "https" unless --dangerous-force-http is passed but got scheme in value "%s" is "%s". To find out more, use "hydra help serve".`, KeyIssuerURL, p.IssuerURL(ctx).String(), p.IssuerURL(ctx).Scheme)
		}

		if len(p.InsecureRedirects(ctx)) > 0 {
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
