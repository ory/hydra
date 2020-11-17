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
import { createClient } from '../helpers'

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
      cy.wrap(createClient(client))
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
