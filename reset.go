package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in development mode."))
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset users", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Hits reset to 0 and database cleared"))
}
