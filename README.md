# Hydra

[![Build Status](https://travis-ci.org/ory-am/hydra.svg)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)

![Hydra](hydra.png)

Hydra is a twelve factor authentication, authorization and account management service, ready for you to use in your micro service architecture.
Hydra is written in go and backed by PostgreSQL or any implementation of [account/storage.go](account/storage.go).

## What is Hydra?

Authentication, authorization and user account management are always lengthy to plan and implement. If you're building a micro service app
in need of these three, you are in the right place.

## Features

Hydra is a RESTful service providing you with things like:

* **Account Management**: Sign up, settings, password recovery
* **Access Control / Policy Management** backed by [ladon](https://github.com/ory-am/ladon)
* Hydra comes with a rich set of **OAuth2** features:
  * Hydra implements OAuth2 as specified at [rfc6749](http://tools.ietf.org/html/rfc6749) and [draft-ietf-oauth-v2-10](http://tools.ietf.org/html/draft-ietf-oauth-v2-10).
  * Hydra uses self-contained Acccess Tokens as suggessted in [rfc6794#section-1.4](http://tools.ietf.org/html/rfc6749#section-1.4) by issuing JSON Web Tokens as specified at
   [https://tools.ietf.org/html/rfc7519](https://tools.ietf.org/html/rfc7519) with [RSASSA-PKCS1-v1_5 with the SHA-256](https://tools.ietf.org/html/rfc7519#section-8) supported reducing overhead significantly. Access
  * Hydra implements **OAuth2 Introspection** as specified in [rfc7662](https://tools.ietf.org/html/rfc7662)

## Good to know

Hydra does not provide a dedicated login page, instead a third party service has to act as a login page
and authenticating the user via the [password grant type](https://aaronparecki.com/articles/2012/07/29/1/oauth2-simplified#others).

## Attributions

Image taken from [here](https://www.flickr.com/photos/pathfinderlinden/7161293044/).