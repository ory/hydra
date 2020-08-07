module.exports = {
  projectName: 'ORY Hydra',
  projectSlug: 'hydra',
  projectTagLine:
    'A cloud native Identity & Access Proxy / API (IAP) and Access Control Decision API that authenticates, authorizes, and mutates incoming HTTP(s) requests. Inspired by the BeyondCorp / Zero Trust white paper. Written in Go.',
  updateTags: [
    {
      image: 'oryd/hydra',
      files: [
        'docs/docs/install.md',
        'docs/docs/configure-deploy.mdx',
        'quickstart.yml',
        'quickstart-cockroach.yml',
        'quickstart-cors.yml',
        'quickstart-debug.yml',
        'quickstart-jwt.yml',
        'quickstart-mysql.yml',
        'quickstart-postgres.yml',
        'quickstart-prometheus.yml',
        'quickstart-tracing.yml'
      ]
    },
    {
      replacer: ({ content, next }) =>
        content.replace(
          /v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?/gi,
          `${next}`
        ),
      files: ['docs/docs/install.md']
    }
  ],
  updateConfig: {
    src: '.schema/config.schema.json',
    dst: './docs/docs/reference/configuration.md'
  }
}
