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

## Introduction

Using **Access Control Policies**, Hydra is able to answer the question:

> **Who** is **able** to do **what** on **something** under a certain **circumstance**

* **Who** An arbitrary unique subject name, for example "ken" or "printer-service.mydomain.com".
* **Able**: The effect which can be either "allow" or "deny".
* **What**: An arbitrary action name, for example "delete", "create" or "scoped:action:something".
* **Something**: An arbitrary unique resource name, for example "something", "resources.articles.1234" or some uniform
    resource name like "urn:isbn:3827370191".
* **Circumstance**: The circumstance under which the policy can be applied. Typically contains information about the
    environment such as the IP Address, the time or date of access, or ownership. (optional)

The evaluation logic follows these rules: 

* By default, all requests are denied.
* An explicit allow overrides this default.
* An explicit deny overrides any allows.

To decide whether access is allowed or not, ORY Hydra uses Access Control Policies represented as JSON. Values `actions`, `subjects`
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

## Warden API

The warden is a HTTP API allowing you to perform access requests.
The warden knows two endpoints:

* `/warden/allowed`: Check if a subject is allowed to do something.
* `/warden/token/allowed`: Check if the subject of a token is allowed to do something.

Both endpoints use Access Control Policies as documented in the next sections to compute the result.
The API endpoints are documented in the [HTTP API Documentation](http://docs.hydra13.apiary.io/#reference/warden:-access-control).

## Groups

Subjects can be grouped together using the [Warden Group API](https://hydra13.docs.apiary.io/#reference/warden/wardengroups).
This allows to create a set of groups (e.g. `admin`, `moderator`) with a specific set of policies attached to them.
A subject (e.g. a human or a service) belonging to one or more groups inherits all policies attached to those groups.

## Best Practices

This sections gives an overview of best practices for access control policies
we developed over the years at ORY.

### Scalability

Access Control Policies not using any regular expressions are quite scalable. Defining many regular expressions and
a lot of policies (50.000+) may have a notable performance impact on CPU, your database and generally increase response
times. This is because regular expressions can not be indexed. Regular expressions have a complexity of `O(n)` in Go,
but that is still getting slow when you define too many.

Try to solve your access control definitions with a few generalized policies, and try to leverage Warden Groups.

### URNs

> “There are only two hard things in Computer Science: cache invalidation and naming things.”
-- Phil Karlton

URN naming is as hard as naming API endpoints. Thankfully, by doing the latter, the former is usually solved as well.
We will explore further best practices in the following sections.

### Scope the Organization Name

A rule of thumb is to prefix resource names with a domain that represents the organization creating the software.

* **Do not:** `<some-id>`
* **Do:** `<organizaion-id>:<some-id>`

### Scope Actions, Resources and Subjects

It is wise to scope actions, resources, and subjects in order to prevent name collisions:

* **Do not:** `myorg.com:<subject-id>`, `myorg.com:<resource-id>`, `myorg.com:<action-id>`
* **Do:** `myorg.com:subjects:<subject-id>`, `myorg.com:resources:<resource-id>`, `myorg.com:actions:<action-id>`
* **Do:** `subjects:myorg.com:<subject-id>`, `resources:myorg.com:<resource-id>`, `actions:myorg.com:<action-id>`

### Multi-Tenant Systems

Multi-tenant systems typically have resources which should not be access by other tenants in the system. This can be
achieved by adding the tenant id to the URN:

* **Do:** `resources:myorg.com:tenants:<tenant-id>:<resource-id>`

In some environments, it is common to have organizations and projects belonging to those organizations. Here, the
following URN semantics can be used:

* **Do:** `resources:myorg.com:organizations:<organization-id>:projects:<project-id>:<resource-id>`

## Conditions & Context

Conditions are defined in policies. Contexts are defined in access control requests. Conditions use contexts and decide
if a policy is responsible for handling the access request at hand.

Conditions are functions returning true or false given a context. Because conditions implement logic,
they must be programmed. ORY Hydra uses conditions defined in [ORY Ladon](https://github.com/ory/ladon/#conditions).
Adding new condition handlers must be done through creating a pull request in the ORY Ladon repository.

A condition has always the same JSON format:

```json
{
  "subjects": ["..."],
  "actions" : ["..."],
  "effect": "allow",
  "resources": ["..."],
  "conditions": {
    "this-key-will-be-matched-with-the-context": {
      "type": "SomeConditionType",
      "options": {
        "some": "configuration options set by the condition type"
      }
    }
  }
}
```

The context in the access request made to ORY Hydra's Warden API must match the specified key in the condition
in order to be evaluated by the condition logic:

```json
{
  "subject": "...",
  "action" : "...",
  "resource": "...",
  "context": {
    "this-key-will-be-matched-with-the-context": { "foo": "bar" }
  }
}
```

### CIDR Condition

The CIDR condition matches CIDR IP Ranges. An exemplary policy definition could look as follows.

```json
{
  "description": "One policy to rule them all.",
  "subjects": ["users:maria"],
  "actions" : ["delete", "create", "update"],
  "effect": "allow",
  "resources": ["resources:articles:<.*>"],
  "conditions": {
    "remoteIPAddress": {
      "type": "CIDRCondition",
      "options": {
        "cidr": "192.168.0.0/16"
      }
    }
  }
}
```

The following access request would be allowed.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "remoteIPAddress": "192.168.0.5"
  }
}
```

The next access request would be denied as the condition is not fulfilled and thus no policy is matched.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "remoteIPAddress": "255.255.0.0"
  }
}
```

The next access request would also be denied as the context is not using the key `remoteIPAddress` but instead `someOtherKey`.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someOtherKey": "192.168.0.5"
  }
}
```

### String Equal Condition

Checks if the value passed in the access request's context is identical with the string that was given initially.

```json
{
  "description": "One policy to rule them all.",
  "subjects": ["users:maria"],
  "actions" : ["delete", "create", "update"],
  "effect": "allow",
  "resources": ["resources:articles:<.*>"],
  "conditions": {
    "someKeyName": {
      "type": "StringEqualCondition",
      "options": {
        "equals": "the-value-should-be-this"
      }
    }
  }
}
```

The following access request would be allowed.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someKeyName": "the-value-should-be-this"
  }
}
```

