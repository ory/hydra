const fs = require('fs')
const path = require('path')

const help = `
  usage:
    node docker-tag.js path/to/config.js $CIRCLE_TAG
`

if (process.argv.length !== 4) {
  if (process.argv[2] === 'help') {
    console.log(help)
    return
  } else if (process.argv.length === 3) {
    console.log('Skipping because tag is empty')
    return
  }

  console.error(help)
  process.exit(1)
  return
}

const config = require(path.resolve(process.argv[2]))
const next = process.argv[3]

const replace = (path, replacer) => {
  const content = fs.readFileSync(path, 'utf8')
  const updated = replacer(content)
  fs.unlinkSync(path)
  fs.writeFileSync(path, updated, 'utf8')
}

config.updateTags.forEach(({ files, image, replacer }) => {
  files.forEach((loc) => {
    replace(loc, (content) => {
      if (replacer) {
        return replacer({
          content,
          next,
          semverRegex: /v[0-9]\.[0-9]\.[0-9](-([0-9a-zA-Z.\-]+)|)/gi
        })
      }

      return content.replace(
        new RegExp(`${image}:v[0-9a-zA-Z.+_-]+`, 'gi'),
        `${image}:${next}`
      )
    })
    console.log('Processed file:', loc)
  })
})
