---
id: index
slug: /
title: Introduction
sidebar_label: Introduction
---

Hydra is an OAuth 2.0 and OpenID Connect Provider. In other words, an
implementation of the OAuth 2.0 Authorization Framework as well as the OpenID
Connect Core 1.0 framework. As such, it issues OAuth 2.0 Access, Refresh, and ID
Tokens that enable third-parties to access your APIs in the name of your users.

## Flexible User Management

One of ORY Hydra's biggest advantages is that unlike other OAuth 2.0
implementations, it implements the OAuth and OpenID Connect standard without
forcing you to use a "Hydra User Management" (login, logout, profile management,
registration), a particular template engine, or a predefined front-end.

This allows you to implement user management and login your way, in your
technology stack, with authentication mechanisms required by your use case
(token-based 2FA, SMS 2FA, etc). You can of course use existing solutions like
[authboss](https://github.com/go-authboss/authboss). It provides you all the benefits of OAuth 2.0
and OpenID Connect while being minimally invasive to your business logic and
technology stack.

## OpenID Certified

ORY Hydra is a
[Certified OpenID Connect Provider Server](https://openid.net/developers/certified/) and
implements all the requirements stated by the OpenID Foundation. In particular,
it correctly implements the various OAuth 2.0 and OpenID Connect flows specified
by the IETF and OpenID Foundation.

## Cryptographic Key Storage

In addition to the OAuth 2.0 functionality, ORY Hydra offers a safe storage for
cryptographic keys (used for example to sign JSON Web Tokens) and can manage
OAuth 2.0 Clients.

## Security First

ORY Hydra's architecture and work flows are designed to neutralize many common
(OWASP TOP TEN) and uncommon attack vectors.
[Learn more](./security-architecture.md).

## High Performance

Hydra has a low CPU and memory footprint, short start-up time, and scales
effortlessly up and down on many platforms including Heroku, Cloud Foundry,
Docker, Google Container Engine, and others.

## Developer Friendly

Hydra is available for all popular platforms including Linux, OSX and Windows.
It ships as a single binary without any additional dependencies. For further
simplicity, it is available as a
[Docker Image](https://hub.docker.com/r/oryd/hydra/).

Hydra also provides a developer-friendly CLI.

## Limitations

Hydra has a few limitations too:

1. ORY Hydra does not manage user accounts, i.e. user registration, password
   reset, user login, sending confirmation emails, etc. In Hydra's architecture,
   the _Identity Provider_ is responsible for this.
2. ORY Hydra doesn't support the OAuth 2.0 Resource Owner Password Credentials
   flow because it is legacy, discouraged, and insecure.

## Is ORY Hydra the right fit for you?

OAuth 2.0 can be used in many environments for various purposes. This list might
help you decide if OAuth 2.0 and Hydra are the right fit for a use case:

1. enable third-party solutions to access to your APIs: This is what an OAuth2
   Provider does, Hydra is a perfect fit.
2. be an Identity Provider like Google, Facebook, or Microsoft: OpenID Connect
   and thus Hydra is a perfect fit.
3. enable your browser, mobile, or wearable applications to access your APIs:
   Running an OAuth2 Provider can work great for this. You don't have to store
   passwords on the device and can revoke access tokens at any time. GMail
   logins work this way.
4. you want to limit what type of information your backend services can read
   from each other. For example, the _comment service_ should only be allowed to
   fetch user profile updates but shouldn't be able to read user passwords.
   OAuth 2.0 might make sense for you.

## Other solutions

If you only need a library or SDKs that implements OAuth 2.0, take a look at
[fosite](https://github.com/ory/fosite) or
[node-oauth2-server](https://github.com/oauthjs/node-oauth2-server).

If you need a fully featured identity solution including user management
and user interfaces, those exist in the cloud as [Ory](https://console.ory.sh) or when self-hosting as
[Keycloak](https://www.keycloak.org) or [Ory Kratos](https://github.com/ory/kratos/) among others.
