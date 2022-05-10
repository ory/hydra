package transaction

import (
	"context"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/gobuffalo/pop/v6"
	"github.com/jmoiron/sqlx"
	"github.com/ory/x/errorsx"
)

type transactionContextType string

const TransactionContextKey transactionContextType = "transactionConnection"

func Transaction(ctx context.Context, conn *pop.Connection, f func(ctx context.Context, c *pop.Connection) error) error {
	isNested := true
	c, ok := ctx.Value(TransactionContextKey).(*pop.Connection)
	if !ok {
		isNested = false

		var err error
		c, err = conn.WithContext(ctx).NewTransaction()

		if err != nil {
			return errorsx.WithStack(err)
		}
	}

	if !isNested && c.Dialect.Name() == "cockroach" { // Only retry the outer closure of cockroach transactions
		return crdb.ExecuteInTx(ctx, sqlxTxAdapter{c.TX.Tx}, func() error {
			return f(context.WithValue(ctx, TransactionContextKey, c), c)
		})
	} else {
		if err := f(context.WithValue(ctx, TransactionContextKey, c), c); err != nil {
			if !isNested {
				if err := c.TX.Rollback(); err != nil {
					return errorsx.WithStack(err)
				}
			}
			return err
		}

		// commit if there is no wrapping transaction
		if !isNested {
			return errorsx.WithStack(c.TX.Commit())
		}
	}

	return nil
}

type sqlxTxAdapter struct {
	*sqlx.Tx
}

var _ crdb.Tx = sqlxTxAdapter{}

func (s sqlxTxAdapter) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.Tx.ExecContext(ctx, query, args...)
	return err
}

func (s sqlxTxAdapter) Commit(ctx context.Context) error {
	return s.Tx.Commit()
}

func (s sqlxTxAdapter) Rollback(ctx context.Context) error {
	return s.Tx.Rollback()
}
