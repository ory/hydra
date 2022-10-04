<h1 align="center"><img src="https://raw.githubusercontent.com/ory/meta/master/static/banners/hydra.svg" alt="Ory Hydra - Open Source OAuth 2 and OpenID Connect server"></h1>

<h4 align="center">
    <a href="https://www.ory.sh/chat">Chat</a> |
    <a href="https://github.com/ory/hydra/discussions">Discussions</a> |
    <a href="http://eepurl.com/di390P">Newsletter</a><br/><br/>
    <a href="https://www.ory.sh/hydra/docs/index">Guide</a> |
    <a href="https://www.ory.sh/hydra/docs/reference/api">API Docs</a> |
    <a href="https://godoc.org/github.com/ory/hydra">Code Docs</a><br/><br/>
    <a href="https://opencollective.com/ory">Support this project!</a><br/><br/>
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
and [SDKs](https://www.ory.sh/docs/hydra/sdk/) for common languages are
provided.

## Ory Cloud

The easiest way to get started with Ory Software is in Ory Cloud! It is
[**free for developers**](https://console.ory.sh/registration?utm_source=github&utm_medium=banner&utm_campaign=hydra-readme),
forever, no credit card required!

Ory Cloud has easy examples, administrative user interfaces, hosted pages (e.g.
for login or registration), support for custom domains, collaborative features
for your colleagues, and much more!

## Get Started

If you're looking to jump straight into it, go ahead:

- **[Run your own OAuth 2.0 Server - step by step guide](https://www.ory.sh/run-oauth2-server-open-source-api-security/)**:
  A in-depth look at setting up Ory Hydra and performing a variety of OAuth 2.0
  Flows.
- [Ory Hydra 5 Minute Tutorial](https://www.ory.sh/docs/hydra/5min-tutorial):
  Set up and use Ory Hydra using Docker Compose in under 5 Minutes. Good for
  hacking a Proof of Concept.
- [Run Ory Hydra in Docker](https://www.ory.sh/docs/hydra/configure-deploy): An
  advanced guide to a fully functional set up with Ory Hydra.
- [Integrating your Login and Consent UI with Ory Hydra](https://www.ory.sh/docs/hydra/oauth2):
  The go-to place if you wish to adopt Ory Hydra in your new or existing stack.

Besides mitigating various attack vectors, such as a compromised database and
OAuth 2.0 weaknesses, Ory Hydra is also able to securely manage JSON Web Keys.
[Click here](https://www.ory.sh/docs/hydra/security-architecture) to read more
about security.

---

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents**

- [What is Ory Hydra?](#what-is-ory-hydra)
  - [Who's using it?](#whos-using-it)
  - [OAuth2 and OpenID Connect: Open Standards!](#oauth2-and-openid-connect-open-standards)
  - [OpenID Connect Certified](#openid-connect-certified)
- [Quickstart](#quickstart)
  - [5 minutes tutorial: Run your very own OAuth2 environment](#5-minutes-tutorial-run-your-very-own-oauth2-environment)
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
    - [Build Docker](#build-docker)
    - [Run the Docker Compose quickstarts](#run-the-docker-compose-quickstarts)
- [Libraries and third-party projects](#libraries-and-third-party-projects)
- [Blog posts & articles](#blog-posts--articles)

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
maintainers. We thank everyone involved - from submitting bug reports and
feature requests, to contributing patches, to sponsoring our work. Our community
is 1000+ strong and growing rapidly. The Ory stack protects 16.000.000.000+ API
requests every month with over 250.000+ active service nodes. We would have
never been able to achieve this without each and everyone of you!

The following list represents companies that have accompanied us along the way
and that have made outstanding contributions to our ecosystem. _If you think
that your company deserves a spot here, reach out to
<a href="mailto:office-muc@ory.sh">office-muc@ory.sh</a> now_!

**Please consider giving back by becoming a sponsor of our open source work on
<a href="https://www.patreon.com/_ory">Patreon</a> or
<a href="https://opencollective.com/ory">Open Collective</a>.**

<table>
    <thead>
        <tr>
            <th>Type</th>
            <th>Name</th>
            <th>Logo</th>
            <th>Website</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>Sponsor</td>
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
            <td>Contributor</td>
            <td>Kyma Project</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/kyma.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/kyma.svg" alt="Kyma Project">
                </picture>
            </td>
            <td><a href="https://kyma-project.io">kyma-project.io</a></td>
        </tr>
        <tr>
            <td>Sponsor</td>
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
            <td>Sponsor</td>
            <td>Cashdeck / All My Funds</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/allmyfunds.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/allmyfunds.svg" alt="All My Funds">
                </picture>
            </td>
            <td><a href="https://cashdeck.com.au/">cashdeck.com.au</a></td>
        </tr>
        <tr>
            <td>Contributor</td>
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
            <td>Adopter *</td>
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
            <td>Adopter *</td>
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
            <td>Adopter *</td>
            <td>DataDetect</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/datadetect.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/datadetect.svg" alt="Datadetect">
                </picture>
            </td>
            <td><a href="https://unifiedglobalarchiving.com/data-detect/">unifiedglobalarchiving.com/data-detect/</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Sainsbury's</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/sainsburys.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/sainsburys.svg" alt="Sainsbury's">
                </picture>
            </td>
            <td><a href="https://www.sainsburys.co.uk/">sainsburys.co.uk</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
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
            <td>Adopter *</td>
            <td>Reyah</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/reyah.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/reyah.svg" alt="Reyah">
                </picture>
            </td>
            <td><a href="https://reyah.eu/">reyah.eu</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Zero</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/commitzero.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/commitzero.svg" alt="Project Zero by Commit">
                </picture>
            </td>
            <td><a href="https://getzero.dev/">getzero.dev</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Padis</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/padis.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/padis.svg" alt="Padis">
                </picture>
            </td>
            <td><a href="https://padis.io/">padis.io</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Cloudbear</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/cloudbear.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/cloudbear.svg" alt="Cloudbear">
                </picture>
            </td>
            <td><a href="https://cloudbear.eu/">cloudbear.eu</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Security Onion Solutions</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/securityonion.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/securityonion.svg" alt="Security Onion Solutions">
                </picture>
            </td>
            <td><a href="https://securityonionsolutions.com/">securityonionsolutions.com</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Factly</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/factly.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/factly.svg" alt="Factly">
                </picture>
            </td>
            <td><a href="https://factlylabs.com/">factlylabs.com</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Nortal</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/nortal.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/nortal.svg" alt="Nortal">
                </picture>
            </td>
            <td><a href="https://nortal.com/">nortal.com</a></td>
        </tr>
        <tr>
            <td>Sponsor</td>
            <td>OrderMyGear</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/ordermygear.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/ordermygear.svg" alt="OrderMyGear">
                </picture>
            </td>
            <td><a href="https://www.ordermygear.com/">ordermygear.com</a></td>
        </tr>
        <tr>
            <td>Sponsor</td>
            <td>Spiri.bo</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/spiribo.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/spiribo.svg" alt="Spiri.bo">
                </picture>
            </td>
            <td><a href="https://spiri.bo/">spiri.bo</a></td>
        </tr>
        <tr>
            <td>Sponsor</td>
            <td>Strivacity</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/strivacity.svg" />
                    <img height="16px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/strivacity.svg" alt="Spiri.bo">
                </picture>
            </td>
            <td><a href="https://strivacity.com/">strivacity.com</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Hanko</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/hanko.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/hanko.svg" alt="Hanko">
                </picture>
            </td>
            <td><a href="https://hanko.io/">hanko.io</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Rabbit</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/rabbit.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/rabbit.svg" alt="Rabbit">
                </picture>
            </td>
            <td><a href="https://rabbit.co.th/">rabbit.co.th</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>inMusic</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/inmusic.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/inmusic.svg" alt="InMusic">
                </picture>
            </td>
            <td><a href="https://inmusicbrands.com/">inmusicbrands.com</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Buhta</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/buhta.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/buhta.svg" alt="Buhta">
                </picture>
            </td>
            <td><a href="https://buhta.com/">buhta.com</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Connctd</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/connctd.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/connctd.svg" alt="Connctd">
                </picture>
            </td>
            <td><a href="https://connctd.com/">connctd.com</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>Paralus</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/paralus.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/paralus.svg" alt="Paralus">
                </picture>
            </td>
            <td><a href="https://www.paralus.io/">paralus.io</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>TIER IV</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/tieriv.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/tieriv.svg" alt="TIER IV">
                </picture>
            </td>
            <td><a href="https://tier4.jp/en/">tier4.jp</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>R2Devops</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/r2devops.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/r2devops.svg" alt="R2Devops">
                </picture>
            </td>
            <td><a href="https://r2devops.io/">r2devops.io</a></td>
        </tr>
        <tr>
            <td>Adopter *</td>
            <td>LunaSec</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/lunasec.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/lunasec.svg" alt="LunaSec">
                </picture>
            </td>
            <td><a href="https://www.lunasec.io/">lunasec.io</a></td>
        </tr>
                <tr>
            <td>Adopter *</td>
            <td>Serlo</td>
            <td align="center">
                <picture>
                    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/ory/meta/master/static/adopters/light/serlo.svg" />
                    <img height="32px" src="https://raw.githubusercontent.com/ory/meta/master/static/adopters/dark/serlo.svg" alt="Serlo">
                </picture>
            </td>
            <td><a href="https://serlo.org/">serlo.org</a></td>
        </tr>
    </tbody>
</table>

We also want to thank all individual contributors

<a href="https://opencollective.com/ory" target="_blank"><img src="https://opencollective.com/ory/contributors.svg?width=890&limit=714&button=false" /></a>

as well as all of our backers

<a href="https://opencollective.com/ory#backers" target="_blank"><img src="https://opencollective.com/ory/backers.svg?width=890"></a>

and past & current supporters (in alphabetical order) on
[Patreon](https://www.patreon.com/_ory): Alexander Alimovs, Billy, Chancy
Kennedy, Drozzy, Edwin Trejos, Howard Edidin, Ken Adler Oz Haven, Stefan Hans,
TheCrealm.

<em>\* Uses one of Ory's major projects in production.</em>

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

### 5 minutes tutorial: Host your own OAuth2 environment

The **[tutorial](https://www.ory.sh/docs/hydra/5min-tutorial)** teaches you to
set up Ory Hydra, a Postgres instance and an exemplary identity provider written
in React using docker-compose. It will take you about 5 minutes to complete the
**[tutorial](https://www.ory.sh/docs/hydra/5min-tutorial)**.

<img src=".github/readme/oauth2-flow.gif" alt="OAuth2 Flow">

<br clear="all">

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
publicly on the forums, the chat, or GitHub and send us an email to
[hi@ory.am](mailto:hi@ory.sh) instead.

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
test-resetdb

# end-to-end tests
make e2e
```

##### Short Tests

It is recommended to use the make file to run your tests using `make quicktest`
, however, you can still use the `go test` command.

**Please note**:

All tests run against a sqlite in-memory database, thus it is required to use
the `-tags sqlite,json1` build tag.

Short tests run fairly quickly. You can either test all of the code at once:

```shell script
go test -v -failfast -short -tags sqlite,json1 ./...
```

or test just a specific module:

```shell script
go test -v -failfast -short -tags sqlite,json1 ./client
```

or a specific test:

```shell script
go test -v -failfast -short -tags sqlite,json1 -run ^TestName$ ./...
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
make docker
docker compose -f quickstart.yml up # ....
```

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
- Visit the [Ory Blog](https://www.ory.sh/blog/) for guides, tutorials and
  articles around Ory Hydra and the Ory ecosystem.
