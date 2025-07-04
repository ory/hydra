// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package profilex

import (
	"os"

	"github.com/pkg/profile"
)

type noop struct{}

// Stop is a noop.
func (p *noop) Stop() {}

// Profile parses the PROFILING environment variable and executes the proper profiling task.
func Profile() interface {
	Stop()
} {
	switch os.Getenv("PROFILING") {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.NoShutdownHook)
	case "mem":
		return profile.Start(profile.MemProfile, profile.NoShutdownHook)
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.NoShutdownHook)
	case "block":
		return profile.Start(profile.BlockProfile, profile.NoShutdownHook)
	}
	return new(noop)
}

// HelpMessage returns a string explaining how profiling works.
func HelpMessage() string {
	return `- PROFILING: Set "PROFILING=cpu" to enable cpu profiling and "PROFILING=mem" to enable memory profiling.
	It is not possible to do both at the same time. Profiling is disabled per default.

	Example: PROFILING=cpu`
}
