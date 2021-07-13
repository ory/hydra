import { prng } from '../../helpers'

describe('OAuth 2.0 JSON Web Token Access Tokens', () => {
  before(function () {
    // this must be a function otherwise this.skip() fails because the context is wrong
    if (
      Cypress.env('jwt_enabled') !== 'true' &&
      !Boolean(Cypress.env('jwt_enabled'))
    ) {
      this.skip()
    }
  })

  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'offline_access',
    redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
    grant_types: ['authorization_code', 'refresh_token']
  })

  it('should return an Access Token in JWT format and validate it and a Refresh Token in opaque format', () => {
    const client = nc()
    cy.authCodeFlow(client, { consent: { scope: ['offline_access'] } })

    cy.request(`${Cypress.env('client_url')}/oauth2/refresh`)
      .its('body')
      .then((body) => {
        const { result, token } = body
        expect(result).to.equal('success')

        expect(token.access_token).to.not.be.empty
        expect(token.refresh_token).to.not.be.empty
        expect(token.access_token.split('.').length).to.equal(3)
        expect(token.refresh_token.split('.').length).to.equal(2)
      })

    cy.request(`${Cypress.env('client_url')}/oauth2/validate-jwt`)
      .its('body')
      .then((body) => {
        console.log(body)
        expect(body.sub).to.eq('foo@bar.com')
        expect(body.client_id).to.eq(client.client_id)
        expect(body.jti).to.not.be.empty
      })
  })
})
