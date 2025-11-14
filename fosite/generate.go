// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/access_request.go github.com/ory/hydra/v2/fosite AccessRequester
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/access_response.go github.com/ory/hydra/v2/fosite AccessResponder
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/access_token_storage.go github.com/ory/hydra/v2/fosite/handler/oauth2 AccessTokenStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/access_token_storage_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 AccessTokenStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/access_token_strategy.go github.com/ory/hydra/v2/fosite/handler/oauth2 AccessTokenStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/access_token_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 AccessTokenStrategyProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_code_storage.go github.com/ory/hydra/v2/fosite/handler/oauth2 AuthorizeCodeStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_code_storage_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 AuthorizeCodeStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_code_strategy.go github.com/ory/hydra/v2/fosite/handler/oauth2 AuthorizeCodeStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_code_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 AuthorizeCodeStrategyProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_endpoint_handler.go github.com/ory/hydra/v2/fosite AuthorizeEndpointHandler
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_endpoint_handlers_provider.go github.com/ory/hydra/v2/fosite AuthorizeEndpointHandlersProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_request.go github.com/ory/hydra/v2/fosite AuthorizeRequester
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/authorize_response.go github.com/ory/hydra/v2/fosite AuthorizeResponder
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/client.go github.com/ory/hydra/v2/fosite Client
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/client_manager.go github.com/ory/hydra/v2/fosite ClientManager
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/oauth2_storage.go github.com/ory/hydra/v2/fosite/handler/oauth2 CoreStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/oauth2_strategy.go github.com/ory/hydra/v2/fosite/handler/oauth2 CoreStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/device_auth_storage.go github.com/ory/hydra/v2/fosite/handler/rfc8628 DeviceAuthStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/device_auth_storage_provider.go github.com/ory/hydra/v2/fosite/handler/rfc8628 DeviceAuthStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/device_code_strategy.go github.com/ory/hydra/v2/fosite/handler/rfc8628 DeviceCodeStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/device_code_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/rfc8628 DeviceCodeStrategyProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/device_rate_limit_strategy.go github.com/ory/hydra/v2/fosite/handler/rfc8628 DeviceRateLimitStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/device_rate_limit_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/rfc8628 DeviceRateLimitStrategyProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/hash.go github.com/ory/hydra/v2/fosite Hasher
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/open_id_connect_token_strategy.go github.com/ory/hydra/v2/fosite/handler/openid OpenIDConnectTokenStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/open_id_connect_token_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/openid OpenIDConnectTokenStrategyProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/open_id_connect_request_storage.go github.com/ory/hydra/v2/fosite/handler/openid OpenIDConnectRequestStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/open_id_connect_request_storage_provider.go github.com/ory/hydra/v2/fosite/handler/openid OpenIDConnectRequestStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/par_storage.go github.com/ory/hydra/v2/fosite PARStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/par_storage_provider.go github.com/ory/hydra/v2/fosite PARStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/pkce_request_storage.go github.com/ory/hydra/v2/fosite/handler/pkce PKCERequestStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/pkce_request_storage_provider.go github.com/ory/hydra/v2/fosite/handler/pkce PKCERequestStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/refresh_token_storage.go github.com/ory/hydra/v2/fosite/handler/oauth2 RefreshTokenStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/refresh_token_storage_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 RefreshTokenStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/refresh_token_strategy.go github.com/ory/hydra/v2/fosite/handler/oauth2 RefreshTokenStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/refresh_token_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 RefreshTokenStrategyProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/request.go github.com/ory/hydra/v2/fosite Requester
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/resource_owner_password_credentials_grant_storage.go github.com/ory/hydra/v2/fosite/handler/oauth2 ResourceOwnerPasswordCredentialsGrantStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/resource_owner_password_credentials_grant_storage_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 ResourceOwnerPasswordCredentialsGrantStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/revocation_handler.go github.com/ory/hydra/v2/fosite RevocationHandler
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/revocation_handlers_provider.go github.com/ory/hydra/v2/fosite RevocationHandlersProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/rfc7523_key_storage.go github.com/ory/hydra/v2/fosite/handler/rfc7523 RFC7523KeyStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/rfc7523_key_storage_provider.go github.com/ory/hydra/v2/fosite/handler/rfc7523 RFC7523KeyStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/storage.go github.com/ory/hydra/v2/fosite Storage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/token_endpoint_handler.go github.com/ory/hydra/v2/fosite TokenEndpointHandler
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/token_introspector.go github.com/ory/hydra/v2/fosite TokenIntrospector
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/token_revocation_storage.go github.com/ory/hydra/v2/fosite/handler/oauth2 TokenRevocationStorage
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/token_revocation_storage_provider.go github.com/ory/hydra/v2/fosite/handler/oauth2 TokenRevocationStorageProvider
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/transactional.go github.com/ory/hydra/v2/fosite Transactional
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/user_code_strategy.go github.com/ory/hydra/v2/fosite/handler/rfc8628 UserCodeStrategy
//go:generate go run go.uber.org/mock/mockgen -package internal -destination internal/user_code_strategy_provider.go github.com/ory/hydra/v2/fosite/handler/rfc8628 UserCodeStrategyProvider
