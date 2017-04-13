FROM golang:1.7.1-alpine
MAINTAINER wlniao <admin@wln.io>

ADD *.sh /
RUN chmod 777 /wln/bin/ngrokd
RUN apk add --no-cache git make openssl
RUN git clone https://github.com/inconshreveable/ngrok.git /ngrok

ENV DOMAIN wlniao.cn
ENV MY_FILES /wln

EXPOSE 4443
EXPOSE 443
EXPOSE 80
VOLUME ["/wln"]
CMD /bin/sh /server.sh