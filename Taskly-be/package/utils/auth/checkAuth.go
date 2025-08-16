package auth

import (
	"encoding/json"
	"fmt"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"

	"github.com/gin-gonic/gin"
)

func CheckAuth(token string) error {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		return fmt.Errorf("don't extract bearer token")
	}
	_, err := VerifyTokenSubject(tokenString)
	if err != nil {
		return fmt.Errorf("Check auth failed", err)
	}
	return nil
}
func CheckAuthForWebsocket(token string) (*model.User, error) {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		return &model.User{}, fmt.Errorf("don't extract bearer token")
	}
	claims, err := VerifyTokenSubject(tokenString)
	if err != nil {
		return &model.User{}, fmt.Errorf("Check auth failed", err)
	}

	subjectStr := claims.Subject
	if subjectStr == "" {
		global.Logger.Sugar().Error("Subject is empty in JWT claims")
		return &model.User{}, err
	}

	var User model.User
	err = json.Unmarshal([]byte(subjectStr), &User)
	if err != nil {
		global.Logger.Sugar().Errorf("JSON unmarshal error: %v", err)
		return &model.User{}, err
	}

	return &User, nil
}

func GetUserFromContext(ctx *gin.Context) *model.UserToken {
	claimsValue := ctx.Request.Context().Value("claims")
	claims, ok := claimsValue.(*PayloadClaims)
	if !ok {
		global.Logger.Sugar().Error("Claims are not of type PayloadClaims")
		return &model.UserToken{}
	}

	return &model.UserToken{
		ID:       claims.UserID,
		UserType: claims.UserType,
	}
}
