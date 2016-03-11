package gorethink

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"gopkg.in/fatih/pool.v2"
)

const maxBadConnRetries = 10

var (
	errPoolClosed = errors.New("gorethink: pool is closed")
	errConnClosed = errors.New("gorethink: conn is closed")
	errConnBusy   = errors.New("gorethink: conn is busy")
)

// A Pool is used to store a pool of connections to a single RethinkDB server
type Pool struct {
	host Host
	opts *ConnectOpts

	pool pool.Pool

	mu      sync.RWMutex // protects following fields
	closed  bool
	maxIdle int
	maxOpen int
}

// NewPool creates a new connection pool for the given host
func NewPool(host Host, opts *ConnectOpts) (*Pool, error) {
	maxOpen := opts.MaxOpen
	if maxOpen == 0 {
		maxOpen = 2
	}

	p, err := pool.NewChannelPool(opts.MaxIdle, maxOpen, func() (net.Conn, error) {
		conn, err := NewConnection(host.String(), opts)
		if err != nil {
			return nil, err
		}

		return conn, err
	})
	if err != nil {
		return nil, err
	}

	return &Pool{
		pool: p,
		host: host,
		opts: opts,
	}, nil
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (p *Pool) Ping() error {
	_, pc, err := p.conn()
	if err != nil {
		return err
	}
	return pc.Close()
}

// Close closes the database, releasing any open resources.
//
// It is rare to Close a Pool, as the Pool handle is meant to be
// long-lived and shared between many goroutines.
func (p *Pool) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return nil
	}

	p.pool.Close()

	return nil
}

func (p *Pool) conn() (*Connection, *pool.PoolConn, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return nil, nil, errPoolClosed
	}

	nc, err := p.pool.Get()
	if err != nil {
		return nil, nil, err
	}

	pc, ok := nc.(*pool.PoolConn)
	if !ok {
		// This should never happen!
		return nil, nil, fmt.Errorf("Invalid connection in pool")
	}

	conn, ok := pc.Conn.(*Connection)
	if !ok {
		// This should never happen!
		return nil, nil, fmt.Errorf("Invalid connection in pool")
	}

	return conn, pc, nil
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle connections are retained.
//
// Deprecated: This value should only be set when connecting
func (p *Pool) SetMaxIdleConns(n int) {
	return
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
//
// Deprecated: This value should only be set when connecting
func (p *Pool) SetMaxOpenConns(n int) {
	return
}

// Query execution functions

// Exec executes a query without waiting for any response.
func (p *Pool) Exec(q Query) error {
	var err error

	for i := 0; i < maxBadConnRetries; i++ {
		var c *Connection
		var pc *pool.PoolConn

		c, pc, err = p.conn()
		if err != nil {
			break
		}
		defer pc.Close()

		_, _, err = c.Query(q)

		if c.isBad() {
			pc.MarkUnusable()
		}

		break
	}

	return err
}

// Query executes a query and waits for the response
func (p *Pool) Query(q Query) (*Cursor, error) {
	var cursor *Cursor
	var err error

	for i := 0; i < maxBadConnRetries; i++ {
		var c *Connection
		var pc *pool.PoolConn

		c, pc, err = p.conn()
		if err != nil {
			break
		}

		_, cursor, err = c.Query(q)

		if err == nil {
			cursor.releaseConn = releaseConn(c, pc)
		}

		break
	}

	return cursor, err
}

// Server returns the server name and server UUID being used by a connection.
func (p *Pool) Server() (ServerResponse, error) {
	var response ServerResponse
	var err error

	for i := 0; i < maxBadConnRetries; i++ {
		var c *Connection
		var pc *pool.PoolConn

		c, pc, err = p.conn()
		if err != nil {
			break
		}

		response, err = c.Server()

		if c.isBad() {
			pc.MarkUnusable()
		}

		break
	}

	return response, err
}

func releaseConn(c *Connection, pc *pool.PoolConn) func() error {
	return func() error {
		if c.isBad() {
			pc.MarkUnusable()
		}

		return pc.Close()
	}
}