The following access request would be denied.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someKeyName": "this-is-a-different-value"
  }
}
```

### String Match Condition

Checks if the value passed in the access request's context matches the regular expression that was given initially.

```json
{
  "description": "One policy to rule them all.",
  "subjects": ["users:maria"],
  "actions" : ["delete", "create", "update"],
  "effect": "allow",
  "resources": ["resources:articles:<.*>"],
  "conditions": {
    "someKeyName": {
      "type": "StringMatchCondition",
      "options": {
        "equals": "regex-pattern-here.+"
      }
    }
  }
}
```

The following access request would be allowed.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someKeyName": "regex-pattern-here-matches"
  }
}
```

The following access request would be denied.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someKeyName": "regex-pattern-here"
  }
}
```

### Subject Condition

Checks if the access request's subject is identical with the string specified in the context.

```json
{
  "description": "One policy to rule them all.",
  "subjects": ["users:maria"],
  "actions" : ["delete", "create", "update"],
  "effect": "allow",
  "resources": ["resources:articles:<.*>"],
  "conditions": {
    "owner": {
      "type": "EqualsSubjectCondition",
      "options": {}
    }
  }
}
```

The following access request would be allowed.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "owner": "users:maria"
  }
}
```

The following access request would be denied.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "owner": "another-user"
  }
}
```

This condition makes more sense when being used with access tokens where the subject is extracted from the token.

### String Pairs Equal Condition

Checks if the value passed in the access request's context contains two-element arrays and that both elements in each pair are equal.

```json
{
  "description": "One policy to rule them all.",
  "subjects": ["users:maria"],
  "actions" : ["delete", "create", "update"],
  "effect": "allow",
  "resources": ["resources:articles:<.*>"],
  "conditions": {
    "someKey": {
      "type": "StringPairsEqualCondition",
      "options": {}
    }
  }
}
```

The following access request would be allowed.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someKey": [
      ["some-arbitrary-pair-value", "some-arbitrary-pair-value"],
      ["some-other-arbitrary-pair-value", "some-other-arbitrary-pair-value"]
    ]
  }
}
```

The following access request would be denied.

```json
{
  "subject": "users:maria",
  "action" : "delete",
  "resource": "resources:articles:12345",
  "context": {
    "someKey": [
      ["some-arbitrary-pair-value", "some-other-arbitrary-pair-value"]
    ]
  }
}
```
