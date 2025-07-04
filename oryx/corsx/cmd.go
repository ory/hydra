// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

// HelpMessage returns a string containing information on setting up this CORS middleware.
func HelpMessage() string {
	return `- CORS_ENABLED: Switch CORS support on (true) or off (false). Default is off (false).

	Example: CORS_ENABLED=true

- CORS_ALLOWED_ORIGINS: A list of origins (comma separated values) a cross-domain request can be executed from.
	If the special * value is present in the list, all origins will be allowed. An origin may contain a wildcard (*)
	to replace 0 or more characters (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
	Only one wildcard can be used per origin. The default value is *.

	Example: CORS_ALLOWED_ORIGINS=http://*.domain.com,http://*.domain2.com

- CORS_ALLOWED_METHODS: A list of methods  (comma separated values) the client is allowed to use with cross-domain
	requests. Default value is simple methods (GET and POST).

	Example: CORS_ALLOWED_METHODS=POST,GET,PUT

- CORS_ALLOWED_CREDENTIALS: Indicates whether the request can include user credentials like cookies, HTTP authentication
	or client side SSL certificates.

	Default: CORS_ALLOWED_CREDENTIALS=false
	Example: CORS_ALLOWED_CREDENTIALS=true

- CORS_DEBUG: Debugging flag adds additional output to debug server side CORS issues.

	Default: CORS_DEBUG=false
	Example: CORS_DEBUG=true

- CORS_MAX_AGE: Indicates how long (in seconds) the results of a preflight request can be cached. The default is 0
	which stands for no max age.

	Default: CORS_MAX_AGE=0
	Example: CORS_MAX_AGE=10

- CORS_ALLOWED_HEADERS: A list of non simple headers (comma separated values) the client is allowed to use with
	cross-domain requests.

- CORS_EXPOSED_HEADERS: Indicates which headers (comma separated values) are safe to expose to the API of a
	CORS API specification.`
}
