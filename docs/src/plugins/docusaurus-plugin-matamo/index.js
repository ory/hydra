/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const path = require('path')

module.exports = function (context) {
  return {
    name: 'docusaurus-plugin-matamo',

    getClientModules() {
      return [path.resolve(__dirname, './analytics')]
    },

    injectHtmlTags() {
      return {
        postBodyTags: [
          `<noscript><p><img src="//sqa-web.ory.sh/np.php?idsite=2&amp;rec=1" style="border:0;" alt="" /></p></noscript>`
        ],
        headTags: [
          {
            tagName: 'script',
            innerHTML: `
var _paq = window._paq = window._paq || [];
/* tracker methods like "setCustomDimension" should be called before "trackPageView" */
_paq.push(["setDocumentTitle", document.domain + "/" + document.title]);
_paq.push(["setCookieDomain", "*.ory.sh"]);
_paq.push(["disableCookies"]);
_paq.push(['trackPageView']);
_paq.push(['enableLinkTracking']);
(function() {
  var u="//sqa-web.ory.sh/";
  _paq.push(['setTrackerUrl', u+'np.php']);
  _paq.push(['setSiteId', '2']);
  var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
  g.type='text/javascript'; g.async=true; g.src=u+'js/np.min.js'; s.parentNode.insertBefore(g,s);
})();
`
          }
        ]
      }
    }
  }
}
