// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetsecure

import (
	"context"
	"os"
	"runtime"
	"testing"
)

type (
	VMProvider interface {
		// JsonnetVM creates a new secure process-isolated Jsonnet VM whose
		// execution is bound to the provided context, i.e.,
		// cancelling the context will terminate the VM process.
		JsonnetVM(context.Context) (VM, error)
	}

	// TestProvider provides a secure VM by running go build on github.
	// com/ory/x/jsonnetsecure/cmd.
	TestProvider struct {
		jsonnetBinary string
		pool          Pool
	}

	// DefaultProvider provides a secure VM by calling the currently
	// running the current binary with the provided subcommand.
	DefaultProvider struct {
		Subcommand string
		Pool       Pool
	}
)

func NewTestProvider(t testing.TB) *TestProvider {
	pool := NewProcessPool(runtime.GOMAXPROCS(0))
	t.Cleanup(pool.Close)
	return &TestProvider{JsonnetTestBinary(t), pool}
}

func (p *TestProvider) JsonnetVM(ctx context.Context) (VM, error) {
	return MakeSecureVM(
		WithProcessPool(p.pool),
		WithJsonnetBinary(p.jsonnetBinary),
	), nil
}

func (p *DefaultProvider) JsonnetVM(ctx context.Context) (VM, error) {
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	return MakeSecureVM(
		WithJsonnetBinary(self),
		WithProcessArgs(p.Subcommand),
		WithProcessPool(p.Pool),
	), nil
}
