// Copyright 2013 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker_test

import (
	"archive/tar"
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func ExampleClient_AttachToContainer() {
	client, err := docker.NewClient("http://localhost:4243")
	if err != nil {
		log.Fatal(err)
	}
	client.SkipServerVersionCheck = true
	// Reading logs from container a84849 and sending them to buf.
	var buf bytes.Buffer
	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    "a84849",
		OutputStream: &buf,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(buf.String())
	buf.Reset()
	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    "a84849",
		OutputStream: &buf,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(buf.String())
}

func ExampleClient_BuildImage() {
	client, err := docker.NewClient("http://localhost:4243")
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write([]byte("FROM base\n"))
	tr.Close()
	opts := docker.BuildImageOptions{
		Name:         "test",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := client.BuildImage(opts); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_AddEventListener() {
	client, err := docker.NewClient("http://localhost:4243")
	if err != nil {
		log.Fatal(err)
	}

	listener := make(chan *docker.APIEvents)
	err = client.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {

		err = client.RemoveEventListener(listener)
		if err != nil {
			log.Fatal(err)
		}

	}()

	timeout := time.After(1 * time.Second)

	for {
		select {
		case msg := <-listener:
			log.Println(msg)
		case <-timeout:
			return
		}
	}

}

func ExampleEnv_Map() {
	e := docker.Env([]string{"A=1", "B=2", "C=3"})
	envs := e.Map()
	for k, v := range envs {
		fmt.Printf("%s=%q\n", k, v)
	}
}

func ExampleEnv_SetJSON() {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Gopher", Age: 4}
	var e docker.Env
	err := e.SetJSON("person", p)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleEnv_GetJSON() {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Gopher", Age: 4}
	var e docker.Env
	e.Set("person", `{"name":"Gopher","age":4}`)
	err := e.GetJSON("person", &p)
	if err != nil {
		log.Fatal(err)
	}
}
