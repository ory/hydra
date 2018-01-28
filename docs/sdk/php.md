## PHP SDK

### Installation

Installation is best done using [composer](https://getcomposer.org/)

```
composer require ory/hydra-sdk
```

If your project doesn't already make use of composer, you will need to include the resulting `vendor/autoload.php` file.

### Configuration

#### OAuth2 configuration

We need OAuth2 capabilities in order to make authorized API calls. You can either write your own OAuth2 mechanism or
use an existing one that has been preconfigured for use with Hydra. Here we use a modified version of the league OAuth2
client that has had this work done for us.

```sh
composer require tulip/oauth2-hydra
```

```php
// Get an access token using your account credentials.
// Note that if you are using the Hydra inside docker as per the getting started docs, the domain will be hydra:4444 from
// within another container.
$provider = new \Hydra\OAuth2\Provider\OAuth2([
    'clientId' => 'admin',
    'clientSecret' => 'demo-password',
    'domain' => 'http://localhost:4444',
]);

try {
    // Get an access token using the client credentials grant.
    // Note that you must separate multiple scopes with a plus (+)
    $accessToken = $provider->getAccessToken(
        'client_credentials', ['scope' => 'hydra.clients']
    );
} catch (\Hydra\Oauth2\Provider\Exception\ConnectionException $e) {
    die("Connection to hydra failed: " . $e->getMessage());
} catch (\Hydra\Oauth2\Provider\Exception\IdentityProviderException $e) {
    die("Failed to get an access token: " . $e->getMessage());
}

```

#### SDK configuration

Using `$accessToken` from the above steps, you may now use the Hydra SDK:

```php
$config = new \Hydra\SDK\Configuration();
$config->setHost('http://localhost:4444');
// Use true in production!
$config->setSSLVerification(false);
$config->setAccessToken($accessToken);

// Pass the config into an ApiClient. You will need this client in the next ste.
$hydraApiClient = new \Hydra\SDK\ApiClient($config);
```

### API Usage

There are several APIs made available, see [../../sdk/php/swagger/README.md](The full API docs) for a list of clients and methods.

For this example, lets use the OAuth2Api to get a list of clients and use the `$hydraApiClient` from above:

```php
$hydraOAuth2Api = new \Hydra\SDK\Api\OAuth2Api($hydraApiClient);

try {
    $clients = $hydraOAuthSDK->listOAuth2Clients();
} catch ( \Hydra\SDK\ApiException $e) {
    if ($e->getCode() == 400) {
        die("Permission denied to get clients. Check the scopes on your access token!");
    }
    die("Failed to get clients: ".$e->getMessage());
}
```
