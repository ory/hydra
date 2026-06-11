// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/pop/v6"
	"github.com/ory/x/sqlcon"
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
	return TransactionWithOptions(ctx, connection, nil, callback)
}

// TransactionWithOptions opens the transaction with given sql.TxOptions, allowing isolation level to be set.
func TransactionWithOptions(ctx context.Context, connection *pop.Connection, opts *sql.TxOptions, callback func(context.Context, *pop.Connection) error) error {
	c := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := c.(*pop.Connection); ok {
			return errors.WithStack(callback(ctx, conn.WithContext(ctx)))
		}
	}

	conn := connection.WithContext(ctx)

	switch conn.Dialect.Name() {
	case "cockroach":
		return conn.Dialect.Lock(func() error {
			tx, err := conn.NewTransactionContextOptions(ctx, opts)
			if err != nil {
				return errors.WithStack(err)
			}
			attempt := 0
			return errors.WithStack(crdb.ExecuteInTx(ctx, sqlxTxAdapter{tx.TX.Tx}, func() error {
				attempt++
				if attempt > 1 {
					c := caller()
					transactionRetries.WithLabelValues(c).Inc()
				}
				return errors.WithStack(callback(WithTransaction(ctx, tx), tx))
			}))
		})
	case "postgres", "mysql":
		// Mirrors pop's Connection#Transaction with opts passed to NewTransactionContextOptions.
		// https://github.com/ory/pop/blob/89126558d36911217a1cc70172ba94ee32692cad/connection.go#L148
		return conn.Dialect.Lock(func() error {
			var err error
			for range MaxTransactionRetries {
				err = func() error {
					cn, err := conn.NewTransactionContextOptions(ctx, opts)
					if err != nil {
						return errors.WithStack(err)
					}
					defer func() {
						if ex := recover(); ex != nil {
							_ = cn.TX.Rollback()
							panic(ex)
						}
					}()
					err = callback(WithTransaction(ctx, cn), cn)
					var dberr error
					if err != nil {
						dberr = cn.TX.Rollback()
						if errors.Is(dberr, sql.ErrTxDone) {
							// Already rolled back by the database (e.g. context cancelled).
							return err
						}
						if dberr != nil && dberr.Error() == "conn closed" {
							// pgx closes the connection on context cancellation before
							// database/sql gets a chance to roll back.
							// See https://github.com/jackc/pgx/issues/2551
							return err
						}
					} else {
						dberr = cn.TX.Commit()
					}
					if dberr != nil {
						return fmt.Errorf("database error on committing or rolling back transaction: %w", dberr)
					}
					return err
				}()
				if err == nil || !errors.Is(sqlcon.HandleError(err), sqlcon.ErrConcurrentUpdate()) {
					return err
				}
			}
			return err
		})
	}

	// SQLite and unknown dialects: opts are ignored; use pop's default
	// transaction path with concurrent-update retry handling.
	var err error
	for range MaxTransactionRetries {
		err = conn.Transaction(func(tx *pop.Connection) error {
			return callback(WithTransaction(ctx, tx), tx)
		})
		if err == nil {
			return nil
		}
		if !errors.Is(sqlcon.HandleError(err), sqlcon.ErrConcurrentUpdate()) {
			return err
		}
	}
	return err
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

// MaxTransactionRetries is the number of times a transaction is retried on
// a concurrent-update conflict before the error is returned to the caller.
const MaxTransactionRetries = 10

var (
	transactionRetries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ory_x_popx_cockroach_transaction_retries_total",
		Help: "Counts the number of automatic CockroachDB transaction retries",
	}, []string{"caller"})
	TransactionRetries prometheus.Collector = transactionRetries
	_                                       = transactionRetries.WithLabelValues(unknownCaller) // make sure the metric is always present
	unknownCaller                           = "unknown"
)

// caller returns the external caller of TransactionWithOptions.
// It skips 8 frames to land just outside the crdb/popx call stack, then
// returns the first frame that is not in this package. The extra scan
// handles the case where Transaction func is the next frame.
func caller() string {
	pc := make([]uintptr, 3)
	n := runtime.Callers(8, pc)
	if n == 0 {
		return unknownCaller
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		if frame.Function != "" && !strings.HasPrefix(frame.Function, "github.com/ory/x/popx.") {
			return frame.Function
		}
		if !more {
			break
		}
	}
	return unknownCaller
}
