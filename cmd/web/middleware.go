package main

import (
	"context"
	"net/http"
)

func (app *application) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := parseAuthorizationHeader(r.Header.Get("Authorization"), "ApiKey ")
		if err != nil {
			// TODO: add logging
			respondWithError(w, http.StatusUnauthorized, "missing api key")
			return
		}

		user, err := app.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid api key")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
