package api

import (
	"net/http"
	"strings"
)

type App struct {
	HealthcheckHandler Healthcheck
}

func NewApp() App {
	return App{
		HealthcheckHandler: Healthcheck{},
	}
}

func (a *App) ServeHttp(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/api", "", 1)

	switch path {
	case "/healthcheck":
		a.HealthcheckHandler.ServeHttp(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
