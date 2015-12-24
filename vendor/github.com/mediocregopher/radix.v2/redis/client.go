package redis

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"
)

const (
	bufSize int = 4096
)

// ErrPipelineEmpty is returned from PipeResp() to indicate that all commands
// which were put into the pipeline have had their responses read
var ErrPipelineEmpty = errors.New("pipeline queue empty")

// Client describes a Redis client.
type Client struct {
	conn         net.Conn
	respReader   *RespReader
	timeout      time.Duration
	pending      []request
	writeScratch []byte
	writeBuf     *bytes.Buffer

	completed, completedHead []*Resp

	// The network/address of the redis instance this client is connected to.
	// These will be wahtever strings were passed into the Dial function when
	// creating this connection
	Network, Addr string

	// The most recent critical network error which occured when either reading
	// or writing. A critical network error is one in which the connection was
	// found to be no longer usable; in essence, any error except a timeout.
	// Close is automatically called on the client when it encounters a critical
	// network error
	LastCritical error
}

// request describes a client's request to the redis server
type request struct {
	cmd  string
	args []interface{}
}

// DialTimeout connects to the given Redis server with the given timeout, which
// will be used as the read/write timeout when communicating with redis
func DialTimeout(network, addr string, timeout time.Duration) (*Client, error) {
	// establish a connection
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}

	completed := make([]*Resp, 0, 10)
	return &Client{
		conn:          conn,
		respReader:    NewRespReader(conn),
		timeout:       timeout,
		writeScratch:  make([]byte, 0, 128),
		writeBuf:      bytes.NewBuffer(make([]byte, 0, 128)),
		completed:     completed,
		completedHead: completed,
		Network:       network,
		Addr:          addr,
	}, nil
}

// Dial connects to the given Redis server.
func Dial(network, addr string) (*Client, error) {
	return DialTimeout(network, addr, time.Duration(0))
}

// Close closes the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Cmd calls the given Redis command.
func (c *Client) Cmd(cmd string, args ...interface{}) *Resp {
	err := c.writeRequest(request{cmd, args})
	if err != nil {
		return newRespIOErr(err)
	}
	return c.ReadResp()
}

// PipeAppend adds the given call to the pipeline queue.
// Use PipeResp() to read the response.
func (c *Client) PipeAppend(cmd string, args ...interface{}) {
	c.pending = append(c.pending, request{cmd, args})
}

// PipeResp returns the reply for the next request in the pipeline queue. Err
// with ErrPipelineEmpty is returned if the pipeline queue is empty.
func (c *Client) PipeResp() *Resp {
	if len(c.completed) > 0 {
		r := c.completed[0]
		c.completed = c.completed[1:]
		return r
	}

	if len(c.pending) == 0 {
		return NewResp(ErrPipelineEmpty)
	}

	nreqs := len(c.pending)
	err := c.writeRequest(c.pending...)
	c.pending = nil
	if err != nil {
		return newRespIOErr(err)
	}
	c.completed = c.completedHead
	for i := 0; i < nreqs; i++ {
		r := c.ReadResp()
		c.completed = append(c.completed, r)
	}

	// At this point c.completed should have something in it
	return c.PipeResp()
}

// PipeClient clears the contents of the current pipeline queue, both commands
// queued by PipeAppend which have yet to be sent and responses which have yet
// to be retrieved through PipeResp. The first returned int will be the number
// of pending commands dropped, the second will be the number of pending
// responses dropped
func (c *Client) PipeClear() (int, int) {
	callCount, replyCount := len(c.pending), len(c.completed)
	if callCount > 0 {
		c.pending = nil
	}
	if replyCount > 0 {
		c.completed = nil
	}
	return callCount, replyCount
}

// ReadResp will read a Resp off of the connection without sending anything
// first (useful after you've sent a SUSBSCRIBE command). This will block until
// a reply is received or the timeout is reached (returning the IOErr). You can
// use IsTimeout to check if the Resp is due to a Timeout
//
// Note: this is a more low-level function, you really shouldn't have to
// actually use it unless you're writing your own pub/sub code
func (c *Client) ReadResp() *Resp {
	if c.timeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
	}
	r := c.respReader.Read()
	if r.IsType(IOErr) && !IsTimeout(r) {
		c.LastCritical = r.Err
		c.Close()
	}
	return r
}

func (c *Client) writeRequest(requests ...request) error {
	if c.timeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
	}
	var err error
outer:
	for i := range requests {
		c.writeBuf.Reset()
		elems := flattenedLength(requests[i].args...) + 1
		_, err = writeArrayHeader(c.writeBuf, c.writeScratch, int64(elems))
		if err != nil {
			break
		}

		_, err = writeTo(c.writeBuf, c.writeScratch, requests[i].cmd, true, true)
		if err != nil {
			break
		}

		for _, arg := range requests[i].args {
			_, err = writeTo(c.writeBuf, c.writeScratch, arg, true, true)
			if err != nil {
				break outer
			}
		}

		if _, err = c.writeBuf.WriteTo(c.conn); err != nil {
			break
		}
	}
	if err != nil {
		c.LastCritical = err
		c.Close()
		return err
	}
	return nil
}

var errBadCmdNoKey = errors.New("bad command, no key")

// KeyFromArgs is a helper function which other library packages which wrap this
// one might find useful. It takes in a set of arguments which might be passed
// into Cmd and returns the first key for the command. Since radix supports
// complicated arguments (like slices, slices of slices, maps, etc...) this is
// not always as straightforward as it might seem, so this helper function is
// provided.
//
// An error is returned if no key can be determined
func KeyFromArgs(args ...interface{}) (string, error) {
	if len(args) == 0 {
		return "", errBadCmdNoKey
	}
	arg := args[0]
	switch argv := arg.(type) {
	case string:
		return argv, nil
	case []byte:
		return string(argv), nil
	default:
		switch reflect.TypeOf(arg).Kind() {
		case reflect.Slice:
			argVal := reflect.ValueOf(arg)
			if argVal.Len() < 1 {
				return "", errBadCmdNoKey
			}
			first := argVal.Index(0).Interface()
			return KeyFromArgs(first)
		case reflect.Map:
			// Maps have no order, we can't possibly choose a key out of one
			return "", errBadCmdNoKey
		default:
			return fmt.Sprint(arg), nil
		}
	}
}
