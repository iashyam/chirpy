package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func sanitizeToken(s, first string) string {
	s = strings.TrimSpace(s)
	n := len(first)
	if len(s) >= n && strings.ToLower(s[:n]) == first {
		s = strings.TrimSpace(s[n:])
	}
	// if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
	// 	s = s[1 : len(s)-1]
	// }
	return s
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("no aothorization headers")
	}

	strippedBearer := sanitizeToken(bearer, "bearer ")
	return strippedBearer, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("no aothorization headers")
	}

	strippedBearer := sanitizeToken(bearer, "apikey ")
	return strippedBearer, nil
}
