// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package servicelocatorx

import (
	"net/http"

	"github.com/ory/x/contextx"

	"google.golang.org/grpc"

	"github.com/ory/x/logrusx"
)

type (
	Options struct {
		logger                 *logrusx.Logger
		contextualizer         contextx.Contextualizer
		httpMiddlewares        []func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
		grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
		grpcStreamInterceptors []grpc.StreamServerInterceptor
	}
	Option func(o *Options)
)

func WithLogger(l *logrusx.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

func WithContextualizer(ctxer contextx.Contextualizer) Option {
	return func(o *Options) {
		o.contextualizer = ctxer
	}
}

func WithHTTPMiddlewares(m ...func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) Option {
	return func(o *Options) {
		o.httpMiddlewares = m
	}
}

func WithGRPCUnaryInterceptors(i ...grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.grpcUnaryInterceptors = i
	}
}

func WithGRPCStreamInterceptors(i ...grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.grpcStreamInterceptors = i
	}
}

func (o *Options) Logger() *logrusx.Logger {
	return o.logger
}

func (o *Options) Contextualizer() contextx.Contextualizer {
	return o.contextualizer
}

func (o *Options) HTTPMiddlewares() []func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return o.httpMiddlewares
}

func (o *Options) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return o.grpcUnaryInterceptors
}

func (o *Options) GRPCStreamInterceptors() []grpc.StreamServerInterceptor {
	return o.grpcStreamInterceptors
}

func NewOptions(options ...Option) *Options {
	o := &Options{
		contextualizer: &contextx.Default{},
	}
	for _, opt := range options {
		opt(o)
	}
	return o
}
