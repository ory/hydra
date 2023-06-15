// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"

	"github.com/ory/x/otelx"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/servicelocatorx"
)

type Options struct {
	preload  bool
	validate bool
	opts     []configx.OptionModifier
	config   *config.DefaultProvider
	// The first default refers to determining the NID at startup; the second default referes to the fact that the Contextualizer may dynamically change the NID.
	skipNetworkInit bool
	replaceTracer   func(*otelx.Tracer) *otelx.Tracer
}

func newOptions(opts []OptionsModifier) *Options {
	o := &Options{
		validate: true,
		preload:  true,
		opts:     []configx.OptionModifier{},
	}
	for _, f := range opts {
		f(o)
	}
	return o
}

func WithConfig(config *config.DefaultProvider) func(o *Options) {
	return func(o *Options) {
		o.config = config
	}
}

func ReplaceTracer(f func(*otelx.Tracer) *otelx.Tracer) func(o *Options) {
	return func(o *Options) {
		o.replaceTracer = f
	}
}

type OptionsModifier func(*Options)

func WithOptions(opts ...configx.OptionModifier) OptionsModifier {
	return func(o *Options) {
		o.opts = append(o.opts, opts...)
	}
}

// DisableValidation validating the config.
//
// This does not affect schema validation!
func DisableValidation() OptionsModifier {
	return func(o *Options) {
		o.validate = false
	}
}

// DisablePreloading will not preload the config.
func DisablePreloading() OptionsModifier {
	return func(o *Options) {
		o.preload = false
	}
}

func SkipNetworkInit() OptionsModifier {
	return func(o *Options) {
		o.skipNetworkInit = true
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
		c, err = config.New(ctx, l, o.opts...)
		if err != nil {
			l.WithError(err).Error("Unable to instantiate configuration.")
			return nil, err
		}
	}

	if o.validate {
		if err := config.Validate(ctx, l, c); err != nil {
			return nil, err
		}
	}

	r, err := NewRegistryFromDSN(ctx, c, l, o.skipNetworkInit, false, ctxter)
	if err != nil {
		l.WithError(err).Error("Unable to create service registry.")
		return nil, err
	}

	if err = r.Init(ctx, o.skipNetworkInit, false, &contextx.Default{}); err != nil {
		l.WithError(err).Error("Unable to initialize service registry.")
		return nil, err
	}

	// Avoid cold cache issues on boot:
	if o.preload {
		CallRegistry(ctx, r)
	}

	c.Source(ctx).SetTracer(ctx, r.Tracer(ctx))
	return r, nil
}
