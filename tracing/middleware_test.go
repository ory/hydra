package tracing

import (
	"net/http"
	"net/http/httptest"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"

	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestTracingServeHttp(t *testing.T) {
	opentracing.SetGlobalTracer(mocktracer.New())
	tracer := &Tracer{
		ServiceName: "Ory Hydra Test",
		Provider:    "Mock Provider",
		tracer:      opentracing.GlobalTracer(),
	}

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
			request := httptest.NewRequest(http.MethodGet, "https://apis.somecompany.com/endpoint", nil)
			next := func(rw http.ResponseWriter, _ *http.Request) {
				rw.WriteHeader(test.httpStatus)
			}

			tracer.ServeHTTP(negroni.NewResponseWriter(httptest.NewRecorder()), request, next)

			mockTracer := tracer.tracer.(*mocktracer.MockTracer)
			spans := mockTracer.FinishedSpans()
			assert.Len(t, spans, 1)
			span := spans[0]

			assert.Equal(t, test.expectedTags, span.Tags())
			mockTracer.Reset()
		})
	}
}
