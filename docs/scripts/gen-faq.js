// gen-faq.js
// generates faq.mdx and faq.module.css from the contents of faq.yaml. See https://github.com/ory/kratos/pull/1039.
const fs = require('fs')
const yaml = require('js-yaml')
const { Remarkable } = require('remarkable')
const path = require('path')
const yamlPath = path.resolve('./faq.yaml')
const prettier = require('prettier')
const prettierStyles = require('ory-prettier-styles')
const config = require('../contrib/config.js')

// Generating FAQ.mdx

if (!fs.existsSync(yamlPath)) {
  //file exists
  console.warn('faq.yaml File does not exists, skipping generating FAQ')
  return 0
}

const faqYaml = fs.readFileSync(yamlPath, 'utf8')
const faq = yaml.load(faqYaml)

const tags = Array.from(new Set(faq.map(({ tags }) => tags).flat(1)))

// which project are we running in?
const project = config.projectSlug

let markdownPage = `---
id: faq
title: Frequently Asked Questions (FAQ)
---
<!-- This file is generated. Please edit /docs/faq.yaml or /docs/scripts/gen-faq.js instead. Changes will be overwritten otherwise -->



import {Question, FaqTags} from '@theme/Faq'

<FaqTags tags={${JSON.stringify(tags)}} initiallyDisabled={[${JSON.stringify(
  project
)}]}/>
<br/><br/>

`
md = new Remarkable()
faq.forEach((el) => {
  markdownPage += `<Question tags={${JSON.stringify(el.tags)}}>\n`
  markdownPage += `${el.tags
    .map((tag) => {
      return '#' + tag
    })
    .join(' ')}
`
  markdownPage += md.render(`**Q**: ${el.q}`)
  markdownPage += md.render(`**A**: ${el.a}`)
  if (el.context) {
    markdownPage += md.render(`context: ${el.context}`)
  }
  markdownPage += `</Question>

<br/>
`
})

fs.writeFileSync(
  path.resolve('./docs/faq.mdx'),
  prettier.format(markdownPage, { ...prettierStyles, parser: 'mdx' })
)

// Generating faq.module.css
const tagList = Array.from(
  new Set(
    faq
      .map((el) => {
        return el.tags
      })
      .flat(1)
  )
)

let generatedCSS = `
.selected {
  background-color: #ffba00;
}

div.question {
  display: none;
}
`

tagList.forEach((tag) => {
  generatedCSS += `
li.selected.${tag} {
    color:red;
}

li.selected.${tag}~.question.${tag} {
    display: inline;
}
`
})

fs.writeFileSync(
  './src/theme/faq.gen.module.css',
  prettier.format(generatedCSS, { ...prettierStyles, parser: 'css' })
)
