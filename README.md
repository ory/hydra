<h1 align="center"><img src="https://raw.githubusercontent.com/ory/meta/master/static/banners/hydra.svg" alt="Ory Hydra - Open Source OAuth 2 and OpenID Connect server"></h1>

<h4 align="center">
    <a href="https://www.ory.sh/chat">Chat</a> |
    <a href="https://github.com/ory/hydra/discussions">Discussions</a> |
    <a href="https://www.ory.sh/l/sign-up-newsletter">Newsletter</a><br/><br/>
    <a href="https://www.ory.sh/hydra/docs/index">Guide</a> |
    <a href="https://www.ory.sh/hydra/docs/reference/api">API Docs</a> |
    <a href="https://godoc.org/github.com/ory/hydra">Code Docs</a><br/><br/>
    <a href="https://console.ory.sh/">Support this project!</a><br/><br/>
    <a href="https://www.ory.sh/jobs/">Work in Open Source, Ory is hiring!</a>
</h4>

---

<p align="left">
    <a href="https://github.com/ory/hydra/actions/workflows/ci.yaml"><img src="https://github.com/ory/hydra/actions/workflows/ci.yaml/badge.svg?branch=master&event=push" alt="CI Tasks for Ory Hydra"></a>
    <a href="https://codecov.io/gh/ory/hydra"><img src="https://codecov.io/gh/ory/hydra/branch/master/graph/badge.svg?token=y4fVk2Of8a"/></a>
    <a href="https://goreportcard.com/report/github.com/ory/hydra"><img src="https://goreportcard.com/badge/github.com/ory/hydra" alt="Go Report Card"></a>
    <a href="https://pkg.go.dev/github.com/ory/hydra"><img src="https://pkg.go.dev/badge/www.github.com/ory/hydra" alt="PkgGoDev"></a>
    <a href="https://bestpractices.coreinfrastructure.org/projects/364"><img src="https://bestpractices.coreinfrastructure.org/projects/364/badge" alt="CII Best Practices"></a>
    <a href="#backers" alt="sponsors on Open Collective"><img src="https://opencollective.com/ory/backers/badge.svg" /></a> <a href="#sponsors" alt="Sponsors on Open Collective"><img src="https://opencollective.com/ory/sponsors/badge.svg" /></a>
    <a href="https://github.com/ory/hydra/blob/master/CODE_OF_CONDUCT.md" alt="Ory Code of Conduct"><img src="https://img.shields.io/badge/ory-code%20of%20conduct-green" /></a>
</p>

