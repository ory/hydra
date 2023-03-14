// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { deleteClients, prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

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
        beforeEach(() => {
          Cypress.Cookies.preserveOnce(
            "oauth2_authentication_session",
            "oauth2_authentication_session_insecure",
            "connect.sid",
          )
        })

        before(() => {
          deleteClients()
        })

        const client = {
          ...nc(),
          backchannel_logout_uri: `${Cypress.env(
            "client_url",
          )}/openid/session/end/bc`,
        }

        it("should log in and remember login without id_token_hint", function () {
          cy.authCodeFlow(
            client,
            {
              login: { remember: true },
              consent: {
                scope: ["openid"],
                remember: true,
              },
            },
            "openid",
          )

          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.true
            })
        })

        it("should show the logout page and complete logout without id_token_hint", () => {
          // cy.request(`${Cypress.env('client_url')}/openid/session/check`)
          //   .its('body')
          //   .then(({ has_session }) => {
          //     expect(has_session).to.be.true;
          //   });

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
      describe.only("Back-Channel", () => {
        beforeEach(() => {
          Cypress.Cookies.preserveOnce(
            "oauth2_authentication_session",
            "oauth2_authentication_session_insecure",
            "connect.sid",
          )
        })

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
          cy.authCodeFlow(
            client,
            {
              login: { remember: true },
              consent: {
                scope: ["openid"],
                remember: true,
              },
            },
            "openid",
          )

          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.true
            })
        })

        it("should show the logout page and complete logout with back-channel", () => {
          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.true
            })

          cy.visit(`${Cypress.env("client_url")}/openid/session/end`, {
            failOnStatusCode: false,
          })

          cy.get("#accept").click()

          cy.get("h1").should(
            "contain",
            "Your log out request however succeeded.",
          )

          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.false
            })
        })
      })

      describe("Front-Channel", () => {
        beforeEach(() => {
          Cypress.Cookies.preserveOnce(
            "oauth2_authentication_session",
            "oauth2_authentication_session_insecure",
            "connect.sid",
          )
        })

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
          cy.authCodeFlow(
            client,
            {
              login: { remember: true },
              consent: {
                scope: ["openid"],
                remember: true,
              },
            },
            "openid",
          )

          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.true
            })
        })

        it("should show the logout page and complete logout with front-channel", () => {
          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.true
            })

          cy.visit(`${Cypress.env("client_url")}/openid/session/end`, {
            failOnStatusCode: false,
          })

          cy.get("#accept").click()

          cy.get("h1").should(
            "contain",
            "Your log out request however succeeded.",
          )

          cy.request(`${Cypress.env("client_url")}/openid/session/check`)
            .its("body")
            .then(({ has_session }) => {
              expect(has_session).to.be.false
            })
        })
      })
    })
  })
})
