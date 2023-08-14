// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ory/x/contextx"

	"github.com/rs/cors"

	"github.com/ory/x/configx"
)

const (
	KeySuffixListenOnHost           = "host"
	KeySuffixListenOnPort           = "port"
	KeySuffixSocketOwner            = "socket.owner"
	KeySuffixSocketGroup            = "socket.group"
	KeySuffixSocketMode             = "socket.mode"
	KeySuffixDisableHealthAccessLog = "request_log.disable_for_health"
)

var (
	PublicInterface ServeInterface = &servePrefix{
		prefix: "serve.public",
	}
	AdminInterface ServeInterface = &servePrefix{
		prefix: "serve.admin",
	}
)

type ServeInterface interface {
	Key(suffix string) string
	String() string
}

type servePrefix struct {
	prefix string
}

func (iface *servePrefix) Key(suffix string) string {
	if suffix == KeyRoot {
		return iface.prefix
	}
	return fmt.Sprintf("%s.%s", iface.prefix, suffix)
}

func (iface *servePrefix) String() string {
	return iface.prefix
}

func (p *DefaultProvider) ListenOn(iface ServeInterface) string {
	host, port := p.host(iface), p.port(iface)
	if strings.HasPrefix(host, "unix:") {
		return host
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func (p *DefaultProvider) SocketPermission(iface ServeInterface) *configx.UnixPermission {
	return &configx.UnixPermission{
		Owner: p.getProvider(contextx.RootContext).String(iface.Key(KeySuffixSocketOwner)),
		Group: p.getProvider(contextx.RootContext).String(iface.Key(KeySuffixSocketGroup)),
		Mode:  os.FileMode(p.getProvider(contextx.RootContext).IntF(iface.Key(KeySuffixSocketMode), 0755)),
	}
}

func (p *DefaultProvider) CORS(ctx context.Context, iface ServeInterface) (cors.Options, bool) {
	return p.getProvider(ctx).CORS(iface.Key(KeyRoot), cors.Options{
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
	})
}

func (p *DefaultProvider) DisableHealthAccessLog(iface ServeInterface) bool {
	return p.getProvider(contextx.RootContext).Bool(iface.Key(KeySuffixDisableHealthAccessLog))
}

func (p *DefaultProvider) host(iface ServeInterface) string {
	return p.getProvider(contextx.RootContext).String(iface.Key(KeySuffixListenOnHost))
}

func (p *DefaultProvider) port(iface ServeInterface) int {
	return p.getProvider(contextx.RootContext).Int(iface.Key(KeySuffixListenOnPort))
}
