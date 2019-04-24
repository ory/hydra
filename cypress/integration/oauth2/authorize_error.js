import { createClient, prng } from '../../helpers'
import qs from 'querystring'

describe('OAuth 2.0 Authorization Endpoint Error Handling', () => {
  it('should return an error when an OAuth 2.0 Client ID is used that does not exist', () => {
    cy.visit(
      `http://127.0.0.1:4000/oauth2/code?client_id=i-do-not-exist&client_secret=i-am-not-correct}`,
      { failOnStatusCode: false },
    )

    cy.location().should(({ search, port }) => {
      const query = qs.parse(search.substr(1))
      expect(query.error).to.equal('invalid_client')

      // Should show ORY Hydra's Error URL because a redirect URL could not be determined
      expect(port).to.equal('4444')
    })
  })

  it('should return an error when an OAuth 2.0 Client requests a scope that is not allowed to be requested', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      scope: 'foo',
      redirect_uris: ['http://127.0.0.1:4000/oauth2/callback'],
      grant_types: ['authorization_code']
    }
    cy.wrap(createClient(c))

    cy.visit(
      `http://127.0.0.1:4000/oauth2/code?client_id=${c.client_id}&client_secret=${c.client_secret}&scope=bar`,
      { failOnStatusCode: false },
    )

    cy.location().should(({ search, port }) => {
      const query = qs.parse(search.substr(1))
      expect(query.error).to.equal('invalid_scope')

      // This is a client error so we expect the client app to show the error
      expect(port).to.equal('4000')
    })
  })

  it('should return an error when an OAuth 2.0 Client requests a response type it is not allowed to call', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      redirect_uris: ['http://127.0.0.1:4000/oauth2/callback'],
      response_types: ['token'] // disallows Authorization Code Grant
    }
    cy.wrap(createClient(c))

    cy.visit(
      `http://127.0.0.1:4000/oauth2/code?client_id=${c.client_id}&client_secret=${c.client_secret}`,
      { failOnStatusCode: false },
    )

    cy.get('body').should('contain', 'unsupported_response_type')
  })

  it('should return an error when an OAuth 2.0 Client requests a grant type it is not allowed to call', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      redirect_uris: ['http://127.0.0.1:4000/oauth2/callback'],
      grant_types: ['client_credentials']
    }
    cy.wrap(createClient(c))

    cy.visit(
      `http://127.0.0.1:4000/oauth2/code?client_id=${c.client_id}&client_secret=${c.client_secret}`,
      { failOnStatusCode: false },
    )

    cy.get('input[name="email"]').type('foo@bar.com')
    cy.get('input[name="password"]').type('foobar')
    cy.get('input[type="submit"]').click()
    cy.get('input[value="Allow access"]').click()

    cy.get('body').should('contain', 'invalid_grant')
  })

  it('should return an error when an OAuth 2.0 Client requests a redirect_uri that is not preregistered', () => {
    const c = {
      client_id: prng(),
      client_secret: prng(),
      redirect_uris: ['http://some-other-domain/not-callback'],
      grant_types: ['client_credentials']
    }
    cy.wrap(createClient(c))

    cy.visit(
      `http://127.0.0.1:4000/oauth2/code?client_id=${c.client_id}&client_secret=${c.client_secret}`,
      { failOnStatusCode: false },
    )

    cy.location().should(({ search, port }) => {
      const query = qs.parse(search.substr(1))
      expect(query.error).to.equal('invalid_request')
      expect(query.error_hint).to.contain('redirect_uri')

      // Should show ORY Hydra's Error URL because a redirect URL could not be determined
      expect(port).to.equal('4444')
    })
  })
})