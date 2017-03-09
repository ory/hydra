// Copyright 2013 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type DockerVersion struct {
	Version   string
	GitCommit string
	GoVersion string
}

func TestVersion(t *testing.T) {
	body := `{
     "Version":"0.2.2",
     "GitCommit":"5a2a5cc+CHANGES",
     "GoVersion":"go1.0.3"
}`
	fakeRT := FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(&fakeRT)
	expected := DockerVersion{
		Version:   "0.2.2",
		GitCommit: "5a2a5cc+CHANGES",
		GoVersion: "go1.0.3",
	}
	version, err := client.Version()
	if err != nil {
		t.Fatal(err)
	}

	if result := version.Get("Version"); result != expected.Version {
		t.Errorf("Version(): Wrong result. Want %#v. Got %#v.", expected.Version, version.Get("Version"))
	}
	if result := version.Get("GitCommit"); result != expected.GitCommit {
		t.Errorf("GitCommit(): Wrong result. Want %#v. Got %#v.", expected.GitCommit, version.Get("GitCommit"))
	}
	if result := version.Get("GoVersion"); result != expected.GoVersion {
		t.Errorf("GoVersion(): Wrong result. Want %#v. Got %#v.", expected.GoVersion, version.Get("GoVersion"))
	}
	req := fakeRT.requests[0]
	if req.Method != "GET" {
		t.Errorf("Version(): wrong request method. Want GET. Got %s.", req.Method)
	}
	u, _ := url.Parse(client.getURL("/version"))
	if req.URL.Path != u.Path {
		t.Errorf("Version(): wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
}

func TestVersionError(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "internal error", status: http.StatusInternalServerError}
	client := newTestClient(fakeRT)
	version, err := client.Version()
	if version != nil {
		t.Errorf("Version(): expected <nil> value, got %#v.", version)
	}
	if err == nil {
		t.Error("Version(): unexpected <nil> error")
	}
}

func TestInfo(t *testing.T) {
	body := `{
     "Containers":11,
     "Images":16,
     "Debug":false,
     "NFd":11,
     "NGoroutines":21,
     "MemoryLimit":true,
     "SwapLimit":false,
     "RegistryConfig":{
       "InsecureRegistryCIDRs":["127.0.0.0/8"],
       "IndexConfigs":{
         "docker.io":{
           "Name":"docker.io",
           "Mirrors":null,
           "Secure":true,
           "Official":true
         }
       },
       "Mirrors":null
     }
}`
	fakeRT := FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(&fakeRT)
	expected := &DockerInfo{
		Containers:  11,
		Images:      16,
		Debug:       false,
		NFd:         11,
		NGoroutines: 21,
		MemoryLimit: true,
		SwapLimit:   false,
		RegistryConfig: &ServiceConfig{
			InsecureRegistryCIDRs: []*NetIPNet{
				{
					Mask: net.CIDRMask(8, 32),
					IP:   net.ParseIP("127.0.0.0").To4(),
				},
			},
			IndexConfigs: map[string]*IndexInfo{
				"docker.io": {
					Name:     "docker.io",
					Secure:   true,
					Official: true,
				},
			},
		},
	}
	info, err := client.Info()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, info) {
		t.Errorf("Info(): Wrong result.\nWant %#v.\nGot %#v.", expected, info)
	}
	req := fakeRT.requests[0]
	if req.Method != "GET" {
		t.Errorf("Info(): Wrong HTTP method. Want GET. Got %s.", req.Method)
	}
	u, _ := url.Parse(client.getURL("/info"))
	if req.URL.Path != u.Path {
		t.Errorf("Info(): Wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
}

func TestInfoError(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "internal error", status: http.StatusInternalServerError}
	client := newTestClient(fakeRT)
	version, err := client.Info()
	if version != nil {
		t.Errorf("Info(): expected <nil> value, got %#v.", version)
	}
	if err == nil {
		t.Error("Info(): unexpected <nil> error")
	}
}

func TestParseRepositoryTag(t *testing.T) {
	var tests = []struct {
		input        string
		expectedRepo string
		expectedTag  string
	}{
		{
			"localhost.localdomain:5000/samalba/hipache:latest",
			"localhost.localdomain:5000/samalba/hipache",
			"latest",
		},
		{
			"localhost.localdomain:5000/samalba/hipache",
			"localhost.localdomain:5000/samalba/hipache",
			"",
		},
		{
			"tsuru/python",
			"tsuru/python",
			"",
		},
		{
			"tsuru/python:2.7",
			"tsuru/python",
			"2.7",
		},
		{
			"busybox@sha256:4a731fb46adc5cefe3ae374a8b6020fc1b6ad667a279647766e9a3cd89f6fa92",
			"busybox",
			"",
		},
		{
			"localhost.localdomain:5000/samalba/hipache:v1@sha256:4a731fb46adc5cefe3ae374a8b6020fc1b6ad667a279647766e9a3cd89f6fa92",
			"localhost.localdomain:5000/samalba/hipache",
			"v1",
		},
	}
	for _, tt := range tests {
		repo, tag := ParseRepositoryTag(tt.input)
		if repo != tt.expectedRepo {
			t.Errorf("ParseRepositoryTag(%q): wrong repository. Want %q. Got %q", tt.input, tt.expectedRepo, repo)
		}
		if tag != tt.expectedTag {
			t.Errorf("ParseRepositoryTag(%q): wrong tag. Want %q. Got %q", tt.input, tt.expectedTag, tag)
		}
	}
}
