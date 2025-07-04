// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

// CORSRequestHeadersSafelist We add the safe list cors accept headers
// https://developer.mozilla.org/en-US/docs/Glossary/CORS-safelisted_request_header
var CORSRequestHeadersSafelist = []string{"Accept", "Content-Type", "Content-Length", "Accept-Language", "Content-Language"}

// CORSResponseHeadersSafelist We add the safe list cors expose headers
// https://developer.mozilla.org/en-US/docs/Glossary/CORS-safelisted_response_header
var CORSResponseHeadersSafelist = []string{"Set-Cookie", "Cache-Control", "Expires", "Last-Modified", "Pragma", "Content-Length", "Content-Language", "Content-Type"}

// CORSDefaultAllowedMethods Default allowed methods
var CORSDefaultAllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// CORSRequestHeadersExtended Extended list of request headers
// these will be concatenated with the safelist
var CORSRequestHeadersExtended = []string{"Authorization", "X-CSRF-TOKEN"}

// CORSResponseHeadersExtended Extended list of response headers
// these will be concatenated with the safelist
var CORSResponseHeadersExtended = []string{}

// CORSDefaultMaxAge max age for cache of preflight request result
// default is 5 seconds
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
var CORSDefaultMaxAge = 5

// CORSAllowCredentials default value for allow credentials
// this is required for cookies to be sent by the browser
// we always want this since we are using cookies for authentication most of the time
var CORSAllowCredentials = true
