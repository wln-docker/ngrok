SELFPATH=$(cd "$(dirname "$0")"; pwd)
GOOS=`go env | grep GOOS | awk -F\" '{print $2}'`
GOARCH=`go env | grep GOARCH | awk -F\" '{print $2}'`


	cd /usr/local/go/src
	GOOS=linux GOARCH=386 ./make.bash
	cd /usr/local/ngrok/
	GOOS=linux GOARCH=386 make release-all

	cd /usr/local/ngrok/
	GOOS=windows GOARCH=386 make release-all

	cd /usr/local/ngrok/
	GOOS=darwin GOARCH=386 make release-all

#	cd /usr/local/go/src
#	GOOS=linux GOARCH=amd64 ./make.bash
	cd /usr/local/ngrok/
	GOOS=linux GOARCH=amd64 make release-all

	cd /usr/local/ngrok/
	GOOS=windows GOARCH=amd64 make release-all
	
	cd /usr/local/ngrok/
	GOOS=darwin GOARCH=amd64 make release-all
	
	cd /usr/local/go/src
	GOOS=linux GOARCH=arm ./make.bash
	cd /usr/local/ngrok/
	GOOS=linux GOARCH=arm make release-all