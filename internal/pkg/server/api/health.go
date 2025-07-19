package api

import "net/http"

type Health struct{}

func (h *Health) router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Health) get(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
