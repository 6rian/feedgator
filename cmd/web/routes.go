package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/6rian/feedgator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	apiRouter := chi.NewRouter()

	apiRouter.Get("/readiness", handleReadiness)
	apiRouter.Get("/err", handleErr)
	apiRouter.Post("/users", app.handleCreateUser)

	router.Mount("/v1", apiRouter)

	return router
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, http.StatusOK, &resp{Status: "ok"})
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal server error")
}

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
