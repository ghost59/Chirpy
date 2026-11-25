package main

import ("net/http"
		"encoding/json"
		"github.com/google/uuid"
		"log"
		"chirpy/internal/auth"
		)

type Hook struct  {
	Event string `json:"event"`
	Data struct {
		User_id string `json:"user_id"`
	}
}

type Confirm struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ChirpyRed bool `json:"is_chirpy_red"`
}



func (a apiConfig) WebHook(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Not authorized")
		return
	}
	if token != a.ApiKey {
		respondWithError(w, 401, "Not authorized")
		return
	}
		
	hook := Hook{}
	err = json.NewDecoder(r.Body).Decode(&hook)
	if err != nil {
		respondWithError(w, 500, "This is trash")
		return
		
	}
	
	if hook.Event != "user.upgraded"{
		respondWithError(w, 204, "Not intersted buddy")
		return
	}

	userID, err := uuid.Parse(hook.Data.User_id)
	if err != nil {
		respondWithError(w, 404, "This ID is not convertable")
		return
	}

	upgraded, err := a.db.UpgradeUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "user can not be found")
		return
	}

	
	log.Println(upgraded)
	respondWithJSON(w, 204, Confirm{})
	
	
}
