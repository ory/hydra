// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlcon"
)

type row struct {
	key   sql.NullString
	value string
}

// NewChangeFeedConnection opens a new connection to the database and enables the CHANGEFEED feature.
//
// The caller is responsible for closing the connection when done.
//
// You must register the `pgx` driver before calling this function:
//
//	import _ "github.com/jackc/pgx/v4/stdlib"
func NewChangeFeedConnection(ctx context.Context, l *logrusx.Logger, dsn string) (*sqlx.DB, error) {
	if !strings.HasPrefix(dsn, "cockroach://") {
		return nil, errors.Errorf("DSN value must be prefixed with a cockroach URI schema")
	}

	_, _, _, _, cleanedDSN := sqlcon.ParseConnectionOptions(l, dsn)
	cleanedDSN = strings.Replace(cleanedDSN, "cockroach://", "postgres://", 1)
	l.WithField("component", "github.com/ory/x/watcherx.NewChangeFeedConnection").Info("Opening watcherx database connection.")
	cx, err := sqlx.Open("pgx", cleanedDSN)
	if err != nil {
		return nil, err
	}

	l.WithField("component", "github.com/ory/x/watcherx.NewChangeFeedConnection").Info("Connection to watcherx database is open.")

	cx.SetMaxIdleConns(1)
	cx.SetMaxOpenConns(1)
	cx.SetConnMaxLifetime(-1)
	cx.SetConnMaxIdleTime(-1)

	l.WithField("component", "github.com/ory/x/watcherx.NewChangeFeedConnection").Info("Trying to ping the watcherx database connection.")

	if err := cx.PingContext(ctx); err != nil {
		return nil, err
	}

	l.WithField("component", "github.com/ory/x/watcherx.NewChangeFeedConnection").Info("Enabling CHANGEFEED on watcherx database connection.")

	// Ensure CHANGEFEED is enabled
	_, err = cx.ExecContext(ctx, "SET CLUSTER SETTING kv.rangefeed.enabled = true")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	l.WithField("component", "github.com/ory/x/watcherx.NewChangeFeedConnection").Info("Initialization of CHANGEFEED is done.")

	return cx, nil
}

const heartBeatInterval = time.Second

