// Copyright 2016 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/fsouza/go-dockerclient"
)

func TestSwarmInit(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/init", bytes.NewReader(nil))
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmInit: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var id string
	err = json.Unmarshal(recorder.Body.Bytes(), &id)
	if err != nil {
		t.Fatalf("SwarmInit: got error. %s", err.Error())
	}
	if id == "" {
		t.Fatal("SwarmInit: id not found")
	}
	if server.swarm == nil {
		t.Fatalf("SwarmInit: expected swarm to be set.")
	}
	if len(server.nodes) != 1 {
		t.Fatalf("SwarmInit: expected node len to be 1, got: %d", len(server.nodes))
	}
	if server.nodes[0].ManagerStatus.Addr != server.SwarmAddress() {
		t.Fatalf("SwarmInit: expected current node to have addr %q, got: %q", server.SwarmAddress(), server.nodes[0].ManagerStatus.Addr)
	}
	if !server.nodes[0].ManagerStatus.Leader {
		t.Fatalf("SwarmInit: expected current node to be leader")
	}
}

func TestSwarmInitDynamicAdvertiseAddrPort(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	data := `{"ListenAddr": "127.0.0.1:0", "AdvertiseAddr": "localhost"}`
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/init", strings.NewReader(data))
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmInit: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server.nodes) != 1 {
		t.Fatalf("SwarmInit: expected node len to be 1, got: %d", len(server.nodes))
	}
	_, port, _ := net.SplitHostPort(server.SwarmAddress())
	expectedAddr := fmt.Sprintf("localhost:%s", port)
	if server.nodes[0].ManagerStatus.Addr != expectedAddr {
		t.Fatalf("SwarmInit: expected current node to have addr %q, got: %q", expectedAddr, server.nodes[0].ManagerStatus.Addr)
	}
}

func TestSwarmInitAlreadyInSwarm(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	server.swarm = &swarm.Swarm{}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/init", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotAcceptable {
		t.Fatalf("SwarmInit: wrong status. Want %d. Got %d.", http.StatusNotAcceptable, recorder.Code)
	}
}

func TestSwarmJoinNoBody(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/join", bytes.NewReader(nil))
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("SwarmJoin: wrong status. Want %d. Got %d.", http.StatusInternalServerError, recorder.Code)
	}
	if server.swarm != nil {
		t.Fatalf("SwarmJoin: expected swarm not to be set.")
	}
}

func TestSwarmJoin(t *testing.T) {
	server1, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server1.Stop()
	server2, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server2.Stop()
	data, err := json.Marshal(swarm.InitRequest{})
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/init", bytes.NewReader(data))
	server1.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmJoin: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	data, err = json.Marshal(swarm.JoinRequest{
		RemoteAddrs: []string{server1.SwarmAddress()},
	})
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/swarm/join", bytes.NewReader(data))
	server2.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmJoin: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if server1.swarm == nil || server2.swarm == nil {
		t.Fatalf("SwarmJoin: expected swarm to be set.")
	}
	if len(server1.nodes) != 2 {
		t.Fatalf("SwarmJoin: expected node len to be 2, got: %d", len(server1.nodes))
	}
	if server1.nodes[0].ManagerStatus.Addr != server1.SwarmAddress() {
		t.Fatalf("SwarmJoin: expected nodes[0] to have addr %q, got: %q", server1.SwarmAddress(), server1.nodes[0].ManagerStatus.Addr)
	}
	if server1.nodes[1].ManagerStatus.Leader {
		t.Fatalf("SwarmInit: expected nodes[1] not to be leader")
	}
	if server1.nodes[1].ManagerStatus.Addr != server2.SwarmAddress() {
		t.Fatalf("SwarmJoin: expected nodes[1] to have addr %q, got: %q", server2.SwarmAddress(), server1.nodes[1].ManagerStatus.Addr)
	}
	if !reflect.DeepEqual(server1.nodes, server2.nodes) {
		t.Fatalf("SwarmJoin: expected nodes to be equal in server1 and server2, got:\n%#v\n%#v", server1.nodes, server2.nodes)
	}
}

