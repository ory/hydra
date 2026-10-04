// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

func sanitizeMethod(m string) string {
	switch m {
	case "GET", "get":
		return "get"
	case "PUT", "put":
		return "put"
	case "HEAD", "head":
		return "head"
	case "POST", "post":
		return "post"
	case "DELETE", "delete":
		return "delete"
	case "CONNECT", "connect":
		return "connect"
	case "OPTIONS", "options":
		return "options"
	case "NOTIFY", "notify":
		return "notify"
	case "TRACE", "trace":
		return "trace"
	case "PATCH", "patch":
		return "patch"
	default:
		return "unknown"
	}
}
