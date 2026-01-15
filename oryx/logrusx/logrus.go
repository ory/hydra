// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package logrusx

import (
	"bytes"
	"cmp"
	_ "embed"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	gelf "github.com/seatgeek/logrus-gelf-formatter"

	"github.com/ory/x/stringsx"
)

type (
	options struct {
		l                         *logrus.Logger
		level                     *logrus.Level
		formatter                 logrus.Formatter
		format                    string
		reportCaller              bool
		exitFunc                  func(int)
		leakSensitive             bool
		redactionText             string
		additionalRedactedHeaders []string
		hooks                     []logrus.Hook
		c                         configurator
	}
	Option           func(*options)
	nullConfigurator struct{}
	configurator     interface {
		Bool(key string) bool
		String(key string) string
		Strings(path string) []string
	}
)

//go:embed config.schema.json
var ConfigSchema string

const ConfigSchemaID = "ory://logging-config"

// AddConfigSchema adds the logging schema to the compiler.
// The interface is specified instead of `jsonschema.Compiler` to allow the use of any jsonschema library fork or version.
func AddConfigSchema(c interface {
	AddResource(url string, r io.Reader) error
},
) error {
	return c.AddResource(ConfigSchemaID, bytes.NewBufferString(ConfigSchema))
}

func newLogger(parent *logrus.Logger, o *options) *logrus.Logger {
	l := parent
	if l == nil {
		l = logrus.New()
	}

	if o.exitFunc != nil {
		l.ExitFunc = o.exitFunc
	}

	for _, hook := range o.hooks {
		l.AddHook(hook)
	}

	setLevel(l, o)
	setFormatter(l, o)

	l.ReportCaller = o.reportCaller || l.IsLevelEnabled(logrus.TraceLevel)
	return l
}

func setLevel(l *logrus.Logger, o *options) {
	if o.level != nil {
		l.Level = *o.level
	} else {
		var err error
		l.Level, err = logrus.ParseLevel(cmp.Or(
			o.c.String("log.level"),
			os.Getenv("LOG_LEVEL")))
		if err != nil {
			l.Level = logrus.InfoLevel
		}
	}
}

func setFormatter(l *logrus.Logger, o *options) {
	if o.formatter != nil {
		l.Formatter = o.formatter
	} else {
		var unknownFormat bool // we first have to set the formatter before we can complain about the unknown format

		format := stringsx.SwitchExact(cmp.Or(o.format, o.c.String("log.format"), os.Getenv("LOG_FORMAT")))
		switch {
		case format.AddCase("json"):
			l.Formatter = &logrus.JSONFormatter{PrettyPrint: false, TimestampFormat: time.RFC3339Nano, DisableHTMLEscape: true}
		case format.AddCase("json_pretty"):
			l.Formatter = &logrus.JSONFormatter{PrettyPrint: true, TimestampFormat: time.RFC3339Nano, DisableHTMLEscape: true}
		case format.AddCase("gelf"):
			l.Formatter = new(gelf.GelfFormatter)
		default:
			unknownFormat = true
			fallthrough
		case format.AddCase("text", ""):
			l.Formatter = &logrus.TextFormatter{
				DisableQuote:     true,
				DisableTimestamp: false,
				FullTimestamp:    true,
			}
		}

		if unknownFormat {
			l.WithError(format.ToUnknownCaseErr()).Warn("got unknown \"log.format\", falling back to \"text\"")
		}
	}
}

func ForceLevel(level logrus.Level) Option {
	return func(o *options) {
		o.level = &level
	}
}

func ForceFormatter(formatter logrus.Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

func WithConfigurator(c configurator) Option {
	return func(o *options) {
		o.c = c
	}
}

func ForceFormat(format string) Option {
	return func(o *options) {
		o.format = format
	}
}

func WithHook(hook logrus.Hook) Option {
	return func(o *options) {
		o.hooks = append(o.hooks, hook)
	}
}

func WithExitFunc(exitFunc func(int)) Option {
	return func(o *options) {
		o.exitFunc = exitFunc
	}
}

func ReportCaller(reportCaller bool) Option {
	return func(o *options) {
		o.reportCaller = reportCaller
	}
}

func UseLogger(l *logrus.Logger) Option {
	return func(o *options) {
		o.l = l
	}
}

func LeakSensitive() Option {
	return func(o *options) {
		o.leakSensitive = true
	}
}

func RedactionText(text string) Option {
	return func(o *options) {
		o.redactionText = text
	}
}

func WithAdditionalRedactedHeaders(headers []string) Option {
	return func(o *options) {
		o.additionalRedactedHeaders = headers
	}
}

func toHeaderMap(headers []string) map[string]struct{} {
	m := make(map[string]struct{}, len(headers))
	for _, h := range headers {
		m[strings.ToLower(h)] = struct{}{}
	}
	return m
}

func (c *nullConfigurator) Bool(_ string) bool        { return false }
func (c *nullConfigurator) String(_ string) string    { return "" }
func (c *nullConfigurator) Strings(_ string) []string { return []string{} }

func newOptions(opts []Option) *options {
	o := new(options)
	o.c = new(nullConfigurator)
	for _, f := range opts {
		f(o)
	}
	return o
}

// New creates a new logger with all the important fields set.
func New(name string, version string, opts ...Option) *Logger {
	o := newOptions(opts)
	return &Logger{
		opts:          opts,
		leakSensitive: o.leakSensitive || o.c.Bool("log.leak_sensitive_values"),
		redactionText: cmp.Or(o.redactionText, `Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`),
		additionalRedactedHeaders: toHeaderMap(func() []string {
			if len(o.additionalRedactedHeaders) > 0 {
				return o.additionalRedactedHeaders
			}
			return o.c.Strings("log.additional_redacted_headers")
		}()),
		Entry: newLogger(o.l, o).WithFields(logrus.Fields{
			"audience": "application", "service_name": name, "service_version": version,
		}),
	}
}

func NewT(t testing.TB, opts ...Option) *Logger {
	opts = append(opts, LeakSensitive(), WithExitFunc(func(code int) {
		t.Fatalf("Logger exited with code %d", code)
	}))
	l := New(t.Name(), "test", opts...)
	l.Logger.Out = t.Output()
	return l
}

func (l *Logger) UseConfig(c configurator) {
	l.leakSensitive = l.leakSensitive || c.Bool("log.leak_sensitive_values")
	l.redactionText = cmp.Or(c.String("log.redaction_text"), l.redactionText)
	newHeaders := toHeaderMap(c.Strings("log.additional_redacted_headers"))
	for k := range newHeaders {
		l.additionalRedactedHeaders[k] = struct{}{}
	}
	o := newOptions(append(l.opts, WithConfigurator(c)))
	setLevel(l.Entry.Logger, o)
	setFormatter(l.Entry.Logger, o)
}

func (l *Logger) ReportError(r *http.Request, code int, err error, args ...interface{}) {
	logger := l.WithError(err).WithRequest(r).WithField("http_response", map[string]interface{}{
		"status_code": code,
	})
	switch {
	case code < 500:
		logger.Info(args...)
	default:
		logger.Error(args...)
	}
}
