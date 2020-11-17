export const prng = () =>
  `${Math.random().toString(36).substring(2)}${Math.random()
    .toString(36)
    .substring(2)}`

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
  fetch(Cypress.env('admin_url') + '/clients', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(client)
  })
    .then(isStatusOk)
    .then((res) => {
      return res.json()
    })
    .then((body) =>
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
  fetch(Cypress.env('admin_url') + '/clients', {
    method: 'GET'
  })
    .then(isStatusOk)
    .then((res) => res.json())
    .then((body = []) => {
      ;(body || []).forEach(({ client_id }) => deleteClient(client_id))
    })

const deleteClient = (client_id) =>
  fetch(Cypress.env('admin_url') + '/clients/' + client_id, {
    method: 'DELETE'
  }).then(isStatusOk)

const getClient = (id) =>
  fetch(Cypress.env('admin_url') + '/clients/' + id)
    .then(isStatusOk)
    .then((res) => res.json())
