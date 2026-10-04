// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { createClient, deleteClients, prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

// The logout flow depends on the login session cookies surviving between the
// sequential tests in each suite. This includes Hydra's browser session cookie:
// without it, /oauth2/sessions/logout skips the confirmation page and redirects
// straight to the fallback success page.
const preserveLogoutSessionCookies = () => {
  Cypress.Cookies.preserveOnce(
    "oauth2_authentication_session",
    "oauth2_authentication_session_insecure",
    "ory_hydra_session_dev",
    "connect.sid",
  )
}

const clearLogoutSessionCookies = () => {
  cy.clearCookies({ domain: null })
}

const rememberLogin = (client, doCreateClient = true) => {
  cy.authCodeFlow(
    client,
    {
      login: { remember: true },
      consent: {
        scope: ["openid"],
        remember: true,
      },
      createClient: doCreateClient,
    },
    "openid",
  )
}

// Poll the client's session state until it reaches the expected value or the
// attempt budget runs out. A cookie that is written or rotated slightly late
// then recovers on the next poll instead of failing the whole run.
const expectSession = (expected, attempts = 20) => {
  cy.request(`${Cypress.env("client_url")}/openid/session/check`)
    .its("body")
    .then(({ has_session }) => {
      if (has_session === expected) {
        return
      }
      if (attempts <= 1) {
        expect(has_session, "has_session").to.equal(expected)
        return
      }
      cy.wait(250)
      expectSession(expected, attempts - 1)
    })
}

const ensureSession = (client, doCreateClient = true) => {
  cy.request(`${Cypress.env("client_url")}/openid/session/check`)
    .its("body")
    .then(({ has_session }) => {
      if (has_session) {
        return
      }
      clearLogoutSessionCookies()
      rememberLogin(client, doCreateClient)
    })

  expectSession(true)
}

accessTokenStrategies.forEach((accessTokenStrategy) => {
  describe("access_token_strategy=" + accessTokenStrategy, function () {
    const nc = () => ({
      client_secret: prng(),
      scope: "openid",
      subject_type: "public",
      redirect_uris: [`${Cypress.env("client_url")}/openid/callback`],
      grant_types: ["authorization_code"],
      access_token_strategy: accessTokenStrategy,
    })

    describe("OpenID Connect Logout", () => {
      before(() => {
        cy.clearCookies({ domain: null })
      })

      after(() => {
        deleteClients()
      })

      describe("logout without id_token_hint", () => {
        beforeEach(preserveLogoutSessionCookies)

        let client

        before(() => {
          deleteClients()

          createClient({
            ...nc(),
            backchannel_logout_uri: `${Cypress.env(
              "client_url",
            )}/openid/session/end/bc`,
          }).then((createdClient) => {
            client = createdClient
          })
        })

        it("should log in and remember login without id_token_hint", function () {
          clearLogoutSessionCookies()

          rememberLogin(client, false)

          expectSession(true)
        })

        it("should show the logout page and complete logout without id_token_hint", () => {
          ensureSession(client, false)

          cy.visit(`${Cypress.env("client_url")}/openid/session/end?simple=1`, {
            failOnStatusCode: false,
          })

          cy.get("#accept").click()

          cy.get("h1").should(
            "contain",
            "Your log out request however succeeded.",
          )
        })

        it("should show the login screen again because we logged out", () => {
          cy.authCodeFlow(
            client,
            {
              login: { remember: false }, // login should have skip false because we removed the session.mak
              consent: {
                scope: ["openid"],
                remember: false,
                skip: true,
              },
              createClient: false,
            },
            "openid",
          )
        })
      })

      // The Back-Channel test should run before the front-channel test because otherwise both tests need a long time to finish.
      describe("Back-Channel", () => {
        beforeEach(preserveLogoutSessionCookies)

        before(() => {
          deleteClients()
        })

        const client = {
          ...nc(),
          backchannel_logout_uri: `${Cypress.env(
            "client_url",
          )}/openid/session/end/bc`,
        }

        it("should log in and remember login with back-channel", function () {
          clearLogoutSessionCookies()

          rememberLogin(client)

          expectSession(true)
        })

        it("should show the logout page and complete logout with back-channel", () => {
          ensureSession(client)

          cy.visit(`${Cypress.env("client_url")}/openid/session/end`, {
            failOnStatusCode: false,
          })

          cy.get("#accept").click()

          cy.get("h1").should(
            "contain",
            "Your log out request however succeeded.",
          )

          expectSession(false)
        })
      })

      describe("Front-Channel", () => {
        beforeEach(preserveLogoutSessionCookies)

        before(() => {
          deleteClients()
        })

        const client = {
          ...nc(),
          frontchannel_logout_uri: `${Cypress.env(
            "client_url",
          )}/openid/session/end/fc`,
        }

        it("should log in and remember login with front-channel", () => {
          clearLogoutSessionCookies()

          rememberLogin(client)

          expectSession(true)
        })

        it("should show the logout page and complete logout with front-channel", () => {
          ensureSession(client)

          cy.visit(`${Cypress.env("client_url")}/openid/session/end`, {
            failOnStatusCode: false,
          })

          cy.get("#accept").click()

          cy.get("h1").should(
            "contain",
            "Your log out request however succeeded.",
          )

          expectSession(false)
        })
      })
    })
  })
})
