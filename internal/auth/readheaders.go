package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func sanitizeToken(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 7 && strings.ToLower(s[:7]) == "bearer " {
		s = strings.TrimSpace(s[7:])
	}
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		s = s[1 : len(s)-1]
	}
	return s
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("no aothorization headers")
	}

	strippedBearer := sanitizeToken(bearer)
	return strippedBearer, nil
}
