// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { createClient, prng, rotateJwks, validateJwt } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

describe("The OAuth 2.0 Refresh Token Grant", function () {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = () => ({
        client_secret: prng(),
        scope: "offline_access openid",
        redirect_uris: [`${Cypress.env("client_url")}/oauth2/callback`],
        grant_types: ["authorization_code", "refresh_token"],
        access_token_strategy: accessTokenStrategy,
      })

      it("should return an Access and Refresh Token and refresh the Access Token", function () {
        const client = nc()
        cy.authCodeFlow(client, {
          consent: {
            scope: ["offline_access"],
            createClient: true,
          },
        })

        cy.request(`${Cypress.env("client_url")}/oauth2/refresh`)
          .its("body")
          .then((body) => {
            const { result, token } = body
            expect(result).to.equal("success")
            expect(token.access_token).to.not.be.empty
            expect(token.refresh_token).to.not.be.empty
          })
      })

      it("should return an Access, ID, and Refresh Token and refresh the Access Token and ID Token", function () {
        const client = nc()
        cy.authCodeFlow(client, {
          consent: {
            scope: ["offline_access", "openid"],
            createClient: true,
          },
        })

        cy.request(`${Cypress.env("client_url")}/oauth2/refresh`)
          .its("body")
          .then((body) => {
            const { result, token } = body
            expect(result).to.equal("success")
            expect(token.access_token).to.not.be.empty
            expect(token.id_token).to.not.be.empty
            expect(token.refresh_token).to.not.be.empty
          })
      })

      it("should revoke Refresh Token on reuse", function () {
        const referrer = `${Cypress.env("client_url")}/empty`
        cy.visit(referrer, {
          failOnStatusCode: false,
        })

        createClient({
          scope: "offline_access",
          redirect_uris: [referrer],
          grant_types: ["authorization_code", "refresh_token"],
          response_types: ["code"],
          token_endpoint_auth_method: "none",
        }).then((client) => {
          cy.authCodeFlowBrowser(client, {
            consent: { scope: ["offline_access"] },
            createClient: false,
          }).then((originalResponse) => {
            expect(originalResponse.status).to.eq(200)
            expect(originalResponse.body.refresh_token).to.not.be.empty

            const originalToken = originalResponse.body.refresh_token

            cy.refreshTokenBrowser(client, originalToken).then(
              (refreshedResponse) => {
                expect(refreshedResponse.status).to.eq(200)
                expect(refreshedResponse.body.refresh_token).to.not.be.empty

                const refreshedToken = refreshedResponse.body.refresh_token

                return cy
                  .refreshTokenBrowser(client, originalToken)
                  .then((response) => {
                    expect(response.status).to.eq(400)
                    expect(response.body.error).to.eq("invalid_grant")
                  })
                  .then(() => cy.refreshTokenBrowser(client, refreshedToken))
                  .then((response) => {
                    expect(response.status).to.eq(400)
                    expect(response.body.error).to.eq("invalid_grant")
                  })
              },
            )
          })
        })
      })

      const validateJwtAndGetKid = (token) =>
        validateJwt(token).then(({ header }) => header.kid)

      it("should refresh the Access and ID Token with newly rotated keys", function () {
        if (
          accessTokenStrategy === "opaque" ||
          (Cypress.env("jwt_enabled") !== "true" &&
            !Boolean(Cypress.env("jwt_enabled")))
        ) {
          this.skip()
        }

        const referrer = `${Cypress.env("client_url")}/empty`
        cy.visit(referrer, {
          failOnStatusCode: false,
        })

        createClient({
          scope: "offline_access openid",
          redirect_uris: [referrer],
          grant_types: ["authorization_code", "refresh_token"],
          response_types: ["code"],
          token_endpoint_auth_method: "none",
        }).then((client) => {
          cy.authCodeFlowBrowser(client, {
            consent: {
              scope: ["offline_access", "openid"],
            },
            createClient: false,
          }).then(({ body: tokensBefore }) => {
            const kidsBefore = {
              accessToken: validateJwtAndGetKid(tokensBefore.access_token),
              idToken: validateJwtAndGetKid(tokensBefore.id_token),
            }

            rotateJwks("hydra.jwt.access-token")
            rotateJwks("hydra.openid.id-token")

            cy.refreshTokenBrowser(client, tokensBefore.refresh_token).then(
              ({ body: tokensAfter }) => {
                const kidsAfter = {
                  accessToken: validateJwtAndGetKid(tokensAfter.access_token),
                  idToken: validateJwtAndGetKid(tokensAfter.id_token),
                }

                expect(kidsAfter.accessToken).to.not.equal(
                  kidsBefore.accessToken,
                )
                expect(kidsAfter.idToken).to.not.equal(kidsBefore.idToken)
              },
            )
          })
        })
      })
    })
  })
})
