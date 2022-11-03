// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { prng } from "../../helpers"

describe("The Clients Admin Interface", function () {
  const nc = () => ({
    scope: "foo openid offline_access",
    grant_types: ["client_credentials"],
  })

  it("should return client_secret with length 26 for newly created clients without client_secret specified", function () {
    const client = nc()

    cy.request(
      "POST",
      Cypress.env("admin_url") + "/clients",
      JSON.stringify(client),
    ).then((response) => {
      expect(response.body.client_secret.length).to.equal(26)
    })
  })
})
