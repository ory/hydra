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

# Access Control Policies

Hydra uses the Access Control Library [Ladon](https://github.com/ory-am/ladon).
For a deep dive, it is a good idea to read the [Ladon Docs](https://github.com/ory-am/ladon#ladon).

In Hydra, policy based access control is when you decide if:

- Aaron (subject) is allowed (effect) to create (action) a new forum post (resource) when accessing the forum website from IP 192.168.178.3 (context).
- Richard (subject) is allowed (effect) to delete (action) a status update (resource) when he is the author (context).

Or, more *generalized:* **Who** is **able** to do **what** on **something** with some **context**.

* **Who (Subject)**: An arbitrary unique subject name, for example "ken" or "printer-service.mydomain.com".
* **Able (Effect)**: The effect which is always "allow" or "deny".
* **What (Action)**: An arbitrary action name, for example "delete", "create" or "scoped:action:something".
* **Something (Resource)**: An arbitrary unique resource name, for example "something", "resources:articles:1234" or some uniform resource name like "urn:isbn:3827370191".
* **Context (Context)**: The current context which may environment information like the IP Address, request date, the resource owner name, the department ken is working in and anything you like.

Policies are JSON documents managed via the [Policy API](http://docs.hdyra.apiary.io/#reference/policies).

```
{
  // A required unique identifier. Used primarily for database retrieval.
  "id": "68819e5a-738b-41ec-b03c-b58a1b19d043",
  
  // A optional human readable description.
  "description": "something humanly readable",
  
  // A subject can be an user or a service. It is the "who" in "who is allowed to do what on something".
  // As you can see here, you can use regular expressions inside < >.
  "subjects": ["user", "<peter|max>"],
    
  
  // Should access be allowed or denied?
  // Note: If multiple policies match an access request, ladon.DenyAccess will always override ladon.AllowAccess
  // and thus deny access.
  "effect": "allow",
  
  // Which resources this policy affects.
  // Again, you can put regular expressions in inside < >.
  "resources": ["articles:<[0-9]+>"],
  
  // Which actions this policy affects. Supports RegExp
  // Again, you can put regular expressions in inside < >.
  "actions": ["create","update"],
  
  // Under which conditions this policy is "active".
  "conditions": {
    "owner": {
     // In this example, the policy is only "active" when the requested subject is the owner of the resource as well.
      "type": "EqualsSubjectCondition",
      "options": {}
    }
   }
}
```

## Examples

This example let's everyone, including anonymous users, read public keys. Anonymous users have no special ID and are
simply empty subject strings in Hydra.

```
{
  "description": "Allow everyone including anonymous users to read JSON Web Keys having Key ID *public*.",
  "subjects": [
    "<.*>"
  ],
  "effect": "allow",
  "resources": [
    "rn:hydra:keys:<[^:]+>:public"
  ],
  "actions": [
    "get"
  ]
}
```

Deny anyone from reading private JWKs:

```
{
  "description": "Explicitly deny everyone reading JSON Web Keys with Key ID *private*.",
  "subjects": [
    "<.*>"
  ],
  "effect": "deny",
  "resources": [
    "rn:hydra:keys:<[^:]+>:private"
  ],
  "actions": [
    "get"
  ]
}
```

## Warden

The Warden is usually called from your own services ("resource providers"), not from third parties. Hydra prevents
third parties from having access to these endpoints per default, but you can change that with custom policies.

The Warden endpoints are documented [here](http://docs.hdyra.apiary.io/#reference/warden:-access-control-for-resource-providers).