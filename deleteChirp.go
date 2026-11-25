package main

import (
	"net/http"
	"github.com/google/uuid"
	"chirpy/internal/auth"
	
)
type delete struct {
	Bearer string `json:"token"`
}

func (a apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	user, err := auth.ValidateJWT(token, a.Secret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	

	
	chirpID := r.PathValue("chirpID")
	parasedUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	chirp, err := a.db.GetChirp(r.Context(), parasedUUID)
	if err != nil {
		respondWithError(w, 404, "not real")
		return
	}

	if chirp.UserID.UUID != user {
		respondWithError(w, 403, "you can't delete this chirp")
		return
	}

	 err = a.db.DeleteChirp(r.Context(), parasedUUID)
	 if err != nil {
	 	respondWithError(w, 404, err.Error())
	 	return
	 }

	 
	 
	 respondWithJSON(w, 204,delete{Bearer:token} )
}
