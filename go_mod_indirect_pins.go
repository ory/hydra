//go:build go_mod_indirect_pins
// +build go_mod_indirect_pins

package main

import (
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/mikefarah/yq/v4"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/stringer"
	_ "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"

	_ "github.com/ory/go-acc"
)
