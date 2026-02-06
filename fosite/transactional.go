// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "context"

// Transactional is an interface that a storage provider has to implement to ensure atomicity for certain flows
// that require transactional semantics. If the storage provider cannot provide transactional semantics,
// it still needs to implement this interface as a no-op.
// It is expected that the storage provider will examine context for an existing transaction
// each time a database operation is to be performed.
//
// An implementation of Transaction should attempt to initiate a new transaction and store that under a unique key
// in the context so it will be passed on by fn to all storage provider methods. If fn returns an error,
// the transaction must be rolled back, otherwise it must be committed. It is permitted to automatically retry
// the transaction in case of specific errors. Fosite will always pass an idempotent function to Transaction.
type Transactional interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
