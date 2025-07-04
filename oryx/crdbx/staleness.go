// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package crdbx

import (
	"net/http"

	"github.com/ory/x/dbal"

	"github.com/ory/pop/v6"

	"github.com/ory/x/sqlcon"
)

// Control API consistency guarantees
//
// swagger:model consistencyRequestParameters
type ConsistencyRequestParameters struct {
	// Read Consistency Level (preview)
	//
	// The read consistency level determines the consistency guarantee for reads:
	//
	// - strong (slow): The read is guaranteed to return the most recent data committed at the start of the read.
	// - eventual (very fast): The result will return data that is about 4.8 seconds old.
	//
	// The default consistency guarantee can be changed in the Ory Network Console or using the Ory CLI with
	// `ory patch project --replace '/previews/default_read_consistency_level="strong"'`.
	//
	// Setting the default consistency level to `eventual` may cause regressions in the future as we add consistency
	// controls to more APIs. Currently, the following APIs will be affected by this setting:
	//
	// - `GET /admin/identities`
	//
	// This feature is in preview and only available in Ory Network.
	//
	// required: false
	// in: query
	Consistency ConsistencyLevel `json:"consistency"`
}

// ConsistencyLevel is the consistency level.
// swagger:enum ConsistencyLevel
type ConsistencyLevel string

const (
	// ConsistencyLevelUnset is the unset / default consistency level.
	ConsistencyLevelUnset ConsistencyLevel = ""
	// ConsistencyLevelStrong is the strong consistency level.
	ConsistencyLevelStrong ConsistencyLevel = "strong"
	// ConsistencyLevelEventual is the eventual consistency level using follower read timestamps.
	ConsistencyLevelEventual ConsistencyLevel = "eventual"
)

// ConsistencyLevelFromRequest extracts the consistency level from a request.
func ConsistencyLevelFromRequest(r *http.Request) ConsistencyLevel {
	return ConsistencyLevelFromString(r.URL.Query().Get("consistency"))
}

// ConsistencyLevelFromString converts a string to a ConsistencyLevel.
// If the string is not recognized or unset, ConsistencyLevelStrong is returned.
func ConsistencyLevelFromString(in string) ConsistencyLevel {
	switch in {
	case string(ConsistencyLevelStrong):
		return ConsistencyLevelStrong
	case string(ConsistencyLevelEventual):
		return ConsistencyLevelEventual
	case string(ConsistencyLevelUnset):
		return ConsistencyLevelUnset
	}
	return ConsistencyLevelStrong
}

// SetTransactionConsistency sets the transaction consistency level for CockroachDB.
func SetTransactionConsistency(c *pop.Connection, level ConsistencyLevel, fallback ConsistencyLevel) error {
	q := getTransactionConsistencyQuery(c.Dialect.Name(), level, fallback)
	if len(q) == 0 {
		return nil
	}

	return sqlcon.HandleError(c.RawQuery(q).Exec())
}

const transactionFollowerReadTimestamp = "SET TRANSACTION AS OF SYSTEM TIME follower_read_timestamp()"

func getTransactionConsistencyQuery(dialect string, level ConsistencyLevel, fallback ConsistencyLevel) string {
	if dialect != dbal.DriverCockroachDB {
		// Only CockroachDB supports this.
		return ""
	}

	switch level {
	case ConsistencyLevelStrong:
		// Nothing to do
		return ""
	case ConsistencyLevelEventual:
		// Jumps to end of function
	case ConsistencyLevelUnset:
		fallthrough
	default:
		if fallback != ConsistencyLevelEventual {
			// Nothing to do
			return ""
		}

		// Jumps to end of function
	}

	return transactionFollowerReadTimestamp
}
