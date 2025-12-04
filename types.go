package main

import (
	"chirpy/internal/database"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             database.Queries
	role           string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type message struct {
	text string
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type ChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}
type CleanedTextReponse struct {
	Body string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Body string `json:"error"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}