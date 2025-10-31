package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

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
	log.Trace().Msg("Begin health")
	defer log.Trace().Msg("End health")

	w.WriteHeader(http.StatusOK)
}
