package config

type ConnectorOptions func(*options)

type options struct {
	UseTracing                  bool
	OmitSQLArgsFromTracingSpans bool
	UseRandomDriverName         bool
	AllowRootTracingSpans       bool
}

// WithTracing will make it so that a wrapped driver is used that supports the OpenTracing API
func WithTracing() ConnectorOptions {
	return func(o *options) {
		o.UseTracing = true
	}
}

// WithOmitSQLArgsFromSpans will make it so that query arguments are omitted from tracing spans
func WithOmitSQLArgsFromSpans() ConnectorOptions {
	return func(o *options) {
		o.OmitSQLArgsFromTracingSpans = true
	}
}

// This option is specifically for writing tests as you can't register a driver with the same name more than once
func WithRandomDriverName() ConnectorOptions {
	return func(o *options) {
		o.UseRandomDriverName = true
	}
}

// WithAllowRoot will make it so that root spans will be created if a trace could not be found in the context
func WithAllowRootTraceSpans() ConnectorOptions {
	return func(o *options) {
		o.AllowRootTracingSpans = true
	}
}
