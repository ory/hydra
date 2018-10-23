package instrumentedsql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// The possible op values passed to the logger and used for child span names
const (
	OpSQLPrepare         = "sql-prepare"
	OpSQLConnExec        = "sql-conn-exec"
	OpSQLConnQuery       = "sql-conn-query"
	OpSQLStmtExec        = "sql-stmt-exec"
	OpSQLStmtQuery       = "sql-stmt-query"
	OpSQLStmtClose       = "sql-stmt-close"
	OpSQLTxBegin         = "sql-tx-begin"
	OpSQLTxCommit        = "sql-tx-commit"
	OpSQLTxRollback      = "sql-tx-rollback"
	OpSQLResLastInsertID = "sql-res-lastInsertId"
	OpSQLResRowsAffected = "sql-res-rowsAffected"
	OpSQLRowsNext        = "sql-rows-next"
	OpSQLPing            = "sql-ping"
	OpSQLDummyPing       = "sql-dummy-ping"
)

type wrappedDriver struct {
	opts
	parent driver.Driver
}

type wrappedConn struct {
	opts
	parent driver.Conn
}

type wrappedTx struct {
	opts
	ctx    context.Context
	parent driver.Tx
}

type wrappedStmt struct {
	opts
	ctx    context.Context
	query  string
	parent driver.Stmt
}

type wrappedResult struct {
	opts
	ctx    context.Context
	parent driver.Result
}

type wrappedRows struct {
	opts
	ctx    context.Context
	parent driver.Rows
}

// WrapDriver will wrap the passed SQL driver and return a new sql driver that uses it and also logs and traces calls using the passed logger and tracer
// The returned driver will still have to be registered with the sql package before it can be used.
//
// Important note: Seeing as the context passed into the various instrumentation calls this package calls,
// Any call without a context passed will not be instrumented. Please be sure to use the ___Context() and BeginTx() function calls added in Go 1.8
// instead of the older calls which do not accept a context.
func WrapDriver(driver driver.Driver, opts ...Opt) driver.Driver {
	d := wrappedDriver{parent: driver}

	for _, opt := range opts {
		opt(&d.opts)
	}

	if d.Logger == nil {
		d.Logger = nullLogger{}
	}
	if d.Tracer == nil {
		d.Tracer = nullTracer{}
	}

	return d
}

func formatArgs(args interface{}) string {
	argsVal := reflect.ValueOf(args)
	if argsVal.Kind() != reflect.Slice {
		return "<unknown>"
	}

	strArgs := make([]string, 0, argsVal.Len())
	for i := 0; i < argsVal.Len(); i++ {
		strArgs = append(strArgs, formatArg(argsVal.Index(i).Interface()))
	}

	return fmt.Sprintf("{%s}", strings.Join(strArgs, ", "))
}

func formatArg(arg interface{}) string {
	strArg := ""
	switch arg := arg.(type) {
	case []uint8:
		strArg = fmt.Sprintf("[%T len:%d]", arg, len(arg))
	case string:
		strArg = fmt.Sprintf("[%T %q]", arg, arg)
	case driver.NamedValue:
		if arg.Name != "" {
			strArg = fmt.Sprintf("[%T %s=%v]", arg.Value, arg.Name, formatArg(arg.Value))
		} else {
			strArg = formatArg(arg.Value)
		}
	default:
		strArg = fmt.Sprintf("[%T %v]", arg, arg)
	}

	return strArg
}

func logQuery(ctx context.Context, opts opts, op, query string, err error, args interface{}, since time.Time) {
	keyvals := []interface{}{
		"query", query,
		"err", err,
		"duration", time.Since(since),
	}

	if !opts.OmitArgs && args != nil {
		keyvals = append(keyvals, "args", formatArgs(args))
	}

	opts.Log(ctx, op, keyvals...)
}

func (d wrappedDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.parent.Open(name)
	if err != nil {
		return nil, err
	}

	return wrappedConn{opts: d.opts, parent: conn}, nil
}