Ory Hydra is a hardened, **OpenID Certified OAuth 2.0 Server and OpenID Connect
Provider** optimized for low-latency, high throughput, and low resource
consumption. Ory Hydra _is not_ an identity provider (user sign up, user login,
password reset flow), but connects to your existing identity provider through a
[login and consent app](https://www.ory.sh/docs/hydra/oauth2#authenticating-users-and-requesting-consent).
Implementing the login and consent app in a different language is easy, and
exemplary consent apps ([Node](https://github.com/ory/hydra-login-consent-node))
and [SDKs](https://www.ory.sh/docs/kratos/sdk/index) for common languages are
provided.

Ory Hydra can use [Ory Kratos](https://github.com/ory/kratos) as its identity
server.

## Ory Hydra on the Ory Network

The [Ory Network](https://www.ory.sh/cloud) is the fastest, most secure and
worry-free way to use Ory's Services. **Ory OAuth2 & OpenID Connect** is powered
by the Ory Hydra open source federation server, and it's fully API-compatible.

The Ory Network provides the infrastructure for modern end-to-end security:

- Identity & credential management scaling to billions of users and devices
- Registration, Login and Account management flows for passkey, biometric,
  social, SSO and multi-factor authentication
- **Pre-built login, registration and account management pages and components**
- **OAuth2 and OpenID provider for single sign on, API access and
  machine-to-machine authorization**
- Low-latency permission checks based on Google's Zanzibar model and with
  built-in support for the Ory Permission Language

It's fully managed, highly available, developer & compliance-friendly!

- GDPR-friendly secure storage with data locality
- Cloud-native APIs, compatible with Ory's Open Source servers
- Comprehensive admin tools with the web-based Ory Console and the Ory Command
  Line Interface (CLI)
- Extensive documentation, straightforward examples and easy-to-follow guides
- Fair, usage-based [pricing](https://www.ory.sh/pricing)

Sign up for a
[**free developer account**](https://console.ory.sh/registration?utm_source=github&utm_medium=banner&utm_campaign=kratos-readme)
today!

## Ory Hydra On-premise support

Are you running Ory Hydra in a mission-critical, commercial environment? The Ory
Enterprise License (OEL) provides enhanced features, security, and expert
support directly from the Ory core maintainers.

Organizations that require advanced features, enhanced security, and
enterprise-grade support for Ory's identity and access management solutions
benefit from the Ory Enterprise License (OEL) as a self-hosted, premium offering
including:

- Additional features not available in the open-source version.
- Regular releases that address CVEs and security vulnerabilities, with strict
  SLAs for patching based on severity.
- Support for advanced scaling and multi-tenancy features.
- Premium support options, including SLAs, direct engineer access, and concierge
  onboarding.
- Access to private Docker registry for a faster, more reliable access to vetted
  enterprise builds.

A valid Ory Enterprise License and access to the Ory Enterprise Docker Registry
are required to use these features. OEL is designed for mission-critical,
production, and global applications where organizations need maximum control and
flexibility over their identity infrastructure. Ory's offering is the only
official program for qualified support from the maintainers. For more
information book a meeting with the Ory team to
**[discuss your needs](https://www.ory.sh/contact/)**!

## Get Started

You can use
[Docker to run Ory Hydra locally](https://www.ory.sh/docs/hydra/5min-tutorial)
or use the Ory CLI to try out Ory Hydra:

```shell
# This example works best in Bash
bash <(curl https://raw.githubusercontent.com/ory/meta/master/install.sh) -b . ory
sudo mv ./ory /usr/local/bin/

# Or with Homebrew installed
brew install ory/tap/cli
```

create a new project (you may also use
[Docker](https://www.ory.sh/docs/hydra/5min-tutorial))

```
ory create project --name "Ory Hydra 2.0 Example"
project_id="{set to the id from output}"
```

and follow the quick & easy steps below.

### OAuth 2.0 Client Credentials / Machine-to-Machine

Create an OAuth 2.0 Client, and run the OAuth 2.0 Client Credentials flow:

```shell
ory create oauth2-client --project $project_id \
    --name "Client Credentials Demo" \
    --grant-type client_credentials
client_id="{set to client id from output}"
client_secret="{set to client secret from output}"

ory perform client-credentials --client-id=$client_id --client-secret=$client_secret --project $project_id
access_token="{set to access token from output}"

ory introspect token $access_token --project $project_id
```

### OAuth 2.0 Authorize Code + OpenID Connect

Try out the OAuth 2.0 Authorize Code grant right away!

By accepting permissions `openid` and `offline_access` at the consent screen,
Ory refreshes and OpenID Connect ID token,

```shell
ory create oauth2-client --project $project_id \
    --name "Authorize Code with OpenID Connect Demo" \
    --grant-type authorization_code,refresh_token \
    --response-type code \
    --redirect-uri http://127.0.0.1:4446/callback
code_client_id="{set to client id from output}"
code_client_secret="{set to client secret from output}"

ory perform authorization-code \
    --project $project_id \
    --client-id $code_client_id \
    --client-secret $code_client_secret
code_access_token="{set to access token from output}"

ory introspect token $code_access_token --project $project_id
```

---

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [What is Ory Hydra?](#what-is-ory-hydra)
  - [Who's using it?](#whos-using-it)
  - [OAuth2 and OpenID Connect: Open Standards!](#oauth2-and-openid-connect-open-standards)
  - [OpenID Connect Certified](#openid-connect-certified)
- [Quickstart](#quickstart)
  - [Installation](#installation)
- [Ecosystem](#ecosystem)
  - [Ory Kratos: Identity and User Infrastructure and Management](#ory-kratos-identity-and-user-infrastructure-and-management)
  - [Ory Hydra: OAuth2 & OpenID Connect Server](#ory-hydra-oauth2--openid-connect-server)
  - [Ory Oathkeeper: Identity & Access Proxy](#ory-oathkeeper-identity--access-proxy)
  - [Ory Keto: Access Control Policies as a Server](#ory-keto-access-control-policies-as-a-server)
- [Security](#security)
  - [Disclosing vulnerabilities](#disclosing-vulnerabilities)
- [Benchmarks](#benchmarks)
- [Telemetry](#telemetry)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [HTTP API documentation](#http-api-documentation)
  - [Upgrading and Changelog](#upgrading-and-changelog)
  - [Command line documentation](#command-line-documentation)
  - [Develop](#develop)
    - [Dependencies](#dependencies)
    - [Formatting Code](#formatting-code)
    - [Running Tests](#running-tests)
      - [Short Tests](#short-tests)
      - [Regular Tests](#regular-tests)
    - [E2E Tests](#e2e-tests)
      - [OpenID Connect Conformity Tests](#openid-connect-conformity-tests)
    - [Build Docker](#build-docker)
    - [Run the Docker Compose quickstarts](#run-the-docker-compose-quickstarts)
    - [Add a new migration](#add-a-new-migration)
- [Libraries and third-party projects](#libraries-and-third-party-projects)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Ory Hydra?

Ory Hydra is a server implementation of the OAuth 2.0 authorization framework
and the OpenID Connect Core 1.0. Existing OAuth2 implementations usually ship as
libraries or SDKs such as
[node-oauth2-server](https://github.com/oauthjs/node-oauth2-server) or
[Ory Fosite](https://github.com/ory/fosite/issues), or as fully featured
identity solutions with user management and user interfaces, such as
[Keycloak](https://www.keycloak.org).

Implementing and using OAuth2 without understanding the whole specification is
challenging and prone to errors, even when SDKs are being used. The primary goal
of Ory Hydra is to make OAuth 2.0 and OpenID Connect 1.0 better accessible.

Ory Hydra implements the flows described in OAuth2 and OpenID Connect 1.0
without forcing you to use a "Hydra User Management" or some template engine or
a predefined front-end. Instead, it relies on HTTP redirection and cryptographic
methods to verify user consent allowing you to use Ory Hydra with any
authentication endpoint, be it [Ory Kratos](https://github.com/ory/kratos),
[authboss](https://github.com/go-authboss/authboss),
[User Frosting](https://www.userfrosting.com/) or your proprietary Java
authentication.

### Who's using it?

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

### OAuth2 and OpenID Connect: Open Standards!

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

## Quickstart

This section is a starter guide to working with Ory Hydra. In-depth docs are
available as well:

- The documentation is available [here](https://www.ory.sh/docs/hydra).
- The REST API documentation is available
  [here](https://www.ory.sh/docs/hydra/sdk/api).

### Installation

Head over to the
[Ory Developer Documentation](https://www.ory.sh/docs/hydra/install) to learn
how to install Ory Hydra on Linux, macOS, Windows, and Docker and how to build
Ory Hydra from source.

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

## Security

_Why should I use Ory Hydra? It's not that hard to implement two OAuth2
endpoints and there are numerous SDKs out there!_

OAuth2 and OAuth2 related specifications are over 400 written pages.
Implementing OAuth2 is easy, getting it right is hard. Ory Hydra is trusted by
companies all around the world, has a vibrant community and faces millions of
requests in production each day. Of course, we also compiled a security guide
with more details on cryptography and security concepts. Read
[the security guide now](https://www.ory.sh/docs/hydra/security-architecture).

### Disclosing vulnerabilities

If you think you found a security vulnerability, please refrain from posting it
publicly on the forums, the chat, or GitHub. You can find all info for
responsible disclosure in our
[security.txt](https://www.ory.sh/.well-known/security.txt).

## Benchmarks

Our continuous integration runs a collection of benchmarks against Ory Hydra.
You can find the results [here](https://www.ory.sh/docs/performance/hydra).

## Telemetry

Our services collect summarized, anonymized data that can optionally be turned
off. Click [here](https://www.ory.sh/docs/ecosystem/sqa) to learn more.

## Documentation

### Guide

The full Ory Hydra documentation is available
[here](https://www.ory.sh/docs/hydra).

### HTTP API documentation

The HTTP API is documented [here](https://www.ory.sh/docs/hydra/sdk/api).

### Upgrading and Changelog

New releases might introduce breaking changes. To help you identify and
incorporate those changes, we document these changes in
[CHANGELOG.md](./CHANGELOG.md).

### Command line documentation

Run `hydra -h` or `hydra help`.

### Develop

We love all contributions! Please read our
[contribution guidelines](./CONTRIBUTING.md).

#### Dependencies

You need Go 1.13+ with `GO111MODULE=on` and (for the test suites):

- Docker and Docker Compose
- Makefile
- NodeJS / npm

It is possible to develop Ory Hydra on Windows, but please be aware that all
guides assume a Unix shell like bash or zsh.

#### Formatting Code

You can format all code using `make format`. Our CI checks if your code is
properly formatted.

#### Running Tests

There are three types of tests you can run:

- Short tests (do not require a SQL database like PostgreSQL)
- Regular tests (do require PostgreSQL, MySQL, CockroachDB)
- End to end tests (do require databases and will use a test browser)

All of the above tests can be run using the makefile. See the commands below.

**Makefile commands**

```shell
# quick tests
make quicktest

# regular tests
make test

# updates all snapshots
make test-refresh

# end-to-end tests
make e2e
```

##### Short Tests

It is recommended to use the make file to run your tests using `make quicktest`
, however, you can still use the `go test` command.

**Please note**:

All tests run against a sqlite in-memory database, thus it is required to use
the `-tags sqlite` build tag.

Short tests run fairly quickly. You can either test all of the code at once:

```shell script
go test -v -failfast -short -tags sqlite ./...
```

or test just a specific module:

```shell script
go test -v -failfast -short -tags sqlite ./client
```

or a specific test:

```shell script
go test -v -failfast -short -tags sqlite -run ^TestName$ ./...
```

##### Regular Tests

Regular tests require a database set up. Our test suite is able to work with
docker directly (using [ory/dockertest](https://github.com/ory/dockertest)) but
we encourage to use the Makefile instead. Using dockertest can bloat the number
of Docker Images on your system and are quite slow. Instead we recommend doing:

```shell script
make test
```

Please be aware that `make test` recreates the databases every time you run
`make test`. This can be annoying if you are trying to fix something very
specific and need the database tests all the time. In that case we suggest that
you initialize the databases with:

```shell script
make test-resetdb
export TEST_DATABASE_MYSQL='mysql://root:secret@(127.0.0.1:3444)/mysql?parseTime=true&multiStatements=true'
export TEST_DATABASE_POSTGRESQL='postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable'
export TEST_DATABASE_COCKROACHDB='cockroach://root@127.0.0.1:3446/defaultdb?sslmode=disable'
```

Then you can run `go test` as often as you'd like:

```shell script
go test -p 1 ./...

# or in a module:
cd client; go test .
```

#### E2E Tests

The E2E tests use [Cypress](https://www.cypress.io) to run full browser tests.
You can execute these tests with:

```
make e2e
```

The runner will not show the Browser window, as it runs in the CI Mode
(background). That makes debugging these type of tests very difficult, but
thankfully you can run the e2e test in the browser which helps with debugging!
Just run:

```shell script
./test/e2e/circle-ci.bash memory --watch

# Or for the JSON Web Token Access Token strategy:
# ./test/e2e/circle-ci.bash memory-jwt --watch
```

or if you would like to test one of the databases:

```shell script
make test-resetdb
export TEST_DATABASE_MYSQL='mysql://root:secret@(127.0.0.1:3444)/mysql?parseTime=true&multiStatements=true'
export TEST_DATABASE_POSTGRESQL='postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable'
export TEST_DATABASE_COCKROACHDB='cockroach://root@127.0.0.1:3446/defaultdb?sslmode=disable'

# You can test against each individual database:
./test/e2e/circle-ci.bash postgres --watch
./test/e2e/circle-ci.bash memory --watch
./test/e2e/circle-ci.bash mysql --watch
# ...
```

Once you run the script, a Cypress window will appear. Hit the button "Run all
Specs"!

The code for these tests is located in
[./cypress/integration](./cypress/integration) and
[./cypress/support](./cypress/support) and
[./cypress/helpers](./cypress/helpers). The website you're seeing is located in
[./test/e2e/oauth2-client](./test/e2e/oauth2-client).

##### OpenID Connect Conformity Tests

To run Ory Hydra against the OpenID Connect conformity suite, run

```shell script
$ test/conformity/start.sh --build
```

and then in a separate shell

```shell script
$ test/conformity/test.sh
```

Running these tests will take a significant amount of time which is why they are
not part of the CI pipeline.

#### Build Docker

You can build a development Docker Image using:

```shell script
make docker
```

#### Run the Docker Compose quickstarts

If you wish to check your code changes against any of the docker-compose
quickstart files, run:

```shell script
docker compose -f quickstart.yml up --build # ....
```

> [!WARNING] If you already have a production image (e.g. `oryd/hydra:v2.2.0`)
> pulled, the above `make docker` command will replace it with a local build of
> the image that is more equivalent to the `-distroless` variant on Docker Hub.
>
> You can pull the production image any time using `docker pull`

#### Add a new migration

1. `mkdir persistence/sql/src/YYYYMMDD000001_migration_name/`
2. Put the migration files into this directory, following the standard naming
   conventions. If you wish to execute different parts of a migration in
   separate transactions, add split marks (lines with the text `--split`) where
   desired. Why this might be necessary is explained in
   https://github.com/gobuffalo/fizz/issues/104.
3. Run `make persistence/sql/migrations/<migration_id>` to generate migration
   fragments.
4. If an update causes the migration to have fewer fragments than the number
   already generated, run
   `make persistence/sql/migrations/<migration_id>-clean`. This is equivalent to
   a `rm` command with the right parameters, but comes with better tab
   completion.
5. Before committing generated migration fragments, run the above clean command
   and generate a fresh copy of migration fragments to make sure the `sql/src`
   and `sql/migrations` directories are consistent.

## Libraries and third-party projects

Official:

- [User Login & Consent Example](https://github.com/ory/hydra-login-consent-node)

Community:

- Visit
  [this document for an overview of community projects and articles](https://www.ory.sh/docs/ecosystem/community)

Developer Blog:

- Visit the [Ory Blog](https://www.ory.sh/blog/) for guides, tutorials and
  articles around Ory Hydra and the Ory ecosystem.
