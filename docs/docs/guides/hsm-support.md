---
id: hsm-support
title: Hardware Security Module support for JSON Web Key Sets
---

The
[PKCS#11 Cryptographic Token Interface Standard](http://docs.oasis-open.org/pkcs11/pkcs11-base/v2.40/os/pkcs11-base-v2.40-os.html),
also known as Cryptoki, is one of the Public Key Cryptography Standards
developed by RSA Security. PKCS#11 defines the interface between an application
and a cryptographic device.

:::note

If a key is not found in the Hardware Security Module, the regular Software Key
Manager with AES-GCM software encryption will be used as a fallback. Storing
keys will always use the Software Key Manager as it is not possible to add keys
to a Hardware Security Module.

:::

PKCS#11 is used as a low-level interface to perform cryptographic operations
without the need for the application to directly interface a device through its
driver. PKCS#11 represents cryptographic devices using a common model referred
to simply as a token. An application can therefore perform cryptographic
operations on any device or token, using the same independent command set.

<a name="hsm-configuration"></a>

### Hardware Security Module configuration

Ory Hydra can be configured using environment variables as well as a
configuration file. For more information on configuration options, open the
configuration documentation:

&gt;&gt; https://www.ory.sh/hydra/docs/reference/configuration &lt;&lt;

```
HSM_ENABLED=true
HSM_LIBRARY=/path/to/hsm-vendor/library.so
HSM_TOKEN_LABEL=hydra
HSM_SLOT=0
HSM_PIN=1234
```

Token that is denoted by environment variables `HSM_TOKEN_LABEL` or `HSM_SLOT`
must preexist and optionally contain RSA or ECDSA key pairs with labels
`hydra.openid.id-token` and `hydra.jwt.access-token` depending on configuration.
**_If keys with these labels don't exist, they will be generated upon
startup._** If both `HSM_TOKEN_LABEL` and `HSM_SLOT` are set, `HSM_TOKEN_LABEL`
takes preference over `HSM_SLOT`. In this case first slot that contains this
label is used. `HSM_LIBRARY` must point to vendor specific PKCS#11 library or
SoftHSM library if you want to [test HSM support](#testing-with-softhsm).

<a name="pkcs11-attribute-mappings"></a>

### PKCS#11 attribute mappings to JSON Web Key Set attributes

When key pair is generated or requested from HSM, the `CKA_LABEL` attribute is
used as JSON Web Key Set name, `CKA_ID` attribute as `kid`. Key usage is
determined by private key attributes, where `CKA_SIGN` and `CKA_DECRYPT` are
mapped to `sig` and `enc` respectively and set as key `use` attribute.
Furthermore, `CKA_ID's` of key pair private/public handles must be identical.
Attribute `alg` is determined from `CKA_KEY_TYPE` and `CKA_ECDSA_PARAMS`.

<a name="supported-key-algorithms"></a>

### Supported key algorithms

Ory Hydra supports generating 4096 bit RSA, ECDSA keys with curves secp256r1 or
secp521r1. As of now PKCS#11 v2.4 doesn't support EdDSA keys using curve
Ed25519. However,
[PKCS#11 v3.0](https://docs.oasis-open.org/pkcs11/pkcs11-curr/v3.0/pkcs11-curr-v3.0.html)
contains support for EdDSA and therefore can be supported in upcoming versions.
Symmetric key algorithms are not supported because it would imply, that shared
HSM is used between server and authenticating client.

<a name="generating-key-pairs"></a>

### Generating key pairs

<a name="initializing-token"></a>

#### Initializing token

Different policies can apply for tokens, therefore HSM configuration expects,
that token where to find or generate keys already exists. Depending on HSM
vendor, tools initializing tokens and generating keys vary. To demonstrate key
pair generation we first initialize token using `pkcs11-tool` (see how to
[setup SoftHSM and OpenSC](#testing-with-softhsm))

```shell
$ pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so --slot 0 --init-token --so-pin 0000 --pin 1234 --init-pin --label hydra

Using slot 0 with a present token (0x2763db07)
Token successfully initialized
User PIN successfully initialized
```

Corresponding Ory Hydra configuration to access this token would be

```
HSM_ENABLED=true
HSM_LIBRARY=/usr/lib/softhsm/libsofthsm2.so
HSM_TOKEN_LABEL=hydra
HSM_SLOT=0
HSM_PIN=1234
```

<a name="generating-key-pair"></a>

#### Generating key pair

Generating RSA keypair for JSON Web Key `hydra.openid.id-token`

```shell
$ pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so \
--pin 1234 --token-label hydra \
--keypairgen --key-type rsa:4096 --usage-sign \
--label hydra.openid.id-token --id 746573742d6b65792d6964

Key pair generated:
Private Key Object; RSA
  label:      hydra.openid.id-token
  ID:         746573742d6b65792d6964
  Usage:      sign
Public Key Object; RSA 4096 bits
  label:      hydra.openid.id-token
  ID:         746573742d6b65792d6964
  Usage:      verify
```

| Parameter                           | Description                                                                                                                                                                                                             |
| :---------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--module`                          | Points to vendor specific PKCS#11 library or SoftHSM library when testing.                                                                                                                                              |
| `--pin 1234`                        | Pin that was used in token initialization to perform token operations.                                                                                                                                                  |
| `--token-label hydra`               | Performs key generation on first slot with label `hydra`. Use `--slot` option instead if you want to specify specific slot.                                                                                             |
| `--label hydra.openid.id-token`     | Sets key pair label attribute `CKA_LABEL` and is used as JSON Web Key Set name.                                                                                                                                         |
| `--id 746573742d6b65792d6964`       | Sets key pair id attribute `CKA_ID` and is used as JSON Web Key Set `kid`. It must be set as a big-endian hexadecimal integer value. `StringToHex("test-key-id") == 746573742d6b65792d6964`                             |
| `--keypairgen`                      | Perform key pair generation on token                                                                                                                                                                                    |
| `--key-type rsa:4096`               | Type and length of the key to generate. Supported values are `rsa:4096`, `EC:secp256r1` or `EC:secp521r1`. Sets `CKA_KEY_TYPE`,`CKA_ECDSA_PARAMS` attributes and is used to determine JSON Web Key Set `alg` attribute. |
| `--usage-sign` or `--usage-decrypt` | Sets private key attribute `CKA_SIGN` or `CKA_DECRYPT` respectively. Used to determine JSON Web Key Set `use` attribute.                                                                                                |

<a name="key-type-mappings"></a>

##### Key type mappings

| Key type     | JWT signing algorithm |
| :----------- | :-------------------- |
| rsa:4096     | RS256                 |
| EC:secp256r1 | ES256                 |
| EC:secp521r1 | ES512                 |

<a name="testing-with-softhsm"></a>

### Testing with SoftHSM

[SoftHSM](https://www.opendnssec.org/softhsm/) is an implementation of a
cryptographic store accessible through a PKCS #11 interface. You can use it to
explore PKCS#11 without having a Hardware Security Module. It is being developed
as a part of the OpenDNSSEC project.

[Follow these instructions to build SoftHSM from source.](https://wiki.opendnssec.org/display/SoftHSMDOCS/SoftHSM+Documentation+v2)

#### Install SoftHSM/OpenSC on Mac OSX

```shell
$ ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" 2> /dev/null
```

```shell
$ brew install softhsm
```

```shell
$ ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" 2> /dev/null
```

```shell
$ brew install opensc
```

#### Install SoftHSM/OpenSC on Ubuntu

```shell
$ sudo apt update
```

```shell
$ sudo apt install softhsm opensc
```

#### Install SoftHSM/OpenSC on Windows

Follow these instructions to install
[SoftHSM](https://github.com/disig/SoftHSM2-for-Windows) and
[OpenSC](https://github.com/OpenSC/OpenSC/wiki) on windows.

#### Run Ory Hydra with HSM using Docker

Alternatively you can use quickstart docker container that setups
SoftHSM/OpenSC, builds and runs Ory Hydra with HSM configuration enabled. You
need to have the latest [Docker](https://www.docker.com) and
[Docker Compose](https://docs.docker.com/compose) version installed. To run
quickstart HSM change into the directory with the Hydra source code and run the
following command:

```shell
$ docker-compose -f quickstart-hsm.yml up --build
```

Following is logged on startup if Hardware Security Module is successfully
configured:

```shell
$ docker logs ory-hydra-example--hydra
time="2021-07-07T12:51:23Z" level=info msg="Hardware Security Module is configured."
time="2021-07-07T12:51:23Z" level=info msg="JSON Web Key Set 'hydra.openid.id-token' does not exist yet, generating new key pair..."
```

#### Run Tests with HSM enabled using Docker

```shell
$ make quicktest-hsm
```
