// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package jsonnetsecure

import "runtime/debug"

func SetVirtualMemoryLimit(limitBytes uint64) error {
	// Tell the Go runtime about the limit.
	debug.SetMemoryLimit(int64(limitBytes)) //nolint:gosec // The number is a compile-time constant.

	// TODO No OS limit for now. Apparently there is a Windows-specific equivalent (Job control)?
	return nil
}
