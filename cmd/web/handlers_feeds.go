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

	// Automatically follow the feed
	ffParams := &database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := app.DB.CreateFeedFollow(r.Context(), *ffParams)
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not follow feed")
		return
	}

	type response struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}

	data := &response{
		Feed:       feed,
		FeedFollow: feedFollow,
	}

	respondWithJSON(w, http.StatusCreated, data)
}

func (app *application) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.DB.GetFeeds(r.Context())
	if err != nil {
		// TODO: add logging
		respondWithError(w, http.StatusInternalServerError, "could not retrieve feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}
