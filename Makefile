# Prepare minerclient archives with binaries and other files.
minerclient:
	astilectron-bundler -v

	rm -rf release
	mkdir -p release/minerclient

	cp -R bin/cpuminer bin/ethminer output/linux-amd64/minerclient \
		release/minerclient
	cd release && tar czf minerclient.amd64.tar.gz minerclient

	rm -rf release/minerclient/*

	cp -R bin/cpuminer.exe bin/ethminer.exe \
		output/windows-amd64/minerclient.exe release/minerclient
	cd release/ && zip -r minerclient.amd64.zip minerclient

	rm -rf release/minerclient
	rm -rf output

.PHONY: minerclient