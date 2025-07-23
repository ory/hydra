// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"net/http"
	"regexp"
	"strings"
	"sync"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type MetricsManager struct {
	prometheusMetrics *Metrics
	routers           struct {
		data []*httprouter.Router
		sync.Mutex
	}
}

func NewMetricsManager(app, version, hash, buildTime string) *MetricsManager {
	return NewMetricsManagerWithPrefix(app, "", version, hash, buildTime)
}

// NewMetricsManagerWithPrefix creates MetricsManager that uses metricsPrefix parameters as a prefix
// for all metrics registered within this middleware. Constants HttpMetrics or GrpcMetrics can be used
// respectively. Setting empty string in metricsPrefix will be equivalent to calling NewMetricsManager.
func NewMetricsManagerWithPrefix(app, metricsPrefix, version, hash, buildTime string) *MetricsManager {
	return &MetricsManager{
		prometheusMetrics: NewMetrics(app, metricsPrefix, version, hash, buildTime),
	}
}

// Main middleware method to collect metrics for Prometheus.
func (pmm *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	pmm.prometheusMetrics.Instrument(rw, next, pmm.getLabelForPath(r))(rw, r)
}

func (pmm *MetricsManager) StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	f := grpcPrometheus.StreamServerInterceptor
	return f(srv, ss, info, handler)
}

func (pmm *MetricsManager) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	f := grpcPrometheus.UnaryServerInterceptor
	return f(ctx, req, info, handler)
}

func (pmm *MetricsManager) Register(server *grpc.Server) {
	grpcPrometheus.Register(server)
}

func (pmm *MetricsManager) RegisterRouter(router *httprouter.Router) {
	pmm.routers.Lock()
	defer pmm.routers.Unlock()
	pmm.routers.data = append(pmm.routers.data, router)
}

var paramPlaceHolderRE = regexp.MustCompile(`\{[a-zA-Z0-9_-]+\}`)

func (pmm *MetricsManager) getLabelForPath(r *http.Request) string {
	// If the request came through a http.ServeMux, it already has a pattern that we
	// can use as a label. We just need to replace all path parameters with a generic
	// placeholder and remove the trailing slash pattern.
	if p := r.Pattern; p != "" {
		return paramPlaceHolderRE.ReplaceAllString(strings.TrimSuffix(p, "/{$}"), "{param}")
	}

	// looking for a match in one of registered routers
	pmm.routers.Lock()
	defer pmm.routers.Unlock()
	for _, router := range pmm.routers.data {
		handler, params, _ := router.Lookup(r.Method, r.URL.Path)
		if handler != nil {
			return reconstructEndpoint(r.URL.Path, params)
		}
	}
	return "{unmatched}"
}

// To reduce cardinality of labels, values of matched path parameters must be replaced with {param}
func reconstructEndpoint(path string, params httprouter.Params) string {
	// if map is empty, then nothing to change in the path
	if len(params) == 0 {
		return path
	}

	// construct a list of parameter values
	paramValues := make(map[string]struct{}, len(params))
	for _, param := range params {
		paramValues[param.Value] = struct{}{}
	}

	parts := strings.Split(path, "/")
	for index, part := range parts {
		if _, ok := paramValues[part]; ok {
			parts[index] = "{param}"
		}
	}

	return strings.Join(parts, "/")
}
