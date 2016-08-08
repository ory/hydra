// Package hydra is an api-only cloud native OAuth2 and OpenID Connect provider that integrates with existing authentication mechanisms:
//
// At first, there was the monolith. The monolith worked well with the bespoke authentication module. Then, the web evolved into an elastic cloud that serves thousands of different user agents in every part of the world.
// Hydra is driven by the need for a scalable, low-latency, in memory Access Control, OAuth2, and OpenID Connect layer that integrates with every identity provider you can imagine.
// Hydra is available through Docker and relies on RethinkDB for persistence. Database drivers are extensible in case you want to use RabbitMQ, MySQL, MongoDB, or some other database instead.
// Hydra is built for high throughput environments. Check out the below siege benchmark on a Macbook Pro Late 2013, connected to RethinkDB validating access tokens.
//
// The official repository is located at https://github.com/ory-am/hydra
package main