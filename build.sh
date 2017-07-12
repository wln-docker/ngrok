#!/bin/sh
set -e

cp -r /cert/ca.crt /ngrok/assets/client/tls/ngrokroot.crt
cp -r /cert/server.crt /ngrok/assets/server/tls/snakeoil.crt
cp -r /cert/server.key /ngrok/assets/server/tls/snakeoil.key

cd /ngrok
make release-server
#GOOS=linux GOARCH=386 make release-client
GOOS=linux GOARCH=amd64 make release-client
#GOOS=windows GOARCH=386 make release-client
#GOOS=windows GOARCH=amd64 make release-client
#GOOS=windows GOARCH=386 make release-server
#GOOS=windows GOARCH=amd64 make release-server
#GOOS=linux GOARCH=amd64 make release-server
#GOOS=darwin GOARCH=386 make release-client
#GOOS=darwin GOARCH=amd64 make release-client
#GOOS=linux GOARCH=arm make release-client

cp -rf /ngrok/bin/* ${MY_FILES}

echo "build ok !"
