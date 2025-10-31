package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

const SERVER_VERSION = "0.0.1"

type Version struct{}

func (v *Version) router(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case http.MethodGet:
		v.get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (v *Version) get(w http.ResponseWriter, _ *http.Request) {
	log.Trace().Msg("Begin version")
	defer log.Trace().Msg("End version")

	w.Write([]byte(SERVER_VERSION + "\n"))
}
