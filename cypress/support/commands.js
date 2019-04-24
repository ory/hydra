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

Cypress.Commands.add('authCodeFlow', (client, {
  override: {
    scope,
    client_id,
    client_secret
  } = {},
  consent: {
    accept: acceptConsent = true,
    skip: skipConsent = false,
    remember: rememberConsent = false,
    scope: acceptScope = []
  } = {},
  login: {
    skip: skipLogin = false,
    remember: rememberLogin = false,
    username = 'foo@bar.com',
    password = 'foobar'
  } = {},
  prompt = '',
} = {}, path = 'oauth2') => {
  cy.wrap(createClient(client))

  cy.visit(
    `http://127.0.0.1:4000/${path}/code?client_id=${client_id || client.client_id}&client_secret=${client_secret || client.client_secret}&scope=${(scope || client.scope).replace(' ', '+')}&prompt=${prompt}`,
    { failOnStatusCode: false },
  )

  if (!skipLogin) {
    cy.get('input[name="email"]').type(username)
    cy.get('input[name="password"]').type(password)
    if (rememberLogin) {
      cy.get('#remember').click()
    }
    cy.get('input[type="submit"]').click()
  }

  if (!skipConsent) {
    acceptScope.forEach((s) => {
      cy.get(`#${s}`).click()
    })

    if (rememberConsent) {
      cy.get('#remember').click()
    }

    if (acceptConsent) {
      cy.get('input[value="Allow access"]').click()
    } else {
      cy.get('input[value="Deny access"]').click()
    }
  }
})
