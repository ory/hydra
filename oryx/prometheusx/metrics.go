// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

type HTTPMetrics struct {
	responseTime    *prometheus.HistogramVec
	totalRequests   *prometheus.CounterVec
	duration        *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	handlerStatuses *prometheus.CounterVec
}

const HTTPPrefix = "http"

func NewHTTPMetrics(app, metricsPrefix, version, hash, date string) *HTTPMetrics {
	labels := map[string]string{
		"app":       app,
		"version":   version,
		"hash":      hash,
		"buildTime": date,
	}

	if metricsPrefix != "" {
		metricsPrefix += "_"
	}

	pm := &HTTPMetrics{
		responseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        metricsPrefix + "response_time_seconds",
				Help:        "Description",
				ConstLabels: labels,
			},
			[]string{"endpoint"},
		),
		totalRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        metricsPrefix + "requests_total",
			Help:        "number of requests",
			ConstLabels: labels,
		}, []string{"code", "method", "endpoint"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        metricsPrefix + "requests_duration_seconds",
			Help:        "duration of a requests in seconds",
			ConstLabels: labels,
		}, []string{"code", "method", "endpoint"}),
		responseSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        metricsPrefix + "response_size_bytes",
			Help:        "size of the responses in bytes",
			ConstLabels: labels,
		}, []string{"code", "method"}),
		requestSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        metricsPrefix + "requests_size_bytes",
			Help:        "size of the requests in bytes",
			ConstLabels: labels,
		}, []string{"code", "method"}),
		handlerStatuses: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        metricsPrefix + "requests_statuses_total",
			Help:        "count number of responses per status",
			ConstLabels: labels,
		}, []string{"method", "status_bucket"}),
	}

	err := prometheus.Register(pm)
	if e := new(prometheus.AlreadyRegisteredError); errors.As(err, e) {
		return pm
	} else if err != nil {
		panic(err)
	}

	grpcPrometheus.EnableHandlingTimeHistogram()

	return pm
}

// Describe implements prometheus Collector interface.
func (h *HTTPMetrics) Describe(in chan<- *prometheus.Desc) {
	h.duration.Describe(in)
	h.totalRequests.Describe(in)
	h.requestSize.Describe(in)
	h.responseSize.Describe(in)
	h.handlerStatuses.Describe(in)
	h.responseTime.Describe(in)
}

// Collect implements prometheus Collector interface.
func (h *HTTPMetrics) Collect(in chan<- prometheus.Metric) {
	h.duration.Collect(in)
	h.totalRequests.Collect(in)
	h.requestSize.Collect(in)
	h.responseSize.Collect(in)
	h.handlerStatuses.Collect(in)
	h.responseTime.Collect(in)
}

func (h *HTTPMetrics) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rr := negroni.NewResponseWriter(rw)
	start := time.Now()

	next(rr, r)

	latency := time.Since(start)
	code := StatusCodeToString(rr.Status())
	method := sanitizeMethod(r.Method)
	endpoint := GetLabelForPattern(r.Pattern)

	h.responseSize.WithLabelValues(code, method).Observe(float64(rr.Size()))
	h.totalRequests.WithLabelValues(code, method, endpoint).Inc()
	h.duration.WithLabelValues(code, method, endpoint).Observe(latency.Seconds())
	h.responseTime.WithLabelValues(endpoint).Observe(latency.Seconds())
	h.requestSize.WithLabelValues(code, method).Observe(float64(computeApproximateRequestSize(r)))

	statusBucket := "unknown"
	switch status := rr.Status(); {
	case status >= 200 && status <= 299:
		statusBucket = "2xx"
	case status >= 300 && status <= 399:
		statusBucket = "3xx"
	case status >= 400 && status <= 499:
		statusBucket = "4xx"
	case status >= 500 && status <= 599:
		statusBucket = "5xx"
	}

	h.handlerStatuses.WithLabelValues(r.Method, statusBucket).Inc()
}

var (
	paramPlaceHolderRE = regexp.MustCompile(`\{[a-zA-Z0-9_-]+}`)
	// patternLabelCache is used to cache the generated label for a given pattern to avoid recomputing it on every request.
	// We use a sync.Map over a ristretto cache as the number of unique patterns is expected to be low and the cache will not grow indefinitely.
	//
	// > The sync.Map type is optimized for two common use cases: (1) when the entry for a given
	// > key is only ever written once but read many times, as in caches that only grow,
	// > or (2) when multiple goroutines read, write, and overwrite entries for disjoint
	// > sets of keys. In these two cases, use of a Map may significantly reduce lock
	// > contention compared to a Go map paired with a separate [Mutex] or [RWMutex].
	// https://pkg.go.dev/sync#Map
	patternLabelCache = sync.Map{}
)

func GetLabelForPattern(pattern string) string {
	if label, ok := patternLabelCache.Load(pattern); ok {
		return label.(string)
	}

	cleanedPattern := pattern
	// remove the method if it is included
	if _, pattern, ok := strings.Cut(cleanedPattern, " "); ok {
		cleanedPattern = pattern
	}
	// trim space, just to be sure
	cleanedPattern = strings.TrimSpace(cleanedPattern)

	label := paramPlaceHolderRE.ReplaceAllString(strings.TrimSuffix(cleanedPattern, "/{$}"), "{param}")
	// Cache the generated label for the pattern to avoid recomputing it on every request.
	patternLabelCache.Store(pattern, label)
	return label
}

// computeApproximateRequestSize is copied from the promhttp package as it is not exported
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
