# How do I know if OAuth2 / Hydra is the right choice for me?

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