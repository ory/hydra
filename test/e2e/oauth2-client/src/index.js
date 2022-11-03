// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

const express = require("express")
const session = require("express-session")
const uuid = require("node-uuid")
const oauth2 = require("simple-oauth2")
const fetch = require("node-fetch")
const ew = require("express-winston")
const winston = require("winston")
const { Issuer } = require("openid-client")
const { URLSearchParams } = require("url")
const bodyParser = require("body-parser")
const jwksClient = require("jwks-rsa")
const jwt = require("jsonwebtoken")

const app = express()

app.use(bodyParser.urlencoded({ extended: true }))

const blacklistedSid = []

const isStatusOk = (res) =>
  res.ok
    ? Promise.resolve(res)
    : Promise.reject(
        new Error(`Received unexpected status code ${res.statusCode}`),
      )

const config = {
  url: process.env.AUTHORIZATION_SERVER_URL || "http://127.0.0.1:5004/",
  public: process.env.PUBLIC_URL || "http://127.0.0.1:5004/",
  admin: process.env.ADMIN_URL || "http://127.0.0.1:5001/",
  port: parseInt(process.env.PORT) || 5003,
}

const redirect_uri = `http://127.0.0.1:${config.port}`

app.use(
  ew.logger({
    transports: [new winston.transports.Console()],
    format: winston.format.combine(
      winston.format.colorize(),
      winston.format.simple(),
    ),
  }),
)

app.use(
  session({
    secret: "804cd9c9-b447-4df0-b9f0-3126893d3a8e",
    resave: false,
    saveUninitialized: true,
    cookie: {
      secure: false,
      httpOnly: true,
    },
  }),
)

const nc = (req) =>
  Issuer.discover(config.public).then((issuer) => {
    // This is necessary when working with docker...
    issuer.metadata.token_endpoint = new URL(
      "/oauth2/token",
      config.public,
    ).toString()
    issuer.metadata.jwks_uri = new URL(
      "/.well-known/jwks.json",
      config.public,
    ).toString()
    issuer.metadata.revocation_endpoint = new URL(
      "/oauth2/revoke",
      config.public,
    ).toString()
    issuer.metadata.introspection_endpoint = new URL(
      "/oauth2/introspect",
      config.admin,
    ).toString()

    return Promise.resolve(
      new issuer.Client({
        ...issuer.metadata,
        ...req.session.oidc_credentials,
      }),
    )
  })

app.get("/oauth2/code", async (req, res) => {
  const credentials = {
    client: {
      id: req.query.client_id,
      secret: req.query.client_secret,
    },
    auth: {
      tokenHost: config.public,
      authorizeHost: config.url,
      tokenPath: "/oauth2/token",
      authorizePath: "/oauth2/auth",
    },
  }

  const state = uuid.v4()
  const scope = req.query.scope || ""

  req.session.credentials = credentials
  req.session.state = state
  req.session.scope = scope.split(" ")

  res.redirect(
    oauth2.create(credentials).authorizationCode.authorizeURL({
      redirect_uri: `${redirect_uri}/oauth2/callback`,
      scope,
      state,
    }),
  )
})

