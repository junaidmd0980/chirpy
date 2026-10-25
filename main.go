package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/junaidmd0980/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if dbURL == "" {
		log.Fatal("DB_URl must set in the environment")
	}

	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	dbQueries := database.New(dbConnection)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux := http.NewServeMux()
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerListChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server is listening on %s...", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
