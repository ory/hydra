// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !sqlite
// +build !sqlite

package sqlcon

// handleSqlite handles the error iff (if and only if) it is an sqlite error
func handleSqlite(_ error) error {
	return nil
}
