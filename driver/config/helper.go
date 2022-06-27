package config

import (
	"context"
	"github.com/pkg/errors"
	"net/url"

	"github.com/ory/x/logrusx"
)

func Validate(ctx context.Context, l *logrusx.Logger, p *DefaultProvider) error {
	if p.IssuerURL(ctx).String() == "" && p.IsDevelopmentMode(ctx) == false {
		l.Errorf("Configuration key `%s` must be set `dev` is `true`. To find out more, use `hydra help serve`.", KeyIssuerURL)
		return errors.New("issuer URL must be set unless development mode is enabled")
	}

	if p.IssuerURL(ctx).Scheme != "https" {
		l.Errorf("Scheme from configuration key `%s` must be `https` unless `dev` is `true`. Got scheme in value `%s` is `%`. To find out more, use `hydra help serve`.", KeyIssuerURL, p.IssuerURL(ctx).String(), p.IssuerURL(ctx).Scheme)
		return errors.New("issuer URL scheme must be HTTPS unless development mode is enabled")
	}

	return nil
}

func urlRoot(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}
