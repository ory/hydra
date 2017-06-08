# FAQ

This file keeps track of questions and discussions from Gitter

## Logout

> Kareem Diaa @kimooz 15:41
Thanks @arekkas. I had two other questions:
1- Is there a way to revoke all access tokens for a certain user("log out user")
2- How can I inform the consent app that this user logged out?

> Aeneas @arekkas 15:42
no this isn't supported currently
\2. you can't because log out and revoking access tokens are two things
and it would require an additional api or something, which makes the consent app harder to write and integrate

> Kareem Diaa @kimooz 15:43
So can you suggest a workaround?
I want implement single sign off

> Aeneas @arekkas 15:44
the user has the access and refresh token right
in his browser or somewhere

> Kareem Diaa @kimooz 15:44
yah

> Aeneas @arekkas 15:44
ok so why not make a request to /oauth2/revoke
and pass that refresh token
(you will probably need a proxy with a client id and secret for that to be possible, but you get the point)

> Kareem Diaa @kimooz 15:46
yah but the moment he refreshes, the client will hit on hydra and then consent where it will find that this user is already logged in
and will return a new token although he should have logged out
ohh so you mean have two requests one for hydra to revoke and one for consent to log out correct?

> Aeneas @arekkas 15:47
yes
