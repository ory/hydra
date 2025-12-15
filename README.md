<h1 align="center">
  <img src="https://raw.githubusercontent.com/ory/meta/master/static/banners/hydra.svg" alt="Ory Hydra - Open Source OAuth 2 and OpenID Connect server">
</h1>

<h4 align="center">
  <a href="https://www.ory.sh/chat">Chat</a> ·
  <a href="https://github.com/ory/hydra/discussions">Discussions</a> ·
  <a href="https://www.ory.sh/l/sign-up-newsletter">Newsletter</a> ·
  <a href="https://www.ory.sh/docs/">Docs</a> ·
  <a href="https://console.ory.sh/">Try Ory Network</a> ·
  <a href="https://www.ory.sh/jobs/">Jobs</a>
</h4>

Ory Hydra is a hardened, OpenID Certified OAuth 2.0 Server and OpenID Connect
Provider optimized for low-latency, high throughput, and low resource
consumption. It connects to your existing identity provider through a login and
consent app, giving you absolute control over the user interface and experience.

---

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [What is Ory Hydra?](#what-is-ory-hydra)
  - [Why Ory Hydra](#why-ory-hydra)
  - [OAuth2 and OpenID Connect: Open Standards](#oauth2-and-openid-connect-open-standards)
  - [OpenID Connect Certified](#openid-connect-certified)
- [Deployment options](#deployment-options)
  - [Use Ory Hydra on the Ory Network](#use-ory-hydra-on-the-ory-network)
  - [Self-host Ory Hydra](#self-host-ory-hydra)
- [Quickstart](#quickstart)
- [Who is using Ory Hydra](#who-is-using-ory-hydra)
- [Ecosystem](#ecosystem)
  - [Ory Kratos: Identity and User Infrastructure and Management](#ory-kratos-identity-and-user-infrastructure-and-management)
  - [Ory Hydra: OAuth2 & OpenID Connect Server](#ory-hydra-oauth2--openid-connect-server)
  - [Ory Oathkeeper: Identity & Access Proxy](#ory-oathkeeper-identity--access-proxy)
  - [Ory Keto: Access Control Policies as a Server](#ory-keto-access-control-policies-as-a-server)
- [Documentation](#documentation)
- [Developing Ory Hydra](#developing-ory-hydra)
- [Security](#security)
  - [Disclosing vulnerabilities](#disclosing-vulnerabilities)
- [Telemetry](#telemetry)
- [Libraries and third-party projects](#libraries-and-third-party-projects)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Ory Hydra?

Ory Hydra is a server implementation of the OAuth 2.0 authorization framework
and the OpenID Connect Core 1.0. It follows
[cloud architecture best practices](https://www.ory.sh/docs/ecosystem/software-architecture-philosophy)
and focuses on:

- OAuth 2.0 and OpenID Connect flows
- Token issuance and validation
- Client management
- Consent and login flow orchestration
- JWKS management
- Low latency and high throughput

We recommend starting with the
[Ory Hydra introduction docs](https://www.ory.sh/docs/hydra) to learn more about
its architecture, feature set, and how it compares to other systems.

### Why Ory Hydra

Ory Hydra is designed to:

- Be a standalone OAuth 2.0 and OpenID Connect server without user management
- Connect to any existing identity provider through a login and consent app
- Give you absolute control over the user interface and experience flows
- Work with any authentication endpoint:
  [Ory Kratos](https://github.com/ory/kratos),
  [authboss](https://github.com/go-authboss/authboss),
  [User Frosting](https://www.userfrosting.com/), or your proprietary system
- Scale to large numbers of clients and tokens
- Fit into modern cloud native environments such as Kubernetes and managed
  platforms

### OAuth2 and OpenID Connect: Open Standards

Ory Hydra implements Open Standards set by the IETF:

- [The OAuth 2.0 Authorization Framework](https://tools.ietf.org/html/rfc6749)
- [OAuth 2.0 Threat Model and Security Considerations](https://tools.ietf.org/html/rfc6819)
- [OAuth 2.0 Token Revocation](https://tools.ietf.org/html/rfc7009)
- [OAuth 2.0 Token Introspection](https://tools.ietf.org/html/rfc7662)
- [OAuth 2.0 for Native Apps](https://tools.ietf.org/html/draft-ietf-oauth-native-apps-10)
- [OAuth 2.0 Dynamic Client Registration Protocol](https://datatracker.ietf.org/doc/html/rfc7591)
- [OAuth 2.0 Dynamic Client Registration Management Protocol](https://datatracker.ietf.org/doc/html/rfc7592)
- [Proof Key for Code Exchange by OAuth Public Clients](https://tools.ietf.org/html/rfc7636)
- [JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7523)

and the OpenID Foundation:

- [OpenID Connect Core 1.0](http://openid.net/specs/openid-connect-core-1_0.html)
- [OpenID Connect Discovery 1.0](https://openid.net/specs/openid-connect-discovery-1_0.html)
- [OpenID Connect Dynamic Client Registration 1.0](https://openid.net/specs/openid-connect-registration-1_0.html)
- [OpenID Connect Front-Channel Logout 1.0](https://openid.net/specs/openid-connect-frontchannel-1_0.html)
- [OpenID Connect Back-Channel Logout 1.0](https://openid.net/specs/openid-connect-backchannel-1_0.html)

### OpenID Connect Certified

Ory Hydra is an OpenID Foundation
[certified OpenID Provider (OP)](http://openid.net/certification/#OPs).

<p align="center">
    <img src="https://github.com/ory/docs/blob/master/docs/hydra/images/oidc-cert.png" alt="Ory Hydra is a certified OpenID Providier" width="256px">
</p>

The following OpenID profiles are certified:

- [Basic OpenID Provider](http://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth)
  (response types `code`)
- [Implicit OpenID Provider](http://openid.net/specs/openid-connect-core-1_0.html#ImplicitFlowAuth)
  (response types `id_token`, `id_token+token`)
- [Hybrid OpenID Provider](http://openid.net/specs/openid-connect-core-1_0.html#HybridFlowAuth)
  (response types `code+id_token`, `code+id_token+token`, `code+token`)
- [OpenID Provider Publishing Configuration Information](https://openid.net/specs/openid-connect-discovery-1_0.html)
- [Dynamic OpenID Provider](https://openid.net/specs/openid-connect-registration-1_0.html)

To obtain certification, we deployed the
[reference user login and consent app](https://github.com/ory/hydra-login-consent-node)
(unmodified) and Ory Hydra v1.0.0.

## Deployment options

You can run Ory Hydra in two main ways:

- As a managed service on the Ory Network
- As a self hosted service under your own control, with or without the Ory
  Enterprise License

### Use Ory Hydra on the Ory Network

The [Ory Network](https://www.ory.sh/cloud) is the fastest way to use Ory
services in production. **Ory OAuth2 & OpenID Connect** is powered by the open
source Ory Hydra server and is API compatible.

The Ory Network provides:

- OAuth2 and OpenID Connect for single sign on, API access, and machine to
  machine authorization
- Identity and credential management that scales to billions of users and
  devices
- Registration, login, and account management flows for passkeys, biometrics,
  social login, SSO, and multi factor authentication
- Prebuilt login, registration, and account management pages and components
- Low latency permission checks based on the Zanzibar model with the Ory
  Permission Language
- GDPR friendly storage with data locality and compliance in mind
- Web based Ory Console and Ory CLI for administration and operations
- Cloud native APIs compatible with the open source servers
- Fair, usage based [pricing](https://www.ory.sh/pricing)

Sign up for a
[free developer account](https://console.ory.sh/registration?utm_source=github&utm_medium=banner&utm_campaign=hydra-readme)
to get started.

### Self-host Ory Hydra

You can run Ory Hydra yourself for full control over infrastructure, deployment,
and customization.

The [install guide](https://www.ory.sh/docs/hydra/install) explains how to:

- Install Hydra on Linux, macOS, Windows, and Docker
- Configure databases such as PostgreSQL, MySQL, and CockroachDB
- Deploy to Kubernetes and other orchestration systems
- Build Hydra from source

This guide uses the open source distribution to get you started without license
requirements. It is a great fit for individuals, researchers, hackers, and
companies that want to experiment, prototype, or run unimportant workloads
without SLAs. You get the full core engine, and you are free to inspect, extend,
and build it from source.

If you run Hydra as part of a business-critical system, for example OAuth2 and
OpenID Connect for all your users, you should use a commercial agreement to
reduce operational and security risk. The **Ory Enterprise License (OEL)**
layers on top of self-hosted Hydra and provides:

- Additional enterprise features that are not available in the open source
  version
- Regular security releases, including CVE patches, with service level
  agreements
- Support for advanced scaling, multi-tenancy, and complex deployments
- Premium support options with SLAs, direct access to engineers, and onboarding
  help
- Access to a private Docker registry with frequent and vetted, up-to-date
  enterprise builds

For guaranteed CVE fixes, current enterprise builds, advanced features, and
support in production, you need a valid
[Ory Enterprise License](https://www.ory.com/ory-enterprise-license) and access
to the Ory Enterprise Docker registry. To learn more,
[contact the Ory team](https://www.ory.sh/contact/).

## Quickstart

Install the [Ory CLI](https://www.ory.sh/docs/guides/cli/installation) and
create a new project to try Ory OAuth2 & OpenID Connect.

```bash
# Install the Ory CLI if you do not have it yet:
bash <(curl https://raw.githubusercontent.com/ory/meta/master/install.sh) -b . ory
sudo mv ./ory /usr/local/bin/

# Sign in or sign up
ory auth

# Create a new project
ory create project --create-workspace "Ory Open Source" --name "GitHub Quickstart" --use-project
```

Try out the OAuth 2.0 Client Credentials flow:

```bash
ory create oauth2-client \
    --name "Client Credentials Demo" \
    --grant-type client_credentials
# Note the client ID and secret from output

ory perform client-credentials \
    --client-id <your-client-id> \
    --client-secret <your-client-secret>
# Note the access token from output

ory introspect token <your-access-token>
```

Try out the OAuth 2.0 Authorize Code + OpenID Connect flow:

```bash
ory create oauth2-client \
    --name "Authorize Code with OpenID Connect Demo" \
    --grant-type authorization_code,refresh_token \
    --response-type code \
    --redirect-uri http://127.0.0.1:4446/callback

ory perform authorization-code \
    --client-id <your-client-id> \
    --client-secret <your-client-secret>
```

## Who is using Ory Hydra

<!--BEGIN ADOPTERS-->

The Ory community stands on the shoulders of individuals, companies, and
maintainers. The Ory team thanks everyone involved - from submitting bug reports
and feature requests, to contributing patches and documentation. The Ory
community counts more than 50.000 members and is growing. The Ory stack protects
7.000.000.000+ API requests every day across thousands of companies. None of
this would have been possible without each and everyone of you!

The following list represents companies that have accompanied us along the way
and that have made outstanding contributions to our ecosystem. _If you think
that your company deserves a spot here, reach out to
<a href="mailto:office@ory.sh">office@ory.sh</a> now_!

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Logo</th>
            <th>Website</th>
            <th>Case Study</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>OpenAI</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/openai.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/openai.svg" alt="OpenAI">
                </picture>
            </td>
            <td><a href="https://openai.com/">openai.com</a></td>
            <td><a href="https://www.ory.sh/case-studies/openai">OpenAI Case Study</a></td>
        </tr>
        <tr>
            <td>Fandom</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/fandom.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/fandom.svg" alt="Fandom">
                </picture>
            </td>
            <td><a href="https://www.fandom.com/">fandom.com</a></td>
            <td><a href="https://www.ory.sh/case-studies/fandom">Fandom Case Study</a></td>
        </tr>
        <tr>
            <td>Lumin</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/lumin.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/lumin.svg" alt="Lumin">
                </picture>
            </td>
            <td><a href="https://www.luminpdf.com/">luminpdf.com</a></td>
            <td><a href="https://www.ory.sh/case-studies/lumin">Lumin Case Study</a></td>
        </tr>
        <tr>
            <td>Sencrop</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/sencrop.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/sencrop.svg" alt="Sencrop">
                </picture>
            </td>
            <td><a href="https://sencrop.com/">sencrop.com</a></td>
            <td><a href="https://www.ory.sh/case-studies/sencrop">Sencrop Case Study</a></td>
        </tr>
        <tr>
            <td>OSINT Industries</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/osint.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/osint.svg" alt="OSINT Industries">
                </picture>
            </td>
            <td><a href="https://www.osint.industries/">osint.industries</a></td>
            <td><a href="https://www.ory.sh/case-studies/osint">OSINT Industries Case Study</a></td>
        </tr>
        <tr>
            <td>HGV</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/hgv.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/hgv.svg" alt="HGV">
                </picture>
            </td>
            <td><a href="https://www.hgv.it/">hgv.it</a></td>
            <td><a href="https://www.ory.sh/case-studies/hgv">HGV Case Study</a></td>
        </tr>
        <tr>
            <td>Maxroll</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/maxroll.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/maxroll.svg" alt="Maxroll">
                </picture>
            </td>
            <td><a href="https://maxroll.gg/">maxroll.gg</a></td>
            <td><a href="https://www.ory.sh/case-studies/maxroll">Maxroll Case Study</a></td>
        </tr>
        <tr>
            <td>Zezam</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/zezam.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/zezam.svg" alt="Zezam">
                </picture>
            </td>
            <td><a href="https://www.zezam.io/">zezam.io</a></td>
            <td><a href="https://www.ory.sh/case-studies/zezam">Zezam Case Study</a></td>
        </tr>
        <tr>
            <td>T.RowePrice</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/troweprice.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/troweprice.svg" alt="T.RowePrice">
                </picture>
            </td>
            <td><a href="https://www.troweprice.com/">troweprice.com</a></td>
        </tr>
        <tr>
            <td>Mistral</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/mistral.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/mistral.svg" alt="Mistral">
                </picture>
            </td>
            <td><a href="https://www.mistral.ai/">mistral.ai</a></td>
        </tr>
        <tr>
            <td>Axel Springer</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/axelspringer.svg" />
                    <img height="22px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/axelspringer.svg" alt="Axel Springer">
                </picture>
            </td>
            <td><a href="https://www.axelspringer.com/">axelspringer.com</a></td>
        </tr>
        <tr>
            <td>Hemnet</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/hemnet.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/hemnet.svg" alt="Hemnet">
                </picture>
            </td>
            <td><a href="https://www.hemnet.se/">hemnet.se</a></td>
        </tr>
        <tr>
            <td>Cisco</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/cisco.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/cisco.svg" alt="Cisco">
                </picture>
            </td>
            <td><a href="https://www.cisco.com/">cisco.com</a></td>
        </tr>
        <tr>
            <td>Presidencia de la República Dominicana</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/republica-dominicana.svg" />
                    <img height="42px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/republica-dominicana.svg" alt="Presidencia de la República Dominicana">
                </picture>
            </td>
            <td><a href="https://www.presidencia.gob.do/">presidencia.gob.do</a></td>
        </tr>
        <tr>
            <td>Moonpig</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/moonpig.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/moonpig.svg" alt="Moonpig">
                </picture>
            </td>
            <td><a href="https://www.moonpig.com/">moonpig.com</a></td>
        </tr>
        <tr>
            <td>Booster</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/booster.svg" />
                    <img height="18px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/booster.svg" alt="Booster">
                </picture>
            </td>
            <td><a href="https://www.choosebooster.com/">choosebooster.com</a></td>
        </tr>
        <tr>
            <td>Zaptec</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/zaptec.svg" />
                    <img height="24px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/zaptec.svg" alt="Zaptec">
                </picture>
            </td>
            <td><a href="https://www.zaptec.com/">zaptec.com</a></td>
        </tr>
        <tr>
            <td>Klarna</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/klarna.svg" />
                    <img height="24px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/klarna.svg" alt="Klarna">
                </picture>
            </td>
            <td><a href="https://www.klarna.com/">klarna.com</a></td>
        </tr>
        <tr>
            <td>Raspberry PI Foundation</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/raspi.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/raspi.svg" alt="Raspberry PI Foundation">
                </picture>
            </td>
            <td><a href="https://www.raspberrypi.org/">raspberrypi.org</a></td>
        </tr>
        <tr>
            <td>Tulip</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/tulip.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/tulip.svg" alt="Tulip Retail">
                </picture>
            </td>
            <td><a href="https://tulip.com/">tulip.com</a></td>
        </tr>
        <tr>
            <td>Hootsuite</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/hootsuite.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/hootsuite.svg" alt="Hootsuite">
                </picture>
            </td>
            <td><a href="https://hootsuite.com/">hootsuite.com</a></td>
        </tr>
        <tr>
            <td>Segment</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/segment.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/segment.svg" alt="Segment">
                </picture>
            </td>
            <td><a href="https://segment.com/">segment.com</a></td>
        </tr>
        <tr>
            <td>Arduino</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/arduino.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/arduino.svg" alt="Arduino">
                </picture>
            </td>
            <td><a href="https://www.arduino.cc/">arduino.cc</a></td>
        </tr>
        <tr>
            <td>Sainsbury's</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/sainsburys.svg" />
                    <img height="24px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/sainsburys.svg" alt="Sainsbury's">
                </picture>
            </td>
            <td><a href="https://www.sainsburys.co.uk/">sainsburys.co.uk</a></td>
        </tr>
        <tr>
            <td>Contraste</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/contraste.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/contraste.svg" alt="Contraste">
                </picture>
            </td>
            <td><a href="https://www.contraste.com/en">contraste.com</a></td>
        </tr>
        <tr>
            <td>inMusic</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/inmusic.svg" />
                    <img height="24px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/inmusic.svg" alt="InMusic">
                </picture>
            </td>
            <td><a href="https://inmusicbrands.com/">inmusicbrands.com</a></td>
        </tr>
        <tr>
            <td>Buhta</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/buhta.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/buhta.svg" alt="Buhta">
                </picture>
            </td>
            <td><a href="https://buhta.com/">buhta.com</a></td>
        </tr>
        </tr>
            <tr>
            <td>Amplitude</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/amplitude.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/amplitude.svg" alt="amplitude.com">
                </picture>
            </td>
            <td><a href="https://amplitude.com/">amplitude.com</a></td>
        </tr>
    <tr>
      <td align="center"><a href="https://tier4.jp/en/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/tieriv.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/tieriv.svg" alt="TIER IV"></picture></a></td>
      <td align="center"><a href="https://kyma-project.io"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/kyma.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/kyma.svg" alt="Kyma Project"></picture></a></td>
      <td align="center"><a href="https://serlo.org/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/serlo.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/serlo.svg" alt="Serlo"></picture></a></td>
      <td align="center"><a href="https://padis.io/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/padis.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/padis.svg" alt="Padis"></picture></a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://cloudbear.eu/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/cloudbear.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/cloudbear.svg" alt="Cloudbear"></picture></a></td>
      <td align="center"><a href="https://securityonionsolutions.com/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/securityonion.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/securityonion.svg" alt="Security Onion Solutions"></picture></a></td>
      <td align="center"><a href="https://factlylabs.com/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/factly.svg" /><img height="24px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/factly.svg" alt="Factly"></picture></a></td>
      <td align="center"><a href="https://cashdeck.com.au/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/allmyfunds.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/allmyfunds.svg" alt="All My Funds"></picture></a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://nortal.com/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/nortal.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/nortal.svg" alt="Nortal"></picture></a></td>
      <td align="center"><a href="https://www.ordermygear.com/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/ordermygear.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/ordermygear.svg" alt="OrderMyGear"></picture></a></td>
      <td align="center"><a href="https://r2devops.io/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/r2devops.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/r2devops.svg" alt="R2Devops"></picture></a></td>
      <td align="center"><a href="https://www.paralus.io/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/paralus.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/paralus.svg" alt="Paralus"></picture></a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://dyrector.io/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/dyrector_io.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/dyrector_io.svg" alt="dyrector.io"></picture></a></td>
      <td align="center"><a href="https://pinniped.dev/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/pinniped.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/pinniped.svg" alt="pinniped.dev"></picture></a></td>
      <td align="center"><a href="https://pvotal.tech/"><picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/pvotal.svg" /><img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/pvotal.svg" alt="pvotal.tech"></picture></a></td>
      <td></td>
    </tr>
    </tbody>
</table>

Many thanks to all individual contributors

<a href="https://opencollective.com/ory" target="_blank"><img src="https://opencollective.com/ory/contributors.svg?width=890&limit=714&button=false" /></a>

<!--END ADOPTERS-->

## Ecosystem

<!--BEGIN ECOSYSTEM-->

We build Ory on several guiding principles when it comes to our architecture
design:

- Minimal dependencies
- Runs everywhere
- Scales without effort
- Minimize room for human and network errors

Ory's architecture is designed to run best on a Container Orchestration system
such as Kubernetes, CloudFoundry, OpenShift, and similar projects. Binaries are
small (5-15MB) and available for all popular processor types (ARM, AMD64, i386)
and operating systems (FreeBSD, Linux, macOS, Windows) without system
dependencies (Java, Node, Ruby, libxml, ...).

### Ory Kratos: Identity and User Infrastructure and Management

[Ory Kratos](https://github.com/ory/kratos) is an API-first Identity and User
Management system that is built according to
[cloud architecture best practices](https://www.ory.sh/docs/next/ecosystem/software-architecture-philosophy).
It implements core use cases that almost every software application needs to
deal with: Self-service Login and Registration, Multi-Factor Authentication
(MFA/2FA), Account Recovery and Verification, Profile, and Account Management.

### Ory Hydra: OAuth2 & OpenID Connect Server

[Ory Hydra](https://github.com/ory/hydra) is an OpenID Certified™ OAuth2 and
OpenID Connect Provider which easily connects to any existing identity system by
writing a tiny "bridge" application. It gives absolute control over the user
interface and user experience flows.

### Ory Oathkeeper: Identity & Access Proxy

[Ory Oathkeeper](https://github.com/ory/oathkeeper) is a BeyondCorp/Zero Trust
Identity & Access Proxy (IAP) with configurable authentication, authorization,
and request mutation rules for your web services: Authenticate JWT, Access
Tokens, API Keys, mTLS; Check if the contained subject is allowed to perform the
request; Encode resulting content into custom headers (`X-User-ID`), JSON Web
Tokens and more!

### Ory Keto: Access Control Policies as a Server

[Ory Keto](https://github.com/ory/keto) is a policy decision point. It uses a
set of access control policies, similar to AWS IAM Policies, in order to
determine whether a subject (user, application, service, car, ...) is authorized
to perform a certain action on a resource.

<!--END ECOSYSTEM-->

## Documentation

The full Ory Hydra documentation is available at
[www.ory.sh/docs/hydra](https://www.ory.sh/docs/hydra), including:

- [Installation guides](https://www.ory.sh/docs/hydra/install)
- [Configuration reference](https://www.ory.sh/docs/hydra/reference/configuration)
- [HTTP API documentation](https://www.ory.sh/docs/hydra/sdk/api)
- [Security architecture](https://www.ory.sh/docs/hydra/security-architecture)
- [Performance benchmarks](https://www.ory.sh/docs/performance/hydra)

For upgrading and changelogs, check
[releases tab](https://github.com/ory/hydra/releases) and
[CHANGELOG.md](./CHANGELOG.md).

## Developing Ory Hydra

See [DEVELOP.md](./DEVELOP.md) for information on:

- Contribution guidelines
- Prerequisites
- Install from source
- Running tests
- Build Docker image
- Preview API documentation

## Security

OAuth2 and OAuth2 related specifications are over 400 written pages.
Implementing OAuth2 is easy, getting it right is hard. Ory Hydra is trusted by
companies all around the world, has a vibrant community and faces millions of
requests in production each day. Read
[the security guide](https://www.ory.sh/docs/hydra/security-architecture) for
more details on cryptography and security concepts.

### Disclosing vulnerabilities

If you think you found a security vulnerability, please refrain from posting it
publicly on the forums, the chat, or GitHub. You can find all info for
responsible disclosure in our
[security.txt](https://www.ory.sh/.well-known/security.txt).

## Telemetry

Our services collect summarized, anonymized data that can optionally be turned
off. Click [here](https://www.ory.sh/docs/ecosystem/sqa) to learn more.

## Libraries and third-party projects

Official:

- [User Login & Consent Example](https://github.com/ory/hydra-login-consent-node)

Community:

- Visit
  [this document for an overview of community projects and articles](https://www.ory.sh/docs/ecosystem/community)

Developer Blog:

- Visit the [Ory Blog](https://www.ory.sh/blog/) for guides, tutorials and
  articles around Ory Hydra and the Ory ecosystem.
