## cli-example seeds

This folder contains scripts and examples for seeding Hydra with a predefined
set of Clients and Policies. You may use this as a starting point for seeding a
new production installation.

#### Contents

| Item            | Description |
| --------------- | ----------- |
| clients         | A list of JSON files that define the Hydra OAuth2 clients that should be created.
| policies        | A list of JSON files that define the Hydra Ladon policies that should be created.
| import-clients  | An example script that demonstrates importing the clients JSON files.
| import-policies | An example script that demonstrates importing the policies JSON files.

## Usage

1. [Install Hydra.][1] This will be the client that you connect to the remote Hydra installation.
2. Obtain the generated super-client credentials logged upon Hydra installation.
3. Connect to Hydra using the super-client, via `hydra token client`.
4. Change all `client_secret` attributes in clients to a very secure and random secret.
5. Import clients, via `import-clients` script.
6. Import policies, via `import-policies` script.
7. Ensure that you have created a new super-client.
8. Delete the original, generated super-client.


[1]: https://ory-am.gitbooks.io/hydra/content/install.html#installing-hydra "Hydra Installation Documentation"
