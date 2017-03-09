// Copyright 2013 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestNewAPIClient(t *testing.T) {
	endpoint := "http://localhost:4243"
	client, err := NewClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	// test native endpoints
	endpoint = nativeRealEndpoint
	client, err = NewClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if !client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be true, got false")
	}
	if client.requestedAPIVersion != nil {
		t.Errorf("Expected requestedAPIVersion to be nil, got %#v.", client.requestedAPIVersion)
	}
}

func newTLSClient(endpoint string) (*Client, error) {
	return NewTLSClient(endpoint,
		"testing/data/cert.pem",
		"testing/data/key.pem",
		"testing/data/ca.pem")
}

func TestNewTSLAPIClient(t *testing.T) {
	endpoint := "https://localhost:4243"
	client, err := newTLSClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if !client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be true, got false")
	}
	if client.requestedAPIVersion != nil {
		t.Errorf("Expected requestedAPIVersion to be nil, got %#v.", client.requestedAPIVersion)
	}
}

func TestNewVersionedClient(t *testing.T) {
	endpoint := "http://localhost:4243"
	client, err := NewVersionedClient(endpoint, "1.12")
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if reqVersion := client.requestedAPIVersion.String(); reqVersion != "1.12" {
		t.Errorf("Wrong requestAPIVersion. Want %q. Got %q.", "1.12", reqVersion)
	}
	if client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be false, got true")
	}
}

func TestNewVersionedClientFromEnv(t *testing.T) {
	endpoint := "tcp://localhost:2376"
	endpointURL := "http://localhost:2376"
	os.Setenv("DOCKER_HOST", endpoint)
	os.Setenv("DOCKER_TLS_VERIFY", "")
	client, err := NewVersionedClientFromEnv("1.12")
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if client.endpointURL.String() != endpointURL {
		t.Errorf("Expected endpointURL %s. Got %s.", endpoint, client.endpoint)
	}
	if reqVersion := client.requestedAPIVersion.String(); reqVersion != "1.12" {
		t.Errorf("Wrong requestAPIVersion. Want %q. Got %q.", "1.12", reqVersion)
	}
	if client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be false, got true")
	}
}

func TestNewVersionedClientFromEnvTLS(t *testing.T) {
	endpoint := "tcp://localhost:2376"
	endpointURL := "https://localhost:2376"
	base, _ := os.Getwd()
	os.Setenv("DOCKER_CERT_PATH", filepath.Join(base, "/testing/data/"))
	os.Setenv("DOCKER_HOST", endpoint)
	os.Setenv("DOCKER_TLS_VERIFY", "1")
	client, err := NewVersionedClientFromEnv("1.12")
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if client.endpointURL.String() != endpointURL {
		t.Errorf("Expected endpointURL %s. Got %s.", endpoint, client.endpoint)
	}
	if reqVersion := client.requestedAPIVersion.String(); reqVersion != "1.12" {
		t.Errorf("Wrong requestAPIVersion. Want %q. Got %q.", "1.12", reqVersion)
	}
	if client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be false, got true")
	}
}

func TestNewTLSVersionedClient(t *testing.T) {
	certPath := "testing/data/cert.pem"
	keyPath := "testing/data/key.pem"
	caPath := "testing/data/ca.pem"
	endpoint := "https://localhost:4243"
	client, err := NewVersionedTLSClient(endpoint, certPath, keyPath, caPath, "1.14")
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if reqVersion := client.requestedAPIVersion.String(); reqVersion != "1.14" {
		t.Errorf("Wrong requestAPIVersion. Want %q. Got %q.", "1.14", reqVersion)
	}
	if client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be false, got true")
	}
}

func TestNewTLSVersionedClientNoClientCert(t *testing.T) {
	certPath := "testing/data/cert_doesnotexist.pem"
	keyPath := "testing/data/key_doesnotexist.pem"
	caPath := "testing/data/ca.pem"
	endpoint := "https://localhost:4243"
	client, err := NewVersionedTLSClient(endpoint, certPath, keyPath, caPath, "1.14")
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
	if reqVersion := client.requestedAPIVersion.String(); reqVersion != "1.14" {
		t.Errorf("Wrong requestAPIVersion. Want %q. Got %q.", "1.14", reqVersion)
	}
	if client.SkipServerVersionCheck {
		t.Error("Expected SkipServerVersionCheck to be false, got true")
	}
}

