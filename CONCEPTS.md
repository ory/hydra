This document gives an overview of the concepts used in go-iam, including:

* Accounts
* Identity
* Groups
* Permissions
* Policies

# Account

A **account** (also known as user) has long living credentials.

*Example:* Bob registeres a new **account** at www.example.com, providing a unique email address and a password (:= credentials).

The account entity:

```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "email": "foo@bar",
  "password": "hashed-password"
}
```

# Identity (work in progress)

A account is an identity, but an identity is not neccessarily a account/user. An identity might also be
a service.

*Example (incomplete!):* Bob is allowed to create articles (*allow* `POST /articles`),
but is not allowed to changes categories of articles (*disallow* `PUT /categories/example-category/example-article`).  
However, when creating an article, he is allowed to additionally choose a category for it.
In this case, the article service delegates (when `POST /articles` allowed, allow `PUT /categories/example-category/example-article`) the permission to choose
a category for

# Groups

Accounts can be grouped together, forming a **group** of users. All users of ag roup
share the same policies. Groups are the R(ole) in RBAC (Role Based Access Control).
Accounts can be part of a unlimited amount of groups.

The group entity:

```json
{
  "id": "c716ef6f-e6c7-4804-bf22-8c3bc7d16dbc",
  "users": ["81328fba-2353-4348-8571-bdc0e15f6ee3"],
  "policies": ["bfcdfae2-cb0e-4c84-94f7-e04d6cac539c"]
}
```

# Permissions

Permissions allow (or disallow) an account, identity or group some type of access (action) to one or more resources.

* Example (User-based): Bob has permission to create and modify all articles.
* Example (Resource-based): The article `example-article` can be modified by Bob, Susan and all members of group `example-group`

The permission entity:

```json
{
  "effect": "Allow",
  "action": ["an:content:article.create"],
  "resources": [
    "rn:content:articles.44efef16-12bc-4752-a0c5-2e768622e46b",
    "rn:content:articles.363bab48-82f1-4b99-ba69-b1cf0e18345e",
    "rn:content:articles.fb33c4e1-2f16-4701-94d1-ee4198968ab4"
  ]
}
```

* `"effect"` (MUST): Can be `"Allow"` or `"Deny"`
* `"action"` (MUST): Is arbitrary. It is recommened to use a layout like `an:<service>:<action>` (an short for *action name*).  
Each key should match `[a-zA-Z0-9\-\.]+`, while, `.` may be used for nesting and `*` for wildcards:
`an:content:article.create` or `an:content:article.modify-timestamp` or `an:content:article.*`
* `"resource"` (OPTIONAL): A collection of arbitrary resource names.  
It is recommened to use a layout like `rn:<service>:<resource-uri>` (rn short for *resource name*).
Each key should match `[a-zA-Z0-9\-\.\*]+`, while `.` replaces `/` and `*` is used for wildcards:
`rn:content:articles.83299f22-5958-469b-9cd4-5d0e25c5a7bb` or `rn:content:articles:*`

# Policies

A policy is a document that provides a formal statement of on or more permissions. Policies are versioned, so you can quickly recover a previous state.

The policy entity:

```json
{
  "id": "bf4f5b8d-3df7-4369-b432-88e685462394",
  "previous": "08071c8e-88de-4744-baae-78e3abfcb924",
  "statements": [
    {
      "effect": "Allow",
      "action": ["an:content:article.*"],
      "resource": "rn:content:articles.83299f22-5958-469b-9cd4-5d0e25c5a7bb"
    },
    {
      "effect": "Allow",
      "action": ["an:content:article.create"],
      "resources": [
        "rn:content:articles.44efef16-12bc-4752-a0c5-2e768622e46b",
        "rn:content:articles.363bab48-82f1-4b99-ba69-b1cf0e18345e",
        "rn:content:articles.fb33c4e1-2f16-4701-94d1-ee4198968ab4"
      ]
    }
  ]
}
```
