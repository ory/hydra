/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment'

export default (function () {
  if (!ExecutionEnvironment.canUseDOM) {
    return null
  }

  return {
    onRouteUpdate({ location }) {
      if (typeof window.gtag !== 'function') {
        return
      }

      const pagePath = location
        ? location.pathname + location.search + location.hash
        : undefined
      window.gtag('config', 'UA-71865250-1', {
        page_path: pagePath
      })
    }
  }
})()
