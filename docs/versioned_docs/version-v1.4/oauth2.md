---
id: oauth2
title: OAuth 2.0 and OpenID Connect
sidebar_label: OAuth 2.0
---

This section describes on a high level what OAuth 2.0 and OpenID Connect 1.0 are
for and how they work.

### OAuth 2.0

[The OAuth 2.0 authorization framework](https://tools.ietf.org/html/rfc6749) is
specified in [IETF RFC 6749](https://tools.ietf.org/html/rfc6749). OAuth 2.0
enables a third-party application to obtain limited access to resources on an
HTTP server on behalf of the owner of those resources.

Why is this important? Without OAuth 2.0, a resource owner who wants to share
resources in their account with a third party would have to share their
credentials with that third party. As an example, let's say you (a resource
owner) have some photos (resources) stored on a social network (the resource
server). Now you want to print them using a third-party printing service. Before
OAuth 2.0 existed, you would have to enter your social network password into the
printing service so that it can access and print your photos. Sharing secret
passwords with third parties is obviously very problematic.

OAuth addresses this problem by introducing:

- the distinction between resource ownership and resource access for clients
- the ability to define fine-grained access privileges (called OAuth scopes)
  instead of full account access for third parties
- an authorization layer and workflow that allows resource owners to grant
  particular clients particular types of access to particular resources.

With OAuth, clients can request access to resources on a server, and the owner
of these resources can grant the requested access together with dedicated
credentials. In our example, you could grant the printing service read-only
access to your photos (only your photos, not your friend list) on the social
network. These credentials come in the form of an access token -- a string
denoting a specific scope, lifetime, and other access attributes. The client
(printing service) can use this access token to request the protected resources
(your photos) from the resource server (the social network).

### What is OpenID Connect 1.0?

OAuth 2.0 is a complex protocol for authorizing access to resources. If all you
need is authentication, OpenID Connect 1.0 enables clients to verify the
identity of the end user based on the authentication performed by an
Authorization Server and obtain basic profile information in an interoperable
and REST-like manner.

OpenID Connect allows clients of all types, including web and mobile, to receive
information about authenticated sessions and end users. The specification is
extensible, allowing participants to add encryption of identity data, discovery
of OpenID Providers, and session management as needed.

There are different work flows for OpenID Connect 1.0. We recommend checking out
the OpenID Connect sandbox at [openidconnect.net](https://openidconnect.net/).

A more detailed introduction of both OAuth 2.0 and OpenID Connect is available
in the following video:

<iframe width="560" height="315" src="https://www.youtube.com/embed/996OiexHze0" frameborder="0" allowfullscreen></iframe>

More details about the various OAuth2 flows can be found in these articles:

- [DigitalOcean: An Introduction to OAuth 2](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2)
- [Aaron Parecki: OAuth2 Simplified](https://aaronparecki.com/2012/07/29/2/oauth2-simplified)
- [Zapier: Chapter 5: Authentication, Part 2](https://zapier.com/learn/apis/chapter-5-authentication-part-2/)

1. ORY Hydra does not manage user accounts. Your application does that. Hydra
   exposes an OAuth 2.0 and OpenID Connect endpoint for the user accounts of
   your application.
2. OAuth Scopes are not access permissions to resources. They entitle an
   external application to act in the name of a user on a particular type of
   resource.

### OAuth 2.0 Scope != Permission

:::info The OAuth2 Scope reflects a permission the user gave to the OAuth2
Application, not a permission the system (e.g. API) gave to that OAuth2
application. Also, the OAuth2 Scope can not be changed without revoking the
token.

:::

A second important concept is the OAuth 2.0 Scope. Many people confuse OAuth 2.0
Scope with internal Access Control like for example Role Based Access Control
(RBAC) or Access Control Lists (ACL). Both concepts cover different aspects of
access control.

Internal access control (RBAC, ACL, etc) decides what a user can do in your
system. An administrator might be allowed to modify everything. A regular user
might only be allowed to read their own messages.

OAuth 2.0 Scopes, on the other hand, describe what a user allowed an external
application (OAuth 2.0 client) to do on his/her behalf. For example, an access
token might grant the external application to see a user's pictures, but not
upload new pictures on his/her behalf (which we assume a user could do herself).

In an extreme case, the user could lie and grant an external application OAuth
scopes that he himself doesn't have permission to ("read all classified
documents"). The OAuth Access Token with those scopes wouldn't help the external
application read those documents because it can only act in the name of the
user, and that user doesn't have these access privileges.

To recap, ORY Hydra's primary feature is implementing the OAuth 2.0 and OpenID
Connect spec, as well as related specs by the IETF and OpenID Foundation.

The next sections explain how to connect your existing user management (user
login, registration, logout, ...) with ORY Hydra in order to become an OAuth 2.0
and OpenID Connect provider like Google, Dropbox, or Facebook.

Again, please be aware that you must know how OAuth 2.0 and OpenID Connect work.
This documentation will not teach you how these protocols work.

### Terminology

To read more natural, this guide uses simpler terminologies like _user_ instead
of _resource owner_. Here is a full list of terms.

1. A **resource owner** is the user account who authorizes an external
   application to access their account. This access is limited (scoped) to
   particular actions (the granted "scopes" like read photos or write messages).
   This guide refers to resource owners as _users_ or _end users_.
2. The **OAuth 2.0 Authorization Server** implements the OAuth 2.0 protocol (and
   optionally OpenID Connect). In our case, this is **ORY Hydra**.
3. The **resource provider** is the service that hosts (provides) the resources.
   These resources (e.g. blog articles, printers, todo lists) are owned by a
   resource owner (user) mentioned above.
4. The **OAuth 2.0 Client** is the _external application_ that wants to access a
   resource owner's resources (read a user's images). To do that, it asks the
   OAuth 2.0 Authorization Server for an access token in a resource owner's
   behalf. The authorization server will ask the user if he/she "is ok with"
   giving that external application e.g. write access to personal images.
5. The **Identity Provider** is a service that allows users to register
   accounts, log in, etc.
6. **User Agent** is usually a browser.
7. **OpenID Connect** is a protocol built on top of OAuth 2.0 for just
   authentication (instead of authorizing access to resources).

A typical OAuth 2.0 flow looks as follows:

1. A developer registers an OAuth 2.0 Client (external application) with the
   Authorization Server (ORY Hydra) the intention to obtain information on
   behalf of a user.
2. The application UI asks the user to authorize the application to access
   information/data on his/her behalf.
3. The user is redirected to the Authorization Server.
4. The Authorization Server confirms the user's identity and asks the user to
   grant the OAuth 2.0 Client certain permissions.
5. The Authorization Server issues tokens that the OAuth 2.0 client uses to
   access resources on the user's behalf.

## Authenticating Users and Requesting Consent

As you already know by now, ORY Hydra does not come with any type of user
management (login, registration, ...). Instead, it relies on the so-called User
Login and Consent Flow. This flow describes a series of redirects where the
user's user agent is redirect to your Login Provider and, once the user is
authenticated, to the Consent Provider. The Login and Consent provider is
implemented by you in a programming language of your choice. You could write,
for example, a NodeJS app that handles HTTP requests to `/login` and `/consent`
and it would thus be your Login & Consent provider.

The flow itself works as follows:

1. The OAuth 2.0 Client initiates an Authorize Code, Hybrid, or Implicit flow.
   The user's user agent is redirected to
   `http://hydra/oauth2/auth?client_id=...&...`.
2. ORY Hydra, if unable to authenticate the user (meaning no session cookie
   exists), redirects the user's user agent to the Login Provider URL. The
   application "sitting" at that URL is implemented by you and typically shows a
   login user interface ("Please enter your username and password"). The URL the
   user is redirected to looks similar to
   `http://login-service/login?login_challenge=1234...`.
3. The Login Provider, once the user has successfully logged in, tells ORY Hydra
   some information about who the user is (e.g. the user's ID) and also that the
   login attempt was successful. This is done using a REST request which
   includes another redirect URL along the lines of
   `http://hydra/oauth2/auth?client_id=...&...&login_verifier=4321`.
4. The user's user agent follows the redirect and lands back at ORY Hydra. Next,
   ORY Hydra redirects the user's user agent to the Consent Provider, hosted
   at - for example - `http://consent-service/consent?consent_challenge=4567...`
5. The Consent Provider shows a user interface which asks the user if he/she
   would like to grant the OAuth 2.0 Client the requested permissions ("OAuth
   2.0 Scope"). You've probably seen this screen around, which is usually
   something similar to: _"Would you like to grant Facebook Image Backup access
   to all your private and public images?"_.
6. The Consent Provider makes another REST request to ORY Hydra to let it know
   which permissions the user authorized, and if the user authorized the request
   at all. The user can usually choose to not grant an application any access to
   his/her personal data. In the response of that REST request, a redirect URL
   is included along the lines of
   `http://hydra/oauth2/auth?client_id=...&...&consent_verifier=7654...`.
7. The user's user agent follows that redirect.
8. Now, the user has successfully authenticated and authorized the application.
   Next, ORY Hydra will run some checks and if everything works out, issue
   access, refresh, and ID tokens.

This flow allows you to take full control of the behaviour of your login system
(e.g. 2FA, passwordless, ...) and consent screen. A well-documented reference
implementation for both the Login and Consent Provider is
[available on GitHub](https://github.com/ory/hydra-login-consent-node).

### The flow from a user's point of view

<iframe width="560" height="315" src="https://www.youtube.com/embed/txUmfORzu8Y" frameborder="0" allowfullscreen></iframe>

### The flow from a network perspective

[![ORY Hydra Login and Consent Flow](https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgT0F1dGgyIENsaWVudC0-Pk9SWSBIeWRyYTogSW5pdGlhdGVzIE9BdXRoMiBBdXRob3JpemUgQ29kZSBvciBJbXBsaWNpdCBGbG93XG4gICAgT1JZIEh5ZHJhLS0-Pk9SWSBIeWRyYTogTm8gZW5kIHVzZXIgc2Vzc2lvbiBhdmFpbGFibGUgKG5vdCBhdXRoZW50aWNhdGVkKVxuICAgIE9SWSBIeWRyYS0-PkxvZ2luIFByb3ZpZGVyOiBSZWRpcmVjdHMgZW5kIHVzZXIgd2l0aCBsb2dpbiBjaGFsbGVuZ2VcbiAgICBMb2dpbiBQcm92aWRlci0tPk9SWSBIeWRyYTogRmV0Y2hlcyBsb2dpbiBpbmZvXG4gICAgTG9naW4gUHJvdmlkZXItLT4-TG9naW4gUHJvdmlkZXI6IEF1dGhlbnRpY2F0ZXMgdXNlciB3aXRoIGNyZWRlbnRpYWxzXG4gICAgTG9naW4gUHJvdmlkZXItLT5PUlkgSHlkcmE6IFRyYW5zbWl0cyBsb2dpbiBpbmZvIGFuZCByZWNlaXZlcyByZWRpcmVjdCB1cmwgd2l0aCBsb2dpbiB2ZXJpZmllclxuICAgIExvZ2luIFByb3ZpZGVyLT4-T1JZIEh5ZHJhOiBSZWRpcmVjdHMgZW5kIHVzZXIgdG8gcmVkaXJlY3QgdXJsIHdpdGggbG9naW4gdmVyaWZpZXJcbiAgICBPUlkgSHlkcmEtLT4-T1JZIEh5ZHJhOiBGaXJzdCB0aW1lIHRoYXQgY2xpZW50IGFza3MgdXNlciBmb3IgcGVybWlzc2lvbnNcbiAgICBPUlkgSHlkcmEtPj5Db25zZW50IFByb3ZpZGVyOiBSZWRpcmVjdHMgZW5kIHVzZXIgd2l0aCBjb25zZW50IGNoYWxsZW5nZVxuICAgIENvbnNlbnQgUHJvdmlkZXItLT5PUlkgSHlkcmE6IEZldGNoZXMgY29uc2VudCBpbmZvICh3aGljaCB1c2VyLCB3aGF0IGFwcCwgd2hhdCBzY29wZXMpXG4gICAgQ29uc2VudCBQcm92aWRlci0tPj5Db25zZW50IFByb3ZpZGVyOiBBc2tzIGZvciBlbmQgdXNlcidzIHBlcm1pc3Npb24gdG8gZ3JhbnQgYXBwbGljYXRpb24gYWNjZXNzXG4gICAgQ29uc2VudCBQcm92aWRlci0tPk9SWSBIeWRyYTogVHJhbnNtaXRzIGNvbnNlbnQgcmVzdWx0IGFuZCByZWNlaXZlcyByZWRpcmVjdCB1cmwgd2l0aCBjb25zZW50IHZlcmlmaWVyXG4gICAgQ29uc2VudCBQcm92aWRlci0-Pk9SWSBIeWRyYTogUmVkaXJlY3RzIHRvIHJlZGlyZWN0IHVybCB3aXRoIGNvbnNlbnQgdmVyaWZpZXJcbiAgICBPUlkgSHlkcmEtLT4-T1JZIEh5ZHJhOiBWZXJpZmllcyBncmFudFxuICAgIE9SWSBIeWRyYS0-Pk9BdXRoMiBDbGllbnQ6IFRyYW5zbWl0cyBhdXRob3JpemF0aW9uIGNvZGUvdG9rZW4iLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9fQ)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgT0F1dGgyIENsaWVudC0-Pk9SWSBIeWRyYTogSW5pdGlhdGVzIE9BdXRoMiBBdXRob3JpemUgQ29kZSBvciBJbXBsaWNpdCBGbG93XG4gICAgT1JZIEh5ZHJhLS0-Pk9SWSBIeWRyYTogTm8gZW5kIHVzZXIgc2Vzc2lvbiBhdmFpbGFibGUgKG5vdCBhdXRoZW50aWNhdGVkKVxuICAgIE9SWSBIeWRyYS0-PkxvZ2luIFByb3ZpZGVyOiBSZWRpcmVjdHMgZW5kIHVzZXIgd2l0aCBsb2dpbiBjaGFsbGVuZ2VcbiAgICBMb2dpbiBQcm92aWRlci0tPk9SWSBIeWRyYTogRmV0Y2hlcyBsb2dpbiBpbmZvXG4gICAgTG9naW4gUHJvdmlkZXItLT4-TG9naW4gUHJvdmlkZXI6IEF1dGhlbnRpY2F0ZXMgdXNlciB3aXRoIGNyZWRlbnRpYWxzXG4gICAgTG9naW4gUHJvdmlkZXItLT5PUlkgSHlkcmE6IFRyYW5zbWl0cyBsb2dpbiBpbmZvIGFuZCByZWNlaXZlcyByZWRpcmVjdCB1cmwgd2l0aCBsb2dpbiB2ZXJpZmllclxuICAgIExvZ2luIFByb3ZpZGVyLT4-T1JZIEh5ZHJhOiBSZWRpcmVjdHMgZW5kIHVzZXIgdG8gcmVkaXJlY3QgdXJsIHdpdGggbG9naW4gdmVyaWZpZXJcbiAgICBPUlkgSHlkcmEtLT4-T1JZIEh5ZHJhOiBGaXJzdCB0aW1lIHRoYXQgY2xpZW50IGFza3MgdXNlciBmb3IgcGVybWlzc2lvbnNcbiAgICBPUlkgSHlkcmEtPj5Db25zZW50IFByb3ZpZGVyOiBSZWRpcmVjdHMgZW5kIHVzZXIgd2l0aCBjb25zZW50IGNoYWxsZW5nZVxuICAgIENvbnNlbnQgUHJvdmlkZXItLT5PUlkgSHlkcmE6IEZldGNoZXMgY29uc2VudCBpbmZvICh3aGljaCB1c2VyLCB3aGF0IGFwcCwgd2hhdCBzY29wZXMpXG4gICAgQ29uc2VudCBQcm92aWRlci0tPj5Db25zZW50IFByb3ZpZGVyOiBBc2tzIGZvciBlbmQgdXNlcidzIHBlcm1pc3Npb24gdG8gZ3JhbnQgYXBwbGljYXRpb24gYWNjZXNzXG4gICAgQ29uc2VudCBQcm92aWRlci0tPk9SWSBIeWRyYTogVHJhbnNtaXRzIGNvbnNlbnQgcmVzdWx0IGFuZCByZWNlaXZlcyByZWRpcmVjdCB1cmwgd2l0aCBjb25zZW50IHZlcmlmaWVyXG4gICAgQ29uc2VudCBQcm92aWRlci0-Pk9SWSBIeWRyYTogUmVkaXJlY3RzIHRvIHJlZGlyZWN0IHVybCB3aXRoIGNvbnNlbnQgdmVyaWZpZXJcbiAgICBPUlkgSHlkcmEtLT4-T1JZIEh5ZHJhOiBWZXJpZmllcyBncmFudFxuICAgIE9SWSBIeWRyYS0-Pk9BdXRoMiBDbGllbnQ6IFRyYW5zbWl0cyBhdXRob3JpemF0aW9uIGNvZGUvdG9rZW4iLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9fQ)

### Implementing a Login & Consent Provider

You should now have a high-level idea of how the login and consent providers
work. Let's get into the details of it.

#### OAuth 2.0 Authorize Code Flow

Before anything happens, the OAuth 2.0 Authorize Code Flow is initiated by an
OAuth 2.0 Client. This usually works by generating a URL in the form of
`https://hydra/oauth2/auth?client_id=1234&scope=foo+bar&response_type=code&...`.
Then, the OAuth 2.0 Client points the end user's user agent to that URL.

Next, the user agent (browser) opens that URL.

#### User Login

As the user agent hits the URL, ORY Hydra checks if a session cookie is set
containing information about a previously successful login. Additionally,
parameters such as `id_token_hint`, `prompt`, and `max_age` are evaluated and
processed.

Next, the user will be redirect to the Login Provider which was set using the
`URLS_LOGIN` environment variable. For example, the user is redirected to
`https://login-provider/login?login_challenge=1234` if
`URLS_LOGIN=https://login-provider/login`. This redirection happens _always_ and
regardless of whether the user has a valid login session or if the user needs to
authenticate.

The service which handles requests to `https://login-provider/login` must first
fetch information on the authentication request using a REST API call. Please be
aware that for reasons of brevity, the following code snippets are pseudo-code.
For a fully working example, check out our reference
[User Login & Consent Provider implementation](https://github.com/ory/hydra-login-consent-node).

The endpoint handler at `/login` **must not remember previous sessions**. This
task is solved by ORY Hydra. If the REST API call tells you to show the login
ui, you **must show it**. If the REST API tells you to not show the login ui,
**you must not show it**. Again, **do not implement any type of session here**.

```
// This is node-js pseudo code and will not work if you copy it 1:1

router.get('/login', function (req, res, next) {
    challenge = req.url.query.login_challenge;

    fetch('https://hydra/oauth2/auth/requests/login?' + querystring.stringify({ login_challenge: challenge })).
        then(function (response) {
            return response.json()
        }).
        then(function (response) {
            // ...
        })
})
```

The server response is a JSON object with the following keys:

```
{
    // Skip, if true, let's us know that ORY Hydra has successfully authenticated the user and we should not show any UI
    "skip": true|false,

    // The user-id of the already authenticated user - only set if skip is true
    "subject": "user-id",

    // The OAuth 2.0 client that initiated the request
    "client": {"id": "...", ...},

    // The initial OAuth 2.0 request url
    "request_url": "https://hydra/oauth2/auth?client_id=1234&scope=foo+bar&response_type=code&...",

    // The OAuth 2.0 Scope requested by the client,
    "requested_scope": ["foo", "bar"],

    // Information on the OpenID Connect request - only required to process if your UI should support these values.
    "oidc_context": {"ui_locales": [...], ...},

	// Context is an optional object which can hold arbitrary data. The data will be made available when fetching the
	// consent request under the "context" field. This is useful in scenarios where login and consent endpoints share
	// data.
    "context": {...}
}
```

For a full documentation on all available keys, please head over to the
[API documentation](https://www.ory.sh/docs/api/hydra/) (make sure to select the
right API version).

Depending of whether or not `skip` is true, you will prompt the user to log in
by showing him/her a username/password form, or by using some other proof of
identity.

If `skip` is true, you **should not** show a user interface but accept the login
request directly by making a REST call. You can use this step to update some
internal count of how often a user logged in, or do some other custom business
logic. But again, do not show the user interface.

To accept the login request, do something along the lines of:

```
// This is node-js pseudo code and will not work if you copy it 1:1

const body = {
    // This is the user ID of the user that authenticated. If `skip` is true, this must be the `subject`
    // value from the `fetch('https://hydra/oauth2/auth/requests/login?' + querystring.stringify({ login_challenge: challenge }))` response:
    //
    // subject = response.subject
    //
    // Otherwise, this can be a value of your choosing:
    subject: "...",

    // If remember is set to true, then the authentication session will be persisted in the user's browser by ORY Hydra. This will set the `skip` flag to true in future requests that are coming from this user. This value has no effect if `skip` was true.
    remember: true|false,

    // The time (in seconds) that the cookie should be valid for. Only has an effect if `remember` is true.
    remember_for: 3600,

    // This value is specified by OpenID connect and optional - it tells OpenID Connect which level of authentication the user performed - for example 2FA or using some biometric data. The concrete values are up to you here.
    acr: ".."
}

fetch('https://hydra/oauth2/auth/requests/login/accept?' + querystring.stringify({ login_challenge: challenge }), {
    method: 'PUT',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' }
}).
    then(function (response) {
        return response.json()
    }).
    then(function (response) {
        // The response will contain a `redirect_to` key which contains the URL where the user's user agent must be redirected to next.
        res.redirect(response.redirect_to);
    })
```

You may also choose to deny the login request. This is possible regardless of
the `skip` value.

```
// This is node-js pseudo code and will not work if you copy it 1:1

const body = {
    error: "...", // This is an error ID like `login_required` or `invalid_request`
    error_description: "..." // This is a more detailed description of the error
}

fetch('https://hydra/oauth2/auth/requests/login/reject?' + querystring.stringify({ login_challenge: challenge }), {
    method: 'PUT',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' }
}).
    then(function (response) {
        return response.json()
    }).
    then(function (response) {
        // The response will contain a `redirect_to` key which contains the URL where the user's user agent must be redirected to next.
        res.redirect(response.redirect_to);
    })
```

#### User Consent

Now that we know who the user is, we must ask the user if he/she wants to grant
the requested permissions to the OAuth 2.0 Client. To do so, we check if the
user has previously granted that exact OAuth 2.0 Client the requested
permissions. If the user has never granted any permissions to the client, or the
client requires new permissions not previously granted, the user must visually
confirm the request.

This works very similar to the User Login Flow. First, the user will be
redirected to the Consent Provider which was set using the `URLS_CONSENT`
environment variable. For example, the user is redirected to
`https://consent-provider/consent?consent_challenge=1234` if
`URLS_CONSENT=https://consent-provider/consent`. This redirection happens
_always_ and regardless of whether the user has a valid login session or if the
user needs to authorize the application or not.

The service which handles requests to `https://consent-provider/consent` must
first fetch information on the consent request using a REST API call. Please be
aware that for reasons of brevity, the following code snippets are pseudo-code.
For a fully working example, check out our reference
[User Login, Logout & Consent Provider implementation](https://github.com/ory/hydra-login-consent-node).

```
// This is node-js pseudo code and will not work if you copy it 1:1

challenge = req.url.query.consent_challenge;

fetch('https://hydra/oauth2/auth/requests/consent?' + querystring.stringify({ consent_challenge: challenge })).
    then(function (response) {
        return response.json()
    }).
    then(function (response) {
        // ...
    })
```

The server response is a JSON object with the following keys:

```
{
    // Skip, if true, let's us know that the client has previously been granted the requested permissions (scope) by the end-user
    "skip": true|false,

    // The user-id of the user that will grant (or deny) the request
    "subject": "user-id",

    // The OAuth 2.0 client that initiated the request
    "client": {"id": "...", ...},

    // The initial OAuth 2.0 request url
    "request_url": "https://hydra/oauth2/auth?client_id=1234&scope=foo+bar&response_type=code&...",

    // The OAuth 2.0 Scope requested by the client.
    "requested_scope": ["foo", "bar"],

    // Contains the access token audience as requested by the OAuth 2.0 Client.
    requested_access_token_audience: ["foo", "bar"]

    // Information on the OpenID Connect request - only required to process if your UI should support these values.
    "oidc_context": {"ui_locales": [...], ...},

    // Contains arbitrary information set by the login endpoint or is empty if not set.
    "context": {...}
}
```

If skip is true, you should not show any user interface to the user. Instead,
you should accept (or deny) the consent request. Typically, you will accept the
request unless you have a very good reason to deny it (e.g. the OAuth 2.0 Client
is banned).

If skip is false and you show the consent screen, you should use the
`requested_scope` array to display a list of permissions which the user must
grant (e.g. using a checkbox). Some people choose to always skip this step if
the OAuth 2.0 Client is a first-party client - meaning that the client is used
by you or your developers in an internal application.

Assuming the user accepts the consent request, the code looks very familiar to
the User Login Flow.

```
// This is node-js pseudo code and will not work if you copy it 1:1

const body = {
    // A list of permissions the user granted to the OAuth 2.0 Client. This can be fewer permissions that initially requested, but are rarely more or other permissions than requested.
    grant_scope: ["foo", "bar"],

	// Sets the audience the user authorized the client to use. Should be a subset of `requested_access_token_audience`.
	grant_access_token_audience: ["foo", "bar"],

    // If remember is set to true, then the consent response will be remembered for future requests. This will set the `skip` flag to true in future requests that are coming from this user for the granted permissions and that particular client. This value has no effect if `skip` was true.
    remember: true|false,

    // The time (in seconds) that the cookie should be valid for. Only has an effect if `remember` is true.
    remember_for: 3600,

    // The session allows you to set additional data in the access and ID tokens.
    session: {
        // Sets session data for the access and refresh token, as well as any future tokens issued by the
        // refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection.
        // If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties
        // can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care!
        access_token: { ... },

        // Sets session data for the OpenID Connect ID token. Keep in mind that the session'id payloads are readable
        // by anyone that has access to the ID Challenge. Use with care! Any information added here will be mirrored at
        // the `/userinfo` endpoint.
        id_token: { ... },
    }
}

fetch('https://hydra/oauth2/auth/requests/consent/accept?' + querystring.stringify({ consent_challenge: challenge }), {
    method: 'PUT',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' }
}).
    then(function (response) {
        return response.json()
    }).
    then(function (response) {
        // The response will contain a `redirect_to` key which contains the URL where the user's user agent must be redirected to next.
        res.redirect(response.redirect_to);
    })
```

You may also choose to deny the consent request. This is possible regardless of
the `skip` value.

```
// This is node-js pseudo code and will not work if you copy it 1:1

const body = {
    // This is an error ID like `consent_required` or `invalid_request`
    error: "...",

    // This is a more detailed description of the error
    error_description: "..."
}

fetch('https://hydra/oauth2/auth/requests/consent/reject?' + querystring.stringify({ consent_challenge: challenge }), {
    method: 'PUT',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' }
}).
    then(function (response) {
        return response.json()
    }).
    then(function (response) {
        // The response will contain a `redirect_to` key which contains the URL where the user's user agent must be redirected to next.
        res.redirect(response.redirect_to);
    })
```

Once the user agent is redirected back, the OAuth 2.0 flow will be finalized.

### Logout

ORY Hydra supports:

- [OpenID Connect Front-Channel Logout 1.0](https://openid.net/specs/openid-connect-frontchannel-1_0.html)
- [OpenID Connect Back-Channel Logout 1.0](https://openid.net/specs/openid-connect-backchannel-1_0.html)

A log out request may be initiated by a RP (Relying Party, alias for OAuth 2.0
Client) or by calling the logout endpoint without any parameters. In both cases,
the high level flow looks as follows:

1. A user-agent (browser) requests the logout endpoint
   (`/oauth2/sessions/logout`). If the request is done on behalf of a RP:

- The URL query MUST contain an ID Token issued by ORY Hydra as the
  `id_token_hint`: `/oauth2/sessions/logout?id_token_hint=...`
- The URL query MAY contain key `post_logout_redirect_uri` indicating where the
  user agent should be redirected after the logout completed successfully. Each
  OAuth 2.0 Client can whitelist a list of URIs that can be used as the value
  using the `post_logout_redirect_uris` metadata field:
  `/oauth2/sessions/logout?id_token_hint=...&post_logout_redirect_uri=https://i-must-be-whitelisted/`
- If `post_logout_redirect_uri` is set, the URL query SHOULD contain a `state`
  value. On successful redirection, this state value will be appended to the
  `post_logout_redirect_uri`. The functionality is equal to the `state`
  parameter when performing OAuth2 flows.

2. The user-agent is redirected to the logout provider URL (configuration item
   `urls.logout`) and contains a challenge:
   `https://my-logout-provider/logout?challenge=...`
3. The logout provider uses the `challenge` query parameter to fetch metadata
   about the request. The logout provider may choose to show a UI where the user
   has to accept the logout request. Alternatively, the logout provider MAY
   choose to silently accept the logout request.
4. To accept the logout request, the logout provider makes a `PUT` call to
   `/oauth2/auth/requests/logout/accept?challenge=...`. No request body is
   required.
5. The response contains a `redirect_to` value where the logout provider
   redirects the user back to.
6. ORY Hydra performs OpenID Connect Front- and Back-Channel logout.
7. The user agent is being redirected to a specified redirect URL. This may
   either be the default redirect URL set by `urls.post_logout_redirect` or to
   the value specified by query parameter `post_logout_redirect_uri`.

[![ORY Hydra Lgoout flow](https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgVXNlciBBZ2VudC0-Pk9SWSBIeWRyYTogQ2FsbHMgbG9nb3V0IGVuZHBvaW50XG4gICAgT1JZIEh5ZHJhLS0-Pk9SWSBIeWRyYTogVmFsaWRhdGVzIGxvZ291dCBlbmRwb2ludFxuICAgIE9SWSBIeWRyYS0-PkxvZ291dCBQcm92aWRlcjogUmVkaXJlY3RzIGVuZCB1c2VyIHdpdGggbG9nb3V0IGNoYWxsZW5nZVxuICAgIExvZ291dCBQcm92aWRlci0tPk9SWSBIeWRyYTogRmV0Y2hlcyBsb2dvdXQgcmVxdWVzdCBpbmZvXG4gICAgTG9nb3V0IFByb3ZpZGVyLS0-PkxvZ291dCBQcm92aWRlcjogQWNxdWlyZXMgdXNlciBjb25zZW50IGZvciBsb2dvdXQgKG9wdGlvbmFsKVxuICAgIExvZ291dCBQcm92aWRlci0tPk9SWSBIeWRyYTogSW5mb3JtcyB0aGF0IGxvZ291dCByZXF1ZXN0IGlzIGdyYW50ZWRcbiAgICBMb2dvdXQgUHJvdmlkZXItPj5PUlkgSHlkcmE6IFJlZGlyZWN0cyBlbmQgdXNlciB0byByZWRpcmVjdCB1cmwgd2l0aCBsb2dvdXQgY2hhbGxlbmdlXG4gICAgT1JZIEh5ZHJhLS0-Pk9SWSBIeWRyYTogUGVyZm9ybXMgbG9nb3V0IHJvdXRpbmVzXG4gICAgT1JZIEh5ZHJhLS0-VXNlciBBZ2VudDogUmVkaXJlY3RzIHRvIHNwZWNpZmllZCByZWRpcmVjdCB1cmwiLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9fQ)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgVXNlciBBZ2VudC0-Pk9SWSBIeWRyYTogQ2FsbHMgbG9nb3V0IGVuZHBvaW50XG4gICAgT1JZIEh5ZHJhLS0-Pk9SWSBIeWRyYTogVmFsaWRhdGVzIGxvZ291dCBlbmRwb2ludFxuICAgIE9SWSBIeWRyYS0-PkxvZ291dCBQcm92aWRlcjogUmVkaXJlY3RzIGVuZCB1c2VyIHdpdGggbG9nb3V0IGNoYWxsZW5nZVxuICAgIExvZ291dCBQcm92aWRlci0tPk9SWSBIeWRyYTogRmV0Y2hlcyBsb2dvdXQgcmVxdWVzdCBpbmZvXG4gICAgTG9nb3V0IFByb3ZpZGVyLS0-PkxvZ291dCBQcm92aWRlcjogQWNxdWlyZXMgdXNlciBjb25zZW50IGZvciBsb2dvdXQgKG9wdGlvbmFsKVxuICAgIExvZ291dCBQcm92aWRlci0tPk9SWSBIeWRyYTogSW5mb3JtcyB0aGF0IGxvZ291dCByZXF1ZXN0IGlzIGdyYW50ZWRcbiAgICBMb2dvdXQgUHJvdmlkZXItPj5PUlkgSHlkcmE6IFJlZGlyZWN0cyBlbmQgdXNlciB0byByZWRpcmVjdCB1cmwgd2l0aCBsb2dvdXQgY2hhbGxlbmdlXG4gICAgT1JZIEh5ZHJhLS0-Pk9SWSBIeWRyYTogUGVyZm9ybXMgbG9nb3V0IHJvdXRpbmVzXG4gICAgT1JZIEh5ZHJhLS0-VXNlciBBZ2VudDogUmVkaXJlY3RzIHRvIHNwZWNpZmllZCByZWRpcmVjdCB1cmwiLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9fQ)

This endpoint does not remove any Access/Refresh Tokens.

#### Logout Provider Example (NodeJS Pseudo-code)

Following step 1 from the flow above, the user-agent is redirected to the logout
provider (e.g. `https://my-logout-provider/logout?challenge=...`). Next, the
logout provider fetches information about the logout request:

```node
// This is node-js pseudo code and will not work if you copy it 1:1

challenge = req.url.query.logout_challenge;

fetch(
  'https://hydra/oauth2/auth/requests/logout?' +
    querystring.stringify({ logout_challenge: challenge })
)
  .then(function (response) {
    return response.json();
  })
  .then(function (response) {
    // ...
  });
```

The server response is a JSON object with the following keys:

```
{
    // The user for whom the logout was request.
    "subject": "user-id",

    // The login session ID that was requested to log out.
    "sid": "abc..",

    // The original request URL.
    "request_url": "https://hydra/oauth2/sessions/logout?id_token_hint=...",

    // True if the request was initiated by a Relying Party (RP) / OAuth 2.0 Client. False otherwise.
    "rp_initiated": true|false
}
```

Next, the logout provider should decide if the end-user should perform a UI
action such as confirming the logout request. It is RECOMMENDED to request
logout confirmation from the end-user when `rp_initiated` is set to true.

When the logout provider decides to accept the logout request, the flow is
completed as follows:

```node
fetch(
  'https://hydra/oauth2/auth/requests/logout/accept?' +
    querystring.stringify({ logout_challenge: challenge }),
  {
    method: 'PUT',
  }
)
  .then(function (response) {
    return response.json();
  })
  .then(function (response) {
    // The response will contain a `redirect_to` key which contains the URL where the user's user agent must be redirected to next.
    res.redirect(response.redirect_to);
  });
```

You can also reject a logout request (e.g. if the user chose to not log out):

```node
fetch(
  'https://hydra/oauth2/auth/requests/logout/reject?' +
    querystring.stringify({ logout_challenge: challenge }),
  {
    method: 'PUT',
  }
).then(function (response) {
  // Now you can do whatever you want - redirect the user back to your home page or whatever comes to mind.
});
```

If the logout request was granted and the user agent redirected back to ORY
Hydra, all OpenID Connect Front-/Back-channel logout flows (if set) will be
performed and the user will be redirect back to his/her final destination.

#### [OpenID Connect Front-Channel Logout 1.0](https://openid.net/specs/openid-connect-frontchannel-1_0.html)

In summary
([read the spec](https://openid.net/specs/openid-connect-frontchannel-1_0.html))
this feature allows an OAuth 2.0 Client to register fields
`frontchannel_logout_uri` and `frontchannel_logout_session_required`.

If `frontchannel_logout_uri` is set to a valid URL (the host, port, path must
all match those of one of the Redirect URIs), ORY Hydra will redirect the
user-agent (typically browser) to that URL after a logout occurred. This allows
the OAuth 2.0 Client Application to log out the end-user in its own system as
well - for example by deleting a Cookie or otherwise invalidating the user
session.

ORY Hydra always appends query parameters values `iss` and `sid` to the
Front-Channel Logout URI, for example:

```
https://rp.example.org/frontchannel_logout
  ?iss=https://server.example.com
  &sid=08a5019c-17e1-4977-8f42-65a12843ea02
```

Each OpenID Connect ID Token is issued with a `sid` claim that will match the
`sid` value from the Front-Channel Logout URI.

ORY Hydra will automatically execute the required HTTP Redirects to make this
work. No extra work is required.

#### [OpenID Connect Back-Channel Logout 1.0](https://openid.net/specs/openid-connect-backchannel-1_0.html)

In summary
([read the spec](https://openid.net/specs/openid-connect-backchannel-1_0.html))
this feature allows an OAuth 2.0 Client to register fields
`backchannel_logout_uri` and `backchannel_logout_session_required`.

If `backchannel_logout_uri` is set to a valid URL, a HTTP Post request with
Content-Type `application/x-www-form-urlencoded` and a `logout_token` will be
made to that URL when a end-user logs out. The `logout_token` is a JWT signed
with the same key that is used to sign OpenID Connect ID Tokens. You should thus
validate the `logout_token` using the ID Token Public Key (can be fetched from
`/.well-known/jwks.json`). The `logout_token` contains the following claims:

- `iss` - Issuer Identifier, as specified in Section 2 of [OpenID.Core].
- `aud` - Audience(s), as specified in Section 2 of [OpenID.Core].
- `iat` - Issued at time, as specified in Section 2 of [OpenID.Core].
- `jti` - Unique identifier for the token, as specified in Section 9 of
  [OpenID.Core].
- `events` - Claim whose value is a JSON object containing the member name
  http://schemas.openid.net/event/backchannel-logout. This declares that the JWT
  is a Logout Token. The corresponding member value MUST be a JSON object and
  SHOULD be the empty JSON object {}.
- `sid` - Session ID - String identifier for a Session. This represents a
  Session of a User Agent or device for a logged-in End-User at an RP. Different
  sid values are used to identify distinct sessions at an OP. The sid value need
  only be unique in the context of a particular issuer. Its contents are opaque
  to the RP. Its syntax is the same as an OAuth 2.0 Client Identifier.

```
{
  "iss": "https://server.example.com",
  "aud": "s6BhdRkqt3",
  "iat": 1471566154,
  "jti": "bWJq",
  "sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
  "events": {
     "http://schemas.openid.net/event/backchannel-logout": {}
   }
}
```

An exemplary Back-Channel Logout Request looks as follows:

```
POST /backchannel_logout HTTP/1.1
Host: rp.example.org
Content-Type: application/x-www-form-urlencoded

logout_token=eyJhbGci ... .eyJpc3Mi ... .T3BlbklE ...
```

The Logout Token must be validated as follows:

- Validate the Logout Token signature in the same way that an ID Token signature
  is validated, with the following refinements.
- Validate the iss, aud, and iat Claims in the same way they are validated in ID
  Tokens.
- Verify that the Logout Token contains a sid Claim.
- Verify that the Logout Token contains an events Claim whose value is JSON
  object containing the member name
  http://schemas.openid.net/event/backchannel-logout.
- Verify that the Logout Token does not contain a nonce Claim.
- Optionally verify that another Logout Token with the same jti value has not
  been recently received.

The endpoint then returns a HTTP 200 OK response. Cache-Control headers should
be set to:

```
Cache-Control: no-cache, no-store
Pragma: no-cache
```

Because the OpenID Connect Back-Channel Logout Flow is not executed using the
user-agent (e.g. Browser) but from ORY Hydra directly, the session cookie of the
end-user will not be available to the OAuth 2.0 Client and the session has to be
invalidated by some other means (e.g. by blacklisting the session ID).

### Revoking consent and login sessions

#### Login

You can revoke login sessions. Revoking a login session will remove all of the
user's cookies at ORY Hydra and will require the user to re-authenticate when
performing the next OAuth 2.0 Authorize Code Flow. Be aware that this option
will remove all cookies from all devices.

Revoking the login sessions of a user is as easy as sending
`DELETE to`/oauth2/auth/sessions/login?subject={subject}`.

This endpoint is not compatible with OpenID Connect Front-/Backchannel logout
and does not revoke any tokens.

#### Consent

You can revoke a user's consent either on a per application basis or for all
applications. Revoking the consent will automatically revoke all related access
and refresh tokens.

Revoking all consent sessions of a user is as easy as sending
`DELETE to`/oauth2/auth/sessions/consent?subject={subject}`.

Revoking the consent sessions of a user for a specific client is as easy as
sending
`DELETE to`/oauth2/auth/sessions/consent?subject={subject}&client={client}`.

## OAuth 2.0

### OAuth 2.0 Scope

The scope of an OAuth 2.0 scope defines the permission the token was granted by
the end-user. For example, a specific token might be allowed to access public
pictures, but not private ones. The granted permissions are established during
the consent screen.

Additionally, ORY Hydra has pre-defined OAuth 2.0 Scope values:

- `offline_access`: Include this scope if you wish to receive a refresh token
- `openid`: Include this scope if you wish to perform an OpenID Connect request.

When performing an OAuth 2.0 Flow where the end-user is involved (e.g. Implicit
or Authorize Code), the granted OAuth 2.0 Scope must be set when accepting the
consent using the `grant_scope` key.

> A OAuth 2.0 Scope **is not a permission**:
>
> - A permission allows an actor to perform a certain action in a system: _Bob
>   is allowed to delete his own photos_.
> - OAuth 2.0 Scope implies that an end-user granted certain privileges to a
>   client: _Bob allowed the OAuth 2.0 Client to delete all users_.
>
> The OAuth 2.0 Scope can be granted without the end-user actually having the
> right permissions. In the examples above, Bob granted an OAuth 2.0 Client the
> permission ("scope") to delete all users in his name. However, since Bob is
> not an administrator, that permission ("access control") is not actually
> granted to Bob. Therefore any request by the OAuth 2.0 Client that tries to
> delete users on behalf of Bob should fail.

### OAuth 2.0 Access Token Audience

The Audience of an Access Token refers to the Resource Servers that this token
is intended for. The audience usually refers to one or more URLs such as

- `https://api.mydomain.com/blog/posts`
- `https://api.mydomain.com/users`

but may also refer to a subset of resources:

- `https://api.mydomain.com/tenants/foo/users`

When performing an OAuth 2.0 Flow where the end-user is involved (e.g. Implicit
or Authorize Code), the granted audience must be set when accepting the consent
using the `grant_access_token_audience` key. In most cases, it is ok to grant
the audience without user-interaction.

### OAuth 2.0 Refresh Tokens

OAuth 2.0 Refresh Tokens are issued only when an Authorize Code Flow
(`response_type=code`) or an OpenID Connect Hybrid Flow with an Authorize Code
Response Type (`response_type=code+...`) is executed. OAuth 2.0 Refresh Tokens
are not returned for Implicit or Client Credentials grants:

Capable of issuing an OAuth 2.0 Refresh Token:

- https://ory-hydra.example/oauth2/auth?response_type=code&...
- https://ory-hydra.example/oauth2/auth?response_type=code+token&...
- https://ory-hydra.example/oauth2/auth?response_type=code+token+id_token&...
- https://ory-hydra.example/oauth2/auth?response_type=code+id_token&...

Will not issue an OAuth 2.0 Refresh Token:

- https://ory-hydra.example/oauth2/auth?response_type=token&...
- https://ory-hydra.example/oauth2/auth?response_type=token+id_token&...
- https://ory-hydra.example/oauth2/token?grant_type=client_redentials&...

Additionally, each OAuth 2.0 Client that wants to request an OAuth 2.0 Refresh
Token must be allowed to request scope `offline_access`. When performing an
OAuth 2.0 Authorize Code Flow, the `offline_access` value must be included in
the requested OAuth 2.0 Scope:

```
https://authorization-server.com/auth
 &scope=offline_access
 ?response_type=code
 &client_id=...
 &redirect_uri=...
 &state=...
```

When accepting the consent request, `offline_access` must be in the list of
`grant_scope`:

```
fetch('https://hydra/oauth2/auth/requests/consent/accept?challenge=' + encodeURIComponent(challenge), {
    method: 'PUT',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' }
}).
const body = {
    grant_scope: ["offline_access"],
}
```

Refresh Token Lifespan can be set using configuration key `ttl.refresh_token`.
If set to -1, Refresh Tokens never expire.

### OAuth 2.0 Token Introspection

OAuth2 Token Introspection is an [IETF](https://tools.ietf.org/html/rfc7662)
standard. It defines a method for a protected resource to query an OAuth 2.0
authorization server to determine the active state of an OAuth 2.0 token and to
determine meta-information about this token. OAuth 2.0 deployments can use this
method to convey information about the authorization context of the token from
the authorization server to the protected resource.

The usage of an access token or client credentials is required to access the
endpoint. ORY Hydra will however accept any valid token or valid credentials as
there is no built-in access control.

You can find more details on this endpoint in the
[ORY Hydra API Docs](https://www.ory.sh/docs/). You can also use the CLI command
`hydra token introspect <token>`.

### OAuth 2.0 Clients

You can manage _OAuth 2.0 clients_ using the cli or the HTTP REST API:

- **CLI:** `hydra help clients`
- **REST:** Read the [API Docs](https://www.ory.sh/docs/hydra/sdk/api)

#### Examples

This section provides a few examples to get you started with the most-used OAuth
2.0 Clients:

##### Authorize Code Flow with Refresh Token

The following command creates an OAuth 2.0 Client capable of executing the
Authorize Code Flow, requesting ID and Refresh Tokens and performing the OAuth
2.0 Refresh Grant:

```sh
hydra clients create \
    --endpoint http://ory-hydra:4445 \
    --id client-id \
    --secret client-secret \
    --grant-types authorization_code,refresh_token \
    --response-types code \
    --scope openid,offline \
    --callbacks http://my-app.com/callback,http://my-other-app.com/callback
```

The OAuth 2.0 Client will be allowed to use values `http://my-app.com/callback`
and `http://my-other-app.com/callback` as `redirect_url`.

> It is expected that the OAuth 2.0 Client sends its credentials using HTTP
> Basic Authorization.

If you wish to send credentials in the POST Body, add the following flag to the
command above:

```
    --token-endpoint-auth-method client_secret_post \
```

The same can be achieved by setting
`"token_endpoint_auth_method": "client_secret_post"` in the request body of
`POST /clients` and `PUT /clients/<id>`.

##### Client Credentials Flow

A client only capable of performing the Client Credentials Flow can be created
as follows:

```
hydra clients create \
    --endpoint http://ory-hydra:4445 \
    --id my-client \
    --secret secret \
    -g client_credentials
```

## OpenID Connect

### Userinfo

The `/userinfo` endpoint returns information on a user given an access token.
Since ORY Hydra is agnostic to any end-user data, the `/userinfo` endpoint
returns only minimal information per default:

```
GET https://ory-hydra:4444/userinfo
Authorization: bearer access-token.xxxx

{
 "acr": "oauth2",
 "sub": "xxx@xxx.com"
}
```

Any information set to the key `session.id_token` during accepting the consent
request will also be included here.

```js

// This is node-js pseudo code and will not work if you copy it 1:1

const body = {
    // grant_scope: ["foo", "bar"],
    // ...
    session: {
        id_token: {
            "foo": "bar"
        },
    }
}

fetch('https://hydra/oauth2/auth/requests/consent/' + challenge + '/accept', {
    method: 'PUT',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' }
}).
    // then(function (response) {
```

By making the `/userinfo` call with a token issued by this consent request, one
would receive:

```
GET https://ory-hydra:4444/userinfo
Authorization: bearer new-access-token.xxxx

{
 "acr": "oauth2",
 "sub": "xxx@xxx.com",
 "foo": "bar"
}
```

You should only include data that has been authorized by the end-user through an
OAuth 2.0 Scope. If an OAuth 2.0 Client, for example, requests the `phone` scope
and the end-user authorizes that scope, the phone number should be added to
`session.id_token`.

> Be aware that the `/userinfo` endpoint is public. Its contents are thus as
> publicly visible as those of ID Tokens. It is therefore imperative to **not
> expose sensitive information without end-user consent.**
