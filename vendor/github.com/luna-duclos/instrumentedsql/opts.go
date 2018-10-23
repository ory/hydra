package instrumentedsql

type opts struct {
	Logger
	Tracer
	OmitArgs           bool
	TraceRowsNext      bool
}

// Opt is a functional option type for the wrapped driver
type Opt func(*opts)

// WithLogger sets the logger of the wrapped driver to the provided logger
func WithLogger(l Logger) Opt {
	return func(o *opts) {
		o.Logger = l
	}
}

// WithTracer sets the tracer of the wrapped driver to the provided tracer
func WithTracer(t Tracer) Opt {
	return func(o *opts) {
		o.Tracer = t
	}
}

// WithOmitArgs will make it so that query arguments are omitted from logging and tracing
func WithOmitArgs() Opt {
	return func(o *opts) {
		o.OmitArgs = true
	}
}

// WithIncludeArgs will make it so that query arguments are included in logging and tracing
// This is the default, but can be used to override WithOmitArgs
func WithIncludeArgs() Opt {
	return func(o *opts) {
		o.OmitArgs = false
	}
}

// WithTraceRowsNext will make it so calls to rows.Next() are traced.
// Those calls are usually incredibly brief, so are by default not traced.
func WithTraceRowsNext() Opt {
	return func(o *opts) {
		o.TraceRowsNext = true
	}
}

// WithNoTraceRowsNext will make it so calls to rows.Next() are traced.
// This is the default, but can be used to override WithTraceRowsNext
func WithNoTraceRowsNext() Opt {
	return func(o *opts) {
		o.TraceRowsNext = false
	}
}