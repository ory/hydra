package instrumentedsql

import "context"

// Logger is the interface needed to be implemented by any logging implementation we use, see also NewFuncLogger
type Logger interface {
	Log(ctx context.Context, msg string, keyvals ...interface{})
}

type nullLogger struct{}

func (nullLogger) Log(ctx context.Context, msg string, keyvals ...interface{}) {}

// LoggerFunc is an adapter which allows a function to be used as a Logger.
type LoggerFunc func(ctx context.Context, msg string, keyvals ...interface{})

// Log calls f(ctx, msg, keyvals...).
func (f LoggerFunc) Log(ctx context.Context, msg string, keyvals ...interface{}) {
	f(ctx, msg, keyvals...)
}
