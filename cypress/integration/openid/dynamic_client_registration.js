describe('OAuth2 / OpenID Connect Dynamic Client Registration', function () {
  it('should return same client_secret given in request for newly created clients with client_secret specified', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        client_name: 'clientName',
        scope: 'foo openid offline_access',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      expect(response.body.registration_access_token).to.not.be.empty
    })
  })

  it('should get client when having a valid client_secret in body', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        client_name: 'clientName',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'GET',
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id,
        headers: {
          Authorization: 'Bearer ' + response.body.registration_access_token
        }
      }).then((response) => {
        expect(response.body.client_name).to.equal('clientName')
      })
    })
  })

  it('should fail for get client when Authorization header is not presented', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'GET',
        failOnStatusCode: false,
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id
      }).then((response) => {
        expect(response.status).to.eq(401)
      })
    })
  })

  it('should update client name', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'PUT',
        failOnStatusCode: false,
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id,
        headers: {
          Authorization: 'Bearer ' + response.body.registration_access_token
        },
        body: {
          client_id: 'clientid',
          client_name: 'clientName2',
          scope: 'foo openid offline_access',
          grant_types: ['client_credentials']
        }
      }).then((response) => {
        expect(response.body.client_name).to.equal('clientName2')
        expect(response.body.client_id).to.not.equal('clientid')
      })
    })
  })
  it('should not be able to choose the secret', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'PUT',
        failOnStatusCode: false,
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id,
        headers: {
          Authorization: 'Bearer ' + response.body.registration_access_token
        },
        body: {
          client_id: 'clientid',
          client_name: 'clientName2',
          client_secret: 'secret',
          scope: 'foo openid offline_access',
          grant_types: ['client_credentials']
        }
      }).then((response) => {
        expect(response.status).to.eq(403)
      })
    })
  })

  it('should fail to update client name when Authorization header is not presented', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        client_name: 'clientName',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'PUT',
        failOnStatusCode: false,
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id
      }).then((response) => {
        expect(response.status).to.eq(401)
      })
    })
  })

  it('should fail to delete client when Authorization header is not presented', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        client_name: 'clientName',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'DELETE',
        failOnStatusCode: false,
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id
      }).then((response) => {
        expect(response.status).to.eq(401)
      })
    })
  })

  it('should delete client when having an valid Authorization header', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/oauth2/register',
      body: {
        client_name: 'clientName',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      cy.request({
        method: 'DELETE',
        failOnStatusCode: false,
        url:
          Cypress.env('public_url') +
          '/oauth2/register/' +
          response.body.client_id,
        headers: {
          Authorization: 'Bearer ' + response.body.registration_access_token
        }
      }).then((response) => {
        expect(response.status).to.eq(204)
      })
    })
  })
})
