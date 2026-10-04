// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"runtime"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/x/prometheusx"

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
			// A transaction's retries are recorded in a single increment so
			// that the trace exemplar's value equals the transaction's retry
			// count, which Grafana renders as the exemplar marker's position.
			// Deferred so that retries are recorded even when the callback
			// panics.
			defer func() {
				if retries := attempt - 1; retries > 0 {
					recordRetries(ctx, retries)
				}
			}()
			return errors.WithStack(crdb.ExecuteInTx(ctx, sqlxTxAdapter{tx.TX.Tx}, func() error {
				attempt++
				return errors.WithStack(callback(WithTransaction(ctx, tx), tx))
			}))
		})
	case "postgres", "mysql", "yugabyte":
		// YugabyteDB speaks the Postgres wire protocol and, like Postgres, may
		// surface serialization failures (SQLSTATE 40001) that the retry loop
		// below handles through sqlcon.ErrConcurrentUpdate.
		//
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
	for attempt := range MaxTransactionRetries {
		err = conn.Transaction(func(tx *pop.Connection) error {
			return callback(WithTransaction(ctx, tx), tx)
		})
		if err == nil {
			return nil
		}
		if !errors.Is(sqlcon.HandleError(err), sqlcon.ErrConcurrentUpdate()) {
			return err
		}
		// Back off with jitter before retrying. SQLite in WAL mode fails a
		// deferred transaction's read-to-write upgrade with
		// SQLITE_BUSY_SNAPSHOT immediately (busy_timeout cannot help a stale
		// snapshot), so an immediate retry under write contention can livelock
		// through the whole retry budget.
		select {
		case <-ctx.Done():
			return err
		case <-time.After(time.Duration(attempt+1) * time.Duration(1+rand.IntN(3)) * time.Millisecond):
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

func (s sqlxTxAdapter) Exec(ctx context.Context, query string, args ...any) error {
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

// recordRetries counts a transaction's automatic CockroachDB restarts. When
// the surrounding request is traced, it also attaches the trace to the span
// (as an event) and to the counter (as an exemplar whose value is the retry
// count) so that retry spikes can be linked to concrete requests.
func recordRetries(ctx context.Context, retries int) {
	c := caller()
	counter := transactionRetries.WithLabelValues(c)
	span := trace.SpanFromContext(ctx)
	span.AddEvent("db.transaction.retry", trace.WithAttributes(
		attribute.String("caller", c),
		attribute.Int("retries", retries),
	))
	prometheusx.AddWithExemplar(ctx, counter, float64(retries))
}

// caller returns the function that opened the transaction: the innermost
// frame that is neither transaction plumbing (this package, pop, and
// crdb.ExecuteInTx) nor one of the thin wrapper methods that services layer
// over popx, so that a retry is attributed to the business operation rather
// than to a generic Transaction wrapper.
func caller() string {
	pc := make([]uintptr, 32)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		if frame.Function != "" && !skipCallerFrame(frame.Function) {
			return frame.Function
		}
		if !more {
			return unknownCaller
		}
	}
}

// skipCallerFrame reports whether a stack frame identifies transaction
// plumbing rather than the transaction's origin: this package and the
// libraries under it, and the thin wrapper methods services layer over popx,
// e.g. hydra's BasePersister.Transaction and RegistrySQL.Transaction
// (fosite's Transactional interface), keto's WithLatestSnapshot, and
// backoffice's WithFollowerReads.
func skipCallerFrame(fn string) bool {
	return strings.HasPrefix(fn, "github.com/ory/x/popx.") ||
		strings.HasPrefix(fn, "github.com/ory/pop/") ||
		strings.HasPrefix(fn, "github.com/cockroachdb/cockroach-go/v2/crdb.ExecuteInTx") ||
		// Deferred recording runs with runtime frames (e.g. runtime.gopanic)
		// on the stack when the transaction callback panicked.
		strings.HasPrefix(fn, "runtime.") ||
		// Wrapper matching is deliberately coarse: a frame that merely
		// mentions one of these names is skipped and the retry is attributed
		// one frame further up, which is still a useful label. "Transaction"
		// also covers TransactionWithOptions and method-value "-fm" symbols.
		strings.Contains(fn, "Transaction") ||
		strings.Contains(fn, "WithLatestSnapshot") ||
		strings.Contains(fn, "WithFollowerReads")
}
