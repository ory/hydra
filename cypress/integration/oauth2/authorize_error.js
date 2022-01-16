import { createClient, prng } from '../../helpers'
import qs from 'querystring'

describe('OAuth 2.0 Authorization Endpoint Error Handling', () => {
  describe('rejecting login and consent requests', () => {
    const nc = () => ({
      client_id: prng(),
      client_secret: prng(),
      scope: 'offline_access openid',
      subject_type: 'public',
      token_endpoint_auth_method: 'client_secret_basic',
      redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
      grant_types: ['authorization_code', 'refresh_token']
    })

    it('should return an error when rejecting login', function () {
      const client = nc()
      cy.authCodeFlow(client, {
        login: { accept: false },
        consent: { skip: true }
      })

      cy.get('body')
        .invoke('text')
        .then((content) => {
          const {
            result,
            error_description,
            token: { access_token, id_token, refresh_token } = {}
          } = JSON.parse(content)

          expect(result).to.equal('error')
          expect(error_description).to.equal(
            'The resource owner denied the request'
          )
          expect(access_token).to.be.undefined
          expect(id_token).to.be.undefined
          expect(refresh_token).to.be.undefined
        })
    })

    it('should return an error when rejecting consent', function () {
      const client = nc()
      cy.authCodeFlow(client, {
        consent: { accept: false }
      })

      cy.get('body')
        .invoke('text')
        .then((content) => {
          const {
            result,
            error_description,
            token: { access_token, id_token, refresh_token } = {}
          } = JSON.parse(content)

          expect(result).to.equal('error')
          expect(error_description).to.equal(
            'The resource owner denied the request'
          )
          expect(access_token).to.be.undefined
          expect(id_token).to.be.undefined
          expect(refresh_token).to.be.undefined
        })
    })
  })

  it('should return an error when an OAuth 2.0 Client ID is used that does not exist', () => {
    cy.visit(
      `${Cypress.env(
        'client_url'
      )}/oauth2/code?client_id=i-do-not-exist&client_secret=i-am-not-correct}`,
      { failOnStatusCode: false }
    )

    cy.location().should(({ search, port }) => {
      const query = qs.parse(search.substr(1))
      expect(query.error).to.equal('invalid_client')

      // Should show Ory Hydra's Error URL because a redirect URL could not be determined
      expect(port).to.equal(Cypress.env('public_port'))
    })
  })

  it('should return an error when an OAuth 2.0 Client requests a scope that is not allowed to be requested', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      scope: 'foo',
      redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
      grant_types: ['authorization_code']
    }
    createClient(c)

    cy.visit(
      `${Cypress.env('client_url')}/oauth2/code?client_id=${
        c.client_id
      }&client_secret=${c.client_secret}&scope=bar`,
      { failOnStatusCode: false }
    )

    cy.location().should(({ search, port }) => {
      const query = qs.parse(search.substr(1))
      expect(query.error).to.equal('invalid_scope')

      // This is a client error so we expect the client app to show the error
      expect(port).to.equal(Cypress.env('client_port'))
    })
  })

  it('should return an error when an OAuth 2.0 Client requests a response type it is not allowed to call', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
      response_types: ['token'] // disallows Authorization Code Grant
    }
    createClient(c)

    cy.visit(
      `${Cypress.env('client_url')}/oauth2/code?client_id=${
        c.client_id
      }&client_secret=${c.client_secret}`,
      { failOnStatusCode: false }
    )

    cy.get('body').should('contain', 'unsupported_response_type')
  })

  it('should return an error when an OAuth 2.0 Client requests a grant type it is not allowed to call', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
      grant_types: ['client_credentials']
    }
    createClient(c)

    cy.visit(
      `${Cypress.env('client_url')}/oauth2/code?client_id=${
        c.client_id
      }&client_secret=${c.client_secret}&scope=`,
      { failOnStatusCode: false }
    )

    cy.get('#email').type('foo@bar.com', { delay: 1 })
    cy.get('#password').type('foobar', { delay: 1 })
    cy.get('#accept').click()
    cy.get('#accept').click()

    cy.get('body').should('contain', 'unauthorized_client')
  })

  it('should return an error when an OAuth 2.0 Client requests a redirect_uri that is not preregistered', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      redirect_uris: ['http://some-other-domain/not-callback'],
      grant_types: ['client_credentials']
    }
    createClient(c)

    cy.visit(
      `${Cypress.env('client_url')}/oauth2/code?client_id=${
        c.client_id
      }&client_secret=${c.client_secret}&scope=`,
      { failOnStatusCode: false }
    )

    cy.location().should(({ search, port }) => {
      const query = qs.parse(search.substr(1))
      console.log(query)
      expect(query.error).to.equal('invalid_request')
      expect(query.error_description).to.contain('redirect_uri')

      // Should show Ory Hydra's Error URL because a redirect URL could not be determined
      expect(port).to.equal(Cypress.env('public_port'))
    })
  })
})
