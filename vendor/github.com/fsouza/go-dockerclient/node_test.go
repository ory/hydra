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

func TestListNodes(t *testing.T) {
	jsonNodes := `[
  {
    "ID": "24ifsmvkjbyhk",
    "Version": {
      "Index": 8
    },
    "CreatedAt": "2016-06-07T20:31:11.853781916Z",
    "UpdatedAt": "2016-06-07T20:31:11.999868824Z",
    "Spec": {
      "Name": "my-node",
      "Role": "manager",
      "Availability": "active",
      "Labels": {
          "foo": "bar"
      }
    },
    "Description": {
      "Hostname": "bf3067039e47",
      "Platform": {
        "Architecture": "x86_64",
        "OS": "linux"
      },
      "Resources": {
        "NanoCPUs": 4000000000,
        "MemoryBytes": 8272408576
      },
      "Engine": {
        "EngineVersion": "1.12.0-dev",
        "Labels": {
            "foo": "bar"
        },
        "Plugins": [
          {
            "Type": "Volume",
            "Name": "local"
          },
          {
            "Type": "Network",
            "Name": "bridge"
          },
          {
            "Type": "Network",
            "Name": "null"
          },
          {
            "Type": "Network",
            "Name": "overlay"
          }
        ]
      }
    },
    "Status": {
      "State": "ready"
    },
    "ManagerStatus": {
      "Leader": true,
      "Reachability": "reachable",
      "Addr": "172.17.0.2:2377"
    }
  }
]`
	var expected []swarm.Node
	err := json.Unmarshal([]byte(jsonNodes), &expected)
	if err != nil {
		t.Fatal(err)
	}
	client := newTestClient(&FakeRoundTripper{message: jsonNodes, status: http.StatusOK})
	nodes, err := client.ListNodes(ListNodesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(nodes, expected) {
		t.Errorf("ListNodes: Expected %#v. Got %#v.", expected, nodes)
	}

}

func TestInspectNode(t *testing.T) {
	jsonNode := `{
  "ID": "24ifsmvkjbyhk",
  "Version": {
    "Index": 8
  },
  "CreatedAt": "2016-06-07T20:31:11.853781916Z",
  "UpdatedAt": "2016-06-07T20:31:11.999868824Z",
  "Spec": {
    "Name": "my-node",
    "Role": "manager",
    "Availability": "active",
    "Labels": {
        "foo": "bar"
    }
  },
  "Description": {
    "Hostname": "bf3067039e47",
    "Platform": {
      "Architecture": "x86_64",
      "OS": "linux"
    },
    "Resources": {
      "NanoCPUs": 4000000000,
      "MemoryBytes": 8272408576
    },
    "Engine": {
      "EngineVersion": "1.12.0-dev",
      "Labels": {
          "foo": "bar"
      },
      "Plugins": [
        {
          "Type": "Volume",
          "Name": "local"
        },
        {
          "Type": "Network",
          "Name": "bridge"
        },
        {
          "Type": "Network",
          "Name": "null"
        },
        {
          "Type": "Network",
          "Name": "overlay"
        }
      ]
    }
  },
  "Status": {
    "State": "ready"
  },
  "ManagerStatus": {
    "Leader": true,
    "Reachability": "reachable",
    "Addr": "172.17.0.2:2377"
  }
}`

	var expected swarm.Node
	err := json.Unmarshal([]byte(jsonNode), &expected)
	if err != nil {
		t.Fatal(err)
	}
	fakeRT := &FakeRoundTripper{message: jsonNode, status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "24ifsmvkjbyhk"
	node, err := client.InspectNode(id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*node, expected) {
		t.Errorf("InspectNode(%q): Expected %#v. Got %#v.", id, expected, node)
	}
	expectedURL, _ := url.Parse(client.getURL("/nodes/24ifsmvkjbyhk"))
	if gotPath := fakeRT.requests[0].URL.Path; gotPath != expectedURL.Path {
		t.Errorf("InspectNode(%q): Wrong path in request. Want %q. Got %q.", id, expectedURL.Path, gotPath)
	}

}

func TestInspectNodeNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such node", status: http.StatusNotFound})
	node, err := client.InspectNode("notfound")
	if node != nil {
		t.Errorf("InspectNode: Expected <nil> task, got %#v", node)
	}
	expected := &NoSuchNode{ID: "notfound"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("InspectNode: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}

func TestUpdateNode(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2"
	opts := UpdateNodeOptions{}
	err := client.UpdateNode(id, opts)
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	if req.Method != "POST" {
		t.Errorf("UpdateNode: wrong HTTP method. Want %q. Got %q.", "POST", req.Method)
	}
	expectedURL, _ := url.Parse(client.getURL("/nodes/" + id + "/update"))
	if gotPath := req.URL.Path; gotPath != expectedURL.Path {
		t.Errorf("UpdateNode: Wrong path in request. Want %q. Got %q.", expectedURL.Path, gotPath)
	}
	expectedContentType := "application/json"
	if contentType := req.Header.Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("UpdateNode: Wrong content-type in request. Want %q. Got %q.", expectedContentType, contentType)
	}
	var out UpdateNodeOptions
	if err := json.NewDecoder(req.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, opts) {
		t.Errorf("UpdateNode: wrong body, got: %#v, want %#v", out, opts)
	}
}

func TestUpdateNodeNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such node", status: http.StatusNotFound})
	err := client.UpdateNode("notfound", UpdateNodeOptions{})
	expected := &NoSuchNode{ID: "notfound"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("UpdateNode: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}

func TestRemoveNode(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "4fa6e0f0c6786287e131c3852c58a2e01cc697a68231826813597e4994f1d6e2"
	err := client.RemoveNode(RemoveNodeOptions{ID: id})
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	if req.Method != "DELETE" {
		t.Errorf("RemoveNode(%q): wrong HTTP method. Want %q. Got %q.", id, "DELETE", req.Method)
	}
	expectedURL, _ := url.Parse(client.getURL("/nodes/" + id))
	if gotPath := req.URL.Path; gotPath != expectedURL.Path {
		t.Errorf("RemoveNode(%q): Wrong path in request. Want %q. Got %q.", id, expectedURL.Path, gotPath)
	}
}

func TestRemoveNodeNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such node", status: http.StatusNotFound})
	err := client.RemoveNode(RemoveNodeOptions{ID: "notfound"})
	expected := &NoSuchNode{ID: "notfound"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("RemoveNode: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}
