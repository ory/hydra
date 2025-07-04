// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package servicelocator

import (
	"context"

	"github.com/urfave/negroni"
	"google.golang.org/grpc"

	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
)

type contextKeyType uint8

const (
	contextKeyHTTPMiddleware contextKeyType = iota + 1
	contextKeyGRPCStreamInterceptors
	contextKeyGRPCUnaryInterceptors
	contextKeyLogger
	contextKeyContextualizer
)

func WithContextualizer(ctx context.Context, c contextx.Contextualizer) context.Context {
	return context.WithValue(ctx, contextKeyContextualizer, c)
}

func WithLogger(ctx context.Context, c *logrusx.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, c)
}

func WithHTTPMiddlewares(ctx context.Context, mws ...negroni.HandlerFunc) context.Context {
	return context.WithValue(ctx, contextKeyHTTPMiddleware, mws)
}

func WithGRPCUnaryInterceptors(ctx context.Context, mws ...grpc.UnaryServerInterceptor) context.Context {
	return context.WithValue(ctx, contextKeyGRPCUnaryInterceptors, mws)
}

func WithGRPCStreamInterceptors(ctx context.Context, mws ...grpc.StreamServerInterceptor) context.Context {
	return context.WithValue(ctx, contextKeyGRPCStreamInterceptors, mws)
}

func Logger(ctx context.Context, fallback *logrusx.Logger) *logrusx.Logger {
	if v, ok := ctx.Value(contextKeyLogger).(*logrusx.Logger); ok {
		return v
	}
	return fallback
}

func Contextualizer(ctx context.Context, fallback contextx.Contextualizer) contextx.Contextualizer {
	if v, ok := ctx.Value(contextKeyContextualizer).(contextx.Contextualizer); ok {
		return v
	}
	return fallback
}

func HTTPMiddlewares(ctx context.Context) []negroni.HandlerFunc {
	if v, ok := ctx.Value(contextKeyHTTPMiddleware).([]negroni.HandlerFunc); ok {
		return v
	}
	return []negroni.HandlerFunc{}
}

func GRPCUnaryInterceptors(ctx context.Context) []grpc.UnaryServerInterceptor {
	if v, ok := ctx.Value(contextKeyGRPCUnaryInterceptors).([]grpc.UnaryServerInterceptor); ok {
		return v
	}
	return []grpc.UnaryServerInterceptor{}
}

func GRPCStreamInterceptors(ctx context.Context) []grpc.StreamServerInterceptor {
	if v, ok := ctx.Value(contextKeyGRPCStreamInterceptors).([]grpc.StreamServerInterceptor); ok {
		return v
	}
	return []grpc.StreamServerInterceptor{}
}
