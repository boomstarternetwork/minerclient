package main

import (
	"html/template"
	"io"
	"net/http"

	"bitbucket.org/boomstarternetwork/minerclient/handler"
	"bitbucket.org/boomstarternetwork/minerclient/miner/minersBundle"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	bindAddress = "127.0.0.1:8080"
)

func main() {
	m := minersBundle.NewMinersBundle()

	h := handler.NewHandler(m)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("ui/templates/*.html")),
	}

	e.Static("/", "ui/static")

	e.GET("/", h.Index)
	e.POST("/start", h.Start)
	e.POST("/stop", h.Stop)
	e.GET("/miner", h.Miner)
	e.GET("/miner/output", h.MinerOutput)

	e.Logger.Fatal(e.Start(bindAddress))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name+".html", data)
}
