# minerclient

It is miner client web server app dedicated to work with [minerserver](https://github.com/boomstarternetwork/minerserver/).

It wraps `cpuminer` and `ethminer` internally and uses them to mine
cryptocurrencies.

## Dependencies

To build release you need to install several dependencies:
```bash
go get -u github.com/mattn/go-colorable
go get -u github.com/asticode/go-astilectron
go get -u github.com/asticode/go-astilectron-bundler/...
```

To build release you need to collect `cpuminer` and `ethminer` binaries both
for linux and windows in `bin` path:
```bash
ls bin/
cpuminer  cpuminer.exe  ethminer  ethminer.exe
```

## Building release

```bash
make release
```
And you will get release archives in `release` path.