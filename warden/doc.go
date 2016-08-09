// Package warden decides if access requests should be allowed or denied. In a scientific taxonomy, the warden
// is classified as a Policy Decision Point. THe warden's primary goal is to implement `github.com/ory-am/hydra/firewall.Firewall`.
//
// This package is structured as follows:
// * handler.go: A HTTP handler capable of validating access tokens.
// * warden_http.go: A Go API using HTTP to validate access tokens.
// * warden_local.go: A Go API using storage managers to validate access tokens.
// * warden_test.go: Functional tests all of the above.
package warden