func TestSwarmJoinWithService(t *testing.T) {
	server1, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server1.Stop()
	server2, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server2.Stop()
	data, err := json.Marshal(swarm.InitRequest{})
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/init", bytes.NewReader(data))
	server1.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmJoin: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	serviceCreateOpts := docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: "test/test",
				},
			},
		},
	}
	buf, err := json.Marshal(serviceCreateOpts)
	if err != nil {
		t.Fatalf("ServiceCreate error: %s", err.Error())
	}
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/services/create", bytes.NewBuffer(buf))
	server1.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmJoin: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	data, err = json.Marshal(swarm.JoinRequest{
		RemoteAddrs: []string{server1.SwarmAddress()},
	})
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/swarm/join", bytes.NewReader(data))
	server2.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmJoin: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server1.services) != len(server2.services) {
		t.Fatalf("SwarmJoin: expected len services to be equal in server1 and server2, got:\n%#v\n%#v", len(server1.services), len(server2.services))
	}
	if !compareServices(server1.services[0], server2.services[0]) {
		t.Fatalf("SwarmJoin: expected services to be equal in server1 and server2, got:\n%#v\n%#v", server1.services[0], server2.services[0])
	}
	if len(server1.tasks) != len(server2.tasks) {
		t.Fatalf("SwarmJoin: expected len tasks to be equal in server1 and server2, got:\n%#v\n%#v", len(server1.tasks), len(server2.tasks))
	}
	if !compareTasks(server1.tasks[0], server2.tasks[0]) {
		t.Fatalf("SwarmJoin: expected tasks to be equal in server1 and server2, got:\n%#v\n%#v", server1.tasks[0], server2.tasks[0])
	}
}

func TestSwarmJoinAlreadyInSwarm(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	server.swarm = &swarm.Swarm{}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/join", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotAcceptable {
		t.Fatalf("SwarmJoin: wrong status. Want %d. Got %d.", http.StatusNotAcceptable, recorder.Code)
	}
}

func TestSwarmLeave(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	server.swarm = &swarm.Swarm{}
	server.swarmServer, _ = newSwarmServer(server, "127.0.0.1:0")
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/leave", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmLeave: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if server.swarm != nil {
		t.Fatalf("SwarmLeave: expected swarm to be nil. Got %+v.", server.swarm)
	}
}

func TestSwarmLeaveNotInSwarm(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/leave", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotAcceptable {
		t.Fatalf("SwarmLeave: wrong status. Want %d. Got %d.", http.StatusNotAcceptable, recorder.Code)
	}
	if server.swarm != nil {
		t.Fatalf("SwarmLeave: expected swarm to be nil. Got %+v.", server.swarm)
	}
}

func TestSwarmInspect(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	expected := &swarm.Swarm{
		ClusterInfo: swarm.ClusterInfo{
			ID: "swarm-id",
		},
	}
	server.swarm = expected
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/swarm", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("SwarmInspect: wrong status. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var swarmInspect *swarm.Swarm
	err = json.Unmarshal(recorder.Body.Bytes(), &swarmInspect)
	if err != nil {
		t.Fatalf("SwarmInspect: got error. %s", err.Error())
	}
	if expected.ClusterInfo.ID != swarmInspect.ClusterInfo.ID {
		t.Fatalf("SwarmInspect: wrong response. Want %+v. Got %+v.", expected, swarmInspect)
	}
}

func TestSwarmInspectNotInSwarm(t *testing.T) {
	server, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	server.buildMuxer()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/swarm", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotAcceptable {
		t.Fatalf("SwarmInspect: wrong status. Want %d. Got %d.", http.StatusNotAcceptable, recorder.Code)
	}
}

