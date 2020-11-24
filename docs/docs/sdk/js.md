---
id: js
title: JavaScript
---

This generator creates TypeScript/JavaScript client that utilizes
[axios](https://github.com/axios/axios). The generated Node module can be used
in the following environments:

Environment

- Node.js
- Webpack
- Browserify

Language level

- ES5 - you must have a Promises/A+ library installed
- ES6

Module system

- CommonJS
- ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition
should be automatically resolved via `package.json`.
([Reference](http://www.typescriptlang.org/docs/handbook/typings-for-npm-packages.html))

### Building

To build and compile the typescript sources to javascript use:

```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following
commands.

_published:_

```
npm install @oryd/hydra-client@v1.9.0-alpha.2 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### API Docs

API docs are available
[here](https://github.com/ory/sdk/blob/master/clients/hydra/typescript/README.md).
Please note that those docs are generated and may introduce bugs if code
examples are used 1:1. Especially the package name is not correct.

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
