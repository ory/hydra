# What will happen if an error occurs during an OAuth2 flow?

The user agent will either, according to spec, be redirected to the OAuth2 client who initiated the request, if possible. If not, the user agent will be redirected to the identity provider
endpoint and an `error` and `error_description` query parameter will be appended to it's URL.
