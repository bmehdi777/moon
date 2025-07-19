package api

import "net/http"

func NewApiMux() *http.ServeMux {
	apiMux := http.NewServeMux()

	apiHealth := Health{}
	apiMux.HandleFunc("/health", apiHealth.router)

	apiVersion := Version{}
	apiMux.HandleFunc("/version", apiVersion.router)

	return apiMux
}
