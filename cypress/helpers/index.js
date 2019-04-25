export const prng = () =>
  `${Math.random()
    .toString(36)
    .substring(2)}${Math.random()
    .toString(36)
    .substring(2)}`;

const isStatusOk = res =>
  res.ok
    ? Promise.resolve(res)
    : Promise.reject(
        new Error(`Received unexpected status code ${res.statusCode}`)
      );

export const createClient = client =>
  fetch(Cypress.env('admin_url') + '/clients', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(client)
  })
    .then(isStatusOk)
    .then(res => {
      return res.json();
    })
    .then(body =>
      fetch(Cypress.env('admin_url') + '/clients/' + client.client_id)
        .then(isStatusOk)
        .then(res => {
          return res.json();
        })
        .then(actual => {
          if (actual.client_id !== body.client_id) {
            return Promise.reject(
              new Error(
                `Expected client_id's to match: ${actual.client_id} !== ${
                  body.client
                }`
              )
            );
          }

          return Promise.resolve(body);
        })
    );
