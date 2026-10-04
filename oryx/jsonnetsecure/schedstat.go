// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetsecure

import (
	"bytes"
	"os"
	"strconv"
	"time"
)

// schedstat holds the scheduler statistics of a process as reported by the
// Linux kernel in /proc/<pid>/schedstat: the CPU time the process actually
// spent running, and the time it spent runnable but waiting for a CPU (due to
// contention or cgroup CPU throttling).
type schedstat struct {
	cpuTime      time.Duration
	runqueueWait time.Duration
}

// readSchedstat returns the cumulative scheduler statistics for pid. It
// reports ok=false on platforms without /proc/<pid>/schedstat (e.g. macOS,
// Windows) and on Linux kernels built without CONFIG_SCHED_INFO.
func readSchedstat(pid int) (schedstat, bool) {
	raw, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/schedstat")
	if err != nil {
		return schedstat{}, false
	}
	return parseSchedstat(raw)
}

// parseSchedstat parses the /proc/<pid>/schedstat format: three
// space-separated integers (run time in nanoseconds, runqueue wait time in
// nanoseconds, number of timeslices).
func parseSchedstat(raw []byte) (schedstat, bool) {
	fields := bytes.Fields(raw)
	if len(fields) < 2 {
		return schedstat{}, false
	}
	cpuTime, err := strconv.ParseInt(string(fields[0]), 10, 64)
	if err != nil || cpuTime < 0 {
		return schedstat{}, false
	}
	runqueueWait, err := strconv.ParseInt(string(fields[1]), 10, 64)
	if err != nil || runqueueWait < 0 {
		return schedstat{}, false
	}
	return schedstat{
		cpuTime:      time.Duration(cpuTime),
		runqueueWait: time.Duration(runqueueWait),
	}, true
}
