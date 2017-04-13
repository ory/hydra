
const request = require('request');

let accessToken;

const getWardenToken = (done) => {
  const wardenCreds = new Buffer(process.env.HYDRA_CLIENT_ID + ':' + process.env.HYDRA_CLIENT_SECRET).toString('base64')
  const tokenForm = { grant_type: 'client_credentials', scope: 'hydra' };
  const reqOptions = {
    url: process.env.HYDRA_URL + '/oauth2/token',
    headers: {
      'Authorization': 'Basic ' + wardenCreds
    },
    method: 'POST',
    form: tokenForm
  };

  return request(reqOptions, (err, res, body) => {
    accessToken = JSON.parse(body).access_token;
    done();
  });
};

const makeWardenReq = (subject, action, resource, context, cb) => {
  const wardenBody = { resource, action, subject, context: context || {} };
  const reqOptions = {
    url: process.env.HYDRA_URL + '/warden/allowed',
    headers: {},
    method: 'POST',
    json: true,
    body: wardenBody
  }

  const tokenCb = () => {
    reqOptions.headers.Authorization = 'Bearer ' + accessToken;
    return request(reqOptions, cb);
  };

  if (!accessToken) {
    return getWardenToken(tokenCb);
  };
  tokenCb();
};

module.exports = {
  makeWardenReq
};

