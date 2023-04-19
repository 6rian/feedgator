package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

	apiRouter.Get("/users", app.requireAuth(app.handleGetUserByApiKey))
	apiRouter.Post("/users", app.handleCreateUser)

	apiRouter.Get("/feeds", app.handleGetFeeds)
	apiRouter.Post("/feeds", app.requireAuth(app.handleCreateFeed))

	apiRouter.Post("/feed_follows", app.requireAuth(app.handleCreateFeedFollow))
	apiRouter.Delete("/feed_follows/{id}", app.requireAuth(app.handleDeleteFeedFollow))

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
