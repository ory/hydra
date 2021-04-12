import { createClient, prng } from '../../helpers'

describe('The OAuth 2.0 Authorization Code Grant', function () {
  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'foo openid offline_access',
    grant_types: ['client_credentials']
  })

  it('should return an Access Token but not Refresh or ID Token for client_credentials flow', function () {
    const client = nc()
    createClient(client)

    cy.request(
      `${Cypress.env('client_url')}/oauth2/cc?client_id=${
        client.client_id
      }&client_secret=${client.client_secret}&scope=${client.scope}`,
      { failOnStatusCode: false }
    )
      .its('body')
      .then((body) => {
        const {
          result,
          token: { access_token, id_token, refresh_token } = {}
        } = body

        expect(result).to.equal('success')
        expect(access_token).to.not.be.empty
        expect(id_token).to.be.undefined
        expect(refresh_token).to.be.undefined
      })
  })
})
