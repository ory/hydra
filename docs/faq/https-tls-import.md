# How can I import TLS certificates?

You can import TLS certificates when running `hydra host`. This can be done by setting the following environment variables:

**Read from file**
- `HTTPS_TLS_CERT_PATH`: The path to the TLS certificate (pem encoded).
- `HTTPS_TLS_KEY_PATH`: The path to the TLS private key (pem encoded).

**Embedded**
- `HTTPS_TLS_CERT`: A pem encoded TLS certificate passed as string. Can be used instead of TLS_CERT_PATH.
- `HTTPS_TLS_KEY`: A pem encoded TLS key passed as string. Can be used instead of TLS_KEY_PATH.

Or by specifying the following flags:

```
--https-tls-cert-path string   Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
--https-tls-key-path string    Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
```