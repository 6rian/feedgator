package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	server := http.Server{
		Addr:    addr,
		Handler: routes(),
	}

	log.Printf("Starting server on %s", addr)
	err = server.ListenAndServe()
	log.Fatal(err)
}
