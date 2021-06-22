import { prng } from '../../helpers'

describe('The Clients Pubic Interface', function () {

  it('should return same client_secret given in request for newly created clients with client_secret specified', function () {
    cy.request({
      method: 'POST',
      url: Cypress.env('public_url') + '/connect/register',
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
      expect(response.body.registration_access_token).to.not.be.empty
    })
  })

  it('should get client when having a valid client_secret in body', function () {
    cy.request({
      method: 'GET',
      url: Cypress.env('public_url') + '/connect/register?client_id=clientid',
      headers: {
        'Authorization' : 'Basic Y2xpZW50aWQ6c2VjcmV0'
      }
    }).then((response) => {
      console.log(response.body)
      expect(response.body.client_name).to.equal('clientName')
    })
  })

  it('should fail for get client when Authorization header is not presented', function () {
    cy.request({
      method: 'GET',
      failOnStatusCode: false,
      url:
        Cypress.env('public_url') + '/connect/register?client_id=clientid'
    }).then((response) => {
      expect(response.status).to.eq(401)
    })
  })

  it('should update client name when having a valid client_secret in body', function () {
    cy.request({
      method: 'PUT',
      url: Cypress.env('public_url') + '/connect/register?client_id=clientid',
      headers: {
        'Authorization' : 'Basic Y2xpZW50aWQ6c2VjcmV0'
      },
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

  it('should fail to update client name when Authorization header is not presented', function () {
    cy.request({
      method: 'PUT',
      failOnStatusCode: false,
      url: Cypress.env('public_url') + '/connect/register?client_id=clientid',
      body: {
        client_id: 'clientid',
        client_name: 'clientName2',
        scope: 'foo openid offline_access',
        grant_types: ['client_credentials']
      }
    }).then((response) => {
      expect(response.status).to.eq(401)
    })
  })

  it('should fail to delete client when Authorization header is not presented', function () {
    cy.request({
      method: 'DELETE',
      failOnStatusCode: false,
      url:
        Cypress.env('public_url') + '/connect/register?client_id=clientid'
    }).then((response) => {
      expect(response.status).to.eq(401)
    })
  })

  it('should delete client when having an valid Authorization header', function () {
    cy.request({
      method: 'DELETE',
      failOnStatusCode: false,
      url: Cypress.env('public_url') + '/connect/register?client_id=clientid',
      headers: {
        'Authorization' : 'Basic Y2xpZW50aWQ6c2VjcmV0'
      }
    }).then((response) => {
      expect(response.status).to.eq(204)
    })
  })
})
