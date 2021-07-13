// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })
import { createClient, prng } from '../helpers'

Cypress.Commands.add(
  'authCodeFlow',
  (
    client,
    {
      override: { scope, client_id, client_secret } = {},
      consent: {
        accept: acceptConsent = true,
        skip: skipConsent = false,
        remember: rememberConsent = false,
        scope: acceptScope = []
      } = {},
      login: {
        accept: acceptLogin = true,
        skip: skipLogin = false,
        remember: rememberLogin = false,
        username = 'foo@bar.com',
        password = 'foobar'
      } = {},
      prompt = '',
      createClient: doCreateClient = true
    } = {},
    path = 'oauth2'
  ) => {
    if (doCreateClient) {
      createClient(client)
    }

    cy.visit(
      `${Cypress.env('client_url')}/${path}/code?client_id=${
        client_id || client.client_id
      }&client_secret=${client_secret || client.client_secret}&scope=${(
        scope || client.scope
      ).replace(' ', '+')}&prompt=${prompt}`,
      { failOnStatusCode: false }
    )

    if (!skipLogin) {
      cy.get('#email').type(username, { delay: 1 })
      cy.get('#password').type(password, { delay: 1 })

      if (rememberLogin) {
        cy.get('#remember').click()
      }

      if (acceptLogin) {
        cy.get('#accept').click()
      } else {
        cy.get('#reject').click()
      }
    }

    if (!skipConsent) {
      acceptScope.forEach((s) => {
        cy.get(`#${s}`).click()
      })

      if (rememberConsent) {
        cy.get('#remember').click()
      }

      if (acceptConsent) {
        cy.get('#accept').click()
      } else {
        cy.get('#reject').click()
      }
    }
  }
)

Cypress.Commands.add(
  'authCodeFlowBrowser',
  (
    client,
    {
      consent: {
        accept: acceptConsent = true,
        skip: skipConsent = false,
        remember: rememberConsent = false,
        scope: acceptScope = []
      } = {},
      login: {
        accept: acceptLogin = true,
        skip: skipLogin = false,
        remember: rememberLogin = false,
        username = 'foo@bar.com',
        password = 'foobar'
      } = {},
      createClient: doCreateClient = true
    } = {}
  ) => {
    if (doCreateClient) {
      createClient(client)
    }

    const codeChallenge = 'QeNVR-BHuB6I2d0HycQzp2qUNNKi_-5QoR4fQSifLH0'
    const codeVerifier =
      'ZmRrenFxZ3pid3A0T0xqY29falJNUS5lWlY4SDBxS182U21uQkhjZ3UuOXpnd3NOak56d2lLMTVYemNNdHdNdlE5TW03WC1RZUlaM0N5R2FhdGRpNW1oVGhjbzVuRFBD'
    const state = prng()

    const authURL = new URL(`${Cypress.env('public_url')}/oauth2/auth`)
    authURL.searchParams.set('response_type', 'code')
    authURL.searchParams.set('client_id', client.client_id)
    authURL.searchParams.set('redirect_uri', client.redirect_uris[0])
    authURL.searchParams.set('scope', client.scope)
    authURL.searchParams.set('state', state)
    authURL.searchParams.set('code_challenge', codeChallenge)
    authURL.searchParams.set('code_challenge_method', 'S256')

    cy.window().then((win) => {
      return win.open(authURL, '_self')
    })

    if (!skipLogin) {
      cy.get('#email').type(username, { delay: 1 })
      cy.get('#password').type(password, { delay: 1 })

      if (rememberLogin) {
        cy.get('#remember').click()
      }

      if (acceptLogin) {
        cy.get('#accept').click()
      } else {
        cy.get('#reject').click()
      }
    }

    if (!skipConsent) {
      acceptScope.forEach((s) => {
        cy.get(`#${s}`).click()
      })

      if (rememberConsent) {
        cy.get('#remember').click()
      }

      if (acceptConsent) {
        cy.get('#accept').click()
      } else {
        cy.get('#reject').click()
      }
    }

    return cy.location('search').then((search) => {
      const callbackParams = new URLSearchParams(search)
      const code = callbackParams.get('code')

      expect(code).to.not.be.empty

      return cy.request({
        url: `${Cypress.env('public_url')}/oauth2/token`,
        method: 'POST',
        form: true,
        body: {
          grant_type: 'authorization_code',
          client_id: client.client_id,
          redirect_uri: client.redirect_uris[0],
          code: code,
          code_verifier: codeVerifier
        }
      })
    })
  }
)

Cypress.Commands.add('refreshTokenBrowser', (client, token) =>
  cy.request({
    url: `${Cypress.env('public_url')}/oauth2/token`,
    method: 'POST',
    form: true,
    body: {
      grant_type: 'refresh_token',
      client_id: client.client_id,
      refresh_token: token
    },
    failOnStatusCode: false
  })
)
