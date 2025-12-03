package main

import (
	"chirpy/internal/database"
	"context"

	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, "<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %v times!</p>\n</body>\n</html>", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) middleWareInc(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		handle.ServeHTTP(w, r)
	})
}

// / api handlers from now now
func (cfg *apiConfig) CreatUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/json")
	email, err := DecodeBody[EmailRequest](r)
	ctx := context.Background()
	if err != nil {
		RespondWithError(w, 400, "Error decoding body")
		return
	}

	user, err := cfg.db.AddUser(ctx, database.AddUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     email.Email,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error putting user in the database %v", err))
		return
	}

	userNew := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	RespondWithJson(w, 201, userNew)
}

func (cfg *apiConfig) CreatChirpHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/json")
	chirpReq, err := DecodeBody[ChirpRequest](r)
	ctx := context.Background()
	if err != nil {
		RespondWithError(w, 400, "Error decoding body")
		return
	}

	err = ValidateBody(*chirpReq)

	if err != nil {

		RespondWithError(w, 400, "Chirp is too long")
	}

	chirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    chirpReq.UserID,
		Body:      chirpReq.Body,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error putting chirp in the database: %v", err))
		return
	}

	chirpNew := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	}

	RespondWithJson(w, 201, chirpNew)
}

/// admin methods here

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("the role is %s\n", cfg.role)

	if cfg.role != "admin" {
		RespondWithError(w, 403, "forbidden experiment")
		return
	}

	err := cfg.db.Reset(context.Background())
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error deleting useres form the database: %v", err))
		return
	}
	m := message{text: "OK"}

	RespondWithJson(w, 200, m)

}
