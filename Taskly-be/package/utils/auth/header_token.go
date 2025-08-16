package auth

import (
	"strings"
)

func ExtractBearerToken(authHeader string) (string, bool) {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}
	return "", false
}
