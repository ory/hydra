# How can I import a custom CA for RethinkDB?

You can do so by specifying environment variables:

- `RETHINK_TLS_CERT_PATH`: The path to the TLS certificate (pem encoded) used to connect to rethinkdb.
- `RETHINK_TLS_CERT`: A pem encoded TLS certificate passed as string. Can be used instead of `RETHINK_TLS_CERT_PATH`.

or via command line flag:

```
--rethink-tls-cert-path string   Path to the certificate file to connect to rethinkdb over TLS (https). You can set RETHINK_TLS_CERT_PATH or RETHINK_TLS_CERT instead.
```