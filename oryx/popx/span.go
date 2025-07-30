// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const (
	MigrationStatusOpName         = "migration-status"
	MigrationInitOpName           = "migration-init"
	MigrationUpOpName             = "migration-up"
	MigrationRunTransactionOpName = "migration-run-transaction"
	MigrationDownOpName           = "migration-down"
)

func startSpan(ctx context.Context, opName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return trace.SpanFromContext(ctx).TracerProvider().Tracer(tracingComponent).Start(ctx, opName, opts...)
}
