package main

import (
	"net/http"

	"github.com/google/uuid"
	"main.go/internal/auth"
	"main.go/internal/database"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting a chirp goes here
	chirpID := r.PathValue("chirpID")
	parsedUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ChirpID format", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret_key)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	Chirp, err := cfg.db.GetChirpID_ByUserID(r.Context(), database.GetChirpID_ByUserIDParams{
		UserID: userID,
		ID:     parsedUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unhautorized", err)
		return
	}

	Deleted_id, err := cfg.db.DeleteChirpById(r.Context(), Chirp.ID)
	if err != nil {
		if Deleted_id == uuid.Nil {
			respondWithError(w, http.StatusNotFound, "Chirp not found or already deleted", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]string{"message": "Chirp deleted successfully"})

}
