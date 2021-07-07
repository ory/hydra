---
id: hsm-support
title: Hardware Security Module support for JSON Web Key Sets
---

The PKCS#11 Cryptographic Token Interface Standard, also known as Cryptoki, is
one of the Public Key Cryptography Standards developed by RSA Security. PKCS#11
defines the interface between an application and a cryptographic device.

PKCS#11 is used as a low-level interface to perform cryptographic operations
without the need for the application to directly interface a device through its
driver. PKCS#11 represents cryptographic devices using a common model referred
to simply as a token. An application can therefore perform cryptographic
operations on any device or token, using the same independent command set.

### HSM configuration

```
HSM_ENABLED=true
HSM_LIBRARY=/path/to/hsm-vendor/library.so
HSM_TOKEN_LABEL=hydra
HSM_SLOT=0
HSM_PIN=1234
```

Token that is denoted by environment variables `HSM_TOKEN_LABEL` or `HSM_SLOT`
must preexist and optionally contain RSA or ECDSA key pairs with labels
`hydra.openid.id-token` and `hydra.jwt.access-token` depending on Hydra
configuration. If keys don't exist, they will be generated upon startup. If both
`HSM_TOKEN_LABEL` and `HSM_SLOT` are set, `HSM_TOKEN_LABEL` takes preference
over `HSM_SLOT`. In this case first slot that contains this label is used.

When Hydra uses key pair from HSM, the ID of HSM object (i.e. `--id` value when
generating key pair with `pkcs11-tool`) is used as `kid` in JSON Web Key Set.
Furthermore, ID's of key pair private/public handles must be identical. Key
usage (i.e. `--usage-sign` or `--usage-decrypt` value when generating key pair
with `pkcs11-tool`) is mapped to either `sig` or `enc` and set as `use` in JSON
Web Key Set.

### Testing with SoftHSM

Change into the directory with the Hydra source code and run the following
command to start the needed containers with SoftHSM support:

```shell
$ docker-compose -f quickstart-hsm.yml up --build
```

On start up, ORY Hydra should inform if HSM is configured. Let's take a look at
the logs:

```shell
$ docker logs ory-hydra-example--hydra
time="2021-07-07T12:51:23Z" level=info msg="Hardware Security Module is configured."
time="2021-07-07T12:51:23Z" level=info msg="Using key pair 'hydra.openid.id-token' from Hardware Security Module."
time="2021-07-07T12:51:23Z" level=info msg="Using key pair 'hydra.jwt.access-token' from Hardware Security Module."
```

### Generating key pairs

Depending on HSM vendor, tools initializing tokens and generating keys vary.
Let's take a look how key pairs are generated in HSM quickstart container using
`pkcs11-tool` from OpenSC:

#### Initializing token

Different policies can apply for tokens, therefore Hydra HSM configuration
expects, that token where to find or generate keys already exists.

Just to demonstrate key pair generation we first initialize token.

```shell
$ pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so --slot 0 --init-token --so-pin 0000 --init-pin --pin 1234 --label hydra

Using slot 0 with a present token (0x2763db07)
Token successfully initialized
User PIN successfully initialized
```

Where parameter `--label hydra` value corresponds to value used in configuration
`HSM_TOKEN_LABEL` and `--pin 1234` to `HSM_PIN`

#### Generating key pair

Generating keypair for JSON Web Key `hydra.openid.id-token`

```shell
$ pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so \
--login --pin 1234 --token-label hydra \
--keypairgen --key-type rsa:4096 --usage-sign \
--label hydra.openid.id-token --id 68796472612e6f70656e69642e69642d746f6b656e

Key pair generated:
Private Key Object; RSA
  label:      hydra.openid.id-token
  ID:         68796472612e6f70656e69642e69642d746f6b656e
  Usage:      sign
Public Key Object; RSA 4096 bits
  label:      hydra.openid.id-token
  ID:         68796472612e6f70656e69642e69642d746f6b656e
  Usage:      verify
```

Where parameter `--id 68796472612e6f70656e69642e69642d746f6b656e` is the value
used as `kid` in JSON Web Key Set. It must be set as a big-endian hexadecimal
integer value.

Generating keypair for JSON Web Key `hydra.jwt.access-token`

```shell
$ pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so \
       --login --pin 1234 --token-label hydra \
       --keypairgen --key-type rsa:4096 --usage-sign \
       --label hydra.jwt.access-token --id 68796472612e6a77742e6163636573732d746f6b656e

Key pair generated:
Private Key Object; RSA
  label:      hydra.jwt.access-token
  ID:         68796472612e6a77742e6163636573732d746f6b656e
  Usage:      sign
Public Key Object; RSA 4096 bits
  label:      hydra.jwt.access-token
  ID:         68796472612e6a77742e6163636573732d746f6b656e
  Usage:      verify
```
