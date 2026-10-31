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
	jwtSecret      string
	polkaKey       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

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
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()

	// ==========================================
	// 1. Fileserver & Static Assets
	// ==========================================
	// Serves static files and client-side web app from the root directory.
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// ==========================================
	// 2. Health & Readiness
	// ==========================================
	// Checks if the server is up and operational.
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// ==========================================
	// 3. Authentication
	// ==========================================
	// Handles user login, session token refreshing, and token revocation.
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	// ==========================================
	// 4. Users
	// ==========================================
	// Handles user account creation and profile updates.
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUserUpdate)

	// ==========================================
	// 5. Chirps (Posts)
	// ==========================================
	// Handles creating, listing, fetching individual posts, and deleting chirps.
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerListChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	// ==========================================
	// 6. Webhooks & Third-Party Integrations
	// ==========================================
	// Handles incoming webhook notifications from external services like Polka.
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUserUpgrade)

	// ==========================================
	// 7. Admin & Metrics
	// ==========================================
	// Provides internal server metrics and database state reset capabilities.
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
