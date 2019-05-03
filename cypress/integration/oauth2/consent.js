import { prng } from '../../helpers';

describe('OAuth 2.0 End-User Authorization', () => {
  const nc = () => ({
    client_id: prng(),
    client_secret: prng(),
    scope: 'offline_access',
    redirect_uris: [`${Cypress.env('client_url')}/oauth2/callback`],
    grant_types: ['authorization_code', 'refresh_token']
  });

  const hasConsent = (client, body) => {
    let found = false;
    body.forEach(({ consent_request: { client: { client_id } } }) => {
      if (client_id === client.client_id) {
        found = true;
      }
    });
    return found;
  };

  it('should check if end user authorization exists', () => {
    const client = nc();
    cy.authCodeFlow(client, {
      consent: {
        scope: ['offline_access'],
        remember: true
      }
    });

    cy.request(
      Cypress.env('admin_url') +
        '/oauth2/auth/sessions/consent?subject=foo@bar.com'
    )
      .its('body')
      .then(body => {
        expect(body.length).to.be.greaterThan(0);
        expect(hasConsent(client, body)).to.be.true;
      });

    cy.request(
      'DELETE',
      Cypress.env('admin_url') +
        '/oauth2/auth/sessions/consent?subject=foo@bar.com'
    );

    cy.request(
      Cypress.env('admin_url') +
        '/oauth2/auth/sessions/consent?subject=foo@bar.com'
    )
      .its('body')
      .then(body => {
        expect(body.length).to.eq(0);
        expect(hasConsent(client, body)).to.be.false;
      });

    cy.request(`${Cypress.env('client_url')}/oauth2/introspect/at`)
      .its('body')
      .then(body => {
        expect(body.result).to.equal('success');
        expect(body.body.active).to.be.false;
      });

    cy.request(`${Cypress.env('client_url')}/oauth2/introspect/rt`)
      .its('body')
      .then(body => {
        expect(body.result).to.equal('success');
        expect(body.body.active).to.be.false;
      });
  });
});
