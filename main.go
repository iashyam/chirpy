package main

//let's start

import (
	// "database/sql"
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = "8080"

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	cfg := apiConfig{}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening the databse from url")
		return
	}

	cfg.role = os.Getenv("PLATFORM")
	cfg.secret_token = os.Getenv("TOKEN_STRING")

	polka_key := os.Getenv("POLKA_KEY")
	cfg.polka_key = polka_key

	dbQueries := database.New(db)
	cfg.db = *dbQueries

	fmt.Println("Server listening at ", port)
	serverMux := http.NewServeMux()
	filepathRoot := ""
	serverMux.Handle("POST /api/users", http.HandlerFunc(cfg.CreatUserHandler))
	serverMux.Handle("PUT /api/users", http.HandlerFunc(cfg.UpdateUserHandler))
	serverMux.Handle("POST /api/login", http.HandlerFunc(cfg.LoginHandler))
	serverMux.Handle("POST /api/refresh", http.HandlerFunc(cfg.RefreshTokenHandler))
	serverMux.Handle("POST /api/revoke", http.HandlerFunc(cfg.RevokeTokenHandler))
	serverMux.Handle("POST /api/chirps", http.HandlerFunc(cfg.CreatChirpHandler))
	serverMux.Handle("GET /api/chirps", http.HandlerFunc(cfg.ListChirpsHandler))
	serverMux.Handle("POST /api/polka/webhooks", http.HandlerFunc(cfg.UpgradeUserHandler))
	serverMux.Handle("GET /api/chirps/{chirpID}", http.HandlerFunc(cfg.GetChirpHandler))
	serverMux.Handle("DELETE /api/chirps/{chirpID}", http.HandlerFunc(cfg.DeleteChirpHandler))
	serverMux.Handle("GET /admin/metrics", &cfg)
	serverMux.Handle("/app/", cfg.middleWareInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serverMux.Handle("POST /admin/reset", http.HandlerFunc(cfg.ResetHandler))
	serverMux.Handle("/app/assets/logo.png", cfg.middleWareInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	localServer := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	err = localServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
