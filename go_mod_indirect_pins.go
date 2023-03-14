// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build go_mod_indirect_pins
// +build go_mod_indirect_pins

package main

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/mikefarah/yq/v4"
	_ "golang.org/x/tools/cmd/goimports"

	_ "github.com/ory/go-acc"
)
