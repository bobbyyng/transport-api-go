package main

import (
	"context"
	"net/http"
	"os"
	"transport-api/internal/logger"
	"transport-api/internal/routes"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	if err := logger.Init("../../logs"); err != nil {
		logger.Fatal("Failed to initialize logger: " + err.Error())
	}

	// Import environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	// Connect to database
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Unable to connect to database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close(context.Background())
	err = db.Ping(context.Background())
	if err != nil {
		logger.Error("Error pinging the database: " + err.Error())
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router, db)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// Start server
	logger.Success("Connected to database")
	logger.Success("Starting server on port " + os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), router); err != nil {
		logger.Fatal("Could not start server: " + err.Error())
	}
}
