# Access Control

OAuth 2.0 is a protocol that allows applications to act on a user's behalf,
without knowing their credentials. While it offers a secure flow for applications
(web apps, mobile apps, touchpoints, IoT) to gain access to your APIs,
it does not specify what access rights users actually have.
The frequent question "is the user actually allowed to access that resource?"
is not answered by OAuth 2.0, nor by OpenID Connect.

For that reason, ORY Hydra offers something we call Access Control Policies.
If you ever worked with Google Cloud IAM or AWS IAM, you probably know
what these policies look like. Access Control Policies are a powerful tool
capable of modeling simple and complex access control environments, such as
simple read/write APIs or complex multi-tenant environments.

The next sections give you an overview and best practices of these principles.

## Access Control Policies

### Introduction

Hydra's Access Control is able to answer the question:

> **Who** is **able** to do **what** on **something** given some **context**

* **Who** An arbitrary unique subject name, for example "ken" or "printer-service.mydomain.com".
* **Able**: The effect which can be either "allow" or "deny".
* **What**: An arbitrary action name, for example "delete", "create" or "scoped:action:something".
* **Something**: An arbitrary unique resource name, for example "something", "resources.articles.1234" or some uniform
    resource name like "urn:isbn:3827370191".
* **Context**: The current context containing information about the environment such as the IP Address, the time or date
    of access, or some other type of context. (optional)

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

### Best Practices

This sections gives an overview of best practices for access control policies
we developed over the years at ORY.

#### URNs

> “There are only two hard things in Computer Science: cache invalidation and naming things.”
-- Phil Karlton

URN naming is as hard as naming API endpoints. Thankfully, by doing the latter, the former is usually solved as well.
We will explore further best practices in the following sections.

##### Scope the Organization Name

A rule of thumb is to prefix resource names with a domain that represents the organization creating the software.

* **Do not:** `<some-id>`
* **Do:** `<organizaion-id>:<some-id>`

##### Scope Actions, Resources and Subjects

It is wise to scope actions, resources, and subjects in order to prevent name collisions:

* **Do not:** `myorg.com:<subject-id>`, `myorg.com:<resource-id>`, `myorg.com:<action-id>`
* **Do:** `myorg.com:subjects:<subject-id>`, `myorg.com:resources:<resource-id>`, `myorg.com:actions:<action-id>`
* **Do:** `subjects:myorg.com:<subject-id>`, `resources:myorg.com:<resource-id>`, `actions:myorg.com:<action-id>`

##### Multi-Tenant Systems

Multi-tenant systems typically have resources which should not be access by other tenants in the system. This can be
achieved by adding the tenant id to the URN:

* **Do:** `resources:myorg.com:tenants:<tenant-id>:<resource-id>`

In some environments, it is common to have organizations and projects belonging to those organizations. Here, the
following URN semantics can be used:

* **Do:** `resources:myorg.com:organizations:<organization-id>:projects:<project-id>:<resource-id>`

## Access Control Decisions: The Warden

The warden is a HTTP API allowing you to perform these access requests.
The warden knows two endpoints:

* `/warden/allowed`: Check if a subject is allowed to do something.
* `/warden/token/allowed`: Check if the subject of a token is allowed to do something.

Both endpoints use policies to compute the result and are documented in the [HTTP API Documentation](http://docs.hydra13.apiary.io/#reference/warden:-access-control).