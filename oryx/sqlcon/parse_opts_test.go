// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlcon

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/logrusx"
)

func TestParseConnectionOptions(t *testing.T) {
	defaultMaxConns, defaultMaxIdleConns, defaultMaxConnIdleTime, defaultMaxConnLifetime := maxParallelism()*2, maxParallelism(), time.Duration(0), time.Duration(0)
	logger := logrusx.New("", "")
	for i, tc := range []struct {
		name, dsn, cleanedDSN            string
		maxConns, maxIdleConns           int
		maxConnIdleTime, maxConnLifetime time.Duration
	}{
		{
			name:            "no parameters",
			dsn:             "postgres://user:pwd@host:port",
			cleanedDSN:      "postgres://user:pwd@host:port",
			maxConns:        defaultMaxConns,
			maxIdleConns:    defaultMaxIdleConns,
			maxConnIdleTime: defaultMaxConnIdleTime,
			maxConnLifetime: defaultMaxConnLifetime,
		},
		{
			name:            "only other parameters",
			dsn:             "postgres://user:pwd@host:port?bar=value&foo=other_value",
			cleanedDSN:      "postgres://user:pwd@host:port?bar=value&foo=other_value",
			maxConns:        defaultMaxConns,
			maxIdleConns:    defaultMaxIdleConns,
			maxConnIdleTime: defaultMaxConnIdleTime,
			maxConnLifetime: defaultMaxConnLifetime,
		},
		{
			name:            "only maxConns",
			dsn:             "postgres://user:pwd@host:port?max_conns=5254",
			cleanedDSN:      "postgres://user:pwd@host:port?",
			maxConns:        5254,
			maxIdleConns:    defaultMaxIdleConns,
			maxConnIdleTime: defaultMaxConnIdleTime,
			maxConnLifetime: defaultMaxConnLifetime,
		},
		{
			name:            "only maxIdleConns",
			dsn:             "postgres://user:pwd@host:port?max_idle_conns=9342",
			cleanedDSN:      "postgres://user:pwd@host:port?",
			maxConns:        defaultMaxConns,
			maxIdleConns:    9342,
			maxConnIdleTime: defaultMaxConnIdleTime,
			maxConnLifetime: defaultMaxConnLifetime,
		},
		{
			name:            "only maxConnIdleTime",
			dsn:             "postgres://user:pwd@host:port?max_conn_idle_time=112s",
			cleanedDSN:      "postgres://user:pwd@host:port?",
			maxConns:        defaultMaxConns,
			maxIdleConns:    defaultMaxIdleConns,
			maxConnIdleTime: 112 * time.Second,
			maxConnLifetime: defaultMaxConnLifetime,
		},
		{
			name:            "only maxConnLifetime",
			dsn:             "postgres://user:pwd@host:port?max_conn_lifetime=112s",
			cleanedDSN:      "postgres://user:pwd@host:port?",
			maxConns:        defaultMaxConns,
			maxIdleConns:    defaultMaxIdleConns,
			maxConnIdleTime: defaultMaxConnIdleTime,
			maxConnLifetime: 112 * time.Second,
		},
		{
			name:            "all parameters and others",
			dsn:             "postgres://user:pwd@host:port?max_conns=5254&max_idle_conns=9342&max_conn_lifetime=112s&bar=value&foo=other_value",
			cleanedDSN:      "postgres://user:pwd@host:port?bar=value&foo=other_value",
			maxConns:        5254,
			maxIdleConns:    9342,
			maxConnIdleTime: defaultMaxConnIdleTime,
			maxConnLifetime: 112 * time.Second,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/name=%s", i, tc.name), func(t *testing.T) {
			maxConns, maxIdleConns, maxConnLifetime, maxConnIdleTime, cleanedDSN := ParseConnectionOptions(logger, tc.dsn)
			assert.Equal(t, tc.maxConns, maxConns)
			assert.Equal(t, tc.maxIdleConns, maxIdleConns)
			assert.Equal(t, tc.maxConnLifetime, maxConnLifetime)
			assert.Equal(t, tc.maxConnIdleTime, maxConnIdleTime)
			assert.Equal(t, tc.cleanedDSN, cleanedDSN)
		})
	}
}

func TestFinalizeDSN(t *testing.T) {
	for i, tc := range []struct {
		dsn, expected string
	}{
		{
			dsn:      "mysql://localhost",
			expected: "mysql://localhost?clientFoundRows=true&multiStatements=true&parseTime=true",
		},
		{
			dsn:      "mysql://localhost?multiStatements=true&parseTime=true&clientFoundRows=false",
			expected: "mysql://localhost?clientFoundRows=true&multiStatements=true&parseTime=true",
		},
		{
			dsn:      "postgres://localhost",
			expected: "postgres://localhost",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			assert.Equal(t, tc.expected, FinalizeDSN(logrusx.New("", ""), tc.dsn))
		})
	}
}
