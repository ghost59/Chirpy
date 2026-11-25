package main

import ("net/http"
		"chirpy/internal/auth"
		"time"
		)

type validRefresh struct {
	Token string `json:"token"`
}

func (a *apiConfig) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "does not exist")
		return
	}

	user, err := a.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Does not exist")
		return 
	}

	newToken, err := auth.MakeJWT(user.ID, a.Secret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "does not exist")
		return
	}

	
	valid := validRefresh{ Token: newToken}

	respondWithJSON(w, 200, valid)
	

	
	
}

func (a *apiConfig) Revoked(w http.ResponseWriter, r *http.Request) {
	
	token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, "does not exist")
			return
		}
    _,  err = a.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "does not exist")
		return
	}

	w.WriteHeader(204)
	
	
}