app.get("/oauth2/callback", async (req, res) => {
  if (req.query.error) {
    res.send(JSON.stringify(Object.assign({ result: "error" }, req.query)))
    return
  }

  if (req.query.state !== req.session.state) {
    res.send(JSON.stringify({ result: "error", error: "states mismatch" }))
    return
  }

  if (!req.query.code) {
    res.send(JSON.stringify({ result: "error", error: "no code given" }))
    return
  }

  oauth2
    .create(req.session.credentials)
    .authorizationCode.getToken({
      redirect_uri: `${redirect_uri}/oauth2/callback`,
      scope: req.session.scope,
      code: req.query.code,
    })
    .then((token) => {
      req.session.oauth2_flow = { token } // code returns {access_token} because why not...
      res.send({ result: "success", token })
    })
    .catch((err) => {
      if (err.data.payload) {
        res.send(JSON.stringify(err.data.payload))
        return
      }
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/oauth2/refresh", function (req, res) {
  oauth2
    .create(req.session.credentials)
    .accessToken.create(req.session.oauth2_flow.token)
    .refresh()
    .then((token) => {
      req.session.oauth2_flow = token // refresh returns {token:{access_token}} because why not...
      res.send({ result: "success", token: token.token })
    })
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/oauth2/revoke", (req, res) => {
  oauth2
    .create(req.session.credentials)
    .accessToken.create(req.session.oauth2_flow.token)
    .revoke(req.query.type || "access_token")
    .then(() => {
      res.sendStatus(201)
    })
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/oauth2/validate-jwt", (req, res) => {
  const client = jwksClient({
    jwksUri: new URL("/.well-known/jwks.json", config.public).toString(),
  })

  jwt.verify(
    req.session.oauth2_flow.token.access_token,
    (header, callback) => {
      client.getSigningKey(header.kid, function (err, key) {
        const signingKey = key.publicKey || key.rsaPublicKey
        callback(null, signingKey)
      })
    },
    (err, decoded) => {
      if (err) {
        console.error(err)
        res.send(400)
        return
      }

      res.send(decoded)
    },
  )
})

app.get("/oauth2/introspect/at", (req, res) => {
  const params = new URLSearchParams()
  params.append("token", req.session.oauth2_flow.token.access_token)

  fetch(new URL("/oauth2/introspect", config.admin).toString(), {
    method: "POST",
    body: params,
  })
    .then(isStatusOk)
    .then((res) => res.json())
    .then((body) => res.json({ result: "success", body }))
    .catch((err) => {
      console.error(err)
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/oauth2/introspect/rt", async (req, res) => {
  const params = new URLSearchParams()
  params.append("token", req.session.oauth2_flow.token.refresh_token)

  fetch(new URL("/oauth2/introspect", config.admin).toString(), {
    method: "POST",
    body: params,
  })
    .then(isStatusOk)
    .then((res) => res.json())
    .then((body) => res.json({ result: "success", body }))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

// client credentials

app.get("/oauth2/cc", (req, res) => {
  const credentials = {
    client: {
      id: req.query.client_id,
      secret: req.query.client_secret,
    },
    auth: {
      tokenHost: config.public,
      tokenPath: "/oauth2/token",
    },
    options: {
      authorizationMethod: "header",
    },
  }

  oauth2
    .create(credentials)
    .clientCredentials.getToken({ scope: req.query.scope.split(" ") })
    .then((token) => {
      res.send({ result: "success", token })
    })
    .catch((err) => {
      if (err.data.payload) {
        res.send(JSON.stringify(err.data.payload))
        return
      }

      res.send(JSON.stringify({ error: err.toString() }))
    })
})

// openid

app.get("/openid/code", async (req, res) => {
  const credentials = {
    client_id: req.query.client_id,
    client_secret: req.query.client_secret,
  }

  const state = uuid.v4()
  const nonce = uuid.v4()
  const scope = req.query.scope || ""

  req.session.oidc_credentials = credentials
  req.session.state = state
  req.session.nonce = nonce
  req.session.scope = scope.split(" ")

  const client = await nc(req)
  const url = client.authorizationUrl({
    redirect_uri: `${redirect_uri}/openid/callback`,
    scope: scope,
    state: state,
    nonce: nonce,
    prompt: req.query.prompt,
  })
  res.redirect(url)
})

app.get("/openid/callback", async (req, res) => {
  if (req.query.error) {
    res.send(JSON.stringify(Object.assign({ result: "error" }, req.query)))
    return
  }

  if (req.query.state !== req.session.state) {
    res.send(JSON.stringify({ result: "error", error: "states mismatch" }))
    return
  }

  if (!req.query.code) {
    res.send(JSON.stringify({ result: "error", error: "no code given" }))
    return
  }

  const client = await nc(req)
  client
    .authorizationCallback(`${redirect_uri}/openid/callback`, req.query, {
      state: req.session.state,
      nonce: req.session.nonce,
      response_type: "code",
    })
    .then((ts) => {
      req.session.openid_token = ts
      req.session.openid_claims = ts.claims
      res.send({ result: "success", token: ts, claims: ts.claims })
    })
    .catch((err) => {
      console.error(err)
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/openid/userinfo", async (req, res) => {
  const client = await nc(req)
  client
    .userinfo(req.session.openid_token.access_token)
    .then((ui) => res.json(ui))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/openid/revoke/at", async (req, res) => {
  const client = await nc(req)
  client
    .revoke(req.session.openid_token.access_token)
    .then(() => res.json({ result: "success" }))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/openid/revoke/rt", async (req, res) => {
  const client = await nc(req)
  client
    .revoke(req.session.openid_token.refresh_token)
    .then(() => res.json({ result: "success" }))
    .catch((err) => {
      res.send(JSON.stringify({ error: err.toString() }))
    })
})

app.get("/openid/session/end", async (req, res) => {
  const client = await nc(req)
  const state = uuid.v4()

  if (req.query.simple) {
    res.redirect(new URL("/oauth2/sessions/logout", config.public).toString())
  } else {
    req.session.logout_state = state
    res.redirect(
      client.endSessionUrl({
        state,
        id_token_hint:
          req.query.id_token_hint || req.session.openid_token.id_token,
      }),
    )
  }
})

app.get("/openid/session/end/fc", async (req, res) => {
  if (req.session.openid_claims.sid !== req.query.sid) {
    res.sendStatus(400)
    return
  }

  if (req.session.openid_claims.iss !== req.query.iss) {
    res.sendStatus(400)
    return
  }

  setTimeout(() => {
    req.session.destroy(() => {
      res.send("ok")
    })
  }, 500)
})

app.post("/openid/session/end/bc", (req, res) => {
  const client = jwksClient({
    jwksUri: new URL("/.well-known/jwks.json", config.public).toString(),
    cache: false,
  })

  jwt.verify(
    req.body.logout_token,
    (header, callback) => {
      client.getSigningKey(header.kid, (err, key) => {
        if (err) {
          console.error(err)
          res.sendStatus(400)
          return
        }

        callback(null, key.publicKey || key.rsaPublicKey)
      })
    },
    (err, decoded) => {
      if (err) {
        console.error(err)
        res.sendStatus(400)
        return
      }

      if (decoded.nonce) {
        console.error("nonce is set but should not be", decoded.nonce)
        res.sendStatus(400)
        return
      }

      if (decoded.sid.length === 0) {
        console.error("sid should be set but is not", decoded.sid)
        res.sendStatus(400)
        return
      }

      if (decoded.iss.indexOf(config.url) === -1) {
        console.error("issuer is mismatching", decoded.iss, config.url)
        res.sendStatus(400)
        return
      }

      blacklistedSid.push(decoded.sid)
      res.send("ok")
    },
  )
})

app.get("/openid/session/check", async (req, res) => {
  const { openid_claims: { sid = "" } = {} } = req.session

  if (blacklistedSid.indexOf(sid) > -1) {
    req.session.destroy(() => {
      res.json({ has_session: false })
    })
    return
  }

  res.json({
    has_session:
      Boolean(req.session.oauth2_flow) ||
      (Boolean(req.session.openid_token) && Boolean(req.session.openid_claims)),
  })
})

app.get("/empty", (req, res) => {
  res.setHeader("Content-Type", "text/html")
  res.send(Buffer.from("<div>Nothing to see here.</div>"))
})

app.listen(config.port, function () {
  console.log(`Listening on port ${config.port}!`)
})
