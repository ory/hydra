---
id: case-study
title: OAuth 2.0 Case Study
---

OAuth2 and OpenID Connect are tricky to understand. It is important to keep in
mind that OAuth2 is a delegation protocol. Let's look at a use case to
understand how Hydra makes sense in new and existing projects.

Let's assume we are running todo24.com, a ToDo list app. ToDo24 has a login
endpoint (todo24.com/login). The login endpoint is written in Node.JS and uses
MongoDB to store user information (email + password + user profile). Of course,
ToDo24 has other services as well: list management (todo24.com/lists: create,
rename, close lists), item management (todo24.com/lists/{list-id}/items: add or
mark an item as solved), and so on. You are using cookie-based user sessions.

Now you decide to use OAuth2 on top of your current infrastructure. There are
many reasons to do this:

- You want to open your APIs to third-party developers. Their apps will be using
  OAuth2 Access Tokens to access your users todo lists.
- You want to build more client applications like a web app, mobile app,
  chat-bot, etc.
- You have cross-origin requests. Making cookies work with cross-origin requests
  weakens or even disables important anti-CSRF measures.

These are only a couple of reasons to use OAuth2. You might decide to use OAuth2
as your only authorization workflow, thus minimizing maintenance overhead while
always being able to support third party applications. With OpenID Connect, you
can delegate authentication as well!

So you decide to implement OAuth2 and use Ory Hydra to do the job. You run Hydra
by adding its Docker image to your cluster. Next, you set up some exemplary
OAuth2 clients. These clients need to access a user's todo lists. To do so, the
client initiates an OAuth2 request. This is where Hydra's
[user login & consent flow](concepts/oauth2.mdx) comes into play. Before Hydra
can issue an access token, we need to know which user is giving consent. To
determine this, Hydra redirects the user agent (browser, mobile device) to
ToDo24's login endpoint alongside with a challenge that contains important
request information. The login endpoint (todo24.com/login) authenticates the
user as usual, for example by username & password, session cookie or other
means. Upon successful authentication, the login endpoint redirects the user
back to Ory Hydra. Next, Ory Hydra needs the user's consent. It redirects the
user agent to the consent endpoint (todo24.com/consent) where the user is asked
something like _"Do you want to grant MyAnalyticsApp read access to your todo
lists? [Yes][no]"_. Once the user gives consent by clicking _Yes_, the consent
endpoint redirects back to Ory Hydra. Hydra validates the request and finally
issues the access, refresh, and ID tokens.

You can validate the access tokens which are sent to your API directly at Ory
Hydra, or use an Identity & Access Proxy like Ory Oathkeeper to do it for you.
