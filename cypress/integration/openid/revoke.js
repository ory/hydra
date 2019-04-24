import { prng } from '../../helpers'

describe('OpenID Connect Token Revokation', () => {
  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'openid offline_access',
    redirect_uris: ['http://127.0.0.1:4000/openid/callback'],
    grant_types: ['authorization_code', 'refresh_token']
  })

  it('should be able to revoke the access token', function () {
    const client = nc()
    cy.authCodeFlow(client, { consent: { scope: ['openid', 'offline_access'] } }, 'openid')

    cy.get('body')
      .invoke('text')
      .then((content) => {
        const { result } = JSON.parse(content)
        expect(result).to.equal('success')
      })

    cy.request('http://127.0.0.1:4000/openid/revoke/at').its('body').then((response) => {
      expect(response.result).to.equal('success')
    })

    cy.request('http://127.0.0.1:4000/openid/userinfo', {failOnStatusCode: false}).its('body').then((response) => {
      expect(response.error).to.contain('request_unauthorized')
    })
  })

  it('should be able to revoke the refresh token', function () {
    const client = nc()
    cy.authCodeFlow(client, { consent: { scope: ['openid', 'offline_access'] } }, 'openid')

    cy.get('body')
      .invoke('text')
      .then((content) => {
        const { result } = JSON.parse(content)
        expect(result).to.equal('success')
      })

    cy.request('http://127.0.0.1:4000/openid/revoke/rt', {failOnStatusCode: false}).its('body').then((response) => {
      expect(response.result).to.equal('success')
    })

    cy.request('http://127.0.0.1:4000/openid/userinfo', {failOnStatusCode: false}).its('body').then((response) => {
      expect(response.error).to.contain('request_unauthorized')
    })
  })
})
