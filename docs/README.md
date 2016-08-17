# What is [Hydra](https://github.com/ory-am/hydra)?

At first, there was the monolith. The monolith worked well with the bespoke authentication module.
Then, the web evolved into an elastic cloud that serves thousands of different user agents
in every part of the world.

Hydra is driven by the need for a **scalable in memory
OAuth2 and OpenID Connect** layer, that integrates with every Identity Provider you can imagine.

Hydra is available through [Docker](https://hub.docker.com/r/oryam/hydra/) and at [GitHub](https://github.com/ory-am/hydra).

### Feature Overview

1. **Availability:** Hydra uses pub/sub to have the latest data available in memory. The in-memory architecture allows for heavy duty workloads.
2. **Scalability:** Hydra scales effortlessly on every platform you can imagine, including Heroku, Cloud Foundry, Docker,
Google Container Engine and many more.
3. **Integration:** Hydra wraps your existing stack like a blanket and keeps it safe. Hydra uses cryptographic tokens for authenticate users and request their consent, no APIs required.
The deprecated php-3.0 authentication service your intern wrote? It works with that too, don't worry.
We wrote an example with React to show you how this could look like: [React.js Identity Provider Example App](https://github.com/ory-am/hydra-idp-react).
4. **Security:** Hydra leverages the security first OAuth2 framework **[Fosite](https://github.com/ory-am/fosite)**,
encrypts important data at rest, and supports HTTP over TLS (https) out of the box.
5. **Ease of use:** Developers and Operators are human. Therefore, Hydra is easy to install and manage. Hydra does not care if you use React, Angular, or Cocoa for your user interface.
To support you even further, there are APIs available for *cryptographic key management, social log on, policy based access control, policy management, and two factor authentication (tbd)*
Hydra is packaged using [Docker](https://hub.docker.com/r/oryam/hydra/).
6. **Open Source:** Hydra is licensed Apache Version 2.0
7. **Professional:** Hydra implements peer reviewed open standards published by [The Internet Engineering Task Force (IETF®)](https://www.ietf.org/) and the [OpenID Foundation](https://openid.net/)
and under supervision of the [LMU Teaching and Research Unit Programming and Modelling Languages](http://www.en.pms.ifi.lmu.de). No funny business.
8. **Real Time:** Operation is a lot easier with real time monitoring. Because Hydra leverages RethinkDB, you get real time monitoring for free.

