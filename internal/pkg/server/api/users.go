package api

import (
	"log"
	"moon/internal/pkg/server/authent"
	"moon/internal/pkg/server/database"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func (a *User) router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.register(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *User) register(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	jwtString := strings.Fields(bearerToken)[1]

	jwt, err := authent.VerifyJwt(jwtString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, err := jwt.Claims.GetSubject()
	if err != nil {
		log.Println("Error while getting sub : ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("Sub received : ", sub)

	_, res := database.FindUserByKCUID(sub, a.DB)
	if res.Error != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if res.RowsAffected != 0 {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	newUser := database.User{
		KCUserID: sub,
	}

	a.DB.Save(&newUser)

	w.WriteHeader(http.StatusOK)
}
