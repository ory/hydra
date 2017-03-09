// Copyright 2016 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/swarm"
)

func TestCreateService(t *testing.T) {
	result := `{
  "Id": "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2"
}`
	var expected swarm.Service
	err := json.Unmarshal([]byte(result), &expected)
	if err != nil {
		t.Fatal(err)
	}
	fakeRT := &FakeRoundTripper{message: result, status: http.StatusOK}
	client := newTestClient(fakeRT)
	opts := CreateServiceOptions{}
	service, err := client.CreateService(opts)
	if err != nil {
		t.Fatal(err)
	}
	id := "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2"
	if service.ID != id {
		t.Errorf("CreateServce: wrong ID. Want %q. Got %q.", id, service.ID)
	}
	req := fakeRT.requests[0]
	if req.Method != "POST" {
		t.Errorf("CreateService: wrong HTTP method. Want %q. Got %q.", "POST", req.Method)
	}
	expectedURL, _ := url.Parse(client.getURL("/services/create"))
	if gotPath := req.URL.Path; gotPath != expectedURL.Path {
		t.Errorf("CreateServices: Wrong path in request. Want %q. Got %q.", expectedURL.Path, gotPath)
	}
	var gotBody Config
	err = json.NewDecoder(req.Body).Decode(&gotBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveService(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2"
	opts := RemoveServiceOptions{ID: id}
	err := client.RemoveService(opts)
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	if req.Method != "DELETE" {
		t.Errorf("RemoveService(%q): wrong HTTP method. Want %q. Got %q.", id, "DELETE", req.Method)
	}
	expectedURL, _ := url.Parse(client.getURL("/services/" + id))
	if gotPath := req.URL.Path; gotPath != expectedURL.Path {
		t.Errorf("RemoveService(%q): Wrong path in request. Want %q. Got %q.", id, expectedURL.Path, gotPath)
	}
}

func TestRemoveServiceNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such service", status: http.StatusNotFound})
	err := client.RemoveService(RemoveServiceOptions{ID: "a2334"})
	expected := &NoSuchService{ID: "a2334"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("RemoveService: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}

func TestUpdateService(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2"
	update := UpdateServiceOptions{Version: 23}
	err := client.UpdateService(id, update)
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	if req.Method != "POST" {
		t.Errorf("UpdateService: wrong HTTP method. Want %q. Got %q.", "POST", req.Method)
	}
	expectedURL, _ := url.Parse(client.getURL("/services/" + id + "/update?version=23"))
	if gotURI := req.URL.RequestURI(); gotURI != expectedURL.RequestURI() {
		t.Errorf("UpdateService: Wrong path in request. Want %q. Got %q.", expectedURL.RequestURI(), gotURI)
	}
	expectedContentType := "application/json"
	if contentType := req.Header.Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("UpdateService: Wrong content-type in request. Want %q. Got %q.", expectedContentType, contentType)
	}
	var out UpdateServiceOptions
	if err := json.NewDecoder(req.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	update.Version = 0
	if !reflect.DeepEqual(out, update) {
		t.Errorf("UpdateService: wrong body\ngot  %#v\nwant %#v", out, update)
	}
}

func TestUpdateServiceNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such service", status: http.StatusNotFound})
	update := UpdateServiceOptions{}
	err := client.UpdateService("notfound", update)
	expected := &NoSuchService{ID: "notfound"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("UpdateService: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}

func TestInspectServiceNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such service", status: http.StatusNotFound})
	service, err := client.InspectService("notfound")
	if service != nil {
		t.Errorf("InspectService: Expected <nil> service, got %#v", service)
	}
	expected := &NoSuchService{ID: "notfound"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("InspectService: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}

func TestInspectService(t *testing.T) {
	jsonService := `{
  "ID": "ak7w3gjqoa3kuz8xcpnyy0pvl",
  "Version": {
    "Index": 95
  },
  "CreatedAt": "2016-06-07T21:10:20.269723157Z",
  "UpdatedAt": "2016-06-07T21:10:20.276301259Z",
  "Spec": {
    "Name": "redis",
    "Task": {
      "ContainerSpec": {
        "Image": "redis"
      },
      "Resources": {
        "Limits": {},
        "Reservations": {}
      },
      "RestartPolicy": {
        "Condition": "ANY"
      },
      "Placement": {}
    },
    "Mode": {
      "Replicated": {
        "Replicas": 1
      }
    },
    "UpdateConfig": {
      "Parallelism": 1
    },
    "EndpointSpec": {
      "Mode": "VIP",
      "Ingress": "PUBLICPORT",
      "ExposedPorts": [
        {
          "Protocol": "tcp",
          "Port": 6379
        }
      ]
    }
  },
  "Endpoint": {
    "Spec": {},
    "ExposedPorts": [
      {
        "Protocol": "tcp",
        "Port": 6379,
        "PublicPort": 30001
      }
    ],
    "VirtualIPs": [
      {
        "NetworkID": "4qvuz4ko70xaltuqbt8956gd1",
        "Addr": "10.255.0.4/16"
      }
    ]
  }
}`
	var expected swarm.Service
	err := json.Unmarshal([]byte(jsonService), &expected)
	if err != nil {
		t.Fatal(err)
	}
	fakeRT := &FakeRoundTripper{message: jsonService, status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "ak7w3gjqoa3kuz8xcpnyy0pvl"
	service, err := client.InspectService(id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*service, expected) {
		t.Errorf("InspectService(%q): Expected %#v. Got %#v.", id, expected, service)
	}
	expectedURL, _ := url.Parse(client.getURL("/services/ak7w3gjqoa3kuz8xcpnyy0pvl"))
	if gotPath := fakeRT.requests[0].URL.Path; gotPath != expectedURL.Path {
		t.Errorf("InspectService(%q): Wrong path in request. Want %q. Got %q.", id, expectedURL.Path, gotPath)
	}
}

func TestListServices(t *testing.T) {
	jsonServices := `[
  {
    "ID": "9mnpnzenvg8p8tdbtq4wvbkcz",
    "Version": {
      "Index": 19
    },
    "CreatedAt": "2016-06-07T21:05:51.880065305Z",
    "UpdatedAt": "2016-06-07T21:07:29.962229872Z",
    "Spec": {
      "Name": "hopeful_cori",
      "TaskTemplate": {
        "ContainerSpec": {
          "Image": "redis"
        },
        "Resources": {
          "Limits": {},
          "Reservations": {}
        },
        "RestartPolicy": {
          "Condition": "ANY"
        },
        "Placement": {}
      },
      "Mode": {
        "Replicated": {
          "Replicas": 1
        }
      },
      "UpdateConfig": {
        "Parallelism": 1
      },
      "EndpointSpec": {
        "Mode": "VIP",
        "Ingress": "PUBLICPORT",
        "ExposedPorts": [
          {
            "Protocol": "tcp",
            "Port": 6379
          }
        ]
      }
    },
    "Endpoint": {
      "Spec": {},
      "ExposedPorts": [
        {
          "Protocol": "tcp",
          "Port": 6379,
          "PublicPort": 30000
        }
      ],
      "VirtualIPs": [
        {
          "NetworkID": "4qvuz4ko70xaltuqbt8956gd1",
          "Addr": "10.255.0.2/16"
        },
        {
          "NetworkID": "4qvuz4ko70xaltuqbt8956gd1",
          "Addr": "10.255.0.3/16"
        }
      ]
    }
  }
]`
	var expected []swarm.Service
	err := json.Unmarshal([]byte(jsonServices), &expected)
	if err != nil {
		t.Fatal(err)
	}
	client := newTestClient(&FakeRoundTripper{message: jsonServices, status: http.StatusOK})
	services, err := client.ListServices(ListServicesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(services, expected) {
		t.Errorf("ListServices: Expected %#v. Got %#v.", expected, services)
	}
}
