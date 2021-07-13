const path = require('path')
const name = process.argv[2]
const fs = require('fs')

const p = path.join(__dirname, '../versions.json')

fs.writeFile(
  p,
  JSON.stringify(require(p).filter((v) => v !== name)),
  function (err) {
    if (err) {
      return console.error(err)
    }
  }
)
