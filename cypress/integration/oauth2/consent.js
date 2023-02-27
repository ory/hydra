// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { createClient, prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

describe("OAuth 2.0 End-User Authorization", () => {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = () => ({
        client_secret: prng(),
        scope: "offline_access",
        redirect_uris: [`${Cypress.env("client_url")}/oauth2/callback`],
        grant_types: ["authorization_code", "refresh_token"],
        access_token_strategy: accessTokenStrategy,
      })

      const hasConsent = (client, body) => {
        let found = false
        body.forEach(
          ({
            consent_request: {
              client: { client_id },
            },
          }) => {
            if (client_id === client.client_id) {
              found = true
            }
          },
        )
        return found
      }

      it("should check if end user authorization exists", () => {
        createClient(nc()).then((client) => {
          cy.authCodeFlow(client, {
            consent: {
              scope: ["offline_access"],
              remember: true,
            },
            createClient: false,
          })

          console.log("got ", { client })

          cy.request(
            Cypress.env("admin_url") +
              "/oauth2/auth/sessions/consent?subject=foo@bar.com",
          )
            .its("body")
            .then((body) => {
              expect(body.length).to.be.greaterThan(0)
              console.log({
                body,
                client,
              })
              expect(hasConsent(client, body)).to.be.true
              body.forEach((consent) => {
                expect(
                  consent.handled_at.match(
                    /^[2-9]\d{3}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z$/,
                  ),
                ).not.to.be.empty
              })
            })

          cy.request(
            "DELETE",
            Cypress.env("admin_url") +
              "/oauth2/auth/sessions/consent?subject=foo@bar.com&all=true",
          )

          cy.request(
            Cypress.env("admin_url") +
              "/oauth2/auth/sessions/consent?subject=foo@bar.com",
          )
            .its("body")
            .then((body) => {
              expect(body.length).to.eq(0)
              expect(hasConsent(client, body)).to.be.false
            })

          cy.request(`${Cypress.env("client_url")}/oauth2/introspect/at`)
            .its("body")
            .then((body) => {
              expect(body.result).to.equal("success")
              expect(body.body.active).to.be.false
            })

          cy.request(`${Cypress.env("client_url")}/oauth2/introspect/rt`)
            .its("body")
            .then((body) => {
              expect(body.result).to.equal("success")
              expect(body.body.active).to.be.false
            })
        })
      })
    })
  })
})
