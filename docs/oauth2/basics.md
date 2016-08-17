# OAuth2 Basics

[This introduction was taken from the Digital Ocean Blog.](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2)

**OAuth 2** is an authorization framework that enables applications to obtain limited access to user accounts on an
HTTP service, such as Facebook, GitHub, and Hydra. It works by delegating user authentication to the service that hosts
the user account, and authorizing third-party applications to access the user account. OAuth 2 provides authorization
flows for web and desktop applications, and mobile devices.

![](../dist/images/abstract_flow.png) 

Here is a more detailed explanation of the steps in the diagram:

* The application requests authorization to access service resources from the user
* If the user authorized the request, the application receives an authorization grant
* The application requests an access token from the authorization server (API) by presenting authentication of its own
identity, and the authorization grant.
* If the application identity is authenticated and the authorization grant is valid, the authorization server (API)
issues an access token to the application. Authorization is complete.
* The application requests the resource from the resource server (API) and presents the access token for authentication.
* If the access token is valid, the resource server (API) serves the resource to the application.

The actual flow of this process will differ depending on the authorization grant type in use, but this is the general idea.

Read more on OAuth2 on [the Digital Ocean Blog](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2).
We also recommend reading [API Security: Deep Dive into OAuth and OpenID Connect](http://nordicapis.com/api-security-oauth-openid-connect-depth/).

**Glossary**
* **The resource owner** is the user who authorizes an application to access their account. The application's access to
the user's account is limited to the "scope" of the authorization granted (e.g. read or write access).
* **Authorization Server (Hydra)** verifies the identity of the user and issues access tokens to the *client application*.
* **Client** is the *application* that wants to access the user's account. Before it may do so, it must be authorized
by the user.
* **Identity Provider** contains a log in user interface and a database of all your users. To integrate Hydra,
you must modify the Identity Provider. It mus be able to generate consent tokens and ask for the user's consent.
* **User Agent** is usually the resource owner's browser.
* **Consent Endpoint** is an app (e.g. NodeJS) that is able to receive consent challenges and create consent tokens.
It must verify the identity of the user that is giving the consent. This can be achieved using Cookie Auth,
HTTP Basic Auth, Login HTML Form, or any other mean of authentication. Upon authentication, the user must be asked
if he consents to allowing the client access to his resources.

## OAuth2 Clients

We already covered some basic OAuth2 concepts [in the Introduction](introduction.html).
You can manage *clients* using the cli or the HTTP REST API.

* **CLI:** `hydra clients -h`
* **REST:** Read the [API Docs](http://docs.hdyra.apiary.io/#reference/oauth2-clients)