// WatchChangeFeed sends changed rows on the channel. To cancel the execution, cancel the context!
//
// Watcher.DispatchNow() does not have an effect in this method.
//
// This function spawns the necessary go-routines to process the change-feed events and deduplicate them.
func WatchChangeFeed(ctx context.Context, cx *sqlx.DB, tableName string, out EventChannel, cursor time.Time) (_ Watcher, err error) {
	c := make(EventChannel)
	InternalDeduplicate(ctx, c, out, 100)

	var rows *sql.Rows
	if cursor.IsZero() {
		rows, err = cx.QueryContext(ctx, fmt.Sprintf("EXPERIMENTAL CHANGEFEED FOR %s RESOLVED = $1, MIN_CHECKPOINT_FREQUENCY = $2", tableName), heartBeatInterval.String(), heartBeatInterval.String())
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		var err error
		rows, err = cx.QueryContext(ctx, fmt.Sprintf("EXPERIMENTAL CHANGEFEED FOR %s WITH CURSOR = $1, RESOLVED = $2, MIN_CHECKPOINT_FREQUENCY = $3", tableName), strconv.Itoa(int(cursor.UnixNano())), heartBeatInterval.String(), heartBeatInterval.String())
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	d := newDispatcher()

	// basically run the watcher in a go routine which gets canceled either by the connection being closed
	// or by calling `"CANCEL QUERY"` below.
	heartBeat := make(chan struct{})

	// The "control" go routine to detect if the changefeed is still alive.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-heartBeat:
				if !ok {
					return
				}
			case <-time.After(heartBeatInterval * 10):
				c <- &ErrorEvent{
					error: errors.New("unable to detect changefeed heartbeat in time"),
				}
			case <-d.trigger:
				d.done <- 0
			}
		}
	}()

	// The "work" go routine to read the changefeed and send events.
	go func() {
		defer func() {
			// we signal that we are done
			close(heartBeat)
		}()

		var r row
		var table sql.NullString

		for rows.Next() {
			if err := errors.WithStack(rows.Scan(&table, &r.key, &r.value)); err != nil {
				c <- &ErrorEvent{
					error: err,
				}
				continue
			}

			keys := gjson.Parse(r.key.String)
			eventSource := keys.Raw

			// For some reason this is an array - maybe because of composite primary keys?
			// See: https://www.cockroachlabs.com/docs/v20.2/changefeed-for.html
			if ka := keys.Array(); len(ka) > 0 {
				ids := make([]string, len(ka))
				for i := range ka {
					ids[i] = ka[i].String()
				}

				eventSource = strings.Join(ids, "/")
			}

			if gjson.Get(r.value, "resolved").Exists() {
				heartBeat <- struct{}{}
				continue
			}

			after := gjson.Get(r.value, "after")
			if after.IsObject() {
				c <- &ChangeEvent{
					data:   []byte(after.Raw),
					source: source(eventSource),
				}
			} else {
				c <- &RemoveEvent{
					source: source(eventSource),
				}
			}
		}

		if err := rows.Err(); err != nil {
			// We can land here (after the row read loop) when the context is closed, which means there is probably
			// no receiver anymore. Let's just try to send the error in case someone is listening.
			select {
			case c <- &ErrorEvent{
				error: err,
			}:
			case <-ctx.Done():
			}
			return
		}

		if err := rows.Close(); err != nil {
			// We can land here (after the row read loop) when the context is closed, which means there is probably
			// no receiver anymore. Let's just try to send the error in case someone is listening.
			select {
			case c <- &ErrorEvent{
				error: err,
			}:
			case <-ctx.Done():
			}
			return
		}

		// no need to close rows here, as they are closed if rows.Next() returns false

		if err := cx.Close(); err != nil {
			// We can land here (after the row read loop) when the context is closed, which means there is probably
			// no receiver anymore. Let's just try to send the error in case someone is listening.
			select {
			case c <- &ErrorEvent{
				error: err,
			}:
			case <-ctx.Done():
			}
			return
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return d, nil
}

// InternalDeduplicate sents events from `events` to the `deduplicated` channel, but
// deduplicates events that are sent multiple times. This is necessary, because
// the CochroachDB changefeed has a atleast-once guarantee for change events,
// meaning that events could be sent multiple times.
//
// For deduplication, the last x `pastEvents` are considered.
func InternalDeduplicate(ctx context.Context, in <-chan Event, out chan<- Event, pastEvents int) {
	go func() {
		previous := newRingBuffer(pastEvents)

		for {
			select {
			case e, ok := <-in:
				if !ok {
					// we only want to close the channel if the input channel is closed
					close(out)
					return
				}
				if previous.Contains(e) {
					// Ignore event
					continue
				} else {
					previous.Add(e)
					out <- e
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

type ringBufferKey [sha256.Size]byte

var emptyKey ringBufferKey

// ringBuffer is a data structure for constant-time set membership (through
// `Contains`) while maintaining constant memory usage by keeping at most
// `capacity` elements.
//
// ringBuffer is not safe for concurrent use.
type ringBuffer struct {
	capacity int
	seen     map[ringBufferKey]struct{} // map for efficient Contains().
	keys     []ringBufferKey            // ring buffer so we can evict events on FIFO basis.
	keyIdx   int                        // index of the next key to be added.
}

func newRingBuffer(capacity int) *ringBuffer {
	return &ringBuffer{
		capacity: capacity,
		seen:     make(map[ringBufferKey]struct{}, capacity),
		keys:     make([]ringBufferKey, capacity),
	}
}

func (r *ringBuffer) key(el fmt.Stringer) ringBufferKey {
	return sha256.Sum256([]byte(el.String()))
}

func (r *ringBuffer) Contains(el fmt.Stringer) bool {
	_, ok := r.seen[r.key(el)]
	return ok
}

func (r *ringBuffer) Add(el fmt.Stringer) {
	// Evict the oldest key.
	if oldestKey := r.keys[r.keyIdx%r.capacity]; oldestKey != emptyKey {
		delete(r.seen, oldestKey)
	}

	key := r.key(el)
	r.seen[key] = struct{}{}
	r.keys[r.keyIdx%r.capacity] = key

	r.keyIdx++
}
