package api

import (
	"errors"
	"moon/internal/pkg/server/authent"
	"moon/internal/pkg/server/database"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
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
	log.Trace().Msg("Begin register")
	defer log.Trace().Msg("End register")

	bearerToken := r.Header.Get("Authorization")
	jwtString := strings.Fields(bearerToken)[1]

	jwt, err := authent.VerifyJwt(jwtString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, err := jwt.Claims.GetSubject()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error while getting sub")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Debug().Msgf("Sub received : %v", sub)

	_, res := database.FindUserByKCUID(sub, a.DB)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Error().Stack().Err(err).Msg("An error occured while finding if user exist by KCUID")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	newUser := database.User{
		KCUserID: sub,
	}

	a.DB.Save(&newUser)

	w.WriteHeader(http.StatusOK)
}
