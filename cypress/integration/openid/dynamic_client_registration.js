import { prng } from '../../helpers'

describe('The Clients Pubic Interface', function () {
  it('should return same client_secret given in request for newly created clients with client_secret specified', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/dyn-clients',
      body: {
        client_id: 'clientid',
        client_name: 'clientName',
        client_secret: 'secret',
        scope: 'foo openid offline_access',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      console.log(response.body)
      expect(response.body.client_secret).to.equal('secret')
    })
  })

  it('should get client when having a valid client_secret in body', function () {
    cy.request({
      method: 'GET',
      url: Cypress.env('public_url') + '/dyn-clients/clientid?secret=secret'
    }).then((response) => {
      console.log(response.body)
      expect(response.body.client_name).to.equal('clientName')
    })
  })

  it('should fail for get client when having an invalid client_secret in body', function () {
    cy.request({
      method: 'GET',
      failOnStatusCode: false,
      url:
        Cypress.env('public_url') + '/dyn-clients/clientid?secret=wrongsecret'
    }).then((response) => {
      expect(response.status).to.eq(401)
    })
  })

  it('should update client name when having a valid client_secret in body', function () {
    cy.request({
      method: 'PUT',
      url: Cypress.env('public_url') + '/dyn-clients/clientid',
      body: {
        client_id: 'clientid',
        client_name: 'clientName2',
        client_secret: 'secret',
        scope: 'foo openid offline_access',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      console.log(response.body)
      expect(response.body.client_name).to.equal('clientName2')
    })
  })

  it('should fail to update client name when having an invalid client_secret in body', function () {
    cy.request({
      method: 'PUT',
      failOnStatusCode: false,
      url: Cypress.env('public_url') + '/dyn-clients/clientid',
      body: {
        client_id: 'clientid',
        client_name: 'clientName2',
        client_secret: 'wrongsecret',
        scope: 'foo openid offline_access',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      expect(response.status).to.eq(401)
    })
  })

  it('should fail to delete client when having an invalid client_secret as parameter body', function () {
    cy.request({
      method: 'DELETE',
      failOnStatusCode: false,
      url:
        Cypress.env('public_url') + '/dyn-clients/clientid?secret=wrongsecret'
    }).then((response) => {
      expect(response.status).to.eq(401)
    })
  })

  it('should delete client when having an valid client_secret as parameter body', function () {
    cy.request({
      method: 'DELETE',
      failOnStatusCode: false,
      url: Cypress.env('public_url') + '/dyn-clients/clientid?secret=secret'
    }).then((response) => {
      expect(response.status).to.eq(204)
    })
  })
})