func TestServiceCreate(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	recorder := httptest.NewRecorder()
	serviceCreateOpts := docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "test",
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image:   "test/test",
					Command: []string{"sh"},
					Args:    []string{"--test"},
					Env:     []string{"ENV=1"},
					User:    "test",
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode: swarm.ResolutionModeVIP,
				Ports: []swarm.PortConfig{{
					Protocol:      swarm.PortConfigProtocolTCP,
					TargetPort:    uint32(80),
					PublishedPort: uint32(80),
				}},
			},
		},
	}
	buf, err := json.Marshal(serviceCreateOpts)
	if err != nil {
		t.Fatalf("ServiceCreate error: %s", err.Error())
	}
	var params io.Reader
	params = bytes.NewBuffer(buf)
	request, _ := http.NewRequest("POST", "/services/create", params)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceCreate: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server.services) != 1 || len(server.tasks) != 1 || len(server.containers) != 1 {
		t.Fatalf("ServiceCreate: wrong item count. Want 1. Got services: %d, tasks: %d, containers: %d.", len(server.services), len(server.tasks), len(server.containers))
	}
	cont := server.containers[0]
	expectedContainer := &docker.Container{
		ID:      cont.ID,
		Created: cont.Created,
		Image:   "test/test",
		Name:    "test-0",
		Config: &docker.Config{
			Entrypoint: []string{"sh"},
			Cmd:        []string{"--test"},
			Env:        []string{"ENV=1"},
		},
		HostConfig: &docker.HostConfig{},
		State: docker.State{
			Running:   true,
			StartedAt: cont.State.StartedAt,
			Pid:       cont.State.Pid,
			ExitCode:  0,
		},
	}
	if !reflect.DeepEqual(cont, expectedContainer) {
		t.Fatalf("ServiceCreate: wrong cont. Want\n%#v\nGot\n%#v", expectedContainer, cont)
	}
	srv := server.services[0]
	expectedService := &swarm.Service{
		ID:   srv.ID,
		Spec: serviceCreateOpts.ServiceSpec,
		Endpoint: swarm.Endpoint{
			Spec:  *serviceCreateOpts.ServiceSpec.EndpointSpec,
			Ports: []swarm.PortConfig{{Protocol: "tcp", TargetPort: 80, PublishedPort: 80}},
		},
	}
	if !reflect.DeepEqual(srv, expectedService) {
		t.Fatalf("ServiceCreate: wrong service. Want\n%#v\nGot\n%#v", expectedService, srv)
	}
	task := server.tasks[0]
	expectedTask := &swarm.Task{
		ID:        task.ID,
		ServiceID: srv.ID,
		NodeID:    server.nodes[0].ID,
		Status: swarm.TaskStatus{
			State: swarm.TaskStateReady,
			ContainerStatus: swarm.ContainerStatus{
				ContainerID: cont.ID,
			},
		},
		DesiredState: swarm.TaskStateReady,
		Spec:         serviceCreateOpts.ServiceSpec.TaskTemplate,
	}
	if !reflect.DeepEqual(task, expectedTask) {
		t.Fatalf("ServiceCreate: wrong task. Want\n%#v\nGot\n%#v", expectedTask, task)
	}
}

func TestServiceCreateDynamicPort(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	recorder := httptest.NewRecorder()
	serviceCreateOpts := docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "test",
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image:   "test/test",
					Command: []string{"sh"},
					Args:    []string{"--test"},
					Env:     []string{"ENV=1"},
					User:    "test",
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode: swarm.ResolutionModeVIP,
				Ports: []swarm.PortConfig{{
					Protocol:      swarm.PortConfigProtocolTCP,
					TargetPort:    uint32(80),
					PublishedPort: uint32(0),
				}},
			},
		},
	}
	buf, err := json.Marshal(serviceCreateOpts)
	if err != nil {
		t.Fatalf("ServiceCreate error: %s", err.Error())
	}
	var params io.Reader
	params = bytes.NewBuffer(buf)
	request, _ := http.NewRequest("POST", "/services/create", params)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceCreate: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server.services) != 1 || len(server.tasks) != 1 || len(server.containers) != 1 {
		t.Fatalf("ServiceCreate: wrong item count. Want 1. Got services: %d, tasks: %d, containers: %d.", len(server.services), len(server.tasks), len(server.containers))
	}
	srv := server.services[0]
	expectedService := &swarm.Service{
		ID:   srv.ID,
		Spec: serviceCreateOpts.ServiceSpec,
		Endpoint: swarm.Endpoint{
			Spec:  *serviceCreateOpts.ServiceSpec.EndpointSpec,
			Ports: []swarm.PortConfig{{Protocol: "tcp", TargetPort: 80, PublishedPort: 30000}},
		},
	}
	if !reflect.DeepEqual(srv, expectedService) {
		t.Fatalf("ServiceCreate: wrong service. Want\n%#v\nGot\n%#v", expectedService, srv)
	}
}

