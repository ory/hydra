const config = require('./contrib/config.js')

const githubRepoName =
  config.projectSlug === 'ecosystem' ? 'docs' : config.projectSlug

const baseUrl = config.baseUrl ? config.baseUrl : `/${config.projectSlug}/docs/`

const links = [
  {
    to: '/',
    activeBasePath: baseUrl,
    label: `Docs`,
    position: 'left'
  },
  {
    href: 'https://www.ory.sh/docs',
    label: 'Ecosystem',
    position: 'left'
  },
  {
    href: 'https://www.ory.sh/blog',
    label: 'Blog',
    position: 'left'
  },
  {
    href: 'https://github.com/ory/${githubRepoName}/discussions',
    label: 'Discussions',
    position: 'left'
  },
  {
    href: 'https://www.ory.sh/chat',
    label: 'Chat',
    position: 'left'
  },
  {
    href: `https://github.com/ory/${githubRepoName}`,
    label: 'GitHub',
    position: 'left'
  }
]

module.exports = {
  title: config.projectName,
  tagline: config.projectTagLine,
  url: `https://www.ory.sh/`,
  baseUrl,
  favicon: 'img/favico.png',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  organizationName: 'ory', // Usually your GitHub org/user name.
  projectName: config.projectSlug, // Usually your repo name.
  themeConfig: {
    prism: {
      theme: require('prism-react-renderer/themes/github'),
      darkTheme: require('prism-react-renderer/themes/dracula'),
      additionalLanguages: ['pug']
    },
    announcementBar: {
      id: 'supportus',
      content:
        config.projectSlug === 'docs'
          ? `Sign up for <a href="${config.newsletter}">important security announcements</a> and if you like the ${config.projectName} give us some ⭐️ on <a target="_blank" rel="noopener noreferrer" href="https://github.com/ory">GitHub</a>!`
          : `Sign up for <a href="${config.newsletter}">important security announcements</a> and if you like ${config.projectName} give it a ⭐️ on <a target="_blank" rel="noopener noreferrer" href="https://github.com/ory/${githubRepoName}">GitHub</a>!`
    },
    algolia: {
      apiKey: '8463c6ece843b377565726bb4ed325b0',
      indexName: 'ory',
      contextualSearch: true,
      searchParameters: {
        facetFilters: [[`tags:${config.projectSlug}`, `tags:docs`]]
      }
    },
    navbar: {
      hideOnScroll: true,
      logo: {
        alt: config.projectName,
        src: `img/logo-${config.projectSlug}.svg`,
        srcDark: `img/logo-${config.projectSlug}.svg`,
        href:
          config.projectSlug === 'docs'
            ? `https://www.ory.sh`
            : `https://www.ory.sh/${config.projectSlug}`
      },
      items: [
        ...links,
        {
          type: 'docsVersionDropdown',
          position: 'right',
          dropdownActiveClassDisabled: true,
          dropdownItemsAfter: [
            {
              to: '/versions',
              label: 'All versions'
            }
          ]
        }
      ]
    },
    footer: {
      style: 'dark',
      copyright: `Copyright © ${new Date().getFullYear()} ORY GmbH`,
      links: [
        {
          title: 'Company',
          items: [
            {
              label: 'Imprint',
              href: 'https://www.ory.sh/imprint'
            },
            {
              label: 'Privacy',
              href: 'https://www.ory.sh/privacy'
            },
            {
              label: 'Terms',
              href: 'https://www.ory.sh/tos'
            }
          ]
        }
      ]
    }
  },
  plugins: [
    [
      '@docusaurus/plugin-content-docs',
      {
        path:
          config.projectSlug === 'docusaurus-template'
            ? 'contrib/docs'
            : 'docs',
        sidebarPath: require.resolve('./contrib/sidebar.js'),
        editUrl: `https://github.com/ory/${githubRepoName}/edit/master/docs`,
        routeBasePath: '/',
        showLastUpdateAuthor: true,
        showLastUpdateTime: true,
        disableVersioning: false
      }
    ],
    '@docusaurus/plugin-content-pages',
    require.resolve('./src/plugins/docusaurus-plugin-matamo'),
    '@docusaurus/plugin-sitemap'
  ],
  themes: [
    [
      '@docusaurus/theme-classic',
      {
        customCss:
          config.projectSlug === 'docusaurus-template'
            ? require.resolve('./contrib/theme.css')
            : require.resolve('./src/css/theme.css')
      }
    ],
    '@docusaurus/theme-search-algolia'
  ]
}