func TestNewTLSVersionedClientInvalidCA(t *testing.T) {
	certPath := "testing/data/cert.pem"
	keyPath := "testing/data/key.pem"
	caPath := "testing/data/key.pem"
	endpoint := "https://localhost:4243"
	_, err := NewVersionedTLSClient(endpoint, certPath, keyPath, caPath, "1.14")
	if err == nil {
		t.Errorf("Expected invalid ca at %s", caPath)
	}
}

func TestNewTLSVersionedClientInvalidCANoClientCert(t *testing.T) {
	certPath := "testing/data/cert_doesnotexist.pem"
	keyPath := "testing/data/key_doesnotexist.pem"
	caPath := "testing/data/key.pem"
	endpoint := "https://localhost:4243"
	_, err := NewVersionedTLSClient(endpoint, certPath, keyPath, caPath, "1.14")
	if err == nil {
		t.Errorf("Expected invalid ca at %s", caPath)
	}
}

func TestNewClientInvalidEndpoint(t *testing.T) {
	cases := []string{
		"htp://localhost:3243", "http://localhost:a",
		"", "http://localhost:8080:8383", "http://localhost:65536",
		"https://localhost:-20",
	}
	for _, c := range cases {
		client, err := NewClient(c)
		if client != nil {
			t.Errorf("Want <nil> client for invalid endpoint, got %#v.", client)
		}
		if !reflect.DeepEqual(err, ErrInvalidEndpoint) {
			t.Errorf("NewClient(%q): Got invalid error for invalid endpoint. Want %#v. Got %#v.", c, ErrInvalidEndpoint, err)
		}
	}
}

func TestNewClientNoSchemeEndpoint(t *testing.T) {
	cases := []string{"localhost", "localhost:8080"}
	for _, c := range cases {
		client, err := NewClient(c)
		if client == nil {
			t.Errorf("Want client for scheme-less endpoint, got <nil>")
		}
		if err != nil {
			t.Errorf("Got unexpected error scheme-less endpoint: %q", err)
		}
	}
}

func TestNewTLSClient(t *testing.T) {
	var tests = []struct {
		endpoint string
		expected string
	}{
		{"tcp://localhost:2376", "https"},
		{"tcp://localhost:2375", "https"},
		{"tcp://localhost:4000", "https"},
		{"http://localhost:4000", "https"},
	}
	for _, tt := range tests {
		client, err := newTLSClient(tt.endpoint)
		if err != nil {
			t.Error(err)
		}
		got := client.endpointURL.Scheme
		if got != tt.expected {
			t.Errorf("endpointURL.Scheme: Got %s. Want %s.", got, tt.expected)
		}
	}
}

func TestEndpoint(t *testing.T) {
	client, err := NewVersionedClient("http://localhost:4243", "1.12")
	if err != nil {
		t.Fatal(err)
	}
	if endpoint := client.Endpoint(); endpoint != client.endpoint {
		t.Errorf("Client.Endpoint(): want %q. Got %q", client.endpoint, endpoint)
	}
}

func TestGetURL(t *testing.T) {
	var tests = []struct {
		endpoint string
		path     string
		expected string
	}{
		{"http://localhost:4243/", "/", "http://localhost:4243/"},
		{"http://localhost:4243", "/", "http://localhost:4243/"},
		{"http://localhost:4243", "/containers/ps", "http://localhost:4243/containers/ps"},
		{"tcp://localhost:4243", "/containers/ps", "http://localhost:4243/containers/ps"},
		{"http://localhost:4243/////", "/", "http://localhost:4243/"},
		{nativeRealEndpoint, "/containers", "/containers"},
	}
	for _, tt := range tests {
		client, _ := NewClient(tt.endpoint)
		client.endpoint = tt.endpoint
		client.SkipServerVersionCheck = true
		got := client.getURL(tt.path)
		if got != tt.expected {
			t.Errorf("getURL(%q): Got %s. Want %s.", tt.path, got, tt.expected)
		}
	}
}

