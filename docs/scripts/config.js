const RefParser = require('json-schema-ref-parser')
const parser = new RefParser()
const jsf = require('json-schema-faker').default
const YAML = require('yaml')
const { pathOr } = require('ramda')
const path = require('path')
const fs = require('fs')
const prettier = require('prettier')
const prettierStyles = require('ory-prettier-styles')

jsf.option({
  alwaysFakeOptionals: true,
  useExamplesValue: true,
  useDefaultValue: true,
  minItems: 1,
  random: () => 0
})

if (process.argv.length !== 3 || process.argv[1] === 'help') {
  console.error(`
  usage:
    node config.js path/to/config.js
`)
  return
}

const config = require(path.resolve(process.argv[2]))

const enhance =
  (schema, parents = []) =>
  (item) => {
    const key = item.key.value

    const path = [
      ...parents.map((parent) => ['properties', parent]),
      ['properties', key]
    ].flat()

    if (['title', 'description'].find((f) => path[path.length - 1] === f)) {
      return
    }

    const comments = [`# ${pathOr(key, [...path, 'title'], schema)} ##`, '']

    const description = pathOr('', [...path, 'description'], schema)
    if (description) {
      comments.push(' ' + description.split('\n').join('\n '), '')
    }

    const defaultValue = pathOr('', [...path, 'default'], schema)
    if (defaultValue || defaultValue === false) {
      comments.push(' Default value: ' + defaultValue, '')
    }

    const enums = pathOr('', [...path, 'enum'], schema)
    if (enums && Array.isArray(enums)) {
      comments.push(
        ' One of:',
        ...YAML.stringify(enums)
          .split('\n')
          .map((i) => ` ${i}`)
      ) // split always returns one empty object so no need for newline
    }

    const min = pathOr('', [...path, 'minimum'], schema)
    if (min || min === 0) {
      comments.push(` Minimum value: ${min}`, '')
    }

    const max = pathOr('', [...path, 'maximum'], schema)
    if (max || max === 0) {
      comments.push(` Maximum value: ${max}`, '')
    }

    const examples = pathOr('', [...path, 'examples'], schema)
    if (examples) {
      comments.push(
        ' Examples:',
        ...YAML.stringify(examples)
          .split('\n')
          .map((i) => ` ${i}`)
      ) // split always returns one empty object so no need for newline
    }

    let hasChildren
    if (item.value.items) {
      item.value.items.forEach((item) => {
        if (item.key) {
          enhance(schema, [...parents, key])(item)
          hasChildren = true
        }
      })
    }

    const showEnvVarBlockForObject = pathOr(
      '',
      [...path, 'showEnvVarBlockForObject'],
      schema
    )
    if (!hasChildren || showEnvVarBlockForObject) {
      const env = [...parents, key].map((i) => i.toUpperCase()).join('_')
      comments.push(
        ' Set this value using environment variables on',
        ' - Linux/macOS:',
        `    $ export ${env}=<value>`,
        ' - Windows Command Line (CMD):',
        `    > set ${env}=<value>`,
        ''
      )

      // Show this if the config property is an object, to call out how to specify the env var
      if (hasChildren) {
        comments.push(
          ' This can be set as an environment variable by supplying it as a JSON object.',
          ''
        )
      }
    }

    item.commentBefore = comments.join('\n')
    item.spaceBefore = true
  }

new Promise((resolve, reject) => {
  parser.dereference(
    require(path.resolve(config.updateConfig.src)),
    (err, result) => (err ? reject(err) : resolve(result))
  )
})
  .then((schema) => {
    const removeAdditionalProperties = (o) => {
      delete o['additionalProperties']
      if (o.properties) {
        Object.keys(o.properties).forEach((key) =>
          removeAdditionalProperties(o.properties[key])
        )
      }
    }

    const enableAll = (o) => {
      if (o.properties) {
        Object.keys(o.properties).forEach((key) => {
          if (key === 'enable') {
            o.properties[key] = true
          }
          enableAll(o.properties[key])
        })
      }
    }

    removeAdditionalProperties(schema)
    enableAll(schema)
    if (schema.definitions) {
      Object.keys(schema.definitions).forEach((key) => {
        removeAdditionalProperties(schema.definitions[key])
        enableAll(schema.definitions[key])
      })
    }

    jsf.option({
      useExamplesValue: true,
      useDefaultValue: false, // do not change this!!
      fixedProbabilities: true,
      alwaysFakeOptionals: true
    })

    const values = jsf.generate(schema)
    const doc = YAML.parseDocument(YAML.stringify(values))

    const comments = [`# ${pathOr(config.projectSlug, ['title'], schema)}`, '']

    const description = pathOr('', ['description'], schema)
    if (description) {
      comments.push(' ' + description)
    }

    doc.commentBefore = comments.join('\n')
    doc.spaceAfter = false
    doc.spaceBefore = false

    doc.contents.items.forEach(enhance(schema, []))

    return Promise.resolve({
      // schema,
      // values,
      yaml: doc.toString()
    })
  })
  .then((out) => {
    const content = `---
id: configuration
title: Configuration
---

<!-- THIS FILE IS BEING AUTO-GENERATED. DO NOT MODIFY IT AS ALL CHANGES WILL BE OVERWRITTEN.
OPEN AN ISSUE IF YOU WOULD LIKE TO MAKE ADJUSTMENTS HERE AND MAINTAINERS WILL HELP YOU LOCATE THE RIGHT
FILE -->

If file \`$HOME/.${config.projectSlug}.yaml\` exists, it will be used as a configuration file which supports all
configuration settings listed below.

You can load the config file from another source using the \`-c path/to/config.yaml\` or \`--config path/to/config.yaml\`
flag: \`${config.projectSlug} --config path/to/config.yaml\`.

Config files can be formatted as JSON, YAML and TOML. Some configuration values support reloading without server restart.
All configuration values can be set using environment variables, as documented below.

This reference configuration documents all keys, also deprecated ones!
It is a reference for all possible configuration values.

If you are looking for an example configuration, it is better to try out the quickstart.

To find out more about edge cases like setting string array values through environmental variables head to the
[Configuring ORY services](https://www.ory.sh/docs/ecosystem/configuring) section.

\`\`\`yaml
${out.yaml}
\`\`\``

    return new Promise((resolve, reject) => {
      fs.writeFile(
        path.resolve(config.updateConfig.dst),
        prettier.format(content, { ...prettierStyles, parser: 'markdown' }),
        'utf8',
        (err) => {
          if (err) {
            reject(err)
            return
          }
          resolve()
        }
      )
    })
  })
  .then(() => {
    console.log('Done!')
  })
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })
