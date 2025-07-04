// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"fmt"
	"testing"

	"github.com/ory/pop/v6"
	"github.com/ory/pop/v6/logging"
)

func formatter(lvl logging.Level, s string, args ...interface{}) string {
	if pop.Debug == false {
		return ""
	}

	if lvl == logging.SQL {
		if len(args) > 0 {
			xargs := make([]string, len(args))
			for i, a := range args {
				switch a.(type) {
				case string:
					xargs[i] = fmt.Sprintf("%q", a)
				default:
					xargs[i] = fmt.Sprintf("%v", a)
				}
			}
			s = fmt.Sprintf("%s - %s | %s", lvl, s, xargs)
		} else {
			s = fmt.Sprintf("%s - %s", lvl, s)
		}
	} else {
		s = fmt.Sprintf(s, args...)
		s = fmt.Sprintf("%s - %s", lvl, s)
	}
	return s
}

func TestingLogger(t testing.TB) func(lvl logging.Level, s string, args ...interface{}) {
	return func(lvl logging.Level, s string, args ...interface{}) {
		if line := formatter(lvl, s, args...); len(line) > 0 {
			t.Log(line)
		}
	}
}

func NullLogger() func(lvl logging.Level, s string, args ...interface{}) {
	return func(lvl logging.Level, s string, args ...interface{}) {
		// do nothing
	}
}
