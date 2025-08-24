package api

import (
	"net/http"

	"gorm.io/gorm"
)

func NewApiMux(db *gorm.DB) *http.ServeMux {
	apiMux := http.NewServeMux()

	apiHealth := Health{}
	apiMux.HandleFunc("/health", apiHealth.router)

	apiVersion := Version{}
	apiMux.HandleFunc("/version", apiVersion.router)

	apiUser := User{
		DB: db,
	}
	apiMux.HandleFunc("/users", apiUser.router)

	return apiMux
}
