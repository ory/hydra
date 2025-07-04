// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ory/x/watcherx"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "expected 1 comand line argument but got %d\n", len(os.Args)-1)
		os.Exit(1)
	}
	c := make(chan watcherx.Event)
	ctx, cancel := context.WithCancel(context.Background())
	_, err := watcherx.WatchFile(ctx, os.Args[1], c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not initialize file watcher: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("watching file %s\n", os.Args[1])
	defer cancel()
	for {
		switch e := (<-c).(type) {
		case *watcherx.ChangeEvent:
			var data []byte
			data, err = io.ReadAll(e.Reader())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "could not read data: %+v\n", err)
				os.Exit(1)
			}
			fmt.Printf("got change event:\nData: %s,\nSrc: %s\n", data, e.Source())
		case *watcherx.RemoveEvent:
			fmt.Printf("got remove event:\nSrc: %s\n", e.Source())
		case *watcherx.ErrorEvent:
			fmt.Printf("got error event:\nError: %s\n", e.Error())
		default:
			fmt.Println("got unknown event")
		}
	}
}
