package transaction

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/gobuffalo/pop/v6"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlcon"
)

func newDB(t *testing.T) *pop.Connection {
	if runtime.GOOS == "windows" {
		t.Skip("CockroachDB test suite does not support windows")
	}

	ts, err := testserver.NewTestServer()
	require.NoError(t, err)

	dsn := ts.PGURL()
	dsn.Scheme = "cockroach:"
	q := dsn.Query()
	q.Set("search_path", "d,public")
	dsn.RawQuery = q.Encode()

	c, err := pop.NewConnection(&pop.ConnectionDetails{URL: dsn.String()})
	require.NoError(t, err)
	require.NoError(t, c.Open())
	return c
}

func TestCockroachTransactionRetryExpectedFailure(t *testing.T) {
	c := newDB(t)
	require.Error(t, crdb.ExecuteTxGenericTest(context.Background(), popWriteSkewTest{c: c, t: t}))
}

func TestCockroachTransactionRetrySuccess(t *testing.T) {
	c := newDB(t)
	require.NoError(t, crdb.ExecuteTxGenericTest(context.Background(), hydraWriteSkewTest{c: c, popWriteSkewTest: popWriteSkewTest{c: c, t: t}}))
}

type table struct {
	ID      int `db:"id"`
	Balance int `db:"balance"`
}

func (t table) TableName() string {
	return "t"
}

type popWriteSkewTest struct {
	t *testing.T
	c *pop.Connection
}

type hydraWriteSkewTest struct {
	popWriteSkewTest
	c *pop.Connection
}

var _ crdb.WriteSkewTest = popWriteSkewTest{}
var _ crdb.WriteSkewTest = hydraWriteSkewTest{}

// ExecuteTx is part of the crdb.WriteSkewTest interface.
func (t hydraWriteSkewTest) ExecuteTx(ctx context.Context, fn func(tx interface{}) error) error {
	return Transaction(ctx, t.c, func(ctx context.Context, tx *pop.Connection) error {
		return fn(tx.WithContext(ctx))
	})
}

func (t popWriteSkewTest) Init(ctx context.Context) error {
	for _, s := range []string{
		"CREATE DATABASE d",
		"CREATE TABLE d.t (id INT PRIMARY KEY, balance INT)",
		"USE d",
		"INSERT INTO d.t (id, balance) VALUES (1, 100), (2, 100)",
	} {
		if err := t.c.RawQuery(s).Exec(); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteTx is part of the crdb.WriteSkewTest interface.
func (t popWriteSkewTest) ExecuteTx(ctx context.Context, fn func(tx interface{}) error) error {
	fmt.Printf("entering...\n")
	return t.c.Transaction(func(tx *pop.Connection) error {
		return fn(tx)
	})
}

// GetBalances is part of the crdb.WriteSkewTest interface.
func (t popWriteSkewTest) GetBalances(ctx context.Context, txi interface{}) (int, int, error) {
	tx := txi.(*pop.Connection).WithContext(ctx)
	var tables []table

	err := tx.RawQuery(`SELECT * FROM d.t WHERE id IN (1, 2);`).All(&tables)
	if err != nil {
		return 0, 0, sqlcon.HandleError(err)
	}

	if len(tables) != 2 {
		err := fmt.Errorf("expected two balances; got %d", len(tables))
		t.t.Logf("Got error: %+v", err)
		return 0, 0, err
	}
	return tables[0].Balance, tables[1].Balance, nil
}

// UpdateBalance is part of the crdb.WriteSkewInterface.
func (t popWriteSkewTest) UpdateBalance(
	ctx context.Context, txi interface{}, acct, delta int,
) error {
	tx := txi.(*pop.Connection).WithContext(ctx)
	err := tx.RawQuery(`UPDATE d.t SET balance=balance+$1 WHERE id=$2;`, delta, acct).Exec()
	t.t.Logf("Got error: %+v", err)
	if err != nil {
		return err
	}
	return nil
}
