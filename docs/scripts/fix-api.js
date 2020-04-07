const fs = require('fs');

if (process.argv.length !== 3 || process.argv[1] === 'help') {
  console.error(`
  usage:
    node fix-api.js path/to/file.md
`);
  process.exit(1);
}

const file = process.argv[2];

fs.readFile(file, (err, b) => {
  if (err) {
    throw err;
  }

  const t = b
    .toString()
    .replace(/^id: api/gim, '')
    .replace(/^title:(.*)/im, 'title: REST API\nid: api') // improve title, add docusaurus id
    .replace(/^language_tabs:.*\n/im, '') // not supported by docusaurus
    .replace(/^toc_footers.*\n/im, '') // not supported by docusaurus
    .replace(/^includes.*\n/im, '') // not supported by docusaurus
    .replace(/^search.*\n/im, '') // not supported by docusaurus
    .replace(/^highlight_theme.*\n/im, '') // not supported by docusaurus
    .replace(/^headingLevel.*\n/im, '') // not supported by docusaurus
    // .replace(/^<h1.*\n/im, '') // remove first headline (this is the title of the main package usually)
    // .replace(/^> Scroll down for example requests and responses.*\n/im, '') // Irrelevant information
    // .replace(/^Base Urls:*\n/im, '') // Irrelevant information, let's replace it with something useful instead!
    // .replace(/^\* <a href="\/">\/<\/a>\n/im, '') // Irrelevant information
    // .replace(/<h2 id="toc([a-zA-Z0-9_\-]+)">([a-zA-Z0-9_\-]+)<\/h2>\n/gim, '## $2')
    // .replace(/<h1 id="ory-([a-zA-Z0-9_\-]+)">([a-zA-Z0-9_\-]+)<\/h2>\n/gim, '## $2')
    .replace(/\n\s*\n/g, '\n\n', -1)
    .replace(/^-(\s.*)\n/gim, '-$1', -1)
    .replace(/\n\n---/gi, '\n---\n\n')
    // .replace(/\n\s*\n```/gi, '\n```')
    // .replace(/^<h3 id="[0-9a-zA-Z0-9\-_.]+-responses">Responses<\/h3>$/gim, '#### Summary',-1)
    // .replace(/^> Example responses/gim, '### Responses',-1)
    // .replace(/^> Body parameter/gim, '### Request body',-1)
    .replace(/^> ([0-9]+) Response$/gim, '###### $1 response', -1);

  fs.writeFile(file, t, (err) => {
    if (err) {
      throw err;
    }
  });
});
