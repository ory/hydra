// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	Requests  uint64        `json:"requests,omitempty"`
	Responses uint64        `json:"responses,omitempty"`
	Latencies map[int64]int `json:"latencies,omitempty"`
}

func (h *Metrics) AddRequest() {
	h.Requests++
}

func (h *Metrics) AddResponse() {
	h.Responses++
}

type HTTPMetrics struct {
	Methods map[string]*Metrics `json:"methods"`
	Status  map[int]*Metrics    `json:"status"`
	Sizes   map[int]*Metrics    `json:"sizes"`
}

func (h *HTTPMetrics) AddMethodResponse(method string) {
	h.addMethod(method, 0, 1)
}

func (h *HTTPMetrics) AddMethodRequest(method string) {
	h.addMethod(method, 1, 0)
}

func (h *Metrics) AddLatency(latency time.Duration) {
	h.Latencies[int64(latency)]++
}

func (h *HTTPMetrics) SizeMetrics(size int) *Metrics {
	if size > 5*1024 {
		size = 5 * 1024
	}

	if _, ok := h.Sizes[size]; !ok {
		h.Sizes[size] = newMetrics()
	}
	return h.Sizes[size]
}

func (h *HTTPMetrics) StatusMetrics(status int) *Metrics {
	if _, ok := h.Status[status]; !ok {
		h.Status[status] = newMetrics()
	}
	return h.Status[status]
}

func (h *HTTPMetrics) MethodMetrics(method string) *Metrics {
	if _, ok := h.Methods[method]; !ok {
		h.Methods[method] = newMetrics()
	}
	return h.Methods[method]
}

func (h *HTTPMetrics) AddStatus(status int) {
	if _, ok := h.Status[status]; !ok {
		h.Status[status] = newMetrics()
	}
	h.Status[status].Responses++
}

func (h *HTTPMetrics) AddSize(size int) {
	h.SizeMetrics(size).Responses++
}

func (h *HTTPMetrics) addMethod(method string, req, res uint64) {
	if _, ok := h.Methods[method]; !ok {
		h.Methods[method] = newMetrics()
	}
	h.Methods[method].Requests = h.Methods[method].Requests + res
	h.Methods[method].Responses = h.Methods[method].Responses + req
}

type PathMetrics struct {
	*Metrics
	*HTTPMetrics
}

type Snapshot struct {
	sync.RWMutex
	*Metrics
	*HTTPMetrics
	Paths          map[string]*PathMetrics `json:"paths"`
	ID             string                  `json:"id"`
	UpTime         int64                   `json:"uptime"`
	start          time.Time               `json:"-"`
	MemorySnapshot *MemorySnapshot         `json:"memory"`
}

type MemorySnapshot struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"totalAlloc"`
	Sys          uint64 `json:"sys"`
	Lookups      uint64 `json:"lookups"`
	Mallocs      uint64 `json:"mallocs"`
	Frees        uint64 `json:"frees"`
	HeapAlloc    uint64 `json:"heapAlloc"`
	HeapSys      uint64 `json:"heapSys"`
	HeapIdle     uint64 `json:"heapIdle"`
	HeapInuse    uint64 `json:"heapInuse"`
	HeapReleased uint64 `json:"heapReleased"`
	HeapObjects  uint64 `json:"heapObjects"`
	NumGC        uint32 `json:"numGC"`
}

func newMetrics() *Metrics {
	return &Metrics{
		Latencies: map[int64]int{},
	}
}

func (sw *Snapshot) GetUpTime() int64 {
	sw.Update()
	sw.RLock()
	defer sw.RUnlock()
	return sw.UpTime
}

func (sw *Snapshot) Update() {
	sw.Lock()
	defer sw.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// sw.MemorySnapshot = &(MemorySnapshot(m))
	sw.MemorySnapshot = &MemorySnapshot{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		Lookups:      m.Lookups,
		Mallocs:      m.Mallocs,
		Frees:        m.Frees,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		HeapObjects:  m.HeapObjects,
		NumGC:        m.NumGC,
	}

	sw.UpTime = int64(time.Now().UTC().Sub(sw.start) / time.Second)

}

func (s *Snapshot) Path(path string) *PathMetrics {
	paths := []string{
		"/.well-known/jwks.json",
		"/.well-known/openid-configuration",
		"/clients",
		"/health",
		"/keys",
		"/oauth2/auth",
		"/oauth2/session",
		"/oauth2/consent",
		"/oauth2/introspect",
		"/oauth2/revoke",
		"/oauth2/token",
		"/policies",
		"/warden/allowed",
		"/warden/groups",
		"/warden/token/allowed",
		"/",
	}

	for _, p := range paths {
		if len(path) >= len(p) && path[:len(p)] == p {
			path = p
			break
		}
	}

	if _, ok := s.Paths[path]; !ok {
		s.Paths[path] = &PathMetrics{
			Metrics:     newMetrics(),
			HTTPMetrics: newHttpMetrics(),
		}
	}

	return s.Paths[path]
}

func newHttpMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		Methods: map[string]*Metrics{},
		Status:  map[int]*Metrics{},
		Sizes:   map[int]*Metrics{},
	}
}

func newPathMetrics() *PathMetrics {
	return &PathMetrics{
		Metrics:     newMetrics(),
		HTTPMetrics: newHttpMetrics(),
	}
}
