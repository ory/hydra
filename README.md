# Hydra

[![Build Status](https://travis-ci.org/ory-am/hydra.svg)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)

![Hydra](hydra.png)

Hydra is a twelve factor authentication, authorization and account management service, ready for you to use in your micro service architecture.
Hydra is written in go and backed by PostgreSQL or any implementation of [account/storage.go](account/storage.go).

## What is Hydra?

Authentication, authorization and user account management are always lengthy to plan and implement. If you're building an app in Go with
a micro service architecture in mind, you have come the right way.

## Features

Hydra is a RESTful service providing you with things like:

* Account management: Sign up, settings, password recovery
* Authorization & Policy management backed by [ladon](https://github.com/ory-am/ladon)
* Authentication: Sign in, OAuth2 Consumer (Google Account, Facebook Login, ...)
* OAuth2 Provider: Hydra implements OAuth2 as specified at [rfc6749](http://tools.ietf.org/html/rfc6749) and [draft-ietf-oauth-v2-10](http://tools.ietf.org/html/draft-ietf-oauth-v2-10)

To keep things sane, Hydra uses JWT with the [ECDSA signature algorithm](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm)
for bearer based http authorization.

## Attributions

Image taken from [here](https://www.flickr.com/photos/pathfinderlinden/7161293044/).