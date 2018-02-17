## JavaScript SDK

### Installation

To install the JavaScript SDK, run:

```
npm install --save ory-hydra-sdk
```

### Configuration

#### Basic configuration

```js
const Hydra = require('ory-hydra-sdk')

// Set this to Hydra's URL
Hydra.ApiClient.instance.basePath = 'http://localhost:4444'

// Configure basic authorization
Hydra.ApiClient.instance.authentications.basic.username = 'client-id'
Hydra.ApiClient.instance.authentications.basic.password = 'client-secret'
```

#### OAuth2 configuration

We need OAuth2 capabilities in order to make authorized API
calls. This currently requires writing your own OAuth2 mechanism.
Thankfully, libraries like `passport-js` and `simple-oauth2` exist.

Here, we will use `simple-oauth2` to configure OAuth2.

```sh
npm i --save simple-oauth2
```

```js
const Hydra = require('ory-hydra-sdk')

// ... configuration from the previous section

const OAuth2 = require('simple-oauth2')

// A list of scopes, tab separated. Use hydra.* to grant
// all hydra scopes
const scope = 'hydra.* some-other-scope'

const oauth2 = OAuth2.create({
  client: {
    id: 'client-id',
    secret: 'client-secret'
  },
  auth: {
    tokenHost: endpoint = 'http://localhost:4444',
    authorizePath: authorizePath = '/oauth2/auth',
    tokenPath: tokenPath = '/oauth2/token'
  },
  // These are important for simple-oauth2 to work properly.
  options: {
    useBodyAuth: false,
    useBasicAuthorizationHeader: true
  }
})

// Next we need to fetch a token, let's wrap that in a
// function called refreshToken
const refreshToken = () => oauth2.clientCredentials
  .getToken({ scope })
  .then((result) => {
    const token = oauth2.accessToken.create(result);
    const hydraClient = Hydra.ApiClient.instance
    hydraClient.authentications.oauth2.accessToken = token.token.access_token
    return Promise.resolve(token)
  })
```

Of course, the `refreshToken` method can be improved. Read more
on this topic [here](https://github.com/lelylan/simple-oauth2#access-token-object).

### API Usage

Let's use `refreshToken` to request a new access token and make
an authorized API call:

```js
refreshToken().then(() => {
  // It is important to note that must not return `hydra.listOAuth2Clients`
  // directly inside a Promise, otherwise you will encounter the "superagent double callback bug".

  const hydra = new Hydra.OAuth2Api()

  // for example, let's fetch all OAuth2 clients
  hydra.listOAuth2Clients((error, data, response) => {
    if (error) {
      // a network error occurred.
      throw error
    } else if (response.statusCode < 200 || response.statusCode >= 400) {
      // an application error occurred.
      throw new Error('Consent endpoint gave status code ' + response.statusCode + ', but status code 200 was expected.')
    }

    console.log(response) // a list of OAuth2 clients.
  })
})
```

### API Docs

API docs are available [here](https://github.com/ory/hydra/blob/master/sdk/js/hydra/swagger/README.md).
Please note that those docs are generated and may introduce bugs if code examples are used 1:1. Especially
the package name is not correct.
