// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"
	"runtime"
	"strings"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgconn"
	pgxconn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/pop/v6"
)

type transactionContextKey int

const transactionKey transactionContextKey = 0

func WithTransaction(ctx context.Context, tx *pop.Connection) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

func InTransaction(ctx context.Context) bool {
	return ctx.Value(transactionKey) != nil
}

func Transaction(ctx context.Context, connection *pop.Connection, callback func(context.Context, *pop.Connection) error) error {
	c := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := c.(*pop.Connection); ok {
			return errors.WithStack(callback(ctx, conn.WithContext(ctx)))
		}
	}

	if connection.Dialect.Name() == "cockroach" {
		return connection.WithContext(ctx).Dialect.Lock(func() error {
			transaction, err := connection.NewTransaction()
			if err != nil {
				return errors.WithStack(err)
			}

			attempt := 0
			return errors.WithStack(crdb.ExecuteInTx(ctx, sqlxTxAdapter{transaction.TX.Tx}, func() error {
				attempt++
				if attempt > 1 {
					caller := caller()
					transactionRetries.WithLabelValues(caller).Inc()
				}
				return errors.WithStack(callback(WithTransaction(ctx, transaction), transaction))
			}))
		})
	}
	if strings.HasPrefix(connection.Dialect.Name(), "postgres") {
		var err error
		for attempt := range 10 {
			_ = attempt
			err = connection.WithContext(ctx).Transaction(func(tx *pop.Connection) error {
				return callback(WithTransaction(ctx, tx), tx)
			})
			if err == nil {
				return nil
			}
			if e := new(pgconn.PgError); errors.As(err, &e) && e.Code == "40001" {
				continue
			}
			if e := new(pgxconn.PgError); errors.As(err, &e) && e.Code == "40001" {
				continue
			}
			return err
		}
		return err
	}

	return errors.WithStack(connection.WithContext(ctx).Transaction(func(tx *pop.Connection) error {
		return errors.WithStack(callback(WithTransaction(ctx, tx), tx))
	}))
}

func GetConnection(ctx context.Context, connection *pop.Connection) *pop.Connection {
	c := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := c.(*pop.Connection); ok {
			return conn.WithContext(ctx)
		}
	}
	return connection.WithContext(ctx)
}

type sqlxTxAdapter struct {
	*sqlx.Tx
}

var _ crdb.Tx = sqlxTxAdapter{}

func (s sqlxTxAdapter) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.Tx.ExecContext(ctx, query, args...)
	return errors.WithStack(err)
}

func (s sqlxTxAdapter) Commit(ctx context.Context) error {
	return errors.WithStack(s.Tx.Commit())
}

func (s sqlxTxAdapter) Rollback(ctx context.Context) error {
	return errors.WithStack(s.Tx.Rollback())
}

var (
	transactionRetries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ory_x_popx_cockroach_transaction_retries_total",
		Help: "Counts the number of automatic CockroachDB transaction retries",
	}, []string{"caller"})
	TransactionRetries prometheus.Collector = transactionRetries
	_                                       = transactionRetries.WithLabelValues(unknownCaller) // make sure the metric is always present
	unknownCaller                           = "unknown"
)

func caller() string {
	pc := make([]uintptr, 3)
	// The number stack frames to skip was determined by putting a breakpoint in
	// ory/kratos and looking for the topmost frame which isn't from ory/x or
	// ory/pop.
	n := runtime.Callers(8, pc)
	if n == 0 {
		return unknownCaller
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		if frame.Function != "" {
			return frame.Function
		}
		if !more {
			break
		}
	}
	return unknownCaller
}
