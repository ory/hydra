// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "context"

// A storage provider that has support for transactions should implement this interface to ensure atomicity for certain flows
// that require transactional semantics. Fosite will call these methods (when atomicity is required) if and only if the storage
// provider has implemented `Transactional`. It is expected that the storage provider will examine context for an existing transaction
// each time a database operation is to be performed.
//
// An implementation of `BeginTX` should attempt to initiate a new transaction and store that under a unique key
// in the context that can be accessible by `Commit` and `Rollback`. The "transactional aware" context will then be
// returned for further propagation, eventually to be consumed by `Commit` or `Rollback` to finish the transaction.
//
// Implementations for `Commit` & `Rollback` should look for the transaction object inside the supplied context using the same
// key used by `BeginTX`. If these methods have been called, it is expected that a txn object should be available in the provided
// context.
type Transactional interface {
	BeginTX(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// MaybeBeginTx is a helper function that can be used to initiate a transaction if the supplied storage
// implements the `Transactional` interface.
func MaybeBeginTx(ctx context.Context, storage interface{}) (context.Context, error) {
	// the type assertion checks whether the dynamic type of `storage` implements `Transactional`
	txnStorage, transactional := storage.(Transactional)
	if transactional {
		return txnStorage.BeginTX(ctx)
	} else {
		return ctx, nil
	}
}

// MaybeCommitTx is a helper function that can be used to commit a transaction if the supplied storage
// implements the `Transactional` interface.
func MaybeCommitTx(ctx context.Context, storage interface{}) error {
	txnStorage, transactional := storage.(Transactional)
	if transactional {
		return txnStorage.Commit(ctx)
	} else {
		return nil
	}
}

// MaybeRollbackTx is a helper function that can be used to rollback a transaction if the supplied storage
// implements the `Transactional` interface.
func MaybeRollbackTx(ctx context.Context, storage interface{}) error {
	txnStorage, transactional := storage.(Transactional)
	if transactional {
		return txnStorage.Rollback(ctx)
	} else {
		return nil
	}
}