func TestGetFakeNativeURL(t *testing.T) {
	var tests = []struct {
		endpoint string
		path     string
		expected string
	}{
		{nativeRealEndpoint, "/", "http://unix.sock/"},
		{nativeRealEndpoint, "/", "http://unix.sock/"},
		{nativeRealEndpoint, "/containers/ps", "http://unix.sock/containers/ps"},
	}
	for _, tt := range tests {
		client, _ := NewClient(tt.endpoint)
		client.endpoint = tt.endpoint
		client.SkipServerVersionCheck = true
		got := client.getFakeNativeURL(tt.path)
		if got != tt.expected {
			t.Errorf("getURL(%q): Got %s. Want %s.", tt.path, got, tt.expected)
		}
	}
}

func TestError(t *testing.T) {
	fakeBody := ioutil.NopCloser(bytes.NewBufferString("bad parameter"))
	resp := &http.Response{
		StatusCode: 400,
		Body:       fakeBody,
	}
	err := newError(resp)
	expected := Error{Status: 400, Message: "bad parameter"}
	if !reflect.DeepEqual(expected, *err) {
		t.Errorf("Wrong error type. Want %#v. Got %#v.", expected, *err)
	}
	message := "API error (400): bad parameter"
	if err.Error() != message {
		t.Errorf("Wrong error message. Want %q. Got %q.", message, err.Error())
	}
}

func TestQueryString(t *testing.T) {
	v := float32(2.4)
	f32QueryString := fmt.Sprintf("w=%s&x=10&y=10.35", strconv.FormatFloat(float64(v), 'f', -1, 64))
	jsonPerson := url.QueryEscape(`{"Name":"gopher","age":4}`)
	var tests = []struct {
		input interface{}
		want  string
	}{
		{&ListContainersOptions{All: true}, "all=1"},
		{ListContainersOptions{All: true}, "all=1"},
		{ListContainersOptions{Before: "something"}, "before=something"},
		{ListContainersOptions{Before: "something", Since: "other"}, "before=something&since=other"},
		{ListContainersOptions{Filters: map[string][]string{"status": {"paused", "running"}}}, "filters=%7B%22status%22%3A%5B%22paused%22%2C%22running%22%5D%7D"},
		{dumb{X: 10, Y: 10.35000}, "x=10&y=10.35"},
		{dumb{W: v, X: 10, Y: 10.35000}, f32QueryString},
		{dumb{X: 10, Y: 10.35000, Z: 10}, "x=10&y=10.35&zee=10"},
		{dumb{v: 4, X: 10, Y: 10.35000}, "x=10&y=10.35"},
		{dumb{T: 10, Y: 10.35000}, "y=10.35"},
		{dumb{Person: &person{Name: "gopher", Age: 4}}, "p=" + jsonPerson},
		{nil, ""},
		{10, ""},
		{"not_a_struct", ""},
	}
	for _, tt := range tests {
		got := queryString(tt.input)
		if got != tt.want {
			t.Errorf("queryString(%v). Want %q. Got %q.", tt.input, tt.want, got)
		}
	}
}

func TestAPIVersions(t *testing.T) {
	var tests = []struct {
		a                              string
		b                              string
		expectedALessThanB             bool
		expectedALessThanOrEqualToB    bool
		expectedAGreaterThanB          bool
		expectedAGreaterThanOrEqualToB bool
	}{
		{"1.11", "1.11", false, true, false, true},
		{"1.10", "1.11", true, true, false, false},
		{"1.11", "1.10", false, false, true, true},

		{"1.11-ubuntu0", "1.11", false, true, false, true},
		{"1.10", "1.11-el7", true, true, false, false},

		{"1.9", "1.11", true, true, false, false},
		{"1.11", "1.9", false, false, true, true},

		{"1.1.1", "1.1", false, false, true, true},
		{"1.1", "1.1.1", true, true, false, false},

		{"2.1", "1.1.1", false, false, true, true},
		{"2.1", "1.3.1", false, false, true, true},
		{"1.1.1", "2.1", true, true, false, false},
		{"1.3.1", "2.1", true, true, false, false},
	}

	for _, tt := range tests {
		a, _ := NewAPIVersion(tt.a)
		b, _ := NewAPIVersion(tt.b)

		if tt.expectedALessThanB && !a.LessThan(b) {
			t.Errorf("Expected %#v < %#v", a, b)
		}
		if tt.expectedALessThanOrEqualToB && !a.LessThanOrEqualTo(b) {
			t.Errorf("Expected %#v <= %#v", a, b)
		}
		if tt.expectedAGreaterThanB && !a.GreaterThan(b) {
			t.Errorf("Expected %#v > %#v", a, b)
		}
		if tt.expectedAGreaterThanOrEqualToB && !a.GreaterThanOrEqualTo(b) {
			t.Errorf("Expected %#v >= %#v", a, b)
		}
	}
}