func (c wrappedConn) Prepare(query string) (driver.Stmt, error) {
	parent, err := c.parent.Prepare(query)
	if err != nil {
		return nil, err
	}

	return wrappedStmt{opts: c.opts, query: query, parent: parent}, nil
}

func (c wrappedConn) Close() error {
	return c.parent.Close()
}

func (c wrappedConn) Begin() (driver.Tx, error) {
	tx, err := c.parent.Begin()
	if err != nil {
		return nil, err
	}

	return wrappedTx{opts: c.opts, parent: tx}, nil
}

func (c wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	span := c.GetSpan(ctx).NewChild(OpSQLTxBegin)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		c.Log(ctx, OpSQLTxBegin, "err", err, "duration", time.Since(start))
	}()

	if connBeginTx, ok := c.parent.(driver.ConnBeginTx); ok {
		tx, err = connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		return wrappedTx{opts: c.opts, ctx: ctx, parent: tx}, nil
	}

	tx, err = c.parent.Begin()
	if err != nil {
		return nil, err
	}

	return wrappedTx{opts: c.opts, ctx: ctx, parent: tx}, nil
}

func (c wrappedConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	span := c.GetSpan(ctx).NewChild(OpSQLPrepare)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		logQuery(ctx, c.opts, OpSQLPrepare, query, err, nil, start)
	}()

	if connPrepareCtx, ok := c.parent.(driver.ConnPrepareContext); ok {
		stmt, err := connPrepareCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}

		return wrappedStmt{opts: c.opts, ctx: ctx, parent: stmt}, nil
	}

	return c.Prepare(query)
}

func (c wrappedConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if execer, ok := c.parent.(driver.Execer); ok {
		res, err := execer.Exec(query, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{opts: c.opts, parent: res}, nil
	}

	return nil, driver.ErrSkip
}

func (c wrappedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	span := c.GetSpan(ctx).NewChild(OpSQLConnExec)
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", query)
	if !c.OmitArgs {
		span.SetLabel("args", formatArgs(args))
	}
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()

		logQuery(ctx, c.opts, OpSQLConnExec, query, err, args, start)
	}()

	if execContext, ok := c.parent.(driver.ExecerContext); ok {
		res, err := execContext.ExecContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{opts: c.opts, ctx: ctx, parent: res}, nil
	}

	// Fallback implementation
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Exec(query, dargs)
}

func (c wrappedConn) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.parent.(driver.Pinger); ok {
		span := c.GetSpan(ctx).NewChild(OpSQLPing)
		span.SetLabel("component", "database/sql")
		start := time.Now()
		defer func() {
			span.SetError(err)
			span.Finish()
			c.Log(ctx, OpSQLPing, "err", err, "duration", time.Since(start))
		}()

		return pinger.Ping(ctx)
	}

	c.Log(ctx, OpSQLDummyPing, "duration", time.Duration(0))

	return nil
}

func (c wrappedConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if queryer, ok := c.parent.(driver.Queryer); ok {
		rows, err := queryer.Query(query, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{opts: c.opts, parent: rows}, nil
	}

	return nil, driver.ErrSkip
}

func (c wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span := c.GetSpan(ctx).NewChild(OpSQLConnQuery)
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", query)
	if !c.OmitArgs {
		span.SetLabel("args", formatArgs(args))
	}
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		logQuery(ctx, c.opts, OpSQLConnQuery, query, err, args, start)
	}()

	if queryerContext, ok := c.parent.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{opts: c.opts, ctx: ctx, parent: rows}, nil
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Query(query, dargs)
}

func (t wrappedTx) Commit() (err error) {
	span := t.GetSpan(t.ctx).NewChild(OpSQLTxCommit)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		t.Log(t.ctx, OpSQLTxCommit, "err", err, "duration", time.Since(start))
	}()

	return t.parent.Commit()
}

