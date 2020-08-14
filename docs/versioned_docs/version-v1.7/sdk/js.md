---
id: js
title: JavaScript
---

To install the JavaScript SDK, run:

```
npm install --save @oryd/hydra-client
```

### Configuration

#### Basic configuration

```js
import { AdminApi } from '@oryd/hydra-client'

// Set this to Hydra's URL
const hydraAdmin = new AdminApi('http://localhost:4445')
```

### API Usage

```js
hydraAdmin.listOAuth2Clients(10, 0).then(({ body }) => {
  body.forEach((client) => {
    console.log(client)
  })
})
```

### API Docs

API docs are available
[here](https://github.com/ory/hydra/blob/master/sdk/js/swagger/README.md).
Please note that those docs are generated and may introduce bugs if code
examples are used 1:1. Especially the package name is not correct.
