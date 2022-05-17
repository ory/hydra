package driver

import (
	"context"

	"github.com/ory/x/configx"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/contextx"
)

type options struct {
	forcedValues map[string]interface{}
	preload      bool
	validate     bool
	opts         []configx.OptionModifier
	// The first default refers to determining the NID at startup; the second default referes to the fact that the Contextualizer may dynamically change the NID.
	skipNetworkInit bool
}

func newOptions() *options {
	return &options{
		validate: true,
		preload:  true,
		opts:     []configx.OptionModifier{},
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

func New(ctx context.Context, opts ...OptionsModifier) Registry {
	o := newOptions()
	for _, f := range opts {
		f(o)
	}

	l := logrusx.New("Ory Hydra", config.Version)
	c, err := config.New(ctx, l, o.opts...)
	if err != nil {
		l.WithError(err).Fatal("Unable to instantiate configuration.")
	}

	if o.validate {
		config.MustValidate(ctx, l, c)
	}

	r, err := NewRegistryFromDSN(ctx, c, l, o.skipNetworkInit, false, &contextx.Default{})
	if err != nil {
		l.WithError(err).Fatal("Unable to create service registry.")
	}

	if err = r.Init(ctx, o.skipNetworkInit, false, &contextx.Default{}); err != nil {
		l.WithError(err).Fatal("Unable to initialize service registry.")
	}

	// Avoid cold cache issues on boot:
	if o.preload {
		CallRegistry(ctx, r)
	}

	c.Source(ctx).SetTracer(ctx, r.Tracer(ctx))

	return r
}
