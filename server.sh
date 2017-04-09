#!/bin/sh
set -e

if [ ! -f "${MY_FILES}/bin/ngrokd" ]; then
    echo "ngrokd is not build,will be build it now..."
    /bin/sh /build.sh
fi

${MY_FILES}/bin/ngrokd -tlsKey=${MY_FILES}/device.key -tlsCrt=${MY_FILES}/device.crt -domain="${DOMAIN}" -httpAddr=":80" -httpsAddr=":443" -tunnelAddr="${TUNNEL_ADDR}"
