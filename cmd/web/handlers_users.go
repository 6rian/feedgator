package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/6rian/feedgator/internal/database"
	"github.com/google/uuid"
)

func (app *application) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	input := reqBody{}
	if err := decoder.Decode(&input); err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusBadRequest, "could not parse request")
		return
	}

	userParams := &database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      input.Name,
	}

	user, err := app.DB.CreateUser(r.Context(), *userParams)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (app *application) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, r.Context().Value("user"))
}
