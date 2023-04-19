package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"
)

func respondWithError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":"%s"}`, msg)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func parseAuthorizationHeader(authString, prefix string) (string, error) {
	if utf8.RuneCountInString(authString) <= utf8.RuneCountInString(prefix) {
		return "", errors.New("auth token missing")
	}

	if !strings.HasPrefix(authString, prefix) {
		return "", errors.New("invalid authorization header")
	}

	token := authString[utf8.RuneCountInString(prefix):]
	return token, nil
}
