// Copyright 2016 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/swarm"
)

func TestInitSwarm(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: `"body"`, status: http.StatusOK}
	client := newTestClient(fakeRT)
	response, err := client.InitSwarm(InitSwarmOptions{})
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	expectedMethod := "POST"
	if req.Method != expectedMethod {
		t.Errorf("InitSwarm: Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
	}
	u, _ := url.Parse(client.getURL("/swarm/init"))
	if req.URL.Path != u.Path {
		t.Errorf("InitSwarm: Wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
	expected := "body"
	if response != expected {
		t.Errorf("InitSwarm: Wrong response. Want %q. Got %q.", expected, response)
	}
}

func TestInitSwarmAlreadyInSwarm(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "", status: http.StatusNotAcceptable})
	_, err := client.InitSwarm(InitSwarmOptions{})
	if err != ErrNodeAlreadyInSwarm {
		t.Errorf("InitSwarm: Wrong error type. Want %#v. Got %#v", ErrNodeAlreadyInSwarm, err)
	}
}

func TestJoinSwarm(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.JoinSwarm(JoinSwarmOptions{})
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	expectedMethod := "POST"
	if req.Method != expectedMethod {
		t.Errorf("JoinSwarm: Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
	}
	u, _ := url.Parse(client.getURL("/swarm/join"))
	if req.URL.Path != u.Path {
		t.Errorf("JoinSwarm: Wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
}

func TestJoinSwarmAlreadyInSwarm(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "", status: http.StatusNotAcceptable})
	err := client.JoinSwarm(JoinSwarmOptions{})
	if err != ErrNodeAlreadyInSwarm {
		t.Errorf("JoinSwarm: Wrong error type. Want %#v. Got %#v", ErrNodeAlreadyInSwarm, err)
	}
}

func TestLeaveSwarm(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	var testData = []struct {
		force       bool
		expectedURI string
	}{
		{false, "/swarm/leave?force=false"},
		{true, "/swarm/leave?force=true"},
	}
	for i, tt := range testData {
		err := client.LeaveSwarm(LeaveSwarmOptions{Force: tt.force})
		if err != nil {
			t.Fatal(err)
		}
		expectedMethod := "POST"
		req := fakeRT.requests[i]
		if req.Method != expectedMethod {
			t.Errorf("LeaveSwarm: Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
		}
		expected, _ := url.Parse(client.getURL(tt.expectedURI))
		if req.URL.String() != expected.String() {
			t.Errorf("LeaveSwarm: Wrong request string. Want %q. Got %q.", expected.String(), req.URL.String())
		}
	}
}

func TestLeaveSwarmNotInSwarm(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "", status: http.StatusNotAcceptable})
	err := client.LeaveSwarm(LeaveSwarmOptions{})
	if err != ErrNodeNotInSwarm {
		t.Errorf("LeaveSwarm: Wrong error type. Want %#v. Got %#v", ErrNodeNotInSwarm, err)
	}
}

func TestUpdateSwarm(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	opts := UpdateSwarmOptions{
		Version:            10,
		RotateManagerToken: true,
		RotateWorkerToken:  false,
	}
	err := client.UpdateSwarm(opts)
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	expectedMethod := "POST"
	if req.Method != expectedMethod {
		t.Errorf("UpdateSwarm: Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
	}
	expectedPath := "/swarm/update"
	if req.URL.Path != expectedPath {
		t.Errorf("UpdateSwarm: Wrong request path. Want %q. Got %q.", expectedPath, req.URL.Path)
	}
	expected := map[string][]string{
		"version":            {"10"},
		"rotateManagerToken": {"true"},
		"rotateWorkerToken":  {"false"},
	}
	got := map[string][]string(req.URL.Query())
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("UpdateSwarm: Wrong request query. Want %v. Got %v", expected, got)
	}
}

func TestUpdateSwarmNotInSwarm(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "", status: http.StatusNotAcceptable})
	err := client.UpdateSwarm(UpdateSwarmOptions{})
	if err != ErrNodeNotInSwarm {
		t.Errorf("UpdateSwarm: Wrong error type. Want %#v. Got %#v", ErrNodeNotInSwarm, err)
	}
}

func TestInspectSwarm(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: `{"ID": "123"}`, status: http.StatusOK}
	client := newTestClient(fakeRT)
	response, err := client.InspectSwarm(nil)
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	expectedMethod := "GET"
	if req.Method != expectedMethod {
		t.Errorf("InspectSwarm: Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
	}
	u, _ := url.Parse(client.getURL("/swarm"))
	if req.URL.Path != u.Path {
		t.Errorf("InspectSwarm: Wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
	expected := swarm.Swarm{ClusterInfo: swarm.ClusterInfo{ID: "123"}}
	if !reflect.DeepEqual(expected, response) {
		t.Errorf("InspectSwarm: Wrong response. Want %#v. Got %#v.", expected, response)
	}
}
