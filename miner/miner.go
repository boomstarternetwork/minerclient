package miner

import (
	"bitbucket.org/boomstarternetwork/minerclient/currency"
	"github.com/getlantern/errors"
)

const PoolBaseAddr = "127.0.0.1"

const (
	BitcoinPoolAddr     = PoolBaseAddr + ":3000"
	BitcoinCashPoolAddr = PoolBaseAddr + ":3001"
	DashPoolAddr        = PoolBaseAddr + ":3002"
	EthereumPoolAddr    = PoolBaseAddr + ":3003"
	LitecoinPoolAddr    = PoolBaseAddr + ":3004"
)

func PoolAddr(c currency.Currency) (string, error) {
	switch c {
	case currency.Bitcoin:
		return BitcoinPoolAddr, nil
	case currency.BitcoinCash:
		return BitcoinCashPoolAddr, nil
	case currency.Dash:
		return DashPoolAddr, nil
	case currency.Ethereum:
		return EthereumPoolAddr, nil
	case currency.Litecoin:
		return LitecoinPoolAddr, nil
	default:
		return "", errors.New("unknown pool address")
	}
}

type Params struct {
	Currency currency.Currency
	Address  string
	Project  string
	Worker   string
}

type Miner interface {
	Params() Params
	SetParams(p Params) error
	Start() error
	Stop()
	ListenOutput() (chan string, error)
	ListenErrors() (chan error, error)
	ListenStop() (chan struct{}, error)
}
