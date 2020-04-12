const fs = require('fs');
const path = require('path');

const help = `
  usage:
    node docker-tag.js path/to/config.js $CIRCLE_TAG
`;

if (process.argv.length !== 4) {
  if (process.argv[2] === 'help') {
    console.log(help);
    return;
  } else if (process.argv.length === 3) {
    console.log('Skipping because tag is empty');
    return;
  }

  console.error(help);
  process.exit(1);
  return;
}

const config = require(path.resolve(process.argv[2]));
const next = process.argv[3];

const replace = (path, replacer) =>
  new Promise((resolve, reject) => {
    fs.readFile(path, 'utf8', (err, data) => {
      if (err) {
        return reject(err);
      }

      fs.writeFile(path, replacer(data), 'utf8', (err) => {
        if (err) {
          return reject(err);
        }
        resolve();
      });
    });
  });

config.updateTags.forEach(({ files, image, replacer }) => {
  files.forEach((loc) => {
    replace(loc, (content) => {
      if (replacer) {
        return replacer({ content, next });
      }

      return content.replace(
        new RegExp(`${image}:v[0-9a-zA-Z\\.\\+\\_\\-]+`, 'g'),
        `${image}:${next}`
      );
    })
      .then(() => {
        console.log('Done!');
      })
      .catch((err) => {
        console.error(err);
        process.exit(1);
      });
  });
});
