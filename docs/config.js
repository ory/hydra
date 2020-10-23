module.exports = {
  projectName: 'ORY Hydra',
  projectSlug: 'hydra',
  projectTagLine:
    'A cloud native Identity & Access Proxy / API (IAP) and Access Control Decision API that authenticates, authorizes, and mutates incoming HTTP(s) requests. Inspired by the BeyondCorp / Zero Trust white paper. Written in Go.',
  updateTags: [
    {
      image: 'oryd/hydra',
      files: ['docs/docs/install.md', 'docs/docs/configure-deploy.mdx']
    },
    {
      replacer: ({ content, next, semverRegex }) =>
        content.replace(semverRegex, `${next}`),
      image: 'oryd/hydra',
      files: ['docs/docs/install.md']
    },
    {
      replacer: ({ content, next, semverRegex }) =>
        content.replace(
          /oryd\/hydra:v[0-9a-zA-Z\.\+\_-]+/g,
          `oryd/hydra:${next}-sqlite`
        ),
      files: ['quickstart.yml']
    },
    {
      image: 'oryd/hydra-login-consent-node',
      files: ['quickstart.yml']
    }
  ],
  updateConfig: {
    src: '.schema/config.schema.json',
    dst: './docs/docs/reference/configuration.md'
  }
}
