// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package servicelocator

import (
	"context"
	"net/http"
	"testing"

	"github.com/urfave/negroni"
	"google.golang.org/grpc"

	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("case=has default contextualizer", func(t *testing.T) {
		assert.Equal(t, &contextx.Default{}, Contextualizer(context.Background(), &contextx.Default{}))
	})

	t.Run("case=overwrites contextualizer", func(t *testing.T) {
		ctxer := &struct {
			contextx.Default
			x string
		}{x: "x"}

		ctx := context.Background()
		ctx = WithContextualizer(ctx, ctxer)
		assert.Equal(t, ctxer, Contextualizer(ctx, nil))
	})

	t.Run("case=Logger", func(t *testing.T) {
		ctx := context.Background()
		expected := logrusx.New("", "")
		assert.EqualValues(t, expected, Logger(ctx, expected))
		assert.EqualValues(t, (*logrusx.Logger)(nil), Logger(ctx, nil))
		assert.EqualValues(t, expected, Logger(WithLogger(ctx, expected), nil))
	})

	t.Run("case=HTTPMiddlewares", func(t *testing.T) {
		ctx := context.Background()
		expected := []negroni.HandlerFunc{func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}}
		assert.Len(t, HTTPMiddlewares(ctx), 0)
		assert.Equal(t, expected, HTTPMiddlewares(WithHTTPMiddlewares(ctx, expected...)))
	})

	t.Run("case=GRPCStreamInterceptors", func(t *testing.T) {
		ctx := context.Background()
		expected := []grpc.StreamServerInterceptor{func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return nil
		}}
		assert.Len(t, GRPCStreamInterceptors(ctx), 0)
		assert.Equal(t, expected, GRPCStreamInterceptors(WithGRPCStreamInterceptors(ctx, expected...)))
	})

	t.Run("case=GRPCStreamInterceptors", func(t *testing.T) {
		ctx := context.Background()
		expected := []grpc.UnaryServerInterceptor{func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, nil
		}}
		assert.Len(t, GRPCUnaryInterceptors(ctx), 0)
		assert.Equal(t, expected, GRPCUnaryInterceptors(WithGRPCUnaryInterceptors(ctx, expected...)))
	})
}
