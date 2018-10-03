package ethminer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"bitbucket.org/boomstarternetwork/minerclient/currency"
	"bitbucket.org/boomstarternetwork/minerclient/miner"
)

type ethminer struct {
	params miner.Params
	cmd    *exec.Cmd
	output chan string
	errors chan error
	stop   chan struct{}
}

func NewEthminer() miner.Miner {
	return &ethminer{}
}

func (m *ethminer) Params() miner.Params {
	return m.params
}

func (m *ethminer) SetParams(p miner.Params) error {
	if err := validateParams(p); err != nil {
		return err
	}
	m.params = p
	return nil
}

func (m *ethminer) Start() error {
	if m.cmd != nil {
		m.Stop()
	}

	login := m.params.Address + "." + m.params.ProjectID
	if len(m.params.Worker) > 0 {
		login = login + "." + m.params.Worker
	}

	poolAddr, err := miner.PoolAddr(currency.Ethereum)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("stratum://%s@%s", login, poolAddr)

	m.cmd = exec.Command(path, "-U", url)

	m.output = make(chan string)
	m.errors = make(chan error)
	m.stop = make(chan struct{})

	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}

	scanwg := sync.WaitGroup{}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func() {
		scanwg.Add(1)
		defer scanwg.Done()
		for stdoutScanner.Scan() {
			m.output <- stdoutScanner.Text()
		}
	}()

	go func() {
		scanwg.Add(1)
		defer scanwg.Done()
		for stderrScanner.Scan() {
			m.output <- stderrScanner.Text()
		}
	}()

	err = m.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		scanwg.Wait()

		if err := stdoutScanner.Err(); err != nil {
			m.errors <- err
		}
		if err := stderrScanner.Err(); err != nil {
			m.errors <- err
		}

		err := m.cmd.Wait()
		if err != nil {
			m.errors <- err
		}

		// wait until possible errors above are handled
		time.Sleep(1 * time.Second)

		close(m.output)
		close(m.errors)
		close(m.stop)

		m.output = nil
		m.errors = nil
		m.stop = nil

		m.cmd = nil
	}()

	return nil
}

func (m *ethminer) Stop() {
	if m.cmd == nil {
		return
	}
	if runtime.GOOS == "windows" {
		// We are killing process because windows doesn't support any kind of
		// interrupt signals for graceful shutdown.
		m.cmd.Process.Kill()
	} else {
		// For other systems it is ok to send interrupt signal.
		m.cmd.Process.Signal(os.Interrupt)
	}
}

func (m *ethminer) ListenOutput() (chan string, error) {
	if m.output == nil {
		return nil, errors.New("output channel is closed")
	}
	return m.output, nil
}

func (m *ethminer) ListenErrors() (chan error, error) {
	if m.errors == nil {
		return nil, errors.New("errors channel is closed")
	}
	return m.errors, nil
}

func (m *ethminer) ListenStop() (chan struct{}, error) {
	if m.stop == nil {
		return nil, errors.New("stop channel is closed")
	}
	return m.stop, nil
}

func validateParams(p miner.Params) error {
	switch p.Currency {
	case currency.Ethereum:
	default:
		return errors.New("unsupported currency")
	}

	if p.Address == "" {
		return errors.New("address shouldn't be empty")
	}

	if p.ProjectID == "" {
		return errors.New("project shouldn't be empty")
	}

	return nil
}
