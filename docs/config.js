module.exports = {
  projectName: 'ORY Hydra',
  projectSlug: 'hydra',
  projectTagLine: 'A cloud native Identity & Access Proxy / API (IAP) and Access Control Decision API that authenticates, authorizes, and mutates incoming HTTP(s) requests. Inspired by the BeyondCorp / Zero Trust white paper. Written in Go.',
  updateTags: [
    {
      image: 'oryd/hydra',
      files: [
        'docs/docs/install.md',
        'docs/docs/configure-deploy.mdx'
      ]
    }
  ],
  updateConfig: {
    src: '.schema/config.schema.json',
    dst: './docs/docs/reference/configuration.md'
  }
};