func TestServiceCreateMultipleServers(t *testing.T) {
	server1, server2 := setUpSwarm(t)
	defer server1.Stop()
	defer server2.Stop()
	_, err := addTestService(server1)
	if err != nil {
		t.Fatal(err)
	}
	if len(server1.services) != 1 {
		t.Fatalf("ServiceCreate: expected services to have len 1, got: %d", len(server1.services))
	}
	if len(server1.services) != len(server2.services) {
		t.Fatalf("ServiceCreate: expected len services to be equal in server1 and server2, got:\n%#v\n%#v", len(server1.services), len(server2.services))
	}
	if !compareServices(server1.services[0], server2.services[0]) {
		t.Fatalf("ServiceCreate: expected services to be equal in server1 and server2, got:\n%#v\n%#v", server1.services[0], server2.services[0])
	}
	if len(server1.tasks) != len(server2.tasks) {
		t.Fatalf("ServiceCreate: expected len tasks to be equal in server1 and server2, got:\n%#v\n%#v", len(server1.tasks), len(server2.tasks))
	}
	if !compareTasks(server1.tasks[0], server2.tasks[0]) {
		t.Fatalf("ServiceCreate: expected tasks to be equal in server1 and server2, got:\n%#v\n%#v", server1.tasks[0], server2.tasks[0])
	}
}

func compareServices(srv1 *swarm.Service, srv2 *swarm.Service) bool {
	srv1.CreatedAt = srv2.CreatedAt
	srv1.UpdatedAt = srv2.UpdatedAt
	if srv1.UpdateStatus != nil {
		srv1.UpdateStatus.StartedAt = srv2.UpdateStatus.StartedAt
		srv1.UpdateStatus.CompletedAt = srv2.UpdateStatus.CompletedAt
	}
	return reflect.DeepEqual(srv1, srv2)
}

func compareTasks(task1 *swarm.Task, task2 *swarm.Task) bool {
	task1.CreatedAt = task2.CreatedAt
	task1.UpdatedAt = task2.UpdatedAt
	task1.Status.Timestamp = task2.Status.Timestamp
	return reflect.DeepEqual(task1, task2)
}

func TestServiceInspect(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/services/"+srv.ID, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceInspect: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var srvInspect swarm.Service
	err = json.Unmarshal(recorder.Body.Bytes(), &srvInspect)
	if err != nil {
		t.Fatalf("ServiceInspect: unable to unmarshal response body: %s", err)
	}
	if !compareServices(srv, &srvInspect) {
		t.Fatalf("ServiceInspect: wrong service. Want\n%#v\nGot\n%#v", srv, &srvInspect)
	}
}

func TestServiceInspectByName(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/services/"+srv.Spec.Name, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceInspect: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var srvInspect swarm.Service
	err = json.Unmarshal(recorder.Body.Bytes(), &srvInspect)
	if err != nil {
		t.Fatalf("ServiceInspect: unable to unmarshal response body: %s", err)
	}
	if !compareServices(srv, &srvInspect) {
		t.Fatalf("ServiceInspect: wrong service. Want\n%#v\nGot\n%#v", srv, &srvInspect)
	}
}

func TestServiceInspectNotFound(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/services/abcd", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("ServiceInspect: wrong status code. Want %d. Got %d.", http.StatusNotFound, recorder.Code)
	}
}

