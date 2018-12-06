# Prepare minerclient archives with binaries and other files.
minerclient:
	astilectron-bundler -v

	rm -rf release
	mkdir -p release/minerclient

	cp -R bin/cpuminer-linux bin/ethminer-linux \
		output/linux-amd64/minerclient release/minerclient
	cd release && tar czf minerclient-linux.amd64.tar.gz minerclient

	rm -rf release/minerclient/*

	cp -R bin/cpuminer.exe bin/ethminer.exe \
		output/windows-amd64/minerclient.exe release/minerclient
	cd release/ && zip -r minerclient-windows.amd64.zip minerclient

	rm -rf release/minerclient/*

	# TODO: compile and add cpuminer-darwin
	cp -R bin/ethminer-darwin \
		output/darwin-amd64/minerclient.app release/minerclient
	cd release/ && zip -r minerclient-darwin.amd64.zip minerclient

	rm -rf release/minerclient
	rm -rf output

.PHONY: minerclient
