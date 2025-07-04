// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package metricsx

import (
	"runtime"
	"sync"
)

// MemoryStatistics is a JSON-able version of runtime.MemStats
type MemoryStatistics struct {
	sync.Mutex
	// Alloc is bytes of allocated heap objects.
	Alloc uint64 `json:"alloc"`
	// TotalAlloc is cumulative bytes allocated for heap objects.
	TotalAlloc uint64 `json:"totalAlloc"`
	// Sys is the total bytes of memory obtained from the OS.
	Sys uint64 `json:"sys"`
	// Lookups is the number of pointer lookups performed by the
	// runtime.
	Lookups uint64 `json:"lookups"`
	// Mallocs is the cumulative count of heap objects allocated.
	// The number of live objects is Mallocs - Frees.
	Mallocs uint64 `json:"mallocs"`
	// Frees is the cumulative count of heap objects freed.
	Frees uint64 `json:"frees"`
	// HeapAlloc is bytes of allocated heap objects.
	HeapAlloc uint64 `json:"heapAlloc"`
	// HeapSys is bytes of heap memory obtained from the OS.
	HeapSys uint64 `json:"heapSys"`
	// HeapIdle is bytes in idle (unused) spans.
	HeapIdle uint64 `json:"heapIdle"`
	// HeapInuse is bytes in in-use spans.
	HeapInuse uint64 `json:"heapInuse"`
	// HeapReleased is bytes of physical memory returned to the OS.
	HeapReleased uint64 `json:"heapReleased"`
	// HeapObjects is the number of allocated heap objects.
	HeapObjects uint64 `json:"heapObjects"`
	// NumGC is the number of completed GC cycles.
	NumGC uint32 `json:"numGC"`
}

// ToMap converts to a map[string]interface{}.
func (ms *MemoryStatistics) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"alloc":          ms.Alloc,
		"totalAlloc":     ms.TotalAlloc,
		"sys":            ms.Sys,
		"lookups":        ms.Lookups,
		"mallocs":        ms.Mallocs,
		"frees":          ms.Frees,
		"heapAlloc":      ms.HeapAlloc,
		"heapSys":        ms.HeapSys,
		"heapIdle":       ms.HeapIdle,
		"heapInuse":      ms.HeapInuse,
		"heapReleased":   ms.HeapReleased,
		"heapObjects":    ms.HeapObjects,
		"numGC":          ms.NumGC,
		"nonInteraction": 1,
	}
}

// Update takes the most recent stats from runtime.
func (ms *MemoryStatistics) Update() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ms.Lock()
	defer ms.Unlock()
	ms.Alloc = m.Alloc
	ms.TotalAlloc = m.TotalAlloc
	ms.Sys = m.Sys
	ms.Lookups = m.Lookups
	ms.Mallocs = m.Mallocs
	ms.Frees = m.Frees
	ms.HeapAlloc = m.HeapAlloc
	ms.HeapSys = m.HeapSys
	ms.HeapIdle = m.HeapIdle
	ms.HeapInuse = m.HeapInuse
	ms.HeapReleased = m.HeapReleased
	ms.HeapObjects = m.HeapObjects
	ms.NumGC = m.NumGC
}
