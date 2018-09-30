# Prepare minerclient archives with binaries and other files.
minerclient:
	rm -rf release
	mkdir -p release/minerclient

	GOARCH=amd64 GOOS=linux go build -o release/minerclient/minerclient .
	cp -R ui bin/cpuminer bin/ethminer release/minerclient
	cd release && tar czf minerclient.amd64.tar.gz minerclient
	rm -rf release/minerclient/*

	GOARCH=amd64 GOOS=windows go build -o release/minerclient/minerclient.exe .
	cp -R ui bin/cpuminer.exe bin/ethminer.exe release/minerclient
	cd release/ && zip -r minerclient.amd64.zip minerclient
	rm -rf release/minerclient

.PHONY: minerclient