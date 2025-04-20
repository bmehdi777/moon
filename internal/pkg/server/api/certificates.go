package api

import (
	"net/http"
)

type Certificates struct{}

func (c *Certificates) handleCreate( w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (c *Certificates) handleGet(id string, w http.ResponseWriter, r *http.Request) {

}
