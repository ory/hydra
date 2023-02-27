// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { createClient, prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

describe("The OAuth 2.0 Authorization Code Grant", function () {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = () => ({
        client_secret: prng(),
        scope: "foo openid offline_access",
        grant_types: ["client_credentials"],
        access_token_strategy: accessTokenStrategy,
      })

      it("should return an Access Token but not Refresh or ID Token for client_credentials flow", function () {
        createClient(nc()).then((client) => {
          cy.request(
            `${Cypress.env("client_url")}/oauth2/cc?client_id=${
              client.client_id
            }&client_secret=${client.client_secret}&scope=${client.scope}`,
            { failOnStatusCode: false },
          )
            .its("body")
            .then((body) => {
              const {
                result,
                token: { access_token, id_token, refresh_token } = {},
              } = body

              expect(result).to.equal("success")
              expect(access_token).to.not.be.empty
              expect(id_token).to.be.undefined
              expect(refresh_token).to.be.undefined
            })
        })
      })
    })
  })
})
