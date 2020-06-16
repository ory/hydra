// +build go_mod_indirect_pins

package main

import (
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/gobuffalo/packr/v2/packr2"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/sqs/goreturns"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/stringer"

	// FIXME pins websocket to 1.4.2
	// FIXME See https://github.com/gobuffalo/buffalo/pull/1999
	_ "github.com/gorilla/websocket"

	_ "github.com/ory/cli"

	_ "github.com/sqs/goreturns"

	_ "github.com/ory/go-acc"
	_ "github.com/ory/x/tools/listx"
)
