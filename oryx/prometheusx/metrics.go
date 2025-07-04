// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"net/http"
	"strconv"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/ory/x/httpx"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics prototypes
type Metrics struct {
	responseTime    *prometheus.HistogramVec
	totalRequests   *prometheus.CounterVec
	duration        *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	handlerStatuses *prometheus.CounterVec
}

const HTTPMetrics = "http"
const GRPCMetrics = "grpc"

// NewMetrics creates new custom Prometheus metrics
func NewMetrics(app, metricsPrefix, version, hash, date string) *Metrics {
	labels := map[string]string{
		"app":       app,
		"version":   version,
		"hash":      hash,
		"buildTime": date,
	}

	if metricsPrefix != "" {
		metricsPrefix += "_"
	}

	pm := &Metrics{
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
func (h *Metrics) Describe(in chan<- *prometheus.Desc) {
	h.duration.Describe(in)
	h.totalRequests.Describe(in)
	h.requestSize.Describe(in)
	h.responseSize.Describe(in)
	h.handlerStatuses.Describe(in)
	h.responseTime.Describe(in)
}

// Collect implements prometheus Collector interface.
func (h *Metrics) Collect(in chan<- prometheus.Metric) {
	h.duration.Collect(in)
	h.totalRequests.Collect(in)
	h.requestSize.Collect(in)
	h.responseSize.Collect(in)
	h.handlerStatuses.Collect(in)
	h.responseTime.Collect(in)
}

func (h Metrics) instrumentHandlerStatusBucket(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)

		status, _ := httpx.GetResponseMeta(rw)

		statusBucket := "unknown"
		switch {
		case status >= 200 && status <= 299:
			statusBucket = "2xx"
		case status >= 300 && status <= 399:
			statusBucket = "3xx"
		case status >= 400 && status <= 499:
			statusBucket = "4xx"
		case status >= 500 && status <= 599:
			statusBucket = "5xx"
		}

		h.handlerStatuses.With(prometheus.Labels{"method": r.Method, "status_bucket": statusBucket}).
			Inc()
	}
}

// Instrument will instrument any http.HandlerFunc with custom metrics
func (h Metrics) Instrument(rw http.ResponseWriter, next http.HandlerFunc, endpoint string) http.HandlerFunc {
	labels := prometheus.Labels{}
	labelsWithEndpoint := prometheus.Labels{"endpoint": endpoint}
	if status, _ := httpx.GetResponseMeta(rw); status != 0 {
		labels = prometheus.Labels{"code": strconv.Itoa(status)}
		labelsWithEndpoint["code"] = labels["code"]
	}
	wrapped := promhttp.InstrumentHandlerResponseSize(h.responseSize.MustCurryWith(labels), next)
	wrapped = promhttp.InstrumentHandlerCounter(h.totalRequests.MustCurryWith(labelsWithEndpoint), wrapped)
	wrapped = promhttp.InstrumentHandlerDuration(h.duration.MustCurryWith(labelsWithEndpoint), wrapped)
	wrapped = promhttp.InstrumentHandlerDuration(h.responseTime.MustCurryWith(prometheus.Labels{"endpoint": endpoint}), wrapped)
	wrapped = promhttp.InstrumentHandlerRequestSize(h.requestSize.MustCurryWith(labels), wrapped)
	wrapped = h.instrumentHandlerStatusBucket(wrapped)

	return wrapped.ServeHTTP
}
