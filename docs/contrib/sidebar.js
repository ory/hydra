const config = require('./config.js')
const fs = require('fs')
const path = require('path')
const request = require('sync-request')
const parser = require('parser-front-matter')
const base = require(path.join(__dirname, 'sidebar.json'))

let sidebar = {
  Welcome: ['index']
}

const cn = path.join(__dirname, '..', 'sidebar.json')
if (fs.existsSync(cn)) {
  sidebar = require(cn)
}

const toHref = (slug, node) => {
  if (node !== null && typeof node === 'object') {
    if (node.type) {
      if (node.type === 'category') {
        return {
          ...node,
          items: node.items.map((n) => toHref(slug, n))
        }
      }
      return node
    }

    Object.entries(node).forEach(([key, value]) => {
      node[key] = toHref(slug, value)
    })
    return node
  } else if (Array.isArray(node)) {
    return node.map((value) => {
      return toHref(slug, value)
    })
  }

  let res = request(
    'GET',
    `https://raw.githubusercontent.com/ory/${slug}/master/docs/docs/${node}.mdx`
  )
  if (res.statusCode === 404) {
    res = request(
      'GET',
      `https://raw.githubusercontent.com/ory/${slug}/master/docs/docs/${node}.md`
    )
  }

  const fm = parser.parseSync(res.getBody().toString())
  const doc = fm.data.slug || node

  return {
    label: fm.data.title,
    type: 'link',
    href: `https://www.ory.sh/${slug}/${slug !== 'docs' ? 'docs/next/' : ''}${
      doc === '/' ? '' : doc
    }`
  }
}

const resolveRefs = (node) => {
  if (node !== null && typeof node == 'object') {
    if (node['$slug']) {
      const slug = node['$slug']
      if (slug === config.projectSlug) {
        return sidebar
      }

      const res = request(
        'GET',
        `https://raw.githubusercontent.com/ory/${slug}/master/docs/sidebar.json`
      )
      const items = JSON.parse(res.getBody().toString())

      return toHref(slug, items)
    }

    Object.entries(node).forEach(([key, value]) => {
      node[key] = resolveRefs(value)
    })
    return node
  }
  return node
}

const result = resolveRefs(base)

module.exports = {
  docs: resolveRefs(result)
}
