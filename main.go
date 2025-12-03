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
	dbQueries := database.New(db)
	cfg.db = *dbQueries

	fmt.Println("Server listening at ", port)
	serverMux := http.NewServeMux()
	filepathRoot := ""
	serverMux.Handle("POST /api/users", http.HandlerFunc(cfg.CreatUserHandler))
	serverMux.Handle("POST /api/chirps", http.HandlerFunc(cfg.CreatChirpHandler))
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
