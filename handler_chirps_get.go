package main

import ("net/http"
		"github.com/google/uuid"
		"database/sql"
		"errors")

func (a apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {		
		
		chirpID := r.PathValue("chirpID")
		parasedUUID, err := uuid.Parse(chirpID)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		chirp, err := a.db.GetChirp(r.Context(), parasedUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows){
			respondWithError(w, 404, err.Error())
			return
		}
			respondWithError(w, 500, err.Error())
			return
		}
		respondWithJSON(w, 200, Chirp{ID: chirp.ID,
			CreatedAt:chirp.CreatedAt,
			UpdatedAt:chirp.UpdatedAt,
			Body:chirp.Body,
			UserID: chirp.UserID.UUID})
		
	
 }
