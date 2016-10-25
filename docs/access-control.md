# Access Control Policies

Besides OAuth2 Token Introspection, Hydra offers Access Control Policies using
the [Ladon](https://github.com/ory-am/ladon) framework. Access Control Policies are used by Hydra internally and exposed
via various HTTP APIs.

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

To decide what the answer is, Hydra uses policy documents which can be represented as JSON. Values `actions`, `subjects`
and `resources` can use regular expressions by encapsulating the expression in `<>`, for example `<.*>`.

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

Now, Hydra is able to answer access requests like the following one:

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

In this case, the access request will be allowed:

1. `users:peter` matches `"subjects": ["users:<[peter|ken]>", "users:maria", "groups:admins"]`, as would `users:ken`, `users:maria` and `group:admins`.
2. `delete` matches `"actions" : ["delete", "<[create|update]>"]` as would `update` and `create`
3. `resource:articles:ladon-introduction` matches `"resources": ["resources:articles:<.*>", "resources:printer"],`
4. `"remoteIP": "192.168.0.5"` matches the [`CIDRCondition`](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
condition that was configured for the field `remoteIP`.

## Access Control Decisions: The Warden

The warden is a HTTP API allowing you to perform these access requests.
The warden knows two endpoints:

* `/warden/allowed`: Check if a subject is allowed to do something.
* `/warden/token/allowed`: Check if the subject of a token is allowed to do something.

Both endpoints use policies to compute the result and are documented in the [HTTP API Documentation](http://docs.hdyra.apiary.io/#reference/warden:-access-control).