func TestTaskInspect(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/tasks/"+task.ID, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskInspect: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskInspect: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect) {
		t.Fatalf("TaskInspect: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
}

func TestTaskInspectNotFound(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/tasks/abcd", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("TaskInspect: wrong status code. Want %d. Got %d.", http.StatusNotFound, recorder.Code)
	}
}

func TestServiceList(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/services", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var srvInspect []swarm.Service
	err = json.Unmarshal(recorder.Body.Bytes(), &srvInspect)
	if err != nil {
		t.Fatalf("ServiceList: unable to unmarshal response body: %s", err)
	}
	if !compareServices(srv, &srvInspect[0]) {
		t.Fatalf("ServiceList: wrong service. Want\n%#v\nGot\n%#v", srv, &srvInspect)
	}
}

func TestServiceListFilterID(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/services?filters={"id":[%q]}`, srv.ID), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var srvInspect []swarm.Service
	err = json.Unmarshal(recorder.Body.Bytes(), &srvInspect)
	if err != nil {
		t.Fatalf("ServiceList: unable to unmarshal response body: %s", err)
	}
	if !compareServices(srv, &srvInspect[0]) {
		t.Fatalf("ServiceList: wrong service. Want\n%#v\nGot\n%#v", srv, &srvInspect)
	}
}

func TestServiceListFilterName(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/services?filters={"name":[%q]}`, srv.Spec.Name), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var srvInspect []swarm.Service
	err = json.Unmarshal(recorder.Body.Bytes(), &srvInspect)
	if err != nil {
		t.Fatalf("ServiceList: unable to unmarshal response body: %s", err)
	}
	if !compareServices(srv, &srvInspect[0]) {
		t.Fatalf("ServiceList: wrong service. Want\n%#v\nGot\n%#v", srv, &srvInspect)
	}
}

func TestServiceListFilterEmpty(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", `/services?filters={"id":["something"]}`, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var srvInspect []swarm.Service
	err = json.Unmarshal(recorder.Body.Bytes(), &srvInspect)
	if err != nil {
		t.Fatalf("ServiceList: unable to unmarshal response body: %s", err)
	}
	if len(srvInspect) != 0 {
		t.Fatalf("ServiceList: expected empty list got %d.", len(srvInspect))
	}
}

func TestTaskList(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/tasks", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
}

func TestTaskListFilterID(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/tasks?filters={"id":[%q]}`, task.ID), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
}

func TestTaskListFilterServiceID(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/tasks?filters={"service":[%q]}`, task.ServiceID), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
}

func TestTaskListFilterServiceName(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/tasks?filters={"service":[%q]}`, srv.Spec.Name), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
}

func TestTaskListFilterMultipleFields(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/tasks?filters={"service":[%q], "id":[%q]}`, srv.Spec.Name, task.ID), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
}

func TestTaskListFilterMultipleFieldsNotFound(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf(`/tasks?filters={"service":[%q], "id":["abc"]}`, srv.Spec.Name), nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if len(taskInspect) != 0 {
		t.Fatalf("TaskList: Want\nempty task list\nGot\n%#v", &taskInspect)
	}
}

func TestTaskListFilterNotFound(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", `/tasks?filters={"id":["something"]}`, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if len(taskInspect) != 0 {
		t.Fatalf("TaskList: expected empty list got %d.", len(taskInspect))
	}
}

func TestTaskListFilterLabel(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	_, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	task := server.tasks[0]
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", `/tasks?filters={"label":["mykey=myvalue"]}`, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	var taskInspect []swarm.Task
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
	request, _ = http.NewRequest("GET", `/tasks?filters={"label":["mykey"]}`, nil)
	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if !compareTasks(task, &taskInspect[0]) {
		t.Fatalf("TaskList: wrong task. Want\n%#v\nGot\n%#v", task, &taskInspect)
	}
	request, _ = http.NewRequest("GET", `/tasks?filters={"label":["otherkey"]}`, nil)
	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TaskList: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &taskInspect)
	if err != nil {
		t.Fatalf("TaskList: unable to unmarshal response body: %s", err)
	}
	if len(taskInspect) != 0 {
		t.Fatalf("TaskList: wrong size. Want 0, Got %d", len(taskInspect))
	}
}

