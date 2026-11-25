package main

import (
		"net/http"
		"sync/atomic"
		"fmt"
		_ "github.com/lib/pq"
		"github.com/joho/godotenv"
		"database/sql"
		"os"
		"log"
		"chirpy/internal/database"
		
								)								

type apiConfig struct {
	filehits atomic.Int32
	db *database.Queries
	plat string
	Secret string
	ApiKey string 
}

func(cfg *apiConfig) MetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc( func(w http.ResponseWriter, req *http.Request) {
		cfg.filehits.Add(1)
		next.ServeHTTP(w, req)
		
		
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	hits := cfg.filehits.Load()
	s := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)
	fmt.Fprintf(w, "%s", s)
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.filehits.Store(0)

	if cfg.plat != "dev"{
		w.WriteHeader(403)
		s := "403 Forbidden"
		w.Write([]byte(s))
		return 
	}
	
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(500)
		s := "failed to delete users"
		w.Write([]byte(s))
		return
	}
	w.WriteHeader(200)
	}



func main() {

	godotenv.Load(".env")
	
	dbURL := os.Getenv("DB_URL")
	plats := os.Getenv("PLATFORM")
	scret := os.Getenv("SECRET")
	key := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

 	deQueries := database.New(db)
 	

	apiCfg := &apiConfig{db:deQueries, plat:plats, Secret:scret, ApiKey: key}
	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler:mux,
	}
		mux.HandleFunc("GET /api/healthz", func( w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	ok := "OK"
		w.Write([]byte(ok))
		
	})

	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.Handle("/app/",http.StripPrefix("/app/" ,apiCfg.MetricsInc(http.FileServer(http.Dir(".")))))
	
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
 	mux.HandleFunc("POST /api/chirps",apiCfg.Chirps)	
	mux.HandleFunc("POST /api/users", apiCfg.UserProfileSetup)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirps)		
	mux.HandleFunc("GET /api/chirps/{chirpID}",apiCfg.GetChirp)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/refresh", apiCfg.Refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.Revoked)
	mux.HandleFunc("PUT /api/users", apiCfg.Update)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks",apiCfg.WebHook)
	server.ListenAndServe()
	
}

