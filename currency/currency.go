package currency

import (
	"errors"
	"strings"
)

type Currency int

const (
	Unknown Currency = iota
	Bitcoin
	BitcoinCash
	Dash
	Ethereum
	Litecoin
)

func (c Currency) ID() string {
	switch c {
	case Bitcoin:
		return "bitcoin"
	case BitcoinCash:
		return "bitcoin-cash"
	case Dash:
		return "dash"
	case Ethereum:
		return "ethereum"
	case Litecoin:
		return "litecoin"
	default:
		return "unknown"
	}
}

func (c Currency) Name() string {
	return c.String()
}

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

func ParseID(s string) (Currency, error) {
	switch s {
	case Bitcoin.ID():
		return Bitcoin, nil
	case BitcoinCash.ID():
		return BitcoinCash, nil
	case Dash.ID():
		return Dash, nil
	case Ethereum.ID():
		return Ethereum, nil
	case Litecoin.ID():
		return Litecoin, nil
	default:
		return Unknown, errors.New("unknown currency: " + s)
	}
}

func (c *Currency) UnmarshalJSON(b []byte) (err error) {
	*c, err = ParseID(strings.Trim(string(b), `"`))
	return err
}

func List() []Currency {
	return []Currency{
		// Bitcoin,
		// BitcoinCash,
		// Dash,
		Ethereum,
		// Litecoin,
	}
}

type Data struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListData() []Data {
	var ds []Data
	for _, c := range List() {
		ds = append(ds, Data{
			ID:   c.ID(),
			Name: c.Name(),
		})
	}
	return ds
}
