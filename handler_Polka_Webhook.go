package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"main.go/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if req.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "", nil)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkakey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	userUUID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UserID format", err)
		return
	}
	Id, err := cfg.db.SetChrispRedByUserId(r.Context(), userUUID)
	if err != nil {
		if Id == uuid.Nil {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't set chirp red", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]string{"message": "User upgraded to Chrisp Red successfully"})

}
