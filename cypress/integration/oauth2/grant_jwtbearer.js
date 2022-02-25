import {
  createClient,
  createGrant,
  deleteGrants,
  deleteClients,
  prng
} from '../../helpers'

const dayjs = require('dayjs')
const isBetween = require('dayjs/plugin/isBetween')
const utc = require('dayjs/plugin/utc')
dayjs.extend(utc)
dayjs.extend(isBetween)

const jwt = require('jsonwebtoken')

let testPublicJwk
let testPrivatePem
let invalidtestPrivatePem
const initTestKeyPairs = async () => {
  const algorithm = {
    name: 'RSASSA-PKCS1-v1_5',
    modulusLength: 2048,
    publicExponent: new Uint8Array([1, 0, 1]),
    hash: 'SHA-256'
  }
  const keys = await crypto.subtle.generateKey(algorithm, true, [
    'sign',
    'verify'
  ])

  // public key to jwk
  const publicJwk = await crypto.subtle.exportKey('jwk', keys.publicKey)
  publicJwk.kid = 'token-service-key'

  // private key to pem
  const exportedPK = await crypto.subtle.exportKey('pkcs8', keys.privateKey)
  const exportedAsBase64 = Buffer.from(exportedPK).toString('base64')
  const privatePem = `-----BEGIN PRIVATE KEY-----\n${exportedAsBase64}\n-----END PRIVATE KEY-----`

  // create another private key to test invalid signatures
  const invalidKeys = await crypto.subtle.generateKey(algorithm, true, [
    'sign',
    'verify'
  ])
  const invalidPK = await crypto.subtle.exportKey(
    'pkcs8',
    invalidKeys.privateKey
  )
  const invalidAsBase64 = Buffer.from(invalidPK).toString('base64')
  const invalidPrivatePem = `-----BEGIN PRIVATE KEY-----\n${invalidAsBase64}\n-----END PRIVATE KEY-----`

  testPublicJwk = publicJwk
  testPrivatePem = privatePem
  invalidtestPrivatePem = invalidPrivatePem
}

describe('The OAuth 2.0 JWT Bearer (RFC 7523) Grant', function () {
  beforeEach(() => {
    deleteGrants()
    deleteClients()
  })

  before(() => {
    return cy.wrap(initTestKeyPairs())
  })

  const tokenUrl = `${Cypress.env('public_url')}/oauth2/token`

  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'foo openid offline_access',
    grant_types: ['urn:ietf:params:oauth:grant-type:jwt-bearer'],
    token_endpoint_auth_method: 'client_secret_post',
    response_types: ['token']
  })

  const gr = (subject) => ({
    issuer: prng(),
    subject: subject,
    allow_any_subject: subject === '',
    scope: ['foo', 'openid', 'offline_access'],
    jwk: testPublicJwk,
    expires_at: dayjs().utc().add(1, 'year').set('millisecond', 0).toISOString()
  })

  const jwtAssertion = (grant, override) => {
    const assert = {
      jti: prng(),
      iss: grant.issuer,
      sub: grant.subject,
      aud: tokenUrl,
      exp: dayjs().utc().add(2, 'minute').set('millisecond', 0).unix(),
      iat: dayjs().utc().subtract(2, 'minute').set('millisecond', 0).unix()
    }
    return { ...assert, ...override }
  }

  it('should return an Access Token when given client credentials and a signed JWT assertion', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(jwtAssertion(grant), testPrivatePem, {
      algorithm: 'RS256'
    })

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      }
    })
      .its('body')
      .then((body) => {
        const { access_token, expires_in, scope, token_type } = body

        expect(access_token).to.not.be.empty
        expect(expires_in).to.not.be.undefined
        expect(scope).to.not.be.empty
        expect(token_type).to.not.be.empty
      })
  })

  it('should return an Error (400) when not given client credentials', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(jwtAssertion(grant), testPrivatePem, {
      algorithm: 'RS256'
    })

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion without a jti', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    var ja = jwtAssertion(grant)
    delete ja['jti']
    const assertion = jwt.sign(ja, testPrivatePem, { algorithm: 'RS256' })

    // first token request should work fine
    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with a duplicated jti', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const jwt1 = jwtAssertion(grant)
    const assertion1 = jwt.sign(jwt1, testPrivatePem, { algorithm: 'RS256' })

    // first token request should work fine
    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion1,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      }
    })
      .its('body')
      .then((body) => {
        const { access_token, expires_in, scope, token_type } = body

        expect(access_token).to.not.be.empty
        expect(expires_in).to.not.be.undefined
        expect(scope).to.not.be.empty
        expect(token_type).to.not.be.empty
      })

    const assertion2 = jwt.sign(
      jwtAssertion(grant, { jti: jwt1['jti'] }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    // the second should fail
    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion2,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion without an iat', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    var ja = jwtAssertion(grant)
    delete ja['iat']
    const assertion = jwt.sign(ja, testPrivatePem, {
      algorithm: 'RS256',
      noTimestamp: true
    })

    // first token request should work fine
    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with an invalid signature', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(jwtAssertion(grant), invalidtestPrivatePem, {
      algorithm: 'RS256'
    })

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with an invalid subject', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, { sub: 'invalid_subject' }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Access Token when given client credentials and a JWT assertion with any subject', function () {
    const client = nc()
    createClient(client)

    const grant = gr('') // allow any subject
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, { sub: 'any-subject-is-valid' }),
      testPrivatePem,
      {
        algorithm: 'RS256'
      }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      }
    })
      .its('body')
      .then((body) => {
        const { access_token, expires_in, scope, token_type } = body

        expect(access_token).to.not.be.empty
        expect(expires_in).to.not.be.undefined
        expect(scope).to.not.be.empty
        expect(token_type).to.not.be.empty
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with an invalid issuer', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, { iss: 'invalid_issuer' }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with an invalid audience', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, { aud: 'invalid_audience' }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with an expired date', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, {
        exp: dayjs().utc().subtract(1, 'minute').set('millisecond', 0).unix()
      }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Error (400) when given client credentials and a JWT assertion with a nbf that is still not valid', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, {
        nbf: dayjs().utc().add(1, 'minute').set('millisecond', 0).unix()
      }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      },
      failOnStatusCode: false
    })
      .its('status')
      .then((status) => {
        expect(status).to.be.equal(400)
      })
  })

  it('should return an Access Token when given client credentials and a JWT assertion with a nbf that is valid', function () {
    const client = nc()
    createClient(client)

    const grant = gr(prng())
    createGrant(grant)

    const assertion = jwt.sign(
      jwtAssertion(grant, {
        nbf: dayjs().utc().subtract(1, 'minute').set('millisecond', 0).unix()
      }),
      testPrivatePem,
      { algorithm: 'RS256' }
    )

    cy.request({
      method: 'POST',
      url: tokenUrl,
      form: true,
      body: {
        grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        assertion: assertion,
        scope: client.scope,
        client_secret: client.client_secret,
        client_id: client.client_id
      }
    })
      .its('body')
      .then((body) => {
        const { access_token, expires_in, scope, token_type } = body

        expect(access_token).to.not.be.empty
        expect(expires_in).to.not.be.undefined
        expect(scope).to.not.be.empty
        expect(token_type).to.not.be.empty
      })
  })
})
