const config = require('./config.js');
const fs = require('fs');
const path = require('path');

const projects = [
  {
    slug: 'kratos',
    name: 'ORY Kratos',
  },
  {
    slug: 'hydra',
    name: 'ORY Hydra',
  },
  {
    slug: 'oathkeeper',
    name: 'ORY Oathkeeper',
  },
  {
    slug: 'keto',
    name: 'ORY Keto',
  },
].filter((item) => config.projectSlug !== item.slug);

let sidebar = {
  Welcome: ['index'],
};

const cn = path.join(__dirname, '..', 'sidebar.js');
if (fs.existsSync(cn)) {
  sidebar = require(cn);
}

projects.forEach((item) => {
  sidebar[item.name] = [
    {
      type: 'link',
      label: 'Home',
      href: `https://www.ory.sh/${item.slug}`,
    },
    {
      type: 'link',
      label: 'Docs',
      href: `https://www.ory.sh/${item.slug}/docs`,
    },
    {
      type: 'link',
      label: 'GitHub',
      href: `https://github.com/ory/${item.slug}`,
    },
  ];
});

module.exports = {
  docs: sidebar,
};
