package main

import (
	"encoding/json"
	"net/http"

	"github.com/amanuelmr/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		} `json:"data"`
	}
    
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if apiKey != apiConfig.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	_, err = apiConfig.dbQueries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}