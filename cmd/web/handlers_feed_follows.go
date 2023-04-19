package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/6rian/feedgator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app *application) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	type params struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	if err := decoder.Decode(&input); err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusBadRequest, "could not parse request")
		return
	}

	user := r.Context().Value("user").(database.User)
	feedId, err := uuid.Parse(input.FeedID)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusUnprocessableEntity, "invalid feed ID")
		return
	}

	feedFollowParams := &database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedId,
	}

	feedFollow, err := app.DB.CreateFeedFollow(r.Context(), *feedFollowParams)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollow)
}

func (app *application) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idAsUUID, err := uuid.Parse(id)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusUnprocessableEntity, "invalid feed follow ID")
		return
	}

	err = app.DB.DeleteFeedFollow(r.Context(), idAsUUID)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not delete feed follow")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handleGetFeedFollowsByUserID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)

	feedFollows, err := app.DB.GetFeedFollowsByUserID(r.Context(), user.ID)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not get feed follows")
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}
