// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import "strconv"

var toString = [600]string{
	100: "100", 101: "101", 102: "102", 103: "103",

	200: "200", 201: "201", 202: "202", 203: "203", 204: "204", 205: "205", 206: "206", 207: "207", 208: "208",
	226: "226",

	300: "300", 301: "301", 302: "302", 303: "303", 304: "304", 305: "305", 307: "307", 308: "308",
	400: "400", 401: "401", 402: "402", 403: "403", 404: "404", 405: "405", 406: "406", 407: "407", 408: "408", 409: "409", 410: "410", 411: "411", 412: "412", 413: "413", 414: "414", 415: "415", 416: "416", 417: "417", 418: "418", 421: "421", 422: "422", 423: "423", 424: "424", 425: "425", 426: "426", 428: "428", 429: "429",
	431: "431",
	451: "451",
	499: "499",

	500: "500", 501: "501", 502: "502", 503: "503", 504: "504", 505: "505", 506: "506", 507: "507", 508: "508", 510: "510", 511: "511",
}

// StatusCodeToString converts an HTTP status code to its string representation.
// It is highly optimized for the well-known HTTP status codes.
func StatusCodeToString(code int) string {
	if code >= 0 && code < len(toString) {
		if s := toString[code]; s != "" {
			return s
		}
	}
	return strconv.Itoa(code)
}
