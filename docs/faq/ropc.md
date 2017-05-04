# Why is the Resource Owner Password Credentials grant not supported?

The following is a copy of the original [comment on GitHub](https://github.com/ory/hydra/pull/297#issuecomment-294282671):

I took a long time for this issue, primarily because I felt very uncomfortable implementing it. The ROCP grant is something from the "dark ages" of OAuth2 and there are suitable replacements for mobile clients, such as public oauth2 clients, which are supported by Hydra: https://tools.ietf.org/html/draft-ietf-oauth-native-apps-09

The OAuth2 Thread Model explicitly states that the ROPC grant is commonly used in legacy/migration scenarios, and

>   This grant type has higher
   risk because it maintains the UID/password anti-pattern.
   Additionally, because the user does not have control over the
   authorization process, clients using this grant type are not limited   by scope but instead have potentially the same capabilities as the
   user themselves.  As there is no authorization step, the ability to
   offer token revocation is bypassed.

> Because passwords are often used for more than 1 service, this
   anti-pattern may also put at risk whatever else is accessible with
   the supplied credential.  Additionally, any easily derived equivalent
   (e.g., joe@example.com and joe@example.net) might easily allow
   someone to guess that the same password can be used elsewhere.

>    Impact: The resource server can only differentiate scope based on the
   access token being associated with a particular client.  The client
   could also acquire long-lived tokens and pass them up to an
   attacker's web service for further abuse.  The client, eavesdroppers,
   or endpoints could eavesdrop the user id and password.

>    o  Except for migration reasons, minimize use of this grant type.

- [source](https://tools.ietf.org/html/rfc6819#section-4.4.3)

Thus, I decided to not implement the ROPC grant in Hydra. Over time, I will add documentation how to deal with mobile scenarios and similar.