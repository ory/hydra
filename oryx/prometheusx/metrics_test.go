// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/herodot"
	"github.com/ory/x/logrusx"

	pbTestproto "github.com/grpc-ecosystem/go-grpc-prometheus/examples/testproto"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ioprometheusclient "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	prometheus "github.com/ory/x/prometheusx"
)

const (
	pingDefaultValue   = "I like kittens."
	countListResponses = 20
)

func TestGRPCMetrics(t *testing.T) {
	testApp := "test_app"
	testPath := "/test/path"

	serverListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "must be able to allocate a port for serverListener")
	pmm := prometheus.NewMetricsManager(testApp, "", "", "")
	server := grpc.NewServer(
		grpc.StreamInterceptor(pmm.StreamServerInterceptor),
		grpc.UnaryInterceptor(pmm.UnaryServerInterceptor),
	)
	pbTestproto.RegisterTestServiceServer(server, &testService{t})

	go func() {
		server.Serve(serverListener)
	}()

	clientConn, err := grpc.Dial(serverListener.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	require.NoError(t, err, "must not error on client Dial")
	testClient := pbTestproto.NewTestServiceClient(clientConn)

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)

	pmm.Register(server)

	_, err = testClient.PingEmpty(ctx, &pbTestproto.Empty{})
	require.NoError(t, err)
	_, err = testClient.PingList(ctx, &pbTestproto.PingRequest{})
	require.NoError(t, err)

	n := negroni.New()

	router := httprouter.New()

	pmm.RegisterRouter(router)
	prometheus.NewHandler(herodot.NewJSONWriter(logrusx.New("Ory X", "test")), "test").SetRoutes(router)

	router.GET(testPath, func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		rw.WriteHeader(http.StatusBadRequest)
	})

	n.UseHandler(router)
	n.Use(pmm)

	ts := httptest.NewServer(n)
	defer ts.Close()

	resp, err := http.Get(ts.URL + testPath)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

	promresp, err := http.Get(ts.URL + prometheus.MetricsPrometheusPath)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, promresp.StatusCode)

	textParser := expfmt.TextParser{}
	text, err := textParser.TextToMetricFamilies(promresp.Body)
	require.NoError(t, err)

	require.EqualValues(t, "grpc_server_handled_total", *text["grpc_server_handled_total"].Name)
	require.EqualValues(t, "Ping", getLabelValue("grpc_method", text["grpc_server_handled_total"].Metric))
	require.EqualValues(t, "mwitkow.testproto.TestService", getLabelValue("grpc_service", text["grpc_server_handled_total"].Metric))
	c, err := GetCounterValue(text["grpc_server_handled_total"].Metric, "PingEmpty", "OK")
	require.NoError(t, err)
	require.EqualValues(t, 1, c)
	c, err = GetCounterValue(text["grpc_server_handled_total"].Metric, "PingList", "OK")
	require.NoError(t, err)
	require.EqualValues(t, 1, c)

	require.EqualValues(t, "grpc_server_msg_sent_total", *text["grpc_server_msg_sent_total"].Name)
	require.EqualValues(t, "Ping", getLabelValue("grpc_method", text["grpc_server_msg_sent_total"].Metric))
	require.EqualValues(t, "mwitkow.testproto.TestService", getLabelValue("grpc_service", text["grpc_server_msg_sent_total"].Metric))

	require.EqualValues(t, "grpc_server_msg_received_total", *text["grpc_server_msg_received_total"].Name)
	require.EqualValues(t, "Ping", getLabelValue("grpc_method", text["grpc_server_msg_received_total"].Metric))
	require.EqualValues(t, "mwitkow.testproto.TestService", getLabelValue("grpc_service", text["grpc_server_msg_received_total"].Metric))

	cancel()
	server.Stop()
	serverListener.Close()
}

func TestHTTPMetrics(t *testing.T) {
	testApp := "test_app"
	testPath := "/test/path"

	n := negroni.New()
	handler := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		prometheus.NewMetrics(testApp, prometheus.HTTPMetrics, "", "", "").Instrument(rw, next, r.RequestURI)(rw, r)
	}
	n.UseFunc(handler)

	router := httprouter.New()
	router.GET(testPath, func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		rw.WriteHeader(http.StatusBadRequest)
	})
	router.GET(prometheus.MetricsPrometheusPath, func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		promhttp.Handler().ServeHTTP(rw, r)
	})
	n.UseHandler(router)

	ts := httptest.NewServer(n)
	defer ts.Close()

	resp, err := http.Get(ts.URL + testPath)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

	promresp, err := http.Get(ts.URL + prometheus.MetricsPrometheusPath)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, promresp.StatusCode)

	textParser := expfmt.TextParser{}
	text, err := textParser.TextToMetricFamilies(promresp.Body)
	require.NoError(t, err)
	require.EqualValues(t, "http_response_time_seconds", *text["http_response_time_seconds"].Name)
	require.EqualValues(t, testPath, getLabelValue("endpoint", text["http_response_time_seconds"].Metric))
	require.EqualValues(t, testApp, getLabelValue("app", text["http_response_time_seconds"].Metric))

	require.EqualValues(t, "http_requests_total", *text["http_requests_total"].Name)
	require.EqualValues(t, "400", getLabelValue("code", text["http_requests_total"].Metric))
	require.EqualValues(t, testPath, getLabelValue("endpoint", text["http_requests_total"].Metric))
	require.EqualValues(t, testApp, getLabelValue("app", text["http_requests_total"].Metric))

	require.EqualValues(t, "http_requests_duration_seconds", *text["http_requests_duration_seconds"].Name)
	require.EqualValues(t, "400", getLabelValue("code", text["http_requests_duration_seconds"].Metric))
	require.EqualValues(t, testPath, getLabelValue("endpoint", text["http_requests_duration_seconds"].Metric))
	require.EqualValues(t, testApp, getLabelValue("app", text["http_requests_duration_seconds"].Metric))

	require.EqualValues(t, "http_response_size_bytes", *text["http_response_size_bytes"].Name)
	require.EqualValues(t, "400", getLabelValue("code", text["http_response_size_bytes"].Metric))
	require.EqualValues(t, testApp, getLabelValue("app", text["http_response_size_bytes"].Metric))

	require.EqualValues(t, "http_requests_size_bytes", *text["http_requests_size_bytes"].Name)
	require.EqualValues(t, "400", getLabelValue("code", text["http_requests_size_bytes"].Metric))
	require.EqualValues(t, testApp, getLabelValue("app", text["http_requests_size_bytes"].Metric))

	require.EqualValues(t, "http_requests_statuses_total", *text["http_requests_statuses_total"].Name)
	require.EqualValues(t, "4xx", getLabelValue("status_bucket", text["http_requests_statuses_total"].Metric))
	require.EqualValues(t, testApp, getLabelValue("app", text["http_requests_statuses_total"].Metric))
}

func getLabelValue(name string, metric []*ioprometheusclient.Metric) string {
	for _, label := range metric[0].Label {
		if *label.Name == name {
			return *label.Value
		}
	}

	return ""
}

func GetCounterValue(metrics []*ioprometheusclient.Metric, lvs ...string) (float64, error) {
	for _, metric := range metrics {
		lvl := len(lvs)
		lvc := 0
		for _, label := range metric.Label {
			for _, lv := range lvs {
				if lv == *label.Value {
					lvc++
				}
			}
		}
		if lvc == lvl {
			return *metric.Counter.Value, nil
		}
	}
	return 0, errors.New("Counter value was not found")
}

type testService struct {
	t *testing.T
}

func (s *testService) PingEmpty(ctx context.Context, _ *pbTestproto.Empty) (*pbTestproto.PingResponse, error) {
	return &pbTestproto.PingResponse{Value: pingDefaultValue, Counter: 42}, nil
}

func (s *testService) Ping(ctx context.Context, ping *pbTestproto.PingRequest) (*pbTestproto.PingResponse, error) {
	// Send user trailers and headers.
	return &pbTestproto.PingResponse{Value: ping.Value, Counter: 42}, nil
}

func (s *testService) PingError(ctx context.Context, ping *pbTestproto.PingRequest) (*pbTestproto.Empty, error) {
	code := codes.Code(ping.ErrorCodeReturned)
	return nil, status.Errorf(code, "Userspace error.")
}

func (s *testService) PingList(ping *pbTestproto.PingRequest, stream pbTestproto.TestService_PingListServer) error {
	if ping.ErrorCodeReturned != 0 {
		return status.Errorf(codes.Code(ping.ErrorCodeReturned), "foobar")
	}
	// Send user trailers and headers.
	for i := 0; i < countListResponses; i++ {
		stream.Send(&pbTestproto.PingResponse{Value: ping.Value, Counter: int32(i)})
	}
	return nil
}
