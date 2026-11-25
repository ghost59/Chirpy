package main

import ("chirpy/internal/database"
		"chirpy/internal/auth"
		"github.com/google/uuid"
		"time"
		"net/http"
		"encoding/json"
		"log")

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Password string `json:"password"`
	Is_chirpy_red bool `json:"is_chirpy_red"`
	
}

func (a *apiConfig) UserProfileSetup(w http.ResponseWriter,r *http.Request) {
	Newuser := database.CreateUserParams{ID: uuid.New(),
										CreatedAt: time.Now(),
										UpdatedAt:time.Now(),}
	err := json.NewDecoder(r.Body).Decode(&Newuser)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error: %s", err)
		return
	}
	hashedPassword, err := auth.HashPassword(Newuser.Password)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	Newuser.Password = hashedPassword
	
	user, err := a.db.CreateUser(r.Context(), Newuser)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	 
	responseUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Password: user.Password,
		Is_chirpy_red: user.IsChirpyRed,
	}

	respondWithJSON(w,201,responseUser)
	

	log.Println(user)
	
	
	
}

type NewPass struct  {
	Email string `json:"email"`
	Password string `json:"password"`
}

type updateUser struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (a *apiConfig) Update(w http.ResponseWriter, r *http.Request) {
	 token, err := auth.GetBearerToken(r.Header)
	 if err != nil {
	 	respondWithError(w, 401, "does not exist")
	 	return
	 }
	 user, err := auth.ValidateJWT(token,a.Secret)
	 
	 if err != nil {
	 	respondWithError(w, 401, "malformed")
	 	return
	 }

	 pass := NewPass{}
	 
	 
	err = json.NewDecoder(r.Body).Decode(&pass)
	if err != nil {
		respondWithError(w, 400, "sucks")
		return
	}

	
	hashedPassword, err := auth.HashPassword(pass.Password)
	if err != nil {
		respondWithError(w, 400, "malformed")
		return
	}
	

	
	updated := database.UpdateUserParams{
		Password: hashedPassword,
		Email: pass.Email,
		ID: user,
	}

	userUp, err := a.db.UpdateUser(r.Context(), updated)

	
	responseUser := updateUser{
		ID: userUp.ID,
		CreatedAt: userUp.CreatedAt,
		UpdatedAt: userUp.UpdatedAt,
		Email: userUp.Email,
	}

	respondWithJSON(w, 200, responseUser)
	 
}
