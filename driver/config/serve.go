package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/cors"

	"github.com/ory/x/configx"
)

const (
	KeySuffixListenOnHost           = "host"
	KeySuffixListenOnPort           = "port"
	KeySuffixSocketOwner            = "socket.owner"
	KeySuffixSocketGroup            = "socket.group"
	KeySuffixSocketMode             = "socket.mode"
	KeySuffixDisableHealthAccessLog = "access_log.disable_for_health"
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

func (p *DefaultProvider) ListenOn(ctx context.Context, iface ServeInterface) string {
	host, port := p.host(ctx, iface), p.port(ctx, iface)
	if strings.HasPrefix(host, "unix:") {
		return host
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func (p *DefaultProvider) SocketPermission(ctx context.Context, iface ServeInterface) *configx.UnixPermission {
	return &configx.UnixPermission{
		Owner: p.getProvider(ctx).String(iface.Key(KeySuffixSocketOwner)),
		Group: p.getProvider(ctx).String(iface.Key(KeySuffixSocketGroup)),
		Mode:  os.FileMode(p.getProvider(ctx).IntF(iface.Key(KeySuffixSocketMode), 0755)),
	}
}

func (p *DefaultProvider) CORS(ctx context.Context, iface ServeInterface) (cors.Options, bool) {
	return p.getProvider(ctx).CORS(iface.Key(KeyRoot), cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
}

func (p *DefaultProvider) DisableHealthAccessLog(ctx context.Context, iface ServeInterface) bool {
	return p.getProvider(ctx).Bool(iface.Key(KeySuffixDisableHealthAccessLog))
}

func (p *DefaultProvider) host(ctx context.Context, iface ServeInterface) string {
	return p.getProvider(ctx).String(iface.Key(KeySuffixListenOnHost))
}

func (p *DefaultProvider) port(ctx context.Context, iface ServeInterface) int {
	return p.getProvider(ctx).Int(iface.Key(KeySuffixListenOnPort))
}
