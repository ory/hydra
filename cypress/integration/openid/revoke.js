// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

describe("OpenID Connect Token Revokation", () => {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = () => ({
        client_secret: prng(),
        scope: "openid offline_access",
        redirect_uris: [`${Cypress.env("client_url")}/openid/callback`],
        grant_types: ["authorization_code", "refresh_token"],
        access_token_strategy: accessTokenStrategy,
      })

      it("should be able to revoke the access token", function () {
        const client = nc()
        cy.authCodeFlow(
          client,
          { consent: { scope: ["openid", "offline_access"] } },
          "openid",
        )

        cy.get("body")
          .invoke("text")
          .then((content) => {
            const { result } = JSON.parse(content)
            expect(result).to.equal("success")
          })

        cy.request(`${Cypress.env("client_url")}/openid/revoke/at`)
          .its("body")
          .then((response) => {
            expect(response.result).to.equal("success")
          })

        cy.request(`${Cypress.env("client_url")}/openid/userinfo`, {
          failOnStatusCode: false,
        })
          .its("body")
          .then((response) => {
            expect(response.error).to.contain("request_unauthorized")
          })
      })

      it("should be able to revoke the refresh token", function () {
        const client = nc()
        cy.authCodeFlow(
          client,
          { consent: { scope: ["openid", "offline_access"] } },
          "openid",
        )

        cy.get("body")
          .invoke("text")
          .then((content) => {
            const { result } = JSON.parse(content)
            expect(result).to.equal("success")
          })

        cy.request(`${Cypress.env("client_url")}/openid/revoke/rt`, {
          failOnStatusCode: false,
        })
          .its("body")
          .then((response) => {
            expect(response.result).to.equal("success")
          })

        cy.request(`${Cypress.env("client_url")}/openid/userinfo`, {
          failOnStatusCode: false,
        })
          .its("body")
          .then((response) => {
            expect(response.error).to.contain("request_unauthorized")
          })
      })
    })
  })
})