func TestServiceDelete(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/services/"+srv.ID, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceDelete: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server.services) != 0 {
		t.Fatalf("ServiceDelete: expected empty services, got %d", len(server.services))
	}
	if len(server.tasks) != 0 {
		t.Fatalf("ServiceDelete: expected empty tasks, got %d", len(server.tasks))
	}
	if len(server.containers) != 0 {
		t.Fatalf("ServiceDelete: expected empty containers, got %d", len(server.containers))
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/services/blahblah", nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("ServiceDelete: wrong status code. Want %d. Got %d.", http.StatusNotFound, recorder.Code)
	}
}

func TestServiceUpdate(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	updateOpts := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "test",
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: "test/test2",
				Args:  []string{"--test2"},
				Env:   []string{"ENV=2"},
				User:  "test",
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Mode: swarm.ResolutionModeVIP,
			Ports: []swarm.PortConfig{{
				Protocol:      swarm.PortConfigProtocolTCP,
				TargetPort:    uint32(80),
				PublishedPort: uint32(80),
			}},
		},
	}
	buf, err := json.Marshal(updateOpts)
	if err != nil {
		t.Fatalf("ServiceUpdate error: %s", err.Error())
	}
	request, _ := http.NewRequest("POST", fmt.Sprintf("/services/%s/update", srv.ID), bytes.NewReader(buf))
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceUpdate: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server.services) != 1 || len(server.tasks) != 1 || len(server.containers) != 1 {
		t.Fatalf("ServiceUpdate: wrong item count. Want 1. Got services: %d, tasks: %d, containers: %d.", len(server.services), len(server.tasks), len(server.containers))
	}
	cont := server.containers[0]
	expectedContainer := &docker.Container{
		ID:      cont.ID,
		Created: cont.Created,
		Image:   "test/test2",
		Name:    "test-0-updated",
		Config: &docker.Config{
			Cmd: []string{"--test2"},
			Env: []string{"ENV=2"},
		},
		HostConfig: &docker.HostConfig{},
		State: docker.State{
			Running:   true,
			StartedAt: cont.State.StartedAt,
			Pid:       cont.State.Pid,
			ExitCode:  0,
		},
	}
	if !reflect.DeepEqual(cont, expectedContainer) {
		t.Fatalf("ServiceUpdate: wrong cont. Want\n%#v\nGot\n%#v", expectedContainer, cont)
	}
	srv = server.services[0]
	expectedService := &swarm.Service{
		ID:   srv.ID,
		Spec: updateOpts,
		Endpoint: swarm.Endpoint{
			Spec:  *updateOpts.EndpointSpec,
			Ports: []swarm.PortConfig{{Protocol: "tcp", TargetPort: 80, PublishedPort: 80}},
		},
	}
	if !reflect.DeepEqual(srv, expectedService) {
		t.Fatalf("ServiceUpdate: wrong service. Want\n%#v\nGot\n%#v", expectedService, srv)
	}
	task := server.tasks[0]
	expectedTask := &swarm.Task{
		ID:        task.ID,
		ServiceID: srv.ID,
		NodeID:    server.nodes[1].ID,
		Status: swarm.TaskStatus{
			State: swarm.TaskStateReady,
			ContainerStatus: swarm.ContainerStatus{
				ContainerID: cont.ID,
			},
		},
		DesiredState: swarm.TaskStateReady,
		Spec:         updateOpts.TaskTemplate,
	}
	if !reflect.DeepEqual(task, expectedTask) {
		t.Fatalf("ServiceUpdate: wrong task. Want\n%#v\nGot\n%#v", expectedTask, task)
	}
}

