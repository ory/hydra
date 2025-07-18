// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"cmp"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"net/url"
	"os"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/tlsx"
)

const (
	ServeConfigSchemaID = "ory://serve-config"
	TLSConfigSchemaID   = "ory://tls-config"
)

//go:embed serve.schema.json
var ServeConfigSchema []byte

//go:embed tls.schema.json
var TLSConfigSchema []byte

type (
	Serve struct {
		Host, WriteListenFile string
		Port                  int
		BaseURL               *url.URL
		Socket                UnixPermission
		TLS                   TLS
		RequestLog            ServeRequestLog
	}
	TLS struct {
		Enabled                                  bool
		AllowTerminationFrom                     []string
		CertBase64, KeyBase64, CertPath, KeyPath string
	}
	ServeRequestLog struct {
		DisableHealth bool
	}
)

func (p *Provider) Serve(prefix string, isDev bool, defaults Serve) *Serve {
	prefix = cleanPrefix(prefix)

	defaults.Socket.Mode = cmp.Or(defaults.Socket.Mode, 0o755)

	serve := Serve{
		Host:            p.StringF(prefix+"host", defaults.Host),
		Port:            p.IntF(prefix+"port", defaults.Port),
		WriteListenFile: p.StringF(prefix+"write_listen_file", defaults.WriteListenFile),
		BaseURL:         p.URIF(prefix+"base_url", defaults.BaseURL),
		Socket: UnixPermission{
			Owner: p.StringF(prefix+"socket.owner", defaults.Socket.Owner),
			Group: p.StringF(prefix+"socket.group", defaults.Socket.Group),
			Mode:  os.FileMode(p.IntF(prefix+"socket.mode", int(defaults.Socket.Mode))),
		},
		TLS: p.TLS(prefix+"tls", defaults.TLS),
		RequestLog: ServeRequestLog{
			DisableHealth: p.BoolF(prefix+"requestlog.disable_health", defaults.RequestLog.DisableHealth),
		},
	}

	if serve.BaseURL == nil {
		serve.BaseURL = &url.URL{
			Scheme: "http",
			Path:   "/",
		}
		if !isDev || serve.TLS.Enabled {
			serve.BaseURL.Scheme = "https"
		}
		host := serve.Host
		if host == "0.0.0.0" || host == "" {
			var err error
			host, err = os.Hostname()
			if err != nil {
				p.logger.WithError(err).Warn("Unable to get hostname from system, falling back to 127.0.0.1.")
				host = "127.0.0.1"
			}
		}
		serve.BaseURL.Host = fmt.Sprintf("%s:%d", host, serve.Port)
	}

	return &serve
}

func (p *Provider) TLS(prefix string, defaults TLS) TLS {
	prefix = cleanPrefix(prefix)

	return TLS{
		Enabled:              p.BoolF(prefix+"enabled", defaults.Enabled),
		AllowTerminationFrom: p.StringsF(prefix+"allow_termination_from", defaults.AllowTerminationFrom),
		CertBase64:           p.StringF(prefix+"cert.base64", defaults.CertBase64),
		KeyBase64:            p.StringF(prefix+"key.base64", defaults.KeyBase64),
		CertPath:             p.StringF(prefix+"cert.path", defaults.CertPath),
		KeyPath:              p.StringF(prefix+"key.path", defaults.KeyPath),
	}
}

func (t *TLS) GetCertFunc(ctx context.Context, l *logrusx.Logger, ifaceName string) (tlsx.CertFunc, error) {
	switch {
	case t.CertBase64 != "" && t.KeyBase64 != "":
		cert, err := tlsx.CertificateFromBase64(t.CertBase64, t.KeyBase64)
		if err != nil {
			return nil, fmt.Errorf("unable to load TLS certificate for interface %s: %w", ifaceName, err)
		}
		l.Infof("Setting up HTTPS for %s", ifaceName)
		return func(*tls.ClientHelloInfo) (*tls.Certificate, error) { return &cert, nil }, nil
	case t.CertPath != "" && t.KeyPath != "":
		errs := make(chan error, 1)
		getCert, err := tlsx.GetCertificate(ctx, t.CertPath, t.KeyPath, errs)
		if err != nil {
			return nil, fmt.Errorf("unable to load TLS certificate for interface %s: %w", ifaceName, err)
		}
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-errs:
					l.WithError(err).Error("Failed to reload TLS certificates, using previous certificates")
				}
			}
		}()
		l.Infof("Setting up HTTPS for %s (automatic certificate reloading active)", ifaceName)
		return getCert, nil
	default:
		l.Infof("TLS has not been configured for %s, skipping", ifaceName)
	}
	return nil, nil
}
