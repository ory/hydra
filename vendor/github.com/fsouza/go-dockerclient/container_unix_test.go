// +build !windows
// Copyright 2016 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

func TestExportContainerViaUnixSocket(t *testing.T) {
	content := "exported container tar content"
	var buf []byte
	out := bytes.NewBuffer(buf)
	tempSocket := tempfile("export_socket")
	defer os.Remove(tempSocket)
	endpoint := "unix://" + tempSocket
	u, _ := parseEndpoint(endpoint, false)
	client := Client{
		HTTPClient:             cleanhttp.DefaultClient(),
		Dialer:                 &net.Dialer{},
		endpoint:               endpoint,
		endpointURL:            u,
		SkipServerVersionCheck: true,
	}
	listening := make(chan string)
	done := make(chan int)
	containerID := "4fa6e0f0c678"
	go runStreamConnServer(t, "unix", tempSocket, listening, done, containerID)
	<-listening // wait for server to start
	opts := ExportContainerOptions{ID: containerID, OutputStream: out}
	err := client.ExportContainer(opts)
	<-done // make sure server stopped
	if err != nil {
		t.Errorf("ExportContainer: caugh error %#v while exporting container, expected nil", err.Error())
	}
	if out.String() != content {
		t.Errorf("ExportContainer: wrong stdout. Want %#v. Got %#v.", content, out.String())
	}
}

func TestStatsTimeoutUnixSocket(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "socket")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	socketPath := filepath.Join(tmpdir, "docker_test.sock")
	t.Logf("socketPath=%s", socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	received := make(chan bool)
	defer l.Close()
	go func() {
		conn, connErr := l.Accept()
		if connErr != nil {
			t.Logf("Failed to accept connection: %s", connErr)
			return
		}
		breader := bufio.NewReader(conn)
		req, connErr := http.ReadRequest(breader)
		if connErr != nil {
			t.Logf("Failed to read request: %s", connErr)
			return
		}
		if req.URL.Path != "/containers/c/stats" {
			t.Logf("Wrong URL path for stats: %q", req.URL.Path)
			return
		}
		received <- true
		time.Sleep(2 * time.Second)
	}()
	client, _ := NewClient("unix://" + socketPath)
	client.SkipServerVersionCheck = true
	errC := make(chan error, 1)
	statsC := make(chan *Stats)
	done := make(chan bool)
	defer close(done)
	go func() {
		errC <- client.Stats(StatsOptions{ID: "c", Stats: statsC, Stream: true, Done: done, Timeout: time.Millisecond})
		close(errC)
	}()
	err = <-errC
	e, ok := err.(net.Error)
	if !ok || !e.Timeout() {
		t.Errorf("Failed to receive timeout error, got %#v", err)
	}
	recvTimeout := 2 * time.Second
	select {
	case <-received:
		return
	case <-time.After(recvTimeout):
		t.Fatalf("Timeout waiting to receive message after %v", recvTimeout)
	}
}
