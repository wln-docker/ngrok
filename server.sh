#!/bin/sh
set -e

if [ ! -f "${MY_FILES}/bin/ngrokd" ]; then
    echo "ngrokd is not build,will be build it now..."
    /bin/sh /build.sh
fi

${MY_FILES}/bin/ngrokd -tlsKey=/cert/server.key -tlsCrt=/cert/server.crt -domain="${DOMAIN}" -httpAddr=":80" -httpsAddr=":443" -tunnelAddr=":4443"
