package miner

import (
	"errors"

	"github.com/boomstarternetwork/minerclient/currency"
)

const PoolBaseAddr = "18.195.144.235"

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
	ProjectID string
	Currency  currency.Currency
	Address   string
	Worker    string
}

type Miner interface {
	Params() Params
	SetParams(p Params) error
	Start() error
	Stop() error
	ListenOutput() (chan string, error)
	ListenErrors() (chan error, error)
	ListenStop() (chan struct{}, error)
}
