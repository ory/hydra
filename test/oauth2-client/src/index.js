const express = require('express')
const session = require('express-session')
const dotenv = require('dotenv')
const uuid = require('node-uuid')
const oauth2 = require('simple-oauth2')
const fetch = require('node-fetch')
const ew = require('express-winston')
const winston = require('winston')
const qs = require('querystring')
const { Issuer } = require('openid-client')
const { URLSearchParams } = require('url')

dotenv.config()

const app = express()

const isStatusOk = (res) => res.ok ? Promise.resolve(res) : Promise.reject(new Error(`Received unexpected status code ${res.statusCode}`))

const config = {
  url: process.env.AUTHORIZATION_SERVER_URL,
  admin: process.env.ADMIN_URL,
  port: parseInt(process.env.PORT)
}
const redirect_uri = `http://127.0.0.1:${config.port}`

app.use(ew.logger({
  transports: [new winston.transports.Console()],
  format: winston.format.combine(
    winston.format.colorize(),
    winston.format.simple()
  ),
}))

app.use(session({
  secret: '804cd9c9-b447-4df0-b9f0-3126893d3a8e',
  resave: false,
  saveUninitialized: true,
  cookie: {
    secure: false,
    httpOnly: true
  }
}))

const nc = (req) => Issuer.discover(config.url).then((issuer) => Promise.resolve(new issuer.Client(req.session.oidc_credentials)))

app.get('/oauth2/code', async (req, res) => {
  const credentials = {
    client: {
      id: req.query.client_id,
      secret: req.query.client_secret,
    },
    auth: {
      tokenHost: config.url,
      tokenPath: '/oauth2/token',
      authorizePath: '/oauth2/auth'
    },
  }

  const state = uuid.v4()
  const scope = req.query.scope || ''

  req.session.credentials = credentials
  req.session.state = state
  req.session.scope = scope.split(' ')

  res.redirect(oauth2.create(credentials).authorizationCode.authorizeURL({
    redirect_uri: `${redirect_uri}/oauth2/callback`,
    scope,
    state
  }))
})

app.get('/oauth2/callback', async (req, res) => {
  if (req.query.error) {
    res.send(JSON.stringify(Object.assign({ result: 'error' }, req.query)))
    return
  }

  if (req.query.state !== req.session.state) {
    res.send(JSON.stringify({ result: 'error', error: 'states mismatch' }))
    return
  }

  if (!req.query.code) {
    res.send(JSON.stringify({ result: 'error', error: 'no code given' }))
    return
  }

  oauth2.create(req.session.credentials).authorizationCode.getToken({
    redirect_uri: `${redirect_uri}/oauth2/callback`,
    scope: req.session.scope,
    code: req.query.code
  }).then((token) => {
    req.session.token = token
    res.send({ result: 'success', token })
  }).catch((err) => {
    if (err.data.payload) {
      res.send(JSON.stringify(err.data.payload))
      return
    }
    res.send(JSON.stringify({ error: err.toString() }))
  })
})

app.get('/oauth2/refresh', function (req, res) {
  oauth2.create(req.session.credentials).accessToken.create(req.session.token).refresh().then((token) => {
    req.session.token = token
    res.send({ result: 'success', token: token.token })
  }).catch((err) => {
    res.send(JSON.stringify({ error: err.toString() }))
  })
})

app.get('/oauth2/revoke', (req, res) => {
  req.session.token.revoke(req.query.type || 'access_token').then(() => {
    res.status(201)
  }).catch((err) => {
    res.send(JSON.stringify({ error: err.toString() }))
  })
})

app.get('/oauth2/introspect/at', (req, res) => {
  const params = new URLSearchParams()
  params.append('token', req.session.token.access_token)

  fetch(new URL('/oauth2/introspect', config.admin).toString(), {
    method: 'POST',
    body: params
  }).then(isStatusOk).then((res) => res.json())
    .then((body) => res.json({ result: 'success', body }))
    .catch((err) => {
      console.log(err)
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get('/oauth2/introspect/rt', async (req, res) => {
  const params = new URLSearchParams()
  params.append('token', req.session.token.refresh_token)

  fetch(new URL('/oauth2/introspect', config.admin).toString(), {
    method: 'POST',
    body: params
  }).then(isStatusOk).then((res) => res.json())
    .then((body) => res.json({ result: 'success', body }))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

// client credentials

app.get('/oauth2/cc', (req, res) => {
  const credentials = {
    client: {
      id: req.query.client_id,
      secret: req.query.client_secret,
    },
    auth: {
      tokenHost: config.url,
      tokenPath: '/oauth2/token',
    },
    options: {
      authorizationMethod: 'header'
    }
  }

  oauth2.create(credentials).clientCredentials.getToken({ scope: req.query.scope.split(' ') }).then((token) => {
    res.send({ result: 'success', token })
  }).catch((err) => {
    if (err.data.payload) {
      res.send(JSON.stringify(err.data.payload))
      return
    }

    res.send(JSON.stringify({ error: err.toString() }))
  })
})

// openid

app.get('/openid/code', async (req, res) => {
  const credentials = {
    client_id: req.query.client_id,
    client_secret: req.query.client_secret
  }

  const state = uuid.v4()
  const nonce = uuid.v4()
  const scope = req.query.scope || ''

  req.session.oidc_credentials = credentials
  req.session.state = state
  req.session.nonce = nonce
  req.session.scope = scope.split(' ')

  const client = await nc(req)
  const url = client.authorizationUrl({
    redirect_uri: `${redirect_uri}/openid/callback`,
    scope: scope,
    state: state,
    nonce: nonce,
    prompt: req.query.prompt
  })
  res.redirect(url)
})

app.get('/openid/callback', async (req, res) => {
  if (req.query.error) {
    res.send(JSON.stringify(Object.assign({ result: 'error' }, req.query)))
    return
  }

  if (req.query.state !== req.session.state) {
    res.send(JSON.stringify({ result: 'error', error: 'states mismatch' }))
    return
  }

  if (!req.query.code) {
    res.send(JSON.stringify({ result: 'error', error: 'no code given' }))
    return
  }

  const client = await nc(req)
  client.authorizationCallback(`${redirect_uri}/openid/callback`, req.query, {
    state: req.session.state,
    nonce: req.session.nonce,
    response_type: 'code'
  }).then((ts) => {
    req.session.openid_token = ts
    res.send({ result: 'success', token: ts, claims: ts.claims })
  }).catch((err) => {
    console.log(err)
    res.send(JSON.stringify({ error: err.toString() }))
  })
})

app.get('/openid/userinfo', async (req, res) => {
  const client = await nc(req)
  client.userinfo(req.session.openid_token.access_token)
    .then((ui) => res.json(ui))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get('/openid/revoke/at', async (req, res) => {
  const client = await nc(req)
  client.revoke(req.session.openid_token.access_token)
    .then(() => res.json({result: 'success'}))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get('/openid/revoke/rt', async (req, res) => {
  const client = await nc(req)
  client.revoke(req.session.openid_token.refresh_token)
    .then(() => res.json({result: 'success'}))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.listen(config.port, function () {
  console.log(`Listening on port ${config.port}!`)
})

