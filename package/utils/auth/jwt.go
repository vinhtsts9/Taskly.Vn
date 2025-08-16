package auth

import (
	"fmt"
	"time"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type PayloadClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserType []string  `json:"user_type"`
	jwt.StandardClaims
}

func GenTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.JWT.API_SECRET_KEY))
}
func CreateToken(userToken model.UserToken) (string, error) {
	timeEx := global.Config.JWT.JWT_EXPIRATION
	if timeEx == "" {
		timeEx = "1h"
	}
	expiration, err := time.ParseDuration(timeEx)
	if err != nil {
		return "", err
	}
	now := time.Now()
	expiresAt := now.Add(expiration)

	return GenTokenJWT(&PayloadClaims{
		UserID:   userToken.ID,
		UserType: []string(userToken.UserType),
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "Vũ Thế Vinh",
		},
	})
}

func CreateRefreshToken(userToken string) (string, error) {
	refreshToken := uuid.New().String()
	return refreshToken, nil
}

func ParseJwtToken(token string) (*PayloadClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &PayloadClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.API_SECRET_KEY), nil
	})
	if claims, ok := parsedToken.Claims.(*PayloadClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, err
}

func VerifyTokenSubject(token string) (*PayloadClaims, error) {
	claims, err := ParseJwtToken(token)
	if err != nil {
		return &PayloadClaims{}, fmt.Errorf("Error verify token ", err)
	}
	if err = claims.Valid(); err != nil {
		return &PayloadClaims{}, fmt.Errorf("Error verify token ", err)

	}
	return claims, nil
}
func GenerateTokens(userToken model.UserToken) (accessToken, refreshToken string, err error) {
	accessToken, err = CreateToken(userToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %v", err)
	}

	refreshToken, err = CreateRefreshToken(userToken.ID.String())
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}