func TestServiceUpdateMoreReplicas(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	srv, err := addTestService(server)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	replicas := uint64(3)
	updateOpts := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "test",
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: "test/test2",
				Args:  []string{"--test2"},
				Env:   []string{"ENV=2"},
				User:  "test",
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Mode: swarm.ResolutionModeVIP,
			Ports: []swarm.PortConfig{{
				Protocol:      swarm.PortConfigProtocolTCP,
				TargetPort:    uint32(80),
				PublishedPort: uint32(80),
			}},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			},
		},
	}
	buf, err := json.Marshal(updateOpts)
	if err != nil {
		t.Fatalf("ServiceUpdate error: %s", err.Error())
	}
	request, _ := http.NewRequest("POST", fmt.Sprintf("/services/%s/update", srv.ID), bytes.NewReader(buf))
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ServiceUpdate: wrong status code. Want %d. Got %d.", http.StatusOK, recorder.Code)
	}
	if len(server.services) != 1 || len(server.tasks) != 3 || len(server.containers) != 3 {
		t.Fatalf("ServiceUpdate: wrong item count. Want 1 service and 3 replicas. Got services: %d, tasks: %d, containers: %d.", len(server.services), len(server.tasks), len(server.containers))
	}
}

func TestServiceUpdateNotFound(t *testing.T) {
	server, unused := setUpSwarm(t)
	defer server.Stop()
	defer unused.Stop()
	recorder := httptest.NewRecorder()
	updateOpts := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "test",
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: "test/test2",
				Args:  []string{"--test2"},
				Env:   []string{"ENV=2"},
				User:  "test",
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Mode: swarm.ResolutionModeVIP,
			Ports: []swarm.PortConfig{{
				Protocol:      swarm.PortConfigProtocolTCP,
				TargetPort:    uint32(80),
				PublishedPort: uint32(80),
			}},
		},
	}
	buf, err := json.Marshal(updateOpts)
	if err != nil {
		t.Fatalf("ServiceUpdate error: %s", err.Error())
	}
	request, _ := http.NewRequest("POST", "/services/pale/update", bytes.NewReader(buf))
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("ServiceUpdate: wrong status code. Want %d. Got %d.", http.StatusNotFound, recorder.Code)
	}
}

func TestNodeList(t *testing.T) {
	srv1, srv2 := setUpSwarm(t)
	defer srv1.Stop()
	defer srv2.Stop()
	for _, srv := range []*DockerServer{srv1, srv2} {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/nodes", nil)
		srv.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("invalid status code: %d", recorder.Code)
		}
		var nodes []swarm.Node
		err := json.NewDecoder(recorder.Body).Decode(&nodes)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(nodes, srv1.nodes) {
			t.Fatalf("expected nodes to equal %#v, got: %#v", srv1.nodes, nodes)
		}
		if !reflect.DeepEqual(nodes, srv2.nodes) {
			t.Fatalf("expected nodes to equal %#v, got: %#v", srv2.nodes, nodes)
		}
	}
}

func TestNodeInfo(t *testing.T) {
	srv1, srv2 := setUpSwarm(t)
	defer srv1.Stop()
	defer srv2.Stop()
	for _, srv := range []*DockerServer{srv1, srv2} {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/nodes/"+srv.nodes[0].ID, nil)
		srv.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("invalid status code: %d", recorder.Code)
		}
		var node swarm.Node
		err := json.NewDecoder(recorder.Body).Decode(&node)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(node, srv1.nodes[0]) {
			t.Fatalf("expected node to equal %#v, got: %#v", srv1.nodes[0], node)
		}
		if !reflect.DeepEqual(node, srv2.nodes[0]) {
			t.Fatalf("expected node to equal %#v, got: %#v", srv2.nodes[0], node)
		}
	}
}

