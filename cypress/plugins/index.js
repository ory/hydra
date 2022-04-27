// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)

module.exports = (_on, config) => {
  // Mimic Go's strconv.ParseBool() behavior which Hydra uses to evaluate boolean
  // configuration values.
  const parseBool = (str) => {
    switch (str) {
      case '1':
      case 't':
      case 'T':
      case 'true':
      case 'TRUE':
      case 'True':
        return true
      default:
        return false
    }
  }

  // If admin basicauth is configured, change admin URL used by Cypress tests.
  if (parseBool(process.env.SERVE_ADMIN_BASIC_AUTH_REQUIRED)) {
    config.env.admin_basic_auth = true
    config.env.admin_username = process.env.SERVE_ADMIN_BASIC_AUTH_USERNAME
    config.env.admin_password = process.env.SERVE_ADMIN_BASIC_AUTH_PASSWORD

    let u = new URL(config.env.admin_url)
    u.username = ''
    u.password = ''
    config.env.admin_url_noauth = u.toString()

    u.username = config.env.admin_username
    u.password = config.env.admin_password
    config.env.admin_url = u.toString()
  }

  return config
}
