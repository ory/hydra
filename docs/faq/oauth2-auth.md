# Should I use OAuth2 tokens for authentication?

OAuth2 tokens are like money. It allows you to buy stuff, but the cashier does not really care if the money is
yours or if you stole it, as long as it's valid money. Depending on what you understand as authentication, this is a yes and no answer:

* **Yes:** You can use access tokens to find out which user ("subject") is performing an action in a resource provider (blog article service, shopping basket, ...).
Coming back to the money example: *You*, the subject, receives a cappuccino from the vendor (resource provider) in exchange for money (access token).
* **No:** Never use access tokens for logging people in, for example `http://myapp.com/login?access_token=...`.
Coming back to the money example: The police officer ("authentication server") will not accept money ("access token") as a proof of identity ("it's really you"). Unless he is corrupt ("vulnerable"), of course.

In the second example ("authentication server"), you must use OpenID Connect ID Tokens.