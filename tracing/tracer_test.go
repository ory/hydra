package tracing_test

import (
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/tracing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type zipkinSpanRequest struct {
	Id            string
	TraceId       string
	Timestamp     uint64
	Name          string
	LocalEndpoint struct {
		ServiceName string
	}
	Tags map[string]string
}

func TestZipkinTracer(t *testing.T) {
	done := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var spans []zipkinSpanRequest
		err = json.Unmarshal(body, &spans)
		assert.NoError(t, err)

		assert.NotEmpty(t, spans[0].Id)
		assert.NotEmpty(t, spans[0].TraceId)
		assert.Equal(t, "testOperation", spans[0].Name)
		assert.Equal(t, "Hydra", spans[0].LocalEndpoint.ServiceName)
		assert.NotNil(t, spans[0].Tags["testTag"])
		assert.Equal(t, "true", spans[0].Tags["testTag"])
	}))
	defer ts.Close()

	tracer := &tracing.Tracer{
		ServiceName: "Hydra",
		ZipkinConfig: &tracing.ZipkinConfig{
			ServerURL: ts.URL,
		},
		Provider: "zipkin",
		Logger:   logrusx.New("Hydra", "1"),
	}
	err := tracer.Setup()
	assert.NoError(t, err)

	span := opentracing.GlobalTracer().StartSpan("testOperation")
	span.SetTag("testTag", true)
	span.Finish()

	select {
	case <-done:
	case <-time.After(time.Millisecond * 1500):
		t.Fatalf("Test server did not receive spans")
	}
}
