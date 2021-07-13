const fs = require('fs')
const path = require('path')

let config = {
  projectName: 'ORY Template',
  projectSlug: 'docusaurus-template',
  projectTagLine: 'Stubbydi dab dub dadada',
  newsletter:
    'https://ory.us10.list-manage.com/subscribe?u=ffb1a878e4ec6c0ed312a3480&id=f605a41b53',
  updateTags: [
    {
      image: 'oryd/docusaurus-template',
      files: ['docs/docs/configure-deploy.md']
    }
  ],
  updateConfig: {
    src: '.schema/config.schema.json',
    dst: './docs/docs/reference/configuration.md'
  },
  enableRedoc: true
}

const cn = path.join(__dirname, '..', 'config.js')
if (fs.existsSync(cn)) {
  config = require(cn)
}

module.exports = config
