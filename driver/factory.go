// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"io/fs"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fositex"
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
		skipNetworkInit bool
		configOpts        []configx.OptionModifier
		config            *config.DefaultProvider
		tracerWrapper     TracerWrapper
		extraMigrations   []fs.FS
		goMigrations      []popx.Migration
		fositexFactories  []fositex.Factory
		registryModifiers []RegistryModifier
		inspect           func(Registry) error
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

func WithConfig(config *config.DefaultProvider) OptionsModifier {
	return func(o *options) {
		o.config = config
	}
}

func WithConfigOptions(opts ...configx.OptionModifier) OptionsModifier {
	return func(o *options) {
		o.configOpts = append(o.configOpts, opts...)
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

func Inspect(f func(Registry) error) OptionsModifier {
	return func(o *options) {
		o.inspect = f
	}
}

func New(ctx context.Context, sl *servicelocatorx.Options, opts []OptionsModifier) (Registry, error) {
	o := newOptions(opts)

	l := sl.Logger()
	if l == nil {
		l = logrusx.New("Ory Hydra", config.Version)
	}

	ctxter := sl.Contextualizer()
	c := o.config
	if c == nil {
		var err error
		c, err = config.New(ctx, l, o.configOpts...)
		if err != nil {
			l.WithError(err).Error("Unable to instantiate configuration.")
			return nil, err
		}
	}

	if !o.noValidate {
		if err := config.Validate(ctx, l, c); err != nil {
			return nil, err
		}
	}

	r, err := NewRegistryWithoutInit(c, l)
	if err != nil {
		l.WithError(err).Error("Unable to create service registry.")
		return nil, err
	}

	if o.tracerWrapper != nil {
		r.WithTracerWrapper(o.tracerWrapper)
	}

	r.WithExtraFositeFactories(o.fositexFactories)

	for _, f := range o.registryModifiers {
		if err := f(r); err != nil {
			return nil, err
		}
	}

	if err = r.Init(ctx, o.skipNetworkInit, false, ctxter, o.extraMigrations, o.goMigrations); err != nil {
		l.WithError(err).Error("Unable to initialize service registry.")
		return nil, err
	}

	// Avoid cold cache issues on boot:
	if !o.noPreload {
		CallRegistry(ctx, r)
	}

	if o.inspect != nil {
		if err := o.inspect(r); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return r, nil
}
