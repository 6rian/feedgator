package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/6rian/feedgator/internal/database"
	"github.com/google/uuid"
)

func (app *application) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"string"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	if err := decoder.Decode(&input); err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusBadRequest, "could not parse request")
		return
	}

	user := r.Context().Value("user").(database.User)

	// TODO: add validation

	feedParams := &database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      input.Name,
		Url:       input.URL,
		UserID:    user.ID,
	}

	feed, err := app.DB.CreateFeed(r.Context(), *feedParams)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not create feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, feed)
}
