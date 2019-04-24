import {prng} from '../../helpers';

const nc = () => ({
  client_id: prng(),
  client_secret: prng(),
  scope: 'openid',
  subject_type: 'public',
  redirect_uris: [`${Cypress.env('client_url')}/openid/callback`],
  grant_types: ['authorization_code']
});

describe('OpenID Connect Logout', () => {
  describe('Front-Channel', () => {
    beforeEach(() => {
      Cypress.Cookies.preserveOnce('oauth2_authentication_session', 'connect.sid')
    })

    it('should log in and remember login with front-channel', function () {
      const client = {
        frontchannel_logout_uri: `${Cypress.env(
          'client_url'
        )}/openid/session/end/fc`,
        ...nc()
      };

      cy.authCodeFlow(
        client,
        {
          login: {remember: true},
          consent: {scope: ['openid'], remember: true}
        },
        'openid'
      );

      cy.request(`${Cypress.env('client_url')}/openid/session/check`).its('body').then(({has_session})=>{
        expect(has_session).to.be.true
      })
    });

    it('should show the logout page and complete logout with front-channel', () => {
      cy.visit(
        `${Cypress.env('client_url')}/openid/session/end`,
        {failOnStatusCode: false}
      );

      cy.get('#accept').click();

      cy.get('h1').should('contain', 'Your log out request however succeeded.')

      cy.request(`${Cypress.env('client_url')}/openid/session/check`).its('body').then(({has_session})=>{
        expect(has_session).to.be.false
      })
    })
  });

  // describe('Back-Channel', () => {
  //   beforeEach(() => {
  //     Cypress.Cookies.preserveOnce('oauth2_authentication_session', 'connect.sid')
  //   })
  //
  //   it('should log in and remember login with back-channel', function () {
  //     const client = {
  //       backchannel_logout_uri: `${Cypress.env(
  //         'client_url'
  //       )}/openid/session/end/bc`,
  //       ...nc()
  //     };
  //
  //     cy.authCodeFlow(
  //       client,
  //       {
  //         login: {remember: true},
  //         consent: {scope: ['openid'], remember: true}
  //       },
  //       'openid'
  //     );
  //
  //     cy.request(`${Cypress.env('client_url')}/openid/session/check`).its('body').then(({has_session})=>{
  //       expect(has_session).to.be.true
  //     })
  //   });
  //
  //   it('should show the logout page and complete logout with back-channel', () => {
  //     cy.visit(
  //       `${Cypress.env('client_url')}/openid/session/end`,
  //       {failOnStatusCode: false}
  //     );
  //
  //     cy.get('#accept').click();
  //
  //     cy.get('h1').should('contain', 'Your log out request however succeeded.')
  //
  //     cy.request(`${Cypress.env('client_url')}/openid/session/check`).its('body').then(({has_session})=>{
  //       expect(has_session).to.be.false
  //     })
  //   })
  // });
});
