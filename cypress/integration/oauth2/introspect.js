// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

describe("OpenID Connect Token Introspection", () => {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = () => ({
        client_secret: prng(),
        scope: "offline_access",
        redirect_uris: [`${Cypress.env("client_url")}/oauth2/callback`],
        grant_types: ["authorization_code", "refresh_token"],
        access_token_strategy: accessTokenStrategy,
      })

      it("should introspect access token", function () {
        const client = nc()
        cy.authCodeFlow(client, {
          consent: {
            scope: ["offline_access"],
            createClient: true,
          },
        })

        cy.get("body")
          .invoke("text")
          .then((content) => {
            const { result } = JSON.parse(content)
            expect(result).to.equal("success")
          })

        cy.request(`${Cypress.env("client_url")}/oauth2/introspect/at`)
          .its("body")
          .then((body) => {
            expect(body.result).to.equal("success")
            expect(body.body.active).to.be.true
            expect(body.body.sub).to.be.equal("foo@bar.com")
            expect(body.body.token_type).to.be.equal("Bearer")
            expect(body.body.token_use).to.be.equal("access_token")
          })
      })

      it("should introspect refresh token", function () {
        const client = nc()
        cy.authCodeFlow(client, {
          consent: {
            scope: ["offline_access"],
            createClient: true,
          },
        })

        cy.get("body")
          .invoke("text")
          .then((content) => {
            const { result } = JSON.parse(content)
            expect(result).to.equal("success")
          })

        cy.request(`${Cypress.env("client_url")}/oauth2/introspect/rt`)
          .its("body")
          .then((body) => {
            expect(body.result).to.equal("success")
            expect(body.body.active).to.be.true
            expect(body.body.sub).to.be.equal("foo@bar.com")
            expect(body.body.token_type).to.be.equal("Bearer")
            expect(body.body.token_use).to.be.equal("refresh_token")
          })
      })
    })
  })
})
