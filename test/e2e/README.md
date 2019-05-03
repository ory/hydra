## Build the e2e docker file:

docker build -t oryd/e2e-env:latest -f Dockerfile-e2e-env .
docker push oryd/e2e-env:latest
