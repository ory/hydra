package driver

import (
	"context"

	"github.com/spf13/pflag"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/driver/config"
)

type options struct {
	forcedValues map[string]interface{}
	preload      bool
	validate     bool
}

func newOptions() *options {
	return &options{
		forcedValues: make(map[string]interface{}),
		validate:     true,
		preload:      true,
	}
}

type OptionsModifier func(*options)

// ForceConfigValue overrides any config values set by one of the providers.
func ForceConfigValue(key string, value interface{}) OptionsModifier {
	return func(o *options) {
		o.forcedValues[key] = value
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

// DisableValidation validating the config.
//
// This does not affect schema validation!
func DisablePreloading() OptionsModifier {
	return func(o *options) {
		o.preload = false
	}
}

func New(flags *pflag.FlagSet, opts ...OptionsModifier) Registry {
	o := newOptions()
	for _, f := range opts {
		f(o)
	}

	l := logrusx.New("ORY Hydra", config.Version)
	c, err := config.New(flags, l)
	if err != nil {
		l.WithError(err).Fatal("Unable to instantiate service registry.")
	}
	l.UseConfig(c.Source())

	for k, v := range o.forcedValues {
		c.Set(k, v)
	}

	if o.validate {
		config.MustValidate(l, c)
	}

	r, err := NewRegistryFromDSN(c, l)
	if err != nil {
		l.WithError(err).Fatal("Unable to instantiate service registry.")
	}

	if err = r.Init(); err != nil {
		l.WithError(err).Fatal("Unable to initialize service registry.")
	}

	// Avoid cold cache issues on boot:
	if o.preload {
		CallRegistry(r)
	}

	c.Source().SetTracer(context.Background(), r.Tracer())

	return r
}
