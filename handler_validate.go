package main

import ("encoding/json"
		"log"
		"net/http"
		"strings"
		)


type chirpRequest  struct {
	Body string `json:"body, omitempty"`
	Password string `json:"password"`
		
}

type chirpResponse struct {

	CleanedBody string `json:cleaned_body`	
}


type Validator struct {
	Valid bool `json:"valid"` 
}

type ErrorResponse struct {
	Error string `json:"error"`
}



func respondWithError(w http.ResponseWriter, code int, msg string){
	errResp := ErrorResponse{ Error: msg}
	
 	respondWithJSON(w, code, errResp)
	
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	pay, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")

	
	
	w.WriteHeader(code)
	w.Write(pay)
	
			
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request){
	params := chirpRequest{}
	if err:= json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding paramters: %s",err)	
		respondWithError(w, 500, err.Error())
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
	

	respondWithJSON(w, 200, chirpResponse{CleanedBody: cleanChip(params.Body, badWords),
	})
	
	
		

	
}


func cleanChip(body string, badWords map[string]struct{}) string{

	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok{
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
