import { prng } from '../../helpers'

describe('OpenID Connect Userinfo', () => {
  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'openid',
    subject_type: 'public',
    token_endpoint_auth_method: 'client_secret_basic',
    redirect_uris: ['http://127.0.0.1:4000/callback'],
    grant_types: ['authorization_code', 'refresh_token']
  })

  it('should successfully call the userinfo endpoint after authorization and return a sid', function () {
    const client = nc()
    cy.oAuth2AuthCodeFlow(client, { consent: { scope: ['openid'] } })

    cy.request('http://127.0.0.1:4000/userinfo').its('body').then((body) => {
      expect(body.sub).to.eq('foo@bar.com')
      expect(body.sid).to.not.be.empty
    })
  })
})
