// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build tools
// +build tools

package fosite

import (
	_ "github.com/gorilla/websocket"
	_ "github.com/mattn/goveralls"

	_ "github.com/ory/go-acc"
)
