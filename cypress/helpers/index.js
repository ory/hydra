export const prng = () => {
  var array = new Uint32Array(2)
  crypto.getRandomValues(array)

  return `${array[0].toString()}${array[1].toString()}`
}

const isStatusOk = (res) =>
  res.ok
    ? Promise.resolve(res)
    : Promise.reject(
        new Error(`Received unexpected status code ${res.statusCode}`)
      )

export const findEndUserAuthorization = (subject) =>
  fetch(
    Cypress.env('admin_url') +
      '/oauth2/auth/sessions/consent?subject=' +
      subject
  )
    .then(isStatusOk)
    .then((res) => res.json())

export const revokeEndUserAuthorization = (subject) =>
  fetch(
    Cypress.env('admin_url') +
      '/oauth2/auth/sessions/consent?subject=' +
      subject,
    { method: 'DELETE' }
  ).then(isStatusOk)

export const createClient = (client) =>
  cy
    .request('POST', Cypress.env('admin_url') + '/clients', client)
    .then(({ body }) =>
      getClient(client.client_id).then((actual) => {
        if (actual.client_id !== body.client_id) {
          return Promise.reject(
            new Error(
              `Expected client_id's to match: ${actual.client_id} !== ${body.client}`
            )
          )
        }

        return Promise.resolve(body)
      })
    )

export const deleteClients = () =>
  cy.request(Cypress.env('admin_url') + '/clients').then(({ body = [] }) => {
    ;(body || []).forEach(({ client_id }) => deleteClient(client_id))
  })

const deleteClient = (client_id) =>
  cy.request('DELETE', Cypress.env('admin_url') + '/clients/' + client_id)

const getClient = (id) =>
  cy
    .request(Cypress.env('admin_url') + '/clients/' + id)
    .then(({ body }) => body)

export const createGrant = (grant) =>
  cy
    .request(
      'POST',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      JSON.stringify(grant)
    )
    .then((response) => {
      const grantID = response.body.id
      getGrant(grantID).then((actual) => {
        if (actual.id !== grantID) {
          return Promise.reject(
            new Error(`Expected id's to match: ${actual.id} !== ${grantID}`)
          )
        }
        return Promise.resolve(response)
      })
    })

export const getGrant = (grantID) =>
  cy
    .request(
      'GET',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + grantID
    )
    .then(({ body }) => body)

export const deleteGrants = () =>
  cy
    .request(Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers')
    .then(({ body = [] }) => {
      ;(body || []).forEach(({ id }) => deleteGrant(id))
    })

const deleteGrant = (id) =>
  cy.request(
    'DELETE',
    Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + id
  )
