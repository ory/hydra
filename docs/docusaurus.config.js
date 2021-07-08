const config = require('./contrib/config.js')
const fs = require('fs')

const githubRepoName =
  config.projectSlug === 'ecosystem' ? 'docs' : config.projectSlug

const baseUrl = config.baseUrl ? config.baseUrl : `/${config.projectSlug}/docs/`

let links = [
  {
    to: 'https://www.ory.sh/',
    label: `Home`,
    position: 'left'
  },
  {
    href: `https://github.com/ory/${githubRepoName}/discussions`,
    label: 'Discussions',
    position: 'right'
  },
  {
    href: 'https://www.ory.sh/chat',
    label: 'Slack',
    position: 'right'
  },
  {
    href: `https://github.com/ory/${githubRepoName}`,
    label: 'GitHub',
    position: 'right'
  }
]

const customCss = [require.resolve('./contrib/theme.css')]

if (fs.existsSync('./src/css/theme.css')) {
  customCss.push(require.resolve('./src/css/theme.css'))
}

const githubPrismTheme = require('prism-react-renderer/themes/github')

const prismThemeLight = {
  ...githubPrismTheme,
  styles: [
    ...githubPrismTheme.styles,
    {
      languages: ['keto-relation-tuples'],
      types: ['namespace'],
      style: {
        color: '#666'
      }
    },
    {
      languages: ['keto-relation-tuples'],
      types: ['object'],
      style: {
        color: '#939'
      }
    },
    {
      languages: ['keto-relation-tuples'],
      types: ['relation'],
      style: {
        color: '#e80'
      }
    },
    {
      languages: ['keto-relation-tuples'],
      types: ['delimiter'],
      style: {
        color: '#555'
      }
    },
    {
      languages: ['keto-relation-tuples'],
      types: ['comment'],
      style: {
        color: '#999'
      }
    },
    {
      languages: ['keto-relation-tuples'],
      types: ['subject'],
      style: {
        color: '#903'
      }
    }
  ]
}

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
      theme: prismThemeLight,
      darkTheme: require('prism-react-renderer/themes/dracula'),
      additionalLanguages: ['pug', 'shell-session']
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
      hideOnScroll: false,
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
        editCurrentVersion: false,
        routeBasePath: '/',
        showLastUpdateAuthor: true,
        showLastUpdateTime: true,
        disableVersioning: false,
        include: ['**/*.md', '**/*.mdx', '**/*.jsx'],
        docLayoutComponent: '@theme/RoutedDocPage'
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
        customCss
      }
    ],
    '@docusaurus/theme-search-algolia',
    'docusaurus-theme-redoc'
  ]
}
