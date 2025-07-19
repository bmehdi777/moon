package api

import "net/http"

type Auth struct {}

func (a *Auth) router(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case http.MethodPost:
		a.generateAuthentToken(w, r)
	case http.MethodDelete:
		a.revokeAuthentToken(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *Auth) generateAuthentToken(w http.ResponseWriter, r *http.Request) {
	// verify the access token is right
	// generate an authent token
	// save it in database

}

func (a *Auth) revokeAuthentToken(w http.ResponseWriter, r *http.Request) {
	// verify the access token is right
	// delete the authent token from the database
}
