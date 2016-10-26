<h1 align="center" style="margin-bottom: 8px">
    <img alt="Hydra - cloud native access control" src="images/logo-essential.png">
</h1>

<h4 style="margin: 0 20% 1.275em" align="center">
Hydra is a runnable server implementation of the <a href="https://tools.ietf.org/html/rfc6749">OAuth 2.0 authorization framework</a>
and the <a href="http://openid.net/specs/openid-connect-core-1_0.html">OpenID Connect Core 1.0</a>.
</h4>

This is the user guide, you might be looking for the [GitHub Repository](http://github.com/ory-am/hydra), the
[HTTP REST API Documentation](http://docs.hdyra.apiary.io/) or
the [Language-Level API Documentation](https://godoc.org/github.com/ory-am/hydra) instead.

<hr>

# Introduction

Welcome to the Hydra documentation. This documentation will

1. teach you what OAuth2 and OpenID Connect are and how Hydra fits in the picture.
2. help you run a Hydra installation on your system using Docker.
3. teach you how to install, configure, run and use Hydra.
3. teach you how to hack and contribute back to Hydra.

Let us begin with the first part, understanding what OAuth2 and OpenID Connect are.

## Introduction to OAuth 2.0 and OpenID Connect

This section will give you some ideas of what OAuth 2.0 and OpenID Connect 1.0 are for. If you
already know what OAuth2 and OpenID Connect are and how they works, you can skip to the next [Section](#introduction-to-hydra).
This section will not explain how the various flows of OAuth2 work and how they look like. We strongly recommend
to read the following articles:

* [DigitalOcean: An Introduction to OAuth 2](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2)
* [Aaron Parecki: OAuth2 Simplified](https://aaronparecki.com/2012/07/29/2/oauth2-simplified)
* [Zapier: Chapter 5: Authentication, Part 2](https://zapier.com/learn/apis/chapter-5-authentication-part-2/)

### What is OAuth 2.0?

[The OAuth 2.0 authorization framework](https://tools.ietf.org/html/rfc6749) is a memo in the
[Request for Comments](https://www.ietf.org/rfc.html) document series published by the
IETF Internet Engineering Task Force (IETF). Memos in the Requests for Comments (RFC) document series
contain technical and organizational notes about the Internet. They cover many aspects of computer
networking, including protocols, procedures, programs, and concepts [...].

The OAuth 2.0 authorization framework enables a third-party
application to obtain limited access to an HTTP service, either on
behalf of a resource owner by orchestrating an approval interaction
between the resource owner and the HTTP service, or by allowing the
third-party application to obtain access on its own behalf.

In the traditional client-server authentication model, the client
requests an access-restricted resource (protected resource) on the
server by authenticating with the server using the resource owner's
credentials. In order to provide third-party applications access to
restricted resources, the resource owner shares its credentials with
the third party. This creates several problems and limitations.

OAuth addresses these issues by introducing an authorization layer
and separating the role of the client from that of the resource
owner.  In OAuth, the client requests access to resources controlled
by the resource owner and hosted by the resource server, and is
issued a different set of credentials than those of the resource
owner.

Instead of using the resource owner's credentials to access protected
resources, the client obtains an access token -- a string denoting a
specific scope, lifetime, and other access attributes.  Access tokens
are issued to third-party clients by an authorization server with the
approval of the resource owner.  The client uses the access token to
access the protected resources hosted by the resource server.

Source: [IETF RFC 6749](https://tools.ietf.org/html/rfc6749)

### OAuth 2.0 Example

An end-user (resource owner) can grant a printing
service (client) access to her protected photos stored at a photo-
sharing service (resource server), without sharing her username and
password with the printing service.  Instead, she authenticates
directly with a server trusted by the photo-sharing service
(authorization server), which issues the printing service delegation-
specific credentials (access token).

Source: [IETF RFC 6749](https://tools.ietf.org/html/rfc6749)

### What is OpenID Connect 1.0?

OpenID Connect 1.0 is a simple identity layer on top of the OAuth 2.0 protocol.
It enables Clients to verify the identity of the End-User based on the authentication performed
by an Authorization Server, as well as to obtain basic profile information about the End-User in
an interoperable and REST-like manner.

As background, the OAuth 2.0 Authorization Framework and OAuth 2.0 Bearer Token
Usage specifications provide a general framework for third-party
applications to obtain and use limited access to HTTP resources.
They define mechanisms to obtain and use Access Tokens to access resources
but do not define standard methods to provide identity information.
Notably, without profiling OAuth 2.0, it is incapable of providing information
about the authentication of an End-User.

OpenID Connect implements authentication as an extension to the OAuth 2.0 authorization process.

Source [OpenID Connect Core 1.0](openid.net/specs/openid-connect-core-1_0.html)


**OpenID Connect 1.0** is a simple identity layer on top of the OAuth 2.0 protocol.
It allows Clients to verify the identity of the End-User based on the authentication performed
by an Authorization Server, as well as to obtain basic profile information about the End-User in an
interoperable and REST-like manner.

OpenID Connect allows clients of all types, including Web-based, mobile, and JavaScript clients,
to request and receive information about authenticated sessions and end-users. The specification
suite is extensible, allowing participants to use optional features such as encryption of identity data,
discovery of OpenID Providers, and session management, when it makes sense for them.

There are different work flows for OpenID Connect 1.0, we recommend checking out the OpenID Connect sandbox at
[openidconnect.net](https://openidconnect.net/).

## Introduction to Hydra

Hydra is a server implementation of the OAuth 2.0 authorization framework and the OpenID Connect Core 1.0. Existing OAuth2
implementations usually ship as libraries or SDKs such as [node-oauth2-server](https://github.com/oauthjs/node-oauth2-server)
or [fosite](https://github.com/ory-am/fosite/issues), or as fully featured identity solutions with user
management and user interfaces, such as [Dex](https://github.com/coreos/dex) or [Okta](https://www.okta.com/).

Implementing and using OAuth2 without understanding the whole specification is challenging and prone to errors, even when
SDKs are being used. The primary goal of Hydra is to make OAuth 2.0 and OpenID Connect 1.0 less painful to set up and easier to use.

Hydra implements the flows described in OAuth2 and OpenID Connect 1.0 without forcing you to use a "Hydra User Management"
or some template engine or a predefined front-end. Instead it relies on HTTP redirection and cryptographic methods
to verify user consent allowing you to use Hydra with any authentication endpoint, be it [authboss](https://github.com/go-authboss/authboss),
[auth0.com](https://auth0.com/) or your proprietary PHP authentication.

Hydra incorporates best practices in the area of the web service technology:

1. Hydra ships as a single binary for all popular platforms including Linux, OSX and Windows, without any additional
dependencies. For further simplicity, Hydra is available as a [Docker Image](https://hub.docker.com/r/oryam/hydra/).
2. Hydra is built security first: architecture and work flows are designed to neutralize various common (OWASP TOP TEN)
and uncommon attack vectors. [Learn more](https://ory-am.gitbooks.io/hydra/content/basics/security.html).
3. Hydra has a low CPU and memory footprint, short start up times and a CLI with developers in mind.
4. Additionally, Hydra is capable of sophisticated access control, suitable for distributed systems or large organization. [Learn more](https://ory-am.gitbooks.io/hydra/content/access-control.html).
5. Hydra scales effortlessly up and down on every platform imaginable, including Heroku, Cloud Foundry, Docker,
Google Container Engine and many more.

Hydra has a limitations too:

1. Hydra is not something that manages user accounts. Hydra does not offer user registration, password reset, user
login, sending confirmation emails. This is what the *Identity Provider* is responsible for.
The communication between Hydra and the Identity Provider is called [*Consent Flow*](https://ory-am.gitbooks.io/hydra/content/oauth2/consent.html).
2. If you are building a simple service for 50-100 registered users, OAuth2 and Hydra will probably be too sophisticated.
3. Hydra does currently not support the OAuth2 resource owner password credentials flow. This will change in the future and is tracked
as issue [#214](https://github.com/ory-am/hydra/issues/214).
4. Hydra has no management frontend. You must manage OAuth2 Clients and other things using the RESTful endpoints or
the command line interface. We are open to having an official Hydra Management frontend.

OAuth2 is used in many areas, for various purposes and supported by all well known programming languages, but it is important
to understand what the vision of OAuth2 is. This non-exclusive list might help you decide, if OAuth 2.0 and Hydra are
the right fit for you.

1. If you want to allow third-party developers accessing your APIs now or in the future, Hydra is the perfect fit. This is what an OAuth2 Provider does.
2. If you want to become a Identity Provider, like Google, Facebook or Microsoft, OpenID Connect and thus Hydra is a perfect fit.
3. Running an OAuth2 Provider works great with browser, mobile and wearable apps, as you can avoid storing user
credentials on the device, phone or wearable and revoke access tokens, and thus access privileges, at any time. Adding
OAuth2 complexity to your environment when you never plan to do (1),
might not be worth it. Our advice: write a pros/cons list.
4. If you have a lot of services and want to limit automated access (think: cronjobs) for those services,
OAuth2 might make sense for you. Example: The comment service is not allowed to read user passwords when fetching
the latest user profile updates.

# OAuth 2.0 Case Study

OAuth2 and OpenID Connect are tricky to understand. It is important to understand that OAuth2 is
a delegation protocol. It makes sense to use Hydra in new and existing projects. A use case covering an existing project
explains how one would use Hydra in a new one as well. So let's look at a use case!

Let's assume we are running a ToDo List App (todo24.com). ToDo24 has a login endpoint (todo24.com/login).
The login endpoint is written in node and uses MongoDB to store user information (email + password + settings). Of course,
todo24 has other services as well: list management (todo24.com/lists/manage: close, create, move), item management (todo24.com/lists/items/manage: mark solved, add), and so on.
You are using cookies to see which user is performing the request.

Now you decide to use OAuth2 on top of your current infrastructure. There are many reasons to do this:
* You want to open up your APIs to third-party developers. Their apps will be using OAuth2 Access Tokens to access a user's to do list.
* You want a mobile client. Because you can not store secrets on devices (they can be reverse engineered and stolen), you use OAuth2 Access Tokens instead.
* You have Cross Origin Requests. Making cookies work with Cross Origin Requests weakens or even disables important anti-CSRF measures.
* You want to write an in-browser client. This is the same case as in a mobile client (you can't store secrets in a browser).

These are only a couple of reasons to use OAuth2. You might decide to use OAuth2 as your single source of authorization, thus maintaining
only one authorization protocol and being able to open up to third party devs in no time. With OpenID Connect, you are able to delegate authentication as well as authorization!

Your decision is final. You want to use OAuth2 and you want Hydra to do the job. You install Hydra in your cluster using docker.
Next, you set up some exemplary OAuth2 clients. Clients can act on their own, but most of the time they need to access a user's todo lists.
To do so, the client initiates an OAuth2 request. This is where [Hydra's authentication flow](https://ory-am.gitbooks.io/hydra/content/oauth2.html#authentication-flow) comes in to play.
Before Hydra can issue an access token, we need to know WHICH user is giving consent. To do so, Hydra redirects the user agent (e.g. browser, mobile device)
to the login endpoint alongside with a challenge that contains an expiry time and other information. The login endpoint (todo24.com/login) authenticates the
user as usual, e.g. by username & password, session cookie or other means. Upon successful authentication, the login endpoint asks for the user's consent:
*"Do you want to grant MyCoolAnalyticsApp read & write access to all your todo lists? [Yes] [No]"*. Once the user clicks *Yes* and gives consent,
the login endpoint redirects back to hydra and appends something called a *consent token*. The consent token is a cryptographically signed
string that contains information about the user, specifically the user's unique id. Hydra validates the signature's trustworthiness
and issues an OAuth2 access token and optionally a refresh or OpenID token.

Every time a request containing an access token hits a resource server (todo24.com/lists/manage), you make a request to Hydra asking who the token's
subject (the user who authorized the client to create a token on its behalf) is and whether the token is valid or not. You may optionally
ask if the token has permission to perform a certain action.