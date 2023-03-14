FROM maven:3-jdk-11

WORKDIR /usr/src/mymaven

RUN wget https://gitlab.com/openid/conformance-suite/-/archive/release-v4.1.4/conformance-suite-release-v4.1.4.zip && \
  unzip conformance-suite-release-v4.1.4.zip -d . && \
  rm conformance-suite-release-v4.1.4.zip && \
  find conformance-suite-release-v4.1.4 -maxdepth 1 -mindepth 1 -exec mv {} . \; && \
  rmdir conformance-suite-release-v4.1.4

RUN mvn -B clean package -DskipTests && \
    apt-get update && apt-get install -y \
    redir ca-certificates

COPY ssl/ory-conformity.crt /etc/ssl/certs/
COPY ssl/ory-conformity.key /etc/ssl/private/
COPY ssl/ory-conformity.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

CMD java -Xdebug -Xrunjdwp:transport=dt_socket,address=*:9999,server=y,suspend=n -jar /usr/src/mymaven/target/fapi-test-suite.jar --fintechlabs.base_url=https://httpd:8443 --fintechlabs.devmode=true --fintechlabs.startredir=true
