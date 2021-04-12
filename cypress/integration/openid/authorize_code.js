import { prng } from '../../helpers'

describe('OpenID Connect Authorize Code Grant', () => {
  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'openid',
    subject_type: 'public',
    token_endpoint_auth_method: 'client_secret_basic',
    redirect_uris: [`${Cypress.env('client_url')}/openid/callback`],
    grant_types: ['authorization_code', 'refresh_token']
  })

  it('should return an access, refresh, and ID token', function () {
    const client = nc()
    cy.authCodeFlow(client, { consent: { scope: ['openid'] } }, 'openid')

    cy.get('body')
      .invoke('text')
      .then((content) => {
        const {
          result,
          token: { access_token, id_token, refresh_token },
          claims: { sub, sid }
        } = JSON.parse(content)

        expect(result).to.equal('success')
        expect(access_token).to.not.be.empty
        expect(id_token).to.not.be.empty
        expect(refresh_token).to.be.undefined

        expect(sub).to.eq('foo@bar.com')
        expect(sid).to.not.be.empty
      })
  })
})
