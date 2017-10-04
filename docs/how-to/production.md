# ORY Hydra in production

This guide aims to help setting up a production system with ORY Hydra.

## ORY Hydra behind a reverse proxy

Although ORY Hydra implements all Go best practices around running public-facing production http servers, we discourage running
ORY Hydra facing the public net directly. We strongly recommend running ORY Hydra behind an API gateway or a load balancer.
It is common to terminate TLS on the edge (gateway / load balancer) and use certificates provided by your infrastructure
provider (e.g. AWS CA) for last mile security.

### TLS Termination

You may also choose to set Hydra to HTTPS mode without actually accepting TLS connections. In that case,
all Hydra URLs are prefixed with `https://`, but the server is actually accepting http. This makes sense if you don't want
last mile security using TLS, and trust your network to properly handle internal traffic. To use this setting, check
for `HTTPS_ALLOW_TERMINATION_FROM` in `hydra help host`.

With TLS termination enabled, ORY Hydra discards all requests unless:
* The request is coming from a trusted IP address set by `HTTPS_ALLOW_TERMINATION_FROM` and the header `X-Forwarded-Proto` is set to `https`.
* The request goes to `/health/status` which does not require TLS termination and that is used to check the health of an instance.

### Routing

It is common to use a router, or API gateway, to route subdomains or paths to a specific service. For example, `https://myservice.com/hydra/`
is routed to `http://10.0.1.213:3912/` where `10.0.1.213` is the host running ORY Hydra. To compute the values for
the consent challenge, ORY Hydra uses the host and path headers from the HTTP request. Therefore, it is important
to set up your API Gateway in such a way, that it passes the public host (in this case `myservice.com`) and the path
without any prefix (in this case `hydra/`). If you use the Mashape Kong API gateway, you can achieve this by setting
`strip_request_path=true` and `preserve_host=true.`
