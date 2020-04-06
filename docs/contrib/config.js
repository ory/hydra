const fs = require('fs');
const path = require('path');

let config = {
  projectName: 'ORY Template',
  projectSlug: 'docusaurus-template',
  projectTagLine: 'Stubbydi dab dub dadada',
  updateTags: [
    {
      image: 'oryd/docusaurus-template',
      files: ['docs/docs/configure-deploy.md'],
    },
  ],
  updateConfig: {
    src: '.schema/config.schema.json',
    dst: './docs/docs/reference/configuration.md',
  },
};

const cn = path.join(__dirname, '..', 'config.js');
if (fs.existsSync(cn)) {
  config = require(cn);
}

module.exports = config;
