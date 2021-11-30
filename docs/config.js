module.exports = {
  projectName: 'Ory Hydra',
  projectSlug: 'hydra',
  newsletter:
    'https://ory.us10.list-manage.com/subscribe?u=ffb1a878e4ec6c0ed312a3480&id=f605a41b53&group[17097][8]=1',
  projectTagLine:
    'A cloud native Identity & Access Proxy / API (IAP) and Access Control Decision API that authenticates, authorizes, and mutates incoming HTTP(s) requests. Inspired by the BeyondCorp / Zero Trust white paper. Written in Go.',
  updateTags: [
    {
      image: 'oryd/hydra',
      files: ['docs/docs/install.md', 'docs/docs/configure-deploy.mdx']
    },
    {
      // replace the docker tags
      image: 'oryd/hydra',
      files: ['docs/docs/install.md']
    },
    {
      // replace the bash curl tag
      replacer: ({ content, next }) =>
        content.replace(/v[0-9].[0-9].[0-9][0-9a-zA-Z.+_-]+/gi, `${next}`),
      image: 'oryd/hydra',
      files: ['docs/docs/install.md']
    },
    {
      replacer: ({ content, next }) =>
        content.replace(
          /oryd\/hydra:v[0-9a-zA-Z.+_-]+/gi,
          `oryd/hydra:${next}-sqlite`
        ),
      files: ['quickstart.yml']
    },
    {
      image: 'oryd/hydra-login-consent-node',
      files: ['quickstart.yml']
    },
    {
      image: 'oryd/hydra',
      files: [
        'quickstart-cockroach.yml',
        'quickstart-mysql.yml',
        'quickstart-postgres.yml'
      ]
    }
  ],
  updateConfig: {
    src: '../spec/config.json',
    dst: './docs/reference/configuration.md'
  }
}
