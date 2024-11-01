FROM maven:3-jdk-11

WORKDIR /usr/src/mymaven
RUN apt-get update && apt-get install -y \
    unzip \
    wget \
    redir \
    ca-certificates

ARG CONFORMANCE_SUITE_VERSION=v4.1.45
RUN wget https://gitlab.com/openid/conformance-suite/-/archive/release-${CONFORMANCE_SUITE_VERSION}/conformance-suite-release-${CONFORMANCE_SUITE_VERSION}.zip && \
  unzip conformance-suite-release-${CONFORMANCE_SUITE_VERSION}.zip -d . && \
  rm conformance-suite-release-${CONFORMANCE_SUITE_VERSION}.zip && \
  find conformance-suite-release-${CONFORMANCE_SUITE_VERSION} -maxdepth 1 -mindepth 1 -exec mv {} . \; && \
  rmdir conformance-suite-release-${CONFORMANCE_SUITE_VERSION} && \
  mvn -B -Dmaven.test.skip -Dpmd.skip clean package

COPY ssl/ory-conformity.crt /etc/ssl/certs/
COPY ssl/ory-conformity.key /etc/ssl/private/
COPY ssl/ory-conformity.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

CMD java -Xdebug -Xrunjdwp:transport=dt_socket,address=*:9999,server=y,suspend=n -jar /usr/src/mymaven/target/fapi-test-suite.jar --fintechlabs.base_url=https://httpd:8443 --fintechlabs.devmode=true --fintechlabs.startredir=true
