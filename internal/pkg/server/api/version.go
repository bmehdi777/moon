package api

import "net/http"

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
	w.Write([]byte(SERVER_VERSION + "\n"))
}
