package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"

	"bitbucket.org/boomstarternetwork/minerclient/currency"
	"bitbucket.org/boomstarternetwork/minerclient/miner"
	"github.com/labstack/echo"
)

type Handler struct {
	projects        map[string]string
	miningProjectID string
	miner           miner.Miner
}

func NewHandler(miner miner.Miner) *Handler {
	return &Handler{
		projects: map[string]string{},
		miner:    miner,
	}
}

const projectsURL = "http://127.0.0.1/projects/list"

type project struct {
	ID   string
	Name string
}

type projectsResponse struct {
	Error  string
	Result []project
}

type currencyData struct {
	ID   string
	Name string
}

type indexPageData struct {
	Projects   []project
	Currencies []currencyData
	ProjectID  string
	Currency   string
	Address    string
	Worker     string
}

func (h *Handler) Index(c echo.Context) error {
	httpsRes, err := http.Get(projectsURL)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "500", 1)
	}

	resJSON, err := ioutil.ReadAll(httpsRes.Body)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "500", 2)
	}

	res := projectsResponse{}

	json.Unmarshal(resJSON, &res)

	if res.Error != "" {
		return c.Render(http.StatusInternalServerError, "500", 3)
	}

	h.projects = map[string]string{}

	for _, p := range res.Result {
		h.projects[p.ID] = p.Name
	}

	var cs []currencyData
	for _, c := range currency.List() {
		cs = append(cs, currencyData{
			ID:   string(c),
			Name: c.String(),
		})
	}

	return c.Render(http.StatusOK, "index",
		indexPageData{
			Projects:   res.Result,
			Currencies: cs,
			ProjectID:  c.QueryParam("project-id"),
			Currency:   c.QueryParam("currency"),
			Address:    c.QueryParam("address"),
			Worker:     c.QueryParam("worker"),
		})
}

func (h *Handler) Start(c echo.Context) error {
	projectID := c.FormValue("project-id")
	if _, exists := h.projects[projectID]; !exists {
		return c.Render(http.StatusNotFound, "404",
			"project-id="+projectID)
	}

	curr, err := currency.Parse(c.FormValue("currency"))
	if err != nil {
		return c.Render(http.StatusBadRequest, "400",
			err.Error()+", code=1")
	}

	h.miner.Stop()

	err = h.miner.SetParams(miner.Params{
		Currency: curr,
		Address:  c.FormValue("address"),
		Project:  projectID,
		Worker:   c.FormValue("worker"),
	})
	if err != nil {
		return c.Render(http.StatusBadRequest, "400",
			err.Error()+", code=2")
	}

	if err := h.miner.Start(); err != nil {
		return c.Render(http.StatusInternalServerError, "500", 3)
	}

	return c.Redirect(http.StatusFound, "/miner")
}

func (h *Handler) Stop(c echo.Context) error {
	projectID := c.FormValue("project-id")
	if _, exists := h.projects[projectID]; !exists {
		return c.Render(http.StatusNotFound, "404",
			"project-id="+projectID)
	}

	h.miner.Stop()

	mp := h.miner.Params()

	qp := url.Values{}
	qp.Set("project-id", projectID)
	qp.Set("currency", string(mp.Currency))
	qp.Set("address", mp.Address)
	qp.Set("worker", mp.Worker)

	return c.Redirect(http.StatusFound, "/?"+qp.Encode())
}

type minerPageData struct {
	Project  project
	Currency string
	Address  string
	Worker   string
}

func (h *Handler) Miner(c echo.Context) error {
	mp := h.miner.Params()

	projectName, _ := h.projects[mp.Project]

	return c.Render(http.StatusOK, "miner",
		minerPageData{
			Project: project{
				ID:   mp.Project,
				Name: projectName,
			},
			Currency: mp.Currency.String(),
			Address:  mp.Address,
			Worker:   mp.Worker,
		})
}

func (h *Handler) MinerOutput(c echo.Context) error {
	lines, err := h.miner.ListenOutput()
	if err != nil {
		return c.Render(http.StatusNotFound, "404", "code=1")
	}

	errs, err := h.miner.ListenErrors()
	if err != nil {
		return c.Render(http.StatusNotFound, "404", "code=2")
	}

	stop, err := h.miner.ListenStop()
	if err != nil {
		return c.Render(http.StatusNotFound, "404", "code=3")
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			select {
			case line, ok := <-lines:
				if !ok {
					continue
				}
				ws.Write([]byte(
					fmt.Sprintf("<p>%s</p>", line)))
			case err, ok := <-errs:
				if !ok {
					continue
				}
				ws.Write([]byte(
					fmt.Sprintf(`<p class="error">%s</p>`, err)))
			case <-stop:
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
