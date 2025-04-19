package api

import "net/http"

type Certificates struct{}

func (c *Certificates) ServeHttp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.handlePost(r, w)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (c *Certificates) handlePost(r *http.Request, w http.ResponseWriter) {

}
