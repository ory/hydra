# What's it good for?

If you are new to OAuth2, this section is for you. To understand what OAuth2 and Hydra are, we will look at what they
**are not good for** first.

1. Hydra is not something that manages user accounts. Hydra does not offer user registration, password reset, user
login, sending confirmation emails. This is what the *Identity Provider* ("login endpoint") is responsible for.
The communication between Hydra and the Identity Provider is called [*Consent Flow*](https://ory-am.gitbooks.io/hydra/content/oauth2/consent.html).
[Auth0.com](https://auth0.com) is an Identity Provider. We might implement this feature at some point and if, it is going to be a different product.
2. If you think running an OAuth2 Provider can solve your user authentication ("log a user in"), Hydra is probably not for you. OAuth2 is a delegation protocol:

  > The OAuth 2.0 authorization framework enables a third-party application *[think: a dropbox app that manages your dropbox photos]*
  to obtain limited access to an HTTP service, either on behalf of *[do you allow "amazing photo app" to access all your photos?]*
  a resource owner *[user]* by orchestrating an approval interaction *[consent flow]* between the resource owner and the
  HTTP service, or by allowing the third-party application *[OAuth2 Client App]* to obtain access on its own behalf. \- *[IETF rfc6749](https://tools.ietf.org/html/rfc6749)*
3. If you are building a simple service for 50-100 registered users, OAuth2 and Hydra will be overkill.
4. Hydra does not support the OAuth2 resource owner password credentials flow.
5. Hydra has no user interface. You must manage OAuth2 Clients and other things using the RESTful endpoints.
A user interface is scheduled to accompany the stable release.

We use the following non-exclusive list to help people decide, if **OAuth2 is the right fit** for them.

1. If you want third-party developers to access your APIs, Hydra is the perfect fit. This is what an OAuth2 Provider does.
2. If you want to become a Identity Provider, like Google, Facebook or Microsoft, OpenID Connect and thus Hydra is a perfect fit.
3. Running an OAuth2 Provider works great with browser, mobile and wearable apps, as you can avoid storing user
credentials on the device, phone or wearable and revoke access tokens, and thus access privileges, at any time. Adding
OAuth2 complexity to your environment when you never plan to do (1),
might not be worth it. Our advice: write a pros/cons list.
4. If you have a lot of services and want to limit automated access (think: cronjobs) for those services,
OAuth2 might make sense for you. Example: The comment service is not allowed to read user passwords when fetching
the latest user profile updates.