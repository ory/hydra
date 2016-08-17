# Access Control

Hydra offers various access control methods. Resource providers (e.g. photo/user/asset/balance/... service) use

1. **Warden Token Validation** to validate access tokens
2. **Warden Access Control with Access Tokens** to validate access tokens and decide
if the token's subject is allowed to perform the request
3. **Warden Access Control without Access Tokens** to decide if any subject is allowed
to perform a request

whereas third party apps (think of a facebook app) use

1. **OAuth2 Token Introspection** to validate access tokens.

There are two common ways to solve access control in a distributed environment (e.g. microservices).

1. Your services are behind a gateway (e.g. access control, rate limiting, and load balancer) 
that does the access control for them. This is known as a "trusted network/subnet".
2. Clients (e.g. Browser) talk to your services
directly. The services are responsible for checking access privileges themselves.

In both cases, you would use on of the warden endpoints.