func TestNodeUpdate(t *testing.T) {
	srv1, srv2 := setUpSwarm(t)
	defer srv1.Stop()
	defer srv2.Stop()
	recorder := httptest.NewRecorder()
	for i, srv := range []*DockerServer{srv1, srv2} {
		key := fmt.Sprintf("l%d", i)
		data, err := json.Marshal(swarm.NodeSpec{
			Annotations: swarm.Annotations{Labels: map[string]string{key: "value"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(data)
		request, _ := http.NewRequest("POST", "/nodes/"+srv.nodes[0].ID+"/update", body)
		srv.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("invalid status code: %d", recorder.Code)
		}
		if srv1.nodes[0].Spec.Annotations.Labels[key] != "value" {
			t.Fatalf("expected node to have label %s", key)
		}
		if srv2.nodes[0].Spec.Annotations.Labels[key] != "value" {
			t.Fatalf("expected node to have label %s", key)
		}
	}
}

func TestNodeDelete(t *testing.T) {
	srv1, srv2 := setUpSwarm(t)
	defer srv1.Stop()
	defer srv2.Stop()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/nodes/"+srv1.nodes[0].ID, nil)
	srv1.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("invalid status code: %d", recorder.Code)
	}
	if len(srv1.nodes) != 1 {
		t.Fatalf("expected len(nodes) to be 1, got %d", len(srv1.nodes))
	}
	if len(srv2.nodes) != 1 {
		t.Fatalf("expected len(nodes) to be 1, got %d", len(srv2.nodes))
	}
}

func setUpSwarm(t *testing.T) (*DockerServer, *DockerServer) {
	server1, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	server2, err := NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/swarm/init", bytes.NewReader(nil))
	server1.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("setUpSwarm: invalid status code swarm init %d", recorder.Code)
	}
	data, err := json.Marshal(swarm.JoinRequest{
		RemoteAddrs: []string{server1.SwarmAddress()},
	})
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/swarm/join", bytes.NewReader(data))
	server2.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("setUpSwarm: invalid status code swarm join %d", recorder.Code)
	}
	return server1, server2
}

func addTestService(server *DockerServer) (*swarm.Service, error) {
	recorder := httptest.NewRecorder()
	serviceCreateOpts := docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "test",
				Labels: map[string]string{
					"mykey": "myvalue",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: "test/test",
					Args:  []string{"--test"},
					Env:   []string{"ENV=1"},
					User:  "test",
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode: swarm.ResolutionModeVIP,
				Ports: []swarm.PortConfig{{
					Protocol:      swarm.PortConfigProtocolTCP,
					TargetPort:    uint32(80),
					PublishedPort: uint32(80),
				}},
			},
		},
	}
	buf, err := json.Marshal(serviceCreateOpts)
	if err != nil {
		return nil, err
	}
	var params io.Reader
	params = bytes.NewBuffer(buf)
	request, _ := http.NewRequest("POST", "/services/create", params)
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", recorder.Code)
	}
	if len(server.services) == 0 {
		return nil, fmt.Errorf("no service created on server")
	}
	if len(server.tasks) == 0 {
		return nil, fmt.Errorf("no tasks created on server")
	}
	return server.services[0], nil
}

func TestMutateTask(t *testing.T) {
	server := DockerServer{failures: make(map[string]string)}
	server.buildMuxer()
	server.tasks = append(server.tasks, &swarm.Task{ID: "id123"})
	newTask := swarm.Task{Status: swarm.TaskStatus{State: swarm.TaskStateFailed}}
	err := server.MutateTask("id123", newTask)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(server.tasks[0], &newTask) {
		t.Errorf("Wrong task after mutation.\nWant %#v.\nGot %#v.",
			newTask, server.tasks[0])
	}
}

func TestMutateTaskNotFound(t *testing.T) {
	server := DockerServer{failures: make(map[string]string)}
	server.buildMuxer()
	newTask := swarm.Task{Status: swarm.TaskStatus{State: swarm.TaskStateFailed}}
	err := server.MutateTask("id123", newTask)
	if err == nil {
		t.Error("Unexpected <nil> error")
	}
	if err.Error() != "task not found" {
		t.Errorf("wrong error message. Want %q. Got %q.", "task not found", err)
	}
}
