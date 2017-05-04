# How to deal with mobile apps?

Authors of apps running on client side (native apps, single page apps, hybrid apps) have for some
time relied on the Resource Owner Password Credentials grant. This is highly discouraged by the IETF, and replaced
with recommendations in [OAuth 2.0 for Native Apps](https://tools.ietf.org/html/draft-ietf-oauth-native-apps-03).

To keep things short, it allows you to perform the normal `authorize_code` flows without supplying a password. Hydra
allows this by setting the public flag, for example:

```sh
hydra clients create \
    --id my-id \
    --is-public \
    -r code,id_token \
    -g authorization_code,refresh_token \
    -a offline,openid \
    -c https://mydomain/callback
```
