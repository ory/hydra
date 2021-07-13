---
id: js
title: JavaScript
---

To install the JavaScript SDK, run:

```
npm install --save @ory/hydra-client
```

### Configuration

#### Basic configuration

```js
import { Configuration, PublicApi, AdminApi } from '@ory/hydra-client'

const hydraPublic = new PublicApi(
  new Configuration({
    basePath: 'https://public.hydra:4444/'
  })
)

const hydraAdmin = new AdminApi(
  new Configuration({
    basePath: 'https://public.hydra:4445/'
  })
)
```

### API Usage

We recommend using TypeScript with auto-completion as API usage is not well
documented currently.
