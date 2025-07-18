// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"

	"github.com/rs/cors"

	"github.com/ory/x/configx"
)

func (p *DefaultProvider) ServePublic(ctx context.Context) *configx.Serve {
	sharedTLS := p.getProvider(ctx).TLS("serve.tls", configx.TLS{})
	c := p.getProvider(ctx).Serve("serve.public", p.IsDevelopmentMode(ctx), configx.Serve{
		Host: "localhost",
		Port: 4444,
		TLS:  sharedTLS,
	})
	return c
}

func (p *DefaultProvider) ServeAdmin(ctx context.Context) *configx.Serve {
	sharedTLS := p.getProvider(ctx).TLS("serve.tls", configx.TLS{})
	return p.getProvider(ctx).Serve("serve.admin", p.IsDevelopmentMode(ctx), configx.Serve{
		Host: "localhost",
		Port: 4445,
		TLS:  sharedTLS,
	})
}

var defaultCORSOptions = cors.Options{
	AllowedOrigins: []string{},
	AllowedMethods: []string{
		"POST",
		"GET",
		"PUT",
		"PATCH",
		"DELETE",
		"CONNECT",
		"HEAD",
		"OPTIONS",
		"TRACE",
	},
	AllowedHeaders: []string{
		"Accept",
		"Content-Type",
		"Content-Length",
		"Accept-Language",
		"Content-Language",
		"Authorization",
	},
	ExposedHeaders: []string{
		"Cache-Control",
		"Expires",
		"Last-Modified",
		"Pragma",
		"Content-Length",
		"Content-Language",
		"Content-Type",
	},
	AllowCredentials: true,
}

func (p *DefaultProvider) CORSPublic(ctx context.Context) (cors.Options, bool) {
	return p.getProvider(ctx).CORS("serve.public", defaultCORSOptions)
}

func (p *DefaultProvider) CORSAdmin(ctx context.Context) (cors.Options, bool) {
	return p.getProvider(ctx).CORS("serve.admin", defaultCORSOptions)
}
