package config

import (
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

func (p *Provider) ListenOn(iface ServeInterface) string {
	host, port := p.host(iface), p.port(iface)
	if strings.HasPrefix(host, "unix:") {
		return host
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func (p *Provider) SocketPermission(iface ServeInterface) *configx.UnixPermission {
	return &configx.UnixPermission{
		Owner: p.p.String(iface.Key(KeySuffixSocketOwner)),
		Group: p.p.String(iface.Key(KeySuffixSocketGroup)),
		Mode:  os.FileMode(p.p.IntF(iface.Key(KeySuffixSocketMode), 0755)),
	}
}

func (p *Provider) CORS(iface ServeInterface) (cors.Options, bool) {
	return p.p.CORS(iface.Key(KeyRoot), cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
}

func (p *Provider) DisableHealthAccessLog(iface ServeInterface) bool {
	return p.p.Bool(iface.Key(KeySuffixDisableHealthAccessLog))
}

func (p *Provider) host(iface ServeInterface) string {
	return p.p.String(iface.Key(KeySuffixListenOnHost))
}

func (p *Provider) port(iface ServeInterface) int {
	return p.p.Int(iface.Key(KeySuffixListenOnPort))
}
