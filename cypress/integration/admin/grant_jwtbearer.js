describe('The JWT-Bearer Grants Admin Interface', () => {
  let d = Cypress.moment().add(1, 'year').milliseconds(0).utc()
  const newGrant = () => ({
    issuer: 'token-service',
    subject: 'bob@example.com',
    expires_at: d.toISOString(),
    scope: ['openid', 'offline'],
    jwk: {
      use: 'sig',
      kty: 'RSA',
      kid: 'token-service-key',
      alg: 'RS256',
      n:
        'ue1_WT_RU6Lc65dmmD7llh9Tcu_Xc909be1Yr5xlHUpkVzacHhSgjliSjUnGCuMo1-m3ILktgt3p86ba6bmIk9fK3nKA7OztDymHuuaYGbJVHhDSKcCBMXGFPcBLxtEns7nvMoQ-lkFN-kYgfSfg0iPGXeRo2Io7phqr54pBaEG_xMK9c-rQ_G3Y9eXn1JREEgQd4OvA2UR9Vc4E-xAYMx7V-ZOvMeKBj9HACE8cllnpKlEKLMo5O5BvkpqA1MeOtzL5jxUUH8D37TJvVQ67VgTs40dRwWwRePfIMDHRJSeJ0KTpkgnX4fmaF2xfi53N8hM9PHzzCtaWrjzm1r1Gyw',
      e: 'AQAB'
    }
  })

  it('should return newly created jwt-bearer grant and grant can be retrieved later', () => {
    const grant = newGrant()
    const start = Cypress.moment().subtract(1, 'minutes').utc()
    const end = Cypress.moment().add(1, 'minutes').utc()
    cy.request(
      'POST',
      Cypress.env('admin_url') + '/grants/jwt-bearer',
      JSON.stringify(grant)
    ).then((response) => {
      const createdAt = Cypress.moment(response.body.created_at)
      const expiresAt = Cypress.moment(response.body.expires_at)
      const grantID = response.body.id

      expect(response.body.issuer).to.equal(grant.issuer)
      expect(response.body.subject).to.equal(grant.subject)
      expect(createdAt.isBetween(start, end)).to.true
      expect(expiresAt.isSame(grant.expires_at)).to.true
      expect(response.body.scope).to.deep.equal(grant.scope)
      expect(response.body.public_key.set).to.equal(grant.issuer)
      expect(response.body.public_key.kid).to.equal(grant.jwk.kid)

      cy.request(
        'GET',
        Cypress.env('admin_url') + '/grants/jwt-bearer/' + grantID
      ).then((response) => {
        expect(response.body.issuer).to.equal(grant.issuer)
        expect(response.body.subject).to.equal(grant.subject)
        expect(response.body.scope).to.deep.equal(grant.scope)
        expect(response.body.public_key.set).to.equal(grant.issuer)
        expect(response.body.public_key.kid).to.equal(grant.jwk.kid)
      })
    })
  })

  it('should return newly created jwt-bearer grant in grants list', () => {
    cy.request('GET', Cypress.env('admin_url') + '/grants/jwt-bearer').then(
      (response) => {
        expect(response.body).to.length(1)
      }
    )
  })

  it('should fail, because the same grant is already exist', () => {
    const grant = newGrant()
    cy.request({
      method: 'POST',
      url: Cypress.env('admin_url') + '/grants/jwt-bearer',
      failOnStatusCode: false,
      body: JSON.stringify(grant)
    }).then((response) => {
      expect(response.status).to.equal(409)
    })
  })

  it('should fail, because trying to create grant with no issuer', () => {
    const grant = newGrant()
    grant.issuer = ''
    cy.request({
      method: 'POST',
      url: Cypress.env('admin_url') + '/grants/jwt-bearer',
      failOnStatusCode: false,
      body: JSON.stringify(grant)
    }).then((response) => {
      expect(response.status).to.equal(400)
    })
  })
})
