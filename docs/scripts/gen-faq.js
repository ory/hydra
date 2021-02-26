// gen-faq.js
// generates faq.mdx and faq.module.css from the contents of faq.yaml. See https://github.com/ory/kratos/pull/1039.
const fs = require('fs')
const yaml = require('js-yaml')
const { Remarkable } = require('remarkable')
const path = require('path')
const yamlPath = path.resolve('./docs/faq.yaml')
const prettier = require('prettier')
const prettierStyles = require('ory-prettier-styles')

// Generating FAQ.mdx

if (!fs.existsSync(yamlPath)) {
  //file exists
  console.warn('.yaml File does not exists, skipping generating FAQ')
  return 0
}

let faqYaml = fs.readFileSync(yamlPath, 'utf8')
let faq = yaml.load(faqYaml)

const tags = Array.from(new Set(faq.map(({ tags }) => tags).flat(1)))

// which project are we running in?
const project = process.env.CIRCLE_PROJECT_REPONAME

let data = `---
id: faq
title: Frequently Asked Questions (FAQ)
---
<!-- This file is generated. Please edit /docs/faq.yaml or /docs/scripts/gen-faq.js instead. Changes will be overwritten otherwise -->



import {Question, Faq} from '@theme/Faq'

<Faq tags={${JSON.stringify(tags)}} switchofftags="${project}"/>
<br/><br/>

`
md = new Remarkable()
faq.forEach((el) => {
  react_tags = el.tags.map((tag) => {
    return tag + '_src-theme-'
  })
  data += `<Question tags="question_src-theme- ${react_tags.join(' ')}">\n`
  data += `    ${el.tags
    .map((tag) => {
      return '#' + tag
    })
    .join(' ')} <br/>\n`
  data += '    ' + md.render(`**Q**: ${el.q}`)
  data += '    ' + md.render(`**A**: ${el.a}\n`)
  if (el.context) {
    data += '    ' + md.render(`context: ${el.context}\n`)
  }
  data += `</Question>\n\n<br/>`
})
// Unfortunatly this is a mix of html/markdown and prettier is either not
// properly formatting html or mixing up the syntax (with the html parser)

fs.writeFileSync(path.resolve('./docs/docs/faq.mdx'), data)

// Generating faq.module.css
const taglist = Array.from(
  new Set(
    faq
      .map((el) => {
        return el.tags
      })
      .flat(1)
  )
)
let css_file = ``

taglist.forEach((tag) => {
  css_file += `
li.selected.${tag} {
    color:red;
}

li.selected.${tag}~.question.${tag} {
    display: inline;
    
}
`
})

fs.writeFileSync('./docs/src/theme/faq.module.gen.css', css_file)
