package gorethink

import (
	"crypto/tls"
	"sync"
	"time"

	p "github.com/dancannon/gorethink/ql2"
)

// A Session represents a connection to a RethinkDB cluster and should be used
// when executing queries.
type Session struct {
	hosts []Host
	opts  *ConnectOpts

	mu      sync.RWMutex
	cluster *Cluster
	closed  bool
}

// ConnectOpts is used to specify optional arguments when connecting to a cluster.
type ConnectOpts struct {
	Address      string        `gorethink:"address,omitempty"`
	Addresses    []string      `gorethink:"addresses,omitempty"`
	Database     string        `gorethink:"database,omitempty"`
	AuthKey      string        `gorethink:"authkey,omitempty"`
	Timeout      time.Duration `gorethink:"timeout,omitempty"`
	WriteTimeout time.Duration `gorethink:"write_timeout,omitempty"`
	ReadTimeout  time.Duration `gorethink:"read_timeout,omitempty"`
	TLSConfig    *tls.Config   `gorethink:"tlsconfig,omitempty"`

	MaxIdle int `gorethink:"max_idle,omitempty"`
	// By default a maximum of 2 connections are opened per host.
	MaxOpen int `gorethink:"max_open,omitempty"`

	// Below options are for cluster discovery, please note there is a high
	// probability of these changing as the API is still being worked on.

	// DiscoverHosts is used to enable host discovery, when true the driver
	// will attempt to discover any new nodes added to the cluster and then
	// start sending queries to these new nodes.
	DiscoverHosts bool `gorethink:"discover_hosts,omitempty"`
	// NodeRefreshInterval is used to determine how often the driver should
	// refresh the status of a node.
	//
	// Deprecated: This function is no longer used due to changes in the
	// way hosts are selected.
	NodeRefreshInterval time.Duration `gorethink:"node_refresh_interval,omitempty"`
	// HostDecayDuration is used by the go-hostpool package to calculate a weighted
	// score when selecting a host. By default a value of 5 minutes is used.
	HostDecayDuration time.Duration

	// Indicates whether the cursors running in this session should use json.Number instead of float64 while
	// unmarshaling documents with interface{}. The default is `false`.
	UseJSONNumber bool
}

func (o *ConnectOpts) toMap() map[string]interface{} {
	return optArgsToMap(o)
}

// Connect creates a new database session. To view the available connection
// options see ConnectOpts.
//
// By default maxIdle and maxOpen are set to 1: passing values greater
// than the default (e.g. MaxIdle: "10", MaxOpen: "20") will provide a
// pool of re-usable connections.
//
// Basic connection example:
//
// 	session, err := r.Connect(r.ConnectOpts{
// 		Host: "localhost:28015",
// 		Database: "test",
// 		AuthKey:  "14daak1cad13dj",
// 	})
//
// Cluster connection example:
//
// 	session, err := r.Connect(r.ConnectOpts{
// 		Hosts: []string{"localhost:28015", "localhost:28016"},
// 		Database: "test",
// 		AuthKey:  "14daak1cad13dj",
// 	})
func Connect(opts ConnectOpts) (*Session, error) {
	var addresses = opts.Addresses
	if len(addresses) == 0 {
		addresses = []string{opts.Address}
	}

	hosts := make([]Host, len(addresses))
	for i, address := range addresses {
		hostname, port := splitAddress(address)
		hosts[i] = NewHost(hostname, port)
	}
	if len(hosts) <= 0 {
		return nil, ErrNoHosts
	}

	// Connect
	s := &Session{
		hosts: hosts,
		opts:  &opts,
	}

	err := s.Reconnect()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// CloseOpts allows calls to the Close function to be configured.
type CloseOpts struct {
	NoReplyWait bool `gorethink:"noreplyWait,omitempty"`
}

func (o *CloseOpts) toMap() map[string]interface{} {
	return optArgsToMap(o)
}

// IsConnected returns true if session has a valid connection.
func (s *Session) IsConnected() bool {
	if s.cluster == nil || s.closed {
		return false
	}
	return s.cluster.IsConnected()
}

// Reconnect closes and re-opens a session.
func (s *Session) Reconnect(optArgs ...CloseOpts) error {
	var err error

	if err = s.Close(optArgs...); err != nil {
		return err
	}

	s.mu.Lock()
	s.cluster, err = NewCluster(s.hosts, s.opts)
	if err != nil {
		return err
	}

	s.closed = false
	s.mu.Unlock()

	return nil
}

// Close closes the session
func (s *Session) Close(optArgs ...CloseOpts) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	if len(optArgs) >= 1 {
		if optArgs[0].NoReplyWait {
			s.mu.Unlock()
			s.NoReplyWait()
			s.mu.Lock()
		}
	}

	if s.cluster != nil {
		s.cluster.Close()
	}
	s.cluster = nil
	s.closed = true

	return nil
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func (s *Session) SetMaxIdleConns(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.opts.MaxIdle = n
	s.cluster.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (s *Session) SetMaxOpenConns(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.opts.MaxOpen = n
	s.cluster.SetMaxOpenConns(n)
}

// NoReplyWait ensures that previous queries with the noreply flag have been
// processed by the server. Note that this guarantee only applies to queries
// run on the given connection
func (s *Session) NoReplyWait() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return ErrConnectionClosed
	}

	return s.cluster.Exec(Query{
		Type: p.Query_NOREPLY_WAIT,
	})
}

// Use changes the default database used
func (s *Session) Use(database string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.opts.Database = database
}

// Query executes a ReQL query using the session to connect to the database
func (s *Session) Query(q Query) (*Cursor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrConnectionClosed
	}

	return s.cluster.Query(q)
}

// Exec executes a ReQL query using the session to connect to the database
func (s *Session) Exec(q Query) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return ErrConnectionClosed
	}

	return s.cluster.Exec(q)
}

// Server returns the server name and server UUID being used by a connection.
func (s *Session) Server() (ServerResponse, error) {
	return s.cluster.Server()
}

// SetHosts resets the hosts used when connecting to the RethinkDB cluster
func (s *Session) SetHosts(hosts []Host) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hosts = hosts
}

func (s *Session) newQuery(t Term, opts map[string]interface{}) (Query, error) {
	return newQuery(t, opts, s.opts)
}
