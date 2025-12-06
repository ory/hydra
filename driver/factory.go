// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"io/fs"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fositex"
	"github.com/ory/hydra/v2/hsm"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/popx"
	"github.com/ory/x/servicelocatorx"
)

type (
	options struct {
		noPreload,
		noValidate,
		autoMigrate,
		skipNetworkInit bool
		configOpts         []configx.OptionModifier
		tracerWrapper      TracerWrapper
		extraMigrations    []fs.FS
		goMigrations       []popx.Migration
		fositexFactories   []fositex.Factory
		registryModifiers  []RegistryModifier
		inspect            func(*RegistrySQL) error
		serviceLocatorOpts []servicelocatorx.Option
		hsmContext         hsm.Context
		kratos             kratos.Client
		fop                fosite.OAuth2Provider
		dbOptsModifier     []func(details *pop.ConnectionDetails)
	}
	OptionsModifier func(*options)

	TracerWrapper func(*otelx.Tracer) *otelx.Tracer
)

func newOptions(opts []OptionsModifier) *options {
	o := &options{}
	for _, f := range opts {
		f(o)
	}
	return o
}

func WithConfigOptions(opts ...configx.OptionModifier) OptionsModifier {
	return func(o *options) {
		o.configOpts = append(o.configOpts, opts...)
	}
}

// WithDBOptionsModifier modifies the pop connection details before the connection is opened.
func WithDBOptionsModifier(f ...func(details *pop.ConnectionDetails)) OptionsModifier {
	return func(o *options) {
		o.dbOptsModifier = append(o.dbOptsModifier, f...)
	}
}

// DisableValidation validating the config.
//
// This does not affect schema validation!
func DisableValidation() OptionsModifier {
	return func(o *options) {
		o.noValidate = true
	}
}

// DisablePreloading will not preload the config.
func DisablePreloading() OptionsModifier {
	return func(o *options) {
		o.noPreload = true
	}
}

func SkipNetworkInit() OptionsModifier {
	return func(o *options) {
		o.skipNetworkInit = true
	}
}

// WithTracerWrapper sets a function that wraps the tracer.
func WithTracerWrapper(wrapper TracerWrapper) OptionsModifier {
	return func(o *options) {
		o.tracerWrapper = wrapper
	}
}

// WithExtraMigrations specifies additional database migration.
func WithExtraMigrations(m ...fs.FS) OptionsModifier {
	return func(o *options) {
		o.extraMigrations = append(o.extraMigrations, m...)
	}
}

func WithGoMigrations(m ...popx.Migration) OptionsModifier {
	return func(o *options) {
		o.goMigrations = append(o.goMigrations, m...)
	}
}

func WithExtraFositeFactories(f ...fositex.Factory) OptionsModifier {
	return func(o *options) {
		o.fositexFactories = append(o.fositexFactories, f...)
	}
}

func Inspect(f func(*RegistrySQL) error) OptionsModifier {
	return func(o *options) {
		o.inspect = f
	}
}

func WithServiceLocatorOptions(opts ...servicelocatorx.Option) OptionsModifier {
	return func(o *options) {
		o.serviceLocatorOpts = append(o.serviceLocatorOpts, opts...)
	}
}

func WithAutoMigrate() OptionsModifier {
	return func(o *options) {
		o.autoMigrate = true
	}
}

func WithHSMContext(h hsm.Context) OptionsModifier {
	return func(o *options) {
		o.hsmContext = h
	}
}

func WithKratosClient(k kratos.Client) OptionsModifier {
	return func(o *options) {
		o.kratos = k
	}
}

func WithOAuth2Provider(p fosite.OAuth2Provider) OptionsModifier {
	return func(o *options) {
		o.fop = p
	}
}

func New(ctx context.Context, opts ...OptionsModifier) (*RegistrySQL, error) {
	o := newOptions(opts)
	sl := servicelocatorx.NewOptions(o.serviceLocatorOpts...)

	l := sl.Logger()
	if l == nil {
		l = logrusx.New("Ory Hydra", config.Version)
	}

	c, err := config.New(ctx, l, sl.Contextualizer(), o.configOpts...)
	if err != nil {
		l.WithError(err).Error("Unable to instantiate configuration.")
		return nil, err
	}

	if !o.noValidate {
		if err := config.Validate(ctx, l, c); err != nil {
			return nil, err
		}
	}

	r, err := newRegistryWithoutInit(c, l)
	if err != nil {
		l.WithError(err).Error("Unable to create service registry.")
		return nil, err
	}

	r.tracerWrapper = o.tracerWrapper
	r.fositeFactories = o.fositexFactories
	r.hsm = o.hsmContext
	r.middlewares = sl.HTTPMiddlewares()
	r.ctxer = sl.Contextualizer()
	r.kratos = o.kratos
	r.fop = o.fop
	r.dbOptsModifier = o.dbOptsModifier

	if err = r.Init(ctx, o.skipNetworkInit, o.autoMigrate, o.extraMigrations, o.goMigrations); err != nil {
		l.WithError(err).Error("Unable to initialize service registry.")
		return nil, err
	}

	for _, f := range o.registryModifiers {
		if err := f(r); err != nil {
			return nil, err
		}
	}

	// Avoid cold cache issues on boot:
	if !o.noPreload {
		callRegistry(ctx, r)
	}

	if o.inspect != nil {
		if err := o.inspect(r); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return r, nil
}
