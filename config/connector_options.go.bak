package config

type ConnectorOptions func(*Options)

type Options struct {
	UseTracing bool

	// unexported as these are specifically used by Hydra for writing tests
	useRandomDriverName   bool
	allowRootTracingSpans bool
	omitSQLArgsFromSpans  bool
}

// WithTracing will make it so that a wrapped driver is used that supports the OpenTracing API
func WithTracing() ConnectorOptions {
	return func(o *Options) {
		o.UseTracing = true
	}
}

// this option is specifically for writing tests as you can't register a driver with the same name more than once
func withRandomDriverName() ConnectorOptions {
	return func(o *Options) {
		o.useRandomDriverName = true
	}
}

// withAllowRoot will make it so that root spans will be created if a trace could not be found in the context
func withAllowRootTraceSpans() ConnectorOptions {
	return func(o *Options) {
		o.allowRootTracingSpans = true
	}
}

// withOmitSQLArgsFromSpans will make it so that query arguments are omitted from tracing spans
func withOmitSQLArgsFromSpans() ConnectorOptions {
	return func(o *Options) {
		o.omitSQLArgsFromSpans = true
	}
}
