package main

import (
	"encoding/json"
	"net/http"

	"github.com/junaidmd0980/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if match == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	} else {
		respondWithJSON(w, http.StatusOK, response{
			User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
		})
	}
}