func (t wrappedTx) Rollback() (err error) {
	span := t.GetSpan(t.ctx).NewChild(OpSQLTxRollback)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		t.Log(t.ctx, OpSQLTxRollback, "err", err, "duration", time.Since(start))
	}()

	return t.parent.Rollback()
}

func (s wrappedStmt) Close() (err error) {
	span := s.GetSpan(s.ctx).NewChild(OpSQLStmtClose)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		s.Log(s.ctx, OpSQLStmtClose, "err", err, "duration", time.Since(start))
	}()

	return s.parent.Close()
}

func (s wrappedStmt) NumInput() int {
	return s.parent.NumInput()
}

func (s wrappedStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	span := s.GetSpan(s.ctx).NewChild(OpSQLStmtExec)
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", formatArgs(args))
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		logQuery(s.ctx, s.opts, OpSQLStmtExec, s.query, err, args, start)
	}()

	res, err = s.parent.Exec(args)
	if err != nil {
		return nil, err
	}

	return wrappedResult{opts: s.opts, ctx: s.ctx, parent: res}, nil
}

func (s wrappedStmt) Query(args []driver.Value) (rows driver.Rows, err error) {
	span := s.GetSpan(s.ctx).NewChild(OpSQLStmtQuery)
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", formatArgs(args))
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		logQuery(s.ctx, s.opts, OpSQLStmtQuery, s.query, err, args, start)
	}()

	rows, err = s.parent.Query(args)
	if err != nil {
		return nil, err
	}

	return wrappedRows{opts: s.opts, ctx: s.ctx, parent: rows}, nil
}

func (s wrappedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	span := s.GetSpan(ctx).NewChild(OpSQLStmtExec)
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", formatArgs(args))
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		logQuery(ctx, s.opts, OpSQLStmtExec, s.query, err, args, start)
	}()

	if stmtExecContext, ok := s.parent.(driver.StmtExecContext); ok {
		res, err := stmtExecContext.ExecContext(ctx, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{opts: s.opts, ctx: ctx, parent: res}, nil
	}

	// Fallback implementation
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Exec(dargs)
}

func (s wrappedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	span := s.GetSpan(ctx).NewChild(OpSQLStmtQuery)
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", formatArgs(args))
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		logQuery(ctx, s.opts, OpSQLStmtQuery, s.query, err, args, start)
	}()

	if stmtQueryContext, ok := s.parent.(driver.StmtQueryContext); ok {
		rows, err := stmtQueryContext.QueryContext(ctx, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{opts: s.opts, ctx: ctx, parent: rows}, nil
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Query(dargs)
}

func (r wrappedResult) LastInsertId() (id int64, err error) {
	span := r.GetSpan(r.ctx).NewChild(OpSQLResLastInsertID)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		r.Log(r.ctx, OpSQLResLastInsertID, "err", err, "duration", time.Since(start))
	}()

	return r.parent.LastInsertId()
}

func (r wrappedResult) RowsAffected() (num int64, err error) {
	span := r.GetSpan(r.ctx).NewChild(OpSQLResRowsAffected)
	span.SetLabel("component", "database/sql")
	start := time.Now()
	defer func() {
		span.SetError(err)
		span.Finish()
		r.Log(r.ctx, OpSQLResRowsAffected, "err", err, "duration", time.Since(start))
	}()

	return r.parent.RowsAffected()
}

func (r wrappedRows) Columns() []string {
	return r.parent.Columns()
}

func (r wrappedRows) Close() error {
	return r.parent.Close()
}

func (r wrappedRows) Next(dest []driver.Value) (err error) {
	if r.opts.TraceRowsNext {
		span := r.GetSpan(r.ctx).NewChild(OpSQLRowsNext)
		span.SetLabel("component", "database/sql")
		defer func() {
			span.SetError(err)
			span.Finish()
		}()
	}

	start := time.Now()
	defer func() {
		r.Log(r.ctx, OpSQLRowsNext, "err", err, "duration", time.Since(start))
	}()

	return r.parent.Next(dest)
}

// namedValueToValue is a helper function copied from the database/sql package
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
