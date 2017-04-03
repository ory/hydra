# Hydra Policy Tests

This directory holds the test suite associated with verifying Bridg's Hydra policies.

The following roles have been covered in these tests (all gateway API endpoints):
1. bridg-admin
2. account-admin
3. account-viewer

## Running Tests

To setup a local testing environment run the following commands:
1. `.citizen/post-build`
2. `docker-compose up -d` (wait for Hydra to start the web server)
3. `.citizen/post-system-start`

To run all tests execute the following command:
`docker-compose -f docker-compose-test.yml run --rm test`

To run a specific test file execute the following command (bridg-admin example:
`docker-compose -f docker-compose-test.yml run --rm test test/bridg_admin_policies.spec.js`

