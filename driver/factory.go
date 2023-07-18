// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/servicelocatorx"
)

type (
	options struct {
		preload  bool
		validate bool
		opts     []configx.OptionModifier
		config   *config.DefaultProvider
		// The first default refers to determining the NID at startup; the second default referes to the fact that the Contextualizer may dynamically change the NID.
		skipNetworkInit bool
		tracerWrapper   TracerWrapper
	}
	TracerWrapper func(*otelx.Tracer) *otelx.Tracer
)

func newOptions() *options {
	return &options{
		validate: true,
		preload:  true,
		opts:     []configx.OptionModifier{},
	}
}

func WithConfig(config *config.DefaultProvider) func(o *options) {
	return func(o *options) {
		o.config = config
	}
}

type OptionsModifier func(*options)

func WithOptions(opts ...configx.OptionModifier) OptionsModifier {
	return func(o *options) {
		o.opts = append(o.opts, opts...)
	}
}

// DisableValidation validating the config.
//
// This does not affect schema validation!
func DisableValidation() OptionsModifier {
	return func(o *options) {
		o.validate = false
	}
}

// DisablePreloading will not preload the config.
func DisablePreloading() OptionsModifier {
	return func(o *options) {
		o.preload = false
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

func New(ctx context.Context, sl *servicelocatorx.Options, opts []OptionsModifier) (Registry, error) {
	o := newOptions()
	for _, f := range opts {
		f(o)
	}

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

	r, err := NewRegistryWithoutInit(c, l)
	if err != nil {
		l.WithError(err).Error("Unable to create service registry.")
		return nil, err
	}

	if o.tracerWrapper != nil {
		r.WithTracerWrapper(o.tracerWrapper)
	}

	if err = r.Init(ctx, o.skipNetworkInit, false, ctxter); err != nil {
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
