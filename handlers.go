package main

import (
	"chirpy/internal/auth"
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
	userRequest, err := DecodeBody[UserRequest](r)
	ctx := context.Background()
	if err != nil {
		RespondWithError(w, 400, "Error decoding body")
		return
	}

	hashedPass, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		RespondWithError(w, 400, "Error hasing passowrd")
		return
	}

	user, err := cfg.db.AddUser(ctx, database.AddUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          userRequest.Email,
		HashedPassword: hashedPass,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error putting user in the database %v", err))
		return
	}

	userNew := UserToUserNew(user)
	RespondWithJson(w, 201, userNew)
}

func (cfg *apiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/json")
	userRequest, err := DecodeBody[UserRequest](r)
	ctx := context.Background()
	if err != nil {
		RespondWithError(w, 400, "Error decoding body")
		return
	}

	hashedPass, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		RespondWithError(w, 400, "Error hasing passowrd")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("Error while confirming login: %v", err))
		return
	}

	uid, err := auth.ValidateJWT(token, cfg.secret_token)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("Unauthorized Action %v", err))
		return
	}

	user, err := cfg.db.UpdateUser(ctx, database.UpdateUserParams{
		ID:             uid,
		Email:          userRequest.Email,
		HashedPassword: hashedPass,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error putting user in the database %v", err))
		return
	}

	userNew := UserToUserNew(user)
	RespondWithJson(w, 200, userNew)
}
func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/json")
	loginRequest, err := DecodeBody[LoginRequest](r)
	ctx := context.Background()
	if err != nil {
		RespondWithError(w, 400, "Error decoding body")
		return
	}

	user, err := cfg.db.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		RespondWithError(w, 400, "Error finding user with that email")
		return
	}

	matchPassword, err := auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, 403, "Incorrect email or password!")
		return
	}

	if !matchPassword {
		RespondWithError(w, 401, "Incorrect email or password!")
		return
	}

	expiresIN := loginRequest.ExpiresInSeconds
	if expiresIN == 0 {
		expiresIN = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret_token, time.Duration(expiresIN)*time.Second)

	if err != nil {
		RespondWithError(w, 400, "Error loggin in")
		return
	}
	retuser := UserToUserNew(user)

	randomString, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf(" Error creating reftok body %v ", err))
		return
	}

	refTok, err := cfg.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     randomString,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		UserID:    user.ID,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf(" Error creating reftok in database %v ", err))
		return
	}
	retuser.Token = token
	retuser.RefreshToken = refTok.Token

	RespondWithJson(w, 200, retuser)
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
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error while confirming login: %v", err))
		return
	}

	// refToken, err := cfg.db.GetRefTokByID(ctx, token)
	// if err != nil {
	// 	RespondWithError(w, 400, fmt.Sprintf("Error while confirming login: %v", err))
	// 	return
	// }

	// if refToken.ExpiresAt.Before(time.Now()) {
	// 	RespondWithError(w, 401, fmt.Sprintf("Error while confirming login: %v", err))
	// 	return
	// }

	uid, err := auth.ValidateJWT(token, cfg.secret_token)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("Unauthorized Action %v", err))
		return
	}

	chirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uid,
		// UserID:    refToken.UserID,
		Body: chirpReq.Body,
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

func (cfg *apiConfig) ListChirpsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "Text/json")

	sort := r.URL.Query().Get("sort")

	var listChirps []database.Chirp
	var err error
	switch sort {
	case "asc":
		listChirps, err = cfg.db.ListChirps(context.Background())
	case "desc":
		listChirps, err = cfg.db.ListChirpsDesc(context.Background())
	default:
		listChirps, err = cfg.db.ListChirps(context.Background())
	}

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error getting chirps from database: %v\n", err))
		return
	}

	author := r.URL.Query().Get("author_id")

	if author != "" {
		author_id, err := uuid.Parse(author)
		if err != nil {
			RespondWithError(w, 400, fmt.Sprintf("Error getting chirps from database: %v\n", err))
			return
		}
		listChirps, err = cfg.db.GetAuthorChrips(context.Background(), author_id)
		if err != nil {
			RespondWithError(w, 400, fmt.Sprintf("Error getting chirps from database: %v\n", err))
			return
		}
		RespondWithJson(w, 200, listChirps)
		return
	}

	RespondWithJson(w, 200, listChirps)
}

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "Text/json")
	chirpId := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		RespondWithError(w, 400, "Wrong uuid")
		return
	}

	chirp, err := cfg.db.GetChipByID(context.Background(), chirpUUID)

	if err != nil {
		RespondWithError(w, 404, fmt.Sprintf("Error getting chirp from database: %v\n", err))
		return
	}

	RespondWithJson(w, 200, chirp)
}

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	w.Header().Set("Content-type", "Text/json")
	chirpId := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		RespondWithError(w, 400, "Wrong uuid")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("Error while confirming login: %v", err))
		return
	}

	uid, err := auth.ValidateJWT(token, cfg.secret_token)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("Unauthorized Action %v", err))
		return
	}
	chirp, err := cfg.db.GetChipByID(context.Background(), chirpUUID)

	if err != nil {
		RespondWithError(w, 404, fmt.Sprintf("Error getting chirp from database: %v\n", err))
		return
	}

	if chirp.UserID != uid {
		RespondWithError(w, 403, fmt.Sprintf("Unauthorized Action: %v", err))
		return
	}

	err = cfg.db.DeleteChirpByID(ctx, chirp.ID)
	if err != nil {
		RespondWithError(w, 404, fmt.Sprintf("Error deleting chirp from database: %v\n", err))
		return
	}

	w.WriteHeader(204)
}

func (cfg *apiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "Text/json")

	bearertok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf(" error getting bearer token %v", err))
		return
	}

	ctx := context.Background()

	reftok, err := cfg.db.GetRefTokByID(ctx, bearertok)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf(" the refresh token doesn't exist %v", err))
		return
	}

	accessToken, err := auth.MakeJWT(reftok.UserID, cfg.secret_token, time.Hour)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf(" error making access token %v", err))
		return
	}

	accessTokenStruct := Token{Token: accessToken}

	RespondWithJson(w, 200, accessTokenStruct)
}

func (cfg *apiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)

	bearertok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf(" error getting bearer token %v", err))
		return
	}

	ctx := context.Background()
	reftok, err := cfg.db.GetRefTokByID(ctx, bearertok)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf(" the refresh token doesn't exist %v", err))
		return
	}

	if reftok.ExpiresAt.Before(time.Now()) {
		RespondWithError(w, 401, " the refresh token is expired ")
		return
	}

	err = cfg.db.RevokeRefTok(ctx, database.RevokeRefTokParams{
		Token:     reftok.Token,
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Fatalf(" can't revoke the token %v ", err)
	}

}

func (cfg *apiConfig) UpgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	upgradeRequest, err := DecodeBody[UpgradeUserRequest](r)
	ctx := context.Background()
	if err != nil {
		w.WriteHeader(400)
		return
	}

	api_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if api_key != cfg.polka_key {
		w.WriteHeader(401)
		return
	}

	if upgradeRequest.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	err = cfg.db.UpgradeUser(ctx, upgradeRequest.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

/// admin methods here

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {

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
