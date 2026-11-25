package main

import ("chirpy/internal/database"
		"github.com/google/uuid"
		"net/http"
		"encoding/json"
		"log"
		"time"
		"chirpy/internal/auth"
		"sort"
		)
type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type parameters struct {
    Body   string    `json:"body"`
}

type Chirpies struct {
	Body string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
func (a *apiConfig) Chirps(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "you are not who you say you are")
		return
	}

	userID, err := auth.ValidateJWT(token, a.Secret)
	if err != nil {
		respondWithError(w, 401, "you are not who you sare you are")
		return
	}

	
	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle":{},
		"sharbert":{},
		"fornax":  {},
	}
	params.Body = cleanChip(params.Body, badWords)
	NewChirp := database.CreateChirpParams{ID: uuid.New(),
											CreatedAt:time.Now(),
											UpdatedAt:time.Now(),
											Body: params.Body,
											UserID: uuid.NullUUID{
												UUID: userID,
												Valid: true,
											},
										}
	chirp, err := a.db.CreateChirp(r.Context(), NewChirp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("%s", err)
		return
	} 
	responseChirp := Chirp{
		ID: chirp.ID,
		CreatedAt:chirp.CreatedAt,
		UpdatedAt:chirp.UpdatedAt,
		Body:chirp.Body,
		UserID: chirp.UserID.UUID,
	}

	
	

	respondWithJSON(w, 201, responseChirp)
	
}


 		
func (a apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
 	s := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	n := r.URL.Query().Get("sort")
	
 	if s != "" {
 		user, err := uuid.Parse(s)
 		if err != nil {
 			respondWithError(w, 400, "Nope")
 			return
 		}
 		chirps , err = a.db.GetChirpsByUser(r.Context(),uuid.NullUUID{
 			UUID:user,
 			Valid: true,
 		})
 		
 	} else {
 		
		chirps, err = a.db.GetChirps(r.Context())
	 	}

	
	if err != nil {
				w.WriteHeader(500)
				return
		}
		
	
	CollectedChirps := []Chirp{}

	for _, chirp := range chirps {
		newChirp := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID.UUID,
			
		}

		CollectedChirps = append(CollectedChirps, newChirp)
		
		
	}

	if n == "desc" {
		
		sort.Slice(CollectedChirps, func(i, j int) bool{return CollectedChirps[i].CreatedAt.After(CollectedChirps[j].CreatedAt)})
	} else {

	sort.Slice(CollectedChirps, func(i, j int) bool{return CollectedChirps[i].CreatedAt.Before(CollectedChirps[j].CreatedAt)})
		}
	respondWithJSON(w, 200, CollectedChirps)
	
	
	
	
} 
