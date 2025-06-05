FROM debian:12-slim
RUN apt-get update \
	&& apt-get install -y apache2 ssl-cert ca-certificates \
	&& apt-get clean
RUN \
	echo 'Listen 8443' > /etc/apache2/ports.conf \
	&& a2enmod headers proxy proxy_ajp proxy_http rewrite ssl \
	&& a2dissite 000-default.conf

COPY httpd/server.conf /etc/apache2/sites-enabled

COPY ssl/ory-conformity.crt /etc/ssl/certs/
COPY ssl/ory-conformity.key /etc/ssl/private/
COPY ssl/ory-conformity.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

ENTRYPOINT ["apachectl", "-DFOREGROUND"]
