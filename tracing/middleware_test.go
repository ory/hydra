package tracing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"

	"github.com/ory/hydra/tracing"
)

var mockedTracer *mocktracer.MockTracer
var tracer *tracing.Tracer = &tracing.Tracer{
	ServiceName: "Ory Hydra Test",
	Provider:    "Mock Provider",
}

func init() {
	mockedTracer = mocktracer.New()
	opentracing.SetGlobalTracer(mockedTracer)
}

func TestTracingServeHttp(t *testing.T) {
	expectedTagsSuccess := map[string]interface{}{
		string(ext.HTTPStatusCode): uint16(200),
		string(ext.HTTPMethod):     "GET",
	}

	expectedTagsError := map[string]interface{}{
		string(ext.HTTPStatusCode): uint16(400),
		string(ext.HTTPMethod):     "GET",
		"error":                    true,
	}

	testCases := []struct {
		httpStatus      int
		testDescription string
		expectedTags    map[string]interface{}
	}{
		{
			testDescription: "success http response",
			httpStatus:      http.StatusOK,
			expectedTags:    expectedTagsSuccess,
		},
		{
			testDescription: "error http response",
			httpStatus:      http.StatusBadRequest,
			expectedTags:    expectedTagsError,
		},
	}

	for _, test := range testCases {
		t.Run(test.testDescription, func(t *testing.T) {
			defer mockedTracer.Reset()
			request := httptest.NewRequest(http.MethodGet, "https://apis.somecompany.com/endpoint", nil)
			next := func(rw http.ResponseWriter, _ *http.Request) {
				rw.WriteHeader(test.httpStatus)
			}

			tracer.ServeHTTP(negroni.NewResponseWriter(httptest.NewRecorder()), request, next)

			spans := mockedTracer.FinishedSpans()
			assert.Len(t, spans, 1)
			span := spans[0]

			assert.Equal(t, test.expectedTags, span.Tags())
		})
	}
}

func TestShouldContinueTraceIfAlreadyPresent(t *testing.T) {
	defer mockedTracer.Reset()
	parentSpan := mockedTracer.StartSpan("some-operation").(*mocktracer.MockSpan)
	ext.SpanKindRPCClient.Set(parentSpan)
	request := httptest.NewRequest(http.MethodGet, "https://apis.somecompany.com/endpoint", nil)
	carrier := opentracing.HTTPHeadersCarrier(request.Header)
	// this request now contains a trace initiated by another service/process (e.g. an edge proxy that fronts Hydra)
	mockedTracer.Inject(parentSpan.Context(), opentracing.HTTPHeaders, carrier)

	next := func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}

	tracer.ServeHTTP(negroni.NewResponseWriter(httptest.NewRecorder()), request, next)

	spans := mockedTracer.FinishedSpans()
	assert.Len(t, spans, 1)
	span := spans[0]

	assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
}
