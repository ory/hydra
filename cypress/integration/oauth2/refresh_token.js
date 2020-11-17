import { prng } from '../../helpers'

describe('OAuth 2.0 JSON Web Token Access Tokens', function () {
  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'offline_access openid',
    redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
    grant_types: ['authorization_code', 'refresh_token']
  })

  it('should return an Access and Refresh Token and refresh the Access Token', function () {
    const client = nc()
    cy.authCodeFlow(client, { consent: { scope: ['offline_access'] } })

    cy.request(`${Cypress.env('client_url')}/oauth2/refresh`)
      .its('body')
      .then((body) => {
        const { result, token } = body
        expect(result).to.equal('success')
        expect(token.access_token).to.not.be.empty
        expect(token.refresh_token).to.not.be.empty
      })
  })

  it('should return an Access, ID, and Refresh Token and refresh the Access Token and ID Token', function () {
    const client = nc()
    cy.authCodeFlow(client, {
      consent: { scope: ['offline_access', 'openid'] }
    })

    cy.request(`${Cypress.env('client_url')}/oauth2/refresh`)
      .its('body')
      .then((body) => {
        const { result, token } = body
        expect(result).to.equal('success')
        expect(token.access_token).to.not.be.empty
        expect(token.id_token).to.not.be.empty
        expect(token.refresh_token).to.not.be.empty
      })
  })
})
