# Why isn't the redirect url working?

Hydra enforces HTTPS for all hosts except localhost. Also make sure that the path is an exact match. `http://localhost:123/`
is not the same as `http://localhost:123`.