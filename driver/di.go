// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fositex"
	"github.com/ory/hydra/v2/hsm"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
)

// WritableRegistry is a deprecated interface that should not be used anymore.
//
// Deprecate this at some point.
type WritableRegistry interface {
	// WithBuildInfo(v, h, d string) Registry

	WithConfig(c *config.DefaultProvider) Registry
	WithContextualizer(ctxer contextx.Contextualizer) Registry
	WithLogger(l *logrusx.Logger) Registry
	WithTracer(t trace.Tracer) Registry
	WithTracerWrapper(TracerWrapper) Registry
	WithKratos(k kratos.Client) Registry
	WithExtraFositeFactories(f []fositex.Factory) Registry
	ExtraFositeFactories() []fositex.Factory
	WithOAuth2Provider(f fosite.OAuth2Provider)
	WithConsentStrategy(c consent.Strategy)
	WithHsmContext(h hsm.Context)
}

type RegistryModifier func(r Registry) error

func WithRegistryModifiers(f ...RegistryModifier) OptionsModifier {
	return func(o *Options) {
		o.registryModifiers = f
	}
}

func RegistryWithHMACSHAStrategy(s func(r Registry) oauth2.CoreStrategy) RegistryModifier {
	return func(r Registry) error {
		switch rt := r.(type) {
		case *RegistrySQL:
			rt.hmacs = s(r)
		default:
			return errors.Errorf("unable to set HMAC strategy on registry of type %T", r)
		}
		return nil
	}
}

func RegistryWithHsmContext(h hsm.Context) RegistryModifier {
	return func(r Registry) error {
		switch rt := r.(type) {
		case *RegistrySQL:
			rt.hsm = h
		default:
			return errors.Errorf("unable to set HMAC strategy on registry of type %T", r)
		}
		return nil
	}
}
