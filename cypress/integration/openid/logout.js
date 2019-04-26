import { deleteClients, prng } from '../../helpers'

const nc = () => ({
  client_id: prng(),
  client_secret: prng(),
  scope: 'openid',
  subject_type: 'public',
  redirect_uris: [`${Cypress.env('client_url')}/openid/callback`],
  grant_types: ['authorization_code']
})

describe('OpenID Connect Logout', () => {
  after(() => {
    cy.wrap(deleteClients())
  })

  // The Back-Channel test should run before the front-channel test because otherwise both tests need a long time to finish.
  describe('Back-Channel', () => {
    beforeEach(() => {
      Cypress.Cookies.preserveOnce(
        'oauth2_authentication_session',
        'connect.sid'
      )
    })

    before(() => {
      cy.wrap(deleteClients())
    })

    const client = {
      ...nc(),
      backchannel_logout_uri: `${Cypress.env(
        'client_url'
      )}/openid/session/end/bc`
    }

    it('should log in and remember login with back-channel', function () {
      cy.authCodeFlow(
        client,
        {
          login: { remember: true },
          consent: { scope: ['openid'], remember: true }
        },
        'openid'
      )

      cy.request(`${Cypress.env('client_url')}/openid/session/check`)
        .its('body')
        .then(({ has_session }) => {
          expect(has_session).to.be.true
        })
    })

    it('should show the logout page and complete logout with back-channel', () => {
      cy.request(`${Cypress.env('client_url')}/openid/session/check`)
        .its('body')
        .then(({ has_session }) => {
          expect(has_session).to.be.true
        })

      cy.visit(`${Cypress.env('client_url')}/openid/session/end`, {
        failOnStatusCode: false
      })

      cy.get('#accept').click()

      cy.get('h1').should('contain', 'Your log out request however succeeded.')

      cy.request(`${Cypress.env('client_url')}/openid/session/check`)
        .its('body')
        .then(({ has_session }) => {
          expect(has_session).to.be.false
        })
    })
  })

  describe('Front-Channel', () => {
    beforeEach(() => {
      Cypress.Cookies.preserveOnce(
        'oauth2_authentication_session',
        'connect.sid'
      )
    })

    before(() => {
      cy.wrap(deleteClients())
    })

    const client = {
      ...nc(),
      frontchannel_logout_uri: `${Cypress.env(
        'client_url'
      )}/openid/session/end/fc`
    }

    it('should log in and remember login with front-channel', () => {
      cy.authCodeFlow(
        client,
        {
          login: { remember: true },
          consent: { scope: ['openid'], remember: true }
        },
        'openid'
      )

      cy.request(`${Cypress.env('client_url')}/openid/session/check`)
        .its('body')
        .then(({ has_session }) => {
          expect(has_session).to.be.true
        })
    })

    it('should show the logout page and complete logout with front-channel', () => {
      cy.request(`${Cypress.env('client_url')}/openid/session/check`)
        .its('body')
        .then(({ has_session }) => {
          expect(has_session).to.be.true
        })

      cy.visit(`${Cypress.env('client_url')}/openid/session/end`, {
        failOnStatusCode: false
      })

      cy.get('#accept').click()

      cy.get('h1').should('contain', 'Your log out request however succeeded.')

      cy.request(`${Cypress.env('client_url')}/openid/session/check`)
        .its('body')
        .then(({ has_session }) => {
          expect(has_session).to.be.false
        })
    })
  })
})
