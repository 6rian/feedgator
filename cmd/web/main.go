package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/6rian/feedgator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type application struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("Could not open database connetion", err)
	}

	dbQueries := database.New(db)

	app := application{
		DB: dbQueries,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	server := http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	log.Printf("Starting server on %s", addr)
	err = server.ListenAndServe()
	log.Fatal(err)
}
