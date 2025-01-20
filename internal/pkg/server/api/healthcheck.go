package api

import "net/http"

type Healthcheck struct{}

func (h *Healthcheck) ServeHttp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *Healthcheck) handleGet(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
