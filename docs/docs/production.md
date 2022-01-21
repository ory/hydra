---
id: production
title: Preparing for Production
---

This document summarizes things you will find useful when going to production.

## ORY Hydra behind an API Gateway

Although ORY Hydra implements all Go best practices around running public-facing
production http servers, we discourage running ORY Hydra facing the public net
directly. We strongly recommend running ORY Hydra behind an API gateway or a
load balancer. It is common to terminate TLS on the edge (gateway / load
balancer) and use certificates provided by your infrastructure provider (e.g.
AWS CA) for last mile security.

### TLS Termination

You may also choose to set Hydra to HTTPS mode without actually accepting TLS
connections. In that case, all Hydra URLs are prefixed with `https://`, but the
server is actually accepting http. This makes sense if you don't want last mile
security using TLS, and trust your network to properly handle internal traffic:

```yaml
serve:
  tls:
    allow_termination_from:
      - 127.0.0.1/32
```

With TLS termination enabled, ORY Hydra discards all requests unless:

- The request is coming from a trusted IP address set by
  `serve.tls.allow_termination_from` and the header `X-Forwarded-Proto` is set
  to `https`.
- The request goes to `/health/alive`, `/health/ready` which does not require
  TLS termination and that is used to check the health of an instance.

When TLS Termination is enabled, you do not need to provide a TLS Certificate
and Private Key.

If you are unable to properly set up TLS Termination, you may want to set the
`--dangerous-force-http` flag. But please be aware that we discourage you from
doing so and that you should know what you're doing.

### Routing

It is common to use a router, or API gateway, to route subdomains or paths to a
specific service. For example, `https://myservice.com/hydra/` is routed to
`http://10.0.1.213:3912/` where `10.0.1.213` is the host running ORY Hydra. To
compute the values for the consent challenge, ORY Hydra uses the host and path
headers from the HTTP request. Therefore, it is important to set up your API
Gateway in such a way, that it passes the public host (in this case
`myservice.com`) and the path without any prefix (in this case `hydra/`). If you
use the Mashape Kong API gateway, you can achieve this by setting
`strip_request_path=true` and `preserve_host=true.`

## Exposing Administrative and Public API Endpoints

ORY Hydra serves APIs via two ports:

- Public port (default 4444)
- Administrative port (default 4445)

The public port can and should be exposed to public internet traffic. That port
handles requests to:

- `/.well-known/jwks.json`
- `/.well-known/openid-configuration`
- `/oauth2/auth`
- `/oauth2/token`
- `/oauth2/revoke`
- `/oauth2/fallbacks/consent`
- `/oauth2/fallbacks/error`
- `/oauth2/sessions/logout`
- `/userinfo`

The administrative port should not be exposed to public internet traffic. If you
want to expose certain endpoints, such as the `/clients` endpoint for OpenID
Connect Dynamic Client Registry, you can do so but you need to properly secure
these endpoints with an API Gateway or Authorization Proxy. Administrative
endpoints include:

- All `/clients` endpoints.
- All `/keys` endpoints.
- All `/health`, `/metrics`, `/version` endpoints.
- All `/oauth2/auth/requests` endpoints.
- Endpoint `/oauth2/introspect`.
- Endpoint `/oauth2/flush`.

None of the administrative endpoints have any built-in access control. You can
do simple `curl` or Postman requests to talk to them.

The Token Introspection endpoint requires authentication. But since there is no
access control, any valid authentication enables the endpoint to be used. If you
need to access this endpoint in production, you should configure your API
Gateway or Application Proxy to restrict which clients have access to the
endpoint.

We generally advise to run ORY Hydra with `hydra serve all` which listens on
both ports in one process. Please be aware that the `memory` backend will not
work in this mode.

### Binding to different interfaces or UNIX sockets

ORY Hydra will bind public and administrative APIs ports to all interfaces.

The interfaces or UNIX sockets used may be specified via environment variables
`PUBLIC_HOST` and `ADMIN_HOST`. Interfaces may be specified as TCP address or as
UNIX socket (giving the absolute path to the socket file prefixed by `unix:`)
like:

- `PUBLIC_HOST=127.0.0.1`
- `ADMIN_HOST="unix:/var/run/hydra/admin_socket"`

ORY Hydra will try to create the socket file during startup and the socket will
be writeable by the user running ORY Hydra. The owner, group and mode of the
socket can be modified:

```yaml
serve:
  admin:
    host: unix:/var/run/hydra/admin_socket
    socket:
      owner: hydra
      group: hydra-admin-api
      mode: 770
```

### Key generation and High Availability environments

Be aware that on the very first launch of the Hydra container(s), a worker
process will perform certain first-time installation tasks, such as generating
[JSON web keys](/hydra/docs/jwks) if they don't already exist.

If you intend on running your production Hydra environment in a highly-available
setup (for example, multiple concurrent containers behind a load-balancer), it's
possible that both containers will generate JWKs at the same time.

Although this isn't a problem, we recommend that you launch your production
environment with just one container to begin with, to complete the initial
seeding of the database.

Once done, you can raise your number of containers to achieve high availability.
