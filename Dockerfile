FROM golang:1.7.1-alpine
MAINTAINER hteen <i@hteen.cn>

RUN apk add --no-cache git make openssl

RUN git clone https://github.com/tutumcloud/ngrok.git /ngrok

ADD *.sh /
ADD ./wln /wln

ENV DOMAIN ngrok.wlniao.com
ENV MY_FILES /wln
ENV TUNNEL_ADDR :4443
ENV HTTP_ADDR :80
ENV HTTPS_ADDR :443

EXPOSE 4443
EXPOSE 80
EXPOSE 443

CMD /bin/sh /server.sh
