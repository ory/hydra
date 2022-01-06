const path = require('path')

module.exports = function (context) {
  return {
    name: 'docusaurus-plugin-ory-web-script',

    getClientModules() {
      return [path.resolve(__dirname, './ory-scripts-loader')]
    }
  }
}
