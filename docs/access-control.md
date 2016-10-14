# Access Control Policies

Besides OAuth2 Token Introspection, Hydra offers Access Control Policies using
the [Ladon](https://github.com/ory-am/ladon) framework. Access Control Policies are used by Hydra internally and exposed
via various HTTP APIs. It is important to understand how policies work so you can set up Hydra with the best possible configuration.
Apart form that, you do not need to use either Access Control Policies nor the Warden in your application.

## Introduction to Access Control Policies

Hydra's Access Control is able to answer the question:

> **Who** is **able** to do **what** on **something** given some **context**

* **Who** An arbitrary unique subject name, for example "ken" or "printer-service.mydomain.com".
* **Able**: The effect which can be either "allow" or "deny".
* **What**: An arbitrary action name, for example "delete", "create" or "scoped:action:something".
* **Something**: An arbitrary unique resource name, for example "something", "resources.articles.1234" or some uniform
    resource name like "urn:isbn:3827370191".
* **Context**: The current context containing information about the environment such as the IP Address,
    request date, the resource owner name, the department ken is working in or any other information you want to pass along.
    (optional)

To decide what the answer is, Hydra uses policy documents which can be represented as JSON

```json
{
  "description": "One policy to rule them all.",
  "subjects": ["users:<[peter|ken]>", "users:maria", "groups:admins"],
  "actions" : ["delete", "<[create|update]>"],
  "effect": "allow",
  "resources": [
    "resources:articles:<.*>",
    "resources:printer"
  ],
  "conditions": {
    "remoteIP": {
        "type": "CIDRCondition",
        "options": {
            "cidr": "192.168.0.1/16"
        }
    }
  }
}
```

and can answer access requests that look like:

```json
{
  "subject": "users:peter",
  "action" : "delete",
  "resource": "resource:articles:ladon-introduction",
  "context": {
    "remoteIP": "192.168.0.5"
  }
}
```

### HTTP Examples (tbd)

```
> curl \
      -X POST \
      -H "Content-Type: application/json" \
      -d@- \
      "https://my-ladon-implementation.localhost/policies" <<EOF
        {
          "description": "One policy to rule them all.",
          "subjects": ["users:<[peter|ken]>", "users:maria", "groups:admins"],
          "actions" : ["delete", "<[create|update]>"],
          "effect": "allow",
          "resources": [
            "resources:articles:<.*>",
            "resources:printer"
          ],
          "conditions": {
            "remoteIP": {
                "type": "CIDRCondition",
                "options": {
                    "cidr": "192.168.0.1/16"
                }
            }
          }
        }
  EOF
```

Then we test if "peter" (ip: "192.168.0.5") is allowed to "delete" the "ladon-introduction" article:

```
> curl \
      -X POST \
      -H "Content-Type: application/json" \
      -d@- \
      "https://my-ladon-implementation.localhost/warden" <<EOF
        {
          "subject": "users:peter",
          "action" : "delete",
          "resource": "resource:articles:ladon-introduction",
          "context": {
            "remoteIP": "192.168.0.5"
          }
        }
  EOF

{
    "allowed": true
}
```

## Access Control Endpoint: The Warden

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

## Warden

The Warden is usually called from your own services ("resource providers"), not from third parties. Hydra prevents
third parties from having access to these endpoints per default, but you can change that with custom policies.

The Warden endpoints are documented [here](http://docs.hdyra.apiary.io/#reference/warden:-access-control-for-resource-providers).