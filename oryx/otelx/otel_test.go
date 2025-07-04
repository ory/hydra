// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package otelx

import (
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	tracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"github.com/ory/x/logrusx"
)

const testTracingComponent = "github.com/ory/x/otelx"

func decodeResponseBody(t *testing.T, r *http.Request) []byte {
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			t.Fatal(err)
		}
	case "deflate":
		var err error
		reader, err = zlib.NewReader(r.Body)
		if err != nil {
			t.Fatal(err)
		}

	default:
		reader = r.Body
	}
	respBody, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	return respBody
}

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

// runTestJaegerAgent starts a mock server listening on a random port for Jaeger spans sent over UDP.
func runTestJaegerAgent(t *testing.T, errs *errgroup.Group, done chan<- struct{}) net.Conn {
	addr := "127.0.0.1:0"

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	require.NoError(t, err)

	srv, err := net.ListenUDP("udp", udpAddr)
	require.NoError(t, err)

	errs.Go(func() error {
		t.Logf("Starting test UDP server for Jaeger spans on %s", srv.LocalAddr().String())

		for {
			buf := make([]byte, 2048)
			_, conn, err := srv.ReadFromUDP(buf)
			if err != nil {
				return err
			}

			if conn == nil {
				continue
			}
			if len(buf) != 0 {
				t.Log("received span!")
				done <- struct{}{}
			}
			break
		}
		return nil
	})

	return srv
}

func TestJaegerTracer(t *testing.T) {
	done := make(chan struct{})
	errs := errgroup.Group{}

	srv := runTestJaegerAgent(t, &errs, done)

	jt, err := New(testTracingComponent, logrusx.New("ory/x", "1"), &Config{
		ServiceName: "Ory X",
		Provider:    "jaeger",
		Providers: ProvidersConfig{
			Jaeger: JaegerConfig{
				LocalAgentAddress: srv.LocalAddr().String(),
				Sampling: JaegerSampling{
					TraceIdRatio: 1,
				},
			},
		},
	})
	require.NoError(t, err)

	trc := jt.Tracer()
	_, span := trc.Start(context.Background(), "testSpan")
	span.SetAttributes(attribute.Bool("testAttribute", true))
	span.End()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatalf("Test server did not receive spans")
	}
	require.NoError(t, errs.Wait())
}

func TestJaegerTracerRespectsParentSamplingDecision(t *testing.T) {
	done := make(chan struct{})
	errs := errgroup.Group{}

	srv := runTestJaegerAgent(t, &errs, done)

	jt, err := New(testTracingComponent, logrusx.New("ory/x", "1"), &Config{
		ServiceName: "Ory X",
		Provider:    "jaeger",
		Providers: ProvidersConfig{
			Jaeger: JaegerConfig{
				LocalAgentAddress: srv.LocalAddr().String(),
				Sampling: JaegerSampling{
					// Effectively disable local sampling.
					TraceIdRatio: 0,
				},
			},
		},
	})
	require.NoError(t, err)

	traceId := strings.Repeat("a", 32)
	spanId := strings.Repeat("b", 16)
	sampledFlag := "1"
	traceHeaders := map[string]string{"uber-trace-id": traceId + ":" + spanId + ":0:" + sampledFlag}

	ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(traceHeaders))
	spanContext := trace.SpanContextFromContext(ctx)

	assert.True(t, spanContext.IsValid())
	assert.True(t, spanContext.IsSampled())
	assert.True(t, spanContext.IsRemote())

	trc := jt.Tracer()
	_, span := trc.Start(ctx, "testSpan", trace.WithLinks(trace.Link{SpanContext: spanContext}))
	span.SetAttributes(attribute.Bool("testAttribute", true))
	span.End()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatalf("Test server did not receive spans")
	}
	require.NoError(t, errs.Wait())
}

func TestZipkinTracer(t *testing.T) {
	done := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var spans []zipkinSpanRequest
		err = json.Unmarshal(body, &spans)

		assert.NoError(t, err)

		assert.NotEmpty(t, spans[0].Id)
		assert.NotEmpty(t, spans[0].TraceId)
		assert.Equal(t, "testspan", spans[0].Name)
		assert.Equal(t, "ory x", spans[0].LocalEndpoint.ServiceName)
		assert.NotNil(t, spans[0].Tags["testTag"])
		assert.Equal(t, "true", spans[0].Tags["testTag"])
	}))
	defer ts.Close()

	zt, err := New(testTracingComponent, logrusx.New("ory/x", "1"), &Config{
		ServiceName: "Ory X",
		Provider:    "zipkin",
		Providers: ProvidersConfig{
			Zipkin: ZipkinConfig{
				ServerURL: ts.URL,
				Sampling: ZipkinSampling{
					SamplingRatio: 1,
				},
			},
		},
	})
	assert.NoError(t, err)

	trc := zt.Tracer()
	_, span := trc.Start(context.Background(), "testspan")
	span.SetAttributes(attribute.Bool("testTag", true))
	span.End()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatalf("Test server did not receive spans")
	}
}

func TestOTLPTracer(t *testing.T) {
	done := make(chan struct{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := decodeResponseBody(t, r)

		var res tracepb.ExportTraceServiceRequest
		err := proto.Unmarshal(body, &res)
		require.NoError(t, err, "must be able to unmarshal traces")

		resourceSpans := res.GetResourceSpans()
		spans := resourceSpans[0].GetScopeSpans()[0].GetSpans()
		assert.Equal(t, len(spans), 1)

		assert.NotEmpty(t, spans[0].GetSpanId())
		assert.NotEmpty(t, spans[0].GetTraceId())
		assert.Equal(t, "testSpan", spans[0].GetName())
		assert.Equal(t, "testAttribute", spans[0].Attributes[0].Key)

		close(done)
	}))
	defer ts.Close()

	tsu, err := url.Parse(ts.URL)
	require.NoError(t, err)

	ot, err := New(testTracingComponent, logrusx.New("ory/x", "1"), &Config{
		ServiceName: "ORY X",
		Provider:    "otel",
		Providers: ProvidersConfig{
			OTLP: OTLPConfig{
				ServerURL: tsu.Host,
				Insecure:  true,
				Sampling: OTLPSampling{
					SamplingRatio: 1,
				},
			},
		},
	})
	assert.NoError(t, err)

	trc := ot.Tracer()
	_, span := trc.Start(context.Background(), "testSpan")
	span.SetAttributes(attribute.Bool("testAttribute", true))
	span.End()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatalf("Test server did not receive spans")
	}
}
