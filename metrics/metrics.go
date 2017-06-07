package metrics

import "time"

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
	if latency > time.Second * 5 {
		latency = time.Second * 5
	}

	// milliseconds / 10
	h.Latencies[int64(latency / 10)]++
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
	*Metrics
	*HTTPMetrics
	Paths map[string]*PathMetrics `json:"paths"`
}

func newMetrics() *Metrics {
	return &Metrics{
		Latencies: map[int64]int{},
	}
}

func (s *Snapshot) Path(path string) *PathMetrics {
	paths := []string{
		"/.well-known/jwks.json",
		"/.well-known/openid-configuration",
		"/clients",
		"/health",
		"/keys",
		"/oauth2/auth",
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
