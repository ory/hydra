const dayjs = require('dayjs')
const isBetween = require('dayjs/plugin/isBetween')
const utc = require('dayjs/plugin/utc')
dayjs.extend(utc)
dayjs.extend(isBetween)

describe('The JWT-Bearer Grants Admin Interface', () => {
  let d = dayjs().utc().add(1, 'year').set('millisecond', 0)
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
      n: 'ue1_WT_RU6Lc65dmmD7llh9Tcu_Xc909be1Yr5xlHUpkVzacHhSgjliSjUnGCuMo1-m3ILktgt3p86ba6bmIk9fK3nKA7OztDymHuuaYGbJVHhDSKcCBMXGFPcBLxtEns7nvMoQ-lkFN-kYgfSfg0iPGXeRo2Io7phqr54pBaEG_xMK9c-rQ_G3Y9eXn1JREEgQd4OvA2UR9Vc4E-xAYMx7V-ZOvMeKBj9HACE8cllnpKlEKLMo5O5BvkpqA1MeOtzL5jxUUH8D37TJvVQ67VgTs40dRwWwRePfIMDHRJSeJ0KTpkgnX4fmaF2xfi53N8hM9PHzzCtaWrjzm1r1Gyw',
      e: 'AQAB'
    }
  })

  beforeEach(() => {
    // Clean up all previous grants
    cy.request(
      'GET',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers'
    ).then((response) => {
      response.body.map(({ id }) => {
        cy.request(
          'delete',
          Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + id
        ).then(() => {})
      })
    })
  })

  it('should return newly created jwt-bearer grant and grant can be retrieved later', () => {
    const grant = newGrant()
    const start = dayjs().subtract(1, 'minutes')
    const end = dayjs().add(1, 'minutes')
    cy.request(
      'POST',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      JSON.stringify(grant)
    ).then((response) => {
      const createdAt = dayjs(response.body.created_at)
      const expiresAt = dayjs(response.body.expires_at)
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
        Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + grantID
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
    // We have exactly one grant
    const grant = newGrant()
    cy.request(
      'POST',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      JSON.stringify(grant)
    ).then(() => {})
    cy.request(
      'GET',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers'
    ).then((response) => {
      expect(response.body).to.length(1)
    })
  })

  it('should fail, because the same grant is already exist', () => {
    const grant = newGrant()
    cy.request({
      method: 'POST',
      url: Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      failOnStatusCode: false,
      body: JSON.stringify(grant)
    }).then((response) => {
      expect(response.status).to.equal(201)
    })

    cy.request({
      method: 'POST',
      url: Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
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
      url: Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      failOnStatusCode: false,
      body: JSON.stringify(grant)
    }).then((response) => {
      expect(response.status).to.equal(400)
    })
  })

  it('should fail, because trying to create grant with no subject and no allow_any_subject flag', () => {
    const grant = newGrant()
    delete grant.subject
    delete grant.allow_any_subject
    cy.request({
      method: 'POST',
      url: Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      failOnStatusCode: false,
      body: JSON.stringify(grant)
    }).then((response) => {
      expect(response.status).to.equal(400)
    })
  })

  it('should return newly created jwt-bearer grant when issuer is allowed to authorize any subject', () => {
    const grant = newGrant()
    delete grant.subject
    grant.allow_any_subject = true
    const start = dayjs().subtract(1, 'minutes')
    const end = dayjs().add(1, 'minutes')
    cy.request(
      'POST',
      Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers',
      JSON.stringify(grant)
    ).then((response) => {
      const createdAt = dayjs(response.body.created_at)
      const expiresAt = dayjs(response.body.expires_at)
      const grantID = response.body.id

      expect(response.body.allow_any_subject).to.equal(grant.allow_any_subject)
      expect(response.body.issuer).to.equal(grant.issuer)
      expect(createdAt.isBetween(start, end)).to.true
      expect(expiresAt.isSame(grant.expires_at)).to.true
      expect(response.body.scope).to.deep.equal(grant.scope)
      expect(response.body.public_key.set).to.equal(grant.issuer)
      expect(response.body.public_key.kid).to.equal(grant.jwk.kid)

      cy.request(
        'GET',
        Cypress.env('admin_url') + '/trust/grants/jwt-bearer/issuers/' + grantID
      ).then((response) => {
        expect(response.body.allow_any_subject).to.equal(
          grant.allow_any_subject
        )
        expect(response.body.issuer).to.equal(grant.issuer)
        expect(response.body.scope).to.deep.equal(grant.scope)
        expect(response.body.public_key.set).to.equal(grant.issuer)
        expect(response.body.public_key.kid).to.equal(grant.jwk.kid)
      })
    })
  })
})
