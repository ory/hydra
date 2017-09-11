# Secure the consent app

This tutorial requires to have read and understood [OAuth 2.0 & OpenID Connect](../oauth2.md).

A consent app should never use the root hydra credentials, and fortunately you can create in two simple steps:

## 1. Create the client in Hydra

A consent app needs to communicate with hydra, so it needs a client:

```json
{
    "id": "YOURCONSENTID",
    "client_secret": "YOURCONSENTSECRET",
    "client_name": "consent",
    "redirect_uris": [],
    "grant_types": [
        "client_credentials"
    ],
    "response_types": [
        "token"
    ],
    "scope": "hydra.keys.get"
}
```

`hydra.keys.get` is the only scope that's strictly required for the consent flow, but you may need to
use other scopes.

To create the client you can save the json configuration on a file ```consent.json``` and then issue the command

```
$ hydra clients import consent.json
```

## 2. Grant permissions to the client

Giving the `hydra.keys.get` scope is not enough. Hydra's warden needs an explicit policy to access hydra's keys.

```json
{
    "actions": [
        "get"
    ] ,
    "conditions": {},
    "description": "Allow consent app to access hydra's keys" ,
    "effect": "allow" ,
    "id": "consent_keys" ,
    "resources": [
        "rn:hydra:keys:hydra.consent.challenge:public"
        "rn:hydra:keys:hydra.consent.response:private"
    ] ,
    "subjects": [
        "YOURCONSENTID"
    ]
}
```

We are granting access explicitedly only to the two strictly necessary keys for the consent flow

To create the policy you can save the json configuration on a file ```policy.json``` and then issue the command

```
$ hydra policies create -f policy.json
```
