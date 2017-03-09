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

func TestListTasks(t *testing.T) {
	jsonTasks := `[
  {
    "ID": "0kzzo1i0y4jz6027t0k7aezc7",
    "Version": {
      "Index": 71
    },
    "CreatedAt": "2016-06-07T21:07:31.171892745Z",
    "UpdatedAt": "2016-06-07T21:07:31.376370513Z",
    "Name": "hopeful_cori",
    "Spec": {
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
    "ServiceID": "9mnpnzenvg8p8tdbtq4wvbkcz",
    "Instance": 1,
    "NodeID": "24ifsmvkjbyhk",
    "ServiceAnnotations": {},
    "Status": {
      "Timestamp": "2016-06-07T21:07:31.290032978Z",
      "State": "FAILED",
      "Message": "execution failed",
      "ContainerStatus": {}
    },
    "DesiredState": "SHUTDOWN",
    "NetworksAttachments": [
      {
        "Network": {
          "ID": "4qvuz4ko70xaltuqbt8956gd1",
          "Version": {
            "Index": 18
          },
          "CreatedAt": "2016-06-07T20:31:11.912919752Z",
          "UpdatedAt": "2016-06-07T21:07:29.955277358Z",
          "Spec": {
            "Name": "ingress",
            "Labels": {
              "com.docker.swarm.internal": "true"
            },
            "DriverConfiguration": {},
            "IPAM": {
              "Driver": {},
              "Configs": [
                {
                  "Family": "UNKNOWN",
                  "Subnet": "10.255.0.0/16"
                }
              ]
            }
          },
          "DriverState": {
            "Name": "overlay",
            "Options": {
              "com.docker.network.driver.overlay.vxlanid_list": "256"
            }
          },
          "IPAM": {
            "Driver": {
              "Name": "default"
            },
            "Configs": [
              {
                "Family": "UNKNOWN",
                "Subnet": "10.255.0.0/16"
              }
            ]
          }
        },
        "Addresses": [
          "10.255.0.10/16"
        ]
      }
    ],
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
  },
  {
    "ID": "1yljwbmlr8er2waf8orvqpwms",
    "Version": {
      "Index": 30
    },
    "CreatedAt": "2016-06-07T21:07:30.019104782Z",
    "UpdatedAt": "2016-06-07T21:07:30.231958098Z",
    "Name": "hopeful_cori",
    "Spec": {
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
    "ServiceID": "9mnpnzenvg8p8tdbtq4wvbkcz",
    "Instance": 1,
    "NodeID": "24ifsmvkjbyhk",
    "ServiceAnnotations": {},
    "Status": {
      "Timestamp": "2016-06-07T21:07:30.202183143Z",
      "State": "FAILED",
      "Message": "execution failed",
      "ContainerStatus": {}
    },
    "DesiredState": "SHUTDOWN",
    "NetworksAttachments": [
      {
        "Network": {
          "ID": "4qvuz4ko70xaltuqbt8956gd1",
          "Version": {
            "Index": 18
          },
          "CreatedAt": "2016-06-07T20:31:11.912919752Z",
          "UpdatedAt": "2016-06-07T21:07:29.955277358Z",
          "Spec": {
            "Name": "ingress",
            "Labels": {
              "com.docker.swarm.internal": "true"
            },
            "DriverConfiguration": {},
            "IPAM": {
              "Driver": {},
              "Configs": [
                {
                  "Family": "UNKNOWN",
                  "Subnet": "10.255.0.0/16"
                }
              ]
            }
          },
          "DriverState": {
            "Name": "overlay",
            "Options": {
              "com.docker.network.driver.overlay.vxlanid_list": "256"
            }
          },
          "IPAM": {
            "Driver": {
              "Name": "default"
            },
            "Configs": [
              {
                "Family": "UNKNOWN",
                "Subnet": "10.255.0.0/16"
              }
            ]
          }
        },
        "Addresses": [
          "10.255.0.5/16"
        ]
      }
    ],
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
	var expected []swarm.Task
	err := json.Unmarshal([]byte(jsonTasks), &expected)
	if err != nil {
		t.Fatal(err)
	}
	client := newTestClient(&FakeRoundTripper{message: jsonTasks, status: http.StatusOK})
	tasks, err := client.ListTasks(ListTasksOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tasks, expected) {
		t.Errorf("ListTasks: Expected %#v. Got %#v.", expected, tasks)
	}

}

func TestInspectTask(t *testing.T) {
	jsonTask := `{
  "ID": "0kzzo1i0y4jz6027t0k7aezc7",
  "Version": {
    "Index": 71
  },
  "CreatedAt": "2016-06-07T21:07:31.171892745Z",
  "UpdatedAt": "2016-06-07T21:07:31.376370513Z",
  "Name": "hopeful_cori",
  "Spec": {
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
  "ServiceID": "9mnpnzenvg8p8tdbtq4wvbkcz",
  "Instance": 1,
  "NodeID": "24ifsmvkjbyhk",
  "ServiceAnnotations": {},
  "Status": {
    "Timestamp": "2016-06-07T21:07:31.290032978Z",
    "State": "FAILED",
    "Message": "execution failed",
    "ContainerStatus": {}
  },
  "DesiredState": "SHUTDOWN",
  "NetworksAttachments": [
    {
      "Network": {
        "ID": "4qvuz4ko70xaltuqbt8956gd1",
        "Version": {
          "Index": 18
        },
        "CreatedAt": "2016-06-07T20:31:11.912919752Z",
        "UpdatedAt": "2016-06-07T21:07:29.955277358Z",
        "Spec": {
          "Name": "ingress",
          "Labels": {
            "com.docker.swarm.internal": "true"
          },
          "DriverConfiguration": {},
          "IPAM": {
            "Driver": {},
            "Configs": [
              {
                "Family": "UNKNOWN",
                "Subnet": "10.255.0.0/16"
              }
            ]
          }
        },
        "DriverState": {
          "Name": "overlay",
          "Options": {
            "com.docker.network.driver.overlay.vxlanid_list": "256"
          }
        },
        "IPAM": {
          "Driver": {
            "Name": "default"
          },
          "Configs": [
            {
              "Family": "UNKNOWN",
              "Subnet": "10.255.0.0/16"
            }
          ]
        }
      },
      "Addresses": [
        "10.255.0.10/16"
      ]
    }
  ],
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
}`

	var expected swarm.Task
	err := json.Unmarshal([]byte(jsonTask), &expected)
	if err != nil {
		t.Fatal(err)
	}
	fakeRT := &FakeRoundTripper{message: jsonTask, status: http.StatusOK}
	client := newTestClient(fakeRT)
	id := "0kzzo1i0y4jz6027t0k7aezc7"
	task, err := client.InspectTask(id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*task, expected) {
		t.Errorf("InspectTask(%q): Expected %#v. Got %#v.", id, expected, task)
	}
	expectedURL, _ := url.Parse(client.getURL("/tasks/0kzzo1i0y4jz6027t0k7aezc7"))
	if gotPath := fakeRT.requests[0].URL.Path; gotPath != expectedURL.Path {
		t.Errorf("InspectTask(%q): Wrong path in request. Want %q. Got %q.", id, expectedURL.Path, gotPath)
	}

}

func TestInspectTaskNotFound(t *testing.T) {
	client := newTestClient(&FakeRoundTripper{message: "no such task", status: http.StatusNotFound})
	task, err := client.InspectTask("notfound")
	if task != nil {
		t.Errorf("InspectTask: Expected <nil> task, got %#v", task)
	}
	expected := &NoSuchTask{ID: "notfound"}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("InspectTask: Wrong error returned. Want %#v. Got %#v.", expected, err)
	}
}
