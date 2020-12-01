/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const path = require('path')

module.exports = function (context) {
  return {
    name: 'docusaurus-plugin-google-analytics',

    getClientModules() {
      return [path.resolve(__dirname, './analytics')]
    },

    injectHtmlTags() {
      return {
        headTags: [
          {
            tagName: 'script',
            innerHTML: `
window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}

gtag('consent', 'default', {
  'analytics_storage': 'allowed',
  'ad_storage': 'denied',
  'ads_data_redaction': true
});

gtag('consent', 'default', {
  'ad_storage': 'denied',
  'analytics_storage': 'denied',
  'ads_data_redaction': true,
  'region': ['BE','BG','CZ','DK','DE','EE','IE','EL','ES','FR','HR','IT','CY','LV','LT','LU','HU','MT','NL','AT','PL','PT','RO','SI','SK','FI','SE','US-CA']
});
            `
          },
          {
            tagName: 'script',
            attributes: {
              async: true,
              src: 'https://www.googletagmanager.com/gtag/js?id=UA-71865250-1'
            }
          },
          {
            tagName: 'script',
            innerHTML: `
window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}
gtag('js', new Date());

gtag('config', 'G-J01VQCC9Y9'); // automatically anonymized
gtag('config', 'UA-71865250-1', { 'anonymize_ip': true });
            `
          }
        ]
      }
    }
  }
}
