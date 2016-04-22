# Ory/Hydra

[![Join the chat at https://gitter.im/ory-am/hydra](https://badges.gitter.im/ory-am/hydra.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)

![Hydra](hydra.png)

![Hydra implements HTTP/2 and TLS.](h2tls.png)

Hydra is being developed at [Ory](https://ory.am) because we need a lightweight and clean IAM solution for our customers.  
Join our [mailinglist](http://eepurl.com/bKT3N9) to stay on top of new developments.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

Authentication, authorization and user account management are always lengthy to plan and implement. If you're building a micro service app
in need of these three, you are in the right place.

## Motivation

We develop Hydra because Hydra we need a lightweight and clean IAM solution for our customers. We believe that security and simplicity come together. This is why Hydra only relies on Google's Go Language, PostgreSQL or RethinkDB and a slim dependency tree. Hydra is the simple, open source alternative to proprietary authorization solutions suited best for your micro service eco system.

*Use it, enjoy it and contribute!*

## Features

## Run hydra-host

### With vagrant

You'll need [Vagrant](https://www.vagrantup.com/), [VirtualBox](https://www.virtualbox.org/) and [Git](https://git-scm.com/)
installed on your system.

```
git clone https://github.com/ory-am/hydra.git
cd hydra
vagrant up
# Get a cup of coffee
```

You should now have a running Hydra instance! Vagrant exposes ports 9000 (HTTPS - Hydra) and 9001 (Postgres) on your localhost.
Open [https://localhost:9000/](https://localhost:9000/) to confirm that Hydra is running. You will probably have to add an exception for the
HTTP certificate because it is self-signed, but after that you should see a 404 error indicating that Hydra is running!

*hydra-host* offers different capabilities for managing your Hydra instance.
Check the [this section](#cli-usage) if you want to find out more.

You can also always access hydra-host through vagrant:

```
# Assuming, that your current working directory is /where/you/cloned/hydra
vagrant ssh
hydra-host help
```

## Security considerations

[rfc6819](https://tools.ietf.org/html/rfc6819) provides good guidelines to keep your apps and environment secure. It is recommended to read:
* [Section 5.3](https://tools.ietf.org/html/rfc6819#section-5.3) on client app security.

## Good to know

This section covers information necessary for understanding how hydra works.

### Deploy with buildpacks (Heroku, Cloud Foundry, ...)

Hydra runs pretty much out of the box when using a Platform as a Service (PaaS).
Here are however a few notes which might assist you in your task:
* Heroku (and probably Cloud Foundry as well) *force* TLS termination, meaning that Hydra must be configured with `DANGEROUSLY_FORCE_HTTP=force`.
* Using bash, you can easily add multi-line environment variables to Heroku using `heroku config:set JWT_PUBLIC_KEY="$(my-public-key.pem)"`.
  This does not work on Windows!



## Attributions

* [Logo source](https://www.flickr.com/photos/pathfinderlinden/7161293044/)
