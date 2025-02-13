// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { prng } from "../../helpers"

const accessTokenStrategies = ["opaque", "jwt"]

describe("The OAuth 2.0 Device Authorization Grant", function () {
  accessTokenStrategies.forEach((accessTokenStrategy) => {
    describe("access_token_strategy=" + accessTokenStrategy, function () {
      const nc = (extradata) => ({
        client_secret: prng(),
        scope: "offline_access openid",
        subject_type: "public",
        token_endpoint_auth_method: "client_secret_basic",
        grant_types: [
          "urn:ietf:params:oauth:grant-type:device_code",
          "refresh_token",
        ],
        access_token_strategy: accessTokenStrategy,
        ...extradata,
      })

      it("should return an Access, Refresh, and ID Token when scope offline_access and openid are granted", function () {
        const client = nc()
        cy.deviceAuthFlow(client, {
          consent: { scope: ["offline_access", "openid"] },
        })

        cy.postDeviceAuthFlow().then((resp) => {
          const {
            result,
            token: { access_token, id_token, refresh_token },
          } = resp.body

          expect(result).to.equal("success")
          expect(access_token).to.not.be.empty
          expect(id_token).to.not.be.empty
          expect(refresh_token).to.not.be.empty
        })
      })

      it("should return an Access and Refresh Token when scope offline_access is granted", function () {
        const client = nc()
        cy.deviceAuthFlow(client, { consent: { scope: ["offline_access"] } })

        cy.postDeviceAuthFlow().then((resp) => {
          console.log(resp)
          const {
            result,
            token: { access_token, id_token, refresh_token },
          } = resp.body

          expect(result).to.equal("success")
          expect(access_token).to.not.be.empty
          expect(id_token).to.be.undefined
          expect(refresh_token).to.not.be.empty
        })
      })

      it("should return an Access and ID Token when scope offline_access is granted", function () {
        const client = nc()
        cy.deviceAuthFlow(client, { consent: { scope: ["openid"] } })

        cy.postDeviceAuthFlow().then((resp) => {
          console.log(resp)
          const {
            result,
            token: { access_token, id_token, refresh_token },
          } = resp.body

          expect(result).to.equal("success")
          expect(access_token).to.not.be.empty
          expect(id_token).to.not.be.empty
          expect(refresh_token).to.be.undefined
        })
      })

      it("should return an Access Token when no scope is granted", function () {
        const client = nc()
        cy.deviceAuthFlow(client, { consent: { scope: [] } })

        cy.postDeviceAuthFlow().then((resp) => {
          console.log(resp)
          const {
            result,
            token: { access_token, id_token, refresh_token },
          } = resp.body

          expect(result).to.equal("success")
          expect(access_token).to.not.be.empty
          expect(id_token).to.be.undefined
          expect(refresh_token).to.be.undefined
        })
      })

      it("should skip consent if the client is confgured thus", function () {
        const client = nc({ skip_consent: true })
        cy.deviceAuthFlow(client, {
          consent: { scope: ["offline_access", "openid"], skip: true },
        })

        cy.postDeviceAuthFlow().then((resp) => {
          console.log(resp)
          const {
            result,
            token: { access_token, id_token, refresh_token },
          } = resp.body

          expect(result).to.equal("success")
          expect(access_token).to.not.be.empty
          expect(id_token).to.not.be.empty
          expect(refresh_token).to.not.be.empty
        })
      })
    })
  })
})
