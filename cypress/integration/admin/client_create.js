import { prng } from '../../helpers'

describe('The Clients Admin Interface', function () {
  const nc = () => ({
    scope: 'foo openid offline_access',
    grant_types: ['client_credentials']
  })

  if (Cypress.env('admin_basic_auth')) {
    it('should return a 403 error if basic auth credentials are not provided', function () {
      const client = nc()

      cy.request({
        method: 'POST',
        url: Cypress.env('admin_url_noauth') + '/clients',
        failOnStatusCode: false,
        body: JSON.stringify(client)
      }).then((response) => {
        expect(response.status).to.equal(401)
      })
    })
  }

  it('should return client_secret with length 26 for newly created clients without client_secret specified', function () {
    const client = nc()

    cy.request(
      'POST',
      Cypress.env('admin_url') + '/clients',
      JSON.stringify(client)
    ).then((response) => {
      console.log(response.body, client)
      expect(response.body.client_secret.length).to.equal(26)
    })
  })
})
