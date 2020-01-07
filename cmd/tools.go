// +build tools

package cmd

import (
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/gobuffalo/packr/packr"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/sdk/swagutil"
	_ "github.com/sqs/goreturns"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/stringer"

	_ "github.com/sqs/goreturns"

	_ "github.com/ory/go-acc"
	_ "github.com/ory/x/tools/listx"
)
