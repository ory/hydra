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
    .request('POST', Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers', JSON.stringify(grant))
    .then((response) => {
      const grantID = response.body.id
      getGrant(grantID).then((actual) => {
        if (actual.id !== grantID) {
          return Promise.reject(
            new Error(
              `Expected id's to match: ${actual.id} !== ${grantID}`
            )
          )
        }
        return Promise.resolve(response)
      })
    })

export const getGrant = (grantID) =>
  cy
    .request('GET', Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + grantID)
    .then(({ body }) => body)

export const deleteGrants = () =>
  cy.request(Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers').then(({ body = [] }) => {
    ;(body || []).forEach(({ id }) => deleteGrant(id))
  })

const deleteGrant = (id) =>
  cy.request('DELETE', Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + id)

export const publicJwk = {
    kid: 'token-service-key',
    kty: 'RSA',
    alg: 'RS256',
    n: 'xbOXL8LDbB8hz4fe6__qpESz5GqX0IjH9lRIywG1xj7_w9UXnds5oZpXp0L4TM7B9j0na6_wIwcnfTlQr1cW3LHJXjzPS19zK5rrvB5eabNhtv4yIyH2DSfkI5J3y0bmfY74_J_rDFtQ1PdpfMzdF5cceYvw05B3Q6naPwPN_86GjOkxBWeBZ1-jL5-7cpbbAfeICjEEBsKDX0j-2ZyKpQ2r4jrxwxDF-J3Xsf6ieRKHggQfG-_xMucz40j7t_s-ttE8LoOm9Mmg0gl6vsfhL9rBvUiW-FLCgCqAKSB9a4JHp4_cgsUUR4TsPrJXTGXDFPoqd63S4ZLkCqOeFLOMUx7zVM_gVyyDIbfXWG2HRt6IbEiU8-A-irw0PtPKKiZ0mue2DT3gbvRJlKpL4RG8Obhlaxzf1eQ9jLx15_DoJt9M8zrK9m99YNRMBeJWwJ-RaUv0odpMkIMawH-ly0IO4Kc6fV2g0PK0f4lBnoHze802Y5SQfN19D3GaL93xlHDTHIsX_q0ICyQzupHjQeFHSa9ku0mA36p40lE3Ejpxjbx1BNAvwozGIE7OuovtUgnaodzpRp5HMrCS5YSGE0LtpTgyEibrG3pA12tSvQW3WDeB8qx4dPBBo917ujdgO23p9ZYm96ohZMUOSR_ItX7n3Q4N6W490YrNgj6c-r9kfWk',
    e: 'AQAB'
}
export const privatePem = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

export const invalidPrivatePem = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`