import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment'

export default (function () {
  if (
    !ExecutionEnvironment.canUseDOM ||
    process.env.NODE_ENV !== 'production'
  ) {
    return null
  }

  const script = document.createElement('script')
  script.src = 'https://www.ory.sh/scripts.js'
  script.onload = () => window.initAnalytics()
  document.body.appendChild(script)

  return {
    onRouteUpdate() {
      if (window && typeof window.initAnalytics === 'function') {
        window.initAnalytics()
      }
    }
  }
})()