func TestPing(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPingFailing(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusInternalServerError}
	client := newTestClient(fakeRT)
	err := client.Ping()
	if err == nil {
		t.Fatal("Expected non nil error, got nil")
	}
	expectedErrMsg := "API error (500): "
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error to be %q, got: %q", expectedErrMsg, err.Error())
	}
}

func TestPingFailingWrongStatus(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusAccepted}
	client := newTestClient(fakeRT)
	err := client.Ping()
	if err == nil {
		t.Fatal("Expected non nil error, got nil")
	}
	expectedErrMsg := "API error (202): "
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error to be %q, got: %q", expectedErrMsg, err.Error())
	}
}

func TestPingErrorWithNativeClient(t *testing.T) {
	srv, cleanup, err := newNativeServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("aaaaaaaaaaa-invalid-aaaaaaaaaaa"))
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	srv.Start()
	defer srv.Close()
	endpoint := nativeBadEndpoint
	client, err := NewClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	err = client.Ping()
	if err == nil {
		t.Fatal("Expected non nil error, got nil")
	}
}

func TestClientStreamTimeoutNotHit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 5; i++ {
			fmt.Fprintf(w, "%d\n", i)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var w bytes.Buffer
	err = client.stream("POST", "/image/create", streamOptions{
		setRawTerminal:    true,
		stdout:            &w,
		inactivityTimeout: 300 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "0\n1\n2\n3\n4\n"
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

func TestClientStreamInactivityTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 5; i++ {
			fmt.Fprintf(w, "%d\n", i)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var w bytes.Buffer
	err = client.stream("POST", "/image/create", streamOptions{
		setRawTerminal:    true,
		stdout:            &w,
		inactivityTimeout: 100 * time.Millisecond,
	})
	if err != ErrInactivityTimeout {
		t.Fatalf("expected request canceled error, got: %s", err)
	}
	expected := "0\n"
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

func TestClientStreamContextDeadline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "abc\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(w, "def\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var w bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err = client.stream("POST", "/image/create", streamOptions{
		setRawTerminal: true,
		stdout:         &w,
		context:        ctx,
	})
	if err != context.DeadlineExceeded {
		t.Fatalf("expected %s, got: %s", context.DeadlineExceeded, err)
	}
	expected := "abc\n"
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

func TestClientStreamContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "abc\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(w, "def\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var w bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()
	err = client.stream("POST", "/image/create", streamOptions{
		setRawTerminal: true,
		stdout:         &w,
		context:        ctx,
	})
	if err != context.Canceled {
		t.Fatalf("expected %s, got: %s", context.Canceled, err)
	}
	expected := "abc\n"
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

var mockPullOutput = `{"status":"Pulling from tsuru/static","id":"latest"}
{"status":"Already exists","progressDetail":{},"id":"a6aa3b66376f"}
{"status":"Pulling fs layer","progressDetail":{},"id":"106572778bf7"}
{"status":"Pulling fs layer","progressDetail":{},"id":"bac681833e51"}
{"status":"Pulling fs layer","progressDetail":{},"id":"7302e23ef08a"}
{"status":"Downloading","progressDetail":{"current":621,"total":621},"progress":"[==================================================\u003e]    621 B/621 B","id":"bac681833e51"}
{"status":"Verifying Checksum","progressDetail":{},"id":"bac681833e51"}
{"status":"Download complete","progressDetail":{},"id":"bac681833e51"}
{"status":"Downloading","progressDetail":{"current":1854,"total":1854},"progress":"[==================================================\u003e] 1.854 kB/1.854 kB","id":"106572778bf7"}
{"status":"Verifying Checksum","progressDetail":{},"id":"106572778bf7"}
{"status":"Download complete","progressDetail":{},"id":"106572778bf7"}
{"status":"Extracting","progressDetail":{"current":1854,"total":1854},"progress":"[==================================================\u003e] 1.854 kB/1.854 kB","id":"106572778bf7"}
{"status":"Extracting","progressDetail":{"current":1854,"total":1854},"progress":"[==================================================\u003e] 1.854 kB/1.854 kB","id":"106572778bf7"}
{"status":"Downloading","progressDetail":{"current":233019,"total":21059403},"progress":"[\u003e                                                  ]   233 kB/21.06 MB","id":"7302e23ef08a"}
{"status":"Downloading","progressDetail":{"current":462395,"total":21059403},"progress":"[=\u003e                                                 ] 462.4 kB/21.06 MB","id":"7302e23ef08a"}
{"status":"Downloading","progressDetail":{"current":8490555,"total":21059403},"progress":"[====================\u003e                              ] 8.491 MB/21.06 MB","id":"7302e23ef08a"}
{"status":"Downloading","progressDetail":{"current":20876859,"total":21059403},"progress":"[=================================================\u003e ] 20.88 MB/21.06 MB","id":"7302e23ef08a"}
{"status":"Verifying Checksum","progressDetail":{},"id":"7302e23ef08a"}
{"status":"Download complete","progressDetail":{},"id":"7302e23ef08a"}
{"status":"Pull complete","progressDetail":{},"id":"106572778bf7"}
{"status":"Extracting","progressDetail":{"current":621,"total":621},"progress":"[==================================================\u003e]    621 B/621 B","id":"bac681833e51"}
{"status":"Extracting","progressDetail":{"current":621,"total":621},"progress":"[==================================================\u003e]    621 B/621 B","id":"bac681833e51"}
{"status":"Pull complete","progressDetail":{},"id":"bac681833e51"}
{"status":"Extracting","progressDetail":{"current":229376,"total":21059403},"progress":"[\u003e                                                  ] 229.4 kB/21.06 MB","id":"7302e23ef08a"}
{"status":"Extracting","progressDetail":{"current":458752,"total":21059403},"progress":"[=\u003e                                                 ] 458.8 kB/21.06 MB","id":"7302e23ef08a"}
{"status":"Extracting","progressDetail":{"current":11239424,"total":21059403},"progress":"[==========================\u003e                        ] 11.24 MB/21.06 MB","id":"7302e23ef08a"}
{"status":"Extracting","progressDetail":{"current":21059403,"total":21059403},"progress":"[==================================================\u003e] 21.06 MB/21.06 MB","id":"7302e23ef08a"}
{"status":"Pull complete","progressDetail":{},"id":"7302e23ef08a"}
{"status":"Digest: sha256:b754472891aa7e33fc0214e3efa988174f2c2289285fcae868b7ec8b6675fc77"}
{"status":"Status: Downloaded newer image for 192.168.50.4:5000/tsuru/static"}
`

func TestClientStreamJSONDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockPullOutput))
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var w bytes.Buffer
	err = client.stream("POST", "/image/create", streamOptions{
		stdout:         &w,
		useJSONDecoder: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := `latest: Pulling from tsuru/static
a6aa3b66376f: Already exists
106572778bf7: Pulling fs layer
bac681833e51: Pulling fs layer
7302e23ef08a: Pulling fs layer
bac681833e51: Verifying Checksum
bac681833e51: Download complete
106572778bf7: Verifying Checksum
106572778bf7: Download complete
7302e23ef08a: Verifying Checksum
7302e23ef08a: Download complete
106572778bf7: Pull complete
bac681833e51: Pull complete
7302e23ef08a: Pull complete
Digest: sha256:b754472891aa7e33fc0214e3efa988174f2c2289285fcae868b7ec8b6675fc77
Status: Downloaded newer image for 192.168.50.4:5000/tsuru/static
`
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

type terminalBuffer struct {
	bytes.Buffer
}

func (b *terminalBuffer) FD() uintptr {
	return os.Stdout.Fd()
}

func (b *terminalBuffer) IsTerminal() bool {
	return true
}

func TestClientStreamJSONDecodeWithTerminal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockPullOutput))
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var w terminalBuffer
	err = client.stream("POST", "/image/create", streamOptions{
		stdout:         &w,
		useJSONDecoder: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "latest: Pulling from tsuru/static\n\n" +
		"\x1b[1A\x1b[1K\x1b[K\ra6aa3b66376f: Already exists \r\x1b[1B\n" +
		"\x1b[1A\x1b[1K\x1b[K\r106572778bf7: Pulling fs layer \r\x1b[1B\n" +
		"\x1b[1A\x1b[1K\x1b[K\rbac681833e51: Pulling fs layer \r\x1b[1B\n" +
		"\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Pulling fs layer \r\x1b[1B\x1b[2A\x1b[1K\x1b[K\rbac681833e51: Downloading [==================================================>]     621B/621B\r\x1b[2B\x1b[2A\x1b[1K\x1b[K\rbac681833e51: Verifying Checksum \r\x1b[2B\x1b[2A\x1b[1K\x1b[K\rbac681833e51: Download complete \r\x1b[2B\x1b[3A\x1b[1K\x1b[K\r106572778bf7: Downloading [==================================================>]  1.854kB/1.854kB\r\x1b[3B\x1b[3A\x1b[1K\x1b[K\r106572778bf7: Verifying Checksum \r\x1b[3B\x1b[3A\x1b[1K\x1b[K\r106572778bf7: Download complete \r\x1b[3B\x1b[3A\x1b[1K\x1b[K\r106572778bf7: Extracting [==================================================>]  1.854kB/1.854kB\r\x1b[3B\x1b[3A\x1b[1K\x1b[K\r106572778bf7: Extracting [==================================================>]  1.854kB/1.854kB\r\x1b[3B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Downloading [>                                                  ]    233kB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Downloading [=>                                                 ]  462.4kB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Downloading [====================>                              ]  8.491MB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Downloading [=================================================> ]  20.88MB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Verifying Checksum \r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Download complete \r\x1b[1B\x1b[3A\x1b[1K\x1b[K\r106572778bf7: Pull complete \r\x1b[3B\x1b[2A\x1b[1K\x1b[K\rbac681833e51: Extracting [==================================================>]     621B/621B\r\x1b[2B\x1b[2A\x1b[1K\x1b[K\rbac681833e51: Extracting [==================================================>]     621B/621B\r\x1b[2B\x1b[2A\x1b[1K\x1b[K\rbac681833e51: Pull complete \r\x1b[2B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Extracting [>                                                  ]  229.4kB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Extracting [=>                                                 ]  458.8kB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Extracting [==========================>                        ]  11.24MB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Extracting [==================================================>]  21.06MB/21.06MB\r\x1b[1B\x1b[1A\x1b[1K\x1b[K\r7302e23ef08a: Pull complete \r\x1b[1BDigest: sha256:b754472891aa7e33fc0214e3efa988174f2c2289285fcae868b7ec8b6675fc77\n" +
		"Status: Downloaded newer image for 192.168.50.4:5000/tsuru/static\n"
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

func TestClientDoContextDeadline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err = client.do("POST", "/image/create", doOptions{
		context: ctx,
	})
	if err != context.DeadlineExceeded {
		t.Fatalf("expected %s, got: %s", context.DeadlineExceeded, err)
	}
}

func TestClientDoContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	_, err = client.do("POST", "/image/create", doOptions{
		context: ctx,
	})
	if err != context.Canceled {
		t.Fatalf("expected %s, got: %s", context.Canceled, err)
	}
}

func TestClientStreamTimeoutNativeClient(t *testing.T) {
	srv, cleanup, err := newNativeServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 5; i++ {
			fmt.Fprintf(w, "%d\n", i)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	srv.Start()
	defer srv.Close()
	client, err := NewClient(nativeProtocol + "://" + srv.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	var w bytes.Buffer
	err = client.stream("POST", "/image/create", streamOptions{
		setRawTerminal:    true,
		stdout:            &w,
		inactivityTimeout: 100 * time.Millisecond,
	})
	if err != ErrInactivityTimeout {
		t.Fatalf("expected request canceled error, got: %s", err)
	}
	expected := "0\n"
	result := w.String()
	if result != expected {
		t.Fatalf("expected stream result %q, got: %q", expected, result)
	}
}

func TestClientDoConcurrentStress(t *testing.T) {
	var reqs []*http.Request
	var mu sync.Mutex
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		reqs = append(reqs, r)
		mu.Unlock()
	})
	var nativeSrvs []*httptest.Server
	for i := 0; i < 3; i++ {
		srv, cleanup, err := newNativeServer(handler)
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup()
		nativeSrvs = append(nativeSrvs, srv)
	}
	var tests = []struct {
		testCase      string
		srv           *httptest.Server
		scheme        string
		withTimeout   bool
		withTLSServer bool
		withTLSClient bool
	}{
		{testCase: "http server", srv: httptest.NewUnstartedServer(handler), scheme: "http"},
		{testCase: "native server", srv: nativeSrvs[0], scheme: nativeProtocol},
		{testCase: "http with timeout", srv: httptest.NewUnstartedServer(handler), scheme: "http", withTimeout: true},
		{testCase: "native with timeout", srv: nativeSrvs[1], scheme: nativeProtocol, withTimeout: true},
		{testCase: "http with tls", srv: httptest.NewUnstartedServer(handler), scheme: "https", withTLSServer: true, withTLSClient: true},
		{testCase: "native with client-only tls", srv: nativeSrvs[2], scheme: nativeProtocol, withTLSServer: false, withTLSClient: nativeProtocol == unixProtocol}, // TLS client only works with unix protocol
	}
	for _, tt := range tests {
		t.Run(tt.testCase, func(t *testing.T) {
			reqs = nil
			var client *Client
			var err error
			endpoint := tt.scheme + "://" + tt.srv.Listener.Addr().String()
			if tt.withTLSServer {
				tt.srv.StartTLS()
			} else {
				tt.srv.Start()
			}
			defer tt.srv.Close()
			if tt.withTLSClient {
				certPEMBlock, certErr := ioutil.ReadFile("testing/data/cert.pem")
				if certErr != nil {
					t.Fatal(certErr)
				}
				keyPEMBlock, certErr := ioutil.ReadFile("testing/data/key.pem")
				if certErr != nil {
					t.Fatal(certErr)
				}
				client, err = NewTLSClientFromBytes(endpoint, certPEMBlock, keyPEMBlock, nil)
			} else {
				client, err = NewClient(endpoint)
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.withTimeout {
				client.SetTimeout(time.Minute)
			}
			n := 50
			wg := sync.WaitGroup{}
			var paths []string
			errsCh := make(chan error, 3*n)
			waiters := make(chan CloseWaiter, n)
			for i := 0; i < n; i++ {
				path := fmt.Sprintf("/%05d", i)
				paths = append(paths, "GET"+path)
				paths = append(paths, "POST"+path)
				paths = append(paths, "HEAD"+path)
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, clientErr := client.do("GET", path, doOptions{})
					if clientErr != nil {
						errsCh <- clientErr
					}
					clientErr = client.stream("POST", path, streamOptions{})
					if clientErr != nil {
						errsCh <- clientErr
					}
					cw, clientErr := client.hijack("HEAD", path, hijackOptions{})
					if clientErr != nil {
						errsCh <- clientErr
					} else {
						waiters <- cw
					}
				}()
			}
			wg.Wait()
			close(errsCh)
			close(waiters)
			for cw := range waiters {
				cw.Wait()
				cw.Close()
			}
			for err = range errsCh {
				t.Error(err)
			}
			var reqPaths []string
			for _, r := range reqs {
				reqPaths = append(reqPaths, r.Method+r.URL.Path)
			}
			sort.Strings(paths)
			sort.Strings(reqPaths)
			if !reflect.DeepEqual(reqPaths, paths) {
				t.Fatalf("expected server request paths to equal %v, got: %v", paths, reqPaths)
			}
		})
	}
}

type FakeRoundTripper struct {
	message  string
	status   int
	header   map[string]string
	requests []*http.Request
}

func (rt *FakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body := strings.NewReader(rt.message)
	rt.requests = append(rt.requests, r)
	res := &http.Response{
		StatusCode: rt.status,
		Body:       ioutil.NopCloser(body),
		Header:     make(http.Header),
	}
	for k, v := range rt.header {
		res.Header.Set(k, v)
	}
	return res, nil
}

func (rt *FakeRoundTripper) Reset() {
	rt.requests = nil
}

type person struct {
	Name string
	Age  int `json:"age"`
}

type dumb struct {
	T      int `qs:"-"`
	v      int
	W      float32
	X      int
	Y      float64
	Z      int     `qs:"zee"`
	Person *person `qs:"p"`
}
