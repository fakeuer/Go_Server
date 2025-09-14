package main

import (
	"net/http"

	"github.com/google/uuid"
	"main.go/internal/database"
)

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	var chirps []database.Chirp
	var err error
	if authorID != "" {
		parsedUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id format", err)
			return
		}
		if sortOrder == "desc" {
			chirps, err = cfg.db.GetChirpsByAuthorIDDesc(r.Context(), parsedUUID)
		} else {
			chirps, err = cfg.db.GetChirpsByAuthorID(r.Context(), parsedUUID)
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps by author_id", err)
			return
		}
	} else if sortOrder == "desc" {
		chirps, err = cfg.db.GetAllChirpsDesc(r.Context())
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	var response []Chirp
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) HandlerChrispbyID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("ChirpID")
	parsedUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ChirpID format", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	})
}
