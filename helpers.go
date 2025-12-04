package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
)

// validate if the chirp is shorter than 140 chars
func ValidateBody(chirp ChirpRequest) error {
	if len(chirp.Body) > 140 {
		return errors.New("chirp is too long")
	}
	return nil
}

func HandleBadWords(s string) string {
	const badWordReplacement string = "****"
	BadWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	for _, badWord := range BadWords {
		re := regexp.MustCompile("(?i)\\b" + regexp.QuoteMeta(badWord) + "\\b")
		s = re.ReplaceAllString(s, badWordReplacement)
	}

	return s
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {

	response := ErrorResponse{Body: msg}
	dat, er := json.Marshal(response)
	if er != nil {
		log.Printf("Error marshalling JSON: %s", er)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("content-type", "text/json")
	w.Write(dat)
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {

	dat, er := json.Marshal(payload)
	if er != nil {
		log.Printf("Error marshalling JSON: %s", er)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("content-type", "text/json")
	w.Write(dat)
}

func DecodeBody[T any](r *http.Request) (*T, error) {
	requestBody := r.Body
	defer r.Body.Close()
	decoder := json.NewDecoder(requestBody)
	var zero, something T
	err := decoder.Decode(&something)

	if err != nil {
		log.Printf("error decoding teh request %v", err)
		return &zero, err
	}

	return &something, nil
}

// database user struct has more that to display
func UserToUserNew(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
