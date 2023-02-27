// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { createClient, prng } from "../../helpers"
import qs from "querystring"

const accessTokenStrategies = ["opaque", "jwt"]

describe("OpenID Connect Prompt", () => {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = () => ({
        client_secret: prng(),
        scope: "openid",
        redirect_uris: [`${Cypress.env("client_url")}/openid/callback`],
        grant_types: ["authorization_code", "refresh_token"],
        access_token_strategy: accessTokenStrategy,
      })

      it("should fail prompt=none when no session exists", function () {
        createClient(nc()).then((client) => {
          cy.visit(
            `${Cypress.env("client_url")}/openid/code?client_id=${
              client.client_id
            }&client_secret=${client.client_secret}&prompt=none`,
            { failOnStatusCode: false },
          )

          cy.location().should(({ search, port }) => {
            const query = qs.parse(search.substr(1))
            expect(query.error).to.equal("login_required")
            expect(port).to.equal(Cypress.env("client_port"))
          })
        })
      })

      it("should pass with prompt=none if both login and consent were remembered", function () {
        createClient(nc()).then((client) => {
          cy.authCodeFlow(
            client,
            {
              login: { remember: true },
              consent: {
                scope: ["openid"],
                remember: true,
              },
              createClient: false,
            },
            "openid",
          )

          cy.request(
            `${Cypress.env("client_url")}/openid/code?client_id=${
              client.client_id
            }&client_secret=${client.client_secret}&scope=openid`,
          )
            .its("body")
            .then((body) => {
              const {
                result,
                token: { access_token },
              } = body
              expect(result).to.equal("success")
              expect(access_token).to.not.be.empty
            })
        })
      })

      it("should require login with prompt=login even when session exists", function () {
        createClient(nc()).then((client) => {
          cy.authCodeFlow(
            client,
            {
              login: { remember: true },
              consent: {
                scope: ["openid"],
                remember: true,
              },
              createClient: false,
            },
            "openid",
          )

          cy.request(
            `${Cypress.env("client_url")}/openid/code?client_id=${
              client.client_id
            }&client_secret=${client.client_secret}&scope=openid&prompt=login`,
          )
            .its("body")
            .then((body) => {
              expect(body).to.contain("Please log in")
            })
        })
      })

      it("should require consent with prompt=consent even when session exists", function () {
        createClient(nc()).then((client) => {
          cy.authCodeFlow(
            client,
            {
              login: { remember: true },
              consent: {
                scope: ["openid"],
                remember: true,
              },
              createClient: false,
            },
            "openid",
          )

          cy.request(
            `${Cypress.env("client_url")}/openid/code?client_id=${
              client.client_id
            }&client_secret=${
              client.client_secret
            }&scope=openid&prompt=consent`,
          )
            .its("body")
            .then((body) => {
              expect(body).to.contain(
                "An application requests access to your data!",
              )
            })
        })
      })
    })
  })
})
