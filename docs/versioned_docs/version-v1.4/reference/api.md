---
title: REST API
id: api
---



Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.

> You are viewing REST API documentation. This documentation is auto-generated from a swagger specification which
itself is generated from annotations in the source code of the project. It is possible that this documentation includes
bugs and that code samples are incomplete or wrong.
>
> If you find issues in the respective documentation, please do not edit the
Markdown files directly (as they are generated) but raise an issue on the project's GitHub presence instead. This documentation
will improve over time with your help! If you have ideas how to improve this part of the documentation, feel free to
share them in a [GitHub issue](https://github.com/ory/docs/issues/new) any time.

## Authentication

- HTTP Authentication, scheme: basic - OAuth 2.0 Authorization.   - Flow: authorizationCode
  - OAuth 2.0 Authorization URL = [ https://hydra.demo.ory.sh/oauth2/auth](https://hydra.demo.ory.sh/oauth2/auth)
  - OAuth 2.0 Token URL = [ https://hydra.demo.ory.sh/oauth2/token](https://hydra.demo.ory.sh/oauth2/token)
  - OAuth 2.0 Scope

    |Scope|Scope Description|
    |---|---|
    |offline|A scope required when requesting refresh tokens (alias for `offline_access`)|
    |offline_access|A scope required when requesting refresh tokens|
    |openid|Request an OpenID Connect ID Token|

<a id="ory-hydra-public-endpoints"></a>
## Public Endpoints

<a id="opIdwellKnown"></a>

### JSON Web Keys Discovery

```
GET /.well-known/jwks.json HTTP/1.1
Accept: application/json

```

This endpoint returns JSON Web Keys to be used as public keys for verifying OpenID Connect ID Tokens and,
if enabled, OAuth 2.0 JWT Access Tokens. This endpoint can be used with client libraries like
[node-jwks-rsa](https://github.com/auth0/node-jwks-rsa) among others.

#### Responses

<a id="json-web-keys-discovery-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|JSONWebKeySet|[JSONWebKeySet](#schemajsonwebkeyset)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-wellKnown">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-wellKnown-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-wellKnown-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-wellKnown-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-wellKnown-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-wellKnown-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-wellKnown-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-wellKnown-shell">

```shell
curl -X GET /.well-known/jwks.json \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-wellKnown-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/.well-known/jwks.json", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-wellKnown-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/.well-known/jwks.json', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-wellKnown-java">

```java
// This sample needs improvement.
URL obj = new URL("/.well-known/jwks.json");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-wellKnown-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/.well-known/jwks.json',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-wellKnown-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/.well-known/jwks.json',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIddiscoverOpenIDConfiguration"></a>

### OpenID Connect Discovery

```
GET /.well-known/openid-configuration HTTP/1.1
Accept: application/json

```

The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll
your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this
flow at https://openid.net/specs/openid-connect-discovery-1_0.html .

Popular libraries for OpenID Connect clients include oidc-client-js (JavaScript), go-oidc (Golang), and others.
For a full list of clients go here: https://openid.net/developers/certified/

#### Responses

<a id="openid-connect-discovery-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|wellKnown|[wellKnown](#schemawellknown)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "authorization_endpoint": "https://playground.ory.sh/ory-hydra/public/oauth2/auth",
  "backchannel_logout_session_supported": true,
  "backchannel_logout_supported": true,
  "claims_parameter_supported": true,
  "claims_supported": [
    "string"
  ],
  "end_session_endpoint": "string",
  "frontchannel_logout_session_supported": true,
  "frontchannel_logout_supported": true,
  "grant_types_supported": [
    "string"
  ],
  "id_token_signing_alg_values_supported": [
    "string"
  ],
  "issuer": "https://playground.ory.sh/ory-hydra/public/",
  "jwks_uri": "https://playground.ory.sh/ory-hydra/public/.well-known/jwks.json",
  "registration_endpoint": "https://playground.ory.sh/ory-hydra/admin/client",
  "request_parameter_supported": true,
  "request_uri_parameter_supported": true,
  "require_request_uri_registration": true,
  "response_modes_supported": [
    "string"
  ],
  "response_types_supported": [
    "string"
  ],
  "revocation_endpoint": "string",
  "scopes_supported": [
    "string"
  ],
  "subject_types_supported": [
    "string"
  ],
  "token_endpoint": "https://playground.ory.sh/ory-hydra/public/oauth2/token",
  "token_endpoint_auth_methods_supported": [
    "string"
  ],
  "userinfo_endpoint": "string",
  "userinfo_signing_alg_values_supported": [
    "string"
  ]
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-discoverOpenIDConfiguration">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-discoverOpenIDConfiguration-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-discoverOpenIDConfiguration-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-discoverOpenIDConfiguration-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-discoverOpenIDConfiguration-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-discoverOpenIDConfiguration-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-discoverOpenIDConfiguration-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-discoverOpenIDConfiguration-shell">

```shell
curl -X GET /.well-known/openid-configuration \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-discoverOpenIDConfiguration-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/.well-known/openid-configuration", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-discoverOpenIDConfiguration-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/.well-known/openid-configuration', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-discoverOpenIDConfiguration-java">

```java
// This sample needs improvement.
URL obj = new URL("/.well-known/openid-configuration");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-discoverOpenIDConfiguration-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/.well-known/openid-configuration',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-discoverOpenIDConfiguration-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/.well-known/openid-configuration',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdisInstanceReady"></a>

### Check readiness status

```
GET /health/ready HTTP/1.1
Accept: application/json

```

This endpoint returns a 200 status code when the HTTP server is up running and the environment dependencies (e.g.
the database) are responsive as well.

If the service supports TLS Edge Termination, this endpoint does not require the
`X-Forwarded-Proto` header to be set.

Be aware that if you are running multiple nodes of this service, the health status will never
refer to the cluster state, only to a single instance.

#### Responses

<a id="check-readiness-status-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|healthStatus|[healthStatus](#schemahealthstatus)|
|503|[Service Unavailable](https://tools.ietf.org/html/rfc7231#section-6.6.4)|healthNotReadyStatus|[healthNotReadyStatus](#schemahealthnotreadystatus)|

##### Examples

###### 200 response

```json
{
  "status": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-isInstanceReady">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-isInstanceReady-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceReady-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceReady-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceReady-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceReady-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceReady-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-isInstanceReady-shell">

```shell
curl -X GET /health/ready \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceReady-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/health/ready", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceReady-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/health/ready', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceReady-java">

```java
// This sample needs improvement.
URL obj = new URL("/health/ready");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceReady-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/health/ready',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceReady-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/health/ready',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdoauthAuth"></a>

### The OAuth 2.0 authorize endpoint

```
GET /oauth2/auth HTTP/1.1
Accept: application/json

```

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows.
OAuth2 is a very popular protocol and a library for your programming language will exists.

To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

#### Responses

<a id="the-oauth-2.0-authorize-endpoint-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|302|[Found](https://tools.ietf.org/html/rfc7231#section-6.4.3)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 401 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-oauthAuth">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-oauthAuth-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauthAuth-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauthAuth-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauthAuth-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauthAuth-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauthAuth-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-oauthAuth-shell">

```shell
curl -X GET /oauth2/auth \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauthAuth-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/oauth2/auth", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauthAuth-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauthAuth-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauthAuth-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/oauth2/auth',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauthAuth-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/oauth2/auth',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdrevokeOAuth2Token"></a>

### Revoke OAuth2 tokens

```
POST /oauth2/revoke HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Accept: application/json

```

Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no
longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token.
Revoking a refresh token also invalidates the access token that was created with it. A token may only be revoked by
the client the token was generated for.

#### Request body

```yaml
token: string

```

<a id="revoke-oauth2-tokens-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» token|body|string|true|none|

#### Responses

<a id="revoke-oauth2-tokens-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 401 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
basic, oauth2
</aside>

#### Code samples

<div class="tabs" id="tab-revokeOAuth2Token">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-revokeOAuth2Token-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeOAuth2Token-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeOAuth2Token-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeOAuth2Token-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeOAuth2Token-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeOAuth2Token-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-revokeOAuth2Token-shell">

```shell
curl -X POST /oauth2/revoke \
  -H 'Content-Type: application/x-www-form-urlencoded' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeOAuth2Token-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/x-www-form-urlencoded"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("POST", "/oauth2/revoke", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeOAuth2Token-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "token": "string"
}';
const headers = {
  'Content-Type': 'application/x-www-form-urlencoded',  'Accept': 'application/json'
}

fetch('/oauth2/revoke', {
  method: 'POST',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeOAuth2Token-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/revoke");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeOAuth2Token-python">

```python
import requests

headers = {
  'Content-Type': 'application/x-www-form-urlencoded',
  'Accept': 'application/json'
}

r = requests.post(
  '/oauth2/revoke',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeOAuth2Token-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/x-www-form-urlencoded',
  'Accept' => 'application/json'
}

result = RestClient.post '/oauth2/revoke',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIddisconnectUser"></a>

### OpenID Connect Front-Backchannel enabled Logout

```
GET /oauth2/sessions/logout HTTP/1.1

```

This endpoint initiates and completes user logout at ORY Hydra and initiates OpenID Connect Front-/Back-channel logout:

https://openid.net/specs/openid-connect-frontchannel-1_0.html
https://openid.net/specs/openid-connect-backchannel-1_0.html

#### Responses

<a id="openid-connect-front-backchannel-enabled-logout-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|302|[Found](https://tools.ietf.org/html/rfc7231#section-6.4.3)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-disconnectUser">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-disconnectUser-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-disconnectUser-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-disconnectUser-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-disconnectUser-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-disconnectUser-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-disconnectUser-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-disconnectUser-shell">

```shell
curl -X GET /oauth2/sessions/logout

```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-disconnectUser-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/oauth2/sessions/logout", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-disconnectUser-node">

```nodejs
const fetch = require('node-fetch');

fetch('/oauth2/sessions/logout', {
  method: 'GET'
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-disconnectUser-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/sessions/logout");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-disconnectUser-python">

```python
import requests

r = requests.get(
  '/oauth2/sessions/logout',
  params={)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-disconnectUser-ruby">

```ruby
require 'rest-client'
require 'json'

result = RestClient.get '/oauth2/sessions/logout',
  params: {}

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdoauth2Token"></a>

### The OAuth 2.0 token endpoint

```
POST /oauth2/token HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Accept: application/json

```

The client makes a request to the token endpoint by sending the
following parameters using the "application/x-www-form-urlencoded" HTTP
request entity-body.

> Do not implement a client for this endpoint yourself. Use a library. There are many libraries
> available for any programming language. You can find a list of libraries here: https://oauth.net/code/
>
> Do note that Hydra SDK does not implement this endpoint properly. Use one of the libraries listed above!

#### Request body

```yaml
grant_type: string
code: string
refresh_token: string
redirect_uri: string
client_id: string

```

<a id="the-oauth-2.0-token-endpoint-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» grant_type|body|string|true|none|
|» code|body|string|false|none|
|» refresh_token|body|string|false|none|
|» redirect_uri|body|string|false|none|
|» client_id|body|string|false|none|

#### Responses

<a id="the-oauth-2.0-token-endpoint-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|oauth2TokenResponse|[oauth2TokenResponse](#schemaoauth2tokenresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "access_token": "string",
  "expires_in": 0,
  "id_token": "string",
  "refresh_token": "string",
  "scope": "string",
  "token_type": "string"
}
```

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
basic, oauth2
</aside>

#### Code samples

<div class="tabs" id="tab-oauth2Token">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-oauth2Token-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauth2Token-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauth2Token-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauth2Token-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauth2Token-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-oauth2Token-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-oauth2Token-shell">

```shell
curl -X POST /oauth2/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauth2Token-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/x-www-form-urlencoded"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("POST", "/oauth2/token", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauth2Token-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "grant_type": "string",
  "code": "string",
  "refresh_token": "string",
  "redirect_uri": "string",
  "client_id": "string"
}';
const headers = {
  'Content-Type': 'application/x-www-form-urlencoded',  'Accept': 'application/json'
}

fetch('/oauth2/token', {
  method: 'POST',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauth2Token-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/token");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauth2Token-python">

```python
import requests

headers = {
  'Content-Type': 'application/x-www-form-urlencoded',
  'Accept': 'application/json'
}

r = requests.post(
  '/oauth2/token',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-oauth2Token-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/x-www-form-urlencoded',
  'Accept' => 'application/json'
}

result = RestClient.post '/oauth2/token',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIduserinfo"></a>

### OpenID Connect Userinfo

```
GET /userinfo HTTP/1.1
Accept: application/json

```

This endpoint returns the payload of the ID Token, including the idTokenExtra values, of
the provided OAuth 2.0 Access Token.

For more information please [refer to the spec](http://openid.net/specs/openid-connect-core-1_0.html#UserInfo).

#### Responses

<a id="openid-connect-userinfo-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|userinfoResponse|[userinfoResponse](#schemauserinforesponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "birthdate": "string",
  "email": "string",
  "email_verified": true,
  "family_name": "string",
  "gender": "string",
  "given_name": "string",
  "locale": "string",
  "middle_name": "string",
  "name": "string",
  "nickname": "string",
  "phone_number": "string",
  "phone_number_verified": true,
  "picture": "string",
  "preferred_username": "string",
  "profile": "string",
  "sub": "string",
  "updated_at": 0,
  "website": "string",
  "zoneinfo": "string"
}
```

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
oauth2
</aside>

#### Code samples

<div class="tabs" id="tab-userinfo">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-userinfo-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-userinfo-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-userinfo-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-userinfo-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-userinfo-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-userinfo-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-userinfo-shell">

```shell
curl -X GET /userinfo \
  -H 'Accept: application/json' \  -H 'Authorization: Bearer {access-token}'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-userinfo-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
        "Authorization": []string{"Bearer {access-token}"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/userinfo", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-userinfo-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json',  'Authorization': 'Bearer {access-token}'
}

fetch('/userinfo', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-userinfo-java">

```java
// This sample needs improvement.
URL obj = new URL("/userinfo");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-userinfo-python">

```python
import requests

headers = {
  'Accept': 'application/json',
  'Authorization': 'Bearer {access-token}'
}

r = requests.get(
  '/userinfo',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-userinfo-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json',
  'Authorization' => 'Bearer {access-token}'
}

result = RestClient.get '/userinfo',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="ory-hydra-administrative-endpoints"></a>
## Administrative Endpoints

<a id="opIdlistOAuth2Clients"></a>

### List OAuth 2.0 Clients

```
GET /clients HTTP/1.1
Accept: application/json

```

This endpoint lists all clients in the database, and never returns client secrets.

OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
The "Link" header is also included in successful responses, which contains one or more links for pagination, formatted like so: '<https://hydra-url/admin/clients?limit={limit}&offset={offset}>; rel="{page}"', where page is one of the following applicable pages: 'first', 'next', 'last', and 'previous'.
Multiple links can be included in this header, and will be separated by a comma.

<a id="list-oauth-2.0-clients-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer(int64)|false|The maximum amount of policies returned.|
|offset|query|integer(int64)|false|The offset from where to start looking.|

#### Responses

<a id="list-oauth-2.0-clients-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A list of clients.|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

<a id="list-oauth-2.0-clients-responseschema"></a>
##### Response Schema

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[oAuth2Client](#schemaoauth2client)]|false|none|none|
|» Client represents an OAuth 2.0 Client.|[oAuth2Client](#schemaoauth2client)|false|none|none|
|»» allowed_cors_origins|[string]|false|none|none|
|»» audience|[string]|false|none|none|
|»» backchannel_logout_session_required|boolean|false|none|Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false.|
|»» backchannel_logout_uri|string|false|none|RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.|
|»» client_id|string|false|none|ClientID  is the id for this client.|
|»» client_name|string|false|none|Name is the human-readable string name of the client to be presented to the end-user during authorization.|
|»» client_secret|string|false|none|Secret is the client's secret. The secret will be included in the create request as cleartext, and then never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users that they need to write the secret down as it will not be made available again.|
|»» client_secret_expires_at|integer(int64)|false|none|SecretExpiresAt is an integer holding the time at which the client secret will expire or 0 if it will not expire. The time is represented as the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time of expiration.  This feature is currently not supported and it's value will always be set to 0.|
|»» client_uri|string|false|none|ClientURI is an URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion.|
|»» contacts|[string]|false|none|none|
|»» created_at|string(date-time)|false|none|CreatedAt returns the timestamp of the client's creation.|
|»» frontchannel_logout_session_required|boolean|false|none|Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false.|
|»» frontchannel_logout_uri|string|false|none|RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the request and to determine which of the potentially multiple sessions is to be logged out; if either is included, both MUST be.|
|»» grant_types|[string]|false|none|none|
|»» jwks|[JoseJSONWebKeySet](#schemajosejsonwebkeyset)|false|none|none|
|»» jwks_uri|string|false|none|URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.|
|»» logo_uri|string|false|none|LogoURI is an URL string that references a logo for the client.|
|»» metadata|[JSONRawMessage](#schemajsonrawmessage)|false|none|none|
|»» owner|string|false|none|Owner is a string identifying the owner of the OAuth 2.0 Client.|
|»» policy_uri|string|false|none|PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.|
|»» post_logout_redirect_uris|[string]|false|none|none|
|»» redirect_uris|[string]|false|none|none|
|»» request_object_signing_alg|string|false|none|JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm.|
|»» request_uris|[string]|false|none|none|
|»» response_types|[string]|false|none|none|
|»» scope|string|false|none|Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens.|
|»» sector_identifier_uri|string|false|none|URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.|
|»» subject_type|string|false|none|SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.|
|»» token_endpoint_auth_method|string|false|none|Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none.|
|»» tos_uri|string|false|none|TermsOfServiceURI is a URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.|
|»» updated_at|string(date-time)|false|none|UpdatedAt returns the timestamp of the last update.|
|»» userinfo_signed_response_alg|string|false|none|JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type.|

##### Examples

###### 200 response

```json
[
  {
    "allowed_cors_origins": [
      "string"
    ],
    "audience": [
      "string"
    ],
    "backchannel_logout_session_required": true,
    "backchannel_logout_uri": "string",
    "client_id": "string",
    "client_name": "string",
    "client_secret": "string",
    "client_secret_expires_at": 0,
    "client_uri": "string",
    "contacts": [
      "string"
    ],
    "created_at": "2020-04-25T11:08:35Z",
    "frontchannel_logout_session_required": true,
    "frontchannel_logout_uri": "string",
    "grant_types": [
      "string"
    ],
    "jwks": {},
    "jwks_uri": "string",
    "logo_uri": "string",
    "metadata": {},
    "owner": "string",
    "policy_uri": "string",
    "post_logout_redirect_uris": [
      "string"
    ],
    "redirect_uris": [
      "string"
    ],
    "request_object_signing_alg": "string",
    "request_uris": [
      "string"
    ],
    "response_types": [
      "string"
    ],
    "scope": "string",
    "sector_identifier_uri": "string",
    "subject_type": "string",
    "token_endpoint_auth_method": "string",
    "tos_uri": "string",
    "updated_at": "2020-04-25T11:08:35Z",
    "userinfo_signed_response_alg": "string"
  }
]
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-listOAuth2Clients">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-listOAuth2Clients-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listOAuth2Clients-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listOAuth2Clients-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listOAuth2Clients-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listOAuth2Clients-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listOAuth2Clients-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-listOAuth2Clients-shell">

```shell
curl -X GET /clients \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listOAuth2Clients-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/clients", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listOAuth2Clients-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/clients', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listOAuth2Clients-java">

```java
// This sample needs improvement.
URL obj = new URL("/clients");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listOAuth2Clients-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/clients',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listOAuth2Clients-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/clients',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdcreateOAuth2Client"></a>

### Create an OAuth 2.0 client

```
POST /clients HTTP/1.1
Content-Type: application/json
Accept: application/json

```

Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.

OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

#### Request body

```json
{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}
```

<a id="create-an-oauth-2.0-client-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[oAuth2Client](#schemaoauth2client)|true|none|

#### Responses

<a id="create-an-oauth-2.0-client-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|oAuth2Client|[oAuth2Client](#schemaoauth2client)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|genericError|[genericError](#schemagenericerror)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 201 response

```json
{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-createOAuth2Client">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-createOAuth2Client-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createOAuth2Client-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createOAuth2Client-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createOAuth2Client-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createOAuth2Client-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createOAuth2Client-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-createOAuth2Client-shell">

```shell
curl -X POST /clients \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createOAuth2Client-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("POST", "/clients", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createOAuth2Client-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/clients', {
  method: 'POST',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createOAuth2Client-java">

```java
// This sample needs improvement.
URL obj = new URL("/clients");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createOAuth2Client-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post(
  '/clients',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createOAuth2Client-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/clients',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetOAuth2Client"></a>

### Get an OAuth 2.0 Client.

```
GET /clients/{id} HTTP/1.1
Accept: application/json

```

Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.

OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

<a id="get-an-oauth-2.0-client.-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|The id of the OAuth 2.0 Client.|

#### Responses

<a id="get-an-oauth-2.0-client.-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|oAuth2Client|[oAuth2Client](#schemaoauth2client)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getOAuth2Client">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getOAuth2Client-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getOAuth2Client-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getOAuth2Client-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getOAuth2Client-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getOAuth2Client-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getOAuth2Client-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getOAuth2Client-shell">

```shell
curl -X GET /clients/{id} \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getOAuth2Client-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/clients/{id}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getOAuth2Client-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/clients/{id}', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getOAuth2Client-java">

```java
// This sample needs improvement.
URL obj = new URL("/clients/{id}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getOAuth2Client-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/clients/{id}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getOAuth2Client-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/clients/{id}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdupdateOAuth2Client"></a>

### Update an OAuth 2.0 Client

```
PUT /clients/{id} HTTP/1.1
Content-Type: application/json
Accept: application/json

```

Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.

OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

#### Request body

```json
{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}
```

<a id="update-an-oauth-2.0-client-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|
|body|body|[oAuth2Client](#schemaoauth2client)|true|none|

#### Responses

<a id="update-an-oauth-2.0-client-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|oAuth2Client|[oAuth2Client](#schemaoauth2client)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-updateOAuth2Client">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-updateOAuth2Client-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateOAuth2Client-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateOAuth2Client-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateOAuth2Client-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateOAuth2Client-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateOAuth2Client-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-updateOAuth2Client-shell">

```shell
curl -X PUT /clients/{id} \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateOAuth2Client-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/clients/{id}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateOAuth2Client-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/clients/{id}', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateOAuth2Client-java">

```java
// This sample needs improvement.
URL obj = new URL("/clients/{id}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateOAuth2Client-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/clients/{id}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateOAuth2Client-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/clients/{id}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIddeleteOAuth2Client"></a>

### Deletes an OAuth 2.0 Client

```
DELETE /clients/{id} HTTP/1.1
Accept: application/json

```

Delete an existing OAuth 2.0 Client by its ID.

OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

<a id="deletes-an-oauth-2.0-client-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|The id of the OAuth 2.0 Client.|

#### Responses

<a id="deletes-an-oauth-2.0-client-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 404 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-deleteOAuth2Client">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-deleteOAuth2Client-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteOAuth2Client-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteOAuth2Client-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteOAuth2Client-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteOAuth2Client-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteOAuth2Client-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-deleteOAuth2Client-shell">

```shell
curl -X DELETE /clients/{id} \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteOAuth2Client-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("DELETE", "/clients/{id}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteOAuth2Client-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/clients/{id}', {
  method: 'DELETE',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteOAuth2Client-java">

```java
// This sample needs improvement.
URL obj = new URL("/clients/{id}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("DELETE");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteOAuth2Client-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.delete(
  '/clients/{id}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteOAuth2Client-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.delete '/clients/{id}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdisInstanceAlive"></a>

### Check alive status

```
GET /health/alive HTTP/1.1
Accept: application/json

```

This endpoint returns a 200 status code when the HTTP server is up running.
This status does currently not include checks whether the database connection is working.

If the service supports TLS Edge Termination, this endpoint does not require the
`X-Forwarded-Proto` header to be set.

Be aware that if you are running multiple nodes of this service, the health status will never
refer to the cluster state, only to a single instance.

#### Responses

<a id="check-alive-status-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|healthStatus|[healthStatus](#schemahealthstatus)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "status": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-isInstanceAlive">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-isInstanceAlive-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceAlive-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceAlive-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceAlive-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceAlive-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-isInstanceAlive-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-isInstanceAlive-shell">

```shell
curl -X GET /health/alive \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceAlive-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/health/alive", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceAlive-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/health/alive', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceAlive-java">

```java
// This sample needs improvement.
URL obj = new URL("/health/alive");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceAlive-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/health/alive',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-isInstanceAlive-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/health/alive',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetJsonWebKeySet"></a>

### Retrieve a JSON Web Key Set

```
GET /keys/{set} HTTP/1.1
Accept: application/json

```

This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.

A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

<a id="retrieve-a-json-web-key-set-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|set|path|string|true|The set|

#### Responses

<a id="retrieve-a-json-web-key-set-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|JSONWebKeySet|[JSONWebKeySet](#schemajsonwebkeyset)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getJsonWebKeySet">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getJsonWebKeySet-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKeySet-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKeySet-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKeySet-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKeySet-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKeySet-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getJsonWebKeySet-shell">

```shell
curl -X GET /keys/{set} \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKeySet-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/keys/{set}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKeySet-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/keys/{set}', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKeySet-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKeySet-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/keys/{set}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKeySet-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/keys/{set}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdupdateJsonWebKeySet"></a>

### Update a JSON Web Key Set

```
PUT /keys/{set} HTTP/1.1
Content-Type: application/json
Accept: application/json

```

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.

A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

#### Request body

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}
```

<a id="update-a-json-web-key-set-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|set|path|string|true|The set|
|body|body|[JSONWebKeySet](#schemajsonwebkeyset)|false|none|

#### Responses

<a id="update-a-json-web-key-set-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|JSONWebKeySet|[JSONWebKeySet](#schemajsonwebkeyset)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-updateJsonWebKeySet">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-updateJsonWebKeySet-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKeySet-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKeySet-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKeySet-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKeySet-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKeySet-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-updateJsonWebKeySet-shell">

```shell
curl -X PUT /keys/{set} \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKeySet-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/keys/{set}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKeySet-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/keys/{set}', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKeySet-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKeySet-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/keys/{set}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKeySet-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/keys/{set}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdcreateJsonWebKeySet"></a>

### Generate a new JSON Web Key

```
POST /keys/{set} HTTP/1.1
Content-Type: application/json
Accept: application/json

```

This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.

A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

#### Request body

```json
{
  "alg": "string",
  "kid": "string",
  "use": "string"
}
```

<a id="generate-a-new-json-web-key-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|set|path|string|true|The set|
|body|body|[jsonWebKeySetGeneratorRequest](#schemajsonwebkeysetgeneratorrequest)|false|none|

#### Responses

<a id="generate-a-new-json-web-key-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|JSONWebKeySet|[JSONWebKeySet](#schemajsonwebkeyset)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 201 response

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-createJsonWebKeySet">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-createJsonWebKeySet-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createJsonWebKeySet-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createJsonWebKeySet-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createJsonWebKeySet-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createJsonWebKeySet-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-createJsonWebKeySet-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-createJsonWebKeySet-shell">

```shell
curl -X POST /keys/{set} \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createJsonWebKeySet-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("POST", "/keys/{set}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createJsonWebKeySet-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "alg": "string",
  "kid": "string",
  "use": "string"
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/keys/{set}', {
  method: 'POST',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createJsonWebKeySet-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createJsonWebKeySet-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post(
  '/keys/{set}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-createJsonWebKeySet-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/keys/{set}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIddeleteJsonWebKeySet"></a>

### Delete a JSON Web Key Set

```
DELETE /keys/{set} HTTP/1.1
Accept: application/json

```

Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.

A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

<a id="delete-a-json-web-key-set-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|set|path|string|true|The set|

#### Responses

<a id="delete-a-json-web-key-set-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 401 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-deleteJsonWebKeySet">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-deleteJsonWebKeySet-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKeySet-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKeySet-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKeySet-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKeySet-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKeySet-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-deleteJsonWebKeySet-shell">

```shell
curl -X DELETE /keys/{set} \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKeySet-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("DELETE", "/keys/{set}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKeySet-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/keys/{set}', {
  method: 'DELETE',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKeySet-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("DELETE");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKeySet-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.delete(
  '/keys/{set}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKeySet-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.delete '/keys/{set}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetJsonWebKey"></a>

### Fetch a JSON Web Key

```
GET /keys/{set}/{kid} HTTP/1.1
Accept: application/json

```

This endpoint returns a singular JSON Web Key, identified by the set and the specific key ID (kid).

<a id="fetch-a-json-web-key-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|kid|path|string|true|The kid of the desired key|
|set|path|string|true|The set|

#### Responses

<a id="fetch-a-json-web-key-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|JSONWebKeySet|[JSONWebKeySet](#schemajsonwebkeyset)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getJsonWebKey">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getJsonWebKey-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKey-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKey-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKey-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKey-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getJsonWebKey-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getJsonWebKey-shell">

```shell
curl -X GET /keys/{set}/{kid} \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKey-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/keys/{set}/{kid}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKey-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/keys/{set}/{kid}', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKey-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}/{kid}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKey-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/keys/{set}/{kid}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getJsonWebKey-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/keys/{set}/{kid}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdupdateJsonWebKey"></a>

### Update a JSON Web Key

```
PUT /keys/{set}/{kid} HTTP/1.1
Content-Type: application/json
Accept: application/json

```

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.

A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

#### Request body

```json
{
  "alg": "RS256",
  "crv": "P-256",
  "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
  "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
  "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
  "e": "AQAB",
  "k": "GawgguFyGrWKav7AX4VKUg",
  "kid": "1603dfe0af8f4596",
  "kty": "RSA",
  "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
  "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
  "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
  "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
  "use": "sig",
  "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
  "x5c": [
    "string"
  ],
  "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
}
```

<a id="update-a-json-web-key-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|kid|path|string|true|The kid of the desired key|
|set|path|string|true|The set|
|body|body|[JSONWebKey](#schemajsonwebkey)|false|none|

#### Responses

<a id="update-a-json-web-key-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|JSONWebKey|[JSONWebKey](#schemajsonwebkey)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "alg": "RS256",
  "crv": "P-256",
  "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
  "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
  "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
  "e": "AQAB",
  "k": "GawgguFyGrWKav7AX4VKUg",
  "kid": "1603dfe0af8f4596",
  "kty": "RSA",
  "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
  "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
  "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
  "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
  "use": "sig",
  "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
  "x5c": [
    "string"
  ],
  "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-updateJsonWebKey">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-updateJsonWebKey-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKey-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKey-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKey-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKey-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-updateJsonWebKey-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-updateJsonWebKey-shell">

```shell
curl -X PUT /keys/{set}/{kid} \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKey-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/keys/{set}/{kid}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKey-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "alg": "RS256",
  "crv": "P-256",
  "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
  "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
  "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
  "e": "AQAB",
  "k": "GawgguFyGrWKav7AX4VKUg",
  "kid": "1603dfe0af8f4596",
  "kty": "RSA",
  "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
  "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
  "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
  "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
  "use": "sig",
  "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
  "x5c": [
    "string"
  ],
  "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/keys/{set}/{kid}', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKey-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}/{kid}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKey-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/keys/{set}/{kid}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-updateJsonWebKey-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/keys/{set}/{kid}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIddeleteJsonWebKey"></a>

### Delete a JSON Web Key

```
DELETE /keys/{set}/{kid} HTTP/1.1
Accept: application/json

```

Use this endpoint to delete a single JSON Web Key.

A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

<a id="delete-a-json-web-key-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|kid|path|string|true|The kid of the desired key|
|set|path|string|true|The set|

#### Responses

<a id="delete-a-json-web-key-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 401 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-deleteJsonWebKey">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-deleteJsonWebKey-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKey-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKey-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKey-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKey-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-deleteJsonWebKey-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-deleteJsonWebKey-shell">

```shell
curl -X DELETE /keys/{set}/{kid} \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKey-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("DELETE", "/keys/{set}/{kid}", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKey-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/keys/{set}/{kid}', {
  method: 'DELETE',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKey-java">

```java
// This sample needs improvement.
URL obj = new URL("/keys/{set}/{kid}");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("DELETE");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKey-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.delete(
  '/keys/{set}/{kid}',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-deleteJsonWebKey-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.delete '/keys/{set}/{kid}',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdprometheus"></a>

### Get snapshot metrics from the Hydra service. If you're using k8s, you can then add annotations to
your deployment like so:

```
GET /metrics/prometheus HTTP/1.1

```

```
metadata:
annotations:
prometheus.io/port: "4445"
prometheus.io/path: "/metrics/prometheus"
```

#### Responses

<a id="get-snapshot-metrics-from-the-hydra-service.-if-you're-using-k8s,-you-can-then-add-annotations-to
your-deployment-like-so:-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-prometheus">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-prometheus-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-prometheus-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-prometheus-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-prometheus-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-prometheus-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-prometheus-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-prometheus-shell">

```shell
curl -X GET /metrics/prometheus

```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-prometheus-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/metrics/prometheus", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-prometheus-node">

```nodejs
const fetch = require('node-fetch');

fetch('/metrics/prometheus', {
  method: 'GET'
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-prometheus-java">

```java
// This sample needs improvement.
URL obj = new URL("/metrics/prometheus");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-prometheus-python">

```python
import requests

r = requests.get(
  '/metrics/prometheus',
  params={)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-prometheus-ruby">

```ruby
require 'rest-client'
require 'json'

result = RestClient.get '/metrics/prometheus',
  params: {}

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetConsentRequest"></a>

### Get consent request information

```
GET /oauth2/auth/requests/consent?consent_challenge=string HTTP/1.1
Accept: application/json

```

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if
the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.

The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to
grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").

The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted
or rejected the request.

<a id="get-consent-request-information-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|consent_challenge|query|string|true|none|

#### Responses

<a id="get-consent-request-information-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|consentRequest|[consentRequest](#schemaconsentrequest)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "acr": "string",
  "challenge": "string",
  "client": {
    "allowed_cors_origins": [
      "string"
    ],
    "audience": [
      "string"
    ],
    "backchannel_logout_session_required": true,
    "backchannel_logout_uri": "string",
    "client_id": "string",
    "client_name": "string",
    "client_secret": "string",
    "client_secret_expires_at": 0,
    "client_uri": "string",
    "contacts": [
      "string"
    ],
    "created_at": "2020-04-25T11:08:35Z",
    "frontchannel_logout_session_required": true,
    "frontchannel_logout_uri": "string",
    "grant_types": [
      "string"
    ],
    "jwks": {},
    "jwks_uri": "string",
    "logo_uri": "string",
    "metadata": {},
    "owner": "string",
    "policy_uri": "string",
    "post_logout_redirect_uris": [
      "string"
    ],
    "redirect_uris": [
      "string"
    ],
    "request_object_signing_alg": "string",
    "request_uris": [
      "string"
    ],
    "response_types": [
      "string"
    ],
    "scope": "string",
    "sector_identifier_uri": "string",
    "subject_type": "string",
    "token_endpoint_auth_method": "string",
    "tos_uri": "string",
    "updated_at": "2020-04-25T11:08:35Z",
    "userinfo_signed_response_alg": "string"
  },
  "context": {},
  "login_challenge": "string",
  "login_session_id": "string",
  "oidc_context": {
    "acr_values": [
      "string"
    ],
    "display": "string",
    "id_token_hint_claims": {
      "property1": {},
      "property2": {}
    },
    "login_hint": "string",
    "ui_locales": [
      "string"
    ]
  },
  "request_url": "string",
  "requested_access_token_audience": [
    "string"
  ],
  "requested_scope": [
    "string"
  ],
  "skip": true,
  "subject": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getConsentRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getConsentRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getConsentRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getConsentRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getConsentRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getConsentRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getConsentRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getConsentRequest-shell">

```shell
curl -X GET /oauth2/auth/requests/consent?consent_challenge=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getConsentRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/oauth2/auth/requests/consent", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getConsentRequest-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/consent?consent_challenge=string', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getConsentRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/consent?consent_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getConsentRequest-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/oauth2/auth/requests/consent',
  params={
    'consent_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getConsentRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/oauth2/auth/requests/consent',
  params: {
    'consent_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdacceptConsentRequest"></a>

### Accept a consent request

```
PUT /oauth2/auth/requests/consent/accept?consent_challenge=string HTTP/1.1
Content-Type: application/json
Accept: application/json

```

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if
the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.

The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to
grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").

The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted
or rejected the request.

This endpoint tells ORY Hydra that the subject has authorized the OAuth 2.0 client to access resources on his/her behalf.
The consent provider includes additional information, such as session data for access and ID tokens, and if the
consent request should be used as basis for future requests.

The response contains a redirect URL which the consent provider should redirect the user-agent to.

#### Request body

```json
{
  "grant_access_token_audience": [
    "string"
  ],
  "grant_scope": [
    "string"
  ],
  "handled_at": "2020-04-25T11:08:35Z",
  "remember": true,
  "remember_for": 0,
  "session": {
    "access_token": {
      "property1": {},
      "property2": {}
    },
    "id_token": {
      "property1": {},
      "property2": {}
    }
  }
}
```

<a id="accept-a-consent-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|consent_challenge|query|string|true|none|
|body|body|[acceptConsentRequest](#schemaacceptconsentrequest)|false|none|

#### Responses

<a id="accept-a-consent-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|completedRequest|[completedRequest](#schemacompletedrequest)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "redirect_to": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-acceptConsentRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-acceptConsentRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptConsentRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptConsentRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptConsentRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptConsentRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptConsentRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-acceptConsentRequest-shell">

```shell
curl -X PUT /oauth2/auth/requests/consent/accept?consent_challenge=string \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptConsentRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/oauth2/auth/requests/consent/accept", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptConsentRequest-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "grant_access_token_audience": [
    "string"
  ],
  "grant_scope": [
    "string"
  ],
  "handled_at": "2020-04-25T11:08:35Z",
  "remember": true,
  "remember_for": 0,
  "session": {
    "access_token": {
      "property1": {},
      "property2": {}
    },
    "id_token": {
      "property1": {},
      "property2": {}
    }
  }
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/consent/accept?consent_challenge=string', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptConsentRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/consent/accept?consent_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptConsentRequest-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/oauth2/auth/requests/consent/accept',
  params={
    'consent_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptConsentRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/oauth2/auth/requests/consent/accept',
  params: {
    'consent_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdrejectConsentRequest"></a>

### Reject a consent request

```
PUT /oauth2/auth/requests/consent/reject?consent_challenge=string HTTP/1.1
Content-Type: application/json
Accept: application/json

```

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if
the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.

The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to
grant or deny the client access to the requested scope ("Application my-dropbox-app wants write access to all your private files").

The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted
or rejected the request.

This endpoint tells ORY Hydra that the subject has not authorized the OAuth 2.0 client to access resources on his/her behalf.
The consent provider must include a reason why the consent was not granted.

The response contains a redirect URL which the consent provider should redirect the user-agent to.

#### Request body

```json
{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}
```

<a id="reject-a-consent-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|consent_challenge|query|string|true|none|
|body|body|[rejectRequest](#schemarejectrequest)|false|none|

#### Responses

<a id="reject-a-consent-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|completedRequest|[completedRequest](#schemacompletedrequest)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "redirect_to": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-rejectConsentRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-rejectConsentRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectConsentRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectConsentRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectConsentRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectConsentRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectConsentRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-rejectConsentRequest-shell">

```shell
curl -X PUT /oauth2/auth/requests/consent/reject?consent_challenge=string \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectConsentRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/oauth2/auth/requests/consent/reject", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectConsentRequest-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/consent/reject?consent_challenge=string', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectConsentRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/consent/reject?consent_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectConsentRequest-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/oauth2/auth/requests/consent/reject',
  params={
    'consent_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectConsentRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/oauth2/auth/requests/consent/reject',
  params: {
    'consent_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetLoginRequest"></a>

### Get a login request

```
GET /oauth2/auth/requests/login?login_challenge=string HTTP/1.1
Accept: application/json

```

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
(sometimes called "identity provider") to authenticate the subject and then tell ORY Hydra now about it. The login
provider is an web-app you write and host, and it must be able to authenticate ("show the subject a login screen")
a subject (in OAuth2 the proper name for subject is "resource owner").

The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

<a id="get-a-login-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|login_challenge|query|string|true|none|

#### Responses

<a id="get-a-login-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|loginRequest|[loginRequest](#schemaloginrequest)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|genericError|[genericError](#schemagenericerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "challenge": "string",
  "client": {
    "allowed_cors_origins": [
      "string"
    ],
    "audience": [
      "string"
    ],
    "backchannel_logout_session_required": true,
    "backchannel_logout_uri": "string",
    "client_id": "string",
    "client_name": "string",
    "client_secret": "string",
    "client_secret_expires_at": 0,
    "client_uri": "string",
    "contacts": [
      "string"
    ],
    "created_at": "2020-04-25T11:08:35Z",
    "frontchannel_logout_session_required": true,
    "frontchannel_logout_uri": "string",
    "grant_types": [
      "string"
    ],
    "jwks": {},
    "jwks_uri": "string",
    "logo_uri": "string",
    "metadata": {},
    "owner": "string",
    "policy_uri": "string",
    "post_logout_redirect_uris": [
      "string"
    ],
    "redirect_uris": [
      "string"
    ],
    "request_object_signing_alg": "string",
    "request_uris": [
      "string"
    ],
    "response_types": [
      "string"
    ],
    "scope": "string",
    "sector_identifier_uri": "string",
    "subject_type": "string",
    "token_endpoint_auth_method": "string",
    "tos_uri": "string",
    "updated_at": "2020-04-25T11:08:35Z",
    "userinfo_signed_response_alg": "string"
  },
  "oidc_context": {
    "acr_values": [
      "string"
    ],
    "display": "string",
    "id_token_hint_claims": {
      "property1": {},
      "property2": {}
    },
    "login_hint": "string",
    "ui_locales": [
      "string"
    ]
  },
  "request_url": "string",
  "requested_access_token_audience": [
    "string"
  ],
  "requested_scope": [
    "string"
  ],
  "session_id": "string",
  "skip": true,
  "subject": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getLoginRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getLoginRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLoginRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLoginRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLoginRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLoginRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLoginRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getLoginRequest-shell">

```shell
curl -X GET /oauth2/auth/requests/login?login_challenge=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLoginRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/oauth2/auth/requests/login", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLoginRequest-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/login?login_challenge=string', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLoginRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/login?login_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLoginRequest-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/oauth2/auth/requests/login',
  params={
    'login_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLoginRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/oauth2/auth/requests/login',
  params: {
    'login_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdacceptLoginRequest"></a>

### Accept a login request

```
PUT /oauth2/auth/requests/login/accept?login_challenge=string HTTP/1.1
Content-Type: application/json
Accept: application/json

```

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
(sometimes called "identity provider") to authenticate the subject and then tell ORY Hydra now about it. The login
provider is an web-app you write and host, and it must be able to authenticate ("show the subject a login screen")
a subject (in OAuth2 the proper name for subject is "resource owner").

The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

This endpoint tells ORY Hydra that the subject has successfully authenticated and includes additional information such as
the subject's ID and if ORY Hydra should remember the subject's subject agent for future authentication attempts by setting
a cookie.

The response contains a redirect URL which the login provider should redirect the user-agent to.

#### Request body

```json
{
  "acr": "string",
  "context": {},
  "force_subject_identifier": "string",
  "remember": true,
  "remember_for": 0,
  "subject": "string"
}
```

<a id="accept-a-login-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|login_challenge|query|string|true|none|
|body|body|[acceptLoginRequest](#schemaacceptloginrequest)|false|none|

#### Responses

<a id="accept-a-login-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|completedRequest|[completedRequest](#schemacompletedrequest)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "redirect_to": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-acceptLoginRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-acceptLoginRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLoginRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLoginRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLoginRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLoginRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLoginRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-acceptLoginRequest-shell">

```shell
curl -X PUT /oauth2/auth/requests/login/accept?login_challenge=string \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLoginRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/oauth2/auth/requests/login/accept", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLoginRequest-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "acr": "string",
  "context": {},
  "force_subject_identifier": "string",
  "remember": true,
  "remember_for": 0,
  "subject": "string"
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/login/accept?login_challenge=string', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLoginRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/login/accept?login_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLoginRequest-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/oauth2/auth/requests/login/accept',
  params={
    'login_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLoginRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/oauth2/auth/requests/login/accept',
  params: {
    'login_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdrejectLoginRequest"></a>

### Reject a login request

```
PUT /oauth2/auth/requests/login/reject?login_challenge=string HTTP/1.1
Content-Type: application/json
Accept: application/json

```

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider
(sometimes called "identity provider") to authenticate the subject and then tell ORY Hydra now about it. The login
provider is an web-app you write and host, and it must be able to authenticate ("show the subject a login screen")
a subject (in OAuth2 the proper name for subject is "resource owner").

The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

This endpoint tells ORY Hydra that the subject has not authenticated and includes a reason why the authentication
was be denied.

The response contains a redirect URL which the login provider should redirect the user-agent to.

#### Request body

```json
{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}
```

<a id="reject-a-login-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|login_challenge|query|string|true|none|
|body|body|[rejectRequest](#schemarejectrequest)|false|none|

#### Responses

<a id="reject-a-login-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|completedRequest|[completedRequest](#schemacompletedrequest)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "redirect_to": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-rejectLoginRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-rejectLoginRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLoginRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLoginRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLoginRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLoginRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLoginRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-rejectLoginRequest-shell">

```shell
curl -X PUT /oauth2/auth/requests/login/reject?login_challenge=string \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLoginRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/oauth2/auth/requests/login/reject", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLoginRequest-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/login/reject?login_challenge=string', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLoginRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/login/reject?login_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLoginRequest-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/oauth2/auth/requests/login/reject',
  params={
    'login_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLoginRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/oauth2/auth/requests/login/reject',
  params: {
    'login_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetLogoutRequest"></a>

### Get a logout request

```
GET /oauth2/auth/requests/logout?logout_challenge=string HTTP/1.1
Accept: application/json

```

Use this endpoint to fetch a logout request.

<a id="get-a-logout-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|logout_challenge|query|string|true|none|

#### Responses

<a id="get-a-logout-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|logoutRequest|[logoutRequest](#schemalogoutrequest)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "request_url": "string",
  "rp_initiated": true,
  "sid": "string",
  "subject": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getLogoutRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getLogoutRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLogoutRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLogoutRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLogoutRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLogoutRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getLogoutRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getLogoutRequest-shell">

```shell
curl -X GET /oauth2/auth/requests/logout?logout_challenge=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLogoutRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/oauth2/auth/requests/logout", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLogoutRequest-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/logout?logout_challenge=string', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLogoutRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/logout?logout_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLogoutRequest-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/oauth2/auth/requests/logout',
  params={
    'logout_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getLogoutRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/oauth2/auth/requests/logout',
  params: {
    'logout_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdacceptLogoutRequest"></a>

### Accept a logout request

```
PUT /oauth2/auth/requests/logout/accept?logout_challenge=string HTTP/1.1
Accept: application/json

```

When a user or an application requests ORY Hydra to log out a user, this endpoint is used to confirm that logout request.
No body is required.

The response contains a redirect URL which the consent provider should redirect the user-agent to.

<a id="accept-a-logout-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|logout_challenge|query|string|true|none|

#### Responses

<a id="accept-a-logout-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|completedRequest|[completedRequest](#schemacompletedrequest)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "redirect_to": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-acceptLogoutRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-acceptLogoutRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLogoutRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLogoutRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLogoutRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLogoutRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-acceptLogoutRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-acceptLogoutRequest-shell">

```shell
curl -X PUT /oauth2/auth/requests/logout/accept?logout_challenge=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLogoutRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/oauth2/auth/requests/logout/accept", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLogoutRequest-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/logout/accept?logout_challenge=string', {
  method: 'PUT',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLogoutRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/logout/accept?logout_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLogoutRequest-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.put(
  '/oauth2/auth/requests/logout/accept',
  params={
    'logout_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-acceptLogoutRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.put '/oauth2/auth/requests/logout/accept',
  params: {
    'logout_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdrejectLogoutRequest"></a>

### Reject a logout request

```
PUT /oauth2/auth/requests/logout/reject?logout_challenge=string HTTP/1.1
Content-Type: application/json
Accept: application/json

```

When a user or an application requests ORY Hydra to log out a user, this endpoint is used to deny that logout request.
No body is required.

The response is empty as the logout provider has to chose what action to perform next.

#### Request body

```json
{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}
```

```yaml
error: string
error_debug: string
error_description: string
error_hint: string
status_code: 0

```

<a id="reject-a-logout-request-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|logout_challenge|query|string|true|none|
|body|body|[rejectRequest](#schemarejectrequest)|false|none|

#### Responses

<a id="reject-a-logout-request-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 404 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-rejectLogoutRequest">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-rejectLogoutRequest-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLogoutRequest-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLogoutRequest-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLogoutRequest-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLogoutRequest-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-rejectLogoutRequest-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-rejectLogoutRequest-shell">

```shell
curl -X PUT /oauth2/auth/requests/logout/reject?logout_challenge=string \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLogoutRequest-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("PUT", "/oauth2/auth/requests/logout/reject", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLogoutRequest-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/oauth2/auth/requests/logout/reject?logout_challenge=string', {
  method: 'PUT',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLogoutRequest-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/requests/logout/reject?logout_challenge=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLogoutRequest-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put(
  '/oauth2/auth/requests/logout/reject',
  params={
    'logout_challenge': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-rejectLogoutRequest-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/oauth2/auth/requests/logout/reject',
  params: {
    'logout_challenge' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdlistSubjectConsentSessions"></a>

### Lists all consent sessions of a subject

```
GET /oauth2/auth/sessions/consent?subject=string HTTP/1.1
Accept: application/json

```

This endpoint lists all subject's granted consent sessions, including client and granted scope.
The "Link" header is also included in successful responses, which contains one or more links for pagination, formatted like so: '<https://hydra-url/admin/oauth2/auth/sessions/consent?subject={user}&limit={limit}&offset={offset}>; rel="{page}"', where page is one of the following applicable pages: 'first', 'next', 'last', and 'previous'.
Multiple links can be included in this header, and will be separated by a comma.

<a id="lists-all-consent-sessions-of-a-subject-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|subject|query|string|true|none|

#### Responses

<a id="lists-all-consent-sessions-of-a-subject-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A list of used consent requests.|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|genericError|[genericError](#schemagenericerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

<a id="lists-all-consent-sessions-of-a-subject-responseschema"></a>
##### Response Schema

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[PreviousConsentSession](#schemapreviousconsentsession)]|false|none|[The response used to return used consent requests same as HandledLoginRequest, just with consent_request exposed as json]|
|» consent_request|[consentRequest](#schemaconsentrequest)|false|none|none|
|»» acr|string|false|none|ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication.|
|»» challenge|string|false|none|Challenge is the identifier ("authorization challenge") of the consent authorization request. It is used to identify the session.|
|»» client|[oAuth2Client](#schemaoauth2client)|false|none|none|
|»»» allowed_cors_origins|[string]|false|none|none|
|»»» audience|[string]|false|none|none|
|»»» backchannel_logout_session_required|boolean|false|none|Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false.|
|»»» backchannel_logout_uri|string|false|none|RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.|
|»»» client_id|string|false|none|ClientID  is the id for this client.|
|»»» client_name|string|false|none|Name is the human-readable string name of the client to be presented to the end-user during authorization.|
|»»» client_secret|string|false|none|Secret is the client's secret. The secret will be included in the create request as cleartext, and then never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users that they need to write the secret down as it will not be made available again.|
|»»» client_secret_expires_at|integer(int64)|false|none|SecretExpiresAt is an integer holding the time at which the client secret will expire or 0 if it will not expire. The time is represented as the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time of expiration.  This feature is currently not supported and it's value will always be set to 0.|
|»»» client_uri|string|false|none|ClientURI is an URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion.|
|»»» contacts|[string]|false|none|none|
|»»» created_at|string(date-time)|false|none|CreatedAt returns the timestamp of the client's creation.|
|»»» frontchannel_logout_session_required|boolean|false|none|Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false.|
|»»» frontchannel_logout_uri|string|false|none|RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the request and to determine which of the potentially multiple sessions is to be logged out; if either is included, both MUST be.|
|»»» grant_types|[string]|false|none|none|
|»»» jwks|[JoseJSONWebKeySet](#schemajosejsonwebkeyset)|false|none|none|
|»»» jwks_uri|string|false|none|URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.|
|»»» logo_uri|string|false|none|LogoURI is an URL string that references a logo for the client.|
|»»» metadata|[JSONRawMessage](#schemajsonrawmessage)|false|none|none|
|»»» owner|string|false|none|Owner is a string identifying the owner of the OAuth 2.0 Client.|
|»»» policy_uri|string|false|none|PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.|
|»»» post_logout_redirect_uris|[string]|false|none|none|
|»»» redirect_uris|[string]|false|none|none|
|»»» request_object_signing_alg|string|false|none|JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm.|
|»»» request_uris|[string]|false|none|none|
|»»» response_types|[string]|false|none|none|
|»»» scope|string|false|none|Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens.|
|»»» sector_identifier_uri|string|false|none|URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.|
|»»» subject_type|string|false|none|SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.|
|»»» token_endpoint_auth_method|string|false|none|Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none.|
|»»» tos_uri|string|false|none|TermsOfServiceURI is a URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.|
|»»» updated_at|string(date-time)|false|none|UpdatedAt returns the timestamp of the last update.|
|»»» userinfo_signed_response_alg|string|false|none|JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type.|
|»» context|[JSONRawMessage](#schemajsonrawmessage)|false|none|none|
|»» login_challenge|string|false|none|LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate a login and consent request in the login & consent app.|
|»» login_session_id|string|false|none|LoginSessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag) this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false) this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back- channel logout. It's value can generally be used to associate consecutive login requests by a certain user.|
|»» oidc_context|[openIDConnectContext](#schemaopenidconnectcontext)|false|none|none|
|»»» acr_values|[string]|false|none|ACRValues is the Authentication AuthorizationContext Class Reference requested in the OAuth 2.0 Authorization request. It is a parameter defined by OpenID Connect and expresses which level of authentication (e.g. 2FA) is required.  OpenID Connect defines it as follows: > Requested Authentication AuthorizationContext Class Reference values. Space-separated string that specifies the acr values that the Authorization Server is being requested to use for processing this Authentication Request, with the values appearing in order of preference. The Authentication AuthorizationContext Class satisfied by the authentication performed is returned as the acr Claim Value, as specified in Section 2. The acr Claim is requested as a Voluntary Claim by this parameter.|
|»»» display|string|false|none|Display is a string value that specifies how the Authorization Server displays the authentication and consent user interface pages to the End-User. The defined values are: page: The Authorization Server SHOULD display the authentication and consent UI consistent with a full User Agent page view. If the display parameter is not specified, this is the default display mode. popup: The Authorization Server SHOULD display the authentication and consent UI consistent with a popup User Agent window. The popup User Agent window should be of an appropriate size for a login-focused dialog and should not obscure the entire window that it is popping up over. touch: The Authorization Server SHOULD display the authentication and consent UI consistent with a device that leverages a touch interface. wap: The Authorization Server SHOULD display the authentication and consent UI consistent with a "feature phone" type display.  The Authorization Server MAY also attempt to detect the capabilities of the User Agent and present an appropriate display.|
|»»» id_token_hint_claims|object|false|none|IDTokenHintClaims are the claims of the ID Token previously issued by the Authorization Server being passed as a hint about the End-User's current or past authenticated session with the Client.|
|»»»» **additionalProperties**|object|false|none|none|
|»»» login_hint|string|false|none|LoginHint hints about the login identifier the End-User might use to log in (if necessary). This hint can be used by an RP if it first asks the End-User for their e-mail address (or other identifier) and then wants to pass that value as a hint to the discovered authorization service. This value MAY also be a phone number in the format specified for the phone_number Claim. The use of this parameter is optional.|
|»»» ui_locales|[string]|false|none|UILocales is the End-User'id preferred languages and scripts for the user interface, represented as a space-separated list of BCP47 [RFC5646] language tag values, ordered by preference. For instance, the value "fr-CA fr en" represents a preference for French as spoken in Canada, then French (without a region designation), followed by English (without a region designation). An error SHOULD NOT result if some or all of the requested locales are not supported by the OpenID Provider.|
|»» request_url|string|false|none|RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters.|
|»» requested_access_token_audience|[string]|false|none|none|
|»» requested_scope|[string]|false|none|none|
|»» skip|boolean|false|none|Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call.|
|»» subject|string|false|none|Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client.|
|» grant_access_token_audience|[string]|false|none|none|
|» grant_scope|[string]|false|none|none|
|» handled_at|[NullTime](#schemanulltime)(date-time)|false|none|none|
|» remember|boolean|false|none|Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope.|
|» remember_for|integer(int64)|false|none|RememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the authorization will be remembered indefinitely.|
|» session|[consentRequestSession](#schemaconsentrequestsession)|false|none|none|
|»» access_token|object|false|none|AccessToken sets session data for the access and refresh token, as well as any future tokens issued by the refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection. If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care!|
|»»» **additionalProperties**|object|false|none|none|
|»» id_token|object|false|none|IDToken sets session data for the OpenID Connect ID token. Keep in mind that the session'id payloads are readable by anyone that has access to the ID Challenge. Use with care!|
|»»» **additionalProperties**|object|false|none|none|

##### Examples

###### 200 response

```json
[
  {
    "consent_request": {
      "acr": "string",
      "challenge": "string",
      "client": {
        "allowed_cors_origins": [
          "string"
        ],
        "audience": [
          "string"
        ],
        "backchannel_logout_session_required": true,
        "backchannel_logout_uri": "string",
        "client_id": "string",
        "client_name": "string",
        "client_secret": "string",
        "client_secret_expires_at": 0,
        "client_uri": "string",
        "contacts": [
          "string"
        ],
        "created_at": "2020-04-25T11:08:35Z",
        "frontchannel_logout_session_required": true,
        "frontchannel_logout_uri": "string",
        "grant_types": [
          "string"
        ],
        "jwks": {},
        "jwks_uri": "string",
        "logo_uri": "string",
        "metadata": {},
        "owner": "string",
        "policy_uri": "string",
        "post_logout_redirect_uris": [
          "string"
        ],
        "redirect_uris": [
          "string"
        ],
        "request_object_signing_alg": "string",
        "request_uris": [
          "string"
        ],
        "response_types": [
          "string"
        ],
        "scope": "string",
        "sector_identifier_uri": "string",
        "subject_type": "string",
        "token_endpoint_auth_method": "string",
        "tos_uri": "string",
        "updated_at": "2020-04-25T11:08:35Z",
        "userinfo_signed_response_alg": "string"
      },
      "context": {},
      "login_challenge": "string",
      "login_session_id": "string",
      "oidc_context": {
        "acr_values": [
          "string"
        ],
        "display": "string",
        "id_token_hint_claims": {
          "property1": {},
          "property2": {}
        },
        "login_hint": "string",
        "ui_locales": [
          "string"
        ]
      },
      "request_url": "string",
      "requested_access_token_audience": [
        "string"
      ],
      "requested_scope": [
        "string"
      ],
      "skip": true,
      "subject": "string"
    },
    "grant_access_token_audience": [
      "string"
    ],
    "grant_scope": [
      "string"
    ],
    "handled_at": "2020-04-25T11:08:35Z",
    "remember": true,
    "remember_for": 0,
    "session": {
      "access_token": {
        "property1": {},
        "property2": {}
      },
      "id_token": {
        "property1": {},
        "property2": {}
      }
    }
  }
]
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-listSubjectConsentSessions">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-listSubjectConsentSessions-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listSubjectConsentSessions-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listSubjectConsentSessions-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listSubjectConsentSessions-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listSubjectConsentSessions-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-listSubjectConsentSessions-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-listSubjectConsentSessions-shell">

```shell
curl -X GET /oauth2/auth/sessions/consent?subject=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listSubjectConsentSessions-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/oauth2/auth/sessions/consent", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listSubjectConsentSessions-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/sessions/consent?subject=string', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listSubjectConsentSessions-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/sessions/consent?subject=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listSubjectConsentSessions-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/oauth2/auth/sessions/consent',
  params={
    'subject': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-listSubjectConsentSessions-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/oauth2/auth/sessions/consent',
  params: {
    'subject' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdrevokeConsentSessions"></a>

### Revokes consent sessions of a subject for a specific OAuth 2.0 Client

```
DELETE /oauth2/auth/sessions/consent?subject=string HTTP/1.1
Accept: application/json

```

This endpoint revokes a subject's granted consent sessions for a specific OAuth 2.0 Client and invalidates all
associated OAuth 2.0 Access Tokens.

<a id="revokes-consent-sessions-of-a-subject-for-a-specific-oauth-2.0-client-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|subject|query|string|true|The subject (Subject) who's consent sessions should be deleted.|
|client|query|string|false|If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID|

#### Responses

<a id="revokes-consent-sessions-of-a-subject-for-a-specific-oauth-2.0-client-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|genericError|[genericError](#schemagenericerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 400 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-revokeConsentSessions">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-revokeConsentSessions-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeConsentSessions-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeConsentSessions-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeConsentSessions-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeConsentSessions-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeConsentSessions-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-revokeConsentSessions-shell">

```shell
curl -X DELETE /oauth2/auth/sessions/consent?subject=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeConsentSessions-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("DELETE", "/oauth2/auth/sessions/consent", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeConsentSessions-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/sessions/consent?subject=string', {
  method: 'DELETE',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeConsentSessions-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/sessions/consent?subject=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("DELETE");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeConsentSessions-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.delete(
  '/oauth2/auth/sessions/consent',
  params={
    'subject': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeConsentSessions-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.delete '/oauth2/auth/sessions/consent',
  params: {
    'subject' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdrevokeAuthenticationSession"></a>

### Invalidates all login sessions of a certain user
Invalidates a subject's authentication session

```
DELETE /oauth2/auth/sessions/login?subject=string HTTP/1.1
Accept: application/json

```

This endpoint invalidates a subject's authentication session. After revoking the authentication session, the subject
has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens and does not work with OpenID Connect
Front- or Back-channel logout.

<a id="invalidates-all-login-sessions-of-a-certain-user
invalidates-a-subject's-authentication-session-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|subject|query|string|true|none|

#### Responses

<a id="invalidates-all-login-sessions-of-a-certain-user
invalidates-a-subject's-authentication-session-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|genericError|[genericError](#schemagenericerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 400 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-revokeAuthenticationSession">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-revokeAuthenticationSession-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeAuthenticationSession-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeAuthenticationSession-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeAuthenticationSession-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeAuthenticationSession-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-revokeAuthenticationSession-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-revokeAuthenticationSession-shell">

```shell
curl -X DELETE /oauth2/auth/sessions/login?subject=string \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeAuthenticationSession-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("DELETE", "/oauth2/auth/sessions/login", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeAuthenticationSession-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/oauth2/auth/sessions/login?subject=string', {
  method: 'DELETE',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeAuthenticationSession-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/auth/sessions/login?subject=string");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("DELETE");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeAuthenticationSession-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.delete(
  '/oauth2/auth/sessions/login',
  params={
    'subject': 'string'},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-revokeAuthenticationSession-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.delete '/oauth2/auth/sessions/login',
  params: {
    'subject' => 'string'}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdflushInactiveOAuth2Tokens"></a>

### Flush Expired OAuth2 Access Tokens

```
POST /oauth2/flush HTTP/1.1
Content-Type: application/json
Accept: application/json

```

This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be
not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted
automatically when performing the refresh flow.

#### Request body

```json
{
  "notAfter": "2020-04-25T11:08:35Z"
}
```

<a id="flush-expired-oauth2-access-tokens-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[flushInactiveOAuth2TokensRequest](#schemaflushinactiveoauth2tokensrequest)|false|none|

#### Responses

<a id="flush-expired-oauth2-access-tokens-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
typically 201.|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 401 response

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-flushInactiveOAuth2Tokens">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-flushInactiveOAuth2Tokens-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-flushInactiveOAuth2Tokens-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-flushInactiveOAuth2Tokens-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-flushInactiveOAuth2Tokens-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-flushInactiveOAuth2Tokens-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-flushInactiveOAuth2Tokens-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-flushInactiveOAuth2Tokens-shell">

```shell
curl -X POST /oauth2/flush \
  -H 'Content-Type: application/json' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-flushInactiveOAuth2Tokens-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("POST", "/oauth2/flush", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-flushInactiveOAuth2Tokens-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "notAfter": "2020-04-25T11:08:35Z"
}';
const headers = {
  'Content-Type': 'application/json',  'Accept': 'application/json'
}

fetch('/oauth2/flush', {
  method: 'POST',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-flushInactiveOAuth2Tokens-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/flush");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-flushInactiveOAuth2Tokens-python">

```python
import requests

headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post(
  '/oauth2/flush',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-flushInactiveOAuth2Tokens-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/oauth2/flush',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdintrospectOAuth2Token"></a>

### Introspect OAuth2 tokens

```
POST /oauth2/introspect HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Accept: application/json

```

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token
is neither expired nor revoked. If a token is active, additional information on the token will be included. You can
set additional data for a token by setting `accessTokenExtra` during the consent flow.

For more information [read this blog post](https://www.oauth.com/oauth2-servers/token-introspection-endpoint/).

#### Request body

```yaml
token: string
scope: string

```

<a id="introspect-oauth2-tokens-parameters"></a>
##### Parameters

|Parameter|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» token|body|string|true|The string value of the token. For access tokens, this|
|» scope|body|string|false|An optional, space separated list of required scopes. If the access token was not granted one of the|

##### Detailed descriptions

**» token**: The string value of the token. For access tokens, this
is the "access_token" value returned from the token endpoint
defined in OAuth 2.0. For refresh tokens, this is the "refresh_token"
value returned.

**» scope**: An optional, space separated list of required scopes. If the access token was not granted one of the
scopes, the result of active will be false.

#### Responses

<a id="introspect-oauth2-tokens-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|oAuth2TokenIntrospection|[oAuth2TokenIntrospection](#schemaoauth2tokenintrospection)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|genericError|[genericError](#schemagenericerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|genericError|[genericError](#schemagenericerror)|

##### Examples

###### 200 response

```json
{
  "active": true,
  "aud": [
    "string"
  ],
  "client_id": "string",
  "exp": 0,
  "ext": {
    "property1": {},
    "property2": {}
  },
  "iat": 0,
  "iss": "string",
  "nbf": 0,
  "obfuscated_subject": "string",
  "scope": "string",
  "sub": "string",
  "token_type": "string",
  "username": "string"
}
```

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
basic, oauth2
</aside>

#### Code samples

<div class="tabs" id="tab-introspectOAuth2Token">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-introspectOAuth2Token-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-introspectOAuth2Token-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-introspectOAuth2Token-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-introspectOAuth2Token-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-introspectOAuth2Token-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-introspectOAuth2Token-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-introspectOAuth2Token-shell">

```shell
curl -X POST /oauth2/introspect \
  -H 'Content-Type: application/x-www-form-urlencoded' \  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-introspectOAuth2Token-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Content-Type": []string{"application/x-www-form-urlencoded"},
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("POST", "/oauth2/introspect", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-introspectOAuth2Token-node">

```nodejs
const fetch = require('node-fetch');
const input = '{
  "token": "string",
  "scope": "string"
}';
const headers = {
  'Content-Type': 'application/x-www-form-urlencoded',  'Accept': 'application/json'
}

fetch('/oauth2/introspect', {
  method: 'POST',
  body: input,
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-introspectOAuth2Token-java">

```java
// This sample needs improvement.
URL obj = new URL("/oauth2/introspect");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-introspectOAuth2Token-python">

```python
import requests

headers = {
  'Content-Type': 'application/x-www-form-urlencoded',
  'Accept': 'application/json'
}

r = requests.post(
  '/oauth2/introspect',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-introspectOAuth2Token-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/x-www-form-urlencoded',
  'Accept' => 'application/json'
}

result = RestClient.post '/oauth2/introspect',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

<a id="opIdgetVersion"></a>

### Get service version

```
GET /version HTTP/1.1
Accept: application/json

```

This endpoint returns the service version typically notated using semantic versioning.

If the service supports TLS Edge Termination, this endpoint does not require the
`X-Forwarded-Proto` header to be set.

#### Responses

<a id="get-service-version-responses"></a>
##### Overview

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|version|[version](#schemaversion)|

##### Examples

###### 200 response

```json
{
  "version": "string"
}
```

<aside class="success">
This operation does not require authentication
</aside>

#### Code samples

<div class="tabs" id="tab-getVersion">
<nav class="tabs-nav">
<ul class="nav nav-tabs au-link-list au-link-list--inline">
<li class="nav-item"><a class="nav-link active" role="tab" href="#tab-getVersion-shell">Shell</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getVersion-go">Go</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getVersion-node">Node.js</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getVersion-java">Java</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getVersion-python">Python</a></li>
<li class="nav-item"><a class="nav-link" role="tab" href="#tab-getVersion-ruby">Ruby</a></li>
</ul>
</nav>
<div class="tab-content">
<div class="tab-pane active" role="tabpanel" id="tab-getVersion-shell">

```shell
curl -X GET /version \
  -H 'Accept: application/json'
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getVersion-go">

```go
package main

import (
    "bytes"
    "net/http"
)

func main() {
    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    var body []byte
    // body = ...

    req, err := http.NewRequest("GET", "/version", bytes.NewBuffer(body))
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getVersion-node">

```nodejs
const fetch = require('node-fetch');

const headers = {
  'Accept': 'application/json'
}

fetch('/version', {
  method: 'GET',
  headers
})
.then(r => r.json())
.then((body) => {
    console.log(body)
})
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getVersion-java">

```java
// This sample needs improvement.
URL obj = new URL("/version");

HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");

int responseCode = con.getResponseCode();

BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream())
);

String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();

System.out.println(response.toString());
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getVersion-python">

```python
import requests

headers = {
  'Accept': 'application/json'
}

r = requests.get(
  '/version',
  params={},
  headers = headers)

print r.json()
```

</div>
<div class="tab-pane" role="tabpanel"  id="tab-getVersion-ruby">

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/version',
  params: {}, headers: headers

p JSON.parse(result)
```

</div>
</div>
</div>

## Schemas

<a id="tocSjsonrawmessage">JSONRawMessage</a>
#### JSONRawMessage

<a id="schemajsonrawmessage"></a>

```json
{}

```

*JSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger.*

#### Properties

*None*

<a id="tocSjsonwebkey">JSONWebKey</a>
#### JSONWebKey

<a id="schemajsonwebkey"></a>

```json
{
  "alg": "RS256",
  "crv": "P-256",
  "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
  "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
  "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
  "e": "AQAB",
  "k": "GawgguFyGrWKav7AX4VKUg",
  "kid": "1603dfe0af8f4596",
  "kty": "RSA",
  "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
  "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
  "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
  "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
  "use": "sig",
  "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
  "x5c": [
    "string"
  ],
  "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
}

```

*It is important that this model object is named JSONWebKey for
"swagger generate spec" to generate only on definition of a
JSONWebKey.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|alg|string|true|none|The "alg" (algorithm) parameter identifies the algorithm intended for use with the key.  The values used should either be registered in the IANA "JSON Web Signature and Encryption Algorithms" registry established by [JWA] or be a value that contains a Collision- Resistant Name.|
|crv|string|false|none|none|
|d|string|false|none|none|
|dp|string|false|none|none|
|dq|string|false|none|none|
|e|string|false|none|none|
|k|string|false|none|none|
|kid|string|true|none|The "kid" (key ID) parameter is used to match a specific key.  This is used, for instance, to choose among a set of keys within a JWK Set during key rollover.  The structure of the "kid" value is unspecified.  When "kid" values are used within a JWK Set, different keys within the JWK Set SHOULD use distinct "kid" values.  (One example in which different keys might use the same "kid" value is if they have different "kty" (key type) values but are considered to be equivalent alternatives by the application using them.)  The "kid" value is a case-sensitive string.|
|kty|string|true|none|The "kty" (key type) parameter identifies the cryptographic algorithm family used with the key, such as "RSA" or "EC". "kty" values should either be registered in the IANA "JSON Web Key Types" registry established by [JWA] or be a value that contains a Collision- Resistant Name.  The "kty" value is a case-sensitive string.|
|n|string|false|none|none|
|p|string|false|none|none|
|q|string|false|none|none|
|qi|string|false|none|none|
|use|string|true|none|Use ("public key use") identifies the intended use of the public key. The "use" parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Values are commonly "sig" (signature) or "enc" (encryption).|
|x|string|false|none|none|
|x5c|[string]|false|none|The "x5c" (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates [RFC5280].  The certificate chain is represented as a JSON array of certificate value strings.  Each string in the array is a base64-encoded (Section 4 of [RFC4648] -- not base64url-encoded) DER [ITU.X690.1994] PKIX certificate value. The PKIX certificate containing the key value MUST be the first certificate.|
|y|string|false|none|none|

<a id="tocSjsonwebkeyset">JSONWebKeySet</a>
#### JSONWebKeySet

<a id="schemajsonwebkeyset"></a>

```json
{
  "keys": [
    {
      "alg": "RS256",
      "crv": "P-256",
      "d": "T_N8I-6He3M8a7X1vWt6TGIx4xB_GP3Mb4SsZSA4v-orvJzzRiQhLlRR81naWYxfQAYt5isDI6_C2L9bdWo4FFPjGQFvNoRX-_sBJyBI_rl-TBgsZYoUlAj3J92WmY2inbA-PwyJfsaIIDceYBC-eX-xiCu6qMqkZi3MwQAFL6bMdPEM0z4JBcwFT3VdiWAIRUuACWQwrXMq672x7fMuaIaHi7XDGgt1ith23CLfaREmJku9PQcchbt_uEY-hqrFY6ntTtS4paWWQj86xLL94S-Tf6v6xkL918PfLSOTq6XCzxvlFwzBJqApnAhbwqLjpPhgUG04EDRrqrSBc5Y1BLevn6Ip5h1AhessBp3wLkQgz_roeckt-ybvzKTjESMuagnpqLvOT7Y9veIug2MwPJZI2VjczRc1vzMs25XrFQ8DpUy-bNdp89TmvAXwctUMiJdgHloJw23Cv03gIUAkDnsTqZmkpbIf-crpgNKFmQP_EDKoe8p_PXZZgfbRri3NoEVGP7Mk6yEu8LjJhClhZaBNjuWw2-KlBfOA3g79mhfBnkInee5KO9mGR50qPk1V-MorUYNTFMZIm0kFE6eYVWFBwJHLKYhHU34DoiK1VP-svZpC2uAMFNA_UJEwM9CQ2b8qe4-5e9aywMvwcuArRkAB5mBIfOaOJao3mfukKAE",
      "dp": "G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
      "dq": "s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
      "e": "AQAB",
      "k": "GawgguFyGrWKav7AX4VKUg",
      "kid": "1603dfe0af8f4596",
      "kty": "RSA",
      "n": "vTqrxUyQPl_20aqf5kXHwDZrel-KovIp8s7ewJod2EXHl8tWlRB3_Rem34KwBfqlKQGp1nqah-51H4Jzruqe0cFP58hPEIt6WqrvnmJCXxnNuIB53iX_uUUXXHDHBeaPCSRoNJzNysjoJ30TIUsKBiirhBa7f235PXbKiHducLevV6PcKxJ5cY8zO286qJLBWSPm-OIevwqsIsSIH44Qtm9sioFikhkbLwoqwWORGAY0nl6XvVOlhADdLjBSqSAeT1FPuCDCnXwzCDR8N9IFB_IjdStFkC-rVt2K5BYfPd0c3yFp_vHR15eRd0zJ8XQ7woBC8Vnsac6Et1pKS59pX6256DPWu8UDdEOolKAPgcd_g2NpA76cAaF_jcT80j9KrEzw8Tv0nJBGesuCjPNjGs_KzdkWTUXt23Hn9QJsdc1MZuaW0iqXBepHYfYoqNelzVte117t4BwVp0kUM6we0IqyXClaZgOI8S-WDBw2_Ovdm8e5NmhYAblEVoygcX8Y46oH6bKiaCQfKCFDMcRgChme7AoE1yZZYsPbaG_3IjPrC4LBMHQw8rM9dWjJ8ImjicvZ1pAm0dx-KHCP3y5PVKrxBDf1zSOsBRkOSjB8TPODnJMz6-jd5hTtZxpZPwPoIdCanTZ3ZD6uRBpTmDwtpRGm63UQs1m5FWPwb0T2IF0",
      "p": "6NbkXwDWUhi-eR55Cgbf27FkQDDWIamOaDr0rj1q0f1fFEz1W5A_09YvG09Fiv1AO2-D8Rl8gS1Vkz2i0zCSqnyy8A025XOcRviOMK7nIxE4OH_PEsko8dtIrb3TmE2hUXvCkmzw9EsTF1LQBOGC6iusLTXepIC1x9ukCKFZQvdgtEObQ5kzd9Nhq-cdqmSeMVLoxPLd1blviVT9Vm8-y12CtYpeJHOaIDtVPLlBhJiBoPKWg3vxSm4XxIliNOefqegIlsmTIa3MpS6WWlCK3yHhat0Q-rRxDxdyiVdG_wzJvp0Iw_2wms7pe-PgNPYvUWH9JphWP5K38YqEBiJFXQ",
      "q": "0A1FmpOWR91_RAWpqreWSavNaZb9nXeKiBo0DQGBz32DbqKqQ8S4aBJmbRhJcctjCLjain-ivut477tAUMmzJwVJDDq2MZFwC9Q-4VYZmFU4HJityQuSzHYe64RjN-E_NQ02TWhG3QGW6roq6c57c99rrUsETwJJiwS8M5p15Miuz53DaOjv-uqqFAFfywN5WkxHbraBcjHtMiQuyQbQqkCFh-oanHkwYNeytsNhTu2mQmwR5DR2roZ2nPiFjC6nsdk-A7E3S3wMzYYFw7jvbWWoYWo9vB40_MY2Y0FYQSqcDzcBIcq_0tnnasf3VW4Fdx6m80RzOb2Fsnln7vKXAQ",
      "qi": "GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
      "use": "sig",
      "x": "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
      "x5c": [
        "string"
      ],
      "y": "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0"
    }
  ]
}

```

*It is important that this model object is named JSONWebKeySet for
"swagger generate spec" to generate only on definition of a
JSONWebKeySet. Since one with the same name is previously defined as
client.Client.JSONWebKeys and this one is last, this one will be
effectively written in the swagger spec.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|keys|[[JSONWebKey](#schemajsonwebkey)]|false|none|The value of the "keys" parameter is an array of JWK values.  By default, the order of the JWK values within the array does not imply an order of preference among them, although applications of JWK Sets can choose to assign a meaning to the order for their purposes, if desired.|

<a id="tocSjosejsonwebkeyset">JoseJSONWebKeySet</a>
#### JoseJSONWebKeySet

<a id="schemajosejsonwebkeyset"></a>

```json
{}

```

#### Properties

*None*

<a id="tocSnulltime">NullTime</a>
#### NullTime

<a id="schemanulltime"></a>

```json
"2020-04-25T11:08:35Z"

```

*NullTime implements sql.NullTime functionality.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|NullTime implements sql.NullTime functionality.|string(date-time)|false|none|none|

<a id="tocSpreviousconsentsession">PreviousConsentSession</a>
#### PreviousConsentSession

<a id="schemapreviousconsentsession"></a>

```json
{
  "consent_request": {
    "acr": "string",
    "challenge": "string",
    "client": {
      "allowed_cors_origins": [
        "string"
      ],
      "audience": [
        "string"
      ],
      "backchannel_logout_session_required": true,
      "backchannel_logout_uri": "string",
      "client_id": "string",
      "client_name": "string",
      "client_secret": "string",
      "client_secret_expires_at": 0,
      "client_uri": "string",
      "contacts": [
        "string"
      ],
      "created_at": "2020-04-25T11:08:35Z",
      "frontchannel_logout_session_required": true,
      "frontchannel_logout_uri": "string",
      "grant_types": [
        "string"
      ],
      "jwks": {},
      "jwks_uri": "string",
      "logo_uri": "string",
      "metadata": {},
      "owner": "string",
      "policy_uri": "string",
      "post_logout_redirect_uris": [
        "string"
      ],
      "redirect_uris": [
        "string"
      ],
      "request_object_signing_alg": "string",
      "request_uris": [
        "string"
      ],
      "response_types": [
        "string"
      ],
      "scope": "string",
      "sector_identifier_uri": "string",
      "subject_type": "string",
      "token_endpoint_auth_method": "string",
      "tos_uri": "string",
      "updated_at": "2020-04-25T11:08:35Z",
      "userinfo_signed_response_alg": "string"
    },
    "context": {},
    "login_challenge": "string",
    "login_session_id": "string",
    "oidc_context": {
      "acr_values": [
        "string"
      ],
      "display": "string",
      "id_token_hint_claims": {
        "property1": {},
        "property2": {}
      },
      "login_hint": "string",
      "ui_locales": [
        "string"
      ]
    },
    "request_url": "string",
    "requested_access_token_audience": [
      "string"
    ],
    "requested_scope": [
      "string"
    ],
    "skip": true,
    "subject": "string"
  },
  "grant_access_token_audience": [
    "string"
  ],
  "grant_scope": [
    "string"
  ],
  "handled_at": "2020-04-25T11:08:35Z",
  "remember": true,
  "remember_for": 0,
  "session": {
    "access_token": {
      "property1": {},
      "property2": {}
    },
    "id_token": {
      "property1": {},
      "property2": {}
    }
  }
}

```

*The response used to return used consent requests
same as HandledLoginRequest, just with consent_request exposed as json*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|consent_request|[consentRequest](#schemaconsentrequest)|false|none|none|
|grant_access_token_audience|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|grant_scope|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|handled_at|[NullTime](#schemanulltime)|false|none|none|
|remember|boolean|false|none|Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope.|
|remember_for|integer(int64)|false|none|RememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the authorization will be remembered indefinitely.|
|session|[consentRequestSession](#schemaconsentrequestsession)|false|none|none|

<a id="tocSstringslicepipedelimiter">StringSlicePipeDelimiter</a>
#### StringSlicePipeDelimiter

<a id="schemastringslicepipedelimiter"></a>

```json
[
  "string"
]

```

*StringSlicePipeDelimiter de/encodes the string slice to/from a SQL string.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|StringSlicePipeDelimiter de/encodes the string slice to/from a SQL string.|[string]|false|none|none|

<a id="tocSacceptconsentrequest">acceptConsentRequest</a>
#### acceptConsentRequest

<a id="schemaacceptconsentrequest"></a>

```json
{
  "grant_access_token_audience": [
    "string"
  ],
  "grant_scope": [
    "string"
  ],
  "handled_at": "2020-04-25T11:08:35Z",
  "remember": true,
  "remember_for": 0,
  "session": {
    "access_token": {
      "property1": {},
      "property2": {}
    },
    "id_token": {
      "property1": {},
      "property2": {}
    }
  }
}

```

*The request payload used to accept a consent request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|grant_access_token_audience|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|grant_scope|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|handled_at|[NullTime](#schemanulltime)|false|none|none|
|remember|boolean|false|none|Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope.|
|remember_for|integer(int64)|false|none|RememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the authorization will be remembered indefinitely.|
|session|[consentRequestSession](#schemaconsentrequestsession)|false|none|none|

<a id="tocSacceptloginrequest">acceptLoginRequest</a>
#### acceptLoginRequest

<a id="schemaacceptloginrequest"></a>

```json
{
  "acr": "string",
  "context": {},
  "force_subject_identifier": "string",
  "remember": true,
  "remember_for": 0,
  "subject": "string"
}

```

*HandledLoginRequest is the request payload used to accept a login request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|acr|string|false|none|ACR sets the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication.|
|context|[JSONRawMessage](#schemajsonrawmessage)|false|none|none|
|force_subject_identifier|string|false|none|ForceSubjectIdentifier forces the "pairwise" user ID of the end-user that authenticated. The "pairwise" user ID refers to the (Pairwise Identifier Algorithm)[http://openid.net/specs/openid-connect-core-1_0.html#PairwiseAlg] of the OpenID Connect specification. It allows you to set an obfuscated subject ("user") identifier that is unique to the client.  Please note that this changes the user ID on endpoint /userinfo and sub claim of the ID Token. It does not change the sub claim in the OAuth 2.0 Introspection.  Per default, ORY Hydra handles this value with its own algorithm. In case you want to set this yourself you can use this field. Please note that setting this field has no effect if `pairwise` is not configured in ORY Hydra or the OAuth 2.0 Client does not expect a pairwise identifier (set via `subject_type` key in the client's configuration).  Please also be aware that ORY Hydra is unable to properly compute this value during authentication. This implies that you have to compute this value on every authentication process (probably depending on the client ID or some other unique value).  If you fail to compute the proper value, then authentication processes which have id_token_hint set might fail.|
|remember|boolean|false|none|Remember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she will not be asked to log in again.|
|remember_for|integer(int64)|false|none|RememberFor sets how long the authentication should be remembered for in seconds. If set to `0`, the authorization will be remembered for the duration of the browser session (using a session cookie).|
|subject|string|true|none|Subject is the user ID of the end-user that authenticated.|

<a id="tocScompletedrequest">completedRequest</a>
#### completedRequest

<a id="schemacompletedrequest"></a>

```json
{
  "redirect_to": "string"
}

```

*The response payload sent when accepting or rejecting a login or consent request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|redirect_to|string|false|none|RedirectURL is the URL which you should redirect the user to once the authentication process is completed.|

<a id="tocSconsentrequest">consentRequest</a>
#### consentRequest

<a id="schemaconsentrequest"></a>

```json
{
  "acr": "string",
  "challenge": "string",
  "client": {
    "allowed_cors_origins": [
      "string"
    ],
    "audience": [
      "string"
    ],
    "backchannel_logout_session_required": true,
    "backchannel_logout_uri": "string",
    "client_id": "string",
    "client_name": "string",
    "client_secret": "string",
    "client_secret_expires_at": 0,
    "client_uri": "string",
    "contacts": [
      "string"
    ],
    "created_at": "2020-04-25T11:08:35Z",
    "frontchannel_logout_session_required": true,
    "frontchannel_logout_uri": "string",
    "grant_types": [
      "string"
    ],
    "jwks": {},
    "jwks_uri": "string",
    "logo_uri": "string",
    "metadata": {},
    "owner": "string",
    "policy_uri": "string",
    "post_logout_redirect_uris": [
      "string"
    ],
    "redirect_uris": [
      "string"
    ],
    "request_object_signing_alg": "string",
    "request_uris": [
      "string"
    ],
    "response_types": [
      "string"
    ],
    "scope": "string",
    "sector_identifier_uri": "string",
    "subject_type": "string",
    "token_endpoint_auth_method": "string",
    "tos_uri": "string",
    "updated_at": "2020-04-25T11:08:35Z",
    "userinfo_signed_response_alg": "string"
  },
  "context": {},
  "login_challenge": "string",
  "login_session_id": "string",
  "oidc_context": {
    "acr_values": [
      "string"
    ],
    "display": "string",
    "id_token_hint_claims": {
      "property1": {},
      "property2": {}
    },
    "login_hint": "string",
    "ui_locales": [
      "string"
    ]
  },
  "request_url": "string",
  "requested_access_token_audience": [
    "string"
  ],
  "requested_scope": [
    "string"
  ],
  "skip": true,
  "subject": "string"
}

```

*Contains information on an ongoing consent request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|acr|string|false|none|ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication.|
|challenge|string|false|none|Challenge is the identifier ("authorization challenge") of the consent authorization request. It is used to identify the session.|
|client|[oAuth2Client](#schemaoauth2client)|false|none|none|
|context|[JSONRawMessage](#schemajsonrawmessage)|false|none|none|
|login_challenge|string|false|none|LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate a login and consent request in the login & consent app.|
|login_session_id|string|false|none|LoginSessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag) this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false) this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back- channel logout. It's value can generally be used to associate consecutive login requests by a certain user.|
|oidc_context|[openIDConnectContext](#schemaopenidconnectcontext)|false|none|none|
|request_url|string|false|none|RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters.|
|requested_access_token_audience|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|requested_scope|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|skip|boolean|false|none|Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call.|
|subject|string|false|none|Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client.|

<a id="tocSconsentrequestsession">consentRequestSession</a>
#### consentRequestSession

<a id="schemaconsentrequestsession"></a>

```json
{
  "access_token": {
    "property1": {},
    "property2": {}
  },
  "id_token": {
    "property1": {},
    "property2": {}
  }
}

```

*Used to pass session data to a consent request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|access_token|object|false|none|AccessToken sets session data for the access and refresh token, as well as any future tokens issued by the refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection. If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care!|
|» **additionalProperties**|object|false|none|none|
|id_token|object|false|none|IDToken sets session data for the OpenID Connect ID token. Keep in mind that the session'id payloads are readable by anyone that has access to the ID Challenge. Use with care!|
|» **additionalProperties**|object|false|none|none|

<a id="tocSflushinactiveoauth2tokensrequest">flushInactiveOAuth2TokensRequest</a>
#### flushInactiveOAuth2TokensRequest

<a id="schemaflushinactiveoauth2tokensrequest"></a>

```json
{
  "notAfter": "2020-04-25T11:08:35Z"
}

```

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|notAfter|string(date-time)|false|none|NotAfter sets after which point tokens should not be flushed. This is useful when you want to keep a history of recently issued tokens for auditing.|

<a id="tocSgenericerror">genericError</a>
#### genericError

<a id="schemagenericerror"></a>

```json
{
  "debug": "The database adapter was unable to find the element",
  "error": "The requested resource could not be found",
  "error_description": "Object with ID 12345 does not exist",
  "status_code": 404
}

```

*Error response*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|debug|string|false|none|Debug contains debug information. This is usually not available and has to be enabled.|
|error|string|true|none|Name is the error name.|
|error_description|string|false|none|Description contains further information on the nature of the error.|
|status_code|integer(int64)|false|none|Code represents the error status code (404, 403, 401, ...).|

<a id="tocShealthnotreadystatus">healthNotReadyStatus</a>
#### healthNotReadyStatus

<a id="schemahealthnotreadystatus"></a>

```json
{
  "errors": {
    "property1": "string",
    "property2": "string"
  }
}

```

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|errors|object|false|none|Errors contains a list of errors that caused the not ready status.|
|» **additionalProperties**|string|false|none|none|

<a id="tocShealthstatus">healthStatus</a>
#### healthStatus

<a id="schemahealthstatus"></a>

```json
{
  "status": "string"
}

```

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|status|string|false|none|Status always contains "ok".|

<a id="tocSjsonwebkeysetgeneratorrequest">jsonWebKeySetGeneratorRequest</a>
#### jsonWebKeySetGeneratorRequest

<a id="schemajsonwebkeysetgeneratorrequest"></a>

```json
{
  "alg": "string",
  "kid": "string",
  "use": "string"
}

```

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|alg|string|true|none|The algorithm to be used for creating the key. Supports "RS256", "ES512", "HS512", and "HS256"|
|kid|string|true|none|The kid of the key to be created|
|use|string|true|none|The "use" (public key use) parameter identifies the intended use of the public key. The "use" parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Valid values are "enc" and "sig".|

<a id="tocSloginrequest">loginRequest</a>
#### loginRequest

<a id="schemaloginrequest"></a>

```json
{
  "challenge": "string",
  "client": {
    "allowed_cors_origins": [
      "string"
    ],
    "audience": [
      "string"
    ],
    "backchannel_logout_session_required": true,
    "backchannel_logout_uri": "string",
    "client_id": "string",
    "client_name": "string",
    "client_secret": "string",
    "client_secret_expires_at": 0,
    "client_uri": "string",
    "contacts": [
      "string"
    ],
    "created_at": "2020-04-25T11:08:35Z",
    "frontchannel_logout_session_required": true,
    "frontchannel_logout_uri": "string",
    "grant_types": [
      "string"
    ],
    "jwks": {},
    "jwks_uri": "string",
    "logo_uri": "string",
    "metadata": {},
    "owner": "string",
    "policy_uri": "string",
    "post_logout_redirect_uris": [
      "string"
    ],
    "redirect_uris": [
      "string"
    ],
    "request_object_signing_alg": "string",
    "request_uris": [
      "string"
    ],
    "response_types": [
      "string"
    ],
    "scope": "string",
    "sector_identifier_uri": "string",
    "subject_type": "string",
    "token_endpoint_auth_method": "string",
    "tos_uri": "string",
    "updated_at": "2020-04-25T11:08:35Z",
    "userinfo_signed_response_alg": "string"
  },
  "oidc_context": {
    "acr_values": [
      "string"
    ],
    "display": "string",
    "id_token_hint_claims": {
      "property1": {},
      "property2": {}
    },
    "login_hint": "string",
    "ui_locales": [
      "string"
    ]
  },
  "request_url": "string",
  "requested_access_token_audience": [
    "string"
  ],
  "requested_scope": [
    "string"
  ],
  "session_id": "string",
  "skip": true,
  "subject": "string"
}

```

*Contains information on an ongoing login request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|challenge|string|false|none|Challenge is the identifier ("login challenge") of the login request. It is used to identify the session.|
|client|[oAuth2Client](#schemaoauth2client)|false|none|none|
|oidc_context|[openIDConnectContext](#schemaopenidconnectcontext)|false|none|none|
|request_url|string|false|none|RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters.|
|requested_access_token_audience|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|requested_scope|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|session_id|string|false|none|SessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag) this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false) this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back- channel logout. It's value can generally be used to associate consecutive login requests by a certain user.|
|skip|boolean|false|none|Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.  This feature allows you to update / set session information.|
|subject|string|false|none|Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. If this value is set and `skip` is true, you MUST include this subject type when accepting the login request, or the request will fail.|

<a id="tocSlogoutrequest">logoutRequest</a>
#### logoutRequest

<a id="schemalogoutrequest"></a>

```json
{
  "request_url": "string",
  "rp_initiated": true,
  "sid": "string",
  "subject": "string"
}

```

*Contains information about an ongoing logout request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|request_url|string|false|none|RequestURL is the original Logout URL requested.|
|rp_initiated|boolean|false|none|RPInitiated is set to true if the request was initiated by a Relying Party (RP), also known as an OAuth 2.0 Client.|
|sid|string|false|none|SessionID is the login session ID that was requested to log out.|
|subject|string|false|none|Subject is the user for whom the logout was request.|

<a id="tocSoauth2client">oAuth2Client</a>
#### oAuth2Client

<a id="schemaoauth2client"></a>

```json
{
  "allowed_cors_origins": [
    "string"
  ],
  "audience": [
    "string"
  ],
  "backchannel_logout_session_required": true,
  "backchannel_logout_uri": "string",
  "client_id": "string",
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "created_at": "2020-04-25T11:08:35Z",
  "frontchannel_logout_session_required": true,
  "frontchannel_logout_uri": "string",
  "grant_types": [
    "string"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "metadata": {},
  "owner": "string",
  "policy_uri": "string",
  "post_logout_redirect_uris": [
    "string"
  ],
  "redirect_uris": [
    "string"
  ],
  "request_object_signing_alg": "string",
  "request_uris": [
    "string"
  ],
  "response_types": [
    "string"
  ],
  "scope": "string",
  "sector_identifier_uri": "string",
  "subject_type": "string",
  "token_endpoint_auth_method": "string",
  "tos_uri": "string",
  "updated_at": "2020-04-25T11:08:35Z",
  "userinfo_signed_response_alg": "string"
}

```

*Client represents an OAuth 2.0 Client.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|allowed_cors_origins|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|audience|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|backchannel_logout_session_required|boolean|false|none|Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false.|
|backchannel_logout_uri|string|false|none|RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.|
|client_id|string|false|none|ClientID  is the id for this client.|
|client_name|string|false|none|Name is the human-readable string name of the client to be presented to the end-user during authorization.|
|client_secret|string|false|none|Secret is the client's secret. The secret will be included in the create request as cleartext, and then never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users that they need to write the secret down as it will not be made available again.|
|client_secret_expires_at|integer(int64)|false|none|SecretExpiresAt is an integer holding the time at which the client secret will expire or 0 if it will not expire. The time is represented as the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time of expiration.  This feature is currently not supported and it's value will always be set to 0.|
|client_uri|string|false|none|ClientURI is an URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion.|
|contacts|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|created_at|string(date-time)|false|none|CreatedAt returns the timestamp of the client's creation.|
|frontchannel_logout_session_required|boolean|false|none|Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false.|
|frontchannel_logout_uri|string|false|none|RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the request and to determine which of the potentially multiple sessions is to be logged out; if either is included, both MUST be.|
|grant_types|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|jwks|[JoseJSONWebKeySet](#schemajosejsonwebkeyset)|false|none|none|
|jwks_uri|string|false|none|URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.|
|logo_uri|string|false|none|LogoURI is an URL string that references a logo for the client.|
|metadata|[JSONRawMessage](#schemajsonrawmessage)|false|none|none|
|owner|string|false|none|Owner is a string identifying the owner of the OAuth 2.0 Client.|
|policy_uri|string|false|none|PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.|
|post_logout_redirect_uris|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|redirect_uris|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|request_object_signing_alg|string|false|none|JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm.|
|request_uris|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|response_types|[StringSlicePipeDelimiter](#schemastringslicepipedelimiter)|false|none|none|
|scope|string|false|none|Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens.|
|sector_identifier_uri|string|false|none|URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.|
|subject_type|string|false|none|SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.|
|token_endpoint_auth_method|string|false|none|Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none.|
|tos_uri|string|false|none|TermsOfServiceURI is a URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.|
|updated_at|string(date-time)|false|none|UpdatedAt returns the timestamp of the last update.|
|userinfo_signed_response_alg|string|false|none|JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type.|

<a id="tocSoauth2tokenintrospection">oAuth2TokenIntrospection</a>
#### oAuth2TokenIntrospection

<a id="schemaoauth2tokenintrospection"></a>

```json
{
  "active": true,
  "aud": [
    "string"
  ],
  "client_id": "string",
  "exp": 0,
  "ext": {
    "property1": {},
    "property2": {}
  },
  "iat": 0,
  "iss": "string",
  "nbf": 0,
  "obfuscated_subject": "string",
  "scope": "string",
  "sub": "string",
  "token_type": "string",
  "username": "string"
}

```

*Introspection contains an access token's session data as specified by IETF RFC 7662, see:*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|active|boolean|true|none|Active is a boolean indicator of whether or not the presented token is currently active.  The specifics of a token's "active" state will vary depending on the implementation of the authorization server and the information it keeps about its tokens, but a "true" value return for the "active" property will generally indicate that a given token has been issued by this authorization server, has not been revoked by the resource owner, and is within its given time window of validity (e.g., after its issuance time and before its expiration time).|
|aud|[string]|false|none|Audience contains a list of the token's intended audiences.|
|client_id|string|false|none|ClientID is aclient identifier for the OAuth 2.0 client that requested this token.|
|exp|integer(int64)|false|none|Expires at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire.|
|ext|object|false|none|Extra is arbitrary data set by the session.|
|» **additionalProperties**|object|false|none|none|
|iat|integer(int64)|false|none|Issued at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued.|
|iss|string|false|none|IssuerURL is a string representing the issuer of this token|
|nbf|integer(int64)|false|none|NotBefore is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token is not to be used before.|
|obfuscated_subject|string|false|none|ObfuscatedSubject is set when the subject identifier algorithm was set to "pairwise" during authorization. It is the `sub` value of the ID Token that was issued.|
|scope|string|false|none|Scope is a JSON string containing a space-separated list of scopes associated with this token.|
|sub|string|false|none|Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token.|
|token_type|string|false|none|TokenType is the introspected token's type, for example `access_token` or `refresh_token`.|
|username|string|false|none|Username is a human-readable identifier for the resource owner who authorized this token.|

<a id="tocSoauth2tokenresponse">oauth2TokenResponse</a>
#### oauth2TokenResponse

<a id="schemaoauth2tokenresponse"></a>

```json
{
  "access_token": "string",
  "expires_in": 0,
  "id_token": "string",
  "refresh_token": "string",
  "scope": "string",
  "token_type": "string"
}

```

*The Access Token Response*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|access_token|string|false|none|none|
|expires_in|integer(int64)|false|none|none|
|id_token|string|false|none|none|
|refresh_token|string|false|none|none|
|scope|string|false|none|none|
|token_type|string|false|none|none|

<a id="tocSopenidconnectcontext">openIDConnectContext</a>
#### openIDConnectContext

<a id="schemaopenidconnectcontext"></a>

```json
{
  "acr_values": [
    "string"
  ],
  "display": "string",
  "id_token_hint_claims": {
    "property1": {},
    "property2": {}
  },
  "login_hint": "string",
  "ui_locales": [
    "string"
  ]
}

```

*Contains optional information about the OpenID Connect request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|acr_values|[string]|false|none|ACRValues is the Authentication AuthorizationContext Class Reference requested in the OAuth 2.0 Authorization request. It is a parameter defined by OpenID Connect and expresses which level of authentication (e.g. 2FA) is required.  OpenID Connect defines it as follows: > Requested Authentication AuthorizationContext Class Reference values. Space-separated string that specifies the acr values that the Authorization Server is being requested to use for processing this Authentication Request, with the values appearing in order of preference. The Authentication AuthorizationContext Class satisfied by the authentication performed is returned as the acr Claim Value, as specified in Section 2. The acr Claim is requested as a Voluntary Claim by this parameter.|
|display|string|false|none|Display is a string value that specifies how the Authorization Server displays the authentication and consent user interface pages to the End-User. The defined values are: page: The Authorization Server SHOULD display the authentication and consent UI consistent with a full User Agent page view. If the display parameter is not specified, this is the default display mode. popup: The Authorization Server SHOULD display the authentication and consent UI consistent with a popup User Agent window. The popup User Agent window should be of an appropriate size for a login-focused dialog and should not obscure the entire window that it is popping up over. touch: The Authorization Server SHOULD display the authentication and consent UI consistent with a device that leverages a touch interface. wap: The Authorization Server SHOULD display the authentication and consent UI consistent with a "feature phone" type display.  The Authorization Server MAY also attempt to detect the capabilities of the User Agent and present an appropriate display.|
|id_token_hint_claims|object|false|none|IDTokenHintClaims are the claims of the ID Token previously issued by the Authorization Server being passed as a hint about the End-User's current or past authenticated session with the Client.|
|» **additionalProperties**|object|false|none|none|
|login_hint|string|false|none|LoginHint hints about the login identifier the End-User might use to log in (if necessary). This hint can be used by an RP if it first asks the End-User for their e-mail address (or other identifier) and then wants to pass that value as a hint to the discovered authorization service. This value MAY also be a phone number in the format specified for the phone_number Claim. The use of this parameter is optional.|
|ui_locales|[string]|false|none|UILocales is the End-User'id preferred languages and scripts for the user interface, represented as a space-separated list of BCP47 [RFC5646] language tag values, ordered by preference. For instance, the value "fr-CA fr en" represents a preference for French as spoken in Canada, then French (without a region designation), followed by English (without a region designation). An error SHOULD NOT result if some or all of the requested locales are not supported by the OpenID Provider.|

<a id="tocSrejectrequest">rejectRequest</a>
#### rejectRequest

<a id="schemarejectrequest"></a>

```json
{
  "error": "string",
  "error_debug": "string",
  "error_description": "string",
  "error_hint": "string",
  "status_code": 0
}

```

*The request payload used to accept a login or consent request.*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|error|string|false|none|none|
|error_debug|string|false|none|none|
|error_description|string|false|none|none|
|error_hint|string|false|none|none|
|status_code|integer(int64)|false|none|none|

<a id="tocSuserinforesponse">userinfoResponse</a>
#### userinfoResponse

<a id="schemauserinforesponse"></a>

```json
{
  "birthdate": "string",
  "email": "string",
  "email_verified": true,
  "family_name": "string",
  "gender": "string",
  "given_name": "string",
  "locale": "string",
  "middle_name": "string",
  "name": "string",
  "nickname": "string",
  "phone_number": "string",
  "phone_number_verified": true,
  "picture": "string",
  "preferred_username": "string",
  "profile": "string",
  "sub": "string",
  "updated_at": 0,
  "website": "string",
  "zoneinfo": "string"
}

```

*The userinfo response*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|birthdate|string|false|none|End-User's birthday, represented as an ISO 8601:2004 [ISO8601‑2004] YYYY-MM-DD format. The year MAY be 0000, indicating that it is omitted. To represent only the year, YYYY format is allowed. Note that depending on the underlying platform's date related function, providing just year can result in varying month and day, so the implementers need to take this factor into account to correctly process the dates.|
|email|string|false|none|End-User's preferred e-mail address. Its value MUST conform to the RFC 5322 [RFC5322] addr-spec syntax. The RP MUST NOT rely upon this value being unique, as discussed in Section 5.7.|
|email_verified|boolean|false|none|True if the End-User's e-mail address has been verified; otherwise false. When this Claim Value is true, this means that the OP took affirmative steps to ensure that this e-mail address was controlled by the End-User at the time the verification was performed. The means by which an e-mail address is verified is context-specific, and dependent upon the trust framework or contractual agreements within which the parties are operating.|
|family_name|string|false|none|Surname(s) or last name(s) of the End-User. Note that in some cultures, people can have multiple family names or no family name; all can be present, with the names being separated by space characters.|
|gender|string|false|none|End-User's gender. Values defined by this specification are female and male. Other values MAY be used when neither of the defined values are applicable.|
|given_name|string|false|none|Given name(s) or first name(s) of the End-User. Note that in some cultures, people can have multiple given names; all can be present, with the names being separated by space characters.|
|locale|string|false|none|End-User's locale, represented as a BCP47 [RFC5646] language tag. This is typically an ISO 639-1 Alpha-2 [ISO639‑1] language code in lowercase and an ISO 3166-1 Alpha-2 [ISO3166‑1] country code in uppercase, separated by a dash. For example, en-US or fr-CA. As a compatibility note, some implementations have used an underscore as the separator rather than a dash, for example, en_US; Relying Parties MAY choose to accept this locale syntax as well.|
|middle_name|string|false|none|Middle name(s) of the End-User. Note that in some cultures, people can have multiple middle names; all can be present, with the names being separated by space characters. Also note that in some cultures, middle names are not used.|
|name|string|false|none|End-User's full name in displayable form including all name parts, possibly including titles and suffixes, ordered according to the End-User's locale and preferences.|
|nickname|string|false|none|Casual name of the End-User that may or may not be the same as the given_name. For instance, a nickname value of Mike might be returned alongside a given_name value of Michael.|
|phone_number|string|false|none|End-User's preferred telephone number. E.164 [E.164] is RECOMMENDED as the format of this Claim, for example, +1 (425) 555-1212 or +56 (2) 687 2400. If the phone number contains an extension, it is RECOMMENDED that the extension be represented using the RFC 3966 [RFC3966] extension syntax, for example, +1 (604) 555-1234;ext=5678.|
|phone_number_verified|boolean|false|none|True if the End-User's phone number has been verified; otherwise false. When this Claim Value is true, this means that the OP took affirmative steps to ensure that this phone number was controlled by the End-User at the time the verification was performed. The means by which a phone number is verified is context-specific, and dependent upon the trust framework or contractual agreements within which the parties are operating. When true, the phone_number Claim MUST be in E.164 format and any extensions MUST be represented in RFC 3966 format.|
|picture|string|false|none|URL of the End-User's profile picture. This URL MUST refer to an image file (for example, a PNG, JPEG, or GIF image file), rather than to a Web page containing an image. Note that this URL SHOULD specifically reference a profile photo of the End-User suitable for displaying when describing the End-User, rather than an arbitrary photo taken by the End-User.|
|preferred_username|string|false|none|Non-unique shorthand name by which the End-User wishes to be referred to at the RP, such as janedoe or j.doe. This value MAY be any valid JSON string including special characters such as @, /, or whitespace.|
|profile|string|false|none|URL of the End-User's profile page. The contents of this Web page SHOULD be about the End-User.|
|sub|string|false|none|Subject - Identifier for the End-User at the IssuerURL.|
|updated_at|integer(int64)|false|none|Time the End-User's information was last updated. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time.|
|website|string|false|none|URL of the End-User's Web page or blog. This Web page SHOULD contain information published by the End-User or an organization that the End-User is affiliated with.|
|zoneinfo|string|false|none|String from zoneinfo [zoneinfo] time zone database representing the End-User's time zone. For example, Europe/Paris or America/Los_Angeles.|

<a id="tocSversion">version</a>
#### version

<a id="schemaversion"></a>

```json
{
  "version": "string"
}

```

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|version|string|false|none|Version is the service's version.|

<a id="tocSwellknown">wellKnown</a>
#### wellKnown

<a id="schemawellknown"></a>

```json
{
  "authorization_endpoint": "https://playground.ory.sh/ory-hydra/public/oauth2/auth",
  "backchannel_logout_session_supported": true,
  "backchannel_logout_supported": true,
  "claims_parameter_supported": true,
  "claims_supported": [
    "string"
  ],
  "end_session_endpoint": "string",
  "frontchannel_logout_session_supported": true,
  "frontchannel_logout_supported": true,
  "grant_types_supported": [
    "string"
  ],
  "id_token_signing_alg_values_supported": [
    "string"
  ],
  "issuer": "https://playground.ory.sh/ory-hydra/public/",
  "jwks_uri": "https://playground.ory.sh/ory-hydra/public/.well-known/jwks.json",
  "registration_endpoint": "https://playground.ory.sh/ory-hydra/admin/client",
  "request_parameter_supported": true,
  "request_uri_parameter_supported": true,
  "require_request_uri_registration": true,
  "response_modes_supported": [
    "string"
  ],
  "response_types_supported": [
    "string"
  ],
  "revocation_endpoint": "string",
  "scopes_supported": [
    "string"
  ],
  "subject_types_supported": [
    "string"
  ],
  "token_endpoint": "https://playground.ory.sh/ory-hydra/public/oauth2/token",
  "token_endpoint_auth_methods_supported": [
    "string"
  ],
  "userinfo_endpoint": "string",
  "userinfo_signing_alg_values_supported": [
    "string"
  ]
}

```

*WellKnown represents important OpenID Connect discovery metadata*

#### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|authorization_endpoint|string|true|none|URL of the OP's OAuth 2.0 Authorization Endpoint.|
|backchannel_logout_session_supported|boolean|false|none|Boolean value specifying whether the OP can pass a sid (session ID) Claim in the Logout Token to identify the RP session with the OP. If supported, the sid Claim is also included in ID Tokens issued by the OP|
|backchannel_logout_supported|boolean|false|none|Boolean value specifying whether the OP supports back-channel logout, with true indicating support.|
|claims_parameter_supported|boolean|false|none|Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support.|
|claims_supported|[string]|false|none|JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list.|
|end_session_endpoint|string|false|none|URL at the OP to which an RP can perform a redirect to request that the End-User be logged out at the OP.|
|frontchannel_logout_session_supported|boolean|false|none|Boolean value specifying whether the OP can pass iss (issuer) and sid (session ID) query parameters to identify the RP session with the OP when the frontchannel_logout_uri is used. If supported, the sid Claim is also included in ID Tokens issued by the OP.|
|frontchannel_logout_supported|boolean|false|none|Boolean value specifying whether the OP supports HTTP-based logout, with true indicating support.|
|grant_types_supported|[string]|false|none|JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports.|
|id_token_signing_alg_values_supported|[string]|true|none|JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT.|
|issuer|string|true|none|URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier. If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL.|
|jwks_uri|string|true|none|URL of the OP's JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server's encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.|
|registration_endpoint|string|false|none|URL of the OP's Dynamic Client Registration Endpoint.|
|request_parameter_supported|boolean|false|none|Boolean value specifying whether the OP supports use of the request parameter, with true indicating support.|
|request_uri_parameter_supported|boolean|false|none|Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support.|
|require_request_uri_registration|boolean|false|none|Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter.|
|response_modes_supported|[string]|false|none|JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports.|
|response_types_supported|[string]|true|none|JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values.|
|revocation_endpoint|string|false|none|URL of the authorization server's OAuth 2.0 revocation endpoint.|
|scopes_supported|[string]|false|none|SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used|
|subject_types_supported|[string]|true|none|JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public.|
|token_endpoint|string|true|none|URL of the OP's OAuth 2.0 Token Endpoint|
|token_endpoint_auth_methods_supported|[string]|false|none|JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0|
|userinfo_endpoint|string|false|none|URL of the OP's UserInfo Endpoint.|
|userinfo_signing_alg_values_supported|[string]|false|none|JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT].|

