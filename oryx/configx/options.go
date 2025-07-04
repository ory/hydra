// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"

	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/logrusx"

	"github.com/knadh/koanf/v2"

	"github.com/ory/x/watcherx"
)

type (
	OptionModifier func(p *Provider)
)

func WithContext(ctx context.Context) OptionModifier {
	return func(p *Provider) {
		for _, o := range ConfigOptionsFromContext(ctx) {
			o(p)
		}
	}
}

func WithConfigFiles(files ...string) OptionModifier {
	return func(p *Provider) {
		p.files = append(p.files, files...)
	}
}

func WithImmutables(immutables ...string) OptionModifier {
	return func(p *Provider) {
		p.immutables = append(p.immutables, immutables...)
	}
}

func WithExceptImmutables(exceptImmutables ...string) OptionModifier {
	return func(p *Provider) {
		p.exceptImmutables = append(p.exceptImmutables, exceptImmutables...)
	}
}

func WithFlags(flags *pflag.FlagSet) OptionModifier {
	return func(p *Provider) {
		p.flags = flags
	}
}

func WithLogger(l *logrusx.Logger) OptionModifier {
	return func(p *Provider) {
		p.logger = l
	}
}

func SkipValidation() OptionModifier {
	return func(p *Provider) {
		p.skipValidation = true
	}
}

func DisableEnvLoading() OptionModifier {
	return func(p *Provider) {
		p.disableEnvLoading = true
	}
}

func WithValue(key string, value interface{}) OptionModifier {
	return func(p *Provider) {
		p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
	}
}

func WithValues(values map[string]interface{}) OptionModifier {
	return func(p *Provider) {
		for key, value := range values {
			p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
		}
	}
}

func WithBaseValues(values map[string]interface{}) OptionModifier {
	return func(p *Provider) {
		for key, value := range values {
			p.baseValues = append(p.baseValues, tuple{Key: key, Value: value})
		}
	}
}

func WithUserProviders(providers ...koanf.Provider) OptionModifier {
	return func(p *Provider) {
		p.userProviders = providers
	}
}

// DEPRECATED without replacement. This option is a no-op.
func OmitKeysFromTracing(keys ...string) OptionModifier {
	return func(*Provider) {}
}

func AttachWatcher(watcher func(event watcherx.Event, err error)) OptionModifier {
	return func(p *Provider) {
		p.onChanges = append(p.onChanges, watcher)
	}
}

func WithLogrusWatcher(l *logrusx.Logger) OptionModifier {
	return AttachWatcher(LogrusWatcher(l))
}

func LogrusWatcher(l *logrusx.Logger) func(e watcherx.Event, err error) {
	return func(e watcherx.Event, err error) {
		l.WithField("file", e.Source()).
			WithField("event_type", fmt.Sprintf("%T", e)).
			Info("A change to a configuration file was detected.")

		if et := new(jsonschema.ValidationError); errors.As(err, &et) {
			l.WithField("event", fmt.Sprintf("%#v", et)).
				Errorf("The changed configuration is invalid and could not be loaded. Rolling back to the last working configuration revision. Please address the validation errors before restarting the process.")
		} else if et := new(ImmutableError); errors.As(err, &et) {
			l.WithError(err).
				WithField("key", et.Key).
				WithField("old_value", fmt.Sprintf("%v", et.From)).
				WithField("new_value", fmt.Sprintf("%v", et.To)).
				Errorf("A configuration value marked as immutable has changed. Rolling back to the last working configuration revision. To reload the values please restart the process.")
		} else if err != nil {
			l.WithError(err).Errorf("An error occurred while watching config file %s", e.Source())
		} else {
			l.WithField("file", e.Source()).
				WithField("event_type", fmt.Sprintf("%T", e)).
				Info("Configuration change processed successfully.")
		}
	}
}

func WithStderrValidationReporter() OptionModifier {
	return func(p *Provider) {
		p.onValidationError = func(k *koanf.Koanf, err error) {
			p.printHumanReadableValidationErrors(k, os.Stderr, err)
		}
	}
}

func WithStandardValidationReporter(w io.Writer) OptionModifier {
	return func(p *Provider) {
		p.onValidationError = func(k *koanf.Koanf, err error) {
			p.printHumanReadableValidationErrors(k, w, err)
		}
	}
}
