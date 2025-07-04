// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package otelx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestShouldNotTraceHealthEndpoint(t *testing.T) {
	testCases := []struct {
		path            string
		testDescription string
	}{
		{
			path:            "health/ready",
			testDescription: "health",
		},
		{
			path:            "admin/alive",
			testDescription: "adminHealth",
		},
		{
			path:            "foo/bar",
			testDescription: "notHealth",
		},
	}
	for _, test := range testCases {
		t.Run(test.testDescription, func(t *testing.T) {
			recorder := tracetest.NewSpanRecorder()
			tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

			req := httptest.NewRequest(http.MethodGet, "https://api.example.com/"+test.path, nil)
			h := NewHandler(negroni.New(), "test op", otelhttp.WithTracerProvider(tp))
			h.ServeHTTP(negroni.NewResponseWriter(httptest.NewRecorder()), req)

			spans := recorder.Ended()
			if strings.Contains(test.path, "health") {
				assert.Len(t, spans, 0)
			} else {
				assert.Len(t, spans, 1)
			}
		})
	}
}

func TestTraceHandlerSpanName(t *testing.T) {
	testCases := []struct {
		path         string
		expectedName string
		opts         []otelhttp.Option
	}{
		{
			path:         "testPath",
			expectedName: "/testPath",
			opts:         []otelhttp.Option{},
		},
		{
			path:         "testPath",
			expectedName: "/overwritten/name",
			opts: []otelhttp.Option{
				otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
					return "/overwritten/name"
				}),
			},
		},
	}
	for _, test := range testCases {
		recorder := tracetest.NewSpanRecorder()
		tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

		opts := append([]otelhttp.Option{
			otelhttp.WithTracerProvider(tp),
		}, test.opts...)

		req := httptest.NewRequest(http.MethodGet, "https://api.example.com/"+test.path, nil)
		h := TraceHandler(negroni.New(), opts...)
		h.ServeHTTP(negroni.NewResponseWriter(httptest.NewRecorder()), req)

		spans := recorder.Ended()
		assert.Len(t, spans, 1)
		assert.Equal(t, test.expectedName, spans[0].Name())
	}
}
