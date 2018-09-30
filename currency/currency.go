package currency

import (
	"github.com/getlantern/errors"
)

type Currency string

const (
	Bitcoin     Currency = "bitcoin"
	BitcoinCash Currency = "bitcoin-cash"
	Dash        Currency = "dash"
	Ethereum    Currency = "ethereum"
	Litecoin    Currency = "litecoin"
	Unknown     Currency = ""
)

func (c Currency) String() string {
	switch c {
	case Bitcoin:
		return "Bitcoin"
	case BitcoinCash:
		return "Bitcoin Cash"
	case Dash:
		return "Dash"
	case Ethereum:
		return "Ethereum"
	case Litecoin:
		return "Litecoin"
	default:
		return "Unknown"
	}
}

func Parse(s string) (Currency, error) {
	switch s {
	case string(Bitcoin):
		return Bitcoin, nil
	case string(BitcoinCash):
		return BitcoinCash, nil
	case string(Dash):
		return Dash, nil
	case string(Ethereum):
		return Ethereum, nil
	case string(Litecoin):
		return Litecoin, nil
	default:
		return Unknown, errors.New("unknown currency")
	}
}

func List() []Currency {
	return []Currency{
		Bitcoin,
		BitcoinCash,
		Dash,
		Ethereum,
		Litecoin,
	}
}
