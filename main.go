package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/boomstarternetwork/minerclient/currency"
	"github.com/boomstarternetwork/minerclient/miner"
	"github.com/boomstarternetwork/minerclient/miner/minersBundle"
	"github.com/pkg/errors"
)

const minerserverAddr = "18.195.144.235:8080"

var m miner.Miner

func main() {

	m = minersBundle.NewMinersBundle()

	// Run bootstrap
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            "minerclient",
			AppIconDefaultPath: "resources/icon.png",
			AppIconDarwinPath:  "resources/icon.icns",
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage: "index.html",
			Options: &astilectron.WindowOptions{
				Center: astilectron.PtrBool(true),
				Width:  astilectron.PtrInt(1024),
				Height: astilectron.PtrInt(600),
			},
			MessageHandler: messageHandler,
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

func messageHandler(w *astilectron.Window, msg bootstrap.MessageIn) (interface{}, error) {
	switch msg.Name {
	case "getProjects":
		astilog.Info("getProjects js message")
		ps, err := getProjects()
		if err != nil {
			return err.Error(), err
		}
		return ps, nil
	case "getCurrencies":
		astilog.Info("getCurrencies js message")
		return currency.ListData(), nil
	case "startMining":
		astilog.Info("startMining js message")
		var p miner.Params
		if err := json.Unmarshal(msg.Payload, &p); err != nil {
			astilog.Errorf("Failed to JSON unmarshal payload: %v", err)
			err = errors.New(
				"failed to unmarshal payload: " + err.Error())
			return err.Error(), err
		}
		err := startMining(w, p)
		if err != nil {
			return err.Error(), err
		}
		return true, nil
	case "stopMining":
		astilog.Info("stopMining js message")
		m.Stop()
		return true, nil
	}
	return nil, errors.New("unknown message")
}

type project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectsResponse struct {
	Error  string    `json:"error"`
	Result []project `json:"result"`
}

const projectsURL = "http://" + minerserverAddr + "/projects/list"

func getProjects() ([]project, error) {
	httpsRes, err := http.Get(projectsURL)
	if err != nil {
		astilog.Errorf("Failed to http get projects URL: %v", err)
		return nil, err
	}

	resJSON, err := ioutil.ReadAll(httpsRes.Body)
	if err != nil {
		astilog.Errorf("Failed to read projects response body: %v", err)
		return nil, err
	}

	res := projectsResponse{}

	err = json.Unmarshal(resJSON, &res)
	if err != nil {
		astilog.Errorf("Failed to unmarshal response JSON: %v", err)
		return nil, err
	}

	if res.Error != "" {
		err := errors.New(res.Error)
		astilog.Errorf("Response error: %v", err)
		return nil, err
	}

	return res.Result, nil
}

func startMining(w *astilectron.Window, p miner.Params) error {
	m.Stop()

	if err := m.SetParams(p); err != nil {
		astilog.Errorf("Failed to set miner params: %v", err)
		return errors.New("failed to set miner params: " + err.Error())
	}

	if err := m.Start(); err != nil {
		astilog.Errorf("Failed to start miner: %v", err)
		return errors.New("failed to start miner: " + err.Error())
	}

	lines, err := m.ListenOutput()
	if err != nil {
		astilog.Errorf("Miner listen output error: %v", err)
		return errors.New("miner listen output error: " + err.Error())
	}

	errs, err := m.ListenErrors()
	if err != nil {
		astilog.Errorf("Miner listen errors error: %v", err)
		return errors.New("miner listen errors error: " + err.Error())
	}

	stop, err := m.ListenStop()
	if err != nil {
		astilog.Errorf("Miner listen stop error: %v", err)
		return errors.New("miner listen stop error: " + err.Error())
	}

	go func() {
		for {
			select {
			case line, ok := <-lines:
				if !ok {
					continue
				}
				lineJSON, _ := json.Marshal(line)
				_ = w.SendMessage(bootstrap.MessageIn{
					Name:    "logLine",
					Payload: json.RawMessage(lineJSON),
				})
			case err, ok := <-errs:
				if !ok {
					continue
				}
				errJSON, _ := json.Marshal(err.Error())
				_ = w.SendMessage(bootstrap.MessageIn{
					Name:    "error",
					Payload: json.RawMessage(errJSON),
				})
			case <-stop:
				astilog.Info("Miner stopped")
				return
			}
		}
	}()

	return nil
}
