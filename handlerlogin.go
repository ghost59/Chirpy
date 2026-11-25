package main

import ("net/http"
		"chirpy/internal/auth"
		"github.com/google/uuid"
		"time"
		"encoding/json"
		"log"
		"chirpy/internal/database"
		)

type UserLog struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	Refresh_token string `json:"refresh_token"`
	Ischirpyred bool `json:"is_chirpy_red"`
}

type checker struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Expires_in_seconds int `json:"expires_in_seconds"`
} 

func (a *apiConfig) Login(w http.ResponseWriter, r *http.Request){
	check := checker{}
	err := json.NewDecoder(r.Body).Decode(&check)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	userLog, err := a.db.UserLookup(r.Context(), check.Email)

	if err != nil {
		respondWithError(w, 401, "Invalide email or passwrd")
		return
	}

	
	valid, err := auth.CheckPasswordHash(check.Password, userLog.Password)
	if err != nil{
		log.Println("Hash validation error: %s", err.Error())
		respondWithError(w,401, "Invalid email or password")
		return
	}

	if  valid == false {
		respondWithError(w, 401, "Invalid email or password")
		return
	}

	Refresh := auth.MakeRefreshToken()

	days := time.Now()
	days = days.Add(time.Hour * 24 * 60)

	 refreshParams := database.CreateRefreshTokenParams{
		Token: Refresh,
		UserID: userLog.ID,
		ExpiresAt: days,
		
	}

	 _, err =  a.db.CreateRefreshToken(r.Context(), refreshParams)

	if err != nil {
		respondWithError(w, 500, "Expired token")
		return 
	} 
	


	token, err :=  auth.MakeJWT(userLog.ID, a.Secret, time.Hour)

	if err != nil {
		respondWithError(w, 500, "Invalide pass")
		return
	}

	

	

	user := UserLog{
		ID: userLog.ID,
		CreatedAt: userLog.CreatedAt,
		UpdatedAt: userLog.UpdatedAt,
		Email: userLog.Email,
		Token: token,
		Refresh_token: Refresh,
		Ischirpyred: userLog.IsChirpyRed,
	}

	respondWithJSON(w, 200, user)
	
		
		
	
}
