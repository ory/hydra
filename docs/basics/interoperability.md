## Interoperability

We did not want to provide you with LDAP, Active Directory, ADFS, SAML-P, SharePoint Apps, ...
integrations which probably won't work well anyway. Instead we decided to rely on cryptographic tokens
(JSON Web Tokens) for authenticating users and getting their consent. This gives you all the freedom you need with
very little effort. JSON Web Tokens are supported by all web programming languages and Hydra's
[JSON Web Key API](jwk.html) offers a nice way to deal with certificates and keys. Your users won't notice the difference.

![OAuth2 Workflow](../images/hydra-authentication.gif)