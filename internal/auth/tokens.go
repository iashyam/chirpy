package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tockenSecret string, expiresIn time.Duration) (string, error) {
	secretKey := []byte(tockenSecret)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	var jwtToken jwt.RegisteredClaims
	var idzero uuid.UUID

	token, err := jwt.ParseWithClaims(tokenString, &jwtToken, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return idzero, fmt.Errorf(" error parsing tokenstring %v", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !(ok) {
		return idzero, fmt.Errorf(" error reading columns ")
	}

	if !token.Valid {
		return idzero, fmt.Errorf(" the token in not valid ")
	}

	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return idzero, fmt.Errorf(" i don't know what the fuck is going on %v", err)
	}
	return uid, nil

}
