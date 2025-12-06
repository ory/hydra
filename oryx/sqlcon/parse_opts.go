// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlcon

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ory/x/logrusx"
)

// ParseConnectionOptions parses values for max_conns, max_idle_conns, max_conn_lifetime from DSNs.
// It also returns the URI without those query parameters.
func ParseConnectionOptions(l *logrusx.Logger, dsn string) (maxConns int, maxIdleConns int, maxConnLifetime, maxIdleConnTime time.Duration, cleanedDSN string) {
	maxConns = maxParallelism() * 2
	maxIdleConns = maxParallelism()
	cleanedDSN = dsn

	hostPath, rawQuery, ok := strings.Cut(dsn, "?")
	if !ok {
		l.
			WithField("sql_max_connections", maxConns).
			WithField("sql_max_idle_connections", maxIdleConns).
			WithField("sql_max_connection_lifetime", maxConnLifetime).
			WithField("sql_max_idle_connection_time", maxIdleConnTime).
			Debugf("No SQL connection options have been defined, falling back to default connection options.")
		return
	}

	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		l.
			WithField("sql_max_connections", maxConns).
			WithField("sql_max_idle_connections", maxIdleConns).
			WithField("sql_max_connection_lifetime", maxConnLifetime).
			WithField("sql_max_idle_connection_time", maxIdleConnTime).
			WithError(err).
			Warnf("Unable to parse SQL DSN query, falling back to default connection options.")
		return
	}

	if v := query.Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.WithError(err).Warnf(`SQL DSN query parameter "max_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxConns)
		} else {
			maxConns = int(s)
		}
		query.Del("max_conns")
	}

	if v := query.Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.WithError(err).Warnf(`SQL DSN query parameter "max_idle_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxIdleConns)
		} else {
			maxIdleConns = int(s)
		}
		query.Del("max_idle_conns")
	}

	if v := query.Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			l.WithError(err).Warnf(`SQL DSN query parameter "max_conn_lifetime" value %v could not be parsed to duration, falling back to default value %d`, v, maxConnLifetime)
		} else {
			maxConnLifetime = s
		}
		query.Del("max_conn_lifetime")
	}

	if v := query.Get("max_conn_idle_time"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			l.WithError(err).Warnf(`SQL DSN query parameter "max_conn_idle_time" value %v could not be parsed to duration, falling back to default value %d`, v, maxIdleConnTime)
		} else {
			maxIdleConnTime = s
		}
		query.Del("max_conn_idle_time")
	}
	cleanedDSN = fmt.Sprintf("%s?%s", hostPath, query.Encode())

	return
}

// FinalizeDSN will return a finalized DSN URI.
func FinalizeDSN(l *logrusx.Logger, dsn string) string {
	if !strings.HasPrefix(dsn, "mysql://") {
		return dsn
	}

	var q url.Values
	hostPath, query, ok := strings.Cut(dsn, "?")

	if !ok {
		q = make(url.Values)
	} else {
		var err error
		q, err = url.ParseQuery(query)
		if err != nil {
			l.WithError(err).Warnf("Unable to parse SQL DSN query, could not finalize the DSN URI.")
			return dsn
		}
	}

	q.Set("multiStatements", "true")
	q.Set("parseTime", "true")

	// This causes an UPDATE to return the number of matching rows instead of
	// the number of rows changed. This ensures compatibility with PostgreSQL and SQLite behavior.
	q.Set("clientFoundRows", "true")

	return fmt.Sprintf("%s?%s", hostPath, q.Encode())
}
