FROM golang:1.7.1-alpine
MAINTAINER hteen <admin@wlniao.com>

ADD *.sh /
ADD ./wln /wln
RUN chmod 777 /wln/bin/ngrokd
RUN apk add --no-cache git make openssl
RUN git clone https://github.com/tutumcloud/ngrok.git /ngrok

ENV DOMAIN ngrok.wlniao.com
ENV MY_FILES /wln
ENV TUNNEL_ADDR :4443

EXPOSE 4443
EXPOSE 80

CMD /bin/sh /server.sh