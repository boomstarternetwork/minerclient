package minersBundle

import (
	"errors"

	"bitbucket.org/boomstarternetwork/minerclient/currency"
	"bitbucket.org/boomstarternetwork/minerclient/miner"
	"bitbucket.org/boomstarternetwork/minerclient/miner/cpuminer"
	"bitbucket.org/boomstarternetwork/minerclient/miner/ethminer"
)

type minersBundle struct {
	params   miner.Params
	cpuminer miner.Miner
	ethminer miner.Miner
}

func NewMinersBundle() miner.Miner {
	return &minersBundle{
		cpuminer: cpuminer.NewCpuminer(),
		ethminer: ethminer.NewEthminer(),
	}
}

func (m *minersBundle) Params() miner.Params {
	return m.params
}

func (m *minersBundle) SetParams(p miner.Params) error {
	switch p.Currency {
	case currency.Bitcoin, currency.BitcoinCash,
		currency.Dash, currency.Litecoin:
		err := m.cpuminer.SetParams(p)
		if err != nil {
			return err
		}
		m.params = p
		return nil
	case currency.Ethereum:
		err := m.ethminer.SetParams(p)
		if err != nil {
			return err
		}
		m.params = p
		return nil
	default:
		return errors.New("unknown currency")
	}
}

func (m *minersBundle) Start() error {
	switch m.params.Currency {
	case currency.Bitcoin, currency.BitcoinCash,
		currency.Dash, currency.Litecoin:
		return m.cpuminer.Start()
	case currency.Ethereum:
		return m.ethminer.Start()
	default:
		return errors.New("unknown currency")
	}
}

func (m *minersBundle) Stop() {
	switch m.params.Currency {
	case currency.Bitcoin, currency.BitcoinCash,
		currency.Dash, currency.Litecoin:
		m.cpuminer.Stop()
	case currency.Ethereum:
		m.ethminer.Stop()
	}
}

func (m *minersBundle) ListenOutput() (chan string, error) {
	switch m.params.Currency {
	case currency.Bitcoin, currency.BitcoinCash,
		currency.Dash, currency.Litecoin:
		return m.cpuminer.ListenOutput()
	case currency.Ethereum:
		return m.ethminer.ListenOutput()
	default:
		return nil, errors.New("unknown currency")
	}
}

func (m *minersBundle) ListenErrors() (chan error, error) {
	switch m.params.Currency {
	case currency.Bitcoin, currency.BitcoinCash,
		currency.Dash, currency.Litecoin:
		return m.cpuminer.ListenErrors()
	case currency.Ethereum:
		return m.ethminer.ListenErrors()
	default:
		return nil, errors.New("unknown currency")
	}
}

func (m *minersBundle) ListenStop() (chan struct{}, error) {
	switch m.params.Currency {
	case currency.Bitcoin, currency.BitcoinCash,
		currency.Dash, currency.Litecoin:
		return m.cpuminer.ListenStop()
	case currency.Ethereum:
		return m.ethminer.ListenStop()
	default:
		return nil, errors.New("unknown currency")
	}
}
