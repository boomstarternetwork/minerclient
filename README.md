# minerclient

It is miner client web server app dedicated to work with [minerserver](https://bitbucket.org/boomstarternetwork/minerserver/).

It wraps `cpuminer` and `ethminer` internally and uses them to mine
cryptocurrencies.

## Building release

To build release you need to collect `cpuminer` and `ethminer` binaries both
for linux and windows in `bin` path:
```bash
ls bin/
cpuminer  cpuminer.exe  ethminer  ethminer.exe
```

Then type:
```bash
make release
```
And you will get release archives in `release` path.