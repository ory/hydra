# Documentation

This directory contains the project's documentation.

## Develop

To change the documentation locally, you need NodeJS installed.
Next, install the dependencies:

```
$ npm
```

### Develop

```
$ npm start
```

This command starts a local development server and open up a browser window. Most changes are reflected live without having to restart the server.

### Build

```
$ npm build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

## Create Documentation

To create a new documentation for a new or existing project, copy the contents (without `node_modules`)
into the project's `./docs` directory. Next you need to create these files (relative to the project's root directory)
directory:

```js
// ./docs/config.js
//
// This file contains the project's name and slug and tag line, for example:
module.exports = {
  projectName: 'ORY Keto',
  projectSlug: 'keto',
  projectTagLine: 'A cloud native access control server providing best-practice patterns (RBAC, ABAC, ACL, AWS IAM Policies, Kubernetes Roles, ...) via REST APIs.',
  updateTags: [
    {
      image: 'oryd/keto',
      files: ['docs/docs/configure-deploy.md']
    }
  ],
  updateConfig: {
    src: '.schema/config.schema.json',
    dst: './docs/docs/reference/configuration.md'
  }
};
```

```js
// ./docs/sidebar.js
//
// This represents the sidebar navigation, for example:
module.exports = {
  Introduction: [
    "index",
    "install",
  ],
};
```

```
// ./docs/src/css/theme.css
// empty file is ok
```

Next, put your markdown files in `./docs/docs`. You may also want to add the CircleCI Orb `ory/docs` to your CI config,
depending on the project type.

## Update Documentation

Check out [docusaurus-template](https://github.com/ory/docusaurus-template) using `git clone git@github.com:ory/docusaurus-template.git docusaurus-template`.
It is important that the directory is named `docusaurus-template`!

Then, make your changes, and run `./update.sh`.
