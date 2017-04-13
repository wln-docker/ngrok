#!/bin/sh
set -e

if [ ! -f "${MY_FILES}/ngrokd" ]; then
    echo "ngrokd is not build,will be build it now..."
    /bin/sh /build.sh
fi
chmod 777 ${MY_FILES}/ngrokd
${MY_FILES}/ngrokd -tlsKey=/cert/server.key -tlsCrt=/cert/server.crt -domain="${DOMAIN}" -httpAddr=":80" -httpsAddr=":443" -tunnelAddr=":4443"